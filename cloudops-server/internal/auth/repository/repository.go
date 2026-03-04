package repository

import "gorm.io/gorm"

// Repositories 聚合auth模块仓储。
type Repositories struct {
	User       *UserRepository
	Role       *RoleRepository
	Menu       *MenuRepository
	Department *DepartmentRepository
}

// New 创建auth模块仓储集合。
func New(db *gorm.DB) *Repositories {
	return &Repositories{
		User:       NewUserRepository(db),
		Role:       NewRoleRepository(db),
		Menu:       NewMenuRepository(db),
		Department: NewDepartmentRepository(db),
	}
}
