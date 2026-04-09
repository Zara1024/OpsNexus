package model

import "time"

// MonitorAlertSummary is the top-level alert center summary payload.
type MonitorAlertSummary struct {
	TotalIncidents       int64  `json:"totalIncidents"`
	OpenIncidents        int64  `json:"openIncidents"`
	ProcessingIncidents  int64  `json:"processingIncidents"`
	ResolvedIncidents    int64  `json:"resolvedIncidents"`
	TotalWebhookLogs     int64  `json:"totalWebhookLogs"`
	CriticalWebhookLogs  int64  `json:"criticalWebhookLogs"`
	TotalNotifyLogs      int64  `json:"totalNotifyLogs"`
	SuccessfulNotifyLogs int64  `json:"successfulNotifyLogs"`
	FailedNotifyLogs     int64  `json:"failedNotifyLogs"`
	TotalNotifyRobots    int64  `json:"totalNotifyRobots"`
	EnabledNotifyRobots  int64  `json:"enabledNotifyRobots"`
	TotalAlertSources    int64  `json:"totalAlertSources"`
	EnabledAlertSources  int64  `json:"enabledAlertSources"`
	LatestAlertTime      string `json:"latestAlertTime"`
	LatestNotifyTime     string `json:"latestNotifyTime"`
}

// MonitorIncidentQuery defines the supported filters for incident records.
type MonitorIncidentQuery struct {
	Keyword      string
	Status       int
	Level        string
	Namespace    string
	WorkloadName string
	PageSize     int
	PageNum      int
}

// MonitorWebhookLogQuery defines the supported filters for webhook alert logs.
type MonitorWebhookLogQuery struct {
	Keyword  string
	Source   string
	Level    string
	Status   string
	PageSize int
	PageNum  int
}

// MonitorNotifyLogQuery defines the supported filters for notify logs.
type MonitorNotifyLogQuery struct {
	Keyword   string
	Status    string
	RobotType string
	PageSize  int
	PageNum   int
}

// MonitorIncident is one incident row from monitor_incident.
type MonitorIncident struct {
	ID            uint   `json:"id"`
	AlertTime     string `json:"alertTime"`
	BusinessLine  string `json:"businessLine"`
	Frequency     string `json:"frequency"`
	AlertDesc     string `json:"alertDesc"`
	AlertLevel    string `json:"alertLevel"`
	IncidentCause string `json:"incidentCause"`
	Department    string `json:"department"`
	Solution      string `json:"solution"`
	DetailURL     string `json:"detailUrl"`
	Handler       string `json:"handler"`
	HandlerID     uint   `json:"handlerId"`
	Status        int    `json:"status"`
	Remark        string `json:"remark"`
	CreateTime    string `json:"createTime"`
	UpdateTime    string `json:"updateTime"`
}

// MonitorWebhookLog is one webhook alert log row.
type MonitorWebhookLog struct {
	ID             uint   `json:"id"`
	Source         string `json:"source"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	Level          string `json:"level"`
	Tags           string `json:"tags"`
	Extra          string `json:"extra"`
	NotifyRobotIDs string `json:"notifyRobotIds"`
	Status         string `json:"status"`
	ErrorMsg       string `json:"errorMsg"`
	NotifyCount    int64  `json:"notifyCount"`
	SuccessCount   int64  `json:"successCount"`
	FailedCount    int64  `json:"failedCount"`
	CreatedAt      string `json:"createdAt"`
}

// MonitorWebhookNotifyLog is one push delivery row joined with webhook context.
type MonitorWebhookNotifyLog struct {
	ID           uint   `json:"id"`
	WebhookLogID uint   `json:"webhookLogId"`
	RobotID      uint   `json:"robotId"`
	RobotName    string `json:"robotName"`
	RobotType    string `json:"robotType"`
	Status       string `json:"status"`
	ErrorMsg     string `json:"errorMsg"`
	CreatedAt    string `json:"createdAt"`
	AlertTitle   string `json:"alertTitle"`
	AlertSource  string `json:"alertSource"`
	AlertLevel   string `json:"alertLevel"`
}

// MonitorNotifyRobot is one notify robot row.
type MonitorNotifyRobot struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Webhook   string `json:"webhook"`
	Secret    string `json:"secret"`
	Server    string `json:"server"`
	Port      int64  `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Nickname  string `json:"nickname"`
	Headers   string `json:"headers"`
	Method    string `json:"method"`
	Template  string `json:"template"`
	Status    int    `json:"status"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// MonitorAlertSource is one alert source row.
type MonitorAlertSource struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Type       int    `json:"type"`
	AppKey     string `json:"appKey"`
	APIBaseURL string `json:"apiBaseUrl"`
	Status     int    `json:"status"`
	Remark     string `json:"remark"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
	KeyID      uint   `json:"keyId"`
	HostID     uint   `json:"hostId"`
}

