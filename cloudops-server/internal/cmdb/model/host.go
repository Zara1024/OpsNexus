package model

import "time"

// Host 定义 CMDB 主机资产。
type Host struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Hostname      string     `gorm:"column:hostname;size:128;not null" json:"hostname"`
	InnerIP       string     `gorm:"column:inner_ip;size:64;not null;default:''" json:"inner_ip"`
	OuterIP       string     `gorm:"column:outer_ip;size:64;not null;default:''" json:"outer_ip"`
	OSType        string     `gorm:"column:os_type;size:64;not null;default:''" json:"os_type"`
	OSVersion     string     `gorm:"column:os_version;size:128;not null;default:''" json:"os_version"`
	Arch          string     `gorm:"column:arch;size:64;not null;default:''" json:"arch"`
	CPUCores      int        `gorm:"column:cpu_cores;not null;default:0" json:"cpu_cores"`
	MemoryGB      float64    `gorm:"column:memory_gb;type:numeric(10,2);not null;default:0" json:"memory_gb"`
	DiskGB        float64    `gorm:"column:disk_gb;type:numeric(10,2);not null;default:0" json:"disk_gb"`
	GroupID       int64      `gorm:"column:group_id;not null;default:0" json:"group_id"`
	CredentialID  int64      `gorm:"column:credential_id;not null;default:0" json:"credential_id"`
	CloudProvider string     `gorm:"column:cloud_provider;size:64;not null;default:''" json:"cloud_provider"`
	InstanceID    string     `gorm:"column:instance_id;size:128;not null;default:''" json:"instance_id"`
	Region        string     `gorm:"column:region;size:64;not null;default:''" json:"region"`
	AgentStatus   string     `gorm:"column:agent_status;size:32;not null;default:unknown" json:"agent_status"`
	Labels        string     `gorm:"column:labels;type:jsonb;not null;default:'{}'" json:"labels"`
	Status        int16      `gorm:"column:status;not null;default:1" json:"status"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"-"`
}

// TableName 指定表名。
func (Host) TableName() string {
	return "cmdb_hosts"
}
