package service

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestResolveServiceTargetPortParsesNumericString(t *testing.T) {
	got := resolveServiceTargetPort("80", 8080)

	if got.Type != intstr.Int {
		t.Fatalf("expected int targetPort type, got %v", got.Type)
	}
	if got.IntValue() != 80 {
		t.Fatalf("expected targetPort 80, got %d", got.IntValue())
	}
}

func TestResolveServiceTargetPortKeepsNamedPort(t *testing.T) {
	got := resolveServiceTargetPort("http", 8080)

	if got.Type != intstr.String {
		t.Fatalf("expected string targetPort type, got %v", got.Type)
	}
	if got.String() != "http" {
		t.Fatalf("expected targetPort http, got %q", got.String())
	}
}