// MonitorAlertManagerQuery defines the source selector for AlertManager actions.
type MonitorAlertManagerQuery struct {
	SourceID uint
}

// MonitorAlertManagerStatus is the normalized AlertManager status payload.
type MonitorAlertManagerStatus struct {
	SourceID      uint   `json:"sourceId"`
	SourceName    string `json:"sourceName"`
	Endpoint      string `json:"endpoint"`
	Available     bool   `json:"available"`
	ClusterName   string `json:"clusterName"`
	ClusterStatus string `json:"clusterStatus"`
	Version       string `json:"version"`
	Revision      string `json:"revision"`
	BuildDate     string `json:"buildDate"`
	Uptime        string `json:"uptime"`
}

// MonitorAlertManagerMatcher is one AlertManager matcher item.
type MonitorAlertManagerMatcher struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	IsRegex bool   `json:"isRegex"`
	IsEqual bool   `json:"isEqual"`
}

// MonitorAlertManagerSilence is one AlertManager silence row.
type MonitorAlertManagerSilence struct {
	SourceID   uint                         `json:"sourceId"`
	SourceName string                       `json:"sourceName"`
	ID         string                       `json:"id"`
	Matchers   []MonitorAlertManagerMatcher `json:"matchers"`
	StartsAt   string                       `json:"startsAt"`
	EndsAt     string                       `json:"endsAt"`
	CreatedBy  string                       `json:"createdBy"`
	Comment    string                       `json:"comment"`
	Status     string                       `json:"status"`
	UpdatedAt  string                       `json:"updatedAt"`
}

// MonitorAlertManagerSilenceCreateRequest defines the payload for creating one silence.
type MonitorAlertManagerSilenceCreateRequest struct {
	SourceID  uint                         `json:"sourceId"`
	Matchers  []MonitorAlertManagerMatcher `json:"matchers" binding:"required"`
	StartsAt  string                       `json:"startsAt"`
	EndsAt    string                       `json:"endsAt"`
	CreatedBy string                       `json:"createdBy"`
	Comment   string                       `json:"comment" binding:"required"`
}

// MonitorAlertManagerReceiverIntegration is one receiver integration item.
type MonitorAlertManagerReceiverIntegration struct {
	Name         string `json:"name"`
	SendResolved bool   `json:"sendResolved"`
}

// MonitorAlertManagerReceiver is one AlertManager receiver row.
type MonitorAlertManagerReceiver struct {
	SourceID     uint                                     `json:"sourceId"`
	SourceName   string                                   `json:"sourceName"`
	Name         string                                   `json:"name"`
	Active       bool                                     `json:"active"`
	Integrations []MonitorAlertManagerReceiverIntegration `json:"integrations"`
}

// MonitorNotifyRobotUpsertRequest defines the write payload for notify robots.
type MonitorNotifyRobotUpsertRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Webhook  string `json:"webhook"`
	Secret   string `json:"secret"`
	Status   int    `json:"status"`
	Remark   string `json:"remark"`
	Server   string `json:"server"`
	Port     int64  `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Headers  string `json:"headers"`
	Method   string `json:"method"`
	Template string `json:"template"`
}

// MonitorNotifyRobotTestRequest defines the payload for one manual test delivery.
type MonitorNotifyRobotTestRequest struct {
	Source  string `json:"source"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Level   string `json:"level"`
}

// MonitorAlertSourceUpsertRequest defines the write payload for alert sources.
type MonitorAlertSourceUpsertRequest struct {
	Name       string `json:"name" binding:"required"`
	Type       int    `json:"type" binding:"required"`
	AppKey     string `json:"appKey"`
	APIBaseURL string `json:"apiBaseUrl"`
	Status     int    `json:"status"`
	Remark     string `json:"remark"`
	KeyID      uint   `json:"keyId"`
	HostID     uint   `json:"hostId"`
}

// MonitorStatusUpdateRequest toggles enabled/disabled state for alert resources.
type MonitorStatusUpdateRequest struct {
	Status int `json:"status" binding:"oneof=0 1"`
}

