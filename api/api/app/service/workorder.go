package service

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"dodevops-api/api/app/model"
	"dodevops-api/common/result"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkOrderService struct {
	db *gorm.DB
}

type scriptReleaseEntity struct {
	ID                uint       `gorm:"column:id;primaryKey"`
	Title             string     `gorm:"column:title"`
	Reason            string     `gorm:"column:reason"`
	BusinessGroupID   uint       `gorm:"column:business_group_id"`
	AppID             uint       `gorm:"column:app_id"`
	AppName           string     `gorm:"column:app_name"`
	AppCode           string     `gorm:"column:app_code"`
	ApplicantID       uint       `gorm:"column:applicant_id"`
	ApplicantName     string     `gorm:"column:applicant_name"`
	ApproverID        uint       `gorm:"column:approver_id"`
	ApproverName      string     `gorm:"column:approver_name"`
	ExecutorID        uint       `gorm:"column:executor_id"`
	ExecutorName      string     `gorm:"column:executor_name"`
	ExecuteDir        string     `gorm:"column:execute_dir"`
	ScriptContent     string     `gorm:"column:script_content"`
	ApprovalStatus    int        `gorm:"column:approval_status"`
	ApprovalTime      *time.Time `gorm:"column:approval_time"`
	ApprovalRemark    string     `gorm:"column:approval_remark"`
	ExecuteStatus     int        `gorm:"column:execute_status"`
	Status            int        `gorm:"column:status"`
	StartTime         *time.Time `gorm:"column:start_time"`
	EndTime           *time.Time `gorm:"column:end_time"`
	Duration          int64      `gorm:"column:duration"`
	JenkinsEnvID      uint       `gorm:"column:jenkins_env_id"`
	BuildNumber       int64      `gorm:"column:build_number"`
	LogURL            string     `gorm:"column:log_url"`
	ErrorMessage      string     `gorm:"column:error_message"`
	CreatedAt         *time.Time `gorm:"column:created_at"`
	UpdatedAt         *time.Time `gorm:"column:updated_at"`
	Parameters        string     `gorm:"column:parameters"`
	ServerHostID      uint       `gorm:"column:server_host_id"`
	PullCodeStartTime *time.Time `gorm:"column:pull_code_start_time"`
	PullCodeEndTime   *time.Time `gorm:"column:pull_code_end_time"`
	ScriptOutput      string     `gorm:"column:script_output"`
}

func (scriptReleaseEntity) TableName() string {
	return "app_sh_release"
}

func NewWorkOrderService(db *gorm.DB) *WorkOrderService {
	return &WorkOrderService{db: db}
}

func (s *WorkOrderService) GetSummary(c *gin.Context) {
	summary := model.WorkOrderSummary{}

	summary.QuickDeploy = int(s.countTable("quick_deployments", "1=1"))
	summary.ScriptRelease = int(s.countTable("app_sh_release", "deleted_at IS NULL"))
	summary.ServiceRelase = int(s.countTable("app_service_release", "deleted_at IS NULL"))
	summary.SQLWorkOrder = int(s.countTable("cmdb_sql_work_order", "1=1"))
	summary.Total = summary.QuickDeploy + summary.ScriptRelease + summary.ServiceRelase + summary.SQLWorkOrder

	summary.Pending = int(
		s.countTable("quick_deployments", "status = 1") +
			s.countTable("app_sh_release", "deleted_at IS NULL AND status = 1") +
			s.countTable("app_service_release", "deleted_at IS NULL AND status = 1") +
			s.countTable("cmdb_sql_work_order", "status IN (1,2)"),
	)
	summary.Running = int(
		s.countTable("quick_deployments", "status = 2") +
			s.countTable("app_sh_release", "deleted_at IS NULL AND status = 6") +
			s.countTable("app_service_release", "deleted_at IS NULL AND status = 6") +
			s.countTable("cmdb_sql_work_order", "status = 4"),
	)
	summary.Success = int(
		s.countTable("quick_deployments", "status = 3") +
			s.countTable("app_sh_release", "deleted_at IS NULL AND status = 2") +
			s.countTable("app_service_release", "deleted_at IS NULL AND status = 2") +
			s.countTable("cmdb_sql_work_order", "status = 5"),
	)
	summary.Failed = int(
		s.countTable("quick_deployments", "status = 4") +
			s.countTable("app_sh_release", "deleted_at IS NULL AND status = 3") +
			s.countTable("app_service_release", "deleted_at IS NULL AND status = 3") +
			s.countTable("cmdb_sql_work_order", "status = 6"),
	)
	summary.Canceled = int(
		s.countTable("quick_deployments", "status = 5") +
			s.countTable("app_sh_release", "deleted_at IS NULL AND status = 5") +
			s.countTable("app_service_release", "deleted_at IS NULL AND status = 5") +
			s.countTable("cmdb_sql_work_order", "status = 3"),
	)

	result.Success(c, summary)
}

