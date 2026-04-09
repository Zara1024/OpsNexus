package controller

import (
	"strconv"
	"strings"

	"dodevops-api/api/knowledge/model"
	"dodevops-api/api/knowledge/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type KnowledgeController struct {
	service *service.KnowledgeService
}

func NewKnowledgeController() *KnowledgeController {
	return &KnowledgeController{
		service: service.NewKnowledgeService(),
	}
}

func (c *KnowledgeController) List(ctx *gin.Context) {
	query := model.KnowledgeQuery{
		Page:     parseInt(ctx.Query("page"), 1),
		PageSize: parseInt(ctx.Query("pageSize"), 10),
		Keyword:  strings.TrimSpace(ctx.Query("keyword")),
		Type:     strings.TrimSpace(ctx.Query("type")),
		Category: strings.TrimSpace(ctx.Query("category")),
		Enabled:  parseInt(ctx.Query("enabled"), -1),
	}
	c.service.List(ctx, query)
}

func (c *KnowledgeController) Detail(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		result.Failed(ctx, 400, "无效的知识库ID")
		return
	}
	c.service.Detail(ctx, uint(id))
}

func (c *KnowledgeController) Create(ctx *gin.Context) {
	var item model.KnowledgeBase
	if err := ctx.ShouldBindJSON(&item); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	c.service.Create(ctx, &item)
}

func (c *KnowledgeController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		result.Failed(ctx, 400, "无效的知识库ID")
		return
	}
	var item model.KnowledgeBase
	if err = ctx.ShouldBindJSON(&item); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	c.service.Update(ctx, uint(id), &item)
}

func (c *KnowledgeController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		result.Failed(ctx, 400, "无效的知识库ID")
		return
	}
	c.service.Delete(ctx, uint(id))
}

func (c *KnowledgeController) Bootstrap(ctx *gin.Context) {
	c.service.Bootstrap(ctx)
}

func parseInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
