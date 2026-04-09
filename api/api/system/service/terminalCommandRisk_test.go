package service

import "testing"

func TestClassifyTerminalCommandRisk(t *testing.T) {
	tests := []struct {
		name                 string
		command              string
		wantLevel            int64
		wantSensitive        bool
		wantRequiresConfirm  bool
		wantLevelLabel       string
	}{
		{
			name:                "empty command is low risk",
			command:             "   ",
			wantLevel:           0,
			wantSensitive:       false,
			wantRequiresConfirm: false,
			wantLevelLabel:      "low",
		},
		{
			name:                "normal ls command is low risk",
			command:             "ls -al /tmp",
			wantLevel:           0,
			wantSensitive:       false,
			wantRequiresConfirm: false,
			wantLevelLabel:      "low",
		},
		{
			name:                "restart command is medium risk",
			command:             "systemctl restart nginx",
			wantLevel:           1,
			wantSensitive:       true,
			wantRequiresConfirm: true,
			wantLevelLabel:      "medium",
		},
		{
			name:                "rm rf is high risk",
			command:             "rm -rf /tmp/demo",
			wantLevel:           2,
			wantSensitive:       true,
			wantRequiresConfirm: true,
			wantLevelLabel:      "high",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assessment := ClassifyTerminalCommandRisk(tt.command)
			if assessment.RiskLevel != tt.wantLevel {
				t.Fatalf("expected risk level %d, got %d", tt.wantLevel, assessment.RiskLevel)
			}
			if assessment.IsSensitive != tt.wantSensitive {
				t.Fatalf("expected sensitive=%v, got %v", tt.wantSensitive, assessment.IsSensitive)
			}
			if assessment.RequiresConfirmation != tt.wantRequiresConfirm {
				t.Fatalf("expected requiresConfirmation=%v, got %v", tt.wantRequiresConfirm, assessment.RequiresConfirmation)
			}
			if assessment.RiskLevelLabel != tt.wantLevelLabel {
				t.Fatalf("expected risk label %q, got %q", tt.wantLevelLabel, assessment.RiskLevelLabel)
			}
		})
	}
}
