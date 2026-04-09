package dao

import (
	"errors"
	"time"

	"dodevops-api/api/monitor/model"
	. "dodevops-api/pkg/db"

	"gorm.io/gorm"
)

// GetMonitorAlertSummary returns the cross-table summary for the alert center.
func GetMonitorAlertSummary() (summary model.MonitorAlertSummary, err error) {
	sql := `
SELECT
    (SELECT COUNT(*) FROM monitor_incident) AS total_incidents,
    (SELECT COUNT(*) FROM monitor_incident WHERE status = 1) AS open_incidents,
    (SELECT COUNT(*) FROM monitor_incident WHERE status = 2) AS processing_incidents,
    (SELECT COUNT(*) FROM monitor_incident WHERE status = 3) AS resolved_incidents,
    (SELECT COUNT(*) FROM monitor_webhook_log) AS total_webhook_logs,
    (SELECT COUNT(*) FROM monitor_webhook_log WHERE LOWER(COALESCE(level, '')) IN ('critical', 'p1', 'p2')) AS critical_webhook_logs,
    (SELECT COUNT(*) FROM monitor_webhook_notify_log) AS total_notify_logs,
    (SELECT COUNT(*) FROM monitor_webhook_notify_log WHERE status = 'success') AS successful_notify_logs,
    (SELECT COUNT(*) FROM monitor_webhook_notify_log WHERE status <> 'success') AS failed_notify_logs,
    (SELECT COUNT(*) FROM monitor_notify_robot) AS total_notify_robots,
    (SELECT COUNT(*) FROM monitor_notify_robot WHERE status = 1) AS enabled_notify_robots,
    (SELECT COUNT(*) FROM monitor_alert_source) AS total_alert_sources,
    (SELECT COUNT(*) FROM monitor_alert_source WHERE status = 1) AS enabled_alert_sources,
    COALESCE((SELECT DATE_FORMAT(MAX(created_at), '%Y-%m-%d %H:%i:%s') FROM monitor_webhook_log), '') AS latest_alert_time,
    COALESCE((SELECT DATE_FORMAT(MAX(created_at), '%Y-%m-%d %H:%i:%s') FROM monitor_webhook_notify_log), '') AS latest_notify_time
`

	err = Db.Raw(sql).Scan(&summary).Error
	return summary, err
}

// GetMonitorIncidentList returns a paginated incident list.
func GetMonitorIncidentList(query model.MonitorIncidentQuery) (list []model.MonitorIncident, count int64, err error) {
	db := Db.Table("monitor_incident")

	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		db = db.Where(
			"COALESCE(alert_desc, '') LIKE ? OR COALESCE(business_line, '') LIKE ? OR COALESCE(handler, '') LIKE ?",
			like, like, like,
		)
	}
	if query.Status >= 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.Level != "" {
		db = db.Where("alert_level = ?", query.Level)
	}
	if query.Namespace != "" {
		like := "%" + query.Namespace + "%"
		db = db.Where(
			"COALESCE(alert_desc, '') LIKE ? OR COALESCE(detail_url, '') LIKE ? OR COALESCE(remark, '') LIKE ?",
			like, like, like,
		)
	}
	if query.WorkloadName != "" {
		like := "%" + query.WorkloadName + "%"
		db = db.Where(
			"COALESCE(alert_desc, '') LIKE ? OR COALESCE(detail_url, '') LIKE ? OR COALESCE(remark, '') LIKE ?",
			like, like, like,
		)
	}

	if err = db.Count(&count).Error; err != nil {
		return list, count, err
	}

	err = db.Select(`
		id,
		COALESCE(DATE_FORMAT(alert_time, '%Y-%m-%d %H:%i:%s'), '') AS alert_time,
		COALESCE(business_line, '') AS business_line,
		COALESCE(frequency, '') AS frequency,
		COALESCE(alert_desc, '') AS alert_desc,
		COALESCE(alert_level, '') AS alert_level,
		COALESCE(incident_cause, '') AS incident_cause,
		COALESCE(department, '') AS department,
		COALESCE(solution, '') AS solution,
		COALESCE(detail_url, '') AS detail_url,
		COALESCE(handler, '') AS handler,
		COALESCE(handler_id, 0) AS handler_id,
		COALESCE(status, 0) AS status,
		COALESCE(remark, '') AS remark,
		COALESCE(DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s'), '') AS create_time,
		COALESCE(DATE_FORMAT(update_time, '%Y-%m-%d %H:%i:%s'), '') AS update_time
	`).
		Order("create_time DESC, id DESC").
		Limit(query.PageSize).
		Offset((query.PageNum - 1) * query.PageSize).
		Scan(&list).Error

	return list, count, err
}

