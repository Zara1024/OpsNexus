package repository

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/audit/model"
)

// Repository 审计日志仓储。
type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

// LogQuery 审计日志查询条件。
type LogQuery struct {
	Page      int
	PageSize  int
	Keyword   string
	Module    string
	Action    string
	StartTime string
	EndTime   string
}

func (r *Repository) Create(ctx context.Context, item *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *Repository) List(query LogQuery) ([]model.AuditLog, int64, error) {
	db := r.buildQuery(query)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	if query.PageSize > 200 {
		query.PageSize = 200
	}

	list := make([]model.AuditLog, 0)
	err := db.Order("id DESC").Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize).Find(&list).Error
	return list, total, err
}

func (r *Repository) ListForExport(query LogQuery) ([]model.AuditLog, error) {
	db := r.buildQuery(query)
	list := make([]model.AuditLog, 0)
	err := db.Order("id DESC").Limit(5000).Find(&list).Error
	return list, err
}

func (r *Repository) ListLoginLogs(query LogQuery) ([]model.AuditLog, int64, error) {
	query.Module = "auth"
	if strings.TrimSpace(query.Action) == "" {
		query.Action = "create"
	}
	return r.List(query)
}

func (r *Repository) buildQuery(query LogQuery) *gorm.DB {
	db := r.db.Model(&model.AuditLog{}).Where("deleted_at IS NULL")

	if keyword := strings.TrimSpace(query.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("username ILIKE ? OR path ILIKE ? OR description ILIKE ?", like, like, like)
	}
	if module := strings.TrimSpace(query.Module); module != "" {
		db = db.Where("module = ?", module)
	}
	if action := strings.TrimSpace(query.Action); action != "" {
		db = db.Where("action = ?", action)
	}
	if start := strings.TrimSpace(query.StartTime); start != "" {
		if t, err := time.Parse(time.RFC3339, start); err == nil {
			db = db.Where("created_at >= ?", t)
		}
	}
	if end := strings.TrimSpace(query.EndTime); end != "" {
		if t, err := time.Parse(time.RFC3339, end); err == nil {
			db = db.Where("created_at <= ?", t)
		}
	}
	return db
}
