package repository

import (
	"gorm.io/gorm"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/model"
)

// GroupRepository 定义分组数据访问能力。
type GroupRepository struct {
	db *gorm.DB
}

// NewGroupRepository 创建分组仓储。
func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

// Create 创建分组。
func (r *GroupRepository) Create(item *model.Group) error {
	return r.db.Create(item).Error
}

// Update 更新分组。
func (r *GroupRepository) Update(id int64, updates map[string]interface{}) error {
	return r.db.Model(&model.Group{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error
}

// SoftDelete 软删除分组。
func (r *GroupRepository) SoftDelete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.Group{}).Error
}

// List 查询全部分组。
func (r *GroupRepository) List() ([]model.Group, error) {
	list := make([]model.Group, 0)
	if err := r.db.Where("deleted_at IS NULL").Order("sort_order ASC, id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
