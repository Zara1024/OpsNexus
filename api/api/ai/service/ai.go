package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	aiModel "dodevops-api/api/ai/model"
	appModel "dodevops-api/api/app/model"
	cmdbModel "dodevops-api/api/cmdb/model"
	k8sModel "dodevops-api/api/k8s/model"
	knowledgeModel "dodevops-api/api/knowledge/model"
	monitorDao "dodevops-api/api/monitor/dao"
	monitorModel "dodevops-api/api/monitor/model"
	systemDao "dodevops-api/api/system/dao"
	systemModel "dodevops-api/api/system/model"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"dodevops-api/pkg/db"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AIService struct{}

func NewAIService() *AIService {
	return &AIService{}
}

func (s *AIService) ListTemplates(c *gin.Context) {
	var list []aiModel.PromptTemplate
	if err := db.Db.Where("enabled = ?", 1).Order("id ASC").Find(&list).Error; err != nil {
		result.Failed(c, 500, "获取提示词模板失败: "+err.Error())
		return
	}
	result.Success(c, list)
}

func (s *AIService) SuggestKnowledge(c *gin.Context, keyword string) {
	keyword = strings.TrimSpace(keyword)
	var list []knowledgeModel.KnowledgeBase
	query := db.Db.Model(&knowledgeModel.KnowledgeBase{}).Where("enabled = ?", 1)
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR content LIKE ? OR keywords LIKE ? OR tags LIKE ?", like, like, like, like)
	}
	if err := query.Order("score DESC, use_count DESC, update_time DESC").Limit(10).Find(&list).Error; err != nil {
		result.Failed(c, 500, "获取知识建议失败: "+err.Error())
		return
	}
	result.Success(c, list)
}

func (s *AIService) RenderPrompt(c *gin.Context, req aiModel.AIPromptRenderRequest) {
	req.TemplateName = strings.TrimSpace(req.TemplateName)
	req.InputMessage = strings.TrimSpace(req.InputMessage)
	req.Intent = strings.TrimSpace(req.Intent)
	if req.TemplateName == "" {
		result.Failed(c, 400, "模板名称不能为空")
		return
	}

	var tmpl aiModel.PromptTemplate
	if err := db.Db.Where("name = ? AND enabled = 1", req.TemplateName).First(&tmpl).Error; err != nil {
		result.Failed(c, 404, "提示词模板不存在")
		return
	}

	if strings.TrimSpace(req.SessionID) == "" {
		req.SessionID = generateSessionID()
	}
	if req.Variables == nil {
		req.Variables = map[string]interface{}{}
	}

	knowledgeItems := make([]map[string]interface{}, 0)
	if len(req.KnowledgeIDs) > 0 {
		var kbItems []knowledgeModel.KnowledgeBase
		if err := db.Db.Where("id IN ?", req.KnowledgeIDs).Find(&kbItems).Error; err == nil {
			for _, item := range kbItems {
				knowledgeItems = append(knowledgeItems, map[string]interface{}{
					"id":       item.ID,
					"title":    item.Title,
					"type":     item.Type,
					"category": item.Category,
					"content":  item.Content,
					"tags":     item.Tags,
				})
				db.Db.Model(&knowledgeModel.KnowledgeBase{}).Where("id = ?", item.ID).UpdateColumn("use_count", gorm.Expr("use_count + ?", 1))
			}
		}
	}

	if req.InputMessage != "" {
		req.Variables["userMessage"] = req.InputMessage
	}
	if req.Intent != "" {
		req.Variables["intent"] = req.Intent
	}
	if len(knowledgeItems) > 0 {
		knowledgeJSON, _ := json.MarshalIndent(knowledgeItems, "", "  ")
		req.Variables["knowledgeContext"] = string(knowledgeJSON)
	}

	rendered := renderPromptTemplate(tmpl.Template, req.Variables)
	systemPrompt := renderPromptTemplate(tmpl.SystemPrompt, req.Variables)

	userID, _ := jwt.GetAdminId(c)
	entitiesJSON, _ := json.Marshal(req.Variables)
	db.Db.Create(&aiModel.AIChatHistory{
		SessionID:  req.SessionID,
		UserID:     userID,
		Role:       "user",
		Message:    firstNonEmpty(req.InputMessage, tmpl.Name),
		Intent:     req.Intent,
		IntentConf: 1,
		Entities:   string(entitiesJSON),
		TaskType:   tmpl.Scene,
		Status:     2,
		CreateTime: util.HTime{Time: time.Now()},
	})
	db.Db.Create(&aiModel.AIChatHistory{
		SessionID:  req.SessionID,
		UserID:     userID,
		Role:       "assistant",
		Message:    rendered,
		Intent:     req.Intent,
		IntentConf: 1,
		Entities:   string(entitiesJSON),
		TaskType:   tmpl.Scene,
		Status:     2,
		CreateTime: util.HTime{Time: time.Now()},
	})

	result.Success(c, aiModel.AIPromptRenderResponse{
		SessionID:      req.SessionID,
		Template:       tmpl,
		RenderedPrompt: rendered,
		SystemPrompt:   systemPrompt,
		KnowledgeItems: knowledgeItems,
	})
}

