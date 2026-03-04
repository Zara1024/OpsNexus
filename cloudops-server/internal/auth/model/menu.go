package model

import "time"

// Menu 定义菜单模型，支持目录/菜单/按钮。
type Menu struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ParentID   int64      `gorm:"column:parent_id;not null" json:"parent_id"`
	MenuName   string     `gorm:"column:menu_name;size:64;not null" json:"menu_name"`
	RoutePath  string     `gorm:"column:route_path;size:128;not null" json:"route_path"`
	Component  string     `gorm:"column:component;size:128;not null" json:"component"`
	Icon       string     `gorm:"column:icon;size:64;not null" json:"icon"`
	SortOrder  int        `gorm:"column:sort_order;not null" json:"sort_order"`
	MenuType   string     `gorm:"column:menu_type;size:16;not null;default:menu" json:"menu_type"`
	Permission string     `gorm:"column:permission;size:128;not null;default:''" json:"permission"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"-"`
}

// TableName 指定菜单表名。
func (Menu) TableName() string {
	return "sys_menus"
}

// MenuTreeNode 定义前端菜单树节点。
type MenuTreeNode struct {
	ID         int64          `json:"id"`
	ParentID   int64          `json:"parent_id"`
	MenuName   string         `json:"menu_name"`
	RoutePath  string         `json:"route_path"`
	Component  string         `json:"component"`
	Icon       string         `json:"icon"`
	SortOrder  int            `json:"sort_order"`
	MenuType   string         `json:"menu_type"`
	Permission string         `json:"permission"`
	Children   []MenuTreeNode `json:"children"`
}
