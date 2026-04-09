package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

func (s *MonitorAlertServiceImpl) ReceiveWebhook(c *gin.Context, payload map[string]interface{}) {
	req, err := normalizeWebhookPayload(payload)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}

	robots, err := dao.GetEnabledMonitorNotifyRobots(req.NotifyRobotIDs)
	if err != nil {
		result.Failed(c, 500, "获取通知机器人失败: "+err.Error())
		return
	}

	logEntry, summary, err := dispatchAlert(req, robots)
	if err != nil {
		result.Failed(c, 500, "处理 webhook 派发失败: "+err.Error())
		return
	}

	result.Success(c, gin.H{
		"webhookLogId": logEntry.ID,
		"status":       summary.Status,
		"notifyCount":  summary.NotifyCount,
		"successCount": summary.SuccessCount,
		"failedCount":  summary.FailedCount,
	})
}

func normalizeWebhookPayload(payload map[string]interface{}) (model.MonitorWebhookReceiveRequest, error) {
	req := model.MonitorWebhookReceiveRequest{
		Source:         strings.TrimSpace(asString(payload["source"])),
		Title:          strings.TrimSpace(asString(payload["title"])),
		Content:        strings.TrimSpace(asString(payload["content"])),
		Level:          strings.TrimSpace(asString(payload["level"])),
		Tags:           payload["tags"],
		Extra:          payload["extra"],
		NotifyRobotIDs: toUintSlice(payload["notifyRobotIds"]),
	}

	if req.Source == "" {
		req.Source = strings.TrimSpace(asString(payload["sourceType"]))
	}

	if req.Title == "" || req.Content == "" {
		if alertManagerReq, ok := parseAlertManagerPayload(payload); ok {
			if req.Source == "" {
				req.Source = alertManagerReq.Source
			}
			if req.Title == "" {
				req.Title = alertManagerReq.Title
			}
			if req.Content == "" {
				req.Content = alertManagerReq.Content
			}
			if req.Level == "" {
				req.Level = alertManagerReq.Level
			}
			if req.Tags == nil {
				req.Tags = alertManagerReq.Tags
			}
			if req.Extra == nil {
				req.Extra = alertManagerReq.Extra
			}
		}
	}

	if req.Title == "" {
		req.Title = firstNonEmpty(
			strings.TrimSpace(asString(payload["message"])),
			strings.TrimSpace(asString(payload["summary"])),
			strings.TrimSpace(asString(payload["alertName"])),
		)
	}
	if req.Content == "" {
		req.Content = firstNonEmpty(
			strings.TrimSpace(asString(payload["message"])),
			strings.TrimSpace(asString(payload["description"])),
			req.Title,
		)
	}
	if req.Level == "" {
		req.Level = firstNonEmpty(
			strings.TrimSpace(asString(payload["severity"])),
			strings.TrimSpace(asString(payload["status"])),
			"info",
		)
	}
	if req.Source == "" {
		req.Source = "custom"
	}
	if req.Tags == nil {
		req.Tags = map[string]interface{}{}
	}
	if req.Extra == nil {
		req.Extra = payload
	}

	if req.Title == "" {
		return req, errors.New("webhook title 不能为空")
	}
	if req.Content == "" {
		return req, errors.New("webhook content 不能为空")
	}
	return req, nil
}

func parseAlertManagerPayload(payload map[string]interface{}) (model.MonitorWebhookReceiveRequest, bool) {
	alertsRaw, ok := payload["alerts"].([]interface{})
	if !ok || len(alertsRaw) == 0 {
		return model.MonitorWebhookReceiveRequest{}, false
	}

	firstAlert, ok := alertsRaw[0].(map[string]interface{})
	if !ok {
		return model.MonitorWebhookReceiveRequest{}, false
	}

	labels := asStringMap(firstAlert["labels"])
	annotations := asStringMap(firstAlert["annotations"])
	alertName := firstNonEmpty(labels["alertname"], annotations["summary"], annotations["description"], "AlertManager webhook")
	severity := firstNonEmpty(labels["severity"], asString(payload["status"]), "info")
	instance := firstNonEmpty(labels["instance"], labels["pod"], labels["job"])
	status := strings.TrimSpace(asString(payload["status"]))
	summary := firstNonEmpty(annotations["summary"], annotations["description"])
	contentParts := make([]string, 0, 4)
	if status != "" {
		contentParts = append(contentParts, "status="+status)
	}
	if summary != "" {
		contentParts = append(contentParts, summary)
	}
	if instance != "" {
		contentParts = append(contentParts, "instance="+instance)
	}
	if len(alertsRaw) > 1 {
		contentParts = append(contentParts, fmt.Sprintf("alerts=%d", len(alertsRaw)))
	}

	return model.MonitorWebhookReceiveRequest{
		Source: "alertmanager",
		Title:  alertName,
		Content: firstNonEmpty(
			strings.Join(contentParts, "\n"),
			alertName,
		),
		Level: severity,
		Tags:  labels,
		Extra: payload,
	}, true
}

