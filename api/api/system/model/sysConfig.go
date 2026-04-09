package model

import "dodevops-api/common/util"

type SysConfig struct {
	ID         uint       `gorm:"column:id;primaryKey" json:"id"`
	ConfigKey  string     `gorm:"column:config_key" json:"configKey"`
	ConfigType string     `gorm:"column:config_type" json:"configType"`
	ConfigData string     `gorm:"column:config_data" json:"configData"`
	Status     int        `gorm:"column:status" json:"status"`
	Remark     string     `gorm:"column:remark" json:"remark"`
	CreateTime util.HTime `gorm:"column:create_time" json:"createTime"`
	UpdateTime util.HTime `gorm:"column:update_time" json:"updateTime"`
}

func (SysConfig) TableName() string {
	return "sys_config"
}

type LDAPConfig struct {
	Enable          bool              `json:"enable"`
	Host            string            `json:"host"`
	Port            int               `json:"port"`
	BaseDN          string            `json:"baseDn"`
	BindUser        string            `json:"bindUser"`
	BindPass        string            `json:"bindPass"`
	AuthFilter      string            `json:"authFilter"`
	CoverAttributes bool              `json:"coverAttributes"`
	TLS             bool              `json:"tls"`
	StartTLS        bool              `json:"startTLS"`
	DefaultRoles    []string          `json:"defaultRoles"`
	DefaultRoleID   uint              `json:"defaultRoleId"`
	GroupFilter     string            `json:"groupFilter"`
	GroupNameAttr   string            `json:"groupNameAttr"`
	RoleMappings    []LDAPRoleMapping `json:"roleMappings"`
	Attributes      LDAPAttributes    `json:"attributes"`
	Remark          string            `json:"remark"`
}

type LDAPRoleMapping struct {
	GroupName string `json:"groupName"`
	RoleID    uint   `json:"roleId"`
}

type LDAPAttributes struct {
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}
