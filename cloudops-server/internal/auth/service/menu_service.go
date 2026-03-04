package service

import (
	"strings"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/model"
	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/repository"
)

// MenuService 提供菜单管理能力。
type MenuService struct {
	repo *repository.MenuRepository
}

// MenuCreateRequest 定义创建菜单请求。
type MenuCreateRequest struct {
	ParentID   int64  `json:"parent_id"`
	MenuName   string `json:"menu_name" binding:"required,max=64"`
	RoutePath  string `json:"route_path" binding:"max=128"`
	Component  string `json:"component" binding:"max=128"`
	Icon       string `json:"icon" binding:"max=64"`
	SortOrder  int    `json:"sort_order"`
	MenuType   string `json:"menu_type" binding:"required,oneof=dir menu button"`
	Permission string `json:"permission" binding:"max=128"`
}

// MenuUpdateRequest 定义更新菜单请求。
type MenuUpdateRequest = MenuCreateRequest

// NewMenuService 创建菜单服务。
func NewMenuService(repo *repository.MenuRepository) *MenuService {
	return &MenuService{repo: repo}
}

// CreateMenu 创建菜单。
func (s *MenuService) CreateMenu(req *MenuCreateRequest) error {
	return s.repo.Create(&model.Menu{
		ParentID:   req.ParentID,
		MenuName:   strings.TrimSpace(req.MenuName),
		RoutePath:  strings.TrimSpace(req.RoutePath),
		Component:  strings.TrimSpace(req.Component),
		Icon:       strings.TrimSpace(req.Icon),
		SortOrder:  req.SortOrder,
		MenuType:   strings.TrimSpace(req.MenuType),
		Permission: strings.TrimSpace(req.Permission),
	})
}

// UpdateMenu 更新菜单。
func (s *MenuService) UpdateMenu(id int64, req *MenuUpdateRequest) error {
	return s.repo.Update(id, map[string]interface{}{
		"parent_id":  req.ParentID,
		"menu_name":  strings.TrimSpace(req.MenuName),
		"route_path": strings.TrimSpace(req.RoutePath),
		"component":  strings.TrimSpace(req.Component),
		"icon":       strings.TrimSpace(req.Icon),
		"sort_order": req.SortOrder,
		"menu_type":  strings.TrimSpace(req.MenuType),
		"permission": strings.TrimSpace(req.Permission),
	})
}

// DeleteMenu 删除菜单。
func (s *MenuService) DeleteMenu(id int64) error {
	return s.repo.SoftDelete(id)
}

// ListMenus 获取菜单列表。
func (s *MenuService) ListMenus() ([]model.Menu, error) {
	return s.repo.List()
}

// MenuTree 获取菜单树。
func (s *MenuService) MenuTree() ([]model.MenuTreeNode, error) {
	menus, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	return buildMenuTree(menus), nil
}
