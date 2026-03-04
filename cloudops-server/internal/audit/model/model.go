package model

import "time"

// AuditLog 审计日志。
type AuditLog struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID       int64      `gorm:"column:user_id;not null;default:0" json:"user_id"`
	Username     string     `gorm:"column:username;size:64;not null;default:'anonymous'" json:"username"`
	IP           string     `gorm:"column:ip;size:64;not null;default:''" json:"ip"`
	Method       string     `gorm:"column:method;size:16;not null" json:"method"`
	Path         string     `gorm:"column:path;size:512;not null" json:"path"`
	RequestBody  string     `gorm:"column:request_body;type:text;not null;default:''" json:"request_body"`
	ResponseCode int        `gorm:"column:response_code;not null;default:0" json:"response_code"`
	Module       string     `gorm:"column:module;size:64;not null;default:''" json:"module"`
	Action       string     `gorm:"column:action;size:64;not null;default:''" json:"action"`
	ResourceType string     `gorm:"column:resource_type;size:64;not null;default:''" json:"resource_type"`
	ResourceID   string     `gorm:"column:resource_id;size:64;not null;default:''" json:"resource_id"`
	Description  string     `gorm:"column:description;size:512;not null;default:''" json:"description"`
	DurationMS   int64      `gorm:"column:duration_ms;not null;default:0" json:"duration_ms"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (AuditLog) TableName() string { return "audit_logs" }
