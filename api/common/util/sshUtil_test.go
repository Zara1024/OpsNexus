package util

import (
	"strings"
	"testing"
)

func TestGetScriptContentUsesLocaleIndependentMemoryProbe(t *testing.T) {
	script := getScriptContent()

	if !strings.Contains(script, "/proc/meminfo") {
		t.Fatalf("expected getScriptContent to read memory from /proc/meminfo, got script:\n%s", script)
	}

	if strings.Contains(script, "awk '/^Mem:/'") {
		t.Fatalf("expected getScriptContent to avoid locale-dependent free output parsing, got script:\n%s", script)
	}
}

func TestNormalizeMemoryCapacityFromMeminfoKiB(t *testing.T) {
	testCases := map[string]string{
		"32827160": "32",
		"16777216": "16",
		"unknown":  "unknown",
		"":         "",
	}

	for raw, expected := range testCases {
		if actual := normalizeMemoryCapacityFromMeminfoKiB(raw); actual != expected {
			t.Fatalf("expected %q to normalize to %q, got %q", raw, expected, actual)
		}
	}
}
