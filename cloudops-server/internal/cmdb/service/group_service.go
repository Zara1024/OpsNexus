package service

import (
	"strings"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/model"
	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/repository"
)

// GroupService 提供业务分组管理能力。
type GroupService struct {
	repo *repository.GroupRepository
}

// GroupCreateRequest 定义创建分组请求。
type GroupCreateRequest struct {
	ParentID    int64  `json:"parent_id"`
	GroupName   string `json:"group_name" binding:"required,max=128"`
	Description string `json:"description" binding:"max=255"`
	SortOrder   int    `json:"sort_order"`
}

// GroupUpdateRequest 定义更新分组请求。
type GroupUpdateRequest = GroupCreateRequest

// NewGroupService 创建业务分组服务。
func NewGroupService(repo *repository.GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

// CreateGroup 创建业务分组。
func (s *GroupService) CreateGroup(req *GroupCreateRequest) error {
	return s.repo.Create(&model.Group{
		ParentID:    req.ParentID,
		GroupName:   strings.TrimSpace(req.GroupName),
		Description: strings.TrimSpace(req.Description),
		SortOrder:   req.SortOrder,
	})
}

// UpdateGroup 更新业务分组。
func (s *GroupService) UpdateGroup(id int64, req *GroupUpdateRequest) error {
	return s.repo.Update(id, map[string]interface{}{
		"parent_id":   req.ParentID,
		"group_name":  strings.TrimSpace(req.GroupName),
		"description": strings.TrimSpace(req.Description),
		"sort_order":  req.SortOrder,
	})
}

// DeleteGroup 删除业务分组。
func (s *GroupService) DeleteGroup(id int64) error {
	return s.repo.SoftDelete(id)
}

// ListGroupTree 查询分组树。
func (s *GroupService) ListGroupTree() ([]model.GroupTreeNode, error) {
	list, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	return buildGroupTree(list), nil
}

func buildGroupTree(groups []model.Group) []model.GroupTreeNode {
	nodeMap := make(map[int64]*model.GroupTreeNode, len(groups))
	roots := make([]model.GroupTreeNode, 0)

	for _, g := range groups {
		item := g
		nodeMap[item.ID] = &model.GroupTreeNode{
			ID:          item.ID,
			ParentID:    item.ParentID,
			GroupName:   item.GroupName,
			Description: item.Description,
			SortOrder:   item.SortOrder,
			Children:    make([]model.GroupTreeNode, 0),
		}
	}

	for _, node := range nodeMap {
		if node.ParentID == 0 {
			roots = append(roots, *node)
			continue
		}
		if parent, ok := nodeMap[node.ParentID]; ok {
			parent.Children = append(parent.Children, *node)
		}
	}
	return roots
}
