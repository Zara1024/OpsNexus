package controller

import (
	"net/http"
	"strings"

	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type MonitorAutomationController struct {
	automationService service.MonitorAutomationServiceInterface
}

func NewMonitorAutomationController() *MonitorAutomationController {
	return &MonitorAutomationController{
		automationService: service.NewMonitorAutomationService(),
	}
}

func (c *MonitorAutomationController) GetAutomationOverview(ctx *gin.Context) {
	c.automationService.GetAutomationOverview(ctx)
}

func (c *MonitorAutomationController) GetHostAlertTemplates(ctx *gin.Context) {
	c.automationService.GetHostAlertTemplates(ctx)
}

func (c *MonitorAutomationController) GetHostAlertRuleList(ctx *gin.Context) {
	c.automationService.GetHostAlertRuleList(ctx)
}

func (c *MonitorAutomationController) CreateHostAlertRule(ctx *gin.Context) {
	var req model.MonitorHostAlertRuleUpsertRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}
	c.automationService.CreateHostAlertRule(ctx, req)
}

func (c *MonitorAutomationController) UpdateHostAlertRule(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.MonitorHostAlertRuleUpsertRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}
	c.automationService.UpdateHostAlertRule(ctx, id, req)
}

func (c *MonitorAutomationController) UpdateHostAlertRuleStatus(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.MonitorStatusUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}
	c.automationService.UpdateHostAlertRuleStatus(ctx, id, req)
}

func (c *MonitorAutomationController) DeleteHostAlertRule(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}
	c.automationService.DeleteHostAlertRule(ctx, id)
}

func (c *MonitorAutomationController) ScanHostAlerts(ctx *gin.Context) {
	c.automationService.ScanHostAlerts(ctx)
}

func (c *MonitorAutomationController) GetDBAlertRuleList(ctx *gin.Context) {
	c.automationService.GetDBAlertRuleList(ctx)
}

func (c *MonitorAutomationController) CreateDBAlertRule(ctx *gin.Context) {
	var req model.MonitorDBAlertRuleUpsertRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}
	c.automationService.CreateDBAlertRule(ctx, req)
}

func (c *MonitorAutomationController) UpdateDBAlertRule(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.MonitorDBAlertRuleUpsertRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}
	c.automationService.UpdateDBAlertRule(ctx, id, req)
}

func (c *MonitorAutomationController) UpdateDBAlertRuleStatus(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.MonitorStatusUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}
	c.automationService.UpdateDBAlertRuleStatus(ctx, id, req)
}

func (c *MonitorAutomationController) DeleteDBAlertRule(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}
	c.automationService.DeleteDBAlertRule(ctx, id)
}

func (c *MonitorAutomationController) GetDBHealthSnapshots(ctx *gin.Context) {
	c.automationService.GetDBHealthSnapshots(ctx)
}

func (c *MonitorAutomationController) ScanDBAlerts(ctx *gin.Context) {
	c.automationService.ScanDBAlerts(ctx)
}

func (c *MonitorAutomationController) GetAutomationEventList(ctx *gin.Context) {
	query := model.MonitorAutomationEventQuery{
		ResourceType: strings.TrimSpace(ctx.Query("resourceType")),
		Status:       strings.TrimSpace(ctx.Query("status")),
		Keyword:      strings.TrimSpace(ctx.Query("keyword")),
		PageSize:     parseOptionalInt(ctx.Query("pageSize"), 10),
		PageNum:      parseOptionalInt(ctx.Query("pageNum"), 1),
	}
	c.automationService.GetAutomationEventList(ctx, query)
}

func (c *MonitorAutomationController) GetSSLDomainList(ctx *gin.Context) {
	c.automationService.GetSSLDomainList(ctx)
}

func (c *MonitorAutomationController) CreateSSLDomain(ctx *gin.Context) {
	var req model.MonitorDomainUpsertRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}
	c.automationService.CreateSSLDomain(ctx, req)
}

func (c *MonitorAutomationController) UpdateSSLDomain(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.MonitorDomainUpsertRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}
	c.automationService.UpdateSSLDomain(ctx, id, req)
}

func (c *MonitorAutomationController) UpdateSSLDomainStatus(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.MonitorStatusUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}
	c.automationService.UpdateSSLDomainStatus(ctx, id, req)
}

func (c *MonitorAutomationController) DeleteSSLDomain(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}
	c.automationService.DeleteSSLDomain(ctx, id)
}

func (c *MonitorAutomationController) GetSSLDomainScheduleList(ctx *gin.Context) {
	c.automationService.GetSSLDomainScheduleList(ctx)
}

func (c *MonitorAutomationController) UpdateSSLDomainSchedule(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.MonitorDomainScheduleUpsertRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}
	c.automationService.UpdateSSLDomainSchedule(ctx, id, req)
}

func (c *MonitorAutomationController) ScanSSLDomains(ctx *gin.Context) {
	c.automationService.ScanSSLDomains(ctx)
}

func (c *MonitorAutomationController) GetSSLCertList(ctx *gin.Context) {
	c.automationService.GetSSLCertList(ctx)
}

func (c *MonitorAutomationController) GetSSLCertDeployLogs(ctx *gin.Context) {
	c.automationService.GetSSLCertDeployLogs(ctx)
}

func (c *MonitorAutomationController) DeploySSLCert(ctx *gin.Context) {
	var req model.MonitorSSLDeployRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}
	c.automationService.DeploySSLCert(ctx, req)
}
