package model

import "dodevops-api/common/util"

type CmdbDevice struct {
	ID            uint       `gorm:"column:id;primaryKey" json:"id"`
	Name          string     `gorm:"column:name;size:100;not null" json:"name"`
	Address       string     `gorm:"column:address;size:128;not null" json:"address"`
	Platform      string     `gorm:"column:platform;size:100" json:"platform"`
	GroupID       uint       `gorm:"column:group_id;not null" json:"groupId"`
	AccountID     uint       `gorm:"column:account_id;not null" json:"accountId"`
	ProtocolGroup string     `gorm:"column:protocol_group;size:50" json:"protocolGroup"`
	Tags          string     `gorm:"column:tags;size:255" json:"tags"`
	IsActive      bool       `gorm:"column:is_active;default:true" json:"isActive"`
	Remark        string     `gorm:"column:remark;size:500" json:"remark"`
	DeviceType    string     `gorm:"column:device_type;size:50" json:"deviceType"`
	SSHPort       int        `gorm:"column:ssh_port;default:22" json:"sshPort"`
	TelnetPort    int        `gorm:"column:telnet_port;default:23" json:"telnetPort"`
	WebURL        string     `gorm:"column:web_url;size:255" json:"webUrl"`
	CreateTime    util.HTime `gorm:"column:create_time" json:"createTime"`
	UpdateTime    util.HTime `gorm:"column:update_time" json:"updateTime"`
}

func (CmdbDevice) TableName() string {
	return "cmdb_device"
}

type CreateCmdbDeviceDto struct {
	Name          string `json:"name" binding:"required"`
	Address       string `json:"address" binding:"required"`
	Platform      string `json:"platform"`
	GroupID       uint   `json:"groupId" binding:"required"`
	AccountID     uint   `json:"accountId" binding:"required"`
	ProtocolGroup string `json:"protocolGroup"`
	Tags          string `json:"tags"`
	IsActive      *bool  `json:"isActive"`
	Remark        string `json:"remark"`
	DeviceType    string `json:"deviceType"`
	SSHPort       int    `json:"sshPort"`
	TelnetPort    int    `json:"telnetPort"`
	WebURL        string `json:"webUrl"`
}

type UpdateCmdbDeviceDto struct {
	ID            uint   `json:"id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Address       string `json:"address" binding:"required"`
	Platform      string `json:"platform"`
	GroupID       uint   `json:"groupId" binding:"required"`
	AccountID     uint   `json:"accountId" binding:"required"`
	ProtocolGroup string `json:"protocolGroup"`
	Tags          string `json:"tags"`
	IsActive      *bool  `json:"isActive"`
	Remark        string `json:"remark"`
	DeviceType    string `json:"deviceType"`
	SSHPort       int    `json:"sshPort"`
	TelnetPort    int    `json:"telnetPort"`
	WebURL        string `json:"webUrl"`
}

type CmdbDeviceIDDto struct {
	ID uint `json:"id" binding:"required"`
}

type BatchCmdbDeviceConnectivityDto struct {
	DeviceIDs []uint `json:"deviceIds" binding:"required"`
}

type CmdbDeviceConnectivityCheckItem struct {
	DeviceID       uint   `json:"deviceId"`
	Name           string `json:"name"`
	Protocol       string `json:"protocol"`
	Address        string `json:"address"`
	DisplayAddress string `json:"displayAddress"`
	Reachable      bool   `json:"reachable"`
	Status         string `json:"status"`
	Reason         string `json:"reason"`
}

type BatchCmdbDeviceConnectivityResult struct {
	Total       int                               `json:"total"`
	Reachable   int                               `json:"reachable"`
	Unreachable int                               `json:"unreachable"`
	Items       []CmdbDeviceConnectivityCheckItem `json:"items"`
}
