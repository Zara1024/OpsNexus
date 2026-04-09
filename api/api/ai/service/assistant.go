package service

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	aiModel "dodevops-api/api/ai/model"
	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	cmdbService "dodevops-api/api/cmdb/service"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"dodevops-api/pkg/db"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const assistantTaskType = "assistant_chat"

var (
	assistantHostIDPattern  = regexp.MustCompile(`(?i)(?:主机|机器|服务器|host)\s*(?:id)?\s*(\d{1,8})`)
	assistantIPPattern      = regexp.MustCompile(`\b\d{1,3}(?:\.\d{1,3}){3}\b`)
	assistantCommandPattern = regexp.MustCompile("`([^`]+)`")
	assistantNamedPatterns  = []*regexp.Regexp{
		regexp.MustCompile(`(?:主机|机器|服务器|host)\s*[:：]?\s*([A-Za-z0-9._:-]+)`),
		regexp.MustCompile(`([A-Za-z0-9][A-Za-z0-9._:-]{2,})\s*(?:主机|机器|服务器)`),
	}
	assistantDangerousFragments = []string{
		";", "&&", "||", "|", ">", "<", "`", "$(", " rm ", "shutdown", "reboot",
		"mkfs", "dd ", "chmod ", "chown ", "useradd", "userdel", "passwd ",
		"systemctl restart", "systemctl stop", "service restart", "service stop",
		"kill ", "killall ", "pkill ", "docker rm", "kubectl delete", "rollback",
	}
)

