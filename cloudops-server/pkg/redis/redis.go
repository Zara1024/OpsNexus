package redis

import (
	"context"
	"fmt"

	red "github.com/redis/go-redis/v9"

	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/config"
)

// NewClient 初始化Redis客户端并进行连通性检查。
func NewClient(cfg config.RedisConfig) (*red.Client, error) {
	client := red.NewClient(&red.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("连接Redis失败: %w", err)
	}

	return client, nil
}
