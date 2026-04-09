package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	aiModel "dodevops-api/api/ai/model"
	cmdbDao "dodevops-api/api/cmdb/dao"
)

type assistantLLMPlan struct {
	Intent           string   `json:"intent"`
	TargetType       string   `json:"target_type"`
	TargetValue      string   `json:"target_value"`
	Command          string   `json:"command"`
	AssistantMessage string   `json:"assistant_message"`
	Suggestions      []string `json:"suggestions"`
}

type assistantLLMRefinement struct {
	AssistantMessage  string   `json:"assistant_message"`
	Suggestions       []string `json:"suggestions"`
	InspectionSummary string   `json:"inspection_summary"`
	InspectionReport  string   `json:"inspection_report"`
}

func (s *AIService) planAssistantActionWithLLM(ctx context.Context, client *aiRuntimeClient, message string) (*assistantLLMPlan, error) {
	if client == nil || !client.IsEnabled() {
		return nil, fmt.Errorf("llm runtime unavailable")
	}

	hostCatalog := buildAssistantHostCatalog(message)
	instructions := strings.Join([]string{
		"You are the OpsNexus AI operations planner.",
		"You must classify the user's request into exactly one intent:",
		"assistant_help, host_list, host_lookup, host_command, inspection_report, alert_analysis, workload_lookup, work_order_lookup, deployment_lookup, template_center, report_archive, risky_action.",
		"Allowed read-only commands are: hostname, date, uptime, df -h, free -m, ss -lntp, uname -a, cat /etc/os-release, ps -ef.",
		"Never output destructive or write operations.",
		"If the user asks for risky actions such as restart, delete, stop, or modify, switch to assistant_help and explain the assistant is currently read-only.",
		"Always respond with JSON only and no markdown fences.",
		"Use this schema: {\"intent\":\"...\",\"target_type\":\"none|host_id|ip|host_name\",\"target_value\":\"...\",\"command\":\"...\",\"assistant_message\":\"...\",\"suggestions\":[\"...\"]}.",
		"assistant_message and suggestions must be in Simplified Chinese.",
	}, "\n")

	input := strings.Join([]string{
		"User message:",
		message,
		"",
		"Known host snapshot:",
		hostCatalog,
	}, "\n")

	raw, err := client.CreateTextResponse(ctx, instructions, input, 400)
	if err != nil {
		return nil, err
	}

	var plan assistantLLMPlan
	if err := decodeAssistantJSON(raw, &plan); err != nil {
		return nil, err
	}

	plan.Intent = normalizeAssistantIntent(plan.Intent)
	plan.TargetType = normalizeAssistantTargetType(plan.TargetType)
	plan.Command = strings.TrimSpace(plan.Command)
	plan.TargetValue = strings.TrimSpace(plan.TargetValue)
	plan.AssistantMessage = strings.TrimSpace(plan.AssistantMessage)
	plan.Suggestions = normalizeSuggestions(plan.Suggestions)
	if plan.Intent == "" {
		return nil, fmt.Errorf("llm did not return a valid intent")
	}
	return &plan, nil
}