func (s *AIService) ListHistory(c *gin.Context) {
	userID, _ := jwt.GetAdminId(c)
	type sessionRow struct {
		SessionID    string    `json:"sessionId"`
		LatestTime   time.Time `json:"latestTime"`
		MessageCount int64     `json:"messageCount"`
	}
	var rows []sessionRow
	if err := db.Db.Table("ai_agent_chat_history").
		Select("session_id, MAX(create_time) AS latest_time, COUNT(*) AS message_count").
		Where("user_id = ?", userID).
		Group("session_id").
		Order("latest_time DESC").
		Limit(20).
		Scan(&rows).Error; err != nil {
		result.Failed(c, 500, "获取 AI 会话历史失败: "+err.Error())
		return
	}
	result.Success(c, rows)
}

func (s *AIService) GetHistoryDetail(c *gin.Context, sessionID string) {
	userID, _ := jwt.GetAdminId(c)
	var rows []aiModel.AIChatHistory
	if err := db.Db.Where("session_id = ? AND user_id = ?", sessionID, userID).Order("create_time ASC, id ASC").Find(&rows).Error; err != nil {
		result.Failed(c, 500, "获取 AI 会话详情失败: "+err.Error())
		return
	}
	result.Success(c, rows)
}

func (s *AIService) Diagnose(c *gin.Context, req aiModel.AIDiagnosisRequest) {
	scene := strings.TrimSpace(strings.ToLower(req.Scene))
	targetID := strings.TrimSpace(req.TargetID)
	keyword := strings.TrimSpace(req.Keyword)
	templateName := strings.TrimSpace(req.TemplateName)
	if scene == "" {
		result.Failed(c, 400, "诊断场景不能为空")
		return
	}

	sessionID := generateSessionID()
	knowledgeItems := make([]map[string]interface{}, 0)
	renderedPrompt := ""
	systemPrompt := ""
	title := ""
	report := ""
	usedFallback := true

	switch scene {
	case "terminal_audit":
		detail, err := systemDao.GetTerminalAuditSessionDetail(targetID)
		if err != nil || detail.Session.SessionID == "" {
			result.Failed(c, 404, "终端审计会话不存在")
			return
		}
		title = "终端审计复盘"
		knowledgeItems = s.mergeKnowledgeItems(
			s.collectKnowledgeItems(strings.TrimSpace(detail.Session.SessionID+" "+detail.Session.LatestCommand+" terminal audit "+keyword)),
			s.loadKnowledgeItemsByIDs(req.KnowledgeIDs),
		)
		renderedPrompt, systemPrompt = s.buildDiagnosisPrompt(firstNonEmpty(templateName, "terminal_audit_review"), map[string]interface{}{
			"sessionId":        detail.Session.SessionID,
			"sessionType":      detail.Session.SessionType,
			"riskLevel":        detail.Session.RiskLevel,
			"latestCommand":    detail.Session.LatestCommand,
			"commandCount":     detail.Session.CommandCount,
			"sensitiveCount":   detail.Session.SensitiveCommandCount,
			"commandTimeline":  detail.Commands,
			"knowledgeContext": stringify(knowledgeItems),
		})
		report = buildTerminalAuditDiagnosis(detail, knowledgeItems)
	case "sql_work_order":
		orderID, err := strconv.Atoi(targetID)
		if err != nil || orderID <= 0 {
			result.Failed(c, 400, "SQL 工单 ID 无效")
			return
		}
		var order cmdbModel.CmdbSQLWorkOrder
		if err = db.Db.First(&order, orderID).Error; err != nil {
			result.Failed(c, 404, "SQL 工单不存在")
			return
		}
		title = "SQL 工单分析"
		knowledgeItems = s.mergeKnowledgeItems(
			s.collectKnowledgeItems(strings.TrimSpace(order.Title+" "+order.SQLContent+" "+order.RollbackHint+" "+keyword)),
			s.loadKnowledgeItemsByIDs(req.KnowledgeIDs),
		)
		renderedPrompt, systemPrompt = s.buildDiagnosisPrompt(firstNonEmpty(templateName, "yaml_change_review"), map[string]interface{}{
			"orderNo":          order.OrderNo,
			"databaseName":     order.DatabaseName,
			"operationType":    order.OperationType,
			"riskLevel":        order.RiskLevel,
			"riskSummary":      order.RiskSummary,
			"rollbackHint":     order.RollbackHint,
			"sqlContent":       order.SQLContent,
			"knowledgeContext": stringify(knowledgeItems),
		})
		report = buildSQLWorkOrderDiagnosis(order, knowledgeItems)
	case "workload_capacity":
		historyID, err := strconv.Atoi(targetID)
		if err != nil || historyID <= 0 {
			result.Failed(c, 400, "容量建议快照 ID 无效")
			return
		}
		context, err := s.loadWorkloadCapacityDiagnosisContext(uint(historyID))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				result.Failed(c, 404, "容量建议快照不存在")
				return
			}
			result.Failed(c, 500, "获取容量建议快照失败: "+err.Error())
			return
		}

		title = "工作负载容量治理诊断"
		knowledgeItems = s.mergeKnowledgeItems(
			s.collectKnowledgeItems(strings.TrimSpace(s.buildWorkloadCapacityKnowledgeKeyword(context)+" "+keyword)),
			s.loadKnowledgeItemsByIDs(req.KnowledgeIDs),
		)
		renderedPrompt, systemPrompt = s.buildWorkloadCapacityDiagnosisPrompt(firstNonEmpty(templateName, "incident_analysis"), context, knowledgeItems)
		report = buildWorkloadCapacityDiagnosis(context, knowledgeItems)
	case "knowledge_search":
		title = "知识推荐分析"
		knowledgeItems = s.mergeKnowledgeItems(
			s.collectKnowledgeItems(keyword),
			s.loadKnowledgeItemsByIDs(req.KnowledgeIDs),
		)
		renderedPrompt = keyword
		systemPrompt = "knowledge_search"
		report = buildKnowledgeSearchDiagnosis(keyword, knowledgeItems)
	case "alert_analysis":
		title = "告警事件分析"
		alertKeyword := firstNonEmpty(targetID, keyword)
		summary, err := monitorDao.GetMonitorAlertSummary()
		if err != nil {
			result.Failed(c, 500, "获取告警摘要失败: "+err.Error())
			return
		}
		incidents, _, err := monitorDao.GetMonitorIncidentList(monitorModel.MonitorIncidentQuery{
			Keyword:  alertKeyword,
			Status:   -1,
			PageSize: 10,
			PageNum:  1,
		})
		if err != nil {
			result.Failed(c, 500, "获取告警事件失败: "+err.Error())
			return
		}
		knowledgeItems = s.mergeKnowledgeItems(
			s.collectKnowledgeItems(strings.TrimSpace(alertKeyword+" alert incident "+summary.LatestAlertTime)),
			s.loadKnowledgeItemsByIDs(req.KnowledgeIDs),
		)
		renderedPrompt, systemPrompt = s.buildDiagnosisPrompt(firstNonEmpty(templateName, "incident_analysis"), map[string]interface{}{
			"alertKeyword":     firstNonEmpty(alertKeyword, "current alerts"),
			"alertSummary":     stringify(summary),
			"alertIncidents":   stringify(incidents),
			"knowledgeContext": stringify(knowledgeItems),
		})
		report = buildAlertAnalysisDiagnosis(summary, incidents, knowledgeItems)
	case "deployment_review":
		deploymentID, err := strconv.Atoi(targetID)
		if err != nil || deploymentID <= 0 {
			result.Failed(c, 400, "发布 ID 无效")
			return
		}
		var deployment appModel.QuickDeployment
		if err = db.Db.Preload("Tasks").First(&deployment, deploymentID).Error; err != nil {
			result.Failed(c, 404, "快速发布记录不存在")
			return
		}
		title = "发布复盘分析"
		taskKeywords := make([]string, 0, len(deployment.Tasks))
		for _, item := range deployment.Tasks {
			taskKeywords = append(taskKeywords, strings.TrimSpace(item.AppName+" "+item.Environment+" "+item.ErrorMessage))
		}
		knowledgeItems = s.mergeKnowledgeItems(
			s.collectKnowledgeItems(strings.TrimSpace(deployment.Title+" "+deployment.Description+" "+strings.Join(taskKeywords, " ")+" "+keyword)),
			s.loadKnowledgeItemsByIDs(req.KnowledgeIDs),
		)
		renderedPrompt, systemPrompt = s.buildDiagnosisPrompt(firstNonEmpty(templateName, "incident_analysis"), map[string]interface{}{
			"deploymentTitle":   deployment.Title,
			"deploymentDesc":    deployment.Description,
			"deploymentStatus":  deployment.Status,
			"deploymentCreated": deployment.CreatedAt.Format("2006-01-02 15:04:05"),
			"deploymentTasks":   stringify(deployment.Tasks),
			"knowledgeContext":  stringify(knowledgeItems),
		})
		report = buildDeploymentReviewDiagnosis(deployment, deployment.Tasks, knowledgeItems)
	case "inspection_report":
		articleID, err := strconv.Atoi(targetID)
		if err != nil || articleID <= 0 {
			result.Failed(c, 400, "巡检报告知识文章 ID 无效")
			return
		}
		var article knowledgeModel.KnowledgeBase
		if err = db.Db.First(&article, articleID).Error; err != nil {
			result.Failed(c, 404, "巡检报告知识文章不存在")
			return
		}
		title = "巡检报告分析"
		knowledgeItems = s.mergeKnowledgeItems(
			s.collectKnowledgeItems(strings.TrimSpace(article.Title+" "+article.Content+" "+article.Tags+" "+keyword)),
			append([]map[string]interface{}{
				{
					"id":       article.ID,
					"title":    article.Title,
					"type":     article.Type,
					"category": article.Category,
					"content":  article.Content,
					"tags":     article.Tags,
				},
			}, s.loadKnowledgeItemsByIDs(req.KnowledgeIDs)...),
		)
		renderedPrompt, systemPrompt = s.buildDiagnosisPrompt(firstNonEmpty(templateName, "incident_analysis"), map[string]interface{}{
			"reportTitle":      article.Title,
			"reportType":       article.Type,
			"reportCategory":   article.Category,
			"reportContent":    article.Content,
			"reportKeywords":   firstNonEmpty(article.Keywords, keyword),
			"knowledgeContext": stringify(knowledgeItems),
		})
		report = buildInspectionReportDiagnosis(article, knowledgeItems)
	default:
		result.Failed(c, 400, "暂不支持的诊断场景")
		return
	}

	s.persistDiagnosisHistory(c, sessionID, scene, title, req, renderedPrompt, report)
	result.Success(c, aiModel.AIDiagnosisResponse{
		SessionID:      sessionID,
		Scene:          scene,
		TargetID:       targetID,
		Title:          title,
		Report:         report,
		RenderedPrompt: renderedPrompt,
		SystemPrompt:   systemPrompt,
		UsedFallback:   usedFallback,
		KnowledgeItems: knowledgeItems,
	})
}

