package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	aiModel "dodevops-api/api/ai/model"
	cmdbService "dodevops-api/api/cmdb/service"
	"dodevops-api/common/result"
	"dodevops-api/pkg/db"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var builtinInspectionTemplates = []aiModel.AIAssistantInspectionTemplate{
	{
		Name:        "Linux 基础巡检",
		Category:    "基础",
		Scope:       "host",
		Description: "检查主机名、系统时间、系统负载、磁盘、内存和监听端口。",
		PromptHint:  "适用于日常主机健康检查与交接巡检。",
		IsBuiltin:   true,
		Checks: []aiModel.AIAssistantInspectionTemplateCheck{
			{Name: "主机名", Command: "hostname"},
			{Name: "系统时间", Command: "date"},
			{Name: "系统负载", Command: "uptime"},
			{Name: "磁盘占用", Command: "df -h"},
			{Name: "内存情况", Command: "free -m"},
			{Name: "监听端口", Command: "ss -lntp"},
		},
	},
	{
		Name:        "Docker 运行态巡检",
		Category:    "容器",
		Scope:       "host",
		Description: "检查 Docker 容器存活、镜像运行情况与资源快照。",
		PromptHint:  "适用于部署宿主机的容器运行态检查。",
		IsBuiltin:   true,
		Checks: []aiModel.AIAssistantInspectionTemplateCheck{
			{Name: "主机名", Command: "hostname"},
			{Name: "容器列表", Command: "docker ps"},
			{Name: "镜像列表", Command: "docker images"},
			{Name: "磁盘占用", Command: "df -h"},
			{Name: "系统负载", Command: "uptime"},
		},
	},
	{
		Name:        "Java 服务巡检",
		Category:    "应用",
		Scope:       "host",
		Description: "检查 Java 进程、端口、内存与系统状态。",
		PromptHint:  "适用于 Java 服务上线前后的运行态核对。",
		IsBuiltin:   true,
		Checks: []aiModel.AIAssistantInspectionTemplateCheck{
			{Name: "Java 进程", Command: "ps -ef"},
			{Name: "监听端口", Command: "ss -lntp"},
			{Name: "内存情况", Command: "free -m"},
			{Name: "磁盘占用", Command: "df -h"},
			{Name: "系统负载", Command: "uptime"},
		},
	},
}

func (s *AIService) ensureBuiltinInspectionTemplates() {
	for _, item := range builtinInspectionTemplates {
		var existing aiModel.AIInspectionTemplateEntity
		err := db.Db.Where("name = ?", item.Name).First(&existing).Error
		if err == nil {
			continue
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			continue
		}

		checksJSON, _ := json.Marshal(item.Checks)
		_ = db.Db.Create(&aiModel.AIInspectionTemplateEntity{
			Name:        item.Name,
			Category:    item.Category,
			Scope:       item.Scope,
			Description: item.Description,
			ChecksJSON:  string(checksJSON),
			PromptHint:  item.PromptHint,
			Enabled:     1,
			IsBuiltin:   1,
			CreatedBy:   "system",
			UpdatedBy:   "system",
		}).Error
	}
}

func (s *AIService) listInspectionTemplates(limit int) []aiModel.AIAssistantInspectionTemplate {
	s.ensureBuiltinInspectionTemplates()

	if limit <= 0 {
		limit = 20
	}
	var rows []aiModel.AIInspectionTemplateEntity
	if err := db.Db.Where("enabled = 1").Order("is_builtin DESC, id ASC").Limit(limit).Find(&rows).Error; err != nil {
		return nil
	}

	items := make([]aiModel.AIAssistantInspectionTemplate, 0, len(rows))
	for _, row := range rows {
		items = append(items, toInspectionTemplateDTO(row))
	}
	return items
}

func toInspectionTemplateDTO(row aiModel.AIInspectionTemplateEntity) aiModel.AIAssistantInspectionTemplate {
	var checks []aiModel.AIAssistantInspectionTemplateCheck
	_ = json.Unmarshal([]byte(row.ChecksJSON), &checks)
	return aiModel.AIAssistantInspectionTemplate{
		ID:          row.ID,
		Name:        row.Name,
		Category:    row.Category,
		Scope:       row.Scope,
		Description: row.Description,
		PromptHint:  row.PromptHint,
		IsBuiltin:   row.IsBuiltin == 1,
		Checks:      checks,
	}
}

func (s *AIService) getInspectionTemplateByID(id uint) (*aiModel.AIAssistantInspectionTemplate, error) {
	s.ensureBuiltinInspectionTemplates()
	var row aiModel.AIInspectionTemplateEntity
	if err := db.Db.Where("id = ? AND enabled = 1", id).First(&row).Error; err != nil {
		return nil, err
	}
	dto := toInspectionTemplateDTO(row)
	return &dto, nil
}