func (s *AIService) ChatAssistant(c *gin.Context, req aiModel.AIAssistantChatRequest) {
	message := strings.TrimSpace(req.Message)
	if message == "" {
		result.Failed(c, 400, "请输入对话内容")
		return
	}

	sessionID := strings.TrimSpace(req.SessionID)
	if sessionID == "" {
		sessionID = generateSessionID()
	}

	runtimeClient := newAIRuntimeClient()
	intent := detectAssistantIntent(message)
	planUsed := false
	fallbackReason := ""
	var llmPlan *assistantLLMPlan
	if runtimeClient != nil && runtimeClient.IsEnabled() {
		plan, err := s.planAssistantActionWithLLM(c.Request.Context(), runtimeClient, message)
		if err != nil {
			fallbackReason = err.Error()
		} else if plan != nil {
			llmPlan = plan
			planUsed = true
			if plan.Intent != "" {
				intent = plan.Intent
			}
		}
	}

	userID, _ := jwt.GetAdminId(c)
	response := aiModel.AIAssistantChatResponse{
		SessionID:      sessionID,
		Intent:         intent,
		Provider:       runtimeClient.Provider(),
		Model:          runtimeClient.Model(),
		UsedLLM:        planUsed,
		FallbackReason: fallbackReason,
		Suggestions:    assistantDefaultSuggestions(),
		Actions:        []aiModel.AIAssistantAction{},
		HostMatches:    []aiModel.AIAssistantHost{},
		ToolSteps:      []aiModel.AIAssistantToolStep{},
	}
	if llmPlan != nil && len(llmPlan.Suggestions) > 0 {
		response.Suggestions = llmPlan.Suggestions
	}
	response = s.attachAssistantDefaults(response, sessionID, userID, intent)

	entities := map[string]interface{}{
		"message": message,
		"intent":  intent,
	}
	s.persistAssistantEntry(userID, sessionID, "user", message, intent, entities)

	executionMessage := composeAssistantExecutionMessage(message, llmPlan)
	switch intent {
	case "risky_action":
		response = s.handleRiskyIntent(c, response, sessionID, userID, executionMessage)
	case "template_center":
		response = s.handleTemplateIntent(response, message)
	case "report_archive":
		response = s.handleReportArchiveIntent(response, userID)
	case "alert_analysis":
		response = s.handleAlertIntent(response, message)
	case "workload_lookup":
		response = s.handleWorkloadIntent(response, message)
	case "work_order_lookup":
		response = s.handleWorkOrderIntent(response, message)
	case "deployment_lookup":
		response = s.handleDeploymentIntent(response, message)
	case "inspection_report":
		response = s.handleInspectionIntent(c, response, executionMessage)
	case "host_command":
		response = s.handleHostCommandIntent(c, response, executionMessage)
	case "host_list":
		response = s.handleHostListIntent(response, message)
	case "host_lookup":
		response = s.handleHostLookupIntent(response, executionMessage)
	default:
		response.AssistantMessage = strings.Join([]string{
			"我可以帮你做这些事情：",
			"1. 查询主机、查看资产、执行只读命令。",
			"2. 做智能巡检并沉淀正式巡检报告。",
			"3. 查询告警、工作负载、工单和发布信息。",
			"4. 对高风险动作进入确认流，而不是直接执行。",
			"",
			"你可以直接试这些说法：",
			"- 查询主机 10.0.0.200",
			"- 查看这台机器的内存情况",
			"- 查看当前告警摘要",
			"- 查看 demo-nginx 的工作负载治理快照",
			"- 使用 Linux 基础巡检为这台主机生成巡检报告",
		}, "\n")
	}

	if llmPlan != nil && strings.TrimSpace(llmPlan.AssistantMessage) != "" && response.CommandResult == nil && response.InspectionResult == nil {
		response.AssistantMessage = strings.TrimSpace(llmPlan.AssistantMessage)
	}
	if runtimeClient != nil && runtimeClient.IsEnabled() {
		response = s.refineAssistantResponseWithLLM(c.Request.Context(), runtimeClient, message, response)
	}
	if response.Context != nil {
		response.Context.LastIntent = response.Intent
		s.saveAssistantContext(sessionID, userID, response.Context)
	}
	if response.InspectionResult != nil {
		s.saveInspectionReport(sessionID, userID, response.InspectionResult)
		response.RecentReports = s.listRecentInspectionReports(userID, 8)
	}
	if response.AssistantMessage == "" {
		response.AssistantMessage = "本次没有解析出明确动作，请换一种更具体的说法，例如主机 IP、主机名、告警关键字或巡检诉求。"
	}

	entities["response"] = map[string]interface{}{
		"intent":              response.Intent,
		"provider":            response.Provider,
		"model":               response.Model,
		"usedLlm":             response.UsedLLM,
		"context":             response.Context,
		"toolSteps":           response.ToolSteps,
		"pendingConfirmation": response.PendingConfirmation,
		"inspection":          response.InspectionResult,
	}
	s.persistAssistantEntry(userID, sessionID, "assistant", buildAssistantHistoryMessage(response), response.Intent, entities)

	result.Success(c, response)
}

func (s *AIService) ListAssistantHistory(c *gin.Context) {
	userID, _ := jwt.GetAdminId(c)
	type sessionRow struct {
		SessionID    string    `json:"sessionId"`
		LatestTime   time.Time `json:"latestTime"`
		MessageCount int64     `json:"messageCount"`
	}
	var rows []sessionRow
	if err := db.Db.Table("ai_agent_chat_history").
		Select("session_id, MAX(create_time) AS latest_time, COUNT(*) AS message_count").
		Where("user_id = ? AND task_type = ?", userID, assistantTaskType).
		Group("session_id").
		Order("latest_time DESC").
		Limit(20).
		Scan(&rows).Error; err != nil {
		result.Failed(c, 500, "获取 AI 助手历史失败: "+err.Error())
		return
	}
	result.Success(c, rows)
}

func (s *AIService) GetAssistantHistoryDetail(c *gin.Context, sessionID string) {
	userID, _ := jwt.GetAdminId(c)
	var rows []aiModel.AIChatHistory
	if err := db.Db.Where("session_id = ? AND user_id = ? AND task_type = ?", sessionID, userID, assistantTaskType).
		Order("create_time ASC, id ASC").
		Find(&rows).Error; err != nil {
		result.Failed(c, 500, "获取 AI 助手会话详情失败: "+err.Error())
		return
	}
	result.Success(c, rows)
}