func (s *AIService) collectKnowledgeItems(keyword string) []map[string]interface{} {
	keyword = strings.TrimSpace(keyword)
	var list []knowledgeModel.KnowledgeBase
	query := db.Db.Model(&knowledgeModel.KnowledgeBase{}).Where("enabled = ?", 1)
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR content LIKE ? OR keywords LIKE ? OR tags LIKE ?", like, like, like, like)
	}
	if err := query.Order("score DESC, use_count DESC, update_time DESC").Limit(5).Find(&list).Error; err != nil {
		return []map[string]interface{}{}
	}

	items := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		items = append(items, map[string]interface{}{
			"id":       item.ID,
			"title":    item.Title,
			"type":     item.Type,
			"category": item.Category,
			"content":  item.Content,
			"tags":     item.Tags,
		})
	}
	return items
}

func (s *AIService) loadKnowledgeItemsByIDs(ids []uint) []map[string]interface{} {
	if len(ids) == 0 {
		return []map[string]interface{}{}
	}

	var list []knowledgeModel.KnowledgeBase
	if err := db.Db.Where("enabled = ? AND id IN ?", 1, ids).Order("score DESC, use_count DESC, update_time DESC").Find(&list).Error; err != nil {
		return []map[string]interface{}{}
	}

	items := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		items = append(items, map[string]interface{}{
			"id":       item.ID,
			"title":    item.Title,
			"type":     item.Type,
			"category": item.Category,
			"content":  item.Content,
			"tags":     item.Tags,
		})
		db.Db.Model(&knowledgeModel.KnowledgeBase{}).Where("id = ?", item.ID).UpdateColumn("use_count", gorm.Expr("use_count + ?", 1))
	}
	return items
}

