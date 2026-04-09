package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"dodevops-api/api/k8s/dao"
	"dodevops-api/api/k8s/model"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultKubectlTimeoutSeconds = 30
	maxKubectlTimeoutSeconds     = 300
)

type K8sKubectlService struct {
	clusterDao *dao.KubeClusterDao
}

type KubectlTerminalSession struct {
	sync.RWMutex
	Conn      *websocket.Conn
	Ctx       context.Context
	cancel    context.CancelFunc
	cmd       *exec.Cmd
	ptyFile   *os.File
	cleanup   func()
	closed    bool
	namespace string
	clusterID uint
	recorder  TerminalAuditRecorder
}

func NewK8sKubectlService(db *gorm.DB) *K8sKubectlService {
	return &K8sKubectlService{
		clusterDao: dao.NewKubeClusterDao(db),
	}
}

func (s *K8sKubectlService) ExecuteKubectlCommand(clusterID uint, req *model.KubectlCommandRequest) (*model.KubectlCommandResponse, error) {
	kubeconfigPath, cleanup, err := s.prepareKubectlEnvironment(clusterID, req.Namespace)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	timeoutSeconds := req.TimeoutSeconds
	if timeoutSeconds <= 0 {
		timeoutSeconds = defaultKubectlTimeoutSeconds
	}
	if timeoutSeconds > maxKubectlTimeoutSeconds {
		timeoutSeconds = maxKubectlTimeoutSeconds
	}

	commandText := normalizeKubectlCommand(req.Command, req.Namespace)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-lc", commandText)
	cmd.Env = append(os.Environ(),
		"KUBECONFIG="+kubeconfigPath,
		"TERM=xterm-256color",
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	startTime := time.Now()
	runErr := cmd.Run()
	duration := time.Since(startTime).Milliseconds()

	exitCode := 0
	success := runErr == nil
	errorMessage := ""
	if runErr != nil {
		success = false
		errorMessage = runErr.Error()
		if exitErr, ok := runErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return &model.KubectlCommandResponse{
		Success:      success,
		Command:      commandText,
		Namespace:    strings.TrimSpace(req.Namespace),
		Stdout:       stdout.String(),
		Stderr:       stderr.String(),
		ExitCode:     exitCode,
		DurationMs:   duration,
		ErrorMessage: errorMessage,
	}, nil
}

func (s *K8sKubectlService) CreateKubectlTerminalSession(clusterID uint, namespace string, conn *websocket.Conn, recorder TerminalAuditRecorder) (*KubectlTerminalSession, error) {
	kubeconfigPath, cleanup, err := s.prepareKubectlEnvironment(clusterID, namespace)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "sh", "-i")
	cmd.Env = append(os.Environ(),
		"KUBECONFIG="+kubeconfigPath,
		"TERM=xterm-256color",
	)

	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{Cols: 120, Rows: 32})
	if err != nil {
		cancel()
		cleanup()
		return nil, fmt.Errorf("create kubectl terminal failed: %w", err)
	}

	session := &KubectlTerminalSession{
		Conn:      conn,
		Ctx:       ctx,
		cancel:    cancel,
		cmd:       cmd,
		ptyFile:   ptmx,
		cleanup:   cleanup,
		namespace: strings.TrimSpace(namespace),
		clusterID: clusterID,
		recorder:  recorder,
	}

	if err := session.writeBanner(); err != nil {
		session.Close()
		return nil, err
	}

	go session.pipeOutput()
	go session.readWebSocketInput()
	go session.waitProcessExit()

	return session, nil
}

func (s *K8sKubectlService) prepareKubectlEnvironment(clusterID uint, namespace string) (string, func(), error) {
	cluster, err := s.clusterDao.GetByID(clusterID)
	if err != nil {
		return "", nil, fmt.Errorf("get cluster failed: %w", err)
	}
	if strings.TrimSpace(cluster.Credential) == "" {
		return "", nil, fmt.Errorf("cluster kubeconfig is empty")
	}

	kubeconfigContent := []byte(cluster.Credential)
	if strings.TrimSpace(namespace) != "" {
		if configObj, loadErr := clientcmd.Load(kubeconfigContent); loadErr == nil {
			if currentContext := configObj.CurrentContext; currentContext != "" {
				if ctxConfig, ok := configObj.Contexts[currentContext]; ok {
					ctxConfig.Namespace = namespace
				}
			}
			if encoded, writeErr := clientcmd.Write(*configObj); writeErr == nil {
				kubeconfigContent = encoded
			}
		}
	}

	tempFile, err := os.CreateTemp("", "opsnexus-kubeconfig-*.yaml")
	if err != nil {
		return "", nil, fmt.Errorf("create temp kubeconfig failed: %w", err)
	}
	if _, err := tempFile.Write(kubeconfigContent); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return "", nil, fmt.Errorf("write temp kubeconfig failed: %w", err)
	}
	if err := tempFile.Chmod(0600); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return "", nil, fmt.Errorf("chmod temp kubeconfig failed: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		os.Remove(tempFile.Name())
		return "", nil, fmt.Errorf("close temp kubeconfig failed: %w", err)
	}

	cleanup := func() {
		_ = os.Remove(tempFile.Name())
	}
	return tempFile.Name(), cleanup, nil
}

