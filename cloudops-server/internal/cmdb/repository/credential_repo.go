package repository

import (
	"gorm.io/gorm"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/model"
)

// CredentialRepository 定义凭据数据访问能力。
type CredentialRepository struct {
	db *gorm.DB
}

// NewCredentialRepository 创建凭据仓储。
func NewCredentialRepository(db *gorm.DB) *CredentialRepository {
	return &CredentialRepository{db: db}
}

// Create 创建凭据。
func (r *CredentialRepository) Create(item *model.Credential) error {
	return r.db.Create(item).Error
}

// Update 更新凭据。
func (r *CredentialRepository) Update(id int64, updates map[string]interface{}) error {
	return r.db.Model(&model.Credential{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error
}

// SoftDelete 软删除凭据。
func (r *CredentialRepository) SoftDelete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.Credential{}).Error
}

// GetByID 查询凭据。
func (r *CredentialRepository) GetByID(id int64) (*model.Credential, error) {
	var item model.Credential
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// List 查询凭据列表。
func (r *CredentialRepository) List() ([]model.Credential, error) {
	list := make([]model.Credential, 0)
	if err := r.db.Where("deleted_at IS NULL").Order("id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
