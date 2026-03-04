package repository

import (
	"gorm.io/gorm"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/model"
)

// DepartmentRepository 定义部门数据访问能力。
type DepartmentRepository struct {
	db *gorm.DB
}

// NewDepartmentRepository 创建部门仓储。
func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

// Create 创建部门。
func (r *DepartmentRepository) Create(dept *model.Department) error {
	return r.db.Create(dept).Error
}

// Update 更新部门。
func (r *DepartmentRepository) Update(id int64, updates map[string]interface{}) error {
	return r.db.Model(&model.Department{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error
}

// SoftDelete 软删除部门。
func (r *DepartmentRepository) SoftDelete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.Department{}).Error
}

// List 查询部门列表。
func (r *DepartmentRepository) List() ([]model.Department, error) {
	depts := make([]model.Department, 0)
	if err := r.db.Where("deleted_at IS NULL").Order("sort_order ASC, id ASC").Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}
