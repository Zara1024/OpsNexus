package service

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"dodevops-api/api/cmdb/model"
	monitormodel "dodevops-api/api/monitor/model"
)

type agentMetricsSnapshot struct {
	Timestamp      int64
	CPU            float64
	Memory         float64
	Disk           float64
	DiskReadKB     float64
	DiskWriteKB    float64
	NetworkReceive float64
	NetworkSend    float64
	Load1min       float64
	Load5min       float64
	Load15min      float64
	TotalProcesses float64
}

var prometheusLabelPattern = regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)="([^"]*)"`)

func buildAgentMetricsSnapshot(body string) (agentMetricsSnapshot, bool) {
	snapshot := agentMetricsSnapshot{Timestamp: time.Now().Unix()}
	scanner := bufio.NewScanner(strings.NewReader(body))
	var cpuValues []float64
	parsedAny := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		metricName, labels := splitPrometheusMetric(parts[0])
		value, err := strconv.ParseFloat(parts[len(parts)-1], 64)
		if err != nil {
			continue
		}

		switch metricName {
		case "system_cpu_usage_percent":
			cpuValues = append(cpuValues, value)
			parsedAny = true
		case "system_memory_usage_percent":
			snapshot.Memory = value
			parsedAny = true
		case "system_disk_usage_percent":
			if labels["mountpoint"] == "/" {
				snapshot.Disk = value
				parsedAny = true
			}
		case "system_disk_read_kb_per_second":
			if labels["device"] == "vda" || labels["device"] == "" {
				snapshot.DiskReadKB = value
				parsedAny = true
			}
		case "system_disk_write_kb_per_second":
			if labels["device"] == "vda" || labels["device"] == "" {
				snapshot.DiskWriteKB = value
				parsedAny = true
			}
		case "system_network_receive_kb_per_second":
			snapshot.NetworkReceive = value
			parsedAny = true
		case "system_network_send_kb_per_second":
			snapshot.NetworkSend = value
			parsedAny = true
		case "system_load_average":
			switch labels["period"] {
			case "1min":
				snapshot.Load1min = value
			case "5min":
				snapshot.Load5min = value
			case "15min":
				snapshot.Load15min = value
			}
			parsedAny = true
		case "system_total_processes":
			snapshot.TotalProcesses = value
			parsedAny = true
		}
	}

	if len(cpuValues) > 0 {
		var total float64
		for _, value := range cpuValues {
			total += value
		}
		snapshot.CPU = total / float64(len(cpuValues))
		if snapshot.CPU <= 1 {
			snapshot.CPU *= 100
		}
		parsedAny = true
	}

	return snapshot, parsedAny
}

func splitPrometheusMetric(raw string) (string, map[string]string) {
	raw = strings.TrimSpace(raw)
	if !strings.Contains(raw, "{") {
		return raw, map[string]string{}
	}
	metricName, labelPart, ok := strings.Cut(raw, "{")
	if !ok {
		return raw, map[string]string{}
	}
	labelPart = strings.TrimSuffix(labelPart, "}")
	labels := map[string]string{}
	for _, match := range prometheusLabelPattern.FindAllStringSubmatch(labelPart, -1) {
		if len(match) == 3 {
			labels[match[1]] = match[2]
		}
	}
	return metricName, labels
}

func (s *MonitorServiceImpl) scrapeAgentMetrics(host model.CmdbHost) (agentMetricsSnapshot, error) {
	if strings.TrimSpace(host.SSHIP) == "" {
		return agentMetricsSnapshot{}, fmt.Errorf("host ssh ip is empty")
	}

	client := &http.Client{Timeout: 4 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://%s:9100/metrics", host.SSHIP))
	if err != nil {
		return agentMetricsSnapshot{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return agentMetricsSnapshot{}, err
	}

	snapshot, ok := buildAgentMetricsSnapshot(string(body))
	if !ok {
		return agentMetricsSnapshot{}, fmt.Errorf("agent metrics payload does not contain target metrics")
	}
	return snapshot, nil
}

func applyAgentSnapshotToSummary(metrics *monitormodel.HostMetrics, snapshot agentMetricsSnapshot) {
	metrics.CPUUsage = snapshot.CPU
	metrics.MemoryUsage = snapshot.Memory
	metrics.DiskUsage = snapshot.Disk
}

func applyAgentSnapshotToHistory(history *monitormodel.AllMetricsHistory, snapshot agentMetricsSnapshot) {
	appendPoint := func(target *[]monitormodel.MetricDataPoint, value float64) {
		*target = append(*target, monitormodel.MetricDataPoint{
			Timestamp: snapshot.Timestamp,
			Value:     value,
		})
	}

	appendPoint(&history.CPU, snapshot.CPU)
	appendPoint(&history.Memory, snapshot.Memory)
	appendPoint(&history.Disk, snapshot.Disk)
	appendPoint(&history.DiskReadKB, snapshot.DiskReadKB)
	appendPoint(&history.DiskWriteKB, snapshot.DiskWriteKB)
	appendPoint(&history.NetworkReceive, snapshot.NetworkReceive)
	appendPoint(&history.NetworkSend, snapshot.NetworkSend)
	appendPoint(&history.Load1min, snapshot.Load1min)
	appendPoint(&history.Load5min, snapshot.Load5min)
	appendPoint(&history.Load15min, snapshot.Load15min)
	appendPoint(&history.TotalProcesses, snapshot.TotalProcesses)
}
