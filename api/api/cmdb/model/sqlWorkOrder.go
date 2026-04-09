package model

import (
	"dodevops-api/common/util"
	"time"
)

type CmdbSQLWorkOrder struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	OrderNo         string     `gorm:"size:64;not null;uniqueIndex" json:"orderNo"`
	Title           string     `gorm:"size:255;not null" json:"title"`
	Reason          string     `gorm:"type:text" json:"reason"`
	DatabaseID      uint       `gorm:"not null;index" json:"databaseId"`
	DatabaseName    string     `gorm:"size:128;not null;index" json:"databaseName"`
	InstanceName    string     `gorm:"size:128;not null" json:"instanceName"`
	InstanceHost    string     `gorm:"size:128;not null" json:"instanceHost"`
	AccountID       uint       `gorm:"not null" json:"accountId"`
	OperationType   string     `gorm:"size:32;not null;index" json:"operationType"`
	SQLContent      string     `gorm:"type:text;not null" json:"sqlContent"`
	RiskLevel       int        `gorm:"default:0;index" json:"riskLevel"`
	RiskSummary     string     `gorm:"type:text" json:"riskSummary"`
	AffectedTables  string     `gorm:"size:500" json:"affectedTables"`
	RollbackSQL     string     `gorm:"type:text" json:"rollbackSql"`
	RollbackHint    string     `gorm:"type:text" json:"rollbackHint"`
	BackupPreview   string     `gorm:"type:longtext" json:"backupPreview"`
	BackupRowCount  int64      `gorm:"default:0" json:"backupRowCount"`
	RequireApproval bool       `gorm:"default:true" json:"requireApproval"`
	ApplicantID     uint       `gorm:"not null;index" json:"applicantId"`
	ApplicantName   string     `gorm:"size:64;not null;index" json:"applicantName"`
	ApproverID      uint       `gorm:"default:0;index" json:"approverId"`
	ApproverName    string     `gorm:"size:64" json:"approverName"`
	ApprovalComment string     `gorm:"type:text" json:"approvalComment"`
	ApprovalTime    *time.Time `json:"approvalTime"`
	ExecutorID      uint       `gorm:"default:0" json:"executorId"`
	ExecutorName    string     `gorm:"size:64" json:"executorName"`
	ExecutionStart  *time.Time `json:"executionStart"`
	ExecutionEnd    *time.Time `json:"executionEnd"`
	ExecutionTime   int64      `gorm:"default:0" json:"executionTime"`
	AffectedRows    int64      `gorm:"default:0" json:"affectedRows"`
	ReturnedRows    int64      `gorm:"default:0" json:"returnedRows"`
	ResultStatus    string     `gorm:"size:32;default:'PENDING'" json:"resultStatus"`
	ResultMessage   string     `gorm:"type:text" json:"resultMessage"`
	ClientIP        string     `gorm:"size:64" json:"clientIp"`
	Status          int        `gorm:"default:1;index" json:"status"`
	CreateTime      util.HTime `gorm:"not null" json:"createTime"`
	UpdateTime      util.HTime `gorm:"not null" json:"updateTime"`
}

func (CmdbSQLWorkOrder) TableName() string {
	return "cmdb_sql_work_order"
}

type CmdbSQLWorkOrderCreateRequest struct {
	DatabaseID   uint   `json:"databaseId" binding:"required"`
	DatabaseName string `json:"databaseName" binding:"required"`
	Title        string `json:"title"`
	Reason       string `json:"reason"`
	SQL          string `json:"sql" binding:"required"`
}

type CmdbSQLWorkOrderActionRequest struct {
	Comment string `json:"comment"`
}

type CmdbSQLWorkOrderQuery struct {
	Page       int
	PageSize   int
	Status     int
	DatabaseID uint
	Keyword    string
}

type CmdbSQLWorkOrderSummary struct {
	Total     int64 `json:"total"`
	Pending   int64 `json:"pending"`
	Approved  int64 `json:"approved"`
	Rejected  int64 `json:"rejected"`
	Executing int64 `json:"executing"`
	Succeeded int64 `json:"succeeded"`
	Failed    int64 `json:"failed"`
	HighRisk  int64 `json:"highRisk"`
}

type CmdbSQLWorkOrderListResponse struct {
	List     []CmdbSQLWorkOrder `json:"list"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"pageSize"`
}
