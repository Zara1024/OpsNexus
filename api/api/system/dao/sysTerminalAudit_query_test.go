package dao

import (
	"strings"
	"testing"

	"dodevops-api/api/system/model"
)

func TestBuildTerminalAuditWhereSupportsHostFilters(t *testing.T) {
	query := model.TerminalAuditQuery{
		HostID:      718,
		HostKeyword: "opsnexus-local-verify",
		RiskLevel:   -1,
	}

	whereSQL, args := buildTerminalAuditWhere(query)

	if !strings.Contains(whereSQL, "ssr.host_id = ?") {
		t.Fatalf("expected host_id filter in where SQL, got %q", whereSQL)
	}
	if !strings.Contains(whereSQL, "COALESCE(ssr.host_name, '') LIKE ?") {
		t.Fatalf("expected host_name filter in where SQL, got %q", whereSQL)
	}
	if !strings.Contains(whereSQL, "COALESCE(ssr.host_ip, '') LIKE ?") {
		t.Fatalf("expected host_ip filter in where SQL, got %q", whereSQL)
	}
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d (%v)", len(args), args)
	}
}
