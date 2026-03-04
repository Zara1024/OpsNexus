package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

// NotifyService 提供通知发送能力。
type NotifyService struct{}

func NewNotifyService() *NotifyService { return &NotifyService{} }

// SendTest 发送测试通知。
func (s *NotifyService) SendTest(channelType string, cfg map[string]string) error {
	message := "[OpsNexus] monitor channel test"
	switch strings.ToLower(channelType) {
	case "webhook", "dingtalk", "feishu", "wechat":
		url := strings.TrimSpace(cfg["webhook"])
		if url == "" {
			url = strings.TrimSpace(cfg["url"])
		}
		if url == "" {
			return fmt.Errorf("webhook 地址不能为空")
		}
		payload, _ := json.Marshal(map[string]any{"msg": message, "timestamp": time.Now().Unix()})
		resp, err := http.Post(url, "application/json", bytes.NewReader(payload))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			return fmt.Errorf("通知失败，状态码: %d", resp.StatusCode)
		}
		return nil
	case "email":
		host := strings.TrimSpace(cfg["host"])
		port := strings.TrimSpace(cfg["port"])
		user := strings.TrimSpace(cfg["username"])
		pass := cfg["password"]
		to := strings.TrimSpace(cfg["to"])
		if host == "" || port == "" || user == "" || pass == "" || to == "" {
			return fmt.Errorf("邮件配置不完整")
		}
		addr := host + ":" + port
		auth := smtp.PlainAuth("", user, pass, host)
		body := []byte("To: " + to + "\r\nSubject: OpsNexus Test\r\n\r\n" + message)
		return smtp.SendMail(addr, auth, user, []string{to}, body)
	default:
		return fmt.Errorf("暂不支持的通知类型: %s", channelType)
	}
}
