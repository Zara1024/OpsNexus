package service

import (
	"encoding/json"
	"fmt"
	"time"

	"dodevops-api/api/configcenter/dao"
	"dodevops-api/api/configcenter/model"
	"dodevops-api/common/util"
	"github.com/robfig/cron/v3"
)

type SyncScheduleService struct {
	dao *dao.SyncScheduleDao
}

func NewSyncScheduleService() *SyncScheduleService {
	return &SyncScheduleService{
		dao: dao.NewSyncScheduleDao(),
	}
}

// Create creates a sync schedule after validating editable fields.
func (s *SyncScheduleService) Create(syncSchedule *model.SyncSchedule) error {
	if err := s.validateCronExpr(syncSchedule.CronExpr); err != nil {
		return fmt.Errorf("invalid cron expression: %v", err)
	}

	if err := s.validateKeyTypes(syncSchedule.KeyTypes); err != nil {
		return fmt.Errorf("invalid keyTypes format: %v", err)
	}

	nextTime, err := s.calculateNextRunTime(syncSchedule.CronExpr)
	if err != nil {
		return fmt.Errorf("failed to calculate next run time: %v", err)
	}
	syncSchedule.NextRunTime = nextTime

	requestedStatus := syncSchedule.Status
	if err := s.dao.Create(syncSchedule); err != nil {
		return err
	}
	if requestedStatus == 0 && syncSchedule.Status != 0 {
		if err := s.dao.UpdateStatus(syncSchedule.ID, 0); err != nil {
			return err
		}
		syncSchedule.Status = 0
	}
	return nil
}

// Update updates only editable fields and returns the full persisted record.
func (s *SyncScheduleService) Update(syncSchedule *model.SyncSchedule) error {
	if _, err := s.dao.GetByID(syncSchedule.ID); err != nil {
		return err
	}

	if err := s.validateCronExpr(syncSchedule.CronExpr); err != nil {
		return fmt.Errorf("invalid cron expression: %v", err)
	}

	if err := s.validateKeyTypes(syncSchedule.KeyTypes); err != nil {
		return fmt.Errorf("invalid keyTypes format: %v", err)
	}

	nextTime, err := s.calculateNextRunTime(syncSchedule.CronExpr)
	if err != nil {
		return fmt.Errorf("failed to calculate next run time: %v", err)
	}
	syncSchedule.NextRunTime = nextTime

	if err := s.dao.Update(syncSchedule); err != nil {
		return err
	}

	updated, err := s.dao.GetByID(syncSchedule.ID)
	if err != nil {
		return err
	}

	*syncSchedule = *updated
	return nil
}

// Delete removes a sync schedule.
func (s *SyncScheduleService) Delete(id uint) error {
	return s.dao.Delete(id)
}

// GetByID fetches a sync schedule by id.
func (s *SyncScheduleService) GetByID(id uint) (*model.SyncSchedule, error) {
	return s.dao.GetByID(id)
}

// ListWithPage lists sync schedules with pagination.
func (s *SyncScheduleService) ListWithPage(page, pageSize int) ([]model.SyncSchedule, int64, error) {
	return s.dao.ListWithPage(page, pageSize)
}

// GetActiveSchedules lists enabled sync schedules.
func (s *SyncScheduleService) GetActiveSchedules() ([]model.SyncSchedule, error) {
	return s.dao.GetByStatus(1)
}

// UpdateLastRunTime updates the last and next execution times.
func (s *SyncScheduleService) UpdateLastRunTime(id uint, lastRunTime time.Time) error {
	schedule, err := s.dao.GetByID(id)
	if err != nil {
		return err
	}

	nextTime, err := s.calculateNextRunTime(schedule.CronExpr)
	if err != nil {
		return fmt.Errorf("failed to calculate next run time: %v", err)
	}

	lastTime := util.HTime{Time: lastRunTime}
	return s.dao.UpdateRunTime(id, &lastTime, nextTime)
}

// UpdateLastRunTimeAndLog updates execution times and sync log together.
func (s *SyncScheduleService) UpdateLastRunTimeAndLog(id uint, lastRunTime time.Time, syncLog string) error {
	schedule, err := s.dao.GetByID(id)
	if err != nil {
		return err
	}

	nextTime, err := s.calculateNextRunTime(schedule.CronExpr)
	if err != nil {
		return fmt.Errorf("failed to calculate next run time: %v", err)
	}

	lastTime := util.HTime{Time: lastRunTime}
	return s.dao.UpdateRunTimeAndLog(id, &lastTime, nextTime, syncLog)
}

// ToggleStatus flips a schedule between enabled and disabled.
func (s *SyncScheduleService) ToggleStatus(id uint, status int) error {
	return s.dao.UpdateStatus(id, status)
}

func (s *SyncScheduleService) validateCronExpr(cronExpr string) error {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(cronExpr)
	return err
}

func (s *SyncScheduleService) validateKeyTypes(keyTypes string) error {
	var types []int
	return json.Unmarshal([]byte(keyTypes), &types)
}

func (s *SyncScheduleService) calculateNextRunTime(cronExpr string) (*util.HTime, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronExpr)
	if err != nil {
		return nil, err
	}

	nextTime := schedule.Next(time.Now())
	hTime := util.HTime{Time: nextTime}
	return &hTime, nil
}

// ParseKeyTypes parses the JSON-encoded cloud vendor list.
func (s *SyncScheduleService) ParseKeyTypes(keyTypes string) ([]int, error) {
	var types []int
	err := json.Unmarshal([]byte(keyTypes), &types)
	return types, err
}
