package model

// Overview 仪表盘概览数据。
type Overview struct {
	Asset           AssetOverview    `json:"asset"`
	Alert           AlertOverview    `json:"alert"`
	Resource        ResourceOverview `json:"resource"`
	RecentOps       []RecentOp       `json:"recent_ops"`
	RecentAlertTop5 []SeverityStat   `json:"recent_alert_top5"`
}

// AssetOverview 资产概览。
type AssetOverview struct {
	HostTotal       int64 `json:"host_total"`
	HostOnline      int64 `json:"host_online"`
	HostOffline     int64 `json:"host_offline"`
	HostAlert       int64 `json:"host_alert"`
	K8sClusterTotal int64 `json:"k8s_cluster_total"`
	DBInstanceTotal int64 `json:"db_instance_total"`
}

// AlertOverview 告警概览。
type AlertOverview struct {
	ActiveTotal   int64          `json:"active_total"`
	TodayHandled  int64          `json:"today_handled"`
	SeverityStats []SeverityStat `json:"severity_stats"`
}

// SeverityStat 严重等级统计。
type SeverityStat struct {
	Severity string `json:"severity"`
	Count    int64  `json:"count"`
}

// ResourceOverview 资源水位。
type ResourceOverview struct {
	CPUAvg    float64 `json:"cpu_avg"`
	MemoryAvg float64 `json:"memory_avg"`
	DiskAvg   float64 `json:"disk_avg"`
}

// Trends 仪表盘趋势数据。
type Trends struct {
	Days          int             `json:"days"`
	AlertTrend    []DailyTrend    `json:"alert_trend"`
	ResourceTrend []ResourceTrend `json:"resource_trend"`
}

// DailyTrend 每日趋势。
type DailyTrend struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// ResourceTrend 资源趋势点。
type ResourceTrend struct {
	Date   string  `json:"date"`
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
	Disk   float64 `json:"disk"`
}

// RecentOp 最近操作日志。
type RecentOp struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Module    string `json:"module"`
	Action    string `json:"action"`
	Path      string `json:"path"`
	CreatedAt string `json:"created_at"`
}
