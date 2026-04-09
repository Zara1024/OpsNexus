package controller

import (
	"dodevops-api/api/k8s/service"
	systemService "dodevops-api/api/system/service"
	"dodevops-api/common/result"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"dodevops-api/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type K8sTerminalController struct {
	service service.IK8sTerminalService
}

func NewK8sTerminalController(db *gorm.DB) *K8sTerminalController {
	return &K8sTerminalController{
		service: service.NewK8sTerminalService(db),
	}
}

// ConnectPodTerminal connects to a Pod terminal over WebSocket.
// @Summary 杩炴帴鍒癙od缁堢
// @Description 閫氳繃WebSocket杩炴帴鍒版寚瀹歅od鐨勭粓绔?
// @Tags K8s瀹瑰櫒缁堢
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 鐢ㄦ埛Token"
// @Param id path int true "闆嗙兢ID"
// @Param namespaceName path string true "鍛藉悕绌洪棿鍚嶇О"
// @Param podName path string true "Pod鍚嶇О"
// @Param containerName query string false "瀹瑰櫒鍚嶇О锛堥粯璁や负Pod涓涓€涓鍣級"
// @Param command query string false "鎵ц鍛戒护锛堥粯璁や负/bin/bash锛?
// @Success 101 "Switching Protocols"
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/terminal [get]
func (ctrl *K8sTerminalController) ConnectPodTerminal(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, http.StatusBadRequest, "鏃犳晥鐨勯泦缇D")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, http.StatusBadRequest, "鍛藉悕绌洪棿鍚嶇О涓嶈兘涓虹┖")
		return
	}

	podName := c.Param("podName")
	if podName == "" {
		result.Failed(c, http.StatusBadRequest, "Pod鍚嶇О涓嶈兘涓虹┖")
		return
	}

	containerName := c.Query("containerName")
	command := c.Query("command")
	if command == "" {
		command = "/bin/bash"
	}

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
		HostName:  fmt.Sprintf("cluster-%d/pod/%s/%s", clusterID, namespaceName, podName),
		HostIP:    "",
		SSHUser:   chooseNonEmpty(containerName, "pod-terminal"),
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
		log.Printf("failed to upgrade connection: %v", err)
		_ = recorder.Close(3, err.Error())
		return
	}
	defer conn.Close()

	stream, err := ctrl.service.CreateK8sWebSocketStream(uint(clusterID), namespaceName, podName, containerName, command, conn, recorder)
	if err != nil {
		log.Printf("failed to create K8s WebSocket stream: %v", err)
		_ = recorder.Close(3, err.Error())
		_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	defer stream.Close()

	<-stream.Ctx.Done()
	log.Printf("K8s terminal connection closed")
}

func chooseNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

// GetPodContainers returns the container list for the target Pod.
// @Summary 鑾峰彇Pod涓殑瀹瑰櫒鍒楄〃
// @Description 鑾峰彇鎸囧畾Pod涓墍鏈夊鍣ㄧ殑鍚嶇О鍒楄〃锛岀敤浜庣粓绔繛鎺ユ椂閫夋嫨瀹瑰櫒
// @Tags K8s瀹瑰櫒缁堢
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 鐢ㄦ埛Token"
// @Param id path int true "闆嗙兢ID"
// @Param namespaceName path string true "鍛藉悕绌洪棿鍚嶇О"
// @Param podName path string true "Pod鍚嶇О"
// @Success 200 {object} result.Result{data=[]string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/containers [get]
func (ctrl *K8sTerminalController) GetPodContainers(c *gin.Context) {
	clusterID, err := strconv.Atoi(c.Param("id"))
	if err != nil || clusterID <= 0 {
		result.Failed(c, http.StatusBadRequest, "鏃犳晥鐨勯泦缇D")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, http.StatusBadRequest, "鍛藉悕绌洪棿鍚嶇О涓嶈兘涓虹┖")
		return
	}

	podName := c.Param("podName")
	if podName == "" {
		result.Failed(c, http.StatusBadRequest, "Pod鍚嶇О涓嶈兘涓虹┖")
		return
	}

	containers, err := ctrl.service.GetPodContainers(uint(clusterID), namespaceName, podName)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			result.Failed(c, http.StatusNotFound, fmt.Sprintf("鑾峰彇瀹瑰櫒鍒楄〃澶辫触: %v", err))
			return
		}

		result.Failed(c, http.StatusInternalServerError, fmt.Sprintf("鑾峰彇瀹瑰櫒鍒楄〃澶辫触: %v", err))
		return
	}

	result.Success(c, containers)
}
