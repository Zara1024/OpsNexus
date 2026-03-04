package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/service"
	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	jwtx "github.com/Zara1024/OpsNexus/cloudops-server/pkg/jwt"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

const (
	jwtClaimsKey      = "jwt_claims"
	headerNewAccess   = "X-New-Access-Token"
	headerNewRefresh  = "X-New-Refresh-Token"
	headerExpose      = "Access-Control-Expose-Headers"
	exposeHeaderValue = "X-New-Access-Token, X-New-Refresh-Token"
)

// JWTAuth 执行JWT鉴权、黑名单校验与自动续期。
func JWTAuth(jwtMgr *jwtx.Manager, authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c.GetHeader("Authorization"))
		if token == "" {
			response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Code, apperrors.ErrUnauthorized.Message)
			c.Abort()
			return
		}

		claims, err := jwtMgr.ParseAccessToken(token)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Code, apperrors.ErrUnauthorized.Message)
			c.Abort()
			return
		}

		if blacklisted, err := authSvc.IsBlacklisted(c.Request.Context(), token); err != nil {
			response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Code, apperrors.ErrUnauthorized.Message)
			c.Abort()
			return
		} else if blacklisted {
			response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Code, apperrors.ErrUnauthorized.Message)
			c.Abort()
			return
		}

		if jwtMgr.NeedAutoRefresh(claims) {
			pair, err := jwtMgr.GenerateTokenPair(claims.UserID, claims.Username, claims.RoleCodes, claims.Permissions)
			if err == nil {
				c.Header(headerNewAccess, pair.AccessToken)
				c.Header(headerNewRefresh, pair.RefreshToken)
				c.Header(headerExpose, exposeHeaderValue)
			}
		}

		c.Set(jwtClaimsKey, claims)
		c.Next()
	}
}

// GetJWTClaims 从上下文读取JWT声明。
func GetJWTClaims(c *gin.Context) (*jwtx.CustomClaims, bool) {
	v, ok := c.Get(jwtClaimsKey)
	if !ok {
		return nil, false
	}
	claims, ok := v.(*jwtx.CustomClaims)
	return claims, ok
}

func extractBearerToken(auth string) string {
	auth = strings.TrimSpace(auth)
	if auth == "" {
		return ""
	}
	if !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return ""
	}
	return strings.TrimSpace(auth[7:])
}
