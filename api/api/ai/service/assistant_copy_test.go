package service

import (
	"os"
	"strings"
	"testing"
)

func TestBuildAIQuickPromptsUseHostIPExamples(t *testing.T) {
	prompts := buildAIQuickPrompts()
	joined := strings.Join(prompts, "\n")

	expectedPhrases := []string{
		"查询主机 10.0.0.200",
		"在主机 10.0.0.200 执行 `free -m`",
	}
	for _, phrase := range expectedPhrases {
		if !strings.Contains(joined, phrase) {
			t.Fatalf("expected quick prompts to include %q, got %q", phrase, joined)
		}
	}

	forbiddenPhrases := []string{
		"主机 12",
		"主机 ID",
	}
	for _, phrase := range forbiddenPhrases {
		if strings.Contains(joined, phrase) {
			t.Fatalf("expected quick prompts to omit %q, got %q", phrase, joined)
		}
	}
}

func TestAssistantUserFacingCopyStopsAdvertisingHostIDs(t *testing.T) {
	source, err := os.ReadFile("assistant.go")
	if err != nil {
		t.Fatalf("failed to read assistant.go: %v", err)
	}

	text := string(source)
	forbiddenPhrases := []string{
		"主机 ID",
		"主机 12",
		"查看主机 %d 的磁盘占用",
		"查看主机 %d 的内存情况",
		"查看主机 %d 的监听端口",
		"为主机 %d 生成巡检报告",
	}
	for _, phrase := range forbiddenPhrases {
		if strings.Contains(text, phrase) {
			t.Fatalf("expected assistant copy to omit %q", phrase)
		}
	}
}
