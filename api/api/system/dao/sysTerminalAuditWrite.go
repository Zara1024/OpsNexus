package dao

import (
	"time"

	"dodevops-api/api/system/model"
	. "dodevops-api/pkg/db"
)

func CreateTerminalAuditSession(session *model.SysSessionRecording) error {
	if session.CreateTime.IsZero() {
		session.CreateTime = time.Now()
	}
	if session.UpdateTime.IsZero() {
		session.UpdateTime = session.CreateTime
	}
	if session.StartTime.IsZero() {
		session.StartTime = session.CreateTime
	}
	return Db.Create(session).Error
}

func UpdateTerminalAuditSession(recordingID uint, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	if _, ok := updates["update_time"]; !ok {
		updates["update_time"] = time.Now()
	}
	return Db.Model(&model.SysSessionRecording{}).
		Where("id = ?", recordingID).
		Updates(updates).Error
}

func CreateTerminalAuditCommand(command *model.SysCommandAudit) error {
	if command.CreateTime.IsZero() {
		command.CreateTime = time.Now()
	}
	if command.ExecuteTime.IsZero() {
		command.ExecuteTime = command.CreateTime
	}
	return Db.Create(command).Error
}
