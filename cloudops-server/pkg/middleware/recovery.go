package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// Recovery 捕获panic并返回统一错误。
func Recovery(log *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Error("请求发生panic", "panic", recovered)
		response.Error(c, http.StatusInternalServerError, apperrors.ErrInternal.Code, apperrors.ErrInternal.Message)
	})
}
