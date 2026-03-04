package model

import "time"

// AlertRule 告警规则。
type AlertRule struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name           string     `gorm:"column:name;size:128;not null" json:"name"`
	Description    string     `gorm:"column:description;size:512;not null;default:''" json:"description"`
	Severity       string     `gorm:"column:severity;size:8;not null;default:P3" json:"severity"`
	Expression     string     `gorm:"column:expression;type:text;not null" json:"expression"`
	Duration       string     `gorm:"column:duration;size:64;not null;default:5m" json:"duration"`
	Labels         string     `gorm:"column:labels;type:jsonb;not null;default:'{}'" json:"labels"`
	Annotations    string     `gorm:"column:annotations;type:jsonb;not null;default:'{}'" json:"annotations"`
	Enabled        bool       `gorm:"column:enabled;not null;default:true" json:"enabled"`
	NotifyChannels string     `gorm:"column:notify_channels;type:jsonb;not null;default:'[]'" json:"notify_channels"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (AlertRule) TableName() string { return "monitor_alert_rules" }

// AlertEvent 告警事件。
type AlertEvent struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RuleID         int64      `gorm:"column:rule_id;not null" json:"rule_id"`
	Fingerprint    string     `gorm:"column:fingerprint;size:128;not null" json:"fingerprint"`
	Status         string     `gorm:"column:status;size:16;not null;default:firing" json:"status"`
	Severity       string     `gorm:"column:severity;size:8;not null;default:P3" json:"severity"`
	StartsAt       time.Time  `gorm:"column:starts_at" json:"starts_at"`
	EndsAt         *time.Time `gorm:"column:ends_at" json:"ends_at"`
	Labels         string     `gorm:"column:labels;type:jsonb;not null;default:'{}'" json:"labels"`
	Annotations    string     `gorm:"column:annotations;type:jsonb;not null;default:'{}'" json:"annotations"`
	AcknowledgedBy int64      `gorm:"column:acknowledged_by;not null;default:0" json:"acknowledged_by"`
	SilencedUntil  *time.Time `gorm:"column:silenced_until" json:"silenced_until"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (AlertEvent) TableName() string { return "monitor_alert_events" }

// NotificationChannel 通知渠道。
type NotificationChannel struct {
	ID              int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name            string     `gorm:"column:name;size:128;not null" json:"name"`
	Type            string     `gorm:"column:type;size:32;not null" json:"type"`
	ConfigEncrypted string     `gorm:"column:config_encrypted;type:text;not null" json:"-"`
	Enabled         bool       `gorm:"column:enabled;not null;default:true" json:"enabled"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (NotificationChannel) TableName() string { return "monitor_notification_channels" }
