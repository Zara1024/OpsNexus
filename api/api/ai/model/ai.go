package model

import "dodevops-api/common/util"

type PromptTemplate struct {
	ID           uint       `gorm:"column:id;primaryKey" json:"id"`
	Name         string     `gorm:"column:name" json:"name"`
	Category     string     `gorm:"column:category" json:"category"`
	Scene        string     `gorm:"column:scene" json:"scene"`
	Template     string     `gorm:"column:template" json:"template"`
	Variables    string     `gorm:"column:variables" json:"variables"`
	SystemPrompt string     `gorm:"column:system_prompt" json:"systemPrompt"`
	Temperature  float64    `gorm:"column:temperature" json:"temperature"`
	MaxTokens    int64      `gorm:"column:max_tokens" json:"maxTokens"`
	ModelID      uint       `gorm:"column:model_id" json:"modelId"`
	Enabled      int        `gorm:"column:enabled" json:"enabled"`
	CreateTime   util.HTime `gorm:"column:create_time" json:"createTime"`
	UpdateTime   util.HTime `gorm:"column:update_time" json:"updateTime"`
}

func (PromptTemplate) TableName() string {
	return "prompt_template"
}

type AIChatHistory struct {
	ID         uint       `gorm:"column:id;primaryKey" json:"id"`
	SessionID  string     `gorm:"column:session_id" json:"sessionId"`
	UserID     uint       `gorm:"column:user_id" json:"userId"`
	Role       string     `gorm:"column:role" json:"role"`
	Message    string     `gorm:"column:message" json:"message"`
	Intent     string     `gorm:"column:intent" json:"intent"`
	IntentConf float64    `gorm:"column:intent_conf" json:"intentConf"`
	Entities   string     `gorm:"column:entities" json:"entities"`
	TaskID     uint       `gorm:"column:task_id" json:"taskId"`
	TaskType   string     `gorm:"column:task_type" json:"taskType"`
	Status     int        `gorm:"column:status" json:"status"`
	ErrorMsg   string     `gorm:"column:error_msg" json:"errorMsg"`
	CreateTime util.HTime `gorm:"column:create_time" json:"createTime"`
}

func (AIChatHistory) TableName() string {
	return "ai_agent_chat_history"
}

type AIPromptRenderRequest struct {
	SessionID    string                 `json:"sessionId"`
	TemplateName string                 `json:"templateName"`
	Intent       string                 `json:"intent"`
	InputMessage string                 `json:"inputMessage"`
	Variables    map[string]interface{} `json:"variables"`
	KnowledgeIDs []uint                 `json:"knowledgeIds"`
}

type AIPromptRenderResponse struct {
	SessionID      string                   `json:"sessionId"`
	Template       PromptTemplate           `json:"template"`
	RenderedPrompt string                   `json:"renderedPrompt"`
	SystemPrompt   string                   `json:"systemPrompt"`
	KnowledgeItems []map[string]interface{} `json:"knowledgeItems"`
}

type AIDiagnosisRequest struct {
	Scene        string `json:"scene"`
	TargetID     string `json:"targetId"`
	Keyword      string `json:"keyword"`
	TemplateName string `json:"templateName"`
	KnowledgeIDs []uint `json:"knowledgeIds"`
}

type AIDiagnosisResponse struct {
	SessionID      string                   `json:"sessionId"`
	Scene          string                   `json:"scene"`
	TargetID       string                   `json:"targetId"`
	Title          string                   `json:"title"`
	Report         string                   `json:"report"`
	RenderedPrompt string                   `json:"renderedPrompt"`
	SystemPrompt   string                   `json:"systemPrompt"`
	UsedFallback   bool                     `json:"usedFallback"`
	KnowledgeItems []map[string]interface{} `json:"knowledgeItems"`
}