func (s *WorkOrderService) GetList(c *gin.Context, query model.WorkOrderQuery) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}

	items := make([]model.WorkOrderItem, 0)
	if query.Type == "" || query.Type == "quick" {
		items = append(items, s.loadQuickDeployments(query)...)
	}
	if query.Type == "" || query.Type == "script" {
		items = append(items, s.loadScriptReleases(query)...)
	}
	if query.Type == "" || query.Type == "service" {
		items = append(items, s.loadServiceReleases(query)...)
	}
	if query.Type == "" || query.Type == "sql" {
		items = append(items, s.loadSQLWorkOrders(query)...)
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].CreatedAt > items[j].CreatedAt
	})

	total := int64(len(items))
	start := (query.Page - 1) * query.PageSize
	if start > len(items) {
		start = len(items)
	}
	end := start + query.PageSize
	if end > len(items) {
		end = len(items)
	}

	result.Success(c, model.WorkOrderListResponse{
		Total: total,
		List:  items[start:end],
	})
}

func (s *WorkOrderService) GetDetail(c *gin.Context, category string, id uint) {
	switch category {
	case "quick":
		s.getQuickDetail(c, id)
	case "script":
		s.getScriptDetail(c, id)
	case "service":
		s.getServiceDetail(c, id)
	case "sql":
		s.getSQLDetail(c, id)
	default:
		result.Failed(c, 400, "不支持的工单类型")
	}
}

func (s *WorkOrderService) CreateScriptWorkOrder(c *gin.Context, req model.ScriptWorkOrderCreateRequest) {
	admin, username := currentAppWorkOrderAdmin(c)
	now := time.Now()
	entity := &scriptReleaseEntity{
		Title:           strings.TrimSpace(req.Title),
		Reason:          strings.TrimSpace(req.Reason),
		BusinessGroupID: chooseUint(req.BusinessGroupID, 1),
		AppID:           req.AppID,
		AppName:         strings.TrimSpace(req.AppName),
		AppCode:         strings.TrimSpace(req.AppCode),
		ApplicantID:     admin.ID,
		ApplicantName:   username,
		ExecuteDir:      strings.TrimSpace(req.ExecuteDir),
		ScriptContent:   strings.TrimSpace(req.ScriptContent),
		ApprovalStatus:  1,
		ExecuteStatus:   1,
		Status:          1,
		BuildNumber:     0,
		ServerHostID:    req.ServerHostID,
		CreatedAt:       &now,
		UpdatedAt:       &now,
	}
	if entity.Title == "" || entity.Reason == "" || entity.AppName == "" || entity.AppCode == "" || entity.ExecuteDir == "" || entity.ScriptContent == "" {
		result.Failed(c, 400, "标题、原因、应用信息、执行目录和脚本内容不能为空")
		return
	}
	if err := s.db.Create(entity).Error; err != nil {
		result.Failed(c, 500, "创建脚本工单失败: "+err.Error())
		return
	}
	result.Success(c, entity)
}

func (s *WorkOrderService) ApproveScriptRelease(c *gin.Context, id uint, req model.WorkOrderActionRequest) {
	entity, err := s.getScriptReleaseByID(id)
	if err != nil {
		result.Failed(c, 404, "脚本工单不存在")
		return
	}
	if !canApproveScriptRelease(entity.Status, entity.ApprovalStatus) {
		result.Failed(c, 400, "当前脚本工单状态不允许审批")
		return
	}
	admin, username := currentAppWorkOrderAdmin(c)
	now := time.Now()
	updates := map[string]interface{}{
		"approver_id":     admin.ID,
		"approver_name":   username,
		"approval_status": 2,
		"approval_time":   now,
		"approval_remark": firstNonEmpty(req.Comment, "同意执行"),
		"updated_at":      now,
	}
	if err = s.db.Table("app_sh_release").Where("id = ?", id).Updates(updates).Error; err != nil {
		result.Failed(c, 500, "脚本工单审批失败: "+err.Error())
		return
	}
	s.getScriptDetail(c, id)
}

