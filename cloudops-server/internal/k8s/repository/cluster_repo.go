package repository

import (
	"github.com/Zara1024/OpsNexus/cloudops-server/internal/k8s/model"
	"gorm.io/gorm"
)

// ClusterRepository 管理集群数据。
type ClusterRepository struct {
	db *gorm.DB
}

// NewClusterRepository 创建仓储。
func NewClusterRepository(db *gorm.DB) *ClusterRepository {
	return &ClusterRepository{db: db}
}

// List 查询集群列表。
func (r *ClusterRepository) List() ([]model.Cluster, error) {
	list := make([]model.Cluster, 0)
	err := r.db.Where("deleted_at IS NULL").Order("id DESC").Find(&list).Error
	return list, err
}

// GetByID 获取集群详情。
func (r *ClusterRepository) GetByID(id int64) (*model.Cluster, error) {
	var item model.Cluster
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Create 新增集群。
func (r *ClusterRepository) Create(item *model.Cluster) error {
	return r.db.Create(item).Error
}

// Update 更新集群。
func (r *ClusterRepository) Update(id int64, updates map[string]interface{}) error {
	return r.db.Model(&model.Cluster{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error
}

// Delete 软删除集群。
func (r *ClusterRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.Cluster{}).Error
}
