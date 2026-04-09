package controller

import (
	"dodevops-api/api/cmdb/service"
	"dodevops-api/common/result"
	"dodevops-api/pkg/jwt"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type noEarlyResponseWriter struct {
	gin.ResponseWriter
	headerWritten bool
}

type hijackedResponseWriter struct {
	conn net.Conn
}

func (w *noEarlyResponseWriter) WriteHeader(code int) {
	if !w.headerWritten {
		w.headerWritten = true
		w.ResponseWriter.WriteHeader(code)
	}
}

func (w *noEarlyResponseWriter) WriteHeaderNow() {
	if !w.headerWritten {
		w.headerWritten = true
		w.ResponseWriter.WriteHeaderNow()
	}
}

func (w *hijackedResponseWriter) Header() http.Header {
	return http.Header{}
}

func (w *hijackedResponseWriter) Write(data []byte) (int, error) {
	return w.conn.Write(data)
}

func (w *hijackedResponseWriter) WriteHeader(statusCode int) {
	statusText := http.StatusText(statusCode)
	if statusText == "" {
		statusText = "status code " + strconv.Itoa(statusCode)
	}
	w.conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText)))
}

type CmdbHostSSHController struct {
	hostSSHService service.CmdbHostSSHServiceInterface
}

func NewCmdbHostSSHController(hostSSHService service.CmdbHostSSHServiceInterface) *CmdbHostSSHController {
	return &CmdbHostSSHController{hostSSHService: hostSSHService}
}

// ConnectTerminal 杩炴帴SSH缁堢
func (c *CmdbHostSSHController) ConnectTerminal(ctx *gin.Context) {
	hostID := ctx.Param("id")
	log.Printf("寮€濮嬪鐞哤ebSocket杩炴帴璇锋眰, hostID: %s", hostID)

	// 瑙ｆ瀽涓绘満ID
	id, err := strconv.ParseUint(hostID, 10, 32)
	if err != nil {
		log.Printf("鏃犳晥鐨勪富鏈篒D: %s, 閿欒: %v", hostID, err)
		ctx.String(http.StatusBadRequest, "鏃犳晥鐨勪富鏈篒D")
		return
	}

	// 浣跨敤鏍囧噯WebSocket鍗囩骇
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// 浼樺厛浠嶶RL鍙傛暟鑾峰彇token
	token := ctx.Query("token")
	if token == "" {
		authHeader := ctx.GetHeader("Authorization")
		// 濡傛灉URL涓病鏈塼oken锛屾鏌uthorization澶?		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("missing token or authorization header")
			ctx.String(http.StatusUnauthorized, "missing token or authorization header")
			return
		}

		// 楠岃瘉Bearer浠ょ墝鏍煎紡
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			log.Println("invalid authorization header format")
			ctx.String(http.StatusUnauthorized, "invalid authorization header format")
			return
		}
		token = authHeader[7:]
	}
	_, err = jwt.ValidateToken(token)
	if err != nil {
		log.Printf("浠ょ墝楠岃瘉澶辫触: %v", err)
		ctx.String(http.StatusUnauthorized, "浠ょ墝楠岃瘉澶辫触")
		return
	}

	// 鍗囩骇WebSocket杩炴帴
	wsConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("WebSocket鍗囩骇澶辫触: %v", err)
		ctx.String(http.StatusInternalServerError, "WebSocket鍗囩骇澶辫触")
		return
	}

	// 杩炴帴SSH缁堢
	log.Println("灏濊瘯杩炴帴SSH缁堢")
	webSSH, err := c.hostSSHService.ConnectTerminal(ctx, uint(id))
	if err != nil {
		log.Printf("SSH杩炴帴澶辫触: %v", err)
		wsConn.Close()
		ctx.String(http.StatusInternalServerError, "SSH杩炴帴澶辫触")
		return
	}

	// 杩炴帴WebSocket
	if err := webSSH.Connect(wsConn); err != nil {
		log.Printf("WebSSH杩炴帴澶辫触: %v", err)
		webSSH.Close()
		wsConn.Close()
		ctx.String(http.StatusInternalServerError, "WebSSH杩炴帴澶辫触")
		return
	}

	// SSH杩炴帴娴嬭瘯閫氳繃
	log.Println("SSH杩炴帴娴嬭瘯閫氳繃")
	log.Printf("SSH杩炴帴娴嬭瘯閫氳繃, WebSocket鍜孲SH杩炴帴宸插缓绔? 杩滅▼鍦板潃: %s", wsConn.RemoteAddr().String())

	// 璁剧疆defer纭繚璧勬簮閲婃斁
	defer func() {
		webSSH.Close()
		wsConn.Close()
	}()

	// 淇濇寔杩炴帴
	select {}
}

