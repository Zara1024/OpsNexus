package router

import (
	"path/filepath"
	"reflect"
	"strings"

	"dodevops-api/api/system/controller"
	"dodevops-api/common/config"
	"dodevops-api/middleware"
	"dodevops-api/pkg/log"
	appRouter "dodevops-api/router/app"
	aiRouter "dodevops-api/router/ai"
	cmdbRouter "dodevops-api/router/cmdb"
	configCenterRouter "dodevops-api/router/configCenter"
	dashboardRouter "dodevops-api/router/dashboard"
	k8sRouter "dodevops-api/router/k8s"
	knowledgeRouter "dodevops-api/router/knowledge"
	monitorRouter "dodevops-api/router/monitor"
	searchRouter "dodevops-api/router/search"
	systemRouter "dodevops-api/router/system"
	taskRouter "dodevops-api/router/task"
	toolRouter "dodevops-api/router/tool"

	agentController "dodevops-api/api/monitor/controller"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter initializes the top-level Gin router.
func InitRouter() *gin.Engine {
	router := gin.New()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			tag := fld.Tag.Get("json")
			if tag == "-" {
				return ""
			}
			return strings.Split(tag, ",")[0]
		})
	}

	router.Use(gin.Recovery())
	router.Use(middleware.Cors())

	uploadDir := config.Config.ImageSettings.UploadDir
	if !strings.HasPrefix(uploadDir, "/") {
		if absPath, err := filepath.Abs(uploadDir); err == nil {
			uploadDir = absPath
			log.Log().Infof("Static upload directory: %s -> %s", config.Config.ImageSettings.UploadDir, uploadDir)
		}
	} else {
		log.Log().Infof("Static upload directory: %s", uploadDir)
	}
	router.Static("/api/v1/upload", uploadDir)
	router.Use(log.CustomGinLogger())

	register(router)
	registerWebSocketRoutes(router)

	return router
}

func register(router *gin.Engine) {
	if config.Config.Server.EnableSwagger {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		log.Log().Info("Swagger API docs enabled")
	} else {
		log.Log().Info("Swagger API docs disabled")
	}

	agentCtrl := agentController.NewAgentController()
	alertCtrl := agentController.NewMonitorAlertController()

	apiGroup := router.Group("/api/v1")
	{
		apiGroup.GET("/captcha", controller.Captcha)
		apiGroup.POST("/login", controller.Login)
		apiGroup.POST("/monitor/agent/heartbeat", agentCtrl.UpdateHeartbeat)
		apiGroup.POST("/monitor/alerts/webhook", alertCtrl.ReceiveWebhook)

		jwtGroup := apiGroup.Group("")
		jwtGroup.Use(middleware.AuthMiddleware())
		jwtGroup.Use(middleware.LogMiddleware())
		{
			systemRouter.RegisterSystemRoutes(jwtGroup)
			cmdbRouter.RegisterCmdbRoutes(jwtGroup)
			configCenterRouter.RegisterConfigCenterRoutes(jwtGroup)
			aiRouter.RegisterAIRoutes(jwtGroup)
			appRouter.RegisterJenkinsRoutes(jwtGroup)
			appRouter.RegisterApplicationRoutes(jwtGroup)
			dashboardRouter.RegisterDashboardRoutes(jwtGroup)
			k8sRouter.RegisterK8sRoutes(jwtGroup)
			knowledgeRouter.RegisterKnowledgeRoutes(jwtGroup)
			monitorRouter.InitMonitorRouter(jwtGroup)
			searchRouter.RegisterSearchRoutes(jwtGroup)
			taskRouter.RegisterTaskRoutes(jwtGroup)
			toolRouter.RegisterToolRoutes(jwtGroup)
		}
	}
}

func registerWebSocketRoutes(router *gin.Engine) {
	taskRouter.RegisterWebSocketRoutes(router.Group(""))
}
