package model

import "time"

type AIAssistantSessionContextEntity struct {
	ID                    uint      `gorm:"column:id;primaryKey" json:"id"`
	SessionID             string    `gorm:"column:session_id;size:64;not null;uniqueIndex" json:"sessionId"`
	UserID                uint      `gorm:"column:user_id;not null;index" json:"userId"`
	CurrentScope          string    `gorm:"column:current_scope;size:32" json:"currentScope"`
	CurrentHostID         uint      `gorm:"column:current_host_id;default:0" json:"currentHostId"`
	CurrentHostName       string    `gorm:"column:current_host_name;size:191" json:"currentHostName"`
	CurrentClusterID      uint      `gorm:"column:current_cluster_id;default:0" json:"currentClusterId"`
	CurrentClusterName    string    `gorm:"column:current_cluster_name;size:191" json:"currentClusterName"`
	CurrentNamespace      string    `gorm:"column:current_namespace;size:191" json:"currentNamespace"`
	CurrentWorkloadType   string    `gorm:"column:current_workload_type;size:64" json:"currentWorkloadType"`
	CurrentWorkloadName   string    `gorm:"column:current_workload_name;size:191" json:"currentWorkloadName"`
	CurrentWorkOrderType  string    `gorm:"column:current_workorder_type;size:32" json:"currentWorkorderType"`
	CurrentWorkOrderID    uint      `gorm:"column:current_workorder_id;default:0" json:"currentWorkorderId"`
	CurrentDeploymentID   uint      `gorm:"column:current_deployment_id;default:0" json:"currentDeploymentId"`
	LastIntent            string    `gorm:"column:last_intent;size:64" json:"lastIntent"`
	Summary               string    `gorm:"column:summary;type:text" json:"summary"`
	PendingConfirmationID uint      `gorm:"column:pending_confirmation_id;default:0" json:"pendingConfirmationId"`
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt             time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (AIAssistantSessionContextEntity) TableName() string {
	return "ai_assistant_session_context"
}

type AIAssistantConfirmationEntity struct {
	ID            uint      `gorm:"column:id;primaryKey" json:"id"`
	SessionID     string    `gorm:"column:session_id;size:64;not null;index" json:"sessionId"`
	UserID        uint      `gorm:"column:user_id;not null;index" json:"userId"`
	Status        string    `gorm:"column:status;size:32;not null;default:'pending';index" json:"status"`
	Scope         string    `gorm:"column:scope;size:32" json:"scope"`
	ActionType    string    `gorm:"column:action_type;size:64;not null" json:"actionType"`
	TargetID      uint      `gorm:"column:target_id;default:0" json:"targetId"`
	TargetName    string    `gorm:"column:target_name;size:191" json:"targetName"`
	Command       string    `gorm:"column:command;type:text" json:"command"`
	Payload       string    `gorm:"column:payload;type:longtext" json:"payload"`
	Summary       string    `gorm:"column:summary;type:text" json:"summary"`
	ResultSummary string    `gorm:"column:result_summary;type:longtext" json:"resultSummary"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime;index" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	ExpiresAt     time.Time `gorm:"column:expires_at;index" json:"expiresAt"`
}

func (AIAssistantConfirmationEntity) TableName() string {
	return "ai_assistant_confirmation"
}

type AIInspectionTemplateEntity struct {
	ID          uint      `gorm:"column:id;primaryKey" json:"id"`
	Name        string    `gorm:"column:name;size:128;not null;uniqueIndex" json:"name"`
	Category    string    `gorm:"column:category;size:64" json:"category"`
	Scope       string    `gorm:"column:scope;size:32;not null;default:'host'" json:"scope"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	ChecksJSON  string    `gorm:"column:checks_json;type:longtext;not null" json:"checksJson"`
	PromptHint  string    `gorm:"column:prompt_hint;type:text" json:"promptHint"`
	Enabled     int       `gorm:"column:enabled;default:1;index" json:"enabled"`
	IsBuiltin   int       `gorm:"column:is_builtin;default:1;index" json:"isBuiltin"`
	CreatedBy   string    `gorm:"column:created_by;size:64" json:"createdBy"`
	UpdatedBy   string    `gorm:"column:updated_by;size:64" json:"updatedBy"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (AIInspectionTemplateEntity) TableName() string {
	return "ai_inspection_template"
}

type AIInspectionReportEntity struct {
	ID           uint      `gorm:"column:id;primaryKey" json:"id"`
	SessionID    string    `gorm:"column:session_id;size:64;index" json:"sessionId"`
	UserID       uint      `gorm:"column:user_id;not null;index" json:"userId"`
	TemplateID   uint      `gorm:"column:template_id;default:0;index" json:"templateId"`
	TemplateName string    `gorm:"column:template_name;size:128" json:"templateName"`
	Scope        string    `gorm:"column:scope;size:32;not null;default:'host'" json:"scope"`
	TargetID     uint      `gorm:"column:target_id;default:0;index" json:"targetId"`
	TargetName   string    `gorm:"column:target_name;size:191;index" json:"targetName"`
	Summary      string    `gorm:"column:summary;type:text" json:"summary"`
	Report       string    `gorm:"column:report;type:longtext" json:"report"`
	CheckResults string    `gorm:"column:check_results;type:longtext" json:"checkResults"`
	Status       string    `gorm:"column:status;size:32;default:'completed';index" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime;index" json:"createdAt"`
}

func (AIInspectionReportEntity) TableName() string {
	return "ai_inspection_report"
}

type AIAssistantContext struct {
	CurrentScope         string `json:"currentScope"`
	CurrentHostID        uint   `json:"currentHostId"`
	CurrentHostName      string `json:"currentHostName"`
	CurrentClusterID     uint   `json:"currentClusterId"`
	CurrentClusterName   string `json:"currentClusterName"`
	CurrentNamespace     string `json:"currentNamespace"`
	CurrentWorkloadType  string `json:"currentWorkloadType"`
	CurrentWorkloadName  string `json:"currentWorkloadName"`
	CurrentWorkorderType string `json:"currentWorkorderType"`
	CurrentWorkorderID   uint   `json:"currentWorkorderId"`
	CurrentDeploymentID  uint   `json:"currentDeploymentId"`
	LastIntent           string `json:"lastIntent"`
	Summary              string `json:"summary"`
}

type AIAssistantToolStep struct {
	Tool    string `json:"tool"`
	Status  string `json:"status"`
	Summary string `json:"summary"`
	Detail  string `json:"detail,omitempty"`
}

type AIAssistantPendingConfirmation struct {
	ID         uint   `json:"id"`
	Status     string `json:"status"`
	ActionType string `json:"actionType"`
	Scope      string `json:"scope"`
	TargetName string `json:"targetName"`
	Command    string `json:"command,omitempty"`
	Summary    string `json:"summary"`
	ExpiresAt  string `json:"expiresAt"`
}

type AIAssistantInspectionTemplateCheck struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

type AIAssistantInspectionTemplate struct {
	ID          uint                                 `json:"id"`
	Name        string                               `json:"name"`
	Category    string                               `json:"category"`
	Scope       string                               `json:"scope"`
	Description string                               `json:"description"`
	PromptHint  string                               `json:"promptHint"`
	IsBuiltin   bool                                 `json:"isBuiltin"`
	Checks      []AIAssistantInspectionTemplateCheck `json:"checks"`
}

type AIAssistantInspectionReportSummary struct {
	ID           uint   `json:"id"`
	TemplateID   uint   `json:"templateId"`
	TemplateName string `json:"templateName"`
	Scope        string `json:"scope"`
	TargetID     uint   `json:"targetId"`
	TargetName   string `json:"targetName"`
	Summary      string `json:"summary"`
	Status       string `json:"status"`
	CreatedAt    string `json:"createdAt"`
}

type AIAssistantConfirmationDecisionRequest struct {
	Decision string `json:"decision"`
}
