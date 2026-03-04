package service

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/model"
	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/repository"
	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
)

// UserService 提供用户管理能力。
type UserService struct {
	repo *repository.UserRepository
}

// UserCreateRequest 定义创建用户请求。
type UserCreateRequest struct {
	Username     string `json:"username" binding:"required,min=3,max=64"`
	Password     string `json:"password" binding:"required,min=8,max=64"`
	Nickname     string `json:"nickname" binding:"required,max=64"`
	Email        string `json:"email" binding:"max=128"`
	Phone        string `json:"phone" binding:"max=32"`
	Status       int16  `json:"status"`
	DepartmentID int64  `json:"department_id"`
}

// UserUpdateRequest 定义更新用户请求。
type UserUpdateRequest struct {
	Nickname     string `json:"nickname" binding:"required,max=64"`
	Email        string `json:"email" binding:"max=128"`
	Phone        string `json:"phone" binding:"max=32"`
	Status       int16  `json:"status"`
	DepartmentID int64  `json:"department_id"`
}

// ChangePasswordRequest 定义修改密码请求。
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=8,max=64"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=64"`
}

// AssignRolesRequest 定义分配角色请求。
type AssignRolesRequest struct {
	RoleIDs []int64 `json:"role_ids"`
}

// UserListResponse 定义分页结果。
type UserListResponse struct {
	List  []model.User `json:"list"`
	Total int64        `json:"total"`
}

// NewUserService 创建用户服务。
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser 创建用户。
func (s *UserService) CreateUser(req *UserCreateRequest) error {
	if _, err := s.repo.GetByUsername(strings.TrimSpace(req.Username)); err == nil {
		return &apperrors.AppError{Code: 10010, HTTPCode: 400, Message: "用户名已存在"}
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	status := req.Status
	if status == 0 {
		status = 1
	}

	return s.repo.Create(&model.User{
		Username:     strings.TrimSpace(req.Username),
		PasswordHash: string(hash),
		Nickname:     strings.TrimSpace(req.Nickname),
		Email:        strings.TrimSpace(req.Email),
		Phone:        strings.TrimSpace(req.Phone),
		Status:       status,
		DepartmentID: req.DepartmentID,
	})
}

// UpdateUser 更新用户。
func (s *UserService) UpdateUser(id int64, req *UserUpdateRequest) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return apperrors.ErrNotFound
	}
	return s.repo.Update(id, map[string]interface{}{
		"nickname":      strings.TrimSpace(req.Nickname),
		"email":         strings.TrimSpace(req.Email),
		"phone":         strings.TrimSpace(req.Phone),
		"status":        req.Status,
		"department_id": req.DepartmentID,
	})
}

// DeleteUser 删除用户。
func (s *UserService) DeleteUser(id int64) error {
	return s.repo.SoftDelete(id)
}

// ListUsers 分页查询用户。
func (s *UserService) ListUsers(page, pageSize int, keyword string) (*UserListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	list, total, err := s.repo.List(page, pageSize, keyword)
	if err != nil {
		return nil, err
	}
	return &UserListResponse{List: list, Total: total}, nil
}

// ChangePassword 修改密码。
func (s *UserService) ChangePassword(id int64, req *ChangePasswordRequest) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return apperrors.ErrNotFound
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return &apperrors.AppError{Code: 10011, HTTPCode: 400, Message: "旧密码不正确"}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.Update(id, map[string]interface{}{"password_hash": string(hash)})
}

// AssignRoles 分配用户角色。
func (s *UserService) AssignRoles(id int64, req *AssignRolesRequest) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return apperrors.ErrNotFound
	}
	return s.repo.ReplaceRoles(id, req.RoleIDs)
}
