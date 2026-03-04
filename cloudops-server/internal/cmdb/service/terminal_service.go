package service

import (
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/model"
	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/repository"
)

// TerminalService 提供 Web 终端与审计能力。
type TerminalService struct {
	hostRepo       *repository.HostRepository
	credentialRepo *repository.CredentialRepository
	sshRecordRepo  *repository.SSHRecordRepository
	credentialSvc  *CredentialService
}

// SSHRecordListResponse 定义 SSH 审计分页结果。
type SSHRecordListResponse struct {
	List  []model.SSHRecord `json:"list"`
	Total int64             `json:"total"`
}

// NewTerminalService 创建终端服务。
func NewTerminalService(hostRepo *repository.HostRepository, credentialRepo *repository.CredentialRepository, sshRecordRepo *repository.SSHRecordRepository, credentialSvc *CredentialService) *TerminalService {
	return &TerminalService{hostRepo: hostRepo, credentialRepo: credentialRepo, sshRecordRepo: sshRecordRepo, credentialSvc: credentialSvc}
}

// ListSSHRecords 查询 SSH 审计记录。
func (s *TerminalService) ListSSHRecords(page, pageSize int) (*SSHRecordListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	list, total, err := s.sshRecordRepo.List(page, pageSize)
	if err != nil {
		return nil, err
	}
	return &SSHRecordListResponse{List: list, Total: total}, nil
}

// StartWebTerminal 开启 WebSocket 终端会话。
func (s *TerminalService) StartWebTerminal(c *gin.Context, userID, hostID int64) error {
	host, err := s.hostRepo.GetByID(hostID)
	if err != nil {
		return err
	}
	credential, err := s.credentialRepo.GetByID(host.CredentialID)
	if err != nil {
		return err
	}
	username, password, privateKey, passphrase, err := s.credentialSvc.ResolvePlainCredential(credential.ID)
	if err != nil {
		return err
	}

	authMethods := make([]ssh.AuthMethod, 0)
	if credential.Type == "key" && strings.TrimSpace(privateKey) != "" {
		var signer ssh.Signer
		if strings.TrimSpace(passphrase) != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(privateKey), []byte(passphrase))
		} else {
			signer, err = ssh.ParsePrivateKey([]byte(privateKey))
		}
		if err != nil {
			return err
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else {
		authMethods = append(authMethods, ssh.Password(password))
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(host.InnerIP, "22"), &ssh.ClientConfig{
		User:            username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	})
	if err != nil {
		return err
	}
	defer client.Close()

	upgrader := websocket.Upgrader{CheckOrigin: func(_ *http.Request) bool { return true }}
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return err
	}
	defer wsConn.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	session.Stderr = session.Stdout

	if err = session.RequestPty("xterm-256color", 40, 160, ssh.TerminalModes{ssh.ECHO: 1}); err != nil {
		return err
	}
	if err = session.Shell(); err != nil {
		return err
	}

	sessionID := uuid.NewString()
	clientIP := c.ClientIP()
	_ = s.sshRecordRepo.Create(&model.SSHRecord{UserID: userID, HostID: hostID, SessionID: sessionID, StartTime: time.Now(), ClientIP: clientIP})
	defer s.sshRecordRepo.EndSession(sessionID)

	readDone := make(chan struct{})
	writeDone := make(chan struct{})

	go func() {
		defer close(readDone)
		buf := make([]byte, 4096)
		for {
			n, readErr := stdout.Read(buf)
			if n > 0 {
				_ = wsConn.WriteMessage(websocket.TextMessage, buf[:n])
			}
			if readErr != nil {
				if readErr != io.EOF {
					_ = wsConn.WriteMessage(websocket.TextMessage, []byte("\r\n连接已断开\r\n"))
				}
				return
			}
		}
	}()

	go func() {
		defer close(writeDone)
		for {
			_, msg, readErr := wsConn.ReadMessage()
			if readErr != nil {
				return
			}
			if _, err = stdin.Write(msg); err != nil {
				return
			}
		}
	}()

	select {
	case <-readDone:
	case <-writeDone:
	case <-time.After(30 * time.Minute):
		_, _ = stdin.Write([]byte("exit\n"))
	}
	return nil
}
