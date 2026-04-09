package service

import (
	"testing"

	"dodevops-api/api/k8s/model"
)

func TestBuildIngressBackendTestResultSuccess(t *testing.T) {
	details := model.IngressBackendTestDetails{
		ServiceName:    "opsnexus-test-svc",
		ServicePort:    80,
		TestPath:       "/",
		TestHost:       "opsnexus-test.local",
		Method:         "GET",
		EndpointsReady: 1,
		EndpointsTotal: 1,
	}

	got := buildIngressBackendTestResult(200, 123, details, nil)

	if !got.Success {
		t.Fatalf("expected success, got %#v", got)
	}
	if got.Status != "成功" {
		t.Fatalf("expected status 成功, got %q", got.Status)
	}
	if got.StatusCode != 200 {
		t.Fatalf("expected statusCode 200, got %d", got.StatusCode)
	}
	if got.ResponseTime != 123 {
		t.Fatalf("expected responseTime 123, got %d", got.ResponseTime)
	}
	if got.Message == "" {
		t.Fatalf("expected non-empty message")
	}
}

func TestBuildIngressBackendTestResultNoReadyEndpoints(t *testing.T) {
	details := model.IngressBackendTestDetails{
		ServiceName:    "opsnexus-test-svc",
		ServicePort:    80,
		TestPath:       "/",
		Method:         "GET",
		EndpointsReady: 0,
		EndpointsTotal: 1,
	}

	got := buildIngressBackendTestResult(0, 0, details, nil)

	if got.Success {
		t.Fatalf("expected failure, got %#v", got)
	}
	if got.Status != "失败" {
		t.Fatalf("expected status 失败, got %q", got.Status)
	}
	if got.Suggestion == "" {
		t.Fatalf("expected suggestion for no ready endpoints")
	}
}
