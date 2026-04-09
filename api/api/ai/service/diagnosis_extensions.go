package service

import (
	"fmt"
	"strings"

	appModel "dodevops-api/api/app/model"
	monitorModel "dodevops-api/api/monitor/model"
)

func buildAlertAnalysisDiagnosis(summary monitorModel.MonitorAlertSummary, incidents []monitorModel.MonitorIncident, knowledgeItems []map[string]interface{}) string {
	lines := []string{
		"# 告警事件分析",
		"",
		fmt.Sprintf("- Open incidents: %d", summary.OpenIncidents),
		fmt.Sprintf("- Processing incidents: %d", summary.ProcessingIncidents),
		fmt.Sprintf("- Critical webhook logs: %d", summary.CriticalWebhookLogs),
		fmt.Sprintf("- Latest alert time: %s", firstNonEmpty(summary.LatestAlertTime, "-")),
		"",
		"## 最近事件",
	}

	if len(incidents) == 0 {
		lines = append(lines, "- 当前未命中具体告警事件，建议结合关键词或当前工作台上下文继续筛选。")
	} else {
		for _, item := range incidents {
			lines = append(lines, fmt.Sprintf("- [%s] %s / %s", firstNonEmpty(item.AlertLevel, "unknown"), firstNonEmpty(item.BusinessLine, "-"), firstNonEmpty(item.AlertDesc, "-")))
		}
	}

	if len(knowledgeItems) > 0 {
		lines = append(lines, "", "## 关联知识")
		for _, item := range knowledgeItems {
			lines = append(lines, fmt.Sprintf("- %s", stringify(item["title"])))
		}
	}

	lines = append(lines, "", "## 建议", "- 先确认是否存在持续中的高等级告警。", "- 再结合最近事件与知识条目判断是噪声、重复告警，还是正在扩大的故障。", "- 若涉及发布、工单或工作负载变更，建议继续联动 AI 助手做上下文复盘。")
	return strings.Join(lines, "\n")
}

func buildDeploymentReviewDiagnosis(deployment appModel.QuickDeployment, tasks []appModel.QuickDeploymentTask, knowledgeItems []map[string]interface{}) string {
	lines := []string{
		"# 发布复盘分析",
		"",
		fmt.Sprintf("- 发布标题: %s", firstNonEmpty(deployment.Title, "-")),
		fmt.Sprintf("- 发布状态: %s", deploymentStatusText(deployment.Status)),
		fmt.Sprintf("- 创建时间: %s", deployment.CreatedAt.Format("2006-01-02 15:04:05")),
		fmt.Sprintf("- 描述: %s", firstNonEmpty(deployment.Description, "-")),
		"",
		"## 任务执行",
	}

	if len(tasks) == 0 {
		lines = append(lines, "- 当前发布没有关联任务记录。")
	} else {
		for _, task := range tasks {
			taskLine := fmt.Sprintf("- %s / %s / %s", firstNonEmpty(task.AppName, "-"), firstNonEmpty(task.Environment, "-"), deploymentTaskStatusText(task.Status))
			if strings.TrimSpace(task.ErrorMessage) != "" {
				taskLine += fmt.Sprintf(" / 错误: %s", task.ErrorMessage)
			}
			lines = append(lines, taskLine)
		}
	}

	if len(knowledgeItems) > 0 {
		lines = append(lines, "", "## 关联知识")
		for _, item := range knowledgeItems {
			lines = append(lines, fmt.Sprintf("- %s", stringify(item["title"])))
		}
	}

	lines = append(lines, "", "## 建议", "- 优先核对失败任务的错误信息与执行顺序。", "- 若存在环境治理或容量建议信息，继续结合工作负载治理上下文判断是否为资源/发布窗口问题。", "- 若需要继续排查，可联动 AI 助手查看相关主机、工作负载或工单上下文。")
	return strings.Join(lines, "\n")
}

func deploymentStatusText(status int) string {
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
		return fmt.Sprintf("未知状态(%d)", status)
	}
}

func deploymentTaskStatusText(status int) string {
	switch status {
	case 1:
		return "未部署"
	case 2:
		return "部署中"
	case 3:
		return "成功"
	case 4:
		return "异常"
	default:
		return fmt.Sprintf("未知状态(%d)", status)
	}
}
