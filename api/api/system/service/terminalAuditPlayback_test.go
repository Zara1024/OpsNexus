package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"dodevops-api/api/system/model"
	"dodevops-api/common/util"
)

func TestParseTerminalAuditPlayback(t *testing.T) {
	content := strings.Join([]string{
		"[2026-03-20 11:10:41] SESSION START demo",
		"[2026-03-20 11:10:41] IN  kubectl get pods\\n",
		"[2026-03-20 11:10:42] OUT pod-a\\npod-b\\n",
		"[2026-03-20 11:10:43] RESIZE cols=120 rows=40",
		"[2026-03-20 11:10:44] SESSION END status=2 err=",
	}, "\n")
	filePath := filepath.Join(t.TempDir(), "recording.log")
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write test recording: %v", err)
	}

	session := model.TerminalAuditSession{
		SessionID:   "demo",
		DataSource:  "recording",
		StorageType: terminalAuditStorageLocal,
		FilePath:    filePath,
		FileSize:    int64(len(content)),
		StartTime:   util.HTime{Time: time.Date(2026, 3, 20, 11, 10, 41, 0, time.Local)},
	}

	resp, err := parseTerminalAuditPlayback(session, model.TerminalAuditPlaybackQuery{
		PageNum:  1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("parse playback: %v", err)
	}
	if !resp.Health.CanPlayback {
		t.Fatalf("expected playback to be available, got %+v", resp.Health)
	}
	if resp.Stats.TotalEvents != 5 {
		t.Fatalf("expected 5 events, got %d", resp.Stats.TotalEvents)
	}
	if resp.Stats.InputEvents != 1 || resp.Stats.OutputEvents != 1 || resp.Stats.ResizeEvents != 1 {
		t.Fatalf("unexpected stats: %+v", resp.Stats)
	}
	if got := resp.Events[1].Content; got != "kubectl get pods\n" {
		t.Fatalf("expected decoded input content, got %q", got)
	}
	if got := resp.Events[2].Content; got != "pod-a\npod-b\n" {
		t.Fatalf("expected decoded output content, got %q", got)
	}
}

func TestBuildTerminalAuditRecordingHealth(t *testing.T) {
	commandOnly := buildTerminalAuditRecordingHealth(model.TerminalAuditSession{
		DataSource: "command",
	})
	if commandOnly.State != terminalAuditPlaybackStateCommandOnly {
		t.Fatalf("expected command_only state, got %s", commandOnly.State)
	}

	filePath := filepath.Join(t.TempDir(), "empty.log")
	if err := os.WriteFile(filePath, []byte{}, 0o644); err != nil {
		t.Fatalf("write empty file: %v", err)
	}
	empty := buildTerminalAuditRecordingHealth(model.TerminalAuditSession{
		DataSource:  "recording",
		StorageType: terminalAuditStorageLocal,
		FilePath:    filePath,
	})
	if empty.State != terminalAuditPlaybackStateEmpty {
		t.Fatalf("expected empty state, got %s", empty.State)
	}
	if empty.CanPlayback {
		t.Fatalf("expected empty recording to be not playable")
	}
}
