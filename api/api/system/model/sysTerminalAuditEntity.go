package model

import "time"

// SysSessionRecording maps to sys_session_recording for write operations.
type SysSessionRecording struct {
	ID              uint       `gorm:"column:id;primaryKey" json:"id"`
	SessionID       string     `gorm:"column:session_id" json:"sessionId"`
	AdminID         uint       `gorm:"column:admin_id" json:"adminId"`
	Username        string     `gorm:"column:username" json:"username"`
	HostID          uint       `gorm:"column:host_id" json:"hostId"`
	HostName        string     `gorm:"column:host_name" json:"hostName"`
	HostIP          string     `gorm:"column:host_ip" json:"hostIp"`
	SSHUser         string     `gorm:"column:ssh_user" json:"sshUser"`
	StartTime       time.Time  `gorm:"column:start_time" json:"startTime"`
	EndTime         *time.Time `gorm:"column:end_time" json:"endTime"`
	Duration        *int64     `gorm:"column:duration" json:"duration"`
	TerminalWidth   int64      `gorm:"column:terminal_width" json:"terminalWidth"`
	TerminalHeight  int64      `gorm:"column:terminal_height" json:"terminalHeight"`
	FilePath        string     `gorm:"column:file_path" json:"filePath"`
	FileSize        int64      `gorm:"column:file_size" json:"fileSize"`
	StorageType     int64      `gorm:"column:storage_type" json:"storageType"`
	OssKey          string     `gorm:"column:oss_key" json:"ossKey"`
	InputCount      int64      `gorm:"column:input_count" json:"inputCount"`
	OutputCount     int64      `gorm:"column:output_count" json:"outputCount"`
	ResizeCount     int64      `gorm:"column:resize_count" json:"resizeCount"`
	CommandCount    int64      `gorm:"column:command_count" json:"commandCount"`
	ClientIP        string     `gorm:"column:client_ip" json:"clientIp"`
	UserAgent       string     `gorm:"column:user_agent" json:"userAgent"`
	RiskLevel       int64      `gorm:"column:risk_level" json:"riskLevel"`
	HasSensitiveCmd bool       `gorm:"column:has_sensitive_cmd" json:"hasSensitiveCmd"`
	Status          int64      `gorm:"column:status" json:"status"`
	ErrorMsg        string     `gorm:"column:error_msg" json:"errorMsg"`
	CreateTime      time.Time  `gorm:"column:create_time" json:"createTime"`
	UpdateTime      time.Time  `gorm:"column:update_time" json:"updateTime"`
	DeleteTime      *time.Time `gorm:"column:delete_time" json:"deleteTime"`
}

func (SysSessionRecording) TableName() string {
	return "sys_session_recording"
}

// SysCommandAudit maps to sys_command_audit for write operations.
type SysCommandAudit struct {
	ID          uint      `gorm:"column:id;primaryKey" json:"id"`
	RecordingID uint      `gorm:"column:recording_id" json:"recordingId"`
	SessionID   string    `gorm:"column:session_id" json:"sessionId"`
	Command     string    `gorm:"column:command" json:"command"`
	Timestamp   float64   `gorm:"column:timestamp" json:"timestamp"`
	Sequence    int64     `gorm:"column:sequence" json:"sequence"`
	IsSensitive bool      `gorm:"column:is_sensitive" json:"isSensitive"`
	RiskLevel   int64     `gorm:"column:risk_level" json:"riskLevel"`
	RiskReason  string    `gorm:"column:risk_reason" json:"riskReason"`
	ExecuteTime time.Time `gorm:"column:execute_time" json:"executeTime"`
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`
}

func (SysCommandAudit) TableName() string {
	return "sys_command_audit"
}