func (s *AIService) handleHostListIntent(response aiModel.AIAssistantChatResponse, message string) aiModel.AIAssistantChatResponse {
	hostDao := cmdbDao.NewCmdbHostDao()
	lower := strings.ToLower(message)

	var hosts []cmdbModel.CmdbHost
	var title string
	switch {
	case strings.Contains(lower, "认证失败"):
		hosts = hostDao.GetCmdbHostsByStatus(3)
		title = "认证失败主机"
	case strings.Contains(lower, "未认证") || strings.Contains(lower, "待认证"):
		hosts = hostDao.GetCmdbHostsByStatus(2)
		title = "待认证主机"
	case strings.Contains(lower, "已认证") || strings.Contains(lower, "在线"):
		hosts = hostDao.GetCmdbHostsByStatus(1)
		title = "已认证主机"
	default:
		pageHosts, total := hostDao.GetCmdbHostListWithPage(1, 8)
		allHosts := hostDao.GetCmdbHostList()
		successCount, pendingCount, failedCount := summarizeHostStatuses(allHosts)
		response.HostMatches = toAssistantHosts(pageHosts)
		appendAssistantToolStep(&response, "host_list", "success", fmt.Sprintf("加载了 %d 台主机概览", total), "")
		response.AssistantMessage = fmt.Sprintf("当前纳管主机共 %d 台，其中已认证 %d 台、待认证 %d 台、认证失败 %d 台。下面先展示最近 8 台机器。", total, successCount, pendingCount, failedCount)
		response.Actions = []aiModel.AIAssistantAction{
			{Label: "查看待认证主机", Message: "查看待认证主机列表", Kind: "list"},
			{Label: "查看认证失败主机", Message: "查看认证失败主机列表", Kind: "list"},
		}
		return response
	}

	response.HostMatches = toAssistantHosts(limitHosts(hosts, 8))
	appendAssistantToolStep(&response, "host_list", "success", fmt.Sprintf("%s共 %d 台", title, len(hosts)), "")
	response.AssistantMessage = fmt.Sprintf("%s共找到 %d 台，下面展示前 %d 台。", title, len(hosts), len(response.HostMatches))
	response.Actions = []aiModel.AIAssistantAction{{Label: "查看全部主机概览", Message: "查看主机列表", Kind: "list"}}
	return response
}

func (s *AIService) handleHostLookupIntent(response aiModel.AIAssistantChatResponse, message string) aiModel.AIAssistantChatResponse {
	hosts, hostHint := resolveAssistantHostsWithContext(message, response.Context)
	if len(hosts) == 0 {
		appendAssistantToolStep(&response, "host_lookup", "failed", "未命中主机", "")
		response.AssistantMessage = "还没有定位到目标主机，请直接给我主机 IP，或者像 `主机 prod-api-01` 这样的名称。"
		response.Actions = []aiModel.AIAssistantAction{
			{Label: "查看主机列表", Message: "查看主机列表", Kind: "list"},
			{Label: "查看待认证主机", Message: "查看待认证主机列表", Kind: "list"},
		}
		return response
	}

	response.HostMatches = toAssistantHosts(limitHosts(hosts, 6))
	appendAssistantToolStep(&response, "host_lookup", "success", fmt.Sprintf("命中 %d 台主机", len(hosts)), firstNonEmpty(hostHint, message))
	if len(hosts) > 1 {
		response.AssistantMessage = fmt.Sprintf("根据 `%s` 找到了 %d 台候选主机。你可以继续指定主机 IP 或主机名，或者直接点下面的动作继续巡检。", firstNonEmpty(hostHint, message), len(hosts))
		return response
	}

	host := response.HostMatches[0]
	response.AssistantMessage = fmt.Sprintf("已定位到主机 `%s`，归属 `%s`，管理地址 `%s`。如果你愿意，我可以继续帮你查看磁盘、内存、端口，或者直接生成巡检报告。", host.HostName, firstNonEmpty(host.GroupName, "未分组"), firstNonEmpty(host.SSHIP, host.PrivateIP, host.PublicIP))
	response.Actions = assistantHostActions(assistantHostReference(host.HostName, host.SSHIP, host.PrivateIP, host.PublicIP))
	enrichContextWithHost(&response, &host)
	return response
}

