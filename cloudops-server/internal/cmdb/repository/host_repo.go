package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/model"
)

// HostRepository 定义主机数据访问能力。
type HostRepository struct {
	db *gorm.DB
}

// NewHostRepository 创建主机仓储。
func NewHostRepository(db *gorm.DB) *HostRepository {
	return &HostRepository{db: db}
}

// HostListFilter 定义主机分页查询条件。
type HostListFilter struct {
	Page     int
	PageSize int
	GroupID  int64
	Status   int16
	Keyword  string
	InnerIP  string
	OuterIP  string
}

// Create 创建主机。
func (r *HostRepository) Create(item *model.Host) error {
	return r.db.Create(item).Error
}

// Update 更新主机。
func (r *HostRepository) Update(id int64, updates map[string]interface{}) error {
	return r.db.Model(&model.Host{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error
}

// SoftDelete 软删除主机。
func (r *HostRepository) SoftDelete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.Host{}).Error
}

// GetByID 根据 ID 查询主机。
func (r *HostRepository) GetByID(id int64) (*model.Host, error) {
	var item model.Host
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// List 分页查询主机。
func (r *HostRepository) List(filter *HostListFilter) ([]model.Host, int64, error) {
	list := make([]model.Host, 0)
	var total int64

	query := r.db.Model(&model.Host{}).Where("deleted_at IS NULL")
	if filter.GroupID > 0 {
		query = query.Where("group_id = ?", filter.GroupID)
	}
	if filter.Status != 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Keyword != "" {
		kw := fmt.Sprintf("%%%s%%", strings.TrimSpace(filter.Keyword))
		query = query.Where("hostname ILIKE ? OR inner_ip ILIKE ? OR outer_ip ILIKE ?", kw, kw, kw)
	}
	if filter.InnerIP != "" {
		query = query.Where("inner_ip ILIKE ?", fmt.Sprintf("%%%s%%", strings.TrimSpace(filter.InnerIP)))
	}
	if filter.OuterIP != "" {
		query = query.Where("outer_ip ILIKE ?", fmt.Sprintf("%%%s%%", strings.TrimSpace(filter.OuterIP)))
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Order("id DESC").Offset(offset).Limit(filter.PageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// BatchDelete 批量删除主机。
func (r *HostRepository) BatchDelete(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Where("id IN ?", ids).Delete(&model.Host{}).Error
}

// BatchUpdateStatus 批量更新状态。
func (r *HostRepository) BatchUpdateStatus(ids []int64, status int16) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Model(&model.Host{}).Where("id IN ? AND deleted_at IS NULL", ids).Updates(map[string]interface{}{"status": status}).Error
}
