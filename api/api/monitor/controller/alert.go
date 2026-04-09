package controller

import (
	"net/http"
	"strconv"
	"strings"

	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/config"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// MonitorAlertController exposes alert center endpoints.
type MonitorAlertController struct {
	alertService service.MonitorAlertServiceInterface
}

// NewMonitorAlertController creates a controller instance for alert center queries.
func NewMonitorAlertController() *MonitorAlertController {
	return &MonitorAlertController{
		alertService: service.NewMonitorAlertService(),
	}
}

// GetAlertSummary returns cross-table alert summary metrics.
func (c *MonitorAlertController) GetAlertSummary(ctx *gin.Context) {
	c.alertService.GetAlertSummary(ctx)
}

// GetIncidentList returns the paginated incident list.
func (c *MonitorAlertController) GetIncidentList(ctx *gin.Context) {
	query := model.MonitorIncidentQuery{
		Keyword:      strings.TrimSpace(ctx.Query("keyword")),
		Status:       parseOptionalInt(ctx.Query("status"), -1),
		Level:        strings.TrimSpace(ctx.Query("level")),
		Namespace:    strings.TrimSpace(ctx.Query("namespace")),
		WorkloadName: strings.TrimSpace(ctx.Query("workloadName")),
		PageSize:     parseOptionalInt(ctx.Query("pageSize"), 10),
		PageNum:      parseOptionalInt(ctx.Query("pageNum"), 1),
	}
	c.alertService.GetIncidentList(ctx, query)
}

// GetWebhookLogList returns the paginated webhook log list.
func (c *MonitorAlertController) GetWebhookLogList(ctx *gin.Context) {
	query := model.MonitorWebhookLogQuery{
		Keyword:  strings.TrimSpace(ctx.Query("keyword")),
		Source:   strings.TrimSpace(ctx.Query("source")),
		Level:    strings.TrimSpace(ctx.Query("level")),
		Status:   strings.TrimSpace(ctx.Query("status")),
		PageSize: parseOptionalInt(ctx.Query("pageSize"), 10),
		PageNum:  parseOptionalInt(ctx.Query("pageNum"), 1),
	}
	c.alertService.GetWebhookLogList(ctx, query)
}

// GetNotifyLogList returns the paginated notify delivery list.
func (c *MonitorAlertController) GetNotifyLogList(ctx *gin.Context) {
	query := model.MonitorNotifyLogQuery{
		Keyword:   strings.TrimSpace(ctx.Query("keyword")),
		Status:    strings.TrimSpace(ctx.Query("status")),
		RobotType: strings.TrimSpace(ctx.Query("robotType")),
		PageSize:  parseOptionalInt(ctx.Query("pageSize"), 10),
		PageNum:   parseOptionalInt(ctx.Query("pageNum"), 1),
	}
	c.alertService.GetNotifyLogList(ctx, query)
}

// GetNotifyRobotList returns all configured notify robots.
func (c *MonitorAlertController) GetNotifyRobotList(ctx *gin.Context) {
	c.alertService.GetNotifyRobotList(ctx)
}

// ReceiveWebhook ingests one alert webhook without JWT.
func (c *MonitorAlertController) ReceiveWebhook(ctx *gin.Context) {
	expectedToken := strings.TrimSpace(config.Config.Monitor.Webhook.Token)
	if expectedToken == "" {
		expectedToken = "webhook-notify-token-2024"
	}

	token := strings.TrimSpace(ctx.GetHeader("X-Webhook-Token"))
	if token == "" {
		token = strings.TrimSpace(ctx.Query("token"))
	}
	if token == "" {
		authHeader := strings.TrimSpace(ctx.GetHeader("Authorization"))
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			token = strings.TrimSpace(authHeader[7:])
		}
	}
	if token != expectedToken {
		result.Failed(ctx, http.StatusUnauthorized, "webhook token 鏍￠獙澶辫触")
		return
	}

	var payload map[string]interface{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "webhook payload 瑙ｆ瀽澶辫触: "+err.Error())
		return
	}

	c.alertService.ReceiveWebhook(ctx, payload)
}

// CreateNotifyRobot creates one notify robot.
func (c *MonitorAlertController) CreateNotifyRobot(ctx *gin.Context) {
	var req model.MonitorNotifyRobotUpsertRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "璇锋眰鍙傛暟閿欒: "+err.Error())
		return
	}
	c.alertService.CreateNotifyRobot(ctx, req)
}

// UpdateNotifyRobot updates one notify robot.
func (c *MonitorAlertController) UpdateNotifyRobot(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.MonitorNotifyRobotUpsertRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "璇锋眰鍙傛暟閿欒: "+err.Error())
		return
	}
	c.alertService.UpdateNotifyRobot(ctx, id, req)
}

