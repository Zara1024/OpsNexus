package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"dodevops-api/api/k8s/model"
	"dodevops-api/api/k8s/service"
	systemService "dodevops-api/api/system/service"
	"dodevops-api/common/result"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type K8sKubectlController struct {
	service *service.K8sKubectlService
}

func NewK8sKubectlController(db *gorm.DB) *K8sKubectlController {
	return &K8sKubectlController{
		service: service.NewK8sKubectlService(db),
	}
}

func (ctrl *K8sKubectlController) ExecuteKubectl(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, http.StatusBadRequest, "鏃犳晥鐨勯泦缇D")
		return
	}

	var req model.KubectlCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "璇锋眰鍙傛暟閿欒: "+err.Error())
		return
	}

	response, err := ctrl.service.ExecuteKubectlCommand(uint(clusterID), &req)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("kubectl 鎵ц澶辫触: %v", err))
		return
	}

	result.Success(c, response)
}

func (ctrl *K8sKubectlController) ConnectKubectlTerminal(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, http.StatusBadRequest, "鏃犳晥鐨勯泦缇D")
		return
	}

	namespace := c.Query("namespace")
	admin, _ := jwt.GetAdmin(c)
	recorder, err := systemService.NewTerminalAuditRecorder(systemService.TerminalAuditRecorderOptions{
		AdminID: func() uint {
			if admin != nil {
				return admin.ID
			}
			return 0
		}(),
		Username: func() string {
			if admin != nil {
				return admin.Username
			}
			return "unknown"
		}(),
		HostID:    0,
		HostName:  fmt.Sprintf("cluster-%d/kubectl", clusterID),
		HostIP:    "",
		SSHUser:   "kubectl",
		ClientIP:  c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Width:     120,
		Height:    32,
	})
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "创建终端审计会话失败: "+err.Error())
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("failed to upgrade kubectl terminal connection: %v", err)
		_ = recorder.Close(3, err.Error())
		return
	}
	defer conn.Close()

	session, err := ctrl.service.CreateKubectlTerminalSession(uint(clusterID), namespace, conn, recorder)
	if err != nil {
		log.Printf("failed to create kubectl terminal session: %v", err)
		_ = recorder.Close(3, err.Error())
		_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	defer session.Close()

	<-session.Ctx.Done()
}
