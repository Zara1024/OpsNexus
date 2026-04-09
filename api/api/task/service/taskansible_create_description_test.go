package service

import (
	"testing"
)

func TestBuildTaskModelPreservesDescription(t *testing.T) {
	req := &CreateTaskRequest{
		Name:        "opsnexus-safe-ansible-ping-20260331",
		Description: "Codex safe ansible ping verification on verify host 718",
		IsRecurring: 0,
		CronExpr:    "",
	}

	task := buildTaskModel(req, map[string][]uint{"web": {718}}, []uint{718}, 1)
	if task.Description != req.Description {
		t.Fatalf("expected task description %q, got %q", req.Description, task.Description)
	}
}
