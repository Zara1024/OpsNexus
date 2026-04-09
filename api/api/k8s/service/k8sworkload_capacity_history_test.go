package service

import (
	"net/url"
	"testing"
	"time"

	"dodevops-api/api/k8s/model"
)

func TestBuildWorkloadCapacityAlertCenterPath(t *testing.T) {
	svc := &K8sWorkloadServiceImpl{}

	query := svc.buildWorkloadCapacityAlertCenterQuery("opsnexus-apply-e2e", "demo-nginx")
	path := svc.buildWorkloadCapacityAlertCenterPath(query)

	parsed, err := url.Parse(path)
	if err != nil {
		t.Fatalf("parse path failed: %v", err)
	}
	if parsed.Path != "/monitor/alert-center" {
		t.Fatalf("unexpected alert center path: %s", parsed.Path)
	}
	if got := parsed.Query().Get("keyword"); got != "demo-nginx" {
		t.Fatalf("expected keyword demo-nginx, got %q", got)
	}
	if got := parsed.Query().Get("namespace"); got != "opsnexus-apply-e2e" {
		t.Fatalf("expected namespace opsnexus-apply-e2e, got %q", got)
	}
	if got := parsed.Query().Get("workloadName"); got != "demo-nginx" {
		t.Fatalf("expected workloadName demo-nginx, got %q", got)
	}
	if got := parsed.Query().Get("source"); got != "capacity-suggestion" {
		t.Fatalf("expected source capacity-suggestion, got %q", got)
	}
}

func TestConvertWorkloadCapacitySuggestionHistory(t *testing.T) {
	svc := &K8sWorkloadServiceImpl{}
	entity := model.WorkloadCapacitySuggestionHistoryEntity{
		ID:                 12,
		ClusterID:          34,
		ClusterName:        "p0-acceptance-k3s",
		NamespaceName:      "opsnexus-apply-e2e",
		WorkloadType:       "deployment",
		WorkloadName:       "demo-nginx",
		TemplateName:       "prediction_suggestion",
		RenderedPrompt:     "rendered",
		SystemPrompt:       "system",
		Report:             "report",
		UsedFallback:       true,
		AlertKeyword:       "demo-nginx",
		AlertCenterPath:    "/monitor/alert-center?keyword=demo-nginx&namespace=opsnexus-apply-e2e&workloadName=demo-nginx&source=capacity-suggestion",
		AlertCenterQuery:   `{"keyword":"demo-nginx","namespace":"opsnexus-apply-e2e","workloadName":"demo-nginx","source":"capacity-suggestion"}`,
		RiskLevel:          "medium",
		RecommendedActions: `[{"priority":"P1","action":"评估并启用 HPA","reason":"当前工作负载仍依赖手工扩缩容","expectedEffect":"降低人工干预成本"}]`,
		RecommendedPolicy:  `{"type":"hpa","minReplicas":"1","maxReplicas":"3","metric":"cpu","targetUtilization":"70"}`,
		WatchMetrics:       `["CPU requests / 实际使用","关联告警数量与恢复时间"]`,
		Autoscaling:        `{"enabled":true,"name":"demo-nginx","minReplicas":1,"maxReplicas":3}`,
		AlertSummary:       `{"openEventCount":1,"resolvedEventCount":0,"incidentCount":2,"criticalCount":1}`,
		GeneratedByID:      89,
		GeneratedBy:        "admin",
		CreatedAt:          time.Date(2026, 3, 20, 21, 0, 0, 0, time.Local),
	}

	response := svc.convertWorkloadCapacitySuggestionHistory(entity)
	if response.HistoryID != entity.ID {
		t.Fatalf("expected historyId %d, got %d", entity.ID, response.HistoryID)
	}
	if response.GeneratedBy != "admin" {
		t.Fatalf("expected generatedBy admin, got %q", response.GeneratedBy)
	}
	if response.AlertCenterQuery.WorkloadName != "demo-nginx" {
		t.Fatalf("expected alertCenterQuery workloadName demo-nginx, got %q", response.AlertCenterQuery.WorkloadName)
	}
	if response.Autoscaling == nil || response.Autoscaling.Name != "demo-nginx" {
		t.Fatalf("expected autoscaling name demo-nginx, got %#v", response.Autoscaling)
	}
	if len(response.RecommendedActions) != 1 {
		t.Fatalf("expected 1 recommended action, got %d", len(response.RecommendedActions))
	}
	if response.AlertSummary.OpenEventCount != 1 || response.AlertSummary.IncidentCount != 2 {
		t.Fatalf("unexpected alert summary: %#v", response.AlertSummary)
	}
	if len(response.WatchMetrics) != 2 {
		t.Fatalf("expected 2 watch metrics, got %d", len(response.WatchMetrics))
	}
}