func (s *AIService) handleHostCommandIntent(c *gin.Context, response aiModel.AIAssistantChatResponse, message string) aiModel.AIAssistantChatResponse {
	hosts, _ := resolveAssistantHostsWithContext(message, response.Context)
	if len(hosts) == 0 {
		appendAssistantToolStep(&response, "host_command", "failed", "缺少目标主机", "")
		response.AssistantMessage = "我识别到了命令诉求，但还缺少目标主机。请补充主机 IP 或主机名。"
		return response
	}
	if len(hosts) > 1 {
		response.HostMatches = toAssistantHosts(limitHosts(hosts, 6))
		appendAssistantToolStep(&response, "host_command", "failed", "命中多台主机", "")
		response.AssistantMessage = "命中了多台主机，先帮你列出来了。请继续指定具体主机 IP 或主机名后，我再执行命令。"
		return response
	}

	host := hosts[0]
	hostReference := assistantHostReference(host.HostName, host.SSHIP, host.PrivateIP, host.PublicIP)
	command := extractAssistantCommand(message)
	if command == "" {
		command = inferReadonlyCommand(message)
	}
	if command == "" {
		response.HostMatches = toAssistantHosts([]cmdbModel.CmdbHost{host})
		appendAssistantToolStep(&response, "host_command", "failed", "没有识别到命令", "")
		response.AssistantMessage = fmt.Sprintf("已经定位到主机 `%s`，但还没有识别到具体命令。你可以直接说“查看这台机器的磁盘占用”或“在主机 %s 执行 `df -h`”。", host.HostName, hostReference)
		response.Actions = assistantHostActions(hostReference)
		return response
	}

	if isRiskyAssistantAction(message) || !isSafeReadonlyCommand(command) {
		userID, _ := jwt.GetAdminId(c)
		return s.handleRiskyIntent(c, response, response.SessionID, userID, fmt.Sprintf("host %d `%s`", host.ID, command))
	}

	output, err := cmdbService.GetCmdbHostSSHService().ExecuteCommand(c, host.ID, command)
	response.HostMatches = toAssistantHosts([]cmdbModel.CmdbHost{host})
	if err != nil {
		appendAssistantToolStep(&response, "host_command", "failed", "命令执行失败", err.Error())
		response.AssistantMessage = fmt.Sprintf("在主机 `%s` 执行 `%s` 失败：%s", host.HostName, command, err.Error())
		response.Actions = assistantHostActions(hostReference)
		return response
	}

	appendAssistantToolStep(&response, "host_command", "success", fmt.Sprintf("已执行命令 `%s`", command), "")
	response.CommandResult = &aiModel.AIAssistantCommandResult{
		HostID:   host.ID,
		HostName: host.HostName,
		Command:  command,
		Output:   output.Output,
		Success:  true,
	}
	response.AssistantMessage = fmt.Sprintf("已在主机 `%s` 执行只读命令 `%s`。如果需要，我还可以基于这台机器继续生成巡检报告。", host.HostName, command)
	response.Actions = assistantHostActions(hostReference)
	if len(response.HostMatches) > 0 {
		enrichContextWithHost(&response, &response.HostMatches[0])
	}
	return response
}