func (s *AIService) mergeKnowledgeItems(groups ...[]map[string]interface{}) []map[string]interface{} {
	merged := make([]map[string]interface{}, 0)
	seen := make(map[string]struct{})
	for _, group := range groups {
		for _, item := range group {
			key := strings.TrimSpace(stringify(item["id"]))
			if key == "" {
				key = strings.TrimSpace(stringify(item["title"]))
			}
			if key == "" {
				continue
			}
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			merged = append(merged, item)
		}
	}
	return merged
}

type workloadCapacityDiagnosisContext struct {
	Snapshot           k8sModel.WorkloadCapacitySuggestionHistoryEntity
	AlertSummary       k8sModel.WorkloadCapacityAlertSummary
	Autoscaling        *k8sModel.WorkloadAutoscalingSummary
	RecommendedActions []k8sModel.WorkloadCapacityRecommendation
	RecommendedPolicy  k8sModel.WorkloadCapacityPolicy
	WatchMetrics       []string
	AlertCenterQuery   k8sModel.WorkloadCapacityAlertCenterQuery
}

func (s *AIService) loadWorkloadCapacityDiagnosisContext(historyID uint) (*workloadCapacityDiagnosisContext, error) {
	var snapshot k8sModel.WorkloadCapacitySuggestionHistoryEntity
	if err := db.Db.First(&snapshot, historyID).Error; err != nil {
		return nil, err
	}

	context := &workloadCapacityDiagnosisContext{
		Snapshot:           snapshot,
		RecommendedActions: []k8sModel.WorkloadCapacityRecommendation{},
		WatchMetrics:       []string{},
	}
	s.unmarshalJSONText(snapshot.AlertSummary, &context.AlertSummary)
	s.unmarshalJSONText(snapshot.Autoscaling, &context.Autoscaling)
	s.unmarshalJSONText(snapshot.RecommendedActions, &context.RecommendedActions)
	s.unmarshalJSONText(snapshot.RecommendedPolicy, &context.RecommendedPolicy)
	s.unmarshalJSONText(snapshot.WatchMetrics, &context.WatchMetrics)
	s.unmarshalJSONText(snapshot.AlertCenterQuery, &context.AlertCenterQuery)
	return context, nil
}

