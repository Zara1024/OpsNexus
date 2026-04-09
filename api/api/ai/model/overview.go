package model

type AIOverviewRuntime struct {
	Enabled         bool   `json:"enabled"`
	Provider        string `json:"provider"`
	Model           string `json:"model"`
	ReasoningEffort string `json:"reasoningEffort"`
	Status          string `json:"status"`
	StatusText      string `json:"statusText"`
	LastError       string `json:"lastError,omitempty"`
	CheckedAt       string `json:"checkedAt,omitempty"`
}

type AIOverviewStats struct {
	PromptTemplates      int64 `json:"promptTemplates"`
	KnowledgeItems       int64 `json:"knowledgeItems"`
	DiagnosisSessions    int64 `json:"diagnosisSessions"`
	AssistantSessions    int64 `json:"assistantSessions"`
	InspectionTemplates  int64 `json:"inspectionTemplates"`
	InspectionReports    int64 `json:"inspectionReports"`
	PendingConfirmations int64 `json:"pendingConfirmations"`
}

type AIDiagnosisSceneOverview struct {
	Value              string `json:"value"`
	Label              string `json:"label"`
	Description        string `json:"description"`
	TargetLabel        string `json:"targetLabel"`
	TargetPlaceholder  string `json:"targetPlaceholder"`
	KeywordPlaceholder string `json:"keywordPlaceholder"`
	TemplateName       string `json:"templateName"`
}

type AIOverviewDomain struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Prompt      string `json:"prompt"`
	Route       string `json:"route,omitempty"`
}

type AIOverviewResponse struct {
	Runtime         AIOverviewRuntime          `json:"runtime"`
	Stats           AIOverviewStats            `json:"stats"`
	DiagnosisScenes []AIDiagnosisSceneOverview `json:"diagnosisScenes"`
	Domains         []AIOverviewDomain         `json:"domains"`
	QuickPrompts    []string                   `json:"quickPrompts"`
}