func (s *AIService) handleInspectionIntent(c *gin.Context, response aiModel.AIAssistantChatResponse, message string) aiModel.AIAssistantChatResponse {
	hosts, _ := resolveAssistantHostsWithContext(message, response.Context)
	if len(hosts) == 0 {
		appendAssistantToolStep(&response, "inspection", "failed", "缺少目标主机", "")
		response.AssistantMessage = "要生成巡检报告，还需要你告诉我具体是哪台主机。请补充主机 IP 或主机名。"
		return response
	}
	if len(hosts) > 1 {
		response.HostMatches = toAssistantHosts(limitHosts(hosts, 6))
		appendAssistantToolStep(&response, "inspection", "failed", "命中多台主机", "")
		response.AssistantMessage = "命中了多台主机，先帮你列出来了。请继续指定具体主机 IP 或主机名，我就开始巡检。"
		return response
	}

	host := hosts[0]
	template := s.getInspectionTemplateForMessage(message)
	checkDefs := []aiModel.AIAssistantInspectionTemplateCheck{}
	if template != nil {
		checkDefs = template.Checks
	}
	if len(checkDefs) == 0 {
		checkDefs = []aiModel.AIAssistantInspectionTemplateCheck{
			{Name: "主机名", Command: "hostname"},
			{Name: "系统时间", Command: "date"},
			{Name: "系统负载", Command: "uptime"},
			{Name: "磁盘占用", Command: "df -h"},
			{Name: "内存情况", Command: "free -m"},
			{Name: "监听端口", Command: "ss -lntp"},
		}
	}

	appendAssistantToolStep(&response, "inspection_template", "success", fmt.Sprintf("使用模板 %s", func() string {
		if template != nil {
			return template.Name
		}
		return "Linux 基础巡检"
	}()), "")

	checks := make([]aiModel.AIAssistantInspectionCheck, 0, len(checkDefs))
	for _, item := range checkDefs {
		check := aiModel.AIAssistantInspectionCheck{
			Name:    item.Name,
			Command: item.Command,
			Status:  "success",
		}
		output, err := cmdbService.GetCmdbHostSSHService().ExecuteCommand(c, host.ID, item.Command)
		if err != nil {
			check.Status = "failed"
			check.Output = err.Error()
		} else {
			check.Output = strings.TrimSpace(output.Output)
		}
		checks = append(checks, check)
	}

	summary, report := buildInspectionSummary(host, checks)
	response.HostMatches = toAssistantHosts([]cmdbModel.CmdbHost{host})
	response.InspectionResult = &aiModel.AIAssistantInspectionResult{
		HostID:   host.ID,
		HostName: host.HostName,
		TemplateID: func() uint {
			if template != nil {
				return template.ID
			}
			return 0
		}(),
		TemplateName: func() string {
			if template != nil {
				return template.Name
			}
			return "Linux 基础巡检"
		}(),
		Summary: summary,
		Report:  report,
		Checks:  checks,
	}
	appendAssistantToolStep(&response, "inspection_execute", "success", fmt.Sprintf("完成 %d 个巡检检查项", len(checks)), "")
	response.AssistantMessage = fmt.Sprintf("已完成主机 `%s` 的智能巡检，并生成巡检报告。你可以继续让我查看某一项指标，或者换一台机器继续巡检。", host.HostName)
	response.Actions = assistantHostActions(assistantHostReference(host.HostName, host.SSHIP, host.PrivateIP, host.PublicIP))
	if len(response.HostMatches) > 0 {
		enrichContextWithHost(&response, &response.HostMatches[0])
	}
	return response
}

func (s *AIService) persistAssistantEntry(userID uint, sessionID, role, message, intent string, entities map[string]interface{}) {
	payload, _ := json.Marshal(entities)
	db.Db.Create(&aiModel.AIChatHistory{
		SessionID:  sessionID,
		UserID:     userID,
		Role:       role,
		Message:    message,
		Intent:     intent,
		IntentConf: 1,
		Entities:   string(payload),
		TaskType:   assistantTaskType,
		Status:     2,
		CreateTime: util.HTime{Time: time.Now()},
	})
}