func (s *AIService) buildWorkloadCapacityKnowledgeKeyword(context *workloadCapacityDiagnosisContext) string {
	if context == nil {
		return ""
	}
	parts := []string{
		context.Snapshot.ClusterName,
		context.Snapshot.NamespaceName,
		context.Snapshot.WorkloadName,
		context.Snapshot.WorkloadType,
		context.Snapshot.RiskLevel,
		"hpa capacity autoscaling kubernetes diagnosis",
	}
	for _, item := range context.RecommendedActions {
		parts = append(parts, item.Action, item.Reason)
	}
	return strings.Join(parts, " ")
}

func (s *AIService) buildWorkloadCapacityDiagnosisPrompt(
	templateName string,
	context *workloadCapacityDiagnosisContext,
	knowledgeItems []map[string]interface{},
) (string, string) {
	if context == nil {
		return "", ""
	}

	templateName = strings.TrimSpace(templateName)
	if templateName == "" {
		templateName = "incident_analysis"
	}

	generatedAt := context.Snapshot.CreatedAt.Format("2006-01-02 15:04:05")
	alertSummaryText := fmt.Sprintf(
		"open=%d,resolved=%d,incidents=%d,critical=%d",
		context.AlertSummary.OpenEventCount,
		context.AlertSummary.ResolvedEventCount,
		context.AlertSummary.IncidentCount,
		context.AlertSummary.CriticalCount,
	)
	recentChanges := firstNonEmpty(
		fmt.Sprintf("容量建议快照 #%d 由 %s 生成于 %s", context.Snapshot.ID, firstNonEmpty(context.Snapshot.GeneratedBy, "unknown"), generatedAt),
		fmt.Sprintf("建议策略=%s", firstNonEmpty(context.RecommendedPolicy.Type, "none")),
	)

	switch templateName {
	case "k8s_event_analysis":
		return s.buildDiagnosisPrompt(templateName, map[string]interface{}{
			"clusterName":   context.Snapshot.ClusterName,
			"namespace":     context.Snapshot.NamespaceName,
			"resourceKind":  context.Snapshot.WorkloadType,
			"resourceName":  context.Snapshot.WorkloadName,
			"podName":       "-",
			"timeRange":     generatedAt,
			"statusSummary": fmt.Sprintf("risk=%s, followUp=%s", context.Snapshot.RiskLevel, context.Snapshot.FollowUpWindow),
			"eventSummary":  stringify(context.AlertSummary),
			"recentLogs":    context.Snapshot.Report,
			"yamlSnippet":   "capacity-suggestion snapshot only",
			"metricSummary": fmt.Sprintf("watchMetrics=%s; autoscaling=%s", stringify(context.WatchMetrics), stringify(context.Autoscaling)),
		})
	case "alert_triage":
		return s.buildDiagnosisPrompt(templateName, map[string]interface{}{
			"alertSource": "capacity-suggestion",
			"alertName":   context.Snapshot.WorkloadName + " 容量治理诊断",
			"alertLevel":  strings.ToUpper(firstNonEmpty(context.Snapshot.RiskLevel, "medium")),
			"alertStatus": func() string {
				if context.AlertSummary.OpenEventCount > 0 {
					return "open"
				}
				return "observed"
			}(),
			"targetName":         context.Snapshot.WorkloadName,
			"clusterName":        context.Snapshot.ClusterName,
			"namespace":          context.Snapshot.NamespaceName,
			"labels":             stringify(context.AlertCenterQuery),
			"annotations":        context.Snapshot.AlertCenterPath,
			"recentChanges":      recentChanges,
			"metricSnapshot":     stringify(context.WatchMetrics),
			"alertHistory":       alertSummaryText,
			"eventAndLogSummary": context.Snapshot.Report,
		})
	case "prediction_suggestion":
		return s.buildDiagnosisPrompt(templateName, map[string]interface{}{
			"targetType":   "k8s-workload",
			"targetName":   context.Snapshot.WorkloadName,
			"clusterName":  context.Snapshot.ClusterName,
			"namespace":    context.Snapshot.NamespaceName,
			"resourceName": context.Snapshot.WorkloadName,
			"metricName":   firstNonEmpty(context.RecommendedPolicy.Metric, "replicas"),
			"metricType":   firstNonEmpty(context.RecommendedPolicy.Type, "manual"),
			"currentValue": stringify(len(context.RecommendedActions)),
			"threshold":    firstNonEmpty(context.RecommendedPolicy.TargetUtilization, "-"),
			"trendType":    "capacity-suggestion-history",
			"changeRate":   "0",
			"daysLeft":     "0",
			"confidence":   context.Snapshot.RiskLevel,
			"resourceSpec": stringify(context.Autoscaling),
			"scalingState": stringify(context.RecommendedPolicy),
			"predictions":  stringify(context.AlertSummary),
		})
	default:
		return s.buildDiagnosisPrompt("incident_analysis", map[string]interface{}{
			"targetType":    "k8s-workload",
			"hostName":      "-",
			"hostIp":        "-",
			"clusterName":   context.Snapshot.ClusterName,
			"namespace":     context.Snapshot.NamespaceName,
			"resourceName":  context.Snapshot.WorkloadName,
			"category":      "capacity-governance",
			"detectTime":    generatedAt,
			"recentChanges": recentChanges,
			"symptoms":      fmt.Sprintf("risk=%s; alertSummary=%s; recommendedActions=%s", context.Snapshot.RiskLevel, alertSummaryText, stringify(context.RecommendedActions)),
			"evidenceData":  fmt.Sprintf("report=%s\n\nprompt=%s\n\nautoscaling=%s\n\nwatchMetrics=%s\n\nknowledgeContext=%s", context.Snapshot.Report, context.Snapshot.RenderedPrompt, stringify(context.Autoscaling), stringify(context.WatchMetrics), stringify(knowledgeItems)),
		})
	}
}

