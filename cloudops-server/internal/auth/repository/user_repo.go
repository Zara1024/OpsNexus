package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/model"
)

// UserRepository 定义用户数据访问能力。
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储。
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户。
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Update 更新用户。
func (r *UserRepository) Update(id int64, updates map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error
}

// SoftDelete 软删除用户。
func (r *UserRepository) SoftDelete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.User{}).Error
}

// GetByID 按ID查询用户。
func (r *UserRepository) GetByID(id int64) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 按用户名查询用户。
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ? AND deleted_at IS NULL", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// List 分页查询用户。
func (r *UserRepository) List(page, pageSize int, keyword string) ([]model.User, int64, error) {
	var (
		users []model.User
		total int64
	)

	query := r.db.Model(&model.User{}).Where("deleted_at IS NULL")
	if keyword != "" {
		kw := fmt.Sprintf("%%%s%%", keyword)
		query = query.Where("username ILIKE ? OR nickname ILIKE ? OR email ILIKE ?", kw, kw, kw)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ReplaceRoles 重置用户角色。
func (r *UserRepository) ReplaceRoles(userID int64, roleIDs []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}
		if len(roleIDs) == 0 {
			return nil
		}

		rows := make([]model.UserRole, 0, len(roleIDs))
		for _, roleID := range roleIDs {
			rows = append(rows, model.UserRole{UserID: userID, RoleID: roleID})
		}
		return tx.Create(&rows).Error
	})
}

// GetRoleCodesByUserID 查询用户角色编码。
func (r *UserRepository) GetRoleCodesByUserID(userID int64) ([]string, error) {
	roleCodes := make([]string, 0)
	err := r.db.Table("sys_roles AS r").
		Select("r.role_code").
		Joins("JOIN sys_user_roles ur ON ur.role_id = r.id AND ur.deleted_at IS NULL").
		Where("ur.user_id = ? AND r.deleted_at IS NULL", userID).
		Order("r.id ASC").
		Scan(&roleCodes).Error
	if err != nil {
		return nil, err
	}
	return roleCodes, nil
}

// GetRoleIDsByUserID 查询用户角色ID集合。
func (r *UserRepository) GetRoleIDsByUserID(userID int64) ([]int64, error) {
	roleIDs := make([]int64, 0)
	err := r.db.Table("sys_user_roles").
		Select("role_id").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("role_id ASC").
		Scan(&roleIDs).Error
	if err != nil {
		return nil, err
	}
	return roleIDs, nil
}