// GetMonitorWebhookLogList returns a paginated webhook alert log list.
func GetMonitorWebhookLogList(query model.MonitorWebhookLogQuery) (list []model.MonitorWebhookLog, count int64, err error) {
	db := Db.Table("monitor_webhook_log")

	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		db = db.Where(
			"COALESCE(title, '') LIKE ? OR COALESCE(content, '') LIKE ? OR COALESCE(source, '') LIKE ? OR COALESCE(error_msg, '') LIKE ?",
			like, like, like, like,
		)
	}
	if query.Source != "" {
		db = db.Where("source = ?", query.Source)
	}
	if query.Level != "" {
		db = db.Where("level = ?", query.Level)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if err = db.Count(&count).Error; err != nil {
		return list, count, err
	}

	err = db.Select(`
		id,
		COALESCE(source, '') AS source,
		COALESCE(title, '') AS title,
		COALESCE(content, '') AS content,
		COALESCE(level, '') AS level,
		COALESCE(tags, '') AS tags,
		COALESCE(extra, '') AS extra,
		COALESCE(notify_robot_ids, '') AS notify_robot_ids,
		COALESCE(status, '') AS status,
		COALESCE(error_msg, '') AS error_msg,
		COALESCE(notify_count, 0) AS notify_count,
		COALESCE(success_count, 0) AS success_count,
		COALESCE(failed_count, 0) AS failed_count,
		COALESCE(DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), '') AS created_at
	`).
		Order("created_at DESC, id DESC").
		Limit(query.PageSize).
		Offset((query.PageNum - 1) * query.PageSize).
		Scan(&list).Error

	return list, count, err
}

// GetMonitorWebhookNotifyLogList returns a paginated notify log list.
func GetMonitorWebhookNotifyLogList(query model.MonitorNotifyLogQuery) (list []model.MonitorWebhookNotifyLog, count int64, err error) {
	db := Db.Table("monitor_webhook_notify_log mwl").
		Joins("LEFT JOIN monitor_webhook_log mw ON mw.id = mwl.webhook_log_id")

	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		db = db.Where(
			"COALESCE(mw.title, '') LIKE ? OR COALESCE(mw.source, '') LIKE ? OR COALESCE(mwl.robot_name, '') LIKE ? OR COALESCE(mwl.error_msg, '') LIKE ?",
			like, like, like, like,
		)
	}
	if query.Status != "" {
		db = db.Where("mwl.status = ?", query.Status)
	}
	if query.RobotType != "" {
		db = db.Where("mwl.robot_type = ?", query.RobotType)
	}

	if err = db.Count(&count).Error; err != nil {
		return list, count, err
	}

	err = db.Select(`
		mwl.id,
		COALESCE(mwl.webhook_log_id, 0) AS webhook_log_id,
		COALESCE(mwl.robot_id, 0) AS robot_id,
		COALESCE(mwl.robot_name, '') AS robot_name,
		COALESCE(mwl.robot_type, '') AS robot_type,
		COALESCE(mwl.status, '') AS status,
		COALESCE(mwl.error_msg, '') AS error_msg,
		COALESCE(DATE_FORMAT(mwl.created_at, '%Y-%m-%d %H:%i:%s'), '') AS created_at,
		COALESCE(mw.title, '') AS alert_title,
		COALESCE(mw.source, '') AS alert_source,
		COALESCE(mw.level, '') AS alert_level
	`).
		Order("mwl.created_at DESC, mwl.id DESC").
		Limit(query.PageSize).
		Offset((query.PageNum - 1) * query.PageSize).
		Scan(&list).Error

	return list, count, err
}

// GetMonitorNotifyRobotList returns the configured notify robots.
func GetMonitorNotifyRobotList() (list []model.MonitorNotifyRobot, err error) {
	err = Db.Table("monitor_notify_robot").
		Select(`
			id,
			COALESCE(name, '') AS name,
			COALESCE(type, '') AS type,
			COALESCE(webhook, '') AS webhook,
			COALESCE(secret, '') AS secret,
			COALESCE(server, '') AS server,
			COALESCE(port, 0) AS port,
			COALESCE(username, '') AS username,
			COALESCE(password, '') AS password,
			COALESCE(nickname, '') AS nickname,
			COALESCE(headers, '') AS headers,
			COALESCE(method, '') AS method,
			COALESCE(template, '') AS template,
			COALESCE(status, 0) AS status,
			COALESCE(remark, '') AS remark,
			COALESCE(DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), '') AS created_at,
			COALESCE(DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s'), '') AS updated_at
		`).
		Order("updated_at DESC, id DESC").
		Scan(&list).Error

	return list, err
}

