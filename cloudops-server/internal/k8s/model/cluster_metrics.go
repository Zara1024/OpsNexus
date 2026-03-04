package model

import "time"

// ClusterMetrics 定义集群指标缓存。
type ClusterMetrics struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ClusterID int64      `gorm:"column:cluster_id;not null;index" json:"cluster_id"`
	CPUUsage  float64    `gorm:"column:cpu_usage;type:numeric(8,2);not null;default:0" json:"cpu_usage"`
	MemUsage  float64    `gorm:"column:mem_usage;type:numeric(8,2);not null;default:0" json:"mem_usage"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`
}

// TableName 指定表名。
func (ClusterMetrics) TableName() string {
	return "k8s_cluster_metrics"
}
