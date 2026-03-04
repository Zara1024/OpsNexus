package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	authhandler "github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/handler"
	authservice "github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/service"
	cmdbhandler "github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/handler"
	cmdbservice "github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/service"
	k8shandler "github.com/Zara1024/OpsNexus/cloudops-server/internal/k8s/handler"
	k8sservice "github.com/Zara1024/OpsNexus/cloudops-server/internal/k8s/service"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/config"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/crypto"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/database"
	jwtx "github.com/Zara1024/OpsNexus/cloudops-server/pkg/jwt"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/logger"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/middleware"
	redisx "github.com/Zara1024/OpsNexus/cloudops-server/pkg/redis"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("加载配置失败: %w", err))
	}

	log := logger.Init(cfg.Log)
	slog.SetDefault(log)

	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Error("初始化数据库失败", "error", err)
		os.Exit(1)
	}

	redisClient, err := redisx.NewClient(cfg.Redis)
	if err != nil {
		log.Error("初始化Redis失败", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	jwtManager := jwtx.NewManager(cfg.JWT.Secret, cfg.JWT.Issuer, 30*time.Minute, 7*24*time.Hour)
	authServices := authservice.New(db, redisClient, jwtManager)
	authHandlers := authhandler.New(authServices)

	aes := crypto.NewAESGCM(cfg.JWT.Secret)
	cmdbServices := cmdbservice.New(db, aes)
	cmdbHandlers := cmdbhandler.New(cmdbServices)
	k8sServices := k8sservice.New(db, aes)
	k8sHandlers := k8shandler.New(k8sServices)

	jwtMiddleware := middleware.JWTAuth(jwtManager, authServices.Auth)
	permissionMiddleware := middleware.Permission
	auditMiddleware := middleware.Audit(log)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(log))
	r.Use(middleware.Recovery(log))
	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	v1 := r.Group("/api/v1")
	authHandlers.RegisterRoutes(v1, jwtMiddleware, permissionMiddleware, auditMiddleware)
	cmdbHandlers.RegisterRoutes(v1, jwtMiddleware, permissionMiddleware, auditMiddleware)
	cmdbHandlers.RegisterWSRoutes(r, jwtMiddleware)
	k8sHandlers.RegisterRoutes(v1, jwtMiddleware, permissionMiddleware, auditMiddleware)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		log.Info("服务启动成功", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("服务启动失败", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("收到退出信号，开始优雅停机")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("服务优雅停机失败", "error", err)
		os.Exit(1)
	}

	log.Info("服务已停止")
}