// MonitorWebhookReceiveRequest defines a generic webhook ingestion payload.
type MonitorWebhookReceiveRequest struct {
	Source         string      `json:"source"`
	Title          string      `json:"title"`
	Content        string      `json:"content"`
	Level          string      `json:"level"`
	Tags           interface{} `json:"tags"`
	Extra          interface{} `json:"extra"`
	NotifyRobotIDs []uint      `json:"notifyRobotIds"`
}

// MonitorNotifyRobotEntity maps to monitor_notify_robot for write operations.
type MonitorNotifyRobotEntity struct {
	ID        uint      `gorm:"column:id;primaryKey" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Type      string    `gorm:"column:type" json:"type"`
	Webhook   string    `gorm:"column:webhook" json:"webhook"`
	Secret    string    `gorm:"column:secret" json:"secret"`
	Status    int       `gorm:"column:status" json:"status"`
	Remark    string    `gorm:"column:remark" json:"remark"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
	Server    string    `gorm:"column:server" json:"server"`
	Port      int64     `gorm:"column:port" json:"port"`
	Username  string    `gorm:"column:username" json:"username"`
	Password  string    `gorm:"column:password" json:"password"`
	Nickname  string    `gorm:"column:nickname" json:"nickname"`
	Headers   string    `gorm:"column:headers" json:"headers"`
	Method    string    `gorm:"column:method" json:"method"`
	Template  string    `gorm:"column:template" json:"template"`
}

func (MonitorNotifyRobotEntity) TableName() string {
	return "monitor_notify_robot"
}

// MonitorAlertSourceEntity maps to monitor_alert_source for write operations.
type MonitorAlertSourceEntity struct {
	ID         uint   `gorm:"column:id;primaryKey" json:"id"`
	Name       string `gorm:"column:name" json:"name"`
	Type       int    `gorm:"column:type" json:"type"`
	AppKey     string `gorm:"column:app_key" json:"appKey"`
	APIBaseURL string `gorm:"column:api_base_url" json:"apiBaseUrl"`
	Status     int    `gorm:"column:status" json:"status"`
	Remark     string `gorm:"column:remark" json:"remark"`
	CreateTime int64  `gorm:"column:create_time" json:"createTime"`
	UpdateTime int64  `gorm:"column:update_time" json:"updateTime"`
	KeyID      uint   `gorm:"column:key_id" json:"keyId"`
	HostID     uint   `gorm:"column:host_id" json:"hostId"`
}

func (MonitorAlertSourceEntity) TableName() string {
	return "monitor_alert_source"
}

// MonitorWebhookLogEntity maps to monitor_webhook_log for write operations.
type MonitorWebhookLogEntity struct {
	ID             uint      `gorm:"column:id;primaryKey" json:"id"`
	Source         string    `gorm:"column:source" json:"source"`
	Title          string    `gorm:"column:title" json:"title"`
	Content        string    `gorm:"column:content" json:"content"`
	Level          string    `gorm:"column:level" json:"level"`
	Tags           string    `gorm:"column:tags" json:"tags"`
	Extra          string    `gorm:"column:extra" json:"extra"`
	NotifyRobotIDs string    `gorm:"column:notify_robot_ids" json:"notifyRobotIds"`
	Status         string    `gorm:"column:status" json:"status"`
	ErrorMsg       string    `gorm:"column:error_msg" json:"errorMsg"`
	NotifyCount    int64     `gorm:"column:notify_count" json:"notifyCount"`
	SuccessCount   int64     `gorm:"column:success_count" json:"successCount"`
	FailedCount    int64     `gorm:"column:failed_count" json:"failedCount"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (MonitorWebhookLogEntity) TableName() string {
	return "monitor_webhook_log"
}

// MonitorWebhookNotifyLogEntity maps to monitor_webhook_notify_log for write operations.
type MonitorWebhookNotifyLogEntity struct {
	ID           uint      `gorm:"column:id;primaryKey" json:"id"`
	WebhookLogID uint      `gorm:"column:webhook_log_id" json:"webhookLogId"`
	RobotID      uint      `gorm:"column:robot_id" json:"robotId"`
	RobotName    string    `gorm:"column:robot_name" json:"robotName"`
	RobotType    string    `gorm:"column:robot_type" json:"robotType"`
	Status       string    `gorm:"column:status" json:"status"`
	ErrorMsg     string    `gorm:"column:error_msg" json:"errorMsg"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (MonitorWebhookNotifyLogEntity) TableName() string {
	return "monitor_webhook_notify_log"
}