// ExecuteCommand 鎵ц鍛戒护
// @Summary 鎵цSSH鍛戒护
// @Description 鍦⊿SH缁堢鎵ц鍛戒护
// @Tags CMDB涓绘満SSH
// @Accept  json
// @Produce  json
// @Param id path uint true "涓绘満ID"
// @Param command query string true "鍛戒护"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostssh/command/{id} [get]
// @Security ApiKeyAuth
func (c *CmdbHostSSHController) PreviewCommandRisk(ctx *gin.Context) {
	hostID := ctx.Param("id")
	command := strings.TrimSpace(ctx.Query("command"))

	id, err := strconv.ParseUint(hostID, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的主机ID")
		return
	}

	assessment, err := c.hostSSHService.PreviewCommandRisk(ctx, uint(id), command)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "命令风险预检查失败: "+err.Error())
		return
	}

	result.Success(ctx, assessment)
}

func (c *CmdbHostSSHController) ExecuteCommand(ctx *gin.Context) {
	hostID := ctx.Param("id")
	command := strings.TrimSpace(ctx.Query("command"))
	if command == "" {
		result.Failed(ctx, http.StatusBadRequest, "命令不能为空")
		return
	}

	riskAck := false
	switch strings.ToLower(strings.TrimSpace(ctx.Query("riskAck"))) {
	case "1", "true", "yes":
		riskAck = true
	}
	confirmedRiskLevel, _ := strconv.ParseInt(ctx.Query("confirmedRiskLevel"), 10, 64)

	id, err := strconv.ParseUint(hostID, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "无效的主机ID")
		return
	}

	assessment, err := c.hostSSHService.PreviewCommandRisk(ctx, uint(id), command)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "命令风险预检查失败: "+err.Error())
		return
	}
	if assessment.RequiresConfirmation && !riskAck {
		ctx.JSON(http.StatusOK, result.Result{
			Code:    http.StatusConflict,
			Message: "命令风险需要确认后才能执行",
			Data:    assessment,
		})
		return
	}
	if riskAck && confirmedRiskLevel > 0 && confirmedRiskLevel != assessment.RiskLevel {
		ctx.JSON(http.StatusOK, result.Result{
			Code:    http.StatusConflict,
			Message: "命令风险级别已变化，请重新确认后再执行",
			Data:    assessment,
		})
		return
	}

	output, err := c.hostSSHService.ExecuteCommand(ctx, uint(id), command)
	if err != nil {
		result.Failed(ctx, http.StatusInternalServerError, "执行命令失败: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"hostId":  hostID,
		"command": command,
		"output":  output,
		"risk":    assessment,
	})
}

