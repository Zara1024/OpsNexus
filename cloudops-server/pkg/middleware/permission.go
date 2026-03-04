package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// Permission 校验接口权限标识。
func Permission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := GetJWTClaims(c)
		if !ok {
			response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Code, apperrors.ErrUnauthorized.Message)
			c.Abort()
			return
		}
		if permission == "" || hasPermission(claims.Permissions, permission) {
			c.Next()
			return
		}
		response.Error(c, http.StatusForbidden, apperrors.ErrForbidden.Code, apperrors.ErrForbidden.Message)
		c.Abort()
	}
}

func hasPermission(permissionList []string, need string) bool {
	for _, p := range permissionList {
		if p == "*:*:*" || p == "system:*:*" || p == "cmdb:*:*" || p == "k8s:*:*" || p == need {
			return true
		}
	}
	return false
}