func (s *WorkOrderService) RejectScriptRelease(c *gin.Context, id uint, req model.WorkOrderActionRequest) {
	entity, err := s.getScriptReleaseByID(id)
	if err != nil {
		result.Failed(c, 404, "脚本工单不存在")
		return
	}
	if !canRejectScriptRelease(entity.Status) {
		result.Failed(c, 400, "当前脚本工单状态不允许驳回")
		return
	}
	admin, username := currentAppWorkOrderAdmin(c)
	now := time.Now()
	updates := map[string]interface{}{
		"approver_id":     admin.ID,
		"approver_name":   username,
		"approval_status": 3,
		"approval_time":   now,
		"approval_remark": firstNonEmpty(req.Comment, "已驳回"),
		"status":          4,
		"updated_at":      now,
	}
	if err = s.db.Table("app_sh_release").Where("id = ?", id).Updates(updates).Error; err != nil {
		result.Failed(c, 500, "脚本工单驳回失败: "+err.Error())
		return
	}
	s.getScriptDetail(c, id)
}

func (s *WorkOrderService) ExecuteScriptRelease(c *gin.Context, id uint, req model.WorkOrderActionRequest) {
	entity, err := s.getScriptReleaseByID(id)
	if err != nil {
		result.Failed(c, 404, "脚本工单不存在")
		return
	}
	if !canExecuteScriptRelease(entity.Status, entity.ApprovalStatus) {
		result.Failed(c, 400, "当前脚本工单状态不允许执行")
		return
	}
	admin, username := currentAppWorkOrderAdmin(c)
	start := time.Now()
	if err = s.db.Table("app_sh_release").Where("id = ?", id).Updates(map[string]interface{}{
		"executor_id":    admin.ID,
		"executor_name":  username,
		"execute_status": 2,
		"status":         6,
		"start_time":     start,
		"updated_at":     start,
	}).Error; err != nil {
		result.Failed(c, 500, "更新脚本工单执行状态失败: "+err.Error())
		return
	}

	output, execErr := runScriptReleaseLocally(entity.ExecuteDir, entity.ScriptContent)
	end := time.Now()
	duration := int64(end.Sub(start).Seconds())
	updates := map[string]interface{}{
		"executor_id":   admin.ID,
		"executor_name": username,
		"end_time":      end,
		"duration":      duration,
		"script_output": output,
		"updated_at":    end,
	}
	if execErr != nil {
		updates["execute_status"] = 3
		updates["status"] = 3
		updates["error_message"] = execErr.Error()
	} else {
		updates["execute_status"] = 6
		updates["status"] = 2
		updates["error_message"] = ""
	}
	if strings.TrimSpace(req.Comment) != "" {
		updates["approval_remark"] = req.Comment
	}
	if err = s.db.Table("app_sh_release").Where("id = ?", id).Updates(updates).Error; err != nil {
		result.Failed(c, 500, "回写脚本工单执行结果失败: "+err.Error())
		return
	}
	if execErr != nil {
		result.Failed(c, 500, "脚本工单执行失败: "+execErr.Error())
		return
	}
	s.getScriptDetail(c, id)
}

func (s *WorkOrderService) countTable(tableName, where string) int64 {
	var count int64
	s.db.Table(tableName).Where(where).Count(&count)
	return count
}