// UploadFile 涓婁紶鏂囦欢鍒癝SH鏈嶅姟鍣?// @Summary 涓婁紶鏂囦欢鍒癝SH鏈嶅姟鍣?// @Description 涓婁紶鏈湴鏂囦欢鍒拌繙绋婼SH鏈嶅姟鍣?// @Tags CMDB涓绘満SSH
// @Accept multipart/form-data
// @Produce json
// @Param id path uint true "涓绘満ID"
// @Param file formData file true "瑕佷笂浼犵殑鏂囦欢"
// @Param destPath formData string true "?????????"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostssh/upload/{id} [post]
// @Security ApiKeyAuth
func (c *CmdbHostSSHController) UploadFile(ctx *gin.Context) {
	hostID := ctx.Param("id")
	file, err := ctx.FormFile("file")
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "鑾峰彇涓婁紶鏂囦欢澶辫触: "+err.Error())
		return
	}

	destPath := ctx.PostForm("destPath")
	if destPath == "" {
		result.Failed(ctx, http.StatusBadRequest, "鐩爣璺緞涓嶈兘涓虹┖")
		return
	}

	tempDir := "/tmp/ssh_uploads"
	log.Printf("鍒涘缓涓存椂鐩綍: %s", tempDir)
	if err := os.MkdirAll(tempDir, 0o755); err != nil {
		log.Printf("鍒涘缓涓存椂鐩綍澶辫触: %v", err)
		result.Failed(ctx, http.StatusInternalServerError, "鍒涘缓涓存椂鐩綍澶辫触: "+err.Error())
		return
	}

	tempFilePath := filepath.Join(tempDir, file.Filename)
	log.Printf("灏濊瘯淇濆瓨鏂囦欢鍒? %s", tempFilePath)
	if _, err := os.Stat(tempFilePath); err == nil {
		log.Printf("鏂囦欢宸插瓨鍦紝灏嗚瑕嗙洊: %s", tempFilePath)
		if err := os.Remove(tempFilePath); err != nil {
			log.Printf("鍒犻櫎宸插瓨鍦ㄦ枃浠跺け璐? %v", err)
			result.Failed(ctx, http.StatusInternalServerError, "鍒犻櫎宸插瓨鍦ㄦ枃浠跺け璐? "+err.Error())
			return
		}
	}

	log.Printf("淇濆瓨鏂囦欢: %s (澶у皬: %d bytes)", file.Filename, file.Size)
	if err := ctx.SaveUploadedFile(file, tempFilePath); err != nil {
		log.Printf("淇濆瓨鏂囦欢澶辫触: %v", err)
		result.Failed(ctx, http.StatusInternalServerError, "淇濆瓨涓存椂鏂囦欢澶辫触: "+err.Error())
		return
	}

	if err := os.Chmod(tempFilePath, 0o644); err != nil {
		log.Printf("璁剧疆鏂囦欢鏉冮檺澶辫触: %v", err)
		result.Failed(ctx, http.StatusInternalServerError, "璁剧疆鏂囦欢鏉冮檺澶辫触: "+err.Error())
		return
	}

	if tempFile, err := os.Open(tempFilePath); err == nil {
		tempFile.Sync()
		tempFile.Close()
	} else {
		log.Printf("鏃犳硶鎵撳紑鏂囦欢杩涜鍚屾: %v", err)
		result.Failed(ctx, http.StatusInternalServerError, "鏃犳硶楠岃瘉鏂囦欢鐘舵€? "+err.Error())
		return
	}

	if _, err := os.Stat(tempFilePath); err != nil {
		log.Printf("鏂囦欢淇濆瓨楠岃瘉澶辫触: %v", err)
		result.Failed(ctx, http.StatusInternalServerError, "鏂囦欢淇濆瓨楠岃瘉澶辫触: "+err.Error())
		return
	}

	id, err := strconv.ParseUint(hostID, 10, 32)
	if err != nil {
		result.Failed(ctx, http.StatusBadRequest, "鏃犳晥鐨勪富鏈篒D")
		return
	}

	log.Printf("寮€濮嬩笂浼犳枃浠跺埌杩滅▼鏈嶅姟鍣? %s -> %s", tempFilePath, destPath)
	err = c.hostSSHService.UploadFile(ctx, uint(id), tempFilePath, destPath)
	if err != nil {
		log.Printf("鏂囦欢涓婁紶澶辫触: %v", err)
		result.Failed(ctx, http.StatusInternalServerError, "鏂囦欢涓婁紶澶辫触: "+err.Error())
		return
	}

	if err := os.Remove(tempFilePath); err != nil {
		log.Printf("鍒犻櫎涓存椂鏂囦欢澶辫触: %v (鏂囦欢鍙兘宸茶鍏朵粬杩涚▼鍒犻櫎)", err)
	} else {
		log.Printf("涓存椂鏂囦欢宸插垹闄? %s", tempFilePath)
	}

	result.Success(ctx, gin.H{
		"hostId":   hostID,
		"fileName": file.Filename,
		"destPath": destPath,
		"message":  "鏂囦欢涓婁紶鎴愬姛",
	})
}
