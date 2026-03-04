package model

import "time"

// Group 定义 CMDB 业务分组。
type Group struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ParentID    int64      `gorm:"column:parent_id;not null;default:0" json:"parent_id"`
	GroupName   string     `gorm:"column:group_name;size:128;not null" json:"group_name"`
	Description string     `gorm:"column:description;size:255;not null;default:''" json:"description"`
	SortOrder   int        `gorm:"column:sort_order;not null;default:0" json:"sort_order"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"-"`
}

// TableName 指定表名。
func (Group) TableName() string {
	return "cmdb_groups"
}

// GroupTreeNode 定义分组树节点。
type GroupTreeNode struct {
	ID          int64           `json:"id"`
	ParentID    int64           `json:"parent_id"`
	GroupName   string          `json:"group_name"`
	Description string          `json:"description"`
	SortOrder   int             `json:"sort_order"`
	Children    []GroupTreeNode `json:"children"`
}