func (s *WorkOrderService) loadQuickDeployments(query model.WorkOrderQuery) []model.WorkOrderItem {
	type row struct {
		ID              uint
		Title           string
		Status          int
		CreatorName     string
		BusinessGroupID uint
		TaskCount       int64
		CreatedAt       time.Time
		UpdatedAt       time.Time
		Duration        int64
	}
	var rows []row
	db := s.db.Table("quick_deployments").Select(`
		id,
		COALESCE(title, '') AS title,
		COALESCE(status, 0) AS status,
		COALESCE(creator_name, '') AS creator_name,
		COALESCE(business_group_id, 0) AS business_group_id,
		COALESCE(task_count, 0) AS task_count,
		created_at,
		updated_at,
		COALESCE(duration, 0) AS duration
	`)
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		db = db.Where("title LIKE ? OR creator_name LIKE ?", like, like)
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}
	db.Order("created_at DESC").Limit(50).Scan(&rows)

	items := make([]model.WorkOrderItem, 0, len(rows))
	for _, item := range rows {
		items = append(items, model.WorkOrderItem{
			Type:            "quick",
			TypeLabel:       "快速发布",
			ID:              item.ID,
			Title:           item.Title,
			ApplicantName:   item.CreatorName,
			CurrentHandler:  item.CreatorName,
			Status:          item.Status,
			StatusText:      quickStatusText(item.Status),
			BusinessGroupID: item.BusinessGroupID,
			CreatedAt:       formatTime(item.CreatedAt),
			UpdatedAt:       formatTime(item.UpdatedAt),
			Duration:        item.Duration,
			DetailHint:      "关联任务数: " + intToString(item.TaskCount),
		})
	}
	return items
}

func (s *WorkOrderService) loadScriptReleases(query model.WorkOrderQuery) []model.WorkOrderItem {
	type row struct {
		ID              uint
		Title           string
		AppName         string
		ApplicantName   string
		ApproverName    string
		ExecutorName    string
		ApprovalStatus  int
		ExecuteStatus   int
		Status          int
		BusinessGroupID uint
		CreatedAt       time.Time
		UpdatedAt       time.Time
		Duration        int64
	}
	var rows []row
	db := s.db.Table("app_sh_release").Select(`
		id,
		COALESCE(title, '') AS title,
		COALESCE(app_name, '') AS app_name,
		COALESCE(applicant_name, '') AS applicant_name,
		COALESCE(approver_name, '') AS approver_name,
		COALESCE(executor_name, '') AS executor_name,
		COALESCE(approval_status, 0) AS approval_status,
		COALESCE(execute_status, 0) AS execute_status,
		COALESCE(status, 0) AS status,
		COALESCE(business_group_id, 0) AS business_group_id,
		created_at,
		updated_at,
		COALESCE(duration, 0) AS duration
	`).Where("deleted_at IS NULL")
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		db = db.Where("title LIKE ? OR app_name LIKE ? OR applicant_name LIKE ?", like, like, like)
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}
	db.Order("created_at DESC").Limit(50).Scan(&rows)

	items := make([]model.WorkOrderItem, 0, len(rows))
	for _, item := range rows {
		items = append(items, model.WorkOrderItem{
			Type:            "script",
			TypeLabel:       "脚本发布",
			ID:              item.ID,
			Title:           item.Title,
			AppName:         item.AppName,
			ApplicantName:   item.ApplicantName,
			CurrentHandler:  firstNonEmpty(item.ExecutorName, item.ApproverName, item.ApplicantName),
			ApprovalStatus:  item.ApprovalStatus,
			ExecuteStatus:   item.ExecuteStatus,
			CanApprove:      canApproveScriptRelease(item.Status, item.ApprovalStatus),
			CanReject:       canRejectScriptRelease(item.Status),
			CanExecute:      canExecuteScriptRelease(item.Status, item.ApprovalStatus),
			Status:          item.Status,
			StatusText:      releaseStatusText(item.Status),
			BusinessGroupID: item.BusinessGroupID,
			CreatedAt:       formatTime(item.CreatedAt),
			UpdatedAt:       formatTime(item.UpdatedAt),
			Duration:        item.Duration,
			DetailHint:      "脚本发布工单",
		})
	}
	return items
}

