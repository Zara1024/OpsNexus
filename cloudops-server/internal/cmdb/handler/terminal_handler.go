package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/middleware"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// ListSSHRecords 查询 SSH 审计记录。
func (h *Handler) ListSSHRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	res, err := h.services.Terminal.ListSSHRecords(page, pageSize)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.Success(c, res)
}

// WebTerminal 处理 WebSocket 终端请求。
func (h *Handler) WebTerminal(c *gin.Context) {
	hostID, err := strconv.ParseInt(c.Param("host_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 20010, "主机ID格式错误")
		return
	}
	claims, ok := middleware.GetJWTClaims(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 10001, "未登录或登录已过期")
		return
	}
	if err = h.services.Terminal.StartWebTerminal(c, claims.UserID, hostID); err != nil {
		handleAppError(c, err)
		return
	}
}
