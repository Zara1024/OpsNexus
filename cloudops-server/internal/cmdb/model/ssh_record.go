package model

import "time"

// SSHRecord 定义 SSH 会话审计记录。
type SSHRecord struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID        int64      `gorm:"column:user_id;not null;default:0" json:"user_id"`
	HostID        int64      `gorm:"column:host_id;not null;default:0" json:"host_id"`
	SessionID     string     `gorm:"column:session_id;size:128;not null" json:"session_id"`
	StartTime     time.Time  `gorm:"column:start_time;not null" json:"start_time"`
	EndTime       *time.Time `gorm:"column:end_time" json:"end_time"`
	ClientIP      string     `gorm:"column:client_ip;size:64;not null;default:''" json:"client_ip"`
	RecordingPath string     `gorm:"column:recording_path;size:255;not null;default:''" json:"recording_path"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"-"`
}

// TableName 指定表名。
func (SSHRecord) TableName() string {
	return "cmdb_ssh_records"
}
