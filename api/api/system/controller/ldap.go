package controller

import (
	"net/http"

	systemModel "dodevops-api/api/system/model"
	systemService "dodevops-api/api/system/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type LDAPController struct {
	service *systemService.LDAPService
}

func NewLDAPController() *LDAPController {
	return &LDAPController{
		service: systemService.NewLDAPService(),
	}
}

func (c *LDAPController) GetConfig(ctx *gin.Context) {
	c.service.GetConfig(ctx)
}

func (c *LDAPController) UpdateConfig(ctx *gin.Context) {
	var config systemModel.LDAPConfig
	if err := ctx.ShouldBindJSON(&config); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	c.service.UpdateConfig(ctx, config)
}

func (c *LDAPController) TestConfig(ctx *gin.Context) {
	var config systemModel.LDAPConfig
	if err := ctx.ShouldBindJSON(&config); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	c.service.TestConfig(ctx, config)
}
