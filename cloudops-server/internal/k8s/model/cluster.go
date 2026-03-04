package model

import "time"

// Cluster 定义K8s集群信息。
type Cluster struct {
	ID                  int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name                string     `gorm:"column:name;size:128;not null" json:"name"`
	APIServerURL        string     `gorm:"column:api_server_url;size:512;not null" json:"api_server_url"`
	KubeconfigEncrypted string     `gorm:"column:kubeconfig_encrypted;type:text;not null" json:"-"`
	Description         string     `gorm:"column:description;size:512;not null;default:''" json:"description"`
	Version             string     `gorm:"column:version;size:64;not null;default:''" json:"version"`
	Status              int16      `gorm:"column:status;not null;default:1" json:"status"`
	NodeCount           int        `gorm:"column:node_count;not null;default:0" json:"node_count"`
	PodCount            int        `gorm:"column:pod_count;not null;default:0" json:"pod_count"`
	CreatedBy           int64      `gorm:"column:created_by;not null;default:0" json:"created_by"`
	CreatedAt           time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt           *time.Time `gorm:"column:deleted_at" json:"-"`
}

// TableName 指定表名。
func (Cluster) TableName() string {
	return "k8s_clusters"
}
