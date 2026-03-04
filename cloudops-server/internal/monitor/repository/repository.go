package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/monitor/model"
)

// Repositories 聚合 monitor 模块仓储。
type Repositories struct {
	AlertRule *AlertRuleRepository
	Alert     *AlertEventRepository
	Channel   *ChannelRepository
}

// New 创建仓储集合。
func New(db *gorm.DB) *Repositories {
	return &Repositories{
		AlertRule: NewAlertRuleRepository(db),
		Alert:     NewAlertEventRepository(db),
		Channel:   NewChannelRepository(db),
	}
}

// AlertRuleRepository 告警规则仓储。
type AlertRuleRepository struct{ db *gorm.DB }

func NewAlertRuleRepository(db *gorm.DB) *AlertRuleRepository { return &AlertRuleRepository{db: db} }

func (r *AlertRuleRepository) List() ([]model.AlertRule, error) {
	list := make([]model.AlertRule, 0)
	err := r.db.Where("deleted_at IS NULL").Order("id DESC").Find(&list).Error
	return list, err
}

func (r *AlertRuleRepository) Create(item *model.AlertRule) error { return r.db.Create(item).Error }

func (r *AlertRuleRepository) Update(id int64, updates map[string]interface{}) error {
	return r.db.Model(&model.AlertRule{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error
}

func (r *AlertRuleRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.AlertRule{}).Error
}

func (r *AlertRuleRepository) GetByID(id int64) (*model.AlertRule, error) {
	var item model.AlertRule
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// AlertEventRepository 告警事件仓储。
type AlertEventRepository struct{ db *gorm.DB }

func NewAlertEventRepository(db *gorm.DB) *AlertEventRepository { return &AlertEventRepository{db: db} }

func (r *AlertEventRepository) ListActive() ([]model.AlertEvent, error) {
	list := make([]model.AlertEvent, 0)
	err := r.db.Where("deleted_at IS NULL AND status = ?", "firing").Order("starts_at DESC").Find(&list).Error
	return list, err
}

func (r *AlertEventRepository) ListHistory() ([]model.AlertEvent, error) {
	list := make([]model.AlertEvent, 0)
	err := r.db.Where("deleted_at IS NULL").Order("starts_at DESC").Limit(500).Find(&list).Error
	return list, err
}

func (r *AlertEventRepository) Ack(id, userID int64) error {
	return r.db.Model(&model.AlertEvent{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{"acknowledged_by": userID, "updated_at": time.Now()}).Error
}

func (r *AlertEventRepository) Silence(id int64, until time.Time) error {
	return r.db.Model(&model.AlertEvent{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{"silenced_until": until, "updated_at": time.Now()}).Error
}

// ChannelRepository 通知渠道仓储。
type ChannelRepository struct{ db *gorm.DB }

func NewChannelRepository(db *gorm.DB) *ChannelRepository { return &ChannelRepository{db: db} }

func (r *ChannelRepository) List() ([]model.NotificationChannel, error) {
	list := make([]model.NotificationChannel, 0)
	err := r.db.Where("deleted_at IS NULL").Order("id DESC").Find(&list).Error
	return list, err
}

func (r *ChannelRepository) Create(item *model.NotificationChannel) error {
	return r.db.Create(item).Error
}

func (r *ChannelRepository) GetByID(id int64) (*model.NotificationChannel, error) {
	var item model.NotificationChannel
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
