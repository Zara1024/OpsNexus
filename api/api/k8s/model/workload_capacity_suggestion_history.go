package model

import "time"

// WorkloadCapacityAlertCenterQuery defines the structured alert center filters.
type WorkloadCapacityAlertCenterQuery struct {
	Keyword      string `json:"keyword,omitempty"`
	Namespace    string `json:"namespace,omitempty"`
	WorkloadName string `json:"workloadName,omitempty"`
	Source       string `json:"source,omitempty"`
}

// WorkloadCapacitySuggestionHistoryQuery defines the history query params.
type WorkloadCapacitySuggestionHistoryQuery struct {
	PageSize int `json:"pageSize"`
	PageNum  int `json:"pageNum"`
}

// WorkloadCapacitySuggestionHistoryEntity stores one capacity suggestion snapshot.
type WorkloadCapacitySuggestionHistoryEntity struct {
	ID                 uint      `gorm:"column:id;primaryKey" json:"id"`
	ClusterID          uint      `gorm:"column:cluster_id;not null;index:idx_capacity_suggestion_target,priority:1" json:"clusterId"`
	ClusterName        string    `gorm:"column:cluster_name;size:191;not null" json:"clusterName"`
	NamespaceName      string    `gorm:"column:namespace_name;size:191;not null;index:idx_capacity_suggestion_target,priority:2" json:"namespace"`
	WorkloadType       string    `gorm:"column:workload_type;size:32;not null;index:idx_capacity_suggestion_target,priority:3" json:"workloadType"`
	WorkloadName       string    `gorm:"column:workload_name;size:191;not null;index:idx_capacity_suggestion_target,priority:4" json:"workloadName"`
	TemplateName       string    `gorm:"column:template_name;size:64" json:"templateName"`
	RenderedPrompt     string    `gorm:"column:rendered_prompt;type:longtext" json:"renderedPrompt"`
	SystemPrompt       string    `gorm:"column:system_prompt;type:longtext" json:"systemPrompt"`
	Report             string    `gorm:"column:report;type:longtext" json:"report"`
	UsedFallback       bool      `gorm:"column:used_fallback;not null;default:true" json:"usedFallback"`
	AlertKeyword       string    `gorm:"column:alert_keyword;size:191" json:"alertKeyword"`
	AlertCenterPath    string    `gorm:"column:alert_center_path;type:text" json:"alertCenterPath"`
	AlertCenterQuery   string    `gorm:"column:alert_center_query;type:longtext" json:"alertCenterQuery"`
	RiskLevel          string    `gorm:"column:risk_level;size:16" json:"riskLevel"`
	RecommendedActions string    `gorm:"column:recommended_actions;type:longtext" json:"recommendedActions"`
	RecommendedPolicy  string    `gorm:"column:recommended_policy;type:text" json:"recommendedPolicy"`
	WatchMetrics       string    `gorm:"column:watch_metrics;type:text" json:"watchMetrics"`
	FollowUpWindow     string    `gorm:"column:follow_up_window;size:64" json:"followUpWindow"`
	Autoscaling        string    `gorm:"column:autoscaling;type:longtext" json:"autoscaling"`
	AlertSummary       string    `gorm:"column:alert_summary;type:longtext" json:"alertSummary"`
	GeneratedByID      uint      `gorm:"column:generated_by_id;default:0;index" json:"generatedById"`
	GeneratedBy        string    `gorm:"column:generated_by;size:64" json:"generatedBy"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime;index" json:"createdAt"`
}

func (WorkloadCapacitySuggestionHistoryEntity) TableName() string {
	return "k8s_workload_capacity_suggestion_history"
}