func (s *WorkOrderService) loadServiceReleases(query model.WorkOrderQuery) []model.WorkOrderItem {
	type row struct {
		ID                   uint
		Title                string
		ApplicantName        string
		OwnerApproverName    string
		SecurityApproverName string
		TestApproverName     string
		Status               int
		BusinessGroupID      uint
		CreatedAt            time.Time
		UpdatedAt            time.Time
		Duration             int64
		ServiceCount         int64
	}
	var rows []row
	db := s.db.Table("app_service_release").Select(`
		id,
		COALESCE(title, '') AS title,
		COALESCE(applicant_name, '') AS applicant_name,
		COALESCE(owner_approver_name, '') AS owner_approver_name,
		COALESCE(security_approver_name, '') AS security_approver_name,
		COALESCE(test_approver_name, '') AS test_approver_name,
		COALESCE(status, 0) AS status,
		COALESCE(business_group_id, 0) AS business_group_id,
		created_at,
		updated_at,
		COALESCE(duration, 0) AS duration,
		COALESCE(service_count, 0) AS service_count
	`).Where("deleted_at IS NULL")
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		db = db.Where("title LIKE ? OR applicant_name LIKE ?", like, like)
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}
	db.Order("created_at DESC").Limit(50).Scan(&rows)

	items := make([]model.WorkOrderItem, 0, len(rows))
	for _, item := range rows {
		items = append(items, model.WorkOrderItem{
			Type:            "service",
			TypeLabel:       "服务上线",
			ID:              item.ID,
			Title:           item.Title,
			ApplicantName:   item.ApplicantName,
			CurrentHandler:  firstNonEmpty(item.OwnerApproverName, item.SecurityApproverName, item.TestApproverName, item.ApplicantName),
			Status:          item.Status,
			StatusText:      releaseStatusText(item.Status),
			BusinessGroupID: item.BusinessGroupID,
			CreatedAt:       formatTime(item.CreatedAt),
			UpdatedAt:       formatTime(item.UpdatedAt),
			Duration:        item.Duration,
			DetailHint:      "关联服务数: " + intToString(item.ServiceCount),
		})
	}
	return items
}

func (s *WorkOrderService) loadSQLWorkOrders(query model.WorkOrderQuery) []model.WorkOrderItem {
	type row struct {
		ID              uint
		Title           string
		InstanceName    string
		DatabaseName    string
		ApplicantName   string
		ApproverName    string
		ExecutorName    string
		Status          int
		RiskLevel       int
		OperationType   string
		RequireApproval bool
		CreateTime      time.Time
		UpdateTime      time.Time
		ExecutionTime   int64
	}
	var rows []row
	db := s.db.Table("cmdb_sql_work_order").Select(`
		id,
		COALESCE(title, '') AS title,
		COALESCE(instance_name, '') AS instance_name,
		COALESCE(database_name, '') AS database_name,
		COALESCE(applicant_name, '') AS applicant_name,
		COALESCE(approver_name, '') AS approver_name,
		COALESCE(executor_name, '') AS executor_name,
		COALESCE(status, 0) AS status,
		COALESCE(risk_level, 0) AS risk_level,
		COALESCE(operation_type, '') AS operation_type,
		COALESCE(require_approval, 0) AS require_approval,
		create_time,
		update_time,
		COALESCE(execution_time, 0) AS execution_time
	`)
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		db = db.Where("title LIKE ? OR database_name LIKE ? OR applicant_name LIKE ? OR sql_content LIKE ?", like, like, like, like)
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}
	db.Order("create_time DESC").Limit(50).Scan(&rows)

	items := make([]model.WorkOrderItem, 0, len(rows))
	for _, item := range rows {
		title := item.Title
		if strings.TrimSpace(title) == "" {
			title = firstNonEmpty(item.OperationType, "SQL") + " " + firstNonEmpty(item.DatabaseName, item.InstanceName, "变更")
		}
		items = append(items, model.WorkOrderItem{
			Type:            "sql",
			TypeLabel:       "SQL 工单",
			ID:              item.ID,
			Title:           title,
			AppName:         firstNonEmpty(item.InstanceName, item.DatabaseName),
			ApplicantName:   item.ApplicantName,
			CurrentHandler:  firstNonEmpty(item.ExecutorName, item.ApproverName, item.ApplicantName),
			Status:          item.Status,
			StatusText:      sqlWorkOrderStatusText(item.Status),
			CreatedAt:       formatTime(item.CreateTime),
			UpdatedAt:       formatTime(item.UpdateTime),
			Duration:        item.ExecutionTime,
			DetailHint:      fmt.Sprintf("%s / %s", firstNonEmpty(item.DatabaseName, "-"), firstNonEmpty(item.OperationType, "SQL")),
			RiskLevel:       sqlWorkOrderRiskLevel(item.RiskLevel),
			RiskText:        sqlWorkOrderRiskText(item.RiskLevel),
			AIDiagnosisPath: buildSQLWorkOrderAIDiagnosisPath(item.ID),
			KnowledgePath:   buildSQLWorkOrderKnowledgePath(title, item.OperationType, item.DatabaseName),
			CanApprove:      item.Status == 1,
			CanReject:       item.Status == 1,
			CanExecute:      item.Status == 2 || (item.Status == 1 && !item.RequireApproval),
		})
	}
	return items
}

