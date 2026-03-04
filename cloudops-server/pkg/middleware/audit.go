package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Audit 记录增删改审计日志。
func Audit(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		claims, _ := GetJWTClaims(c)
		userID := int64(0)
		username := "anonymous"
		if claims != nil {
			userID = claims.UserID
			username = claims.Username
		}

		log.Info("审计日志",
			"request_id", GetRequestID(c),
			"user_id", userID,
			"username", username,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start).String(),
		)
	}
}
