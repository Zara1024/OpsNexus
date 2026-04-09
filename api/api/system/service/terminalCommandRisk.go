package service

import "strings"

const (
	terminalCommandRiskLow    int64 = 0
	terminalCommandRiskMedium int64 = 1
	terminalCommandRiskHigh   int64 = 2
)

type TerminalCommandRiskAssessment struct {
	Command              string `json:"command"`
	IsSensitive          bool   `json:"isSensitive"`
	RiskLevel            int64  `json:"riskLevel"`
	RiskLevelLabel       string `json:"riskLevelLabel"`
	RiskLevelText        string `json:"riskLevelText"`
	Reason               string `json:"reason"`
	RequiresConfirmation bool   `json:"requiresConfirmation"`
}

func ClassifyTerminalCommandRisk(command string) TerminalCommandRiskAssessment {
	assessment := TerminalCommandRiskAssessment{
		Command:        strings.TrimSpace(command),
		RiskLevel:      terminalCommandRiskLow,
		RiskLevelLabel: "low",
		RiskLevelText:  "低风险",
	}

	normalized := strings.ToLower(strings.TrimSpace(command))
	if normalized == "" {
		return assessment
	}

	highRiskKeywords := []string{
		"rm -rf",
		"mkfs",
		"fdisk",
		"dd if=",
		"shutdown",
		"reboot",
		"poweroff",
		"userdel",
		"drop database",
		"drop table",
		"truncate ",
		"passwd ",
		"chmod 777",
		"kubectl delete",
		"kubectl exec",
		"swapon ",
		"systemctl stop",
		"systemctl disable",
		"kill -9",
	}
	for _, keyword := range highRiskKeywords {
		if strings.Contains(normalized, keyword) {
			assessment.IsSensitive = true
			assessment.RiskLevel = terminalCommandRiskHigh
			assessment.RiskLevelLabel = "high"
			assessment.RiskLevelText = "高风险"
			assessment.Reason = "包含高风险操作关键词"
			assessment.RequiresConfirmation = true
			return assessment
		}
	}

	sensitiveKeywords := []string{
		"curl ",
		"wget ",
		"scp ",
		"rsync ",
		"systemctl restart",
		"kubectl rollout restart",
		"kubectl apply",
		"helm upgrade",
		"helm uninstall",
		"chmod ",
		"chown ",
		"mv ",
		"cp ",
		"vim ",
		"vi ",
	}
	for _, keyword := range sensitiveKeywords {
		if strings.Contains(normalized, keyword) {
			assessment.IsSensitive = true
			assessment.RiskLevel = terminalCommandRiskMedium
			assessment.RiskLevelLabel = "medium"
			assessment.RiskLevelText = "中风险"
			assessment.Reason = "包含敏感操作关键词"
			assessment.RequiresConfirmation = true
			return assessment
		}
	}

	return assessment
}
