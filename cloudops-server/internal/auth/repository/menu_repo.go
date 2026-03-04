package repository

import (
	"gorm.io/gorm"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/model"
)

// MenuRepository 定义菜单数据访问能力。
type MenuRepository struct {
	db *gorm.DB
}

// NewMenuRepository 创建菜单仓储。
func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

// Create 创建菜单。
func (r *MenuRepository) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

// Update 更新菜单。
func (r *MenuRepository) Update(id int64, updates map[string]interface{}) error {
	return r.db.Model(&model.Menu{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error
}

// SoftDelete 软删除菜单。
func (r *MenuRepository) SoftDelete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.Menu{}).Error
}

// GetByID 按ID查询菜单。
func (r *MenuRepository) GetByID(id int64) (*model.Menu, error) {
	var menu model.Menu
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&menu).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

// List 查询全部菜单。
func (r *MenuRepository) List() ([]model.Menu, error) {
	menus := make([]model.Menu, 0)
	if err := r.db.Where("deleted_at IS NULL").Order("sort_order ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// ListByRoleIDs 查询角色可访问菜单。
func (r *MenuRepository) ListByRoleIDs(roleIDs []int64) ([]model.Menu, error) {
	menus := make([]model.Menu, 0)
	if len(roleIDs) == 0 {
		return menus, nil
	}
	err := r.db.Table("sys_menus AS m").
		Select("DISTINCT m.*").
		Joins("JOIN sys_role_menus rm ON rm.menu_id = m.id AND rm.deleted_at IS NULL").
		Where("m.deleted_at IS NULL AND rm.role_id IN ?", roleIDs).
		Order("m.sort_order ASC, m.id ASC").
		Scan(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

// ListPermissionsByRoleIDs 查询角色权限标识。
func (r *MenuRepository) ListPermissionsByRoleIDs(roleIDs []int64) ([]string, error) {
	permissions := make([]string, 0)
	if len(roleIDs) == 0 {
		return permissions, nil
	}
	err := r.db.Table("sys_menus AS m").
		Select("DISTINCT m.permission").
		Joins("JOIN sys_role_menus rm ON rm.menu_id = m.id AND rm.deleted_at IS NULL").
		Where("m.deleted_at IS NULL AND rm.role_id IN ? AND m.permission <> ''", roleIDs).
		Order("m.permission ASC").
		Scan(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}
