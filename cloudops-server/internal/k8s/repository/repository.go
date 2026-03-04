package repository

import "gorm.io/gorm"

// Repositories 聚合k8s模块仓储。
type Repositories struct {
	Cluster *ClusterRepository
}

// New 创建仓储集合。
func New(db *gorm.DB) *Repositories {
	return &Repositories{
		Cluster: NewClusterRepository(db),
	}
}