func (s *WorkOrderService) getQuickDetail(c *gin.Context, id uint) {
	var basic map[string]interface{}
	var items []map[string]interface{}
	s.db.Table("quick_deployments").Where("id = ?", id).Take(&basic)
	s.db.Table("quick_deployment_tasks").Where("deployment_id = ?", id).Order("execute_order ASC").Find(&items)
	result.Success(c, model.WorkOrderDetail{
		Type:       "quick",
		TypeLabel:  "快速发布",
		ID:         id,
		Title:      stringValue(basic["title"]),
		Status:     intValue(basic["status"]),
		StatusText: quickStatusText(intValue(basic["status"])),
		Basic:      basic,
		Items:      items,
	})
}

func (s *WorkOrderService) getScriptDetail(c *gin.Context, id uint) {
	var basic map[string]interface{}
	s.db.Table("app_sh_release").Where("id = ?", id).Take(&basic)
	result.Success(c, model.WorkOrderDetail{
		Type:       "script",
		TypeLabel:  "脚本发布",
		ID:         id,
		Title:      stringValue(basic["title"]),
		Status:     intValue(basic["status"]),
		StatusText: releaseStatusText(intValue(basic["status"])),
		CanApprove: canApproveScriptRelease(intValue(basic["status"]), intValue(basic["approval_status"])),
		CanReject:  canRejectScriptRelease(intValue(basic["status"])),
		CanExecute: canExecuteScriptRelease(intValue(basic["status"]), intValue(basic["approval_status"])),
		Basic:      basic,
		Items:      []map[string]interface{}{},
	})
}

func (s *WorkOrderService) getServiceDetail(c *gin.Context, id uint) {
	var basic map[string]interface{}
	var items []map[string]interface{}
	s.db.Table("app_service_release").Where("id = ?", id).Take(&basic)
	s.db.Table("app_service_release_item").Where("release_id = ?", id).Order("execute_order ASC").Find(&items)
	result.Success(c, model.WorkOrderDetail{
		Type:       "service",
		TypeLabel:  "服务上线",
		ID:         id,
		Title:      stringValue(basic["title"]),
		Status:     intValue(basic["status"]),
		StatusText: releaseStatusText(intValue(basic["status"])),
		Basic:      basic,
		Items:      items,
	})
}

func (s *WorkOrderService) getSQLDetail(c *gin.Context, id uint) {
	var basic map[string]interface{}
	s.db.Table("cmdb_sql_work_order").Where("id = ?", id).Take(&basic)
	title := stringValue(basic["title"])
	if title == "" {
		title = "SQL 工单"
	}
	basic["risk_text"] = sqlWorkOrderRiskText(intValue(basic["risk_level"]))
	basic["risk_level_label"] = sqlWorkOrderRiskLevel(intValue(basic["risk_level"]))
	basic["ai_diagnosis_path"] = buildSQLWorkOrderAIDiagnosisPath(id)
	basic["knowledge_path"] = buildSQLWorkOrderKnowledgePath(title, stringValue(basic["operation_type"]), stringValue(basic["database_name"]))
	result.Success(c, model.WorkOrderDetail{
		Type:       "sql",
		TypeLabel:  "SQL 工单",
		ID:         id,
		Title:      title,
		Status:     intValue(basic["status"]),
		StatusText: sqlWorkOrderStatusText(intValue(basic["status"])),
		CanApprove: intValue(basic["status"]) == 1,
		CanReject:  intValue(basic["status"]) == 1,
		CanExecute: intValue(basic["status"]) == 2 || (intValue(basic["status"]) == 1 && !boolValue(basic["require_approval"])),
		Basic:      basic,
		Items:      []map[string]interface{}{},
	})
}

func quickStatusText(status int) string {
	switch status {
	case 1:
		return "待发布"
	case 2:
		return "发布中"
	case 3:
		return "发布成功"
	case 4:
		return "发布失败"
	case 5:
		return "已取消"
	default:
		return "未知"
	}
}