func resolveAssistantHosts(message string) ([]cmdbModel.CmdbHost, string) {
	hostDao := cmdbDao.NewCmdbHostDao()
	groupDao := cmdbDao.NewCmdbGroupDao()

	if matches := assistantHostIDPattern.FindStringSubmatch(message); len(matches) > 1 {
		if id, err := strconv.Atoi(matches[1]); err == nil && id > 0 {
			host, err := hostDao.GetCmdbHostById(uint(id))
			if err == nil {
				host.Group, _ = groupDao.GetCmdbGroupById(host.GroupID)
				return []cmdbModel.CmdbHost{host}, matches[1]
			}
		}
	}

	if ip := assistantIPPattern.FindString(message); ip != "" {
		return hostDao.GetCmdbHostsByIP(ip), ip
	}

	for _, pattern := range assistantNamedPatterns {
		if matches := pattern.FindStringSubmatch(message); len(matches) > 1 {
			keyword := strings.Trim(matches[1], "，。,.：: ")
			if keyword != "" && !strings.EqualFold(keyword, "id") {
				return hostDao.GetCmdbHostsByHostNameLike(keyword), keyword
			}
		}
	}

	return nil, ""
}

func toAssistantHosts(hosts []cmdbModel.CmdbHost) []aiModel.AIAssistantHost {
	items := make([]aiModel.AIAssistantHost, 0, len(hosts))
	for _, host := range hosts {
		groupName := strings.TrimSpace(host.Group.Name)
		items = append(items, aiModel.AIAssistantHost{
			ID:          host.ID,
			HostName:    firstNonEmpty(host.HostName, host.Name),
			GroupName:   groupName,
			PrivateIP:   host.PrivateIP,
			PublicIP:    host.PublicIP,
			SSHIP:       host.SSHIP,
			OS:          host.OS,
			CPU:         host.CPU,
			Memory:      host.Memory,
			Disk:        host.Disk,
			Status:      host.Status,
			StatusText:  assistantStatusText(host.Status),
			SupportsSSH: strings.TrimSpace(host.SSHIP) != "" && strings.TrimSpace(host.SSHName) != "",
		})
	}
	return items
}

func assistantStatusText(status int) string {
	switch status {
	case 1:
		return "已认证"
	case 2:
		return "待认证"
	case 3:
		return "认证失败"
	default:
		return "未知状态"
	}
}

func summarizeHostStatuses(hosts []cmdbModel.CmdbHost) (int, int, int) {
	var successCount, pendingCount, failedCount int
	for _, host := range hosts {
		switch host.Status {
		case 1:
			successCount++
		case 2:
			pendingCount++
		case 3:
			failedCount++
		}
	}
	return successCount, pendingCount, failedCount
}

func assistantDefaultSuggestions() []string {
	return []string{
		"查询主机 10.0.0.200",
		"查看当前告警摘要",
		"查看 demo-nginx 的工作负载治理快照",
		"查看最近发布记录",
		"使用 Linux 基础巡检为这台主机生成巡检报告",
	}
}

func assistantHostReference(hostName, sshIP, privateIP, publicIP string) string {
	return firstNonEmpty(sshIP, privateIP, publicIP, hostName)
}

func assistantHostActions(hostReference string) []aiModel.AIAssistantAction {
	return []aiModel.AIAssistantAction{
		{Label: "查看磁盘", Message: fmt.Sprintf("查看主机 %s 的磁盘占用", hostReference), Kind: "command"},
		{Label: "查看内存", Message: fmt.Sprintf("查看主机 %s 的内存情况", hostReference), Kind: "command"},
		{Label: "查看端口", Message: fmt.Sprintf("查看主机 %s 的监听端口", hostReference), Kind: "command"},
		{Label: "生成巡检报告", Message: fmt.Sprintf("为主机 %s 生成巡检报告", hostReference), Kind: "inspection"},
	}
}

