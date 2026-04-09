package dao

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	cmdbModel "dodevops-api/api/cmdb/model"
	"dodevops-api/api/monitor/model"
	. "dodevops-api/pkg/db"

	"gorm.io/gorm"
)

func GetMonitorAutomationOverview() (overview model.MonitorAutomationOverview, err error) {
	sql := `
SELECT
	(SELECT COUNT(*) FROM monitor_host_alert_rule) AS host_rule_count,
	(SELECT COUNT(*) FROM monitor_db_alert_rule) AS database_rule_count,
	(SELECT COUNT(*) FROM monitor_alert_event WHERE status = 'open') AS open_event_count,
	(SELECT COUNT(*) FROM monitor_alert_event WHERE status = 'resolved') AS resolved_event_count,
	(SELECT COUNT(*) FROM monitor_db_health_snapshot WHERE available = 1) AS database_healthy_count,
	(SELECT COUNT(*) FROM monitor_db_health_snapshot) AS database_total_count,
	(SELECT COUNT(*) FROM monitor_domain WHERE status = 1) AS domain_total_count,
	(SELECT COUNT(*) FROM monitor_domain WHERE status = 1 AND is_alive = 1) AS domain_alive_count,
	(SELECT COUNT(*) FROM monitor_domain WHERE status = 1 AND ssl_days_left IS NOT NULL AND ssl_days_left <= 30) AS expiring_domain_count,
	(SELECT COUNT(*) FROM monitor_ssl_cert_deploy_log) AS deploy_log_total_count
`
	err = Db.Raw(sql).Scan(&overview).Error
	return overview, err
}

func GetMonitorHostAlertRuleList() (list []model.MonitorHostAlertRuleEntity, err error) {
	err = Db.Order("update_time DESC, id DESC").Find(&list).Error
	return list, err
}

func GetEnabledMonitorHostAlertRules() (list []model.MonitorHostAlertRuleEntity, err error) {
	err = Db.Where("status = ?", 1).Order("id ASC").Find(&list).Error
	return list, err
}