func releaseStatusText(status int) string {
	switch status {
	case 1:
		return "流程进行中"
	case 2:
		return "执行成功"
	case 3:
		return "执行失败"
	case 4:
		return "已驳回"
	case 5:
		return "已取消"
	case 6:
		return "执行中"
	default:
		return "未知"
	}
}

func sqlWorkOrderStatusText(status int) string {
	switch status {
	case 1:
		return "待审批"
	case 2:
		return "待执行"
	case 3:
		return "已驳回"
	case 4:
		return "执行中"
	case 5:
		return "执行成功"
	case 6:
		return "执行失败"
	default:
		return "未知"
	}
}

func sqlWorkOrderRiskLevel(level int) string {
	switch level {
	case 2:
		return "high"
	case 1:
		return "medium"
	default:
		return "low"
	}
}

func sqlWorkOrderRiskText(level int) string {
	switch level {
	case 2:
		return "高风险"
	case 1:
		return "中风险"
	default:
		return "低风险"
	}
}

func buildSQLWorkOrderAIDiagnosisPath(id uint) string {
	if id == 0 {
		return ""
	}
	values := url.Values{}
	values.Set("scene", "sql_work_order")
	values.Set("targetId", fmt.Sprintf("%d", id))
	values.Set("keyword", "sql rollback ddl review")
	values.Set("templateName", "yaml_change_review")
	values.Set("source", "work-order-center")
	values.Set("autoRun", "1")
	return "/ai/diagnosis?" + values.Encode()
}

func buildSQLWorkOrderKnowledgePath(title, operationType, databaseName string) string {
	values := url.Values{}
	values.Set("type", "change")
	values.Set("keyword", strings.TrimSpace(firstNonEmpty(title, operationType, databaseName, "sql change")))
	return "/knowledge/base?" + values.Encode()
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format("2006-01-02 15:04:05")
}

func intToString(value int64) string {
	return fmt.Sprintf("%d", value)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func intValue(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint64:
		return int(v)
	default:
		return 0
	}
}

func stringValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case []uint8:
		return string(v)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func boolValue(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case int32:
		return v != 0
	case int64:
		return v != 0
	case uint:
		return v != 0
	case uint64:
		return v != 0
	case []uint8:
		return string(v) == "1"
	default:
		return false
	}
}

func canApproveScriptRelease(status, approvalStatus int) bool {
	return status == 1 && approvalStatus != 2 && approvalStatus != 3
}

func canRejectScriptRelease(status int) bool {
	return status == 1
}

func canExecuteScriptRelease(status, approvalStatus int) bool {
	return status == 1 && approvalStatus == 2
}

func (s *WorkOrderService) getScriptReleaseByID(id uint) (*scriptReleaseEntity, error) {
	var entity scriptReleaseEntity
	err := s.db.Table("app_sh_release").Where("id = ?", id).Take(&entity).Error
	return &entity, err
}

func currentAppWorkOrderAdmin(c *gin.Context) (jwtAdmin systemUser, username string) {
	admin, err := jwt.GetAdmin(c)
	if err != nil || admin == nil {
		return systemUser{}, "unknown"
	}
	return systemUser{ID: admin.ID, Username: admin.Username}, firstNonEmpty(admin.Username, admin.Nickname, "unknown")
}

type systemUser struct {
	ID       uint
	Username string
}

func chooseUint(value, fallback uint) uint {
	if value > 0 {
		return value
	}
	return fallback
}

func runScriptReleaseLocally(executeDir, scriptContent string) (string, error) {
	executeDir = strings.TrimSpace(executeDir)
	if executeDir == "" {
		return "", fmt.Errorf("执行目录为空")
	}
	if _, err := os.Stat(executeDir); err != nil {
		return "", fmt.Errorf("执行目录不存在: %v", err)
	}

	tmpFile, err := os.CreateTemp("", "opsnexus-script-release-*.sh")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	script := "#!/bin/bash\nset -e\n" + scriptContent + "\n"
	if _, err = tmpFile.WriteString(script); err != nil {
		_ = tmpFile.Close()
		return "", err
	}
	if err = tmpFile.Close(); err != nil {
		return "", err
	}
	if err = os.Chmod(tmpFile.Name(), 0o700); err != nil {
		return "", err
	}

	command := fmt.Sprintf("cd %q && bash %q", executeDir, filepath.ToSlash(tmpFile.Name()))
	cmd := exec.Command("bash", "-lc", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
