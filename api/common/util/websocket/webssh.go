package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

// AuditRecorder captures terminal input/output for session audit.
type AuditRecorder interface {
	RecordInput(data []byte)
	RecordOutput(data []byte)
	RecordResize(cols, rows int)
	Close(status int64, errMsg string) error
}

// WebSSH 简化的 WebSSH 实现.
type WebSSH struct {
	conn      *websocket.Conn
	client    *ssh.Client
	session   *ssh.Session
	stdinPipe io.WriteCloser
	cancel    context.CancelFunc
	recorder  AuditRecorder
}

// GetStdinPipe 获取输入管道.
func (w *WebSSH) GetStdinPipe() io.WriteCloser {
	if w.stdinPipe == nil && w.session != nil {
		var err error
		w.stdinPipe, err = w.session.StdinPipe()
		if err != nil {
			log.Printf("初始化 SSH 输入管道失败: %v", err)
			return nil
		}
	}
	return w.stdinPipe
}

// GetSession 获取 SSH 会话.
func (w *WebSSH) GetSession() *ssh.Session {
	return w.session
}

// SetAuditRecorder attaches an optional session recorder.
func (w *WebSSH) SetAuditRecorder(recorder AuditRecorder) {
	w.recorder = recorder
}

// NewWebSSH 创建新的 WebSSH 连接(密码认证).
func NewWebSSH(host, port, user, password string) (*WebSSH, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
				if len(questions) == 0 {
					return []string{}, nil
				}
				return []string{password}, nil
			}),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
		Config: ssh.Config{
			Ciphers: []string{
				"aes128-ctr",
				"aes192-ctr",
				"aes256-ctr",
				"aes128-gcm@openssh.com",
				"arcfour256",
				"arcfour128",
				"arcfour",
				"aes128-cbc",
				"3des-cbc",
			},
		},
	}

	return newWebSSHWithConfig(host, port, config)
}

// NewWebSSHWithAuth 创建新的 WebSSH 连接(自定义认证方法).
func NewWebSSHWithAuth(host, port, user string, authMethod ssh.AuthMethod) (*WebSSH, error) {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
		Config: ssh.Config{
			Ciphers: []string{
				"aes128-ctr",
				"aes192-ctr",
				"aes256-ctr",
				"aes128-gcm@openssh.com",
				"arcfour256",
				"arcfour128",
				"arcfour",
				"aes128-cbc",
				"3des-cbc",
			},
		},
	}

	webSSH, err := newWebSSHWithConfig(host, port, config)
	if err != nil {
		return nil, err
	}

	session, err := webSSH.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("创建 SSH 会话失败: %v", err)
	}
	webSSH.session = session

	stdinPipe, err := session.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("创建输入管道失败: %v", err)
	}
	webSSH.stdinPipe = stdinPipe

	return webSSH, nil
}

func newWebSSHWithConfig(host, port string, config *ssh.ClientConfig) (*WebSSH, error) {
	addr := net.JoinHostPort(host, port)
	log.Printf("尝试连接 SSH 服务器: %s", addr)

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Printf("SSH 连接失败详情: %v", err)
		return nil, fmt.Errorf("failed to dial SSH server: %v", err)
	}

	return &WebSSH{client: client}, nil
}

// Connect 建立 SSH 会话并处理 WebSocket 连接.
func (w *WebSSH) Connect(wsConn *websocket.Conn) error {
	session, err := w.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	w.session = session
	w.conn = wsConn

	w.stdinPipe, err = session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %v", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %v", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
		ssh.ONLCR:         1,
		ssh.OCRNL:         0,
		ssh.INLCR:         0,
		ssh.IGNCR:         0,
		ssh.ICRNL:         1,
		ssh.OPOST:         1,
		ssh.ONLRET:        0,
		ssh.ONOCR:         0,
	}
	width := 160
	height := 40
	if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
		return fmt.Errorf("request pty failed: %v", err)
	}
	if err := session.Shell(); err != nil {
		return fmt.Errorf("start shell failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel

	go w.handleInput(ctx)
	go w.handleOutput(ctx, stdout)
	return nil
}

func (w *WebSSH) handleInput(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, data, err := w.conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				if w.recorder != nil {
					_ = w.recorder.Close(3, err.Error())
				}
				return
			}

			var msg map[string]interface{}
			if err = json.Unmarshal(data, &msg); err == nil {
				if msg["type"] == "resize" {
					cols, _ := strconv.Atoi(fmt.Sprintf("%v", msg["cols"]))
					rows, _ := strconv.Atoi(fmt.Sprintf("%v", msg["rows"]))
					if cols > 0 && rows > 0 {
						_ = w.session.WindowChange(rows, cols)
						if w.recorder != nil {
							w.recorder.RecordResize(cols, rows)
						}
					}
					continue
				}
			}

			if w.recorder != nil {
				w.recorder.RecordInput(data)
			}
			if _, err = w.stdinPipe.Write(data); err != nil {
				log.Printf("SSH write error: %v", err)
				if w.recorder != nil {
					_ = w.recorder.Close(3, err.Error())
				}
				return
			}
		}
	}
}

func (w *WebSSH) handleOutput(ctx context.Context, stdout io.Reader) {
	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := stdout.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("SSH read error: %v", err)
					if w.recorder != nil {
						_ = w.recorder.Close(3, err.Error())
					}
				}
				return
			}
			if n > 0 {
				if w.recorder != nil {
					w.recorder.RecordOutput(buf[:n])
				}
				if err = w.conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
					log.Printf("WebSocket write error: %v", err)
					if w.recorder != nil {
						_ = w.recorder.Close(3, err.Error())
					}
					return
				}
			}
		}
	}
}

// Close 关闭连接.
func (w *WebSSH) Close() error {
	if w.cancel != nil {
		w.cancel()
	}
	if w.session != nil {
		_ = w.session.Close()
	}
	if w.recorder != nil {
		_ = w.recorder.Close(2, "")
	}
	if w.client != nil {
		return w.client.Close()
	}
	return nil
}