func (s *AIService) unmarshalJSONText(raw string, target interface{}) {
	if target == nil || strings.TrimSpace(raw) == "" || strings.TrimSpace(raw) == "null" {
		return
	}
	_ = json.Unmarshal([]byte(raw), target)
}

func (s *AIService) buildDiagnosisPrompt(templateName string, variables map[string]interface{}) (string, string) {
	if templateName == "" {
		return "", ""
	}
	var tmpl aiModel.PromptTemplate
	if err := db.Db.Where("name = ? AND enabled = 1", templateName).First(&tmpl).Error; err != nil {
		return "", ""
	}
	return renderPromptTemplate(tmpl.Template, variables), renderPromptTemplate(tmpl.SystemPrompt, variables)
}

func (s *AIService) persistDiagnosisHistory(c *gin.Context, sessionID, scene, title string, req aiModel.AIDiagnosisRequest, renderedPrompt, report string) {
	userID, _ := jwt.GetAdminId(c)
	entitiesJSON, _ := json.Marshal(req)
	now := util.HTime{Time: time.Now()}
	db.Db.Create(&aiModel.AIChatHistory{
		SessionID:  sessionID,
		UserID:     userID,
		Role:       "user",
		Message:    firstNonEmpty(req.Keyword, req.TargetID, title),
		Intent:     scene,
		IntentConf: 1,
		Entities:   string(entitiesJSON),
		TaskType:   scene,
		Status:     2,
		CreateTime: now,
	})
	db.Db.Create(&aiModel.AIChatHistory{
		SessionID:  sessionID,
		UserID:     userID,
		Role:       "assistant",
		Message:    report,
		Intent:     scene,
		IntentConf: 1,
		Entities:   string(entitiesJSON),
		TaskType:   scene,
		Status:     2,
		CreateTime: now,
	})
}

