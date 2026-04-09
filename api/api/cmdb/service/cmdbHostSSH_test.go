package service

import "testing"

func TestSanitizeCommandOutputRemovesAnsiAndNormalizesNewlines(t *testing.T) {
	input := []byte("\u001b[32mline1\u001b[0m\r\nline2\r\n")

	output := sanitizeCommandOutput(input)

	if output != "line1\nline2\n" {
		t.Fatalf("expected sanitized output with normalized newlines, got %q", output)
	}
}

func TestSanitizeCommandOutputKeepsPlainTerminalLayout(t *testing.T) {
	input := []byte(".\n..\n.profile\n")

	output := sanitizeCommandOutput(input)

	if output != ".\n..\n.profile\n" {
		t.Fatalf("expected plain output to remain unchanged, got %q", output)
	}
}

func TestSanitizeCommandOutputRemovesTemplateArtifacts(t *testing.T) {
	input := []byte(".\n .Destination}}{{end}}\n.profile\n")

	output := sanitizeCommandOutput(input)

	if output != ".\n.profile\n" {
		t.Fatalf("expected template artifact line removed, got %q", output)
	}
}

func TestNormalizeNonInteractiveShellCommandExpandsCommonLsAliases(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    string
	}{
		{
			name:    "ll with args expands to ls -al",
			command: "ll /opt",
			want:    "ls -al /opt",
		},
		{
			name:    "ll without args expands to ls -al",
			command: "ll",
			want:    "ls -al",
		},
		{
			name:    "non alias command remains unchanged",
			command: "pwd",
			want:    "pwd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeNonInteractiveShellCommand(tt.command)
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
