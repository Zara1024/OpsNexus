package service

import (
	"testing"
	"time"

	"dodevops-api/api/task/model"
)

func TestBuildTaskHistoryRecordsReturnsCurrentTaskSnapshot(t *testing.T) {
	start := time.Date(2026, 3, 24, 12, 39, 45, 0, time.Local)
	finish := time.Date(2026, 3, 24, 13, 12, 23, 0, time.Local)
	task := &model.TaskAnsible{
		ID:            53,
		Name:          "opsnexus-test-ansible-1774327185",
		Status:        3,
		CreatedAt:     start,
		UpdatedAt:     finish,
		TotalDuration: 1,
		Works: []model.TaskAnsibleWork{
			{ID: 99, EntryFileName: "site.yml", Status: 3, Duration: 1},
		},
	}

	records := buildTaskHistoryRecords(task, 0)
	if len(records) != 1 {
		t.Fatalf("expected 1 synthesized history record, got %d", len(records))
	}

	record := records[0]
	if record.ID != task.ID {
		t.Fatalf("expected history record ID %d, got %d", task.ID, record.ID)
	}
	if record.Status != task.Status {
		t.Fatalf("expected history record status %d, got %d", task.Status, record.Status)
	}
	if record.OperatorName != "System" {
		t.Fatalf("expected fallback operator name System, got %q", record.OperatorName)
	}
	if record.FinishedAt == nil || !record.FinishedAt.Equal(finish) {
		t.Fatalf("expected finishedAt %v, got %#v", finish, record.FinishedAt)
	}
	if len(record.Works) != 1 || record.Works[0].ID != 99 {
		t.Fatalf("expected synthesized works to preserve current task works, got %#v", record.Works)
	}
}

func TestBuildTaskHistoryRecordsFiltersByStatus(t *testing.T) {
	task := &model.TaskAnsible{
		ID:        52,
		Name:      "opsnexus-test-ansible",
		Status:    4,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	records := buildTaskHistoryRecords(task, 3)
	if len(records) != 0 {
		t.Fatalf("expected no history records when status filter mismatches, got %d", len(records))
	}
}

func TestBuildTaskHistoryDetailReturnsWorkHistories(t *testing.T) {
	task := &model.TaskAnsible{
		ID:   49,
		Name: "test2",
		Works: []model.TaskAnsibleWork{
			{ID: 49, EntryFileName: "01-linux-os.yaml.yml", Status: 3, Duration: 17},
		},
	}

	detail, ok := buildTaskHistoryDetail(task, 49)
	if !ok {
		t.Fatal("expected synthesized history detail for current task snapshot")
	}
	if detail.HistoryID != 49 {
		t.Fatalf("expected historyID 49, got %d", detail.HistoryID)
	}
	if len(detail.WorkHistories) != 1 {
		t.Fatalf("expected 1 work history item, got %d", len(detail.WorkHistories))
	}
	if detail.WorkHistories[0].WorkID != 49 || detail.WorkHistories[0].HostName != "01-linux-os.yaml.yml" {
		t.Fatalf("unexpected synthesized work history payload: %#v", detail.WorkHistories[0])
	}
}
