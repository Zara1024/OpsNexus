package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	agentModel "dodevops-api/api/monitor/model"
	monitormodel "dodevops-api/api/monitor/model"
)

const monitorAgentHeartbeatTTL = 5 * time.Minute

func resolveMonitorCollectionStatus(hostOnline bool, prometheusReachable bool, agent *agentModel.Agent, hasSamples bool) monitormodel.HostMetrics {
	status := monitormodel.HostMetrics{}
	if !hostOnline {
		status.CollectionStatus = monitormodel.MonitorCollectionOffline
		status.UnavailableReason = "主机离线，无法采集监控数据"
		return status
	}
	if hasSamples {
		status.DataAvailable = true
		status.CollectionStatus = monitormodel.MonitorCollectionReady
		return status
	}
	if !prometheusReachable {
		status.CollectionStatus = monitormodel.MonitorCollectionUnavailable
		status.UnavailableReason = "Prometheus 未接入或不可达"
		return status
	}
	if !isMonitorAgentHealthy(agent, time.Now()) {
		status.CollectionStatus = monitormodel.MonitorCollectionUnavailable
		status.UnavailableReason = "Agent 未部署、未启动或心跳已过期"
		return status
	}
	status.CollectionStatus = monitormodel.MonitorCollectionUnavailable
	status.UnavailableReason = "Prometheus 暂未抓取到该主机的监控样本"
	return status
}

func applyMonitorStatusToSummary(metrics *monitormodel.HostMetrics, hostOnline bool, prometheusReachable bool, agent *agentModel.Agent, hasSamples bool) {
	status := resolveMonitorCollectionStatus(hostOnline, prometheusReachable, agent, hasSamples)
	metrics.DataAvailable = status.DataAvailable
	metrics.CollectionStatus = status.CollectionStatus
	metrics.UnavailableReason = status.UnavailableReason
}

func applyMonitorStatusToHistory(history *monitormodel.AllMetricsHistory, hostOnline bool, prometheusReachable bool, agent *agentModel.Agent) {
	status := resolveMonitorCollectionStatus(hostOnline, prometheusReachable, agent, allMetricsHistoryHasSamples(history))
	history.DataAvailable = status.DataAvailable
	history.CollectionStatus = status.CollectionStatus
	history.UnavailableReason = status.UnavailableReason
}

func allMetricsHistoryHasSamples(history *monitormodel.AllMetricsHistory) bool {
	if history == nil {
		return false
	}
	return len(history.CPU) > 0 ||
		len(history.Memory) > 0 ||
		len(history.Disk) > 0 ||
		len(history.DiskReadKB) > 0 ||
		len(history.DiskWriteKB) > 0 ||
		len(history.NetworkReceive) > 0 ||
		len(history.NetworkSend) > 0 ||
		len(history.Load1min) > 0 ||
		len(history.Load5min) > 0 ||
		len(history.Load15min) > 0 ||
		len(history.TotalProcesses) > 0
}

func isMonitorAgentHealthy(agent *agentModel.Agent, now time.Time) bool {
	if agent == nil || agent.Status != agentModel.AgentStatusRunning {
		return false
	}
	if agent.LastHeartbeat.Time.IsZero() {
		return false
	}
	return now.Sub(agent.LastHeartbeat.Time) <= monitorAgentHeartbeatTTL
}

func (s *MonitorServiceImpl) getMonitorAgentByHostID(hostID uint) (*agentModel.Agent, error) {
	if s.agentDao == nil {
		return nil, fmt.Errorf("agent dao not initialized")
	}
	agent, err := s.agentDao.GetByHostID(hostID)
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func (s *MonitorServiceImpl) isPrometheusReachable() bool {
	baseURL := strings.TrimRight(strings.TrimSpace(s.prometheusURL), "/")
	if baseURL == "" {
		return false
	}

	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest(http.MethodGet, baseURL+"/-/healthy", nil)
	if err != nil {
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