func (s *AIService) refineAssistantResponseWithLLM(ctx context.Context, client *aiRuntimeClient, message string, response aiModel.AIAssistantChatResponse) aiModel.AIAssistantChatResponse {
	if client == nil || !client.IsEnabled() {
		return response
	}

	payload, _ := json.MarshalIndent(response, "", "  ")
	instructions := strings.Join([]string{
		"You are the OpsNexus AI operations assistant.",
		"Rewrite the final user-facing response based on the structured tool results.",
		"Never fabricate data that does not exist in the payload.",
		"Keep the tone concise, calm, and operational.",
		"If inspection_result exists, write a high-quality inspection summary and a markdown inspection report in Simplified Chinese.",
		"If only host_matches or command_result exists, summarize the key findings and suggest next safe follow-up actions.",
		"Return JSON only with schema:",
		"{\"assistant_message\":\"...\",\"suggestions\":[\"...\"],\"inspection_summary\":\"...\",\"inspection_report\":\"...\"}.",
	}, "\n")

	input := strings.Join([]string{
		"Original user message:",
		message,
		"",
		"Structured tool result:",
		string(payload),
	}, "\n")

	raw, err := client.CreateTextResponse(ctx, instructions, input, 900)
	if err != nil {
		response.FallbackReason = firstNonEmpty(response.FallbackReason, err.Error())
		return response
	}

	var refinement assistantLLMRefinement
	if err := decodeAssistantJSON(raw, &refinement); err != nil {
		response.FallbackReason = firstNonEmpty(response.FallbackReason, err.Error())
		return response
	}

	if strings.TrimSpace(refinement.AssistantMessage) != "" {
		response.AssistantMessage = strings.TrimSpace(refinement.AssistantMessage)
	}
	if normalized := normalizeSuggestions(refinement.Suggestions); len(normalized) > 0 {
		response.Suggestions = normalized
	}
	if response.InspectionResult != nil {
		if strings.TrimSpace(refinement.InspectionSummary) != "" {
			response.InspectionResult.Summary = strings.TrimSpace(refinement.InspectionSummary)
		}
		if strings.TrimSpace(refinement.InspectionReport) != "" {
			response.InspectionResult.Report = strings.TrimSpace(refinement.InspectionReport)
		}
	}
	return response
}

func buildAssistantHostCatalog(message string) string {
	hostDao := cmdbDao.NewCmdbHostDao()
	candidates, _ := resolveAssistantHosts(message)
	if len(candidates) == 0 {
		candidates, _ = hostDao.GetCmdbHostListWithPage(1, 12)
	}

	if len(candidates) == 0 {
		return "No host records available."
	}

	items := toAssistantHosts(limitHosts(candidates, 12))
	lines := make([]string, 0, len(items)+1)
	lines = append(lines, "Use these hosts only as hints, not as proof of availability.")
	for _, item := range items {
		lines = append(lines, fmt.Sprintf(
			"- id=%d, host=%s, ssh_ip=%s, private_ip=%s, public_ip=%s, group=%s, status=%s",
			item.ID,
			firstNonEmpty(item.HostName, "-"),
			firstNonEmpty(item.SSHIP, "-"),
			firstNonEmpty(item.PrivateIP, "-"),
			firstNonEmpty(item.PublicIP, "-"),
			firstNonEmpty(item.GroupName, "-"),
			firstNonEmpty(item.StatusText, "-"),
		))
	}
	return strings.Join(lines, "\n")
}

func normalizeAssistantIntent(intent string) string {
	intent = strings.TrimSpace(strings.ToLower(intent))
	switch intent {
	case "assistant_help", "host_list", "host_lookup", "host_command", "inspection_report", "alert_analysis", "workload_lookup", "work_order_lookup", "deployment_lookup", "template_center", "report_archive", "risky_action":
		return intent
	default:
		return ""
	}
}

func normalizeAssistantTargetType(targetType string) string {
	targetType = strings.TrimSpace(strings.ToLower(targetType))
	switch targetType {
	case "host_id", "ip", "host_name":
		return targetType
	default:
		return "none"
	}
}

func normalizeSuggestions(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[string]struct{})
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func decodeAssistantJSON(raw string, target interface{}) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fmt.Errorf("empty llm response")
	}
	if err := json.Unmarshal([]byte(raw), target); err == nil {
		return nil
	}

	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start >= 0 && end > start {
		return json.Unmarshal([]byte(raw[start:end+1]), target)
	}
	return fmt.Errorf("llm response is not valid json")
}

func composeAssistantExecutionMessage(message string, plan *assistantLLMPlan) string {
	if plan == nil {
		return message
	}

	composed := strings.TrimSpace(message)
	if plan.TargetValue != "" && !strings.Contains(composed, plan.TargetValue) {
		switch plan.TargetType {
		case "host_id":
			composed = strings.TrimSpace(composed + " host " + plan.TargetValue)
		case "ip":
			composed = strings.TrimSpace(composed + " " + plan.TargetValue)
		case "host_name":
			composed = strings.TrimSpace(composed + " host " + plan.TargetValue)
		}
	}
	if plan.Command != "" && !strings.Contains(composed, plan.Command) {
		composed = strings.TrimSpace(composed + " `" + plan.Command + "`")
	}
	return composed
}