func detectAssistantIntent(message string) string {
	lower := strings.ToLower(message)
	switch {
	case isRiskyAssistantAction(message):
		return "risky_action"
	case assistantTemplatePattern.MatchString(message):
		return "template_center"
	case assistantReportArchivePattern.MatchString(message):
		return "report_archive"
	case assistantAlertPattern.MatchString(message):
		return "alert_analysis"
	case assistantWorkloadPattern.MatchString(message):
		return "workload_lookup"
	case assistantWorkOrderPattern.MatchString(message):
		return "work_order_lookup"
	case assistantDeploymentPattern.MatchString(message):
		return "deployment_lookup"
	case containsAny(lower, "巡检", "报告", "inspection", "health check", "检查"):
		return "inspection_report"
	case extractAssistantCommand(message) != "" || containsAny(lower, "执行", "运行", "命令", "exec", "磁盘", "disk", "内存", "memory", "负载", "load", "端口", "port", "监听", "系统", "os", "内核", "时间", "date"):
		return "host_command"
	case containsAny(lower, "列表", "list", "多少台", "概览", "汇总", "待认证", "认证失败", "已认证", "在线"):
		return "host_list"
	case assistantHostIDPattern.MatchString(message) || assistantIPPattern.MatchString(message) || containsAny(lower, "主机", "机器", "服务器", "host", "这台机器", "当前主机"):
		return "host_lookup"
	default:
		return "assistant_help"
	}
}

