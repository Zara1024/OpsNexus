package service

import (
	"net/url"
	"testing"

	"dodevops-api/api/app/model"
)

func TestNormalizeAppEnvDeployConfig(t *testing.T) {
	svc := &ApplicationService{}
	cfg := svc.normalizeAppEnvDeployConfig(model.AppEnvDeployConfig{
		ClusterID:                34,
		Namespace:                " opsnexus-apply-e2e ",
		WorkloadType:             "Deployment",
		WorkloadName:             " demo-nginx ",
		ReleaseGovernanceEnabled: true,
	})

	if cfg.Namespace != "opsnexus-apply-e2e" {
		t.Fatalf("expected trimmed namespace, got %q", cfg.Namespace)
	}
	if cfg.WorkloadType != "deployment" {
		t.Fatalf("expected normalized workload type deployment, got %q", cfg.WorkloadType)
	}
	if cfg.WorkloadName != "demo-nginx" {
		t.Fatalf("expected trimmed workload name, got %q", cfg.WorkloadName)
	}
	if cfg.PreCheckMode != "observe" {
		t.Fatalf("expected default pre_check_mode observe, got %q", cfg.PreCheckMode)
	}
}

func TestBuildReleaseGovernanceSummaryBlocksStrictWhenMappingMissing(t *testing.T) {
	svc := &ApplicationService{}
	summary := svc.buildReleaseGovernanceSummary(model.AppEnvDeployConfig{
		ReleaseGovernanceEnabled: true,
		PreCheckMode:             "strict",
	})

	if summary == nil {
		t.Fatal("expected release governance summary, got nil")
	}
	if !summary.Blocking {
		t.Fatalf("expected blocking summary, got %#v", summary)
	}
	if summary.BlockingReason == "" {
		t.Fatalf("expected blocking reason, got %#v", summary)
	}
}

func TestBuildAppEnvPaths(t *testing.T) {
	svc := &ApplicationService{}
	cfg := model.AppEnvDeployConfig{
		ClusterID:    34,
		Namespace:    "opsnexus-apply-e2e",
		WorkloadType: "deployment",
		WorkloadName: "demo-nginx",
	}

	k8sPath := svc.buildAppEnvK8sWorkloadPath(cfg, "governance")
	parsedK8s, err := url.Parse(k8sPath)
	if err != nil {
		t.Fatalf("parse k8s path failed: %v", err)
	}
	if parsedK8s.Path != "/k8s/workload" {
		t.Fatalf("unexpected k8s path: %s", parsedK8s.Path)
	}
	if parsedK8s.Query().Get("clusterId") != "34" || parsedK8s.Query().Get("namespace") != "opsnexus-apply-e2e" || parsedK8s.Query().Get("name") != "demo-nginx" || parsedK8s.Query().Get("action") != "governance" {
		t.Fatalf("unexpected k8s path query: %s", parsedK8s.RawQuery)
	}

	aiPath := svc.buildAppEnvAIDiagnosisPath(12, "demo-nginx", "opsnexus-apply-e2e")
	parsedAI, err := url.Parse(aiPath)
	if err != nil {
		t.Fatalf("parse ai path failed: %v", err)
	}
	if parsedAI.Path != "/ai/diagnosis" {
		t.Fatalf("unexpected ai path: %s", parsedAI.Path)
	}
	if parsedAI.Query().Get("scene") != "workload_capacity" || parsedAI.Query().Get("targetId") != "12" || parsedAI.Query().Get("workloadName") != "demo-nginx" {
		t.Fatalf("unexpected ai path query: %s", parsedAI.RawQuery)
	}
}
