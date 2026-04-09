package service

import (
	"strings"
	"time"

	"dodevops-api/api/task/dao"
	"dodevops-api/api/task/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IConfigAnsibleService interface {
	Create(c *gin.Context, req *CreateConfigRequest)
	Update(c *gin.Context, id uint, req *UpdateConfigRequest)
	Delete(c *gin.Context, id uint)
	Get(c *gin.Context, id uint)
	List(c *gin.Context, page, size int, name string, configType int)
}

type ConfigAnsibleServiceImpl struct {
	dao *dao.ConfigAnsibleDao
}

func NewConfigAnsibleService(db *gorm.DB) IConfigAnsibleService {
	return &ConfigAnsibleServiceImpl{
		dao: dao.NewConfigAnsibleDao(db),
	}
}

type CreateConfigRequest struct {
	Name    string `json:"name" binding:"required"`
	Type    int    `json:"type" binding:"required,oneof=1 2 3 4"`
	Content string `json:"content" binding:"required"`
	Remark  string `json:"remark"`
}

type UpdateConfigRequest struct {
	Name    string `json:"name"`
	Type    int    `json:"type"`
	Content string `json:"content"`
	Remark  string `json:"remark"`
}

func (s *ConfigAnsibleServiceImpl) Create(c *gin.Context, req *CreateConfigRequest) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		result.Failed(c, 400, "配置名称不能为空")
		return
	}

	existing, err := s.dao.List(1, 20, name, 0)
	if err != nil {
		result.Failed(c, 500, "检查配置名称失败: "+err.Error())
		return
	}
	for _, item := range existing.List {
		if strings.EqualFold(strings.TrimSpace(item.Name), name) {
			result.Failed(c, 400, "配置名称已存在")
			return
		}
	}

	username := firstNonEmptyConfig(strings.TrimSpace(c.GetString("username")), "system")
	config := &model.ConfigAnsible{
		Name:      name,
		Type:      req.Type,
		Content:   req.Content,
		Remark:    req.Remark,
		CreatedBy: username,
		UpdatedBy: username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.dao.Create(config); err != nil {
		result.Failed(c, 500, "创建配置失败: "+err.Error())
		return
	}
	result.Success(c, config)
}

func (s *ConfigAnsibleServiceImpl) Update(c *gin.Context, id uint, req *UpdateConfigRequest) {
	config, err := s.dao.GetByID(id)
	if err != nil {
		result.Failed(c, 404, "配置不存在")
		return
	}

	if req.Name != "" {
		config.Name = strings.TrimSpace(req.Name)
	}
	if req.Type > 0 {
		config.Type = req.Type
	}
	if req.Content != "" {
		config.Content = req.Content
	}
	config.Remark = req.Remark
	config.UpdatedBy = firstNonEmptyConfig(strings.TrimSpace(c.GetString("username")), "system")
	config.UpdatedAt = time.Now()

	if err := s.dao.Update(config); err != nil {
		result.Failed(c, 500, "更新配置失败: "+err.Error())
		return
	}
	result.Success(c, config)
}

func (s *ConfigAnsibleServiceImpl) Delete(c *gin.Context, id uint) {
	if err := s.dao.Delete(id); err != nil {
		result.Failed(c, 500, "删除配置失败: "+err.Error())
		return
	}
	result.Success(c, nil)
}

func (s *ConfigAnsibleServiceImpl) Get(c *gin.Context, id uint) {
	config, err := s.dao.GetByID(id)
	if err != nil {
		result.Failed(c, 404, "配置不存在")
		return
	}
	result.Success(c, config)
}

func (s *ConfigAnsibleServiceImpl) List(c *gin.Context, page, size int, name string, configType int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	data, err := s.dao.List(page, size, strings.TrimSpace(name), configType)
	if err != nil {
		result.Failed(c, 500, "获取配置列表失败: "+err.Error())
		return
	}
	result.Success(c, data)
}

func firstNonEmptyConfig(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
