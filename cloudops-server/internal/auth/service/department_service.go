package service

import (
	"strings"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/model"
	"github.com/Zara1024/OpsNexus/cloudops-server/internal/auth/repository"
)

// DepartmentService 提供部门管理能力。
type DepartmentService struct {
	repo *repository.DepartmentRepository
}

// DepartmentCreateRequest 定义创建部门请求。
type DepartmentCreateRequest struct {
	ParentID  int64  `json:"parent_id"`
	DeptName  string `json:"dept_name" binding:"required,max=64"`
	SortOrder int    `json:"sort_order"`
	Leader    string `json:"leader" binding:"max=64"`
	Phone     string `json:"phone" binding:"max=32"`
	Email     string `json:"email" binding:"max=128"`
	Status    int16  `json:"status"`
}

// DepartmentUpdateRequest 定义更新部门请求。
type DepartmentUpdateRequest = DepartmentCreateRequest

// NewDepartmentService 创建部门服务。
func NewDepartmentService(repo *repository.DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

// CreateDepartment 创建部门。
func (s *DepartmentService) CreateDepartment(req *DepartmentCreateRequest) error {
	status := req.Status
	if status == 0 {
		status = 1
	}
	return s.repo.Create(&model.Department{
		ParentID:  req.ParentID,
		DeptName:  strings.TrimSpace(req.DeptName),
		SortOrder: req.SortOrder,
		Leader:    strings.TrimSpace(req.Leader),
		Phone:     strings.TrimSpace(req.Phone),
		Email:     strings.TrimSpace(req.Email),
		Status:    status,
	})
}

// UpdateDepartment 更新部门。
func (s *DepartmentService) UpdateDepartment(id int64, req *DepartmentUpdateRequest) error {
	return s.repo.Update(id, map[string]interface{}{
		"parent_id":  req.ParentID,
		"dept_name":  strings.TrimSpace(req.DeptName),
		"sort_order": req.SortOrder,
		"leader":     strings.TrimSpace(req.Leader),
		"phone":      strings.TrimSpace(req.Phone),
		"email":      strings.TrimSpace(req.Email),
		"status":     req.Status,
	})
}

// DeleteDepartment 删除部门。
func (s *DepartmentService) DeleteDepartment(id int64) error {
	return s.repo.SoftDelete(id)
}

// ListDepartments 查询部门。
func (s *DepartmentService) ListDepartments() ([]model.Department, error) {
	return s.repo.List()
}