// GetMonitorAlertSourceList returns the configured alert sources.
func GetMonitorAlertSourceList() (list []model.MonitorAlertSource, err error) {
	err = Db.Table("monitor_alert_source").
		Select(`
			id,
			COALESCE(name, '') AS name,
			COALESCE(type, 0) AS type,
			COALESCE(app_key, '') AS app_key,
			COALESCE(api_base_url, '') AS api_base_url,
			COALESCE(status, 0) AS status,
			COALESCE(remark, '') AS remark,
			COALESCE(create_time, '') AS create_time,
			COALESCE(update_time, '') AS update_time,
			COALESCE(key_id, 0) AS key_id,
			COALESCE(host_id, 0) AS host_id
		`).
		Order("update_time DESC, id DESC").
		Scan(&list).Error

	return list, err
}

// GetMonitorAlertSourceByID returns one alert source by id.
func GetMonitorAlertSourceByID(id uint) (source *model.MonitorAlertSource, err error) {
	source = &model.MonitorAlertSource{}
	err = Db.Table("monitor_alert_source").
		Select(`
			id,
			COALESCE(name, '') AS name,
			COALESCE(type, 0) AS type,
			COALESCE(app_key, '') AS app_key,
			COALESCE(api_base_url, '') AS api_base_url,
			COALESCE(status, 0) AS status,
			COALESCE(remark, '') AS remark,
			COALESCE(create_time, '') AS create_time,
			COALESCE(update_time, '') AS update_time,
			COALESCE(key_id, 0) AS key_id,
			COALESCE(host_id, 0) AS host_id
		`).
		Where("id = ?", id).
		Take(source).Error
	return source, err
}

// GetMonitorAlertSourcesByType returns alert sources filtered by type and status.
func GetMonitorAlertSourcesByType(sourceType int, enabledOnly bool) (list []model.MonitorAlertSource, err error) {
	db := Db.Table("monitor_alert_source").
		Select(`
			id,
			COALESCE(name, '') AS name,
			COALESCE(type, 0) AS type,
			COALESCE(app_key, '') AS app_key,
			COALESCE(api_base_url, '') AS api_base_url,
			COALESCE(status, 0) AS status,
			COALESCE(remark, '') AS remark,
			COALESCE(create_time, '') AS create_time,
			COALESCE(update_time, '') AS update_time,
			COALESCE(key_id, 0) AS key_id,
			COALESCE(host_id, 0) AS host_id
		`).
		Where("type = ?", sourceType)
	if enabledOnly {
		db = db.Where("status = ?", 1)
	}

	err = db.Order("update_time DESC, id DESC").Scan(&list).Error
	return list, err
}

// CreateMonitorNotifyRobot inserts a notify robot record.
func CreateMonitorNotifyRobot(robot *model.MonitorNotifyRobotEntity) error {
	now := time.Now()
	robot.CreatedAt = now
	robot.UpdatedAt = now
	return Db.Create(robot).Error
}

// UpdateMonitorNotifyRobot updates a notify robot record.
func UpdateMonitorNotifyRobot(robot *model.MonitorNotifyRobotEntity) error {
	var existing model.MonitorNotifyRobotEntity
	if err := Db.First(&existing, robot.ID).Error; err != nil {
		return err
	}

	robot.CreatedAt = existing.CreatedAt
	robot.UpdatedAt = time.Now()
	return Db.Model(&existing).Updates(map[string]interface{}{
		"name":       robot.Name,
		"type":       robot.Type,
		"webhook":    robot.Webhook,
		"secret":     robot.Secret,
		"status":     robot.Status,
		"remark":     robot.Remark,
		"updated_at": robot.UpdatedAt,
		"server":     robot.Server,
		"port":       robot.Port,
		"username":   robot.Username,
		"password":   robot.Password,
		"nickname":   robot.Nickname,
		"headers":    robot.Headers,
		"method":     robot.Method,
		"template":   robot.Template,
	}).Error
}

