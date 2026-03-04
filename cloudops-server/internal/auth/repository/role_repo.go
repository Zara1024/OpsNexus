package repository

import (
	"gorm.io/gorm"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/model"
)

// RoleRepository 定义角色数据访问能力。
type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓储。
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// Create 创建角色。
func (r *RoleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

// Update 更新角色。
func (r *RoleRepository) Update(id int64, updates map[string]interface{}) error {
	return r.db.Model(&model.Role{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error
}

// SoftDelete 软删除角色。
func (r *RoleRepository) SoftDelete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.Role{}).Error
}

// GetByID 按ID查询角色。
func (r *RoleRepository) GetByID(id int64) (*model.Role, error) {
	var role model.Role
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// List 查询角色列表。
func (r *RoleRepository) List() ([]model.Role, error) {
	roles := make([]model.Role, 0)
	if err := r.db.Where("deleted_at IS NULL").Order("id DESC").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// ListByIDs 按ID批量查询角色。
func (r *RoleRepository) ListByIDs(roleIDs []int64) ([]model.Role, error) {
	roles := make([]model.Role, 0)
	if len(roleIDs) == 0 {
		return roles, nil
	}
	if err := r.db.Where("deleted_at IS NULL AND id IN ?", roleIDs).Order("id ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// ReplaceMenus 重置角色菜单。
func (r *RoleRepository) ReplaceMenus(roleID int64, menuIDs []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
			return err
		}
		if len(menuIDs) == 0 {
			return nil
		}
		rows := make([]model.RoleMenu, 0, len(menuIDs))
		for _, menuID := range menuIDs {
			rows = append(rows, model.RoleMenu{RoleID: roleID, MenuID: menuID})
		}
		return tx.Create(&rows).Error
	})
}

// GetMenuIDsByRoleID 查询角色菜单ID。
func (r *RoleRepository) GetMenuIDsByRoleID(roleID int64) ([]int64, error) {
	menuIDs := make([]int64, 0)
	err := r.db.Table("sys_role_menus").
		Select("menu_id").
		Where("role_id = ? AND deleted_at IS NULL", roleID).
		Order("menu_id ASC").
		Scan(&menuIDs).Error
	if err != nil {
		return nil, err
	}
	return menuIDs, nil
}
