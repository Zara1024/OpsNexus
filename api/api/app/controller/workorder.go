package controller

import (
	"strconv"

	"dodevops-api/api/app/model"
	"dodevops-api/api/app/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkOrderController struct {
	service *service.WorkOrderService
}

func NewWorkOrderController(db *gorm.DB) *WorkOrderController {
	return &WorkOrderController{
		service: service.NewWorkOrderService(db),
	}
}

func (c *WorkOrderController) GetSummary(ctx *gin.Context) {
	c.service.GetSummary(ctx)
}

func (c *WorkOrderController) GetList(ctx *gin.Context) {
	query := model.WorkOrderQuery{
		Page:     parseIntOrDefault(ctx.Query("page"), 1),
		PageSize: parseIntOrDefault(ctx.Query("pageSize"), 10),
		Type:     ctx.Query("type"),
		Status:   parseIntOrDefault(ctx.Query("status"), 0),
		Keyword:  ctx.Query("keyword"),
	}
	c.service.GetList(ctx, query)
}

func (c *WorkOrderController) GetDetail(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		result.Failed(ctx, 400, "无效的工单ID")
		return
	}
	c.service.GetDetail(ctx, ctx.Param("type"), uint(id))
}

func (c *WorkOrderController) CreateScript(ctx *gin.Context) {
	var req model.ScriptWorkOrderCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	c.service.CreateScriptWorkOrder(ctx, req)
}

func (c *WorkOrderController) ApproveScript(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		result.Failed(ctx, 400, "无效的工单ID")
		return
	}
	var req model.WorkOrderActionRequest
	_ = ctx.ShouldBindJSON(&req)
	c.service.ApproveScriptRelease(ctx, uint(id), req)
}

func (c *WorkOrderController) RejectScript(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		result.Failed(ctx, 400, "无效的工单ID")
		return
	}
	var req model.WorkOrderActionRequest
	_ = ctx.ShouldBindJSON(&req)
	c.service.RejectScriptRelease(ctx, uint(id), req)
}

func (c *WorkOrderController) ExecuteScript(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		result.Failed(ctx, 400, "无效的工单ID")
		return
	}
	var req model.WorkOrderActionRequest
	_ = ctx.ShouldBindJSON(&req)
	c.service.ExecuteScriptRelease(ctx, uint(id), req)
}

func parseIntOrDefault(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
