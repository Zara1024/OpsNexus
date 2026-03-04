package service

import (
	"strings"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/model"
	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/repository"
)

// RoleService 提供角色管理能力。
type RoleService struct {
	repo *repository.RoleRepository
}

// RoleCreateRequest 定义创建角色请求。
type RoleCreateRequest struct {
	RoleCode    string `json:"role_code" binding:"required,max=64"`
	RoleName    string `json:"role_name" binding:"required,max=64"`
	Description string `json:"description" binding:"max=255"`
}

// RoleUpdateRequest 定义更新角色请求。
type RoleUpdateRequest struct {
	RoleName    string `json:"role_name" binding:"required,max=64"`
	Description string `json:"description" binding:"max=255"`
}

// AssignMenusRequest 定义分配菜单请求。
type AssignMenusRequest struct {
	MenuIDs []int64 `json:"menu_ids"`
}

// NewRoleService 创建角色服务。
func NewRoleService(repo *repository.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

// CreateRole 创建角色。
func (s *RoleService) CreateRole(req *RoleCreateRequest) error {
	return s.repo.Create(&model.Role{
		RoleCode:    strings.TrimSpace(req.RoleCode),
		RoleName:    strings.TrimSpace(req.RoleName),
		Description: strings.TrimSpace(req.Description),
	})
}

// UpdateRole 更新角色。
func (s *RoleService) UpdateRole(id int64, req *RoleUpdateRequest) error {
	return s.repo.Update(id, map[string]interface{}{
		"role_name":   strings.TrimSpace(req.RoleName),
		"description": strings.TrimSpace(req.Description),
	})
}

// DeleteRole 删除角色。
func (s *RoleService) DeleteRole(id int64) error {
	return s.repo.SoftDelete(id)
}

// ListRoles 查询角色列表。
func (s *RoleService) ListRoles() ([]model.Role, error) {
	return s.repo.List()
}

// AssignMenus 分配角色菜单。
func (s *RoleService) AssignMenus(id int64, req *AssignMenusRequest) error {
	return s.repo.ReplaceMenus(id, req.MenuIDs)
}