func GetMonitorHostAlertRuleByID(id uint) (*model.MonitorHostAlertRuleEntity, error) {
	var item model.MonitorHostAlertRuleEntity
	if err := Db.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func CreateMonitorHostAlertRule(item *model.MonitorHostAlertRuleEntity) error {
	now := time.Now()
	item.CreateTime = now
	item.UpdateTime = now
	return Db.Create(item).Error
}

func UpdateMonitorHostAlertRule(item *model.MonitorHostAlertRuleEntity) error {
	item.UpdateTime = time.Now()
	return Db.Model(&model.MonitorHostAlertRuleEntity{}).
		Where("id = ?", item.ID).
		Updates(map[string]interface{}{
			"name":             item.Name,
			"host_id":          item.HostID,
			"metric_key":       item.MetricKey,
			"operator":         item.Operator,
			"threshold_value":  item.ThresholdValue,
			"severity":         item.Severity,
			"status":           item.Status,
			"notify_robot_ids": item.NotifyRobotIDs,
			"remark":           item.Remark,
			"update_time":      item.UpdateTime,
		}).Error
}

func UpdateMonitorHostAlertRuleStatus(id uint, status int) error {
	return Db.Model(&model.MonitorHostAlertRuleEntity{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"update_time": time.Now(),
		}).Error
}

func UpdateMonitorHostAlertRuleLastScan(id uint, scanAt time.Time) error {
	return Db.Model(&model.MonitorHostAlertRuleEntity{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_scan_at": scanAt,
			"update_time":  time.Now(),
		}).Error
}

func DeleteMonitorHostAlertRule(id uint) error {
	return Db.Delete(&model.MonitorHostAlertRuleEntity{}, id).Error
}

func GetMonitorDBAlertRuleList() (list []model.MonitorDBAlertRuleEntity, err error) {
	err = Db.Order("update_time DESC, id DESC").Find(&list).Error
	return list, err
}

func GetEnabledMonitorDBAlertRules() (list []model.MonitorDBAlertRuleEntity, err error) {
	err = Db.Where("status = ?", 1).Order("id ASC").Find(&list).Error
	return list, err
}

func GetMonitorDBAlertRuleByID(id uint) (*model.MonitorDBAlertRuleEntity, error) {
	var item model.MonitorDBAlertRuleEntity
	if err := Db.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func CreateMonitorDBAlertRule(item *model.MonitorDBAlertRuleEntity) error {
	now := time.Now()
	item.CreateTime = now
	item.UpdateTime = now
	return Db.Create(item).Error
}

func UpdateMonitorDBAlertRule(item *model.MonitorDBAlertRuleEntity) error {
	item.UpdateTime = time.Now()
	return Db.Model(&model.MonitorDBAlertRuleEntity{}).
		Where("id = ?", item.ID).
		Updates(map[string]interface{}{
			"name":             item.Name,
			"database_id":      item.DatabaseID,
			"metric_key":       item.MetricKey,
			"operator":         item.Operator,
			"threshold_value":  item.ThresholdValue,
			"severity":         item.Severity,
			"status":           item.Status,
			"notify_robot_ids": item.NotifyRobotIDs,
			"remark":           item.Remark,
			"update_time":      item.UpdateTime,
		}).Error
}

func UpdateMonitorDBAlertRuleStatus(id uint, status int) error {
	return Db.Model(&model.MonitorDBAlertRuleEntity{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"update_time": time.Now(),
		}).Error
}

func UpdateMonitorDBAlertRuleLastScan(id uint, scanAt time.Time) error {
	return Db.Model(&model.MonitorDBAlertRuleEntity{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_scan_at": scanAt,
			"update_time":  time.Now(),
		}).Error
}

func DeleteMonitorDBAlertRule(id uint) error {
	return Db.Delete(&model.MonitorDBAlertRuleEntity{}, id).Error
}

func GetMonitorAutomationEventList(query model.MonitorAutomationEventQuery) (list []model.MonitorAlertEventEntity, total int64, err error) {
	db := Db.Model(&model.MonitorAlertEventEntity{})
	if strings.TrimSpace(query.ResourceType) != "" {
		db = db.Where("resource_type = ?", strings.TrimSpace(query.ResourceType))
	}
	if strings.TrimSpace(query.Status) != "" {
		db = db.Where("status = ?", strings.TrimSpace(query.Status))
	}
	if keyword := strings.TrimSpace(query.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where(
			"resource_name LIKE ? OR title LIKE ? OR summary LIKE ? OR detail LIKE ?",
			like, like, like, like,
		)
	}
	if err = db.Count(&total).Error; err != nil {
		return list, total, err
	}
	err = db.Order("last_triggered_at DESC, id DESC").
		Limit(query.PageSize).
		Offset((query.PageNum - 1) * query.PageSize).
		Find(&list).Error
	return list, total, err
}

func GetLatestOpenMonitorAutomationEvent(fingerprint string) (*model.MonitorAlertEventEntity, error) {
	var item model.MonitorAlertEventEntity
	err := Db.Where("fingerprint = ? AND status = ?", fingerprint, model.MonitorEventStatusOpen).
		Order("id DESC").
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func CreateMonitorAutomationEvent(item *model.MonitorAlertEventEntity) error {
	now := time.Now()
	if item.CreateTime.IsZero() {
		item.CreateTime = now
	}
	item.UpdateTime = now
	return Db.Create(item).Error
}

func UpdateMonitorAutomationEvent(item *model.MonitorAlertEventEntity) error {
	item.UpdateTime = time.Now()
	return Db.Model(&model.MonitorAlertEventEntity{}).
		Where("id = ?", item.ID).
		Updates(map[string]interface{}{
			"summary":            item.Summary,
			"detail":             item.Detail,
			"severity":           item.Severity,
			"status":             item.Status,
			"operator":           item.Operator,
			"threshold_value":    item.ThresholdValue,
			"current_value":      item.CurrentValue,
			"recovery_value":     item.RecoveryValue,
			"dedup_count":        item.DedupCount,
			"notify_robot_ids":   item.NotifyRobotIDs,
			"incident_id":        item.IncidentID,
			"webhook_log_id":     item.WebhookLogID,
			"recovery_log_id":    item.RecoveryLogID,
			"labels":             item.Labels,
			"first_triggered_at": item.FirstTriggeredAt,
			"last_triggered_at":  item.LastTriggeredAt,
			"recovered_at":       item.RecoveredAt,
			"update_time":        item.UpdateTime,
		}).Error
}

func CreateMonitorIncident(item *model.MonitorIncidentEntity) error {
	return Db.Create(item).Error
}

func UpdateMonitorIncidentResolved(id uint, solution, remark string) error {
	now := time.Now()
	return Db.Model(&model.MonitorIncidentEntity{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      3,
			"solution":    solution,
			"remark":      remark,
			"update_time": now,
		}).Error
}

func UpsertMonitorDBHealthSnapshot(item *model.MonitorDBHealthSnapshotEntity) error {
	item.UpdateTime = time.Now()
	var existing model.MonitorDBHealthSnapshotEntity
	err := Db.Where("database_id = ?", item.DatabaseID).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			item.CreateTime = time.Now()
			return Db.Create(item).Error
		}
		return err
	}
	return Db.Model(&existing).
		Updates(map[string]interface{}{
			"database_name": item.DatabaseName,
			"database_type": item.DatabaseType,
			"host":          item.Host,
			"port":          item.Port,
			"available":     item.Available,
			"latency_ms":    item.LatencyMs,
			"error_msg":     item.ErrorMsg,
			"last_check_at": item.LastCheckAt,
			"update_time":   item.UpdateTime,
		}).Error
}

func GetMonitorDBHealthSnapshotList() (list []model.MonitorDBHealthSnapshotEntity, err error) {
	err = Db.Order("last_check_at DESC, id DESC").Find(&list).Error
	return list, err
}

func GetMonitorDomainList() (list []model.MonitorDomainEntity, err error) {
	err = Db.Order("update_time DESC, id DESC").Find(&list).Error
	return list, err
}

func GetEnabledMonitorDomains() (list []model.MonitorDomainEntity, err error) {
	err = Db.Where("status = ?", 1).Order("id ASC").Find(&list).Error
	return list, err
}

func GetMonitorDomainByID(id uint) (*model.MonitorDomainEntity, error) {
	var item model.MonitorDomainEntity
	if err := Db.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func CreateMonitorDomain(item *model.MonitorDomainEntity) error {
	now := time.Now()
	item.CreateTime = now
	item.UpdateTime = now
	return Db.Create(item).Error
}

func UpdateMonitorDomain(item *model.MonitorDomainEntity) error {
	item.UpdateTime = time.Now()
	return Db.Model(&model.MonitorDomainEntity{}).
		Where("id = ?", item.ID).
		Updates(map[string]interface{}{
			"domain":      item.Domain,
			"tags":        item.Tags,
			"remark":      item.Remark,
			"status":      item.Status,
			"update_time": item.UpdateTime,
		}).Error
}

func UpdateMonitorDomainStatus(id uint, status int64) error {
	return Db.Model(&model.MonitorDomainEntity{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"update_time": time.Now(),
		}).Error
}

func UpdateMonitorDomainCheckResult(item *model.MonitorDomainEntity) error {
	item.UpdateTime = time.Now()
	return Db.Model(&model.MonitorDomainEntity{}).
		Where("id = ?", item.ID).
		Updates(map[string]interface{}{
			"is_alive":      item.IsAlive,
			"status_code":   item.StatusCode,
			"response_time": item.ResponseTime,
			"ssl_expire_at": item.SSLExpireAt,
			"ssl_days_left": item.SSLDaysLeft,
			"ssl_issuer":    item.SSLIssuer,
			"last_check_at": item.LastCheckAt,
			"error_msg":     item.ErrorMsg,
			"update_time":   item.UpdateTime,
		}).Error
}

func DeleteMonitorDomain(id uint) error {
	return Db.Delete(&model.MonitorDomainEntity{}, id).Error
}

func GetMonitorDomainScheduleList() (list []model.MonitorDomainScheduleEntity, err error) {
	err = Db.Order("id ASC").Find(&list).Error
	return list, err
}

func GetEnabledMonitorDomainSchedules() (list []model.MonitorDomainScheduleEntity, err error) {
	err = Db.Where("enabled = ?", true).Order("id ASC").Find(&list).Error
	return list, err
}

func GetMonitorDomainScheduleByID(id uint) (*model.MonitorDomainScheduleEntity, error) {
	var item model.MonitorDomainScheduleEntity
	if err := Db.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func UpdateMonitorDomainSchedule(item *model.MonitorDomainScheduleEntity) error {
	item.UpdateTime = time.Now()
	return Db.Model(&model.MonitorDomainScheduleEntity{}).
		Where("id = ?", item.ID).
		Updates(map[string]interface{}{
			"enabled":             item.Enabled,
			"cron_expr":           item.CronExpr,
			"next_run_at":         item.NextRunAt,
			"last_run_at":         item.LastRunAt,
			"status":              item.Status,
			"notify_enabled":      item.NotifyEnabled,
			"notify_robot_id":     item.NotifyRobotID,
			"expire_alert_days":   item.ExpireAlertDays,
			"scan_timeout_ms":     item.ScanTimeoutMs,
			"auto_renew_enabled":  item.AutoRenewEnabled,
			"auto_deploy_enabled": item.AutoDeployEnabled,
			"deploy_host_id":      item.DeployHostID,
			"deploy_path":         item.DeployPath,
			"reload_command":      item.ReloadCommand,
			"update_time":         item.UpdateTime,
		}).Error
}

func UpdateMonitorDomainScheduleRuntime(id uint, status string, lastRunAt, nextRunAt *time.Time) error {
	updates := map[string]interface{}{
		"status":      status,
		"update_time": time.Now(),
	}
	if lastRunAt != nil {
		updates["last_run_at"] = *lastRunAt
	}
	if nextRunAt != nil {
		updates["next_run_at"] = *nextRunAt
	}
	return Db.Model(&model.MonitorDomainScheduleEntity{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func GetMonitorSSLCertList() (list []model.MonitorSSLCertEntity, err error) {
	err = Db.Order("expire_time ASC, id DESC").Find(&list).Error
	return list, err
}

func GetMonitorSSLCertByID(id uint) (*model.MonitorSSLCertEntity, error) {
	var item model.MonitorSSLCertEntity
	if err := Db.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func GetMonitorSSLCertDeployLogList(limit int) (list []model.MonitorSSLCertDeployLogEntity, err error) {
	db := Db.Order("create_time DESC, id DESC")
	if limit > 0 {
		db = db.Limit(limit)
	}
	err = db.Find(&list).Error
	return list, err
}

func CreateMonitorSSLCertDeployLog(item *model.MonitorSSLCertDeployLogEntity) error {
	now := time.Now()
	item.CreateTime = now
	item.UpdateTime = now
	return Db.Create(item).Error
}

func UpdateMonitorSSLCertDeployLog(item *model.MonitorSSLCertDeployLogEntity) error {
	item.UpdateTime = time.Now()
	return Db.Model(&model.MonitorSSLCertDeployLogEntity{}).
		Where("id = ?", item.ID).
		Updates(map[string]interface{}{
			"status":       item.Status,
			"backup_files": item.BackupFiles,
			"deploy_files": item.DeployFiles,
			"logs":         item.Logs,
			"error_msg":    item.ErrorMsg,
			"update_time":  item.UpdateTime,
		}).Error
}

func GetCmdbHostsByIDs(ids []uint) (list []cmdbModel.CmdbHost, err error) {
	if len(ids) == 0 {
		return []cmdbModel.CmdbHost{}, nil
	}
	err = Db.Where("id IN ?", ids).Find(&list).Error
	return list, err
}

func EncodeUintList(values []uint) string {
	if len(values) == 0 {
		return "[]"
	}
	data, err := json.Marshal(values)
	if err != nil {
		return "[]"
	}
	return string(data)
}
