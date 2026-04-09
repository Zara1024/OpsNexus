package controller

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveHostTemplatePathPrefersExistingCandidate(t *testing.T) {
	tempRoot := t.TempDir()
	existing := filepath.Join(tempRoot, "api", "upload", "xlsl", "host.xlsx")
	if err := os.MkdirAll(filepath.Dir(existing), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(existing, []byte("ok"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	resolved, err := resolveExistingPath([]string{
		filepath.Join(tempRoot, "upload", "xlsl", "host.xlsx"),
		existing,
	})
	if err != nil {
		t.Fatalf("resolveExistingPath returned error: %v", err)
	}

	if resolved != existing {
		t.Fatalf("expected %q, got %q", existing, resolved)
	}
}

func TestResolveExistingPathReturnsErrorWhenMissing(t *testing.T) {
	_, err := resolveExistingPath([]string{
		filepath.Join(t.TempDir(), "missing.xlsx"),
	})
	if err == nil {
		t.Fatalf("expected error when all candidates are missing")
	}
}