func calculateDispatchStatus(req model.MonitorWebhookReceiveRequest, robotCount int, successCount, failedCount int64) string {
	if robotCount == 0 {
		if len(req.NotifyRobotIDs) > 0 {
			return "failed"
		}
		return "success"
	}
	if failedCount == 0 {
		return "success"
	}
	if successCount == 0 {
		return "failed"
	}
	return "partial"
}

func sendAlertToRobot(robot *model.MonitorNotifyRobotEntity, req model.MonitorWebhookReceiveRequest) error {
	if robot.Type == "email" {
		return sendAlertEmail(robot, req)
	}
	return sendAlertWebhook(robot, req)
}

func sendAlertWebhook(robot *model.MonitorNotifyRobotEntity, req model.MonitorWebhookReceiveRequest) error {
	if strings.TrimSpace(robot.Webhook) == "" {
		return errors.New("webhook 地址为空")
	}

	method := strings.ToUpper(strings.TrimSpace(robot.Method))
	if method == "" {
		method = "POST"
	}

	body, contentType, err := buildWebhookBody(robot, req)
	if err != nil {
		return err
	}

	var bodyReader io.Reader
	if method != "GET" {
		bodyReader = bytes.NewReader(body)
	}

	httpReq, err := http.NewRequest(method, robot.Webhook, bodyReader)
	if err != nil {
		return err
	}
	if contentType != "" && method != "GET" {
		httpReq.Header.Set("Content-Type", contentType)
	}

	headers, err := parseHeaderJSON(robot.Headers)
	if err != nil {
		return err
	}
	for key, value := range headers {
		httpReq.Header.Set(key, value)
	}
	if strings.TrimSpace(robot.Secret) != "" && robot.Type == "webhook" {
		httpReq.Header.Set("X-Webhook-Secret", strings.TrimSpace(robot.Secret))
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 300))
		return fmt.Errorf("downstream status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return nil
}

func sendAlertEmail(robot *model.MonitorNotifyRobotEntity, req model.MonitorWebhookReceiveRequest) error {
	if strings.TrimSpace(robot.Server) == "" || robot.Port <= 0 {
		return errors.New("SMTP 服务配置不完整")
	}
	recipients := splitRecipients(robot.Webhook)
	if len(recipients) == 0 {
		return errors.New("邮件收件地址为空")
	}

	body := renderAlertTemplate(robot.Template, req)
	if body == "" {
		body = buildAlertMessage(req)
	}

	from := robot.Username
	if strings.TrimSpace(robot.Nickname) != "" {
		from = fmt.Sprintf("%s <%s>", strings.TrimSpace(robot.Nickname), robot.Username)
	}
	headers := []string{
		"From: " + from,
		"To: " + strings.Join(recipients, ","),
		"Subject: " + req.Title,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}

	addr := fmt.Sprintf("%s:%d", strings.TrimSpace(robot.Server), robot.Port)
	auth := smtp.PlainAuth("", robot.Username, robot.Password, strings.TrimSpace(robot.Server))
	return smtp.SendMail(addr, auth, robot.Username, recipients, []byte(strings.Join(headers, "\r\n")))
}