func (s *AIService) getInspectionTemplateForMessage(message string) *aiModel.AIAssistantInspectionTemplate {
	templates := s.listInspectionTemplates(20)
	lower := strings.ToLower(message)
	for _, template := range templates {
		nameLower := strings.ToLower(template.Name)
		categoryLower := strings.ToLower(template.Category)
		if strings.Contains(lower, nameLower) || strings.Contains(lower, categoryLower) {
			t := template
			return &t
		}
	}
	if len(templates) > 0 {
		t := templates[0]
		return &t
	}
	return nil
}

func (s *AIService) saveInspectionReport(sessionID string, userID uint, result *aiModel.AIAssistantInspectionResult) {
	if result == nil {
		return
	}
	checkJSON, _ := json.Marshal(result.Checks)
	_ = db.Db.Create(&aiModel.AIInspectionReportEntity{
		SessionID:    sessionID,
		UserID:       userID,
		TemplateID:   result.TemplateID,
		TemplateName: result.TemplateName,
		Scope:        "host",
		TargetID:     result.HostID,
		TargetName:   result.HostName,
		Summary:      result.Summary,
		Report:       result.Report,
		CheckResults: string(checkJSON),
		Status:       "completed",
	}).Error
}

func (s *AIService) listRecentInspectionReports(userID uint, limit int) []aiModel.AIAssistantInspectionReportSummary {
	if limit <= 0 {
		limit = 10
	}
	var rows []aiModel.AIInspectionReportEntity
	if err := db.Db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&rows).Error; err != nil {
		return nil
	}

	items := make([]aiModel.AIAssistantInspectionReportSummary, 0, len(rows))
	for _, row := range rows {
		items = append(items, aiModel.AIAssistantInspectionReportSummary{
			ID:           row.ID,
			TemplateID:   row.TemplateID,
			TemplateName: row.TemplateName,
			Scope:        row.Scope,
			TargetID:     row.TargetID,
			TargetName:   row.TargetName,
			Summary:      row.Summary,
			Status:       row.Status,
			CreatedAt:    row.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return items
}

func (s *AIService) loadAssistantContext(sessionID string, userID uint) *aiModel.AIAssistantContext {
	var row aiModel.AIAssistantSessionContextEntity
	if err := db.Db.Where("session_id = ? AND user_id = ?", sessionID, userID).First(&row).Error; err != nil {
		return nil
	}
	return &aiModel.AIAssistantContext{
		CurrentScope:         row.CurrentScope,
		CurrentHostID:        row.CurrentHostID,
		CurrentHostName:      row.CurrentHostName,
		CurrentClusterID:     row.CurrentClusterID,
		CurrentClusterName:   row.CurrentClusterName,
		CurrentNamespace:     row.CurrentNamespace,
		CurrentWorkloadType:  row.CurrentWorkloadType,
		CurrentWorkloadName:  row.CurrentWorkloadName,
		CurrentWorkorderType: row.CurrentWorkOrderType,
		CurrentWorkorderID:   row.CurrentWorkOrderID,
		CurrentDeploymentID:  row.CurrentDeploymentID,
		LastIntent:           row.LastIntent,
		Summary:              row.Summary,
	}
}

func (s *AIService) saveAssistantContext(sessionID string, userID uint, context *aiModel.AIAssistantContext) {
	if context == nil {
		return
	}

	var row aiModel.AIAssistantSessionContextEntity
	err := db.Db.Where("session_id = ? AND user_id = ?", sessionID, userID).First(&row).Error
	if err == gorm.ErrRecordNotFound {
		row.SessionID = sessionID
		row.UserID = userID
	} else if err != nil {
		return
	}

	row.CurrentScope = context.CurrentScope
	row.CurrentHostID = context.CurrentHostID
	row.CurrentHostName = context.CurrentHostName
	row.CurrentClusterID = context.CurrentClusterID
	row.CurrentClusterName = context.CurrentClusterName
	row.CurrentNamespace = context.CurrentNamespace
	row.CurrentWorkloadType = context.CurrentWorkloadType
	row.CurrentWorkloadName = context.CurrentWorkloadName
	row.CurrentWorkOrderType = context.CurrentWorkorderType
	row.CurrentWorkOrderID = context.CurrentWorkorderID
	row.CurrentDeploymentID = context.CurrentDeploymentID
	row.LastIntent = context.LastIntent
	row.Summary = context.Summary

	if row.ID == 0 {
		_ = db.Db.Create(&row).Error
		return
	}
	_ = db.Db.Save(&row).Error
}

func (s *AIService) createPendingConfirmation(sessionID string, userID uint, scope, actionType string, targetID uint, targetName, command string, payload interface{}, summary string) *aiModel.AIAssistantPendingConfirmation {
	payloadJSON, _ := json.Marshal(payload)
	entity := &aiModel.AIAssistantConfirmationEntity{
		SessionID:  sessionID,
		UserID:     userID,
		Status:     "pending",
		Scope:      scope,
		ActionType: actionType,
		TargetID:   targetID,
		TargetName: targetName,
		Command:    command,
		Payload:    string(payloadJSON),
		Summary:    summary,
		ExpiresAt:  time.Now().Add(15 * time.Minute),
	}
	if err := db.Db.Create(entity).Error; err != nil {
		return nil
	}
	return toPendingConfirmationDTO(entity)
}

func toPendingConfirmationDTO(row *aiModel.AIAssistantConfirmationEntity) *aiModel.AIAssistantPendingConfirmation {
	if row == nil {
		return nil
	}
	return &aiModel.AIAssistantPendingConfirmation{
		ID:         row.ID,
		Status:     row.Status,
		ActionType: row.ActionType,
		Scope:      row.Scope,
		TargetName: row.TargetName,
		Command:    row.Command,
		Summary:    row.Summary,
		ExpiresAt:  row.ExpiresAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *AIService) getPendingConfirmation(userID uint, id uint) (*aiModel.AIAssistantConfirmationEntity, error) {
	var row aiModel.AIAssistantConfirmationEntity
	if err := db.Db.Where("id = ? AND user_id = ?", id, userID).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (s *AIService) finalizeConfirmation(row *aiModel.AIAssistantConfirmationEntity, status string, summary string) {
	if row == nil {
		return
	}
	row.Status = status
	row.ResultSummary = summary
	row.UpdatedAt = time.Now()
	_ = db.Db.Save(row).Error
}

func (s *AIService) ListInspectionTemplates(c *gin.Context) {
	result.Success(c, s.listInspectionTemplates(50))
}

func (s *AIService) ListInspectionReports(c *gin.Context) {
	userID, _ := jwt.GetAdminId(c)
	result.Success(c, s.listRecentInspectionReports(userID, 20))
}

func (s *AIService) DecideConfirmation(c *gin.Context, id uint, decision string) {
	userID, _ := jwt.GetAdminId(c)
	row, err := s.getPendingConfirmation(userID, id)
	if err != nil {
		result.Failed(c, 404, "待确认任务不存在")
		return
	}
	if row.Status != "pending" {
		result.Success(c, gin.H{
			"id":            row.ID,
			"status":        row.Status,
			"resultSummary": row.ResultSummary,
		})
		return
	}
	if time.Now().After(row.ExpiresAt) {
		s.finalizeConfirmation(row, "expired", "确认任务已过期")
		result.Failed(c, 400, "确认任务已过期")
		return
	}

	switch strings.ToLower(strings.TrimSpace(decision)) {
	case "cancel", "reject":
		s.finalizeConfirmation(row, "cancelled", "已取消执行")
		result.Success(c, gin.H{
			"id":            row.ID,
			"status":        row.Status,
			"resultSummary": row.ResultSummary,
		})
		return
	case "approve", "confirm":
		summary, execErr := s.executeConfirmedAction(c, row)
		if execErr != nil {
			s.finalizeConfirmation(row, "failed", execErr.Error())
			result.Failed(c, 500, execErr.Error())
			return
		}
		s.finalizeConfirmation(row, "executed", summary)
		result.Success(c, gin.H{
			"id":            row.ID,
			"status":        row.Status,
			"resultSummary": row.ResultSummary,
		})
		return
	default:
		result.Failed(c, 400, "不支持的确认动作")
	}
}

func (s *AIService) executeConfirmedAction(c *gin.Context, row *aiModel.AIAssistantConfirmationEntity) (string, error) {
	if row == nil {
		return "", fmt.Errorf("invalid confirmation")
	}

	switch row.Scope {
	case "host":
		if strings.TrimSpace(row.Command) == "" || row.TargetID == 0 {
			return "", fmt.Errorf("确认任务缺少执行命令或目标主机")
		}
		output, err := cmdbService.GetCmdbHostSSHService().ExecuteCommand(c, row.TargetID, row.Command)
		if err != nil {
			return "", fmt.Errorf("执行高风险命令失败: %v", err)
		}
		return fmt.Sprintf("已执行 `%s`，输出摘要：%s", row.Command, truncateForReport(output.Output, 300)), nil
	default:
		return "", fmt.Errorf("当前确认任务暂不支持自动执行")
	}
}
