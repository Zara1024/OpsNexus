package service

import (
	"strings"
	"testing"
	"time"

	appModel "dodevops-api/api/app/model"
	monitorModel "dodevops-api/api/monitor/model"
)

func TestDiagnosisSceneCatalogIncludesNewOperationalScenes(t *testing.T) {
	scenes := buildDiagnosisSceneCatalog()
	if len(scenes) == 0 {
		t.Fatalf("expected diagnosis scenes, got none")
	}

	foundAlert := false
	foundDeployment := false
	for _, item := range scenes {
		if item.Value == "alert_analysis" {
			foundAlert = true
		}
		if item.Value == "deployment_review" {
			foundDeployment = true
		}
	}

	if !foundAlert {
		t.Fatalf("expected alert_analysis scene to exist")
	}
	if !foundDeployment {
		t.Fatalf("expected deployment_review scene to exist")
	}
}

func TestBuildAIOverviewResponseRuntimeStatus(t *testing.T) {
	response := buildAIOverviewResponse(aiOverviewSnapshot{
		RuntimeEnabled:    true,
		RuntimeProvider:   "openai",
		RuntimeModel:      "gpt-5.4",
		ReasoningEffort:   "medium",
		RuntimeReachable:  true,
		RuntimeCheckedAt:  time.Date(2026, 3, 24, 10, 0, 0, 0, time.UTC),
		PromptTemplates:   3,
		KnowledgeItems:    8,
		DiagnosisSessions: 2,
		AssistantSessions: 4,
	})

	if response.Runtime.Status != "ready" {
		t.Fatalf("expected runtime status ready, got %q", response.Runtime.Status)
	}
	if response.Runtime.LastError != "" {
		t.Fatalf("expected no runtime error when probe succeeded, got %q", response.Runtime.LastError)
	}
	if response.Runtime.Provider != "openai" {
		t.Fatalf("expected provider openai, got %q", response.Runtime.Provider)
	}
	if response.Stats.KnowledgeItems != 8 {
		t.Fatalf("expected knowledge count 8, got %d", response.Stats.KnowledgeItems)
	}
	if len(response.DiagnosisScenes) == 0 {
		t.Fatalf("expected diagnosis scenes in overview")
	}
	if len(response.Domains) == 0 {
		t.Fatalf("expected domain cards in overview")
	}
}

func TestBuildAIOverviewResponseUsesDegradedStatusWhenProbeFails(t *testing.T) {
	response := buildAIOverviewResponse(aiOverviewSnapshot{
		RuntimeEnabled:   true,
		RuntimeProvider:  "openai",
		RuntimeModel:     "gpt-5.4",
		ReasoningEffort:  "xhigh",
		RuntimeReachable: false,
		RuntimeCheckedAt: time.Date(2026, 3, 24, 10, 5, 0, 0, time.UTC),
		RuntimeLastError: "invalid character '<' looking for beginning of value",
	})

	if response.Runtime.Status != "degraded" {
		t.Fatalf("expected runtime status degraded, got %q", response.Runtime.Status)
	}
	if response.Runtime.LastError == "" {
		t.Fatalf("expected runtime last error to be populated")
	}
	if response.Runtime.CheckedAt == "" {
		t.Fatalf("expected runtime checkedAt to be populated")
	}
}

func TestBuildAlertAnalysisDiagnosisContainsSummaryAndIncidents(t *testing.T) {
	report := buildAlertAnalysisDiagnosis(
		monitorModel.MonitorAlertSummary{
			OpenIncidents:       3,
			CriticalWebhookLogs: 1,
			LatestAlertTime:     "2026-03-22 20:00:00",
		},
		[]monitorModel.MonitorIncident{
			{
				AlertLevel:   "critical",
				BusinessLine: "k8s",
				AlertDesc:    "kubelet memory high",
			},
		},
		nil,
	)

	if !strings.Contains(report, "Open incidents: 3") {
		t.Fatalf("expected report to include open incident count, got %q", report)
	}
	if !strings.Contains(report, "kubelet memory high") {
		t.Fatalf("expected report to include incident detail, got %q", report)
	}
}

func TestBuildDeploymentReviewDiagnosisContainsTaskResult(t *testing.T) {
	report := buildDeploymentReviewDiagnosis(
		appModel.QuickDeployment{
			Title:       "prod release",
			Description: "release for api",
			Status:      4,
			CreatedAt:   time.Date(2026, 3, 22, 20, 0, 0, 0, time.UTC),
		},
		[]appModel.QuickDeploymentTask{
			{
				AppName:      "opsnexus-api",
				Environment:  "prod",
				Status:       4,
				ErrorMessage: "jenkins timeout",
			},
		},
		nil,
	)

	if !strings.Contains(report, "prod release") {
		t.Fatalf("expected deployment title in report, got %q", report)
	}
	if !strings.Contains(report, "jenkins timeout") {
		t.Fatalf("expected task error in report, got %q", report)
	}
}
