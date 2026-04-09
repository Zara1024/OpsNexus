package dao

import (
	"time"

	"dodevops-api/api/configcenter/model"
	"dodevops-api/common"
	"dodevops-api/common/util"
	"gorm.io/gorm"
)

type SyncScheduleDao struct {
	db *gorm.DB
}

func NewSyncScheduleDao() *SyncScheduleDao {
	return &SyncScheduleDao{
		db: common.GetDB(),
	}
}

// Create creates a sync schedule record.
func (d *SyncScheduleDao) Create(syncSchedule *model.SyncSchedule) error {
	return d.db.Select("*").Create(syncSchedule).Error
}

// Update persists only editable fields so immutable/runtime fields are preserved.
func (d *SyncScheduleDao) Update(syncSchedule *model.SyncSchedule) error {
	syncSchedule.UpdatedAt = util.HTime{Time: time.Now()}

	return d.db.Model(&model.SyncSchedule{}).Where("id = ?", syncSchedule.ID).Updates(map[string]interface{}{
		"name":          syncSchedule.Name,
		"cron_expr":     syncSchedule.CronExpr,
		"key_types":     syncSchedule.KeyTypes,
		"status":        syncSchedule.Status,
		"next_run_time": syncSchedule.NextRunTime,
		"remark":        syncSchedule.Remark,
		"updated_at":    syncSchedule.UpdatedAt,
	}).Error
}

// Delete removes a sync schedule record.
func (d *SyncScheduleDao) Delete(id uint) error {
	return d.db.Delete(&model.SyncSchedule{}, id).Error
}

// GetByID fetches a sync schedule by id.
func (d *SyncScheduleDao) GetByID(id uint) (*model.SyncSchedule, error) {
	var syncSchedule model.SyncSchedule
	err := d.db.First(&syncSchedule, id).Error
	if err != nil {
		return nil, err
	}
	return &syncSchedule, nil
}

// ListWithPage returns sync schedules with pagination.
func (d *SyncScheduleDao) ListWithPage(page, pageSize int) ([]model.SyncSchedule, int64, error) {
	var syncSchedules []model.SyncSchedule
	var total int64

	offset := (page - 1) * pageSize

	if err := d.db.Model(&model.SyncSchedule{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := d.db.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&syncSchedules).Error
	if err != nil {
		return nil, 0, err
	}

	return syncSchedules, total, nil
}

// GetByStatus fetches sync schedules by status.
func (d *SyncScheduleDao) GetByStatus(status int) ([]model.SyncSchedule, error) {
	var syncSchedules []model.SyncSchedule
	err := d.db.Where("status = ?", status).Find(&syncSchedules).Error
	return syncSchedules, err
}

// UpdateStatus updates only the status field.
func (d *SyncScheduleDao) UpdateStatus(id uint, status int) error {
	return d.db.Model(&model.SyncSchedule{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateRunTime updates execution times.
func (d *SyncScheduleDao) UpdateRunTime(id uint, lastRunTime, nextRunTime *util.HTime) error {
	updates := map[string]interface{}{
		"last_run_time": lastRunTime,
		"next_run_time": nextRunTime,
	}
	return d.db.Model(&model.SyncSchedule{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateRunTimeAndLog updates execution times and sync log.
func (d *SyncScheduleDao) UpdateRunTimeAndLog(id uint, lastRunTime, nextRunTime *util.HTime, syncLog string) error {
	updates := map[string]interface{}{
		"last_run_time": lastRunTime,
		"next_run_time": nextRunTime,
		"sync_log":      syncLog,
	}
	return d.db.Model(&model.SyncSchedule{}).Where("id = ?", id).Updates(updates).Error
}

// GetAll fetches all sync schedules.
func (d *SyncScheduleDao) GetAll() ([]model.SyncSchedule, error) {
	var syncSchedules []model.SyncSchedule
	err := d.db.Find(&syncSchedules).Error
	return syncSchedules, err
}
