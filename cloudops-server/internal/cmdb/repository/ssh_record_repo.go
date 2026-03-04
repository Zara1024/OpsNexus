package repository

import (
	"gorm.io/gorm"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/model"
)

// SSHRecordRepository 定义 SSH 审计记录数据访问能力。
type SSHRecordRepository struct {
	db *gorm.DB
}

// NewSSHRecordRepository 创建 SSH 审计记录仓储。
func NewSSHRecordRepository(db *gorm.DB) *SSHRecordRepository {
	return &SSHRecordRepository{db: db}
}

// Create 创建 SSH 审计记录。
func (r *SSHRecordRepository) Create(item *model.SSHRecord) error {
	return r.db.Create(item).Error
}

// EndSession 结束 SSH 会话。
func (r *SSHRecordRepository) EndSession(sessionID string) error {
	return r.db.Model(&model.SSHRecord{}).
		Where("session_id = ? AND deleted_at IS NULL", sessionID).
		Update("end_time", gorm.Expr("NOW()")).Error
}

// List 查询 SSH 审计记录列表。
func (r *SSHRecordRepository) List(page, pageSize int) ([]model.SSHRecord, int64, error) {
	list := make([]model.SSHRecord, 0)
	var total int64
	query := r.db.Model(&model.SSHRecord{}).Where("deleted_at IS NULL")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
