package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/rs/zerolog"

	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/config"
)

// Init 初始化slog与zerolog，统一输出结构化日志。
func Init(cfg config.LogConfig) *slog.Logger {
	zerolog.SetGlobalLevel(parseZeroLevel(cfg.Level))
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: parseSlogLevel(cfg.Level)}))
}

// parseSlogLevel 解析slog日志等级。
func parseSlogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// parseZeroLevel 解析zerolog日志等级。
func parseZeroLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
