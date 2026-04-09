package middleware

import (
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"dodevops-api/pkg/jwt"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		if strings.EqualFold(c.GetHeader("Upgrade"), "websocket") {
			token := c.Query("token")
			if token == "" {
				authHeader := c.GetHeader("Authorization")
				if strings.HasPrefix(authHeader, "Bearer ") {
					token = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}

			if token == "" {
				abortWebSocketAuth(c, http.StatusUnauthorized, int(result.ApiCode.NOAUTH), result.ApiCode.GetMessage(result.ApiCode.NOAUTH))
				return
			}

			mc, err := jwt.ValidateToken(token)
			if err != nil {
				abortWebSocketAuth(c, http.StatusUnauthorized, int(result.ApiCode.INVALIDTOKEN), result.ApiCode.GetMessage(result.ApiCode.INVALIDTOKEN))
				return
			}

			c.Set(constant.ContextKeyUserObj, mc)
			c.Next()
			return
		}

		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			token := c.Query("token")
			if token != "" {
				mc, err := jwt.ValidateToken(token)
				if err == nil {
					c.Set(constant.ContextKeyUserObj, mc)
					c.Next()
					return
				}

				fmt.Printf("[DEBUG] Token validation failed: %v, token length: %d, secret: %s\n", err, len(token), string(jwt.Secret))
			}

			result.Failed(c, int(result.ApiCode.NOAUTH), result.ApiCode.GetMessage(result.ApiCode.NOAUTH))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			result.Failed(c, int(result.ApiCode.AUTHFORMATERROR), result.ApiCode.GetMessage(result.ApiCode.AUTHFORMATERROR))
			c.Abort()
			return
		}

		mc, err := jwt.ValidateToken(parts[1])
		if err != nil {
			result.Failed(c, int(result.ApiCode.INVALIDTOKEN), result.ApiCode.GetMessage(result.ApiCode.INVALIDTOKEN))
			c.Abort()
			return
		}

		c.Set(constant.ContextKeyUserObj, mc)
		c.Next()
	}
}

func abortWebSocketAuth(c *gin.Context, status, code int, message string) {
	c.AbortWithStatusJSON(status, result.Result{
		Code:    code,
		Message: message,
		Data:    gin.H{},
	})
}
