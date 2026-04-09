package service

import (
	"context"
	"strings"
	"time"

	aiModel "dodevops-api/api/ai/model"
	knowledgeModel "dodevops-api/api/knowledge/model"
	"dodevops-api/common/result"
	"dodevops-api/pkg/db"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type aiOverviewSnapshot struct {
	RuntimeEnabled       bool
	RuntimeProvider      string
	RuntimeModel         string
	ReasoningEffort      string
	RuntimeReachable     bool
	RuntimeLastError     string
	RuntimeCheckedAt     time.Time
	PromptTemplates      int64
	KnowledgeItems       int64
	DiagnosisSessions    int64
	AssistantSessions    int64
	InspectionTemplates  int64
	InspectionReports    int64
	PendingConfirmations int64
}

func (s *AIService) ListOverview(c *gin.Context) {
	userID, _ := jwt.GetAdminId(c)
	result.Success(c, buildAIOverviewResponse(s.collectAIOverviewSnapshot(c.Request.Context(), userID)))
}

func (s *AIService) collectAIOverviewSnapshot(ctx context.Context, userID uint) aiOverviewSnapshot {
	s.ensureBuiltinInspectionTemplates()

	runtimeClient := newAIRuntimeClient()
	snapshot := aiOverviewSnapshot{}
	if runtimeClient != nil {
		snapshot.RuntimeEnabled = runtimeClient.IsEnabled()
		snapshot.RuntimeProvider = runtimeClient.Provider()
		snapshot.RuntimeModel = runtimeClient.Model()
		snapshot.ReasoningEffort = runtimeClient.reasoningEffort
		if snapshot.RuntimeEnabled {
			probe := runtimeClient.Probe(ctx)
			snapshot.RuntimeReachable = probe.Reachable
			snapshot.RuntimeLastError = probe.LastError
			snapshot.RuntimeCheckedAt = probe.CheckedAt
		}
	}

	db.Db.Model(&aiModel.PromptTemplate{}).Where("enabled = ?", 1).Count(&snapshot.PromptTemplates)
	db.Db.Model(&knowledgeModel.KnowledgeBase{}).Where("enabled = ?", 1).Count(&snapshot.KnowledgeItems)
	db.Db.Model(&aiModel.AIInspectionTemplateEntity{}).Where("enabled = ?", 1).Count(&snapshot.InspectionTemplates)
	db.Db.Model(&aiModel.AIInspectionReportEntity{}).Where("user_id = ?", userID).Count(&snapshot.InspectionReports)
	db.Db.Model(&aiModel.AIAssistantConfirmationEntity{}).Where("user_id = ? AND status = ?", userID, "pending").Count(&snapshot.PendingConfirmations)

	snapshot.AssistantSessions = countDistinctAISessions(userID, assistantTaskType, false)
	snapshot.DiagnosisSessions = countDistinctAISessions(userID, assistantTaskType, true)
	return snapshot
}

func countDistinctAISessions(userID uint, taskType string, exclude bool) int64 {
	var count int64
	query := db.Db.Table("ai_agent_chat_history").Where("user_id = ?", userID)
	if strings.TrimSpace(taskType) != "" {
		if exclude {
			query = query.Where("task_type <> ?", taskType)
		} else {
			query = query.Where("task_type = ?", taskType)
		}
	}
	query.Distinct("session_id").Count(&count)
	return count
}

func buildDiagnosisSceneCatalog() []aiModel.AIDiagnosisSceneOverview {
	return []aiModel.AIDiagnosisSceneOverview{
		{
			Value:              "terminal_audit",
			Label:              "终端审计复盘",
			Description:        "复盘终端操作轨迹、敏感命令和审计上下文。",
			TargetLabel:        "目标 ID",
			TargetPlaceholder:  "请输入终端会话 ID",
			KeywordPlaceholder: "例如：terminal audit / kubectl exec",
			TemplateName:       "terminal_audit_review",
		},
		{
			Value:              "sql_work_order",
			Label:              "SQL 工单分析",
			Description:        "分析 SQL 风险、回滚建议和审批上下文。",
			TargetLabel:        "目标 ID",
			TargetPlaceholder:  "请输入 SQL 工单 ID",
			KeywordPlaceholder: "例如：sql rollback / ddl review",
			TemplateName:       "yaml_change_review",
		},
		{
			Value:              "alert_analysis",
			Label:              "告警事件分析",
			Description:        "聚合告警摘要、事件列表和知识上下文生成分析建议。",
			TargetLabel:        "告警关键词",
			TargetPlaceholder:  "请输入告警关键词，留空则分析当前摘要",
			KeywordPlaceholder: "例如：mysql / kubelet / 磁盘",
			TemplateName:       "incident_analysis",
		},
		{
			Value:              "deployment_review",
			Label:              "发布复盘分析",
			Description:        "基于快速发布记录和任务执行结果做发布复盘。",
			TargetLabel:        "发布 ID",
			TargetPlaceholder:  "请输入快速发布 ID",
			KeywordPlaceholder: "例如：rollback / 发布失败 / 环境治理",
			TemplateName:       "incident_analysis",
		},
		{
			Value:              "inspection_report",
			Label:              "巡检报告分析",
			Description:        "基于巡检知识文章沉淀风险总结和治理建议。",
			TargetLabel:        "知识文章 ID",
			TargetPlaceholder:  "请输入巡检知识文章 ID",
			KeywordPlaceholder: "例如：inspection report / 巡检 runbook",
			TemplateName:       "incident_analysis",
		},
		{
			Value:              "workload_capacity",
			Label:              "容量治理",
			Description:        "分析 HPA、容量建议快照和工作负载风险。",
			TargetLabel:        "快照 ID",
			TargetPlaceholder:  "请输入容量建议快照 ID",
			KeywordPlaceholder: "例如：demo-nginx hpa capacity autoscaling",
			TemplateName:       "incident_analysis",
		},
		{
			Value:              "knowledge_search",
			Label:              "知识检索增强",
			Description:        "按关键词召回知识库上下文并形成诊断建议。",
			TargetLabel:        "知识文章 ID",
			TargetPlaceholder:  "可选：输入知识文章 ID",
			KeywordPlaceholder: "例如：postgres rollback knowledge",
			TemplateName:       "incident_analysis",
		},
	}
}

func buildAIWorkspaceDomains() []aiModel.AIOverviewDomain {
	return []aiModel.AIOverviewDomain{
		{
			Key:         "host",
			Label:       "主机巡检与只读排障",
			Description: "查询主机、执行只读命令、生成巡检报告。",
			Prompt:      "查询主机 10.0.0.200 并生成基础巡检报告",
			Route:       "/ai/assistant",
		},
		{
			Key:         "k8s",
			Label:       "K8s 工作负载诊断",
			Description: "结合集群、命名空间和工作负载上下文生成分析结果。",
			Prompt:      "查看 demo-nginx 的工作负载治理快照并给出建议",
			Route:       "/ai/assistant",
		},
		{
			Key:         "alert",
			Label:       "告警分诊",
			Description: "聚合告警摘要与事件列表，输出分诊和处置建议。",
			Prompt:      "分析当前告警摘要并给出处置建议",
			Route:       "/ai/diagnosis?scene=alert_analysis",
		},
		{
			Key:         "deployment",
			Label:       "发布复盘",
			Description: "基于发布任务和执行结果分析失败原因与回滚策略。",
			Prompt:      "复盘快速发布 1 的执行结果并给出回滚建议",
			Route:       "/ai/diagnosis?scene=deployment_review",
		},
		{
			Key:         "knowledge",
			Label:       "知识增强",
			Description: "把知识库、历史报告和 AI 诊断串成可复用经验。",
			Prompt:      "结合知识库总结最近一次巡检风险",
			Route:       "/ai/diagnosis?scene=knowledge_search",
		},
	}
}

func buildAIQuickPrompts() []string {
	return []string{
		"查询主机 10.0.0.200",
		"查看主机 10.0.0.200 的磁盘占用",
		"在主机 10.0.0.200 执行 `free -m`",
		"为主机 10.0.0.200 生成巡检报告",
		"分析当前告警摘要并给出处置建议",
		"结合知识库总结最近一次巡检风险",
	}
}

func buildAIOverviewResponse(snapshot aiOverviewSnapshot) aiModel.AIOverviewResponse {
	runtimeStatus, runtimeStatusText := resolveAIOverviewRuntimeStatus(snapshot)

	return aiModel.AIOverviewResponse{
		Runtime: aiModel.AIOverviewRuntime{
			Enabled:         snapshot.RuntimeEnabled,
			Provider:        firstNonEmpty(snapshot.RuntimeProvider, "openai"),
			Model:           firstNonEmpty(snapshot.RuntimeModel, "gpt-5.4"),
			ReasoningEffort: firstNonEmpty(snapshot.ReasoningEffort, "medium"),
			Status:          runtimeStatus,
			StatusText:      runtimeStatusText,
			LastError:       strings.TrimSpace(snapshot.RuntimeLastError),
			CheckedAt:       formatAIOverviewCheckedAt(snapshot.RuntimeCheckedAt),
		},
		Stats: aiModel.AIOverviewStats{
			PromptTemplates:      snapshot.PromptTemplates,
			KnowledgeItems:       snapshot.KnowledgeItems,
			DiagnosisSessions:    snapshot.DiagnosisSessions,
			AssistantSessions:    snapshot.AssistantSessions,
			InspectionTemplates:  snapshot.InspectionTemplates,
			InspectionReports:    snapshot.InspectionReports,
			PendingConfirmations: snapshot.PendingConfirmations,
		},
		DiagnosisScenes: buildDiagnosisSceneCatalog(),
		Domains:         buildAIWorkspaceDomains(),
		QuickPrompts:    buildAIQuickPrompts(),
	}
}

func resolveAIOverviewRuntimeStatus(snapshot aiOverviewSnapshot) (string, string) {
	if !snapshot.RuntimeEnabled {
		return "fallback", "未配置可用大模型密钥，当前工作台将优先使用内置回退逻辑。"
	}
	if snapshot.RuntimeReachable {
		return "ready", "大模型链路可用，当前工作台可调用实时 LLM 规划与总结能力。"
	}
	return "degraded", "大模型已配置，但最近一次实时探测失败，当前将回退到内置逻辑。"
}

func formatAIOverviewCheckedAt(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Local().Format("2006-01-02 15:04:05")
}
