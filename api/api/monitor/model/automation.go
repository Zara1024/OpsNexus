package model

import "time"

const (
	MonitorEventStatusOpen     = "open"
	MonitorEventStatusResolved = "resolved"

	MonitorRuleCategoryHost = "host"
	MonitorRuleCategoryDB   = "database"
	MonitorRuleCategorySSL  = "ssl"
)

// MonitorHostAlertRuleEntity stores host alert threshold rules.
type MonitorHostAlertRuleEntity struct {
	ID             uint      `gorm:"column:id;primaryKey" json:"id"`
	Name           string    `gorm:"column:name;size:191;not null" json:"name"`
	HostID         uint      `gorm:"column:host_id;not null" json:"hostId"`
	MetricKey      string    `gorm:"column:metric_key;size:64;not null" json:"metricKey"`
	Operator       string    `gorm:"column:operator;size:16;not null" json:"operator"`
	ThresholdValue float64   `gorm:"column:threshold_value;type:decimal(20,4);not null" json:"thresholdValue"`
	Severity       string    `gorm:"column:severity;size:16;not null" json:"severity"`
	Status         int       `gorm:"column:status;default:1" json:"status"`
	NotifyRobotIDs string    `gorm:"column:notify_robot_ids;type:text" json:"notifyRobotIds"`
	Remark         string    `gorm:"column:remark;type:text" json:"remark"`
	LastScanAt     time.Time `gorm:"column:last_scan_at" json:"lastScanAt"`
	CreateTime     time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime     time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

func (MonitorHostAlertRuleEntity) TableName() string {
	return "monitor_host_alert_rule"
}

// MonitorDBAlertRuleEntity stores database alert threshold rules.
type MonitorDBAlertRuleEntity struct {
	ID             uint      `gorm:"column:id;primaryKey" json:"id"`
	Name           string    `gorm:"column:name;size:191;not null" json:"name"`
	DatabaseID     uint      `gorm:"column:database_id;not null" json:"databaseId"`
	MetricKey      string    `gorm:"column:metric_key;size:64;not null" json:"metricKey"`
	Operator       string    `gorm:"column:operator;size:16;not null" json:"operator"`
	ThresholdValue float64   `gorm:"column:threshold_value;type:decimal(20,4);not null" json:"thresholdValue"`
	Severity       string    `gorm:"column:severity;size:16;not null" json:"severity"`
	Status         int       `gorm:"column:status;default:1" json:"status"`
	NotifyRobotIDs string    `gorm:"column:notify_robot_ids;type:text" json:"notifyRobotIds"`
	Remark         string    `gorm:"column:remark;type:text" json:"remark"`
	LastScanAt     time.Time `gorm:"column:last_scan_at" json:"lastScanAt"`
	CreateTime     time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime     time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

func (MonitorDBAlertRuleEntity) TableName() string {
	return "monitor_db_alert_rule"
}

// MonitorAlertEventEntity stores trigger/recovery lifecycle for automation alerts.
type MonitorAlertEventEntity struct {
	ID               uint       `gorm:"column:id;primaryKey" json:"id"`
	RuleCategory     string     `gorm:"column:rule_category;size:32;not null" json:"ruleCategory"`
	RuleID           uint       `gorm:"column:rule_id;default:0" json:"ruleId"`
	ResourceType     string     `gorm:"column:resource_type;size:32;not null" json:"resourceType"`
	ResourceID       uint       `gorm:"column:resource_id;not null" json:"resourceId"`
	ResourceName     string     `gorm:"column:resource_name;size:191;not null" json:"resourceName"`
	Fingerprint      string     `gorm:"column:fingerprint;size:191;not null;index:idx_monitor_alert_event_fingerprint" json:"fingerprint"`
	EventKey         string     `gorm:"column:event_key;size:64;not null" json:"eventKey"`
	Title            string     `gorm:"column:title;size:255;not null" json:"title"`
	Summary          string     `gorm:"column:summary;type:text" json:"summary"`
	Detail           string     `gorm:"column:detail;type:text" json:"detail"`
	Severity         string     `gorm:"column:severity;size:16;not null" json:"severity"`
	Status           string     `gorm:"column:status;size:16;not null;index:idx_monitor_alert_event_status" json:"status"`
	Operator         string     `gorm:"column:operator;size:16" json:"operator"`
	ThresholdValue   float64    `gorm:"column:threshold_value;type:decimal(20,4)" json:"thresholdValue"`
	CurrentValue     float64    `gorm:"column:current_value;type:decimal(20,4)" json:"currentValue"`
	RecoveryValue    float64    `gorm:"column:recovery_value;type:decimal(20,4)" json:"recoveryValue"`
	DedupCount       int64      `gorm:"column:dedup_count;default:1" json:"dedupCount"`
	NotifyRobotIDs   string     `gorm:"column:notify_robot_ids;type:text" json:"notifyRobotIds"`
	IncidentID       uint       `gorm:"column:incident_id;default:0" json:"incidentId"`
	WebhookLogID     uint       `gorm:"column:webhook_log_id;default:0" json:"webhookLogId"`
	RecoveryLogID    uint       `gorm:"column:recovery_log_id;default:0" json:"recoveryLogId"`
	Labels           string     `gorm:"column:labels;type:text" json:"labels"`
	FirstTriggeredAt time.Time  `gorm:"column:first_triggered_at;not null" json:"firstTriggeredAt"`
	LastTriggeredAt  time.Time  `gorm:"column:last_triggered_at;not null" json:"lastTriggeredAt"`
	RecoveredAt      *time.Time `gorm:"column:recovered_at" json:"recoveredAt"`
	CreateTime       time.Time  `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime       time.Time  `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

func (MonitorAlertEventEntity) TableName() string {
	return "monitor_alert_event"
}

// MonitorIncidentEntity maps to monitor_incident for internal alert lifecycle updates.
type MonitorIncidentEntity struct {
	ID             uint      `gorm:"column:id;primaryKey" json:"id"`
	AlertTime      time.Time `gorm:"column:alert_time;not null" json:"alertTime"`
	BusinessLine   string    `gorm:"column:business_line;type:longtext" json:"businessLine"`
	Frequency      string    `gorm:"column:frequency;type:longtext" json:"frequency"`
	AlertDesc      string    `gorm:"column:alert_desc;type:text" json:"alertDesc"`
	AlertLevel     string    `gorm:"column:alert_level;size:191" json:"alertLevel"`
	IncidentCause  string    `gorm:"column:incident_cause;type:text" json:"incidentCause"`
	Department     string    `gorm:"column:department;type:longtext" json:"department"`
	Solution       string    `gorm:"column:solution;type:text" json:"solution"`
	DetailURL      string    `gorm:"column:detail_url;type:longtext" json:"detailUrl"`
	Handler        string    `gorm:"column:handler;type:longtext" json:"handler"`
	HandlerID      uint      `gorm:"column:handler_id" json:"handlerId"`
	Status         int       `gorm:"column:status;default:1" json:"status"`
	Remark         string    `gorm:"column:remark;type:text" json:"remark"`
	CreateTime     time.Time `gorm:"column:create_time;not null" json:"createTime"`
	UpdateTime     time.Time `gorm:"column:update_time" json:"updateTime"`
	BusinessLineID uint      `gorm:"column:business_line_id" json:"businessLineId"`
}

func (MonitorIncidentEntity) TableName() string {
	return "monitor_incident"
}

// MonitorDBHealthSnapshotEntity stores the latest database health snapshot.
type MonitorDBHealthSnapshotEntity struct {
	ID           uint      `gorm:"column:id;primaryKey" json:"id"`
	DatabaseID   uint      `gorm:"column:database_id;uniqueIndex:uk_monitor_db_health_snapshot_database" json:"databaseId"`
	DatabaseName string    `gorm:"column:database_name;size:191;not null" json:"databaseName"`
	DatabaseType int       `gorm:"column:database_type;not null" json:"databaseType"`
	Host         string    `gorm:"column:host;size:191" json:"host"`
	Port         int       `gorm:"column:port" json:"port"`
	Available    int       `gorm:"column:available;default:0" json:"available"`
	LatencyMs    int64     `gorm:"column:latency_ms;default:0" json:"latencyMs"`
	ErrorMsg     string    `gorm:"column:error_msg;type:text" json:"errorMsg"`
	LastCheckAt  time.Time `gorm:"column:last_check_at;not null" json:"lastCheckAt"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

func (MonitorDBHealthSnapshotEntity) TableName() string {
	return "monitor_db_health_snapshot"
}

// MonitorDomainEntity maps to the existing monitor_domain table.
type MonitorDomainEntity struct {
	ID           uint       `gorm:"column:id;primaryKey" json:"id"`
	Domain       string     `gorm:"column:domain;size:255;not null" json:"domain"`
	Tags         string     `gorm:"column:tags;size:500" json:"tags"`
	Remark       string     `gorm:"column:remark;type:text" json:"remark"`
	Status       int64      `gorm:"column:status;default:1" json:"status"`
	IsAlive      int64      `gorm:"column:is_alive;default:0" json:"isAlive"`
	StatusCode   int64      `gorm:"column:status_code" json:"statusCode"`
	ResponseTime int64      `gorm:"column:response_time" json:"responseTime"`
	SSLExpireAt  *time.Time `gorm:"column:ssl_expire_at" json:"sslExpireAt"`
	SSLDaysLeft  int64      `gorm:"column:ssl_days_left" json:"sslDaysLeft"`
	SSLIssuer    string     `gorm:"column:ssl_issuer;size:255" json:"sslIssuer"`
	LastCheckAt  *time.Time `gorm:"column:last_check_at" json:"lastCheckAt"`
	ErrorMsg     string     `gorm:"column:error_msg;type:text" json:"errorMsg"`
	CreateTime   time.Time  `gorm:"column:create_time" json:"createTime"`
	UpdateTime   time.Time  `gorm:"column:update_time" json:"updateTime"`
}

func (MonitorDomainEntity) TableName() string {
	return "monitor_domain"
}

// MonitorDomainScheduleEntity maps to the existing domain scan schedule and adds extensibility fields.
type MonitorDomainScheduleEntity struct {
	ID                uint       `gorm:"column:id;primaryKey" json:"id"`
	Enabled           bool       `gorm:"column:enabled" json:"enabled"`
	CronExpr          string     `gorm:"column:cron_expr;size:100" json:"cronExpr"`
	NextRunAt         *time.Time `gorm:"column:next_run_at" json:"nextRunAt"`
	LastRunAt         *time.Time `gorm:"column:last_run_at" json:"lastRunAt"`
	Status            string     `gorm:"column:status;size:50" json:"status"`
	NotifyEnabled     bool       `gorm:"column:notify_enabled" json:"notifyEnabled"`
	NotifyRobotID     uint       `gorm:"column:notify_robot_id" json:"notifyRobotId"`
	ExpireAlertDays   int64      `gorm:"column:expire_alert_days;default:30" json:"expireAlertDays"`
	ScanTimeoutMs     int64      `gorm:"column:scan_timeout_ms;default:8000" json:"scanTimeoutMs"`
	AutoRenewEnabled  bool       `gorm:"column:auto_renew_enabled;default:0" json:"autoRenewEnabled"`
	AutoDeployEnabled bool       `gorm:"column:auto_deploy_enabled;default:0" json:"autoDeployEnabled"`
	DeployHostID      uint       `gorm:"column:deploy_host_id" json:"deployHostId"`
	DeployPath        string     `gorm:"column:deploy_path;size:500" json:"deployPath"`
	ReloadCommand     string     `gorm:"column:reload_command;size:500" json:"reloadCommand"`
	CreateTime        time.Time  `gorm:"column:create_time" json:"createTime"`
	UpdateTime        time.Time  `gorm:"column:update_time" json:"updateTime"`
}

func (MonitorDomainScheduleEntity) TableName() string {
	return "monitor_domain_schedule"
}

// MonitorSSLCertEntity maps to the existing monitor_ssl_cert table.
type MonitorSSLCertEntity struct {
	ID             uint       `gorm:"column:id;primaryKey" json:"id"`
	Domain         string     `gorm:"column:domain;size:191;not null" json:"domain"`
	AliyunConfigID uint       `gorm:"column:aliyun_config_id" json:"aliyunConfigId"`
	OrderID        string     `gorm:"column:order_id;size:191" json:"orderId"`
	CertID         int64      `gorm:"column:cert_id" json:"certId"`
	CertName       string     `gorm:"column:cert_name;type:longtext" json:"certName"`
	ProductCode    string     `gorm:"column:product_code;size:191" json:"productCode"`
	Status         int64      `gorm:"column:status" json:"status"`
	ValidateType   string     `gorm:"column:validate_type;size:191" json:"validateType"`
	ValidateInfo   string     `gorm:"column:validate_info;type:text" json:"validateInfo"`
	Cert           string     `gorm:"column:cert;type:text" json:"-"`
	PrivateKey     string     `gorm:"column:private_key;type:text" json:"-"`
	IssueTime      *time.Time `gorm:"column:issue_time" json:"issueTime"`
	ExpireTime     *time.Time `gorm:"column:expire_time" json:"expireTime"`
	DaysLeft       int64      `gorm:"column:days_left" json:"daysLeft"`
	ErrorMsg       string     `gorm:"column:error_msg;type:text" json:"errorMsg"`
	CreateTime     time.Time  `gorm:"column:create_time" json:"createTime"`
	UpdateTime     time.Time  `gorm:"column:update_time" json:"updateTime"`
	CertSource     string     `gorm:"column:cert_source;size:191" json:"certSource"`
	CAProvider     string     `gorm:"column:ca_provider;size:191" json:"caProvider"`
	IssuerCert     string     `gorm:"column:issuer_cert;type:text" json:"-"`
	Algorithm      string     `gorm:"column:algorithm;size:191" json:"algorithm"`
}

func (MonitorSSLCertEntity) TableName() string {
	return "monitor_ssl_cert"
}

// MonitorSSLCertDeployLogEntity maps to the existing monitor_ssl_cert_deploy_log table.
type MonitorSSLCertDeployLogEntity struct {
	ID          uint      `gorm:"column:id;primaryKey" json:"id"`
	CertID      uint      `gorm:"column:cert_id;not null" json:"certId"`
	Domain      string    `gorm:"column:domain;size:255" json:"domain"`
	HostID      uint      `gorm:"column:host_id;not null" json:"hostId"`
	HostName    string    `gorm:"column:host_name;size:255" json:"hostName"`
	DeployPath  string    `gorm:"column:deploy_path;size:500" json:"deployPath"`
	Status      int       `gorm:"column:status;default:1" json:"status"`
	BackupFiles string    `gorm:"column:backup_files;type:text" json:"backupFiles"`
	DeployFiles string    `gorm:"column:deploy_files;type:text" json:"deployFiles"`
	Logs        string    `gorm:"column:logs;type:text" json:"logs"`
	ErrorMsg    string    `gorm:"column:error_msg;type:text" json:"errorMsg"`
	CreateTime  time.Time `gorm:"column:create_time;not null" json:"createTime"`
	UpdateTime  time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (MonitorSSLCertDeployLogEntity) TableName() string {
	return "monitor_ssl_cert_deploy_log"
}

// MonitorHostAlertRuleUpsertRequest defines host rule write payload.
type MonitorHostAlertRuleUpsertRequest struct {
	Name           string  `json:"name" binding:"required"`
	HostID         uint    `json:"hostId" binding:"required"`
	MetricKey      string  `json:"metricKey" binding:"required"`
	Operator       string  `json:"operator" binding:"required"`
	ThresholdValue float64 `json:"thresholdValue"`
	Severity       string  `json:"severity" binding:"required"`
	Status         int     `json:"status"`
	NotifyRobotIDs []uint  `json:"notifyRobotIds"`
	Remark         string  `json:"remark"`
}

// MonitorDBAlertRuleUpsertRequest defines database rule write payload.
type MonitorDBAlertRuleUpsertRequest struct {
	Name           string  `json:"name" binding:"required"`
	DatabaseID     uint    `json:"databaseId" binding:"required"`
	MetricKey      string  `json:"metricKey" binding:"required"`
	Operator       string  `json:"operator" binding:"required"`
	ThresholdValue float64 `json:"thresholdValue"`
	Severity       string  `json:"severity" binding:"required"`
	Status         int     `json:"status"`
	NotifyRobotIDs []uint  `json:"notifyRobotIds"`
	Remark         string  `json:"remark"`
}

// MonitorAutomationEventQuery defines automation event filters.
type MonitorAutomationEventQuery struct {
	ResourceType string
	Status       string
	Keyword      string
	PageSize     int
	PageNum      int
}

// MonitorDomainUpsertRequest defines domain write payload.
type MonitorDomainUpsertRequest struct {
	Domain string `json:"domain" binding:"required"`
	Tags   string `json:"tags"`
	Remark string `json:"remark"`
	Status int64  `json:"status"`
}

// MonitorDomainScheduleUpsertRequest defines domain schedule write payload.
type MonitorDomainScheduleUpsertRequest struct {
	Enabled           bool   `json:"enabled"`
	CronExpr          string `json:"cronExpr"`
	NotifyEnabled     bool   `json:"notifyEnabled"`
	NotifyRobotID     uint   `json:"notifyRobotId"`
	ExpireAlertDays   int64  `json:"expireAlertDays"`
	ScanTimeoutMs     int64  `json:"scanTimeoutMs"`
	AutoRenewEnabled  bool   `json:"autoRenewEnabled"`
	AutoDeployEnabled bool   `json:"autoDeployEnabled"`
	DeployHostID      uint   `json:"deployHostId"`
	DeployPath        string `json:"deployPath"`
	ReloadCommand     string `json:"reloadCommand"`
}

// MonitorSSLDeployRequest defines one manual SSL deployment task.
type MonitorSSLDeployRequest struct {
	CertID        uint   `json:"certId" binding:"required"`
	HostID        uint   `json:"hostId" binding:"required"`
	DeployPath    string `json:"deployPath" binding:"required"`
	ReloadCommand string `json:"reloadCommand"`
}

// MonitorAlertRuleTemplate is one built-in alert rule template.
type MonitorAlertRuleTemplate struct {
	Name           string  `json:"name"`
	MetricKey      string  `json:"metricKey"`
	Operator       string  `json:"operator"`
	ThresholdValue float64 `json:"thresholdValue"`
	Severity       string  `json:"severity"`
	Description    string  `json:"description"`
}

// MonitorAutomationOverview aggregates high-level metrics for the automation page.
type MonitorAutomationOverview struct {
	HostRuleCount        int64                              `json:"hostRuleCount"`
	DatabaseRuleCount    int64                              `json:"databaseRuleCount"`
	OpenEventCount       int64                              `json:"openEventCount"`
	ResolvedEventCount   int64                              `json:"resolvedEventCount"`
	DatabaseHealthyCount int64                              `json:"databaseHealthyCount"`
	DatabaseTotalCount   int64                              `json:"databaseTotalCount"`
	DomainTotalCount     int64                              `json:"domainTotalCount"`
	DomainAliveCount     int64                              `json:"domainAliveCount"`
	ExpiringDomainCount  int64                              `json:"expiringDomainCount"`
	DeployLogTotalCount  int64                              `json:"deployLogTotalCount"`
	TotalRobotCount      int64                              `json:"totalRobotCount"`
	EnabledRobotCount    int64                              `json:"enabledRobotCount"`
	TotalNotifyLogCount  int64                              `json:"totalNotifyLogCount"`
	FailedNotifyLogCount int64                              `json:"failedNotifyLogCount"`
	RecentEvents         []MonitorAutomationWorkbenchEvent  `json:"recentEvents,omitempty"`
	RecentActions        []MonitorAutomationWorkbenchAction `json:"recentActions,omitempty"`
	RiskTips             []string                           `json:"riskTips,omitempty"`
	RecommendedActions   []MonitorAutomationWorkbenchAction `json:"recommendedActions,omitempty"`
	AlertCenterPath      string                             `json:"alertCenterPath,omitempty"`
	AlertHistoryPath     string                             `json:"alertHistoryPath,omitempty"`
	AlertNotifyPath      string                             `json:"alertNotifyPath,omitempty"`
}

type MonitorAutomationWorkbenchEvent struct {
	Title        string `json:"title"`
	Severity     string `json:"severity,omitempty"`
	Status       string `json:"status,omitempty"`
	Summary      string `json:"summary,omitempty"`
	ResourceType string `json:"resourceType,omitempty"`
	ResourceName string `json:"resourceName,omitempty"`
	OccurredAt   string `json:"occurredAt,omitempty"`
}

type MonitorAutomationWorkbenchAction struct {
	Title   string `json:"title"`
	Status  string `json:"status,omitempty"`
	Summary string `json:"summary,omitempty"`
	Path    string `json:"path,omitempty"`
}
