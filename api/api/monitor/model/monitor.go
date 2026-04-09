package model

type HostMetrics struct {
	CPUUsage          float64 `json:"cpuUsage"`
	MemoryUsage       float64 `json:"memoryUsage"`
	DiskUsage         float64 `json:"diskUsage"`
	OnlineStatus      int     `json:"onlineStatus"`
	DataAvailable     bool    `json:"dataAvailable"`
	CollectionStatus  string  `json:"collectionStatus"`
	UnavailableReason string  `json:"unavailableReason"`
}

const (
	MonitorCollectionReady       = "ready"
	MonitorCollectionOffline     = "offline"
	MonitorCollectionUnavailable = "unavailable"
)

type PrometheusQueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values [][]interface{}   `json:"values"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

type HostMetricHistory struct {
	HostID   uint              `json:"hostId"`
	Metric   string            `json:"metric"`
	TimeData []MetricDataPoint `json:"timeData"`
}

type MetricDataPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

type AllMetricsHistory struct {
	HostID            uint              `json:"hostId"`
	DataAvailable     bool              `json:"dataAvailable"`
	CollectionStatus  string            `json:"collectionStatus"`
	UnavailableReason string            `json:"unavailableReason"`
	CPU               []MetricDataPoint `json:"cpu"`
	Memory            []MetricDataPoint `json:"memory"`
	Disk              []MetricDataPoint `json:"disk"`
	DiskReadKB        []MetricDataPoint `json:"diskReadKB"`
	DiskWriteKB       []MetricDataPoint `json:"diskWriteKB"`
	NetworkReceive    []MetricDataPoint `json:"networkReceive"`
	NetworkSend       []MetricDataPoint `json:"networkSend"`
	Load1min          []MetricDataPoint `json:"load1min"`
	Load5min          []MetricDataPoint `json:"load5min"`
	Load15min         []MetricDataPoint `json:"load15min"`
	TotalProcesses    []MetricDataPoint `json:"totalProcesses"`
}

type ProcessMetrics struct {
	PID        uint              `json:"pid"`
	Name       string            `json:"name"`
	CPUPercent []MetricDataPoint `json:"cpuPercent"`
	MemPercent []MetricDataPoint `json:"memPercent"`
	Host       string            `json:"host"`
}

type ProcessInfo struct {
	PID        uint    `json:"pid"`
	Name       string  `json:"name"`
	CPUPercent float64 `json:"cpuPercent"`
	MemPercent float64 `json:"memPercent"`
	Host       string  `json:"host"`
}

type TopProcessesResult struct {
	HostID     uint          `json:"hostId"`
	HostName   string        `json:"hostName"`
	TopCPU     []ProcessInfo `json:"topCPU"`
	TopMemory  []ProcessInfo `json:"topMemory"`
	UpdateTime int64         `json:"updateTime"`
	DataAvailable     bool   `json:"dataAvailable"`
	CollectionStatus  string `json:"collectionStatus"`
	UnavailableReason string `json:"unavailableReason"`
}

type PortInfo struct {
	Port     string  `json:"port"`
	PID      string  `json:"pid"`
	Service  string  `json:"service"`
	Status   int     `json:"status"`
	CPUUsage float64 `json:"cpuUsage"`
	MemUsage float64 `json:"memUsage"`
}

type HostPortsResult struct {
	HostID     uint       `json:"hostId"`
	HostName   string     `json:"hostName"`
	Ports      []PortInfo `json:"ports"`
	Total      int        `json:"total"`
	UpdateTime int64      `json:"updateTime"`
	DataAvailable     bool   `json:"dataAvailable"`
	CollectionStatus  string `json:"collectionStatus"`
	UnavailableReason string `json:"unavailableReason"`
}
