package model

import "time"

// Role 定义角色模型。
type Role struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RoleCode    string     `gorm:"column:role_code;size:64;uniqueIndex;not null" json:"role_code"`
	RoleName    string     `gorm:"column:role_name;size:64;not null" json:"role_name"`
	Description string     `gorm:"column:description;size:255;not null" json:"description"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"-"`
}

// TableName 指定角色表名。
func (Role) TableName() string {
	return "sys_roles"
}

// RoleMenu 定义角色菜单关联模型。
type RoleMenu struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RoleID    int64      `gorm:"column:role_id;not null" json:"role_id"`
	MenuID    int64      `gorm:"column:menu_id;not null" json:"menu_id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`
}

// TableName 指定角色菜单关联表名。
func (RoleMenu) TableName() string {
	return "sys_role_menus"
}