func normalizeKubectlCommand(command string, namespace string) string {
	trimmed := strings.TrimSpace(command)
	if trimmed == "" {
		return "kubectl"
	}

	if !strings.HasPrefix(trimmed, "kubectl") {
		trimmed = "kubectl " + trimmed
	}

	if namespace == "" || hasKubectlNamespaceFlag(trimmed) {
		return trimmed
	}

	rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "kubectl"))
	return fmt.Sprintf("kubectl -n %s %s", namespace, rest)
}

func hasKubectlNamespaceFlag(command string) bool {
	return strings.Contains(command, " -n ") ||
		strings.Contains(command, " --namespace ") ||
		strings.Contains(command, " --namespace=") ||
		strings.Contains(command, " -n=")
}

func (s *KubectlTerminalSession) writeBanner() error {
	namespaceHint := "current-context"
	if s.namespace != "" {
		namespaceHint = s.namespace
	}

	message := fmt.Sprintf(
		"\x1B[1;32mOpsNexus kubectl terminal connected\x1B[0m\r\n\x1B[1;34mCluster ID: %d\x1B[0m\r\n\x1B[1;33mNamespace: %s\x1B[0m\r\nType kubectl commands directly.\r\n\r\n",
		s.clusterID,
		namespaceHint,
	)
	return s.Conn.WriteJSON(K8sMessage{
		Operation: "stdout",
		Data:      message,
	})
}

func (s *KubectlTerminalSession) pipeOutput() {
	buffer := make([]byte, 4096)
	for {
		n, err := s.ptyFile.Read(buffer)
		if n > 0 {
			if s.recorder != nil {
				s.recorder.RecordOutput(buffer[:n])
			}
			if writeErr := s.Conn.WriteJSON(K8sMessage{
				Operation: "stdout",
				Data:      string(buffer[:n]),
			}); writeErr != nil {
				if s.recorder != nil {
					_ = s.recorder.Close(terminalAuditStatusAborted, writeErr.Error())
				}
				_ = s.Close()
				return
			}
		}

		if err != nil {
			if err != io.EOF {
				_ = s.Conn.WriteJSON(K8sMessage{
					Operation: "stderr",
					Data:      err.Error(),
				})
				if s.recorder != nil {
					_ = s.recorder.Close(terminalAuditStatusAborted, err.Error())
				}
			}
			_ = s.Close()
			return
		}
	}
}

func (s *KubectlTerminalSession) readWebSocketInput() {
	defer s.Close()

	for {
		if s.IsClosed() {
			return
		}

		var message K8sMessage
		if err := s.Conn.ReadJSON(&message); err != nil {
			return
		}

		switch message.Operation {
		case "stdin":
			input, ok := message.Data.(string)
			if !ok {
				continue
			}
			if s.recorder != nil {
				s.recorder.RecordInput([]byte(input))
			}
			if _, err := s.ptyFile.Write([]byte(input)); err != nil {
				if s.recorder != nil {
					_ = s.recorder.Close(terminalAuditStatusAborted, err.Error())
				}
				return
			}
		case "resize":
			cols, rows := extractResizeDimensions(message)
			if cols > 0 && rows > 0 {
				if s.recorder != nil {
					s.recorder.RecordResize(cols, rows)
				}
				_ = pty.Setsize(s.ptyFile, &pty.Winsize{
					Cols: uint16(cols),
					Rows: uint16(rows),
				})
			}
		}
	}
}

func extractResizeDimensions(message K8sMessage) (int, int) {
	if message.Cols > 0 && message.Rows > 0 {
		return message.Cols, message.Rows
	}

	dataMap, ok := message.Data.(map[string]interface{})
	if !ok {
		return 0, 0
	}

	cols, _ := dataMap["cols"].(float64)
	rows, _ := dataMap["rows"].(float64)
	return int(cols), int(rows)
}

func (s *KubectlTerminalSession) waitProcessExit() {
	_ = s.cmd.Wait()
	_ = s.Close()
}

func (s *KubectlTerminalSession) Close() error {
	s.Lock()
	defer s.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true
	if s.cancel != nil {
		s.cancel()
	}
	if s.cmd != nil && s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
	}
	if s.ptyFile != nil {
		_ = s.ptyFile.Close()
	}
	if s.recorder != nil {
		_ = s.recorder.Close(terminalAuditStatusCompleted, "")
	}
	if s.cleanup != nil {
		s.cleanup()
	}
	return nil
}

func (s *KubectlTerminalSession) IsClosed() bool {
	s.RLock()
	defer s.RUnlock()
	return s.closed
}
