package model

import (
	"dodevops-api/common/util"
)

// CmdbSQL 定义CMDB中的数据库模型
type CmdbSQL struct {
	ID              uint       `gorm:"column:id;primaryKey" json:"id"`                          // 主键ID
	Name            string     `gorm:"column:name;size:100;not null" json:"name"`               // 资产名称/实例名
	Address         string     `gorm:"column:address;size:128" json:"address"`                  // 资产地址
	Platform        string     `gorm:"column:platform;size:100" json:"platform"`                // 平台/版本
	DefaultDatabase string     `gorm:"column:default_database;size:100" json:"defaultDatabase"` // 默认schema/数据库
	Type            int        `gorm:"column:type;type:integer;not null" json:"type"`           // 数据库类型(1=MySQL 2=PostgreSQL 3=Redis 4=MongoDB 5=Elasticsearch)
	AccountID       uint       `gorm:"column:account_id;not null" json:"accountId"`             // 所属账号ID
	GroupID         uint       `gorm:"column:group_id;not null" json:"groupId"`                 // 所属业务组ID
	ProtocolGroup   string     `gorm:"column:protocol_group;size:100" json:"protocolGroup"`     // 协议组
	Tags            string     `gorm:"column:tags;size:255" json:"tags"`                        // 标签(多个标签用逗号分隔)
	IsActive        bool       `gorm:"column:is_active;default:true" json:"isActive"`           // 是否启用
	Remark          string     `gorm:"column:description;size:500" json:"remark"`               // 备注
	Description     string     `gorm:"-" json:"description,omitempty"`                          // 兼容旧字段
	CreatedAt       util.HTime `gorm:"column:created_at" json:"createdAt"`                      // 创建时间
	UpdatedAt       util.HTime `gorm:"column:updated_at" json:"updatedAt"`                      // 更新时间
}

func (CmdbSQL) TableName() string {
	return "cmdb_sql"
}

type CreateCmdbSQLDto struct {
	Name            string `json:"name" binding:"required"`
	Address         string `json:"address" binding:"required"`
	Platform        string `json:"platform"`
	DefaultDatabase string `json:"defaultDatabase"`
	Type            int    `json:"type" binding:"required"`
	AccountID       uint   `json:"accountId" binding:"required"`
	GroupID         uint   `json:"groupId" binding:"required"`
	ProtocolGroup   string `json:"protocolGroup"`
	Tags            string `json:"tags"`
	IsActive        *bool  `json:"isActive"`
	Remark          string `json:"remark"`
	Description     string `json:"description"`
}

type UpdateCmdbSQLDto struct {
	ID              uint   `json:"id" binding:"required"`
	Name            string `json:"name" binding:"required"`
	Address         string `json:"address" binding:"required"`
	Platform        string `json:"platform"`
	DefaultDatabase string `json:"defaultDatabase"`
	Type            int    `json:"type" binding:"required"`
	AccountID       uint   `json:"accountId" binding:"required"`
	GroupID         uint   `json:"groupId" binding:"required"`
	ProtocolGroup   string `json:"protocolGroup"`
	Tags            string `json:"tags"`
	IsActive        *bool  `json:"isActive"`
	Remark          string `json:"remark"`
	Description     string `json:"description"`
}