func buildTerminalAuditDiagnosis(detail systemModel.TerminalAuditSessionDetail, knowledgeItems []map[string]interface{}) string {
	lines := []string{
		"# 终端审计复盘报告",
		"",
		fmt.Sprintf("- 会话ID: %s", detail.Session.SessionID),
		fmt.Sprintf("- 终端类型: %s", firstNonEmpty(detail.Session.SessionType, "unknown")),
		fmt.Sprintf("- 风险等级: %d", detail.Session.RiskLevel),
		fmt.Sprintf("- 命令数: %d", detail.Session.CommandCount),
		fmt.Sprintf("- 敏感命令数: %d", detail.Session.SensitiveCommandCount),
		"",
		"## 关键命令",
	}
	for index, command := range detail.Commands {
		if index >= 5 {
			break
		}
		lines = append(lines, fmt.Sprintf("%d. %s", index+1, strings.TrimSpace(command.Command)))
	}
	lines = append(lines,
		"",
		"## 初步判断",
		"- 先结合命令时间线确认是否存在批量粘贴、危险命令或异常退出。",
		"- 若录像状态正常，建议继续比对输入输出事件流。",
		"",
		"## 推荐知识",
	)
	if len(knowledgeItems) == 0 {
		lines = append(lines, "- 暂未命中关联知识，建议补充 Runbook。")
	} else {
		for _, item := range knowledgeItems {
			lines = append(lines, fmt.Sprintf("- %s", stringify(item["title"])))
		}
	}
	return strings.Join(lines, "\n")
}

func buildSQLWorkOrderDiagnosis(order cmdbModel.CmdbSQLWorkOrder, knowledgeItems []map[string]interface{}) string {
	lines := []string{
		"# SQL 工单分析报告",
		"",
		fmt.Sprintf("- 工单号: %s", order.OrderNo),
		fmt.Sprintf("- 数据库: %s", order.DatabaseName),
		fmt.Sprintf("- 操作类型: %s", order.OperationType),
		fmt.Sprintf("- 风险等级: %d", order.RiskLevel),
		fmt.Sprintf("- 当前状态: %d", order.Status),
		"",
		"## 风险摘要",
		firstNonEmpty(order.RiskSummary, "无"),
		"",
		"## 回滚建议",
		firstNonEmpty(order.RollbackHint, order.RollbackSQL, "无"),
		"",
		"## 推荐知识",
	}
	if len(knowledgeItems) == 0 {
		lines = append(lines, "- 暂未命中关联知识，建议补充 SQL 变更 Runbook。")
	} else {
		for _, item := range knowledgeItems {
			lines = append(lines, fmt.Sprintf("- %s", stringify(item["title"])))
		}
	}
	return strings.Join(lines, "\n")
}

