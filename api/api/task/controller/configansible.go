package controller

import (
	"strconv"

	"dodevops-api/api/task/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type ConfigAnsibleController struct {
	service service.IConfigAnsibleService
}

func NewConfigAnsibleController(service service.IConfigAnsibleService) *ConfigAnsibleController {
	return &ConfigAnsibleController{service: service}
}

func (c *ConfigAnsibleController) Create(ctx *gin.Context) {
	var req service.CreateConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	c.service.Create(ctx, &req)
}

func (c *ConfigAnsibleController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		result.Failed(ctx, 400, "无效的配置ID")
		return
	}

	var req service.UpdateConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	c.service.Update(ctx, uint(id), &req)
}

func (c *ConfigAnsibleController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		result.Failed(ctx, 400, "无效的配置ID")
		return
	}
	c.service.Delete(ctx, uint(id))
}

func (c *ConfigAnsibleController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		result.Failed(ctx, 400, "无效的配置ID")
		return
	}
	c.service.Get(ctx, uint(id))
}

func (c *ConfigAnsibleController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	name := ctx.Query("name")
	configType, _ := strconv.Atoi(ctx.Query("type"))

	c.service.List(ctx, page, size, name, configType)
}