// UpdateMonitorNotifyRobotStatus toggles a notify robot.
func UpdateMonitorNotifyRobotStatus(id uint, status int) error {
	return Db.Model(&model.MonitorNotifyRobotEntity{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// DeleteMonitorNotifyRobot deletes a notify robot record.
func DeleteMonitorNotifyRobot(id uint) error {
	return Db.Delete(&model.MonitorNotifyRobotEntity{}, id).Error
}

// CreateMonitorAlertSource inserts an alert source record.
func CreateMonitorAlertSource(source *model.MonitorAlertSourceEntity) error {
	now := time.Now().Unix()
	source.CreateTime = now
	source.UpdateTime = now
	return Db.Create(source).Error
}

// UpdateMonitorAlertSource updates an alert source record.
func UpdateMonitorAlertSource(source *model.MonitorAlertSourceEntity) error {
	var existing model.MonitorAlertSourceEntity
	if err := Db.First(&existing, source.ID).Error; err != nil {
		return err
	}

	source.CreateTime = existing.CreateTime
	source.UpdateTime = time.Now().Unix()
	return Db.Model(&existing).Updates(map[string]interface{}{
		"name":         source.Name,
		"type":         source.Type,
		"app_key":      source.AppKey,
		"api_base_url": source.APIBaseURL,
		"status":       source.Status,
		"remark":       source.Remark,
		"update_time":  source.UpdateTime,
		"key_id":       source.KeyID,
		"host_id":      source.HostID,
	}).Error
}

// UpdateMonitorAlertSourceStatus toggles an alert source.
func UpdateMonitorAlertSourceStatus(id uint, status int) error {
	return Db.Model(&model.MonitorAlertSourceEntity{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"update_time": time.Now().Unix(),
		}).Error
}

// DeleteMonitorAlertSource deletes an alert source record.
func DeleteMonitorAlertSource(id uint) error {
	return Db.Delete(&model.MonitorAlertSourceEntity{}, id).Error
}

// GetEnabledMonitorNotifyRobots returns enabled robots, optionally filtered by ids.
func GetEnabledMonitorNotifyRobots(robotIDs []uint) (list []model.MonitorNotifyRobotEntity, err error) {
	db := Db.Model(&model.MonitorNotifyRobotEntity{}).Where("status = ?", 1)
	if len(robotIDs) > 0 {
		db = db.Where("id IN ?", robotIDs)
	}
	err = db.Order("id ASC").Find(&list).Error
	return list, err
}

// GetMonitorNotifyRobotByID returns one notify robot.
func GetMonitorNotifyRobotByID(id uint) (*model.MonitorNotifyRobotEntity, error) {
	var robot model.MonitorNotifyRobotEntity
	if err := Db.First(&robot, id).Error; err != nil {
		return nil, err
	}
	return &robot, nil
}

// CreateMonitorWebhookLog inserts one webhook log row.
func CreateMonitorWebhookLog(logEntry *model.MonitorWebhookLogEntity) error {
	if logEntry.CreatedAt.IsZero() {
		logEntry.CreatedAt = time.Now()
	}
	return Db.Create(logEntry).Error
}

// UpdateMonitorWebhookLogResult updates the final delivery result of one webhook log.
func UpdateMonitorWebhookLogResult(id uint, status, errorMsg, notifyRobotIDs string, notifyCount, successCount, failedCount int64) error {
	return Db.Model(&model.MonitorWebhookLogEntity{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":           status,
			"error_msg":        errorMsg,
			"notify_robot_ids": notifyRobotIDs,
			"notify_count":     notifyCount,
			"success_count":    successCount,
			"failed_count":     failedCount,
		}).Error
}

// CreateMonitorWebhookNotifyLog inserts one delivery log row.
func CreateMonitorWebhookNotifyLog(logEntry *model.MonitorWebhookNotifyLogEntity) error {
	if logEntry.CreatedAt.IsZero() {
		logEntry.CreatedAt = time.Now()
	}
	return Db.Create(logEntry).Error
}

// EnsureAlertResourcesExist verifies target records before update/delete.
func EnsureAlertResourcesExist(resourceType string, id uint) error {
	switch resourceType {
	case "robot":
		var robot model.MonitorNotifyRobotEntity
		if err := Db.First(&robot, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("通知机器人不存在")
			}
			return err
		}
	case "source":
		var source model.MonitorAlertSourceEntity
		if err := Db.First(&source, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("告警源不存在")
			}
			return err
		}
	}
	return nil
}
