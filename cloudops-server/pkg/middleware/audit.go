package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditRecord 定义审计记录。
type AuditRecord struct {
	UserID       int64
	Username     string
	IP           string
	Method       string
	Path         string
	RequestBody  string
	ResponseCode int
	Module       string
	Action       string
	ResourceType string
	ResourceID   string
	Description  string
	DurationMS   int64
}

// AuditRecorder 定义审计记录回调。
type AuditRecorder func(ctx context.Context, record *AuditRecord) error

// Audit 记录增删改审计日志。
func Audit(log *slog.Logger, recorder AuditRecorder) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestBody := captureRequestBody(c)

		c.Next()

		claims, _ := GetJWTClaims(c)
		userID := int64(0)
		username := "anonymous"
		if claims != nil {
			userID = claims.UserID
			username = claims.Username
		}

		module, action, resourceType, resourceID := parseRouteMeta(c.Request.Method, c.Request.URL.Path)
		durationMS := time.Since(start).Milliseconds()
		record := &AuditRecord{
			UserID:       userID,
			Username:     username,
			IP:           clientIP(c),
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			RequestBody:  requestBody,
			ResponseCode: c.Writer.Status(),
			Module:       module,
			Action:       action,
			ResourceType: resourceType,
			ResourceID:   resourceID,
			Description:  fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
			DurationMS:   durationMS,
		}

		if recorder != nil {
			if err := recorder(c.Request.Context(), record); err != nil {
				log.Warn("写入审计日志失败", "error", err)
			}
		}

		log.Info("审计日志",
			"request_id", GetRequestID(c),
			"user_id", userID,
			"username", username,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start).String(),
		)
	}
}

func captureRequestBody(c *gin.Context) string {
	if c.Request == nil || c.Request.Body == nil {
		return ""
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	raw := strings.TrimSpace(string(body))
	if raw == "" {
		return ""
	}
	lower := strings.ToLower(raw)
	if strings.Contains(lower, "password") {
		return "{\"masked\":true}"
	}
	if len(raw) > 2000 {
		return raw[:2000]
	}
	return raw
}

func parseRouteMeta(method, path string) (module, action, resourceType, resourceID string) {
	module = "unknown"
	action = "query"
	resourceType = ""
	resourceID = ""

	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) >= 3 && segments[0] == "api" && segments[1] == "v1" {
		module = segments[2]
	}
	if len(segments) >= 4 {
		resourceType = segments[3]
	}
	if len(segments) >= 5 {
		resourceID = segments[4]
	}

	switch strings.ToUpper(method) {
	case "POST":
		action = "create"
	case "PUT", "PATCH":
		action = "update"
	case "DELETE":
		action = "delete"
	default:
		action = "query"
	}
	return
}

func clientIP(c *gin.Context) string {
	ip := strings.TrimSpace(c.ClientIP())
	if ip == "" {
		return "unknown"
	}
	return ip
}
