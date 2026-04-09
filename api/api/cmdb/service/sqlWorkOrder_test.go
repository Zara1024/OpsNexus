package service

import "testing"

func TestAnalyzeSQLWorkOrder(t *testing.T) {
	level, summary, requireApproval, tables, rollbackSQL, rollbackHint := analyzeSQLWorkOrder("UPDATE knowledge_base SET title = 'x' WHERE id = 1;")
	if level != sqlRiskMedium {
		t.Fatalf("expected medium risk, got %d", level)
	}
	if !requireApproval {
		t.Fatalf("expected update to require approval")
	}
	if tables != "knowledge_base" {
		t.Fatalf("expected affected table knowledge_base, got %q", tables)
	}
	if rollbackSQL == "" || rollbackHint == "" || summary == "" {
		t.Fatalf("expected rollback guidance and summary to be populated")
	}

	level, summary, _, _, _, _ = analyzeSQLWorkOrder("DELETE FROM knowledge_base;")
	if level != sqlRiskHigh {
		t.Fatalf("expected high risk for delete without where, got %d", level)
	}
	if summary == "" {
		t.Fatalf("expected delete summary to be populated")
	}

	level, _, _, tables, _, _ = analyzeSQLWorkOrder("CREATE TABLE IF NOT EXISTS opsnexus_p2_sql_ticket_demo (id BIGINT);")
	if level != sqlRiskMedium {
		t.Fatalf("expected medium risk for create table, got %d", level)
	}
	if tables != "opsnexus_p2_sql_ticket_demo" {
		t.Fatalf("expected extracted create table name, got %q", tables)
	}
}

func TestExtractSQLMutationTarget(t *testing.T) {
	tableName, whereClause, ok := extractSQLMutationTarget("DELETE FROM knowledge_base WHERE id = 1;")
	if !ok {
		t.Fatalf("expected delete statement to be recognized")
	}
	if tableName != "knowledge_base" {
		t.Fatalf("unexpected table name: %s", tableName)
	}
	if whereClause != "id = 1" {
		t.Fatalf("unexpected where clause: %s", whereClause)
	}
}

func TestDetectRedisWorkOrderOperationType(t *testing.T) {
	if got := detectRedisWorkOrderOperationType("SET demo value"); got != "SET" {
		t.Fatalf("expected SET operation type, got %q", got)
	}
	if got := detectRedisWorkOrderOperationType("  flushall "); got != "FLUSHALL" {
		t.Fatalf("expected FLUSHALL operation type, got %q", got)
	}
}

func TestAnalyzeRedisWorkOrder(t *testing.T) {
	level, summary, requireApproval, keys, rollbackSQL, rollbackHint := analyzeRedisWorkOrder("SET app:feature true")
	if level != sqlRiskMedium {
		t.Fatalf("expected medium risk for redis write command, got %d", level)
	}
	if !requireApproval {
		t.Fatalf("expected redis write command to require approval")
	}
	if keys == "" || summary == "" || rollbackHint == "" {
		t.Fatalf("expected redis work order metadata to be populated")
	}
	if rollbackSQL == "" {
		t.Fatalf("expected rollback guidance placeholder for redis command")
	}
}