func buildWebhookBody(robot *model.MonitorNotifyRobotEntity, req model.MonitorWebhookReceiveRequest) ([]byte, string, error) {
	rendered := renderAlertTemplate(robot.Template, req)
	message := rendered
	if message == "" {
		message = buildAlertMessage(req)
	}

	if robot.Type == "webhook" && isJSONDocument(rendered) {
		return []byte(rendered), "application/json", nil
	}

	switch robot.Type {
	case "feishu":
		body, err := json.Marshal(map[string]interface{}{
			"msg_type": "text",
			"content": map[string]string{
				"text": message,
			},
		})
		return body, "application/json", err
	case "dingtalk":
		body, err := json.Marshal(map[string]interface{}{
			"msgtype": "text",
			"text": map[string]string{
				"content": message,
			},
		})
		return body, "application/json", err
	case "wechat":
		body, err := json.Marshal(map[string]interface{}{
			"msgtype": "text",
			"text": map[string]string{
				"content": message,
			},
		})
		return body, "application/json", err
	case "teams":
		if isJSONDocument(rendered) {
			return []byte(rendered), "application/json", nil
		}
		body, err := json.Marshal(map[string]string{"text": message})
		return body, "application/json", err
	default:
		body, err := json.Marshal(map[string]interface{}{
			"source":  req.Source,
			"title":   req.Title,
			"content": req.Content,
			"level":   req.Level,
			"message": message,
			"tags":    req.Tags,
			"extra":   req.Extra,
		})
		return body, "application/json", err
	}
}

func buildAlertMessage(req model.MonitorWebhookReceiveRequest) string {
	lines := []string{
		"[OpsNexus Alert]",
		"source: " + req.Source,
		"title: " + req.Title,
		"level: " + req.Level,
		"content: " + req.Content,
		"time: " + time.Now().Format("2006-01-02 15:04:05"),
	}
	return strings.Join(lines, "\n")
}

func renderAlertTemplate(template string, req model.MonitorWebhookReceiveRequest) string {
	template = strings.TrimSpace(template)
	if template == "" {
		return ""
	}

	values := map[string]string{
		"title":   req.Title,
		"content": req.Content,
		"level":   req.Level,
		"source":  req.Source,
		"tags":    encodeJSONValue(req.Tags),
		"extra":   encodeJSONValue(req.Extra),
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	}

	rendered := template
	for key, value := range values {
		rendered = strings.ReplaceAll(rendered, "{{"+key+"}}", value)
	}
	return rendered
}

func parseHeaderJSON(raw string) (map[string]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return map[string]string{}, nil
	}

	var headers map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &headers); err != nil {
		return nil, errors.New("机器人自定义 Header 不是合法 JSON")
	}

	result := make(map[string]string, len(headers))
	for key, value := range headers {
		result[key] = asString(value)
	}
	return result, nil
}

func splitRecipients(raw string) []string {
	items := strings.Split(raw, ",")
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func toUintSlice(value interface{}) []uint {
	values, ok := value.([]interface{})
	if !ok {
		return []uint{}
	}

	result := make([]uint, 0, len(values))
	for _, item := range values {
		switch actual := item.(type) {
		case float64:
			if actual > 0 {
				result = append(result, uint(actual))
			}
		case int:
			if actual > 0 {
				result = append(result, uint(actual))
			}
		case string:
			actual = strings.TrimSpace(actual)
			if actual == "" {
				continue
			}
			var parsed uint64
			if _, err := fmt.Sscanf(actual, "%d", &parsed); err == nil && parsed > 0 {
				result = append(result, uint(parsed))
			}
		}
	}
	return result
}

func asString(value interface{}) string {
	switch actual := value.(type) {
	case nil:
		return ""
	case string:
		return actual
	default:
		return fmt.Sprint(actual)
	}
}

func asStringMap(value interface{}) map[string]string {
	source, ok := value.(map[string]interface{})
	if !ok {
		return map[string]string{}
	}

	result := make(map[string]string, len(source))
	for key, item := range source {
		result[key] = strings.TrimSpace(asString(item))
	}
	return result
}

func encodeJSONValue(value interface{}) string {
	if value == nil {
		return ""
	}
	data, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(data)
}

func truncateText(value string, limit int) string {
	if limit <= 0 || len(value) <= limit {
		return value
	}
	return value[:limit]
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func isJSONDocument(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	if !strings.HasPrefix(value, "{") && !strings.HasPrefix(value, "[") {
		return false
	}
	return json.Valid([]byte(value))
}