func buildWorkloadCapacityDiagnosis(context *workloadCapacityDiagnosisContext, knowledgeItems []map[string]interface{}) string {
	if context == nil {
		return ""
	}

	lines := []string{
		"# 工作负载容量治理诊断报告",
		"",
		fmt.Sprintf("- 快照ID: %d", context.Snapshot.ID),
		fmt.Sprintf("- 集群: %s", context.Snapshot.ClusterName),
		fmt.Sprintf("- 命名空间: %s", context.Snapshot.NamespaceName),
		fmt.Sprintf("- 工作负载: %s / %s", context.Snapshot.WorkloadType, context.Snapshot.WorkloadName),
		fmt.Sprintf("- 风险等级: %s", context.Snapshot.RiskLevel),
		fmt.Sprintf("- 建议复查窗口: %s", firstNonEmpty(context.Snapshot.FollowUpWindow, "-")),
		fmt.Sprintf("- 快照时间: %s", context.Snapshot.CreatedAt.Format("2006-01-02 15:04:05")),
		"",
		"## HPA 现状",
	}

	if context.Autoscaling == nil || !context.Autoscaling.Enabled {
		lines = append(lines, "- 当前未配置 HPA")
	} else {
		lines = append(lines,
			fmt.Sprintf("- HPA: %s", firstNonEmpty(context.Autoscaling.Name, context.Snapshot.WorkloadName)),
			fmt.Sprintf("- 副本范围: %d - %d", context.Autoscaling.MinReplicas, context.Autoscaling.MaxReplicas),
			fmt.Sprintf("- 当前/期望副本: %d/%d", context.Autoscaling.CurrentReplicas, context.Autoscaling.DesiredReplicas),
		)
		if len(context.Autoscaling.Warnings) > 0 {
			lines = append(lines, "- 治理提示:")
			for _, warning := range context.Autoscaling.Warnings {
				lines = append(lines, fmt.Sprintf("  - %s", warning))
			}
		}
	}

	lines = append(lines,
		"",
		"## 关联告警",
		fmt.Sprintf("- 未恢复自动化事件: %d", context.AlertSummary.OpenEventCount),
		fmt.Sprintf("- 已恢复自动化事件: %d", context.AlertSummary.ResolvedEventCount),
		fmt.Sprintf("- 关联 incident: %d", context.AlertSummary.IncidentCount),
		fmt.Sprintf("- 高优先级计数: %d", context.AlertSummary.CriticalCount),
		"",
		"## 建议动作",
	)
	if len(context.RecommendedActions) == 0 {
		lines = append(lines, "- 暂无建议动作")
	} else {
		for _, item := range context.RecommendedActions {
			lines = append(lines, fmt.Sprintf("- %s %s：%s；预期收益：%s", item.Priority, item.Action, item.Reason, item.ExpectedEffect))
		}
	}

	lines = append(lines, "", "## 持续观察")
	if len(context.WatchMetrics) == 0 {
		lines = append(lines, "- 暂无持续观察指标")
	} else {
		for _, metric := range context.WatchMetrics {
			lines = append(lines, fmt.Sprintf("- %s", metric))
		}
	}

	lines = append(lines, "", "## 原始容量建议")
	lines = append(lines, context.Snapshot.Report)

	lines = append(lines, "", "## 推荐知识")
	if len(knowledgeItems) == 0 {
		lines = append(lines, "- 暂未命中关联知识，建议补充 HPA/KEDA/容量治理 Runbook。")
	} else {
		for _, item := range knowledgeItems {
			lines = append(lines, fmt.Sprintf("- %s", stringify(item["title"])))
		}
	}
	return strings.Join(lines, "\n")
}

func buildKnowledgeSearchDiagnosis(keyword string, knowledgeItems []map[string]interface{}) string {
	lines := []string{"# 知识推荐结果", "", fmt.Sprintf("- 检索关键字: %s", firstNonEmpty(keyword, "-")), "", "## 推荐条目"}
	if len(knowledgeItems) == 0 {
		lines = append(lines, "- 暂无匹配知识。")
	} else {
		for _, item := range knowledgeItems {
			lines = append(lines, fmt.Sprintf("- %s", stringify(item["title"])))
		}
	}
	return strings.Join(lines, "\n")
}

func buildInspectionReportDiagnosis(article knowledgeModel.KnowledgeBase, knowledgeItems []map[string]interface{}) string {
	lines := []string{
		"# 巡检报告分析",
		"",
		fmt.Sprintf("- 巡检文章 ID: %d", article.ID),
		fmt.Sprintf("- 标题: %s", firstNonEmpty(article.Title, "-")),
		fmt.Sprintf("- 类型: %s", firstNonEmpty(article.Type, "-")),
		fmt.Sprintf("- 分类: %s", firstNonEmpty(article.Category, "-")),
		"",
		"## 报告摘要",
		firstNonEmpty(article.Content, "-"),
		"",
		"## 建议动作",
		"- 先确认巡检项中是否存在高风险告警、发布变更或资源异常。",
		"- 将需要人工处理的项沉淀为工单或补充到知识库 SOP。",
		"- 如涉及重复问题，建议同步建立 AI 诊断模板或固定排查脚本。",
		"",
		"## 推荐知识",
	}
	if len(knowledgeItems) == 0 {
		lines = append(lines, "- 暂未命中关联知识，建议补充巡检 Runbook 与历史案例。")
	} else {
		for _, item := range knowledgeItems {
			lines = append(lines, fmt.Sprintf("- %s", stringify(item["title"])))
		}
	}
	return strings.Join(lines, "\n")
}

func renderPromptTemplate(template string, variables map[string]interface{}) string {
	output := template
	for key, value := range variables {
		placeholder := "{{" + key + "}}"
		output = strings.ReplaceAll(output, placeholder, stringify(value))
	}
	return output
}

func stringify(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	default:
		data, _ := json.Marshal(v)
		return string(data)
	}
}

func generateSessionID() string {
	return time.Now().Format("20060102150405") + "-" + randomSuffix()
}

func randomSuffix() string {
	return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
