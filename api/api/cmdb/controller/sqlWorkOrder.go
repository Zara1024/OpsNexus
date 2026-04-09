package controller

import (
	"strconv"
	"strings"

	"dodevops-api/api/cmdb/model"
	"dodevops-api/api/cmdb/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type SQLWorkOrderController struct {
	service *service.SQLWorkOrderService
}

func NewSQLWorkOrderController() *SQLWorkOrderController {
	return &SQLWorkOrderController{
		service: service.NewSQLWorkOrderServiceWithDB(),
	}
}

func (c *SQLWorkOrderController) GetSummary(ctx *gin.Context) {
	c.service.GetSummary(ctx)
}

func (c *SQLWorkOrderController) List(ctx *gin.Context) {
	query := model.CmdbSQLWorkOrderQuery{
		Page:     parseIntSQLWorkOrder(ctx.Query("page"), 1),
		PageSize: parseIntSQLWorkOrder(ctx.Query("pageSize"), 10),
		Status:   parseIntSQLWorkOrder(ctx.Query("status"), 0),
		Keyword:  strings.TrimSpace(ctx.Query("keyword")),
	}
	if databaseID := parseIntSQLWorkOrder(ctx.Query("databaseId"), 0); databaseID > 0 {
		query.DatabaseID = uint(databaseID)
	}
	c.service.List(ctx, query)
}

func (c *SQLWorkOrderController) Detail(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		result.Failed(ctx, 400, "无效的 SQL 工单 ID")
		return
	}
	c.service.Detail(ctx, uint(id))
}

func (c *SQLWorkOrderController) Create(ctx *gin.Context) {
	var req model.CmdbSQLWorkOrderCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	c.service.Create(ctx, req)
}

func (c *SQLWorkOrderController) Approve(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		result.Failed(ctx, 400, "无效的 SQL 工单 ID")
		return
	}
	var req model.CmdbSQLWorkOrderActionRequest
	_ = ctx.ShouldBindJSON(&req)
	c.service.Approve(ctx, uint(id), req)
}

func (c *SQLWorkOrderController) Reject(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		result.Failed(ctx, 400, "无效的 SQL 工单 ID")
		return
	}
	var req model.CmdbSQLWorkOrderActionRequest
	_ = ctx.ShouldBindJSON(&req)
	c.service.Reject(ctx, uint(id), req)
}

func (c *SQLWorkOrderController) Execute(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		result.Failed(ctx, 400, "无效的 SQL 工单 ID")
		return
	}
	var req model.CmdbSQLWorkOrderActionRequest
	_ = ctx.ShouldBindJSON(&req)
	c.service.Execute(ctx, uint(id), req)
}

func parseIntSQLWorkOrder(value string, fallback int) int {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
