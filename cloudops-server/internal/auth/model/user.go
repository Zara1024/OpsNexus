package model

import "time"

// User 定义系统用户模型。
type User struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username     string     `gorm:"column:username;size:64;uniqueIndex;not null" json:"username"`
	PasswordHash string     `gorm:"column:password_hash;size:255;not null" json:"-"`
	Nickname     string     `gorm:"column:nickname;size:64;not null" json:"nickname"`
	Email        string     `gorm:"column:email;size:128;not null" json:"email"`
	Phone        string     `gorm:"column:phone;size:32;not null" json:"phone"`
	Status       int16      `gorm:"column:status;not null" json:"status"`
	DepartmentID int64      `gorm:"column:department_id;not null;default:0" json:"department_id"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"-"`
}

// TableName 指定用户表名。
func (User) TableName() string {
	return "sys_users"
}

// UserRole 定义用户与角色关联。
type UserRole struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    int64      `gorm:"column:user_id;not null" json:"user_id"`
	RoleID    int64      `gorm:"column:role_id;not null" json:"role_id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`
}

// TableName 指定用户角色关联表名。
func (UserRole) TableName() string {
	return "sys_user_roles"
}
