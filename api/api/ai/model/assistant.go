package model

type AIAssistantChatRequest struct {
	SessionID string `json:"sessionId"`
	Message   string `json:"message"`
}

type AIAssistantAction struct {
	Label   string `json:"label"`
	Message string `json:"message"`
	Kind    string `json:"kind"`
}

type AIAssistantHost struct {
	ID          uint   `json:"id"`
	HostName    string `json:"hostName"`
	GroupName   string `json:"groupName"`
	PrivateIP   string `json:"privateIp"`
	PublicIP    string `json:"publicIp"`
	SSHIP       string `json:"sshIp"`
	OS          string `json:"os"`
	CPU         string `json:"cpu"`
	Memory      string `json:"memory"`
	Disk        string `json:"disk"`
	Status      int    `json:"status"`
	StatusText  string `json:"statusText"`
	SupportsSSH bool   `json:"supportsSsh"`
}

type AIAssistantCommandResult struct {
	HostID   uint   `json:"hostId"`
	HostName string `json:"hostName"`
	Command  string `json:"command"`
	Output   string `json:"output"`
	Success  bool   `json:"success"`
}

type AIAssistantInspectionCheck struct {
	Name    string `json:"name"`
	Command string `json:"command"`
	Status  string `json:"status"`
	Output  string `json:"output"`
}

type AIAssistantInspectionResult struct {
	HostID       uint                         `json:"hostId"`
	HostName     string                       `json:"hostName"`
	TemplateID   uint                         `json:"templateId"`
	TemplateName string                       `json:"templateName"`
	Summary      string                       `json:"summary"`
	Report       string                       `json:"report"`
	Checks       []AIAssistantInspectionCheck `json:"checks"`
}

type AIAssistantChatResponse struct {
	SessionID           string                               `json:"sessionId"`
	Intent              string                               `json:"intent"`
	Provider            string                               `json:"provider,omitempty"`
	Model               string                               `json:"model,omitempty"`
	UsedLLM             bool                                 `json:"usedLlm"`
	FallbackReason      string                               `json:"fallbackReason,omitempty"`
	AssistantMessage    string                               `json:"assistantMessage"`
	Suggestions         []string                             `json:"suggestions"`
	Actions             []AIAssistantAction                  `json:"actions"`
	Context             *AIAssistantContext                  `json:"context,omitempty"`
	ToolSteps           []AIAssistantToolStep                `json:"toolSteps,omitempty"`
	PendingConfirmation *AIAssistantPendingConfirmation      `json:"pendingConfirmation,omitempty"`
	AvailableTemplates  []AIAssistantInspectionTemplate      `json:"availableTemplates,omitempty"`
	RecentReports       []AIAssistantInspectionReportSummary `json:"recentReports,omitempty"`
	HostMatches         []AIAssistantHost                    `json:"hostMatches"`
	CommandResult       *AIAssistantCommandResult            `json:"commandResult"`
	InspectionResult    *AIAssistantInspectionResult         `json:"inspectionResult"`
}
