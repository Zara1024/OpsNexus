package service

import "testing"

func TestNormalizeKubectlCommandAddsPrefixAndNamespace(t *testing.T) {
	got := normalizeKubectlCommand("get pods", "opsnexus")
	want := "kubectl -n opsnexus get pods"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestNormalizeKubectlCommandKeepsExistingNamespace(t *testing.T) {
	got := normalizeKubectlCommand("kubectl --namespace kube-system get pods", "opsnexus")
	want := "kubectl --namespace kube-system get pods"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestNormalizeKubectlCommandDefaultsToKubectl(t *testing.T) {
	got := normalizeKubectlCommand("   ", "")
	if got != "kubectl" {
		t.Fatalf("expected default kubectl command, got %q", got)
	}
}

func TestHasKubectlNamespaceFlag(t *testing.T) {
	if !hasKubectlNamespaceFlag("kubectl -n default get pods") {
		t.Fatal("expected short namespace flag to be detected")
	}
	if !hasKubectlNamespaceFlag("kubectl --namespace=kube-system get pods") {
		t.Fatal("expected long namespace flag to be detected")
	}
	if hasKubectlNamespaceFlag("kubectl get pods") {
		t.Fatal("did not expect namespace flag to be detected")
	}
}
