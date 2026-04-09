package service

import (
	"context"
	"fmt"
	"io"
	"sync"

	"dodevops-api/api/k8s/dao"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

const (
	terminalAuditStatusCompleted = 2
	terminalAuditStatusAborted   = 3
)

// TerminalAuditRecorder captures terminal input/output for audit persistence.
type TerminalAuditRecorder interface {
	RecordInput(data []byte)
	RecordOutput(data []byte)
	RecordResize(cols, rows int)
	Close(status int64, errMsg string) error
}

// IK8sTerminalService 容器终端服务接口.
type IK8sTerminalService interface {
	CreateK8sWebSocketStream(clusterID uint, namespaceName, podName, containerName, command string, conn *websocket.Conn, recorder TerminalAuditRecorder) (*K8sWebSocketStream, error)
	GetPodContainers(clusterID uint, namespaceName, podName string) ([]string, error)
}

// K8sMessage WebSocket 消息结构.
type K8sMessage struct {
	Operation string      `json:"operation"`
	Data      interface{} `json:"data"`
	Cols      int         `json:"cols,omitempty"`
	Rows      int         `json:"rows,omitempty"`
}

// K8sWebSocketStream K8s WebSocket 流处理.
type K8sWebSocketStream struct {
	sync.RWMutex
	Conn     *websocket.Conn
	executor remotecommand.Executor
	Ctx      context.Context
	cancel   context.CancelFunc
	closed   bool
	reader   *io.PipeReader
	writer   *io.PipeWriter
	recorder TerminalAuditRecorder
}

type terminalConn struct {
	stream *K8sWebSocketStream
}

func (tc *terminalConn) Read(p []byte) (n int, err error) {
	return tc.stream.reader.Read(p)
}

func (tc *terminalConn) Write(p []byte) (n int, err error) {
	return tc.stream.WriteToWebSocket(p)
}

// Close closes the WebSocket stream and finalizes audit state.
func (kws *K8sWebSocketStream) Close() error {
	kws.Lock()
	defer kws.Unlock()

	if kws.closed {
		return nil
	}

	kws.closed = true
	if kws.cancel != nil {
		kws.cancel()
	}
	if kws.recorder != nil {
		_ = kws.recorder.Close(terminalAuditStatusCompleted, "")
	}
	return nil
}

func (kws *K8sWebSocketStream) closeWithError(err error) {
	kws.Lock()
	defer kws.Unlock()
	if kws.closed {
		return
	}
	kws.closed = true
	if kws.cancel != nil {
		kws.cancel()
	}
	if kws.recorder != nil && err != nil {
		_ = kws.recorder.Close(terminalAuditStatusAborted, err.Error())
	}
}

// IsClosed checks whether the stream is already closed.
func (kws *K8sWebSocketStream) IsClosed() bool {
	kws.RLock()
	defer kws.RUnlock()
	return kws.closed
}

// WriteToWebSocket writes pod terminal output back to the WebSocket.
func (kws *K8sWebSocketStream) WriteToWebSocket(p []byte) (n int, err error) {
	if kws.IsClosed() {
		return 0, io.EOF
	}

	if kws.recorder != nil {
		kws.recorder.RecordOutput(p)
	}

	message := K8sMessage{
		Operation: "stdout",
		Data:      string(p),
	}

	if err = kws.Conn.WriteJSON(message); err != nil {
		kws.closeWithError(err)
		return 0, err
	}

	return len(p), nil
}

// ReadFromWebSocket reads pod terminal input from the WebSocket.
func (kws *K8sWebSocketStream) ReadFromWebSocket() {
	defer func() {
		_ = kws.Close()
		if kws.writer != nil {
			_ = kws.writer.Close()
		}
	}()

	for {
		if kws.IsClosed() {
			return
		}

		var message K8sMessage
		err := kws.Conn.ReadJSON(&message)
		if err != nil {
			kws.closeWithError(err)
			return
		}

		switch message.Operation {
		case "stdin":
			input, ok := message.Data.(string)
			if !ok {
				continue
			}
			if kws.recorder != nil {
				kws.recorder.RecordInput([]byte(input))
			}
			if _, err = kws.writer.Write([]byte(input)); err != nil {
				kws.closeWithError(err)
				return
			}
		case "resize":
			if kws.recorder != nil {
				kws.recorder.RecordResize(message.Cols, message.Rows)
			}
		}
	}
}

// K8sTerminalService 容器终端服务实现.
type K8sTerminalService struct {
	dao *dao.KubeClusterDao
}

func NewK8sTerminalService(db *gorm.DB) IK8sTerminalService {
	return &K8sTerminalService{
		dao: dao.NewKubeClusterDao(db),
	}
}

// CreateK8sWebSocketStream creates a pod exec stream and wires audit hooks.
func (s *K8sTerminalService) CreateK8sWebSocketStream(clusterID uint, namespaceName, podName, containerName, command string, conn *websocket.Conn, recorder TerminalAuditRecorder) (*K8sWebSocketStream, error) {
	cluster, err := s.dao.GetByID(clusterID)
	if err != nil {
		return nil, fmt.Errorf("获取集群信息失败: %v", err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Credential))
	if err != nil {
		return nil, fmt.Errorf("创建 Kubernetes 配置失败: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建 Kubernetes 客户端失败: %v", err)
	}

	if containerName == "" {
		pod, getErr := clientset.CoreV1().Pods(namespaceName).Get(context.Background(), podName, metav1.GetOptions{})
		if getErr != nil {
			return nil, fmt.Errorf("获取 Pod 信息失败: %v", getErr)
		}
		if len(pod.Spec.Containers) == 0 {
			return nil, fmt.Errorf("Pod 中没有找到容器")
		}
		containerName = pod.Spec.Containers[0].Name
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespaceName).
		SubResource("exec").
		Param("container", containerName).
		Param("command", command).
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "true")

	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return nil, fmt.Errorf("创建 executor 失败: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	reader, writer := io.Pipe()

	stream := &K8sWebSocketStream{
		Conn:     conn,
		executor: executor,
		Ctx:      ctx,
		cancel:   cancel,
		reader:   reader,
		writer:   writer,
		recorder: recorder,
	}

	termConn := &terminalConn{stream: stream}
	go stream.ReadFromWebSocket()

	go func() {
		defer func() {
			cancel()
			if reader != nil {
				_ = reader.Close()
			}
		}()

		streamErr := executor.StreamWithContext(ctx, remotecommand.StreamOptions{
			Stdin:  termConn,
			Stdout: termConn,
			Stderr: termConn,
			Tty:    true,
		})
		if streamErr != nil {
			stream.closeWithError(streamErr)
		}
	}()

	return stream, nil
}

// GetPodContainers returns all containers in one pod.
func (s *K8sTerminalService) GetPodContainers(clusterID uint, namespaceName, podName string) ([]string, error) {
	cluster, err := s.dao.GetByID(clusterID)
	if err != nil {
		return nil, fmt.Errorf("获取集群信息失败: %v", err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Credential))
	if err != nil {
		return nil, fmt.Errorf("创建 Kubernetes 配置失败: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建 Kubernetes 客户端失败: %v", err)
	}

	pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取 Pod 信息失败: %v", err)
	}

	var containers []string
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	for _, container := range pod.Spec.InitContainers {
		containers = append(containers, container.Name+" (init)")
	}

	return containers, nil
}
