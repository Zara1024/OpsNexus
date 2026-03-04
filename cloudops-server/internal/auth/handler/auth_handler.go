package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/service"
	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/middleware"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// Login 处理登录请求。
func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
		CaptchaID string `json:"captcha_id" binding:"required"`
		Captcha   string `json:"captcha" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	res, err := h.services.Auth.Login(c.Request.Context(), &service.LoginRequest{
		Username:  req.Username,
		Password:  req.Password,
		CaptchaID: req.CaptchaID,
		Captcha:   req.Captcha,
	})
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, res)
}

// GetCaptcha 获取验证码。
func (h *Handler) GetCaptcha(c *gin.Context) {
	res, err := h.services.Auth.GetCaptcha(c.Request.Context())
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, res)
}

// Refresh 刷新Token。
func (h *Handler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.ErrValidation.Code, apperrors.ErrValidation.Message)
		return
	}
	pair, err := h.services.Auth.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{"token": pair})
}

// Logout 处理登出请求。
func (h *Handler) Logout(c *gin.Context) {
	claims, ok := middleware.GetJWTClaims(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Code, apperrors.ErrUnauthorized.Message)
		return
	}
	auth := c.GetHeader("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer"))
	if err := h.services.Auth.Logout(c.Request.Context(), strings.TrimSpace(token), claims); err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

// UserInfo 获取当前用户信息。
func (h *Handler) UserInfo(c *gin.Context) {
	claims, ok := middleware.GetJWTClaims(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, apperrors.ErrUnauthorized.Code, apperrors.ErrUnauthorized.Message)
		return
	}
	data, err := h.services.Auth.GetUserInfo(claims.UserID)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, data)
}
