package service

import (
	"testing"
	"time"

	agentModel "dodevops-api/api/monitor/model"
	monitormodel "dodevops-api/api/monitor/model"
	"dodevops-api/common/util"
)

func TestResolveMonitorCollectionStatusPrometheusUnavailable(t *testing.T) {
	status := resolveMonitorCollectionStatus(true, false, nil, false)

	if status.DataAvailable {
		t.Fatalf("expected unavailable status when prometheus is down")
	}
	if status.CollectionStatus != monitormodel.MonitorCollectionUnavailable {
		t.Fatalf("expected unavailable collection status, got %q", status.CollectionStatus)
	}
	if status.UnavailableReason == "" {
		t.Fatalf("expected unavailable reason when prometheus is down")
	}
}

func TestResolveMonitorCollectionStatusAgentUnavailable(t *testing.T) {
	status := resolveMonitorCollectionStatus(true, false, nil, false)

	if status.DataAvailable {
		t.Fatalf("expected unavailable status when agent is missing")
	}
	if status.CollectionStatus != monitormodel.MonitorCollectionUnavailable {
		t.Fatalf("expected unavailable collection status, got %q", status.CollectionStatus)
	}
	if status.UnavailableReason == "" {
		t.Fatalf("expected unavailable reason when agent is missing")
	}
}

func TestResolveMonitorCollectionStatusReadyWithSamples(t *testing.T) {
	agent := &agentModel.Agent{
		Status:        agentModel.AgentStatusRunning,
		LastHeartbeat: util.HTime{Time: time.Now()},
	}
	status := resolveMonitorCollectionStatus(true, true, agent, true)

	if !status.DataAvailable {
		t.Fatalf("expected data available when monitoring chain is healthy")
	}
	if status.CollectionStatus != monitormodel.MonitorCollectionReady {
		t.Fatalf("expected ready collection status, got %q", status.CollectionStatus)
	}
	if status.UnavailableReason != "" {
		t.Fatalf("expected no unavailable reason when data is available, got %q", status.UnavailableReason)
	}
}

func TestResolveMonitorCollectionStatusAllowsDirectSamplesWithoutPrometheus(t *testing.T) {
	agent := &agentModel.Agent{
		Status:        agentModel.AgentStatusRunning,
		LastHeartbeat: util.HTime{Time: time.Now()},
	}
	status := resolveMonitorCollectionStatus(true, false, agent, true)

	if !status.DataAvailable {
		t.Fatalf("expected direct agent samples to mark data available")
	}
	if status.CollectionStatus != monitormodel.MonitorCollectionReady {
		t.Fatalf("expected ready collection status with direct samples, got %q", status.CollectionStatus)
	}
}

func TestAllMetricsHistoryHasSamples(t *testing.T) {
	history := &monitormodel.AllMetricsHistory{
		CPU: []monitormodel.MetricDataPoint{{Timestamp: 1774320000, Value: 12.5}},
	}

	if !allMetricsHistoryHasSamples(history) {
		t.Fatalf("expected history with cpu points to count as having samples")
	}
}

func mustMonitorHeartbeat(value string) time.Time {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	return parsed
}