// UpdateNotifyRobotStatus updates one notify robot status.
func (c *MonitorAlertController) UpdateNotifyRobotStatus(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.MonitorStatusUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "璇锋眰鍙傛暟閿欒: "+err.Error())
		return
	}
	c.alertService.UpdateNotifyRobotStatus(ctx, id, req)
}

// DeleteNotifyRobot deletes one notify robot.
func (c *MonitorAlertController) DeleteNotifyRobot(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}
	c.alertService.DeleteNotifyRobot(ctx, id)
}

// TestNotifyRobot sends one manual test delivery for the target robot.
func (c *MonitorAlertController) TestNotifyRobot(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.MonitorNotifyRobotTestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}
	c.alertService.TestNotifyRobot(ctx, id, req)
}

// GetAlertSourceList returns all configured alert sources.
func (c *MonitorAlertController) GetAlertSourceList(ctx *gin.Context) {
	c.alertService.GetAlertSourceList(ctx)
}

// CreateAlertSource creates one alert source.
func (c *MonitorAlertController) CreateAlertSource(ctx *gin.Context) {
	var req model.MonitorAlertSourceUpsertRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "璇锋眰鍙傛暟閿欒: "+err.Error())
		return
	}
	c.alertService.CreateAlertSource(ctx, req)
}

// UpdateAlertSource updates one alert source.
func (c *MonitorAlertController) UpdateAlertSource(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.MonitorAlertSourceUpsertRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "璇锋眰鍙傛暟閿欒: "+err.Error())
		return
	}
	c.alertService.UpdateAlertSource(ctx, id, req)
}

// UpdateAlertSourceStatus updates one alert source status.
func (c *MonitorAlertController) UpdateAlertSourceStatus(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.MonitorStatusUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "璇锋眰鍙傛暟閿欒: "+err.Error())
		return
	}
	c.alertService.UpdateAlertSourceStatus(ctx, id, req)
}

// DeleteAlertSource deletes one alert source.
func (c *MonitorAlertController) DeleteAlertSource(ctx *gin.Context) {
	id, ok := parseUintParam(ctx, "id")
	if !ok {
		return
	}
	c.alertService.DeleteAlertSource(ctx, id)
}

// GetAlertManagerStatus returns one AlertManager status snapshot.
func (c *MonitorAlertController) GetAlertManagerStatus(ctx *gin.Context) {
	c.alertService.GetAlertManagerStatus(ctx, parseAlertManagerQuery(ctx))
}

// GetAlertManagerSilenceList returns one AlertManager silence list.
func (c *MonitorAlertController) GetAlertManagerSilenceList(ctx *gin.Context) {
	c.alertService.GetAlertManagerSilenceList(ctx, parseAlertManagerQuery(ctx))
}

// CreateAlertManagerSilence creates one AlertManager silence.
func (c *MonitorAlertController) CreateAlertManagerSilence(ctx *gin.Context) {
	var req model.MonitorAlertManagerSilenceCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, http.StatusBadRequest, "璇锋眰鍙傛暟閿欒: "+err.Error())
		return
	}
	c.alertService.CreateAlertManagerSilence(ctx, req)
}

// DeleteAlertManagerSilence expires one AlertManager silence.
func (c *MonitorAlertController) DeleteAlertManagerSilence(ctx *gin.Context) {
	silenceID := strings.TrimSpace(ctx.Param("silenceId"))
	if silenceID == "" {
		result.Failed(ctx, http.StatusBadRequest, "silenceId is required")
		return
	}
	c.alertService.DeleteAlertManagerSilence(ctx, parseAlertManagerQuery(ctx), silenceID)
}

// GetAlertManagerReceiverList returns one AlertManager receiver list.
func (c *MonitorAlertController) GetAlertManagerReceiverList(ctx *gin.Context) {
	c.alertService.GetAlertManagerReceiverList(ctx, parseAlertManagerQuery(ctx))
}

func parseOptionalInt(value string, fallback int) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseUintParam(ctx *gin.Context, key string) (uint, bool) {
	id, err := strconv.ParseUint(strings.TrimSpace(ctx.Param(key)), 10, 32)
	if err != nil || id == 0 {
		result.Failed(ctx, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return uint(id), true
}

func parseAlertManagerQuery(ctx *gin.Context) model.MonitorAlertManagerQuery {
	return model.MonitorAlertManagerQuery{
		SourceID: parseOptionalUint(ctx.Query("sourceId")),
	}
}

func parseOptionalUint(value string) uint {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}

	parsed, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0
	}
	return uint(parsed)
}