func extractAssistantCommand(message string) string {
	if matches := assistantCommandPattern.FindStringSubmatch(message); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func inferReadonlyCommand(message string) string {
	lower := strings.ToLower(message)
	switch {
	case containsAny(lower, "磁盘", "disk", "容量", "空间"):
		return "df -h"
	case containsAny(lower, "内存", "memory"):
		return "free -m"
	case containsAny(lower, "负载", "load", "uptime"):
		return "uptime"
	case containsAny(lower, "端口", "port", "监听"):
		return "ss -lntp"
	case containsAny(lower, "时间", "date", "clock"):
		return "date"
	case containsAny(lower, "系统", "os", "内核", "kernel"):
		return "uname -a"
	default:
		return ""
	}
}

func isSafeReadonlyCommand(command string) bool {
	trimmed := strings.TrimSpace(strings.ToLower(command))
	if trimmed == "" {
		return false
	}
	for _, fragment := range assistantDangerousFragments {
		if strings.Contains(trimmed, fragment) {
			return false
		}
	}
	allowedPrefixes := []string{
		"hostname",
		"date",
		"uptime",
		"df -h",
		"free -m",
		"ss -lntp",
		"uname -a",
		"cat /etc/os-release",
		"ps -ef",
		"docker ps",
		"docker images",
	}
	for _, prefix := range allowedPrefixes {
		if strings.HasPrefix(trimmed, prefix) {
			return true
		}
	}
	return false
}

func buildInspectionSummary(host cmdbModel.CmdbHost, checks []aiModel.AIAssistantInspectionCheck) (string, string) {
	riskItems := make([]string, 0)
	failedChecks := make([]string, 0)
	evidence := make([]string, 0, len(checks))

	for _, check := range checks {
		if check.Status != "success" {
			failedChecks = append(failedChecks, check.Name)
			evidence = append(evidence, fmt.Sprintf("### %s\n- 执行命令: `%s`\n- 结果: %s\n", check.Name, check.Command, check.Output))
			continue
		}
		if check.Name == "磁盘占用" {
			if usages := collectHighUsages(check.Output, 80); len(usages) > 0 {
				riskItems = append(riskItems, "磁盘分区使用率偏高: "+strings.Join(usages, ", "))
			}
		}
		evidence = append(evidence, fmt.Sprintf("### %s\n- 执行命令: `%s`\n```\n%s\n```\n", check.Name, check.Command, truncateForReport(check.Output, 1200)))
	}

	summaryParts := []string{
		fmt.Sprintf("主机 `%s` 巡检完成", firstNonEmpty(host.HostName, host.Name)),
		fmt.Sprintf("认证状态: %s", assistantStatusText(host.Status)),
	}
	if len(riskItems) == 0 {
		summaryParts = append(summaryParts, "当前未发现明显高风险项")
	} else {
		summaryParts = append(summaryParts, "发现关注项: "+strings.Join(riskItems, "；"))
	}
	if len(failedChecks) > 0 {
		summaryParts = append(summaryParts, "以下检查未完成: "+strings.Join(failedChecks, "、"))
	}

	lines := []string{
		"# 智能巡检报告",
		fmt.Sprintf("- 主机: %s", firstNonEmpty(host.HostName, host.Name)),
		fmt.Sprintf("- 巡检时间: %s", time.Now().Format("2006-01-02 15:04:05")),
		fmt.Sprintf("- 资产分组: %s", strings.TrimSpace(host.Group.Name)),
		fmt.Sprintf("- 管理地址: %s", firstNonEmpty(host.SSHIP, host.PrivateIP, host.PublicIP)),
		fmt.Sprintf("- 操作系统: %s", firstNonEmpty(host.OS, "未知")),
		"",
		"## 巡检结论",
		"- " + strings.Join(summaryParts, "；"),
	}
	if len(riskItems) > 0 {
		lines = append(lines, "", "## 风险提示")
		for _, item := range riskItems {
			lines = append(lines, "- "+item)
		}
	}
	if len(failedChecks) > 0 {
		lines = append(lines, "", "## 未完成检查")
		for _, item := range failedChecks {
			lines = append(lines, "- "+item)
		}
	}
	lines = append(lines, "", "## 巡检证据", strings.Join(evidence, "\n"))
	return strings.Join(summaryParts, "；"), strings.Join(lines, "\n")
}

func collectHighUsages(output string, threshold int) []string {
	lines := strings.Split(output, "\n")
	alerts := make([]string, 0)
	pattern := regexp.MustCompile(`(\d+)%`)
	for _, line := range lines {
		matches := pattern.FindStringSubmatch(line)
		if len(matches) < 2 {
			continue
		}
		usage, err := strconv.Atoi(matches[1])
		if err != nil || usage < threshold {
			continue
		}
		alerts = append(alerts, strings.TrimSpace(line))
	}
	return alerts
}

func truncateForReport(value string, maxLen int) string {
	value = strings.TrimSpace(value)
	if len(value) <= maxLen {
		return value
	}
	return value[:maxLen] + "\n... (truncated)"
}

func buildAssistantHistoryMessage(response aiModel.AIAssistantChatResponse) string {
	parts := []string{strings.TrimSpace(response.AssistantMessage)}
	if response.CommandResult != nil && strings.TrimSpace(response.CommandResult.Output) != "" {
		parts = append(parts, fmt.Sprintf("命令 `%s` 输出:\n%s", response.CommandResult.Command, truncateForReport(response.CommandResult.Output, 1600)))
	}
	if response.InspectionResult != nil && strings.TrimSpace(response.InspectionResult.Report) != "" {
		parts = append(parts, response.InspectionResult.Report)
	}
	if response.PendingConfirmation != nil {
		parts = append(parts, fmt.Sprintf("待确认任务 #%d: %s", response.PendingConfirmation.ID, response.PendingConfirmation.Summary))
	}
	return strings.Join(parts, "\n\n")
}

func limitHosts(hosts []cmdbModel.CmdbHost, limit int) []cmdbModel.CmdbHost {
	if len(hosts) <= limit {
		return hosts
	}
	return hosts[:limit]
}

func containsAny(source string, values ...string) bool {
	for _, value := range values {
		if strings.Contains(source, strings.ToLower(value)) {
			return true
		}
	}
	return false
}
