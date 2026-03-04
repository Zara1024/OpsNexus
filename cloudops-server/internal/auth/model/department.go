package model

import "time"

// Department 定义部门模型。
type Department struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ParentID  int64      `gorm:"column:parent_id;not null" json:"parent_id"`
	DeptName  string     `gorm:"column:dept_name;size:64;not null" json:"dept_name"`
	SortOrder int        `gorm:"column:sort_order;not null;default:0" json:"sort_order"`
	Leader    string     `gorm:"column:leader;size:64;not null;default:''" json:"leader"`
	Phone     string     `gorm:"column:phone;size:32;not null;default:''" json:"phone"`
	Email     string     `gorm:"column:email;size:128;not null;default:''" json:"email"`
	Status    int16      `gorm:"column:status;not null;default:1" json:"status"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`
}

// TableName 指定部门表名。
func (Department) TableName() string {
	return "sys_departments"
}
