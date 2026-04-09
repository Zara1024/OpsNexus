package model

type WorkOrderSummary struct {
	Total         int `json:"total"`
	Pending       int `json:"pending"`
	Running       int `json:"running"`
	Success       int `json:"success"`
	Failed        int `json:"failed"`
	Canceled      int `json:"canceled"`
	QuickDeploy   int `json:"quickDeploy"`
	ScriptRelease int `json:"scriptRelease"`
	ServiceRelase int `json:"serviceRelease"`
	SQLWorkOrder  int `json:"sqlWorkOrder"`
}

type WorkOrderQuery struct {
	Page     int
	PageSize int
	Type     string
	Status   int
	Keyword  string
}

type WorkOrderItem struct {
	Type            string `json:"type"`
	TypeLabel       string `json:"typeLabel"`
	ID              uint   `json:"id"`
	Title           string `json:"title"`
	AppName         string `json:"appName"`
	ApplicantName   string `json:"applicantName"`
	CurrentHandler  string `json:"currentHandler"`
	Status          int    `json:"status"`
	StatusText      string `json:"statusText"`
	ApprovalStatus  int    `json:"approvalStatus"`
	ExecuteStatus   int    `json:"executeStatus"`
	CanApprove      bool   `json:"canApprove"`
	CanReject       bool   `json:"canReject"`
	CanExecute      bool   `json:"canExecute"`
	BusinessGroupID uint   `json:"businessGroupId"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
	Duration        int64  `json:"duration"`
	DetailHint      string `json:"detailHint"`
	RiskLevel       string `json:"riskLevel,omitempty"`
	RiskText        string `json:"riskText,omitempty"`
	AIDiagnosisPath string `json:"aiDiagnosisPath,omitempty"`
	KnowledgePath   string `json:"knowledgePath,omitempty"`
}

type WorkOrderListResponse struct {
	Total int64           `json:"total"`
	List  []WorkOrderItem `json:"list"`
}

type WorkOrderDetail struct {
	Type       string                   `json:"type"`
	TypeLabel  string                   `json:"typeLabel"`
	ID         uint                     `json:"id"`
	Title      string                   `json:"title"`
	Status     int                      `json:"status"`
	StatusText string                   `json:"statusText"`
	CanApprove bool                     `json:"canApprove"`
	CanReject  bool                     `json:"canReject"`
	CanExecute bool                     `json:"canExecute"`
	Basic      map[string]interface{}   `json:"basic"`
	Items      []map[string]interface{} `json:"items"`
}

type ScriptWorkOrderCreateRequest struct {
	Title           string `json:"title" binding:"required"`
	Reason          string `json:"reason" binding:"required"`
	BusinessGroupID uint   `json:"businessGroupId"`
	AppID           uint   `json:"appId"`
	AppName         string `json:"appName" binding:"required"`
	AppCode         string `json:"appCode" binding:"required"`
	ExecuteDir      string `json:"executeDir" binding:"required"`
	ScriptContent   string `json:"scriptContent" binding:"required"`
	ServerHostID    uint   `json:"serverHostId"`
}

type WorkOrderActionRequest struct {
	Comment string `json:"comment"`
}
