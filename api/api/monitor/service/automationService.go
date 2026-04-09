package service

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	configService "dodevops-api/api/configcenter/service"
	monitorDao "dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common"
	"dodevops-api/common/config"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

const (
	hostScanTimeout      = 15 * time.Second
	databaseScanTimeout  = 10 * time.Second
	defaultDomainTimeout = 8 * time.Second
	defaultExpireDays    = int64(30)
)

type MonitorAutomationServiceInterface interface {
	GetAutomationOverview(c *gin.Context)
	GetHostAlertTemplates(c *gin.Context)
	GetHostAlertRuleList(c *gin.Context)
	CreateHostAlertRule(c *gin.Context, req model.MonitorHostAlertRuleUpsertRequest)
	UpdateHostAlertRule(c *gin.Context, id uint, req model.MonitorHostAlertRuleUpsertRequest)
	UpdateHostAlertRuleStatus(c *gin.Context, id uint, req model.MonitorStatusUpdateRequest)
	DeleteHostAlertRule(c *gin.Context, id uint)
	ScanHostAlerts(c *gin.Context)
	GetDBAlertRuleList(c *gin.Context)
	CreateDBAlertRule(c *gin.Context, req model.MonitorDBAlertRuleUpsertRequest)
	UpdateDBAlertRule(c *gin.Context, id uint, req model.MonitorDBAlertRuleUpsertRequest)
	UpdateDBAlertRuleStatus(c *gin.Context, id uint, req model.MonitorStatusUpdateRequest)
	DeleteDBAlertRule(c *gin.Context, id uint)
	GetDBHealthSnapshots(c *gin.Context)
	ScanDBAlerts(c *gin.Context)
	GetAutomationEventList(c *gin.Context, query model.MonitorAutomationEventQuery)
	GetSSLDomainList(c *gin.Context)
	CreateSSLDomain(c *gin.Context, req model.MonitorDomainUpsertRequest)
	UpdateSSLDomain(c *gin.Context, id uint, req model.MonitorDomainUpsertRequest)
	UpdateSSLDomainStatus(c *gin.Context, id uint, req model.MonitorStatusUpdateRequest)
	DeleteSSLDomain(c *gin.Context, id uint)
	GetSSLDomainScheduleList(c *gin.Context)
	UpdateSSLDomainSchedule(c *gin.Context, id uint, req model.MonitorDomainScheduleUpsertRequest)
	ScanSSLDomains(c *gin.Context)
	GetSSLCertList(c *gin.Context)
	GetSSLCertDeployLogs(c *gin.Context)
	DeploySSLCert(c *gin.Context, req model.MonitorSSLDeployRequest)
	RunHostAlertScan(ctx context.Context) error
	RunDBAlertScan(ctx context.Context) error
	RunSSLDomainScan(ctx context.Context) error
}

type MonitorAutomationServiceImpl struct {
	hostDao      cmdbDao.CmdbHostDao
	databaseDao  *cmdbDao.CmdbSQLDao
	alertService *MonitorAlertServiceImpl
}

func NewMonitorAutomationService() MonitorAutomationServiceInterface {
	return &MonitorAutomationServiceImpl{
		hostDao:      cmdbDao.NewCmdbHostDao(),
		databaseDao:  cmdbDao.NewCmdbSQLDao(common.GetDB()),
		alertService: &MonitorAlertServiceImpl{},
	}
}

func (s *MonitorAutomationServiceImpl) GetAutomationOverview(c *gin.Context) {
	overview, err := monitorDao.GetMonitorAutomationOverview()
	if err != nil {
		result.Failed(c, 500, "获取监控深化概览失败: "+err.Error())
		return
	}
	s.enrichAutomationOverview(&overview)
	result.Success(c, overview)
}

func (s *MonitorAutomationServiceImpl) enrichAutomationOverview(overview *model.MonitorAutomationOverview) {
	if overview == nil {
		return
	}

	dbConn := common.GetDB()
	dbConn.Table("monitor_notify_robot").Count(&overview.TotalRobotCount)
	dbConn.Table("monitor_notify_robot").Where("status = ?", 1).Count(&overview.EnabledRobotCount)
	dbConn.Table("monitor_webhook_notify_log").Count(&overview.TotalNotifyLogCount)
	dbConn.Table("monitor_webhook_notify_log").Where("status <> ?", "success").Count(&overview.FailedNotifyLogCount)

	overview.RecentEvents = s.loadAutomationWorkbenchRecentEvents()
	overview.RecentActions = s.loadAutomationWorkbenchRecentActions()
	overview.RiskTips = s.buildAutomationRiskTips(*overview)
	overview.RecommendedActions = s.buildAutomationRecommendedActions(*overview)
	overview.AlertCenterPath = "/monitor/alert-center"
	overview.AlertHistoryPath = "/monitor/alert-history"
	overview.AlertNotifyPath = "/monitor/alert-notify"
}

func (s *MonitorAutomationServiceImpl) loadAutomationWorkbenchRecentEvents() []model.MonitorAutomationWorkbenchEvent {
	type row struct {
		Title           string
		Severity        string
		Status          string
		Summary         string
		ResourceType    string
		ResourceName    string
		LastTriggeredAt string
	}

	var rows []row
	_ = common.GetDB().Raw(`
		SELECT
			COALESCE(title, '') AS title,
			COALESCE(severity, '') AS severity,
			COALESCE(status, '') AS status,
			COALESCE(summary, '') AS summary,
			COALESCE(resource_type, '') AS resource_type,
			COALESCE(resource_name, '') AS resource_name,
			COALESCE(DATE_FORMAT(last_triggered_at, '%Y-%m-%d %H:%i:%s'), '') AS last_triggered_at
		FROM monitor_alert_event
		ORDER BY last_triggered_at DESC, id DESC
		LIMIT 5
	`).Scan(&rows).Error

	list := make([]model.MonitorAutomationWorkbenchEvent, 0, len(rows))
	for _, item := range rows {
		list = append(list, model.MonitorAutomationWorkbenchEvent{
			Title:        item.Title,
			Severity:     item.Severity,
			Status:       item.Status,
			Summary:      item.Summary,
			ResourceType: item.ResourceType,
			ResourceName: item.ResourceName,
			OccurredAt:   item.LastTriggeredAt,
		})
	}
	return list
}

func (s *MonitorAutomationServiceImpl) loadAutomationWorkbenchRecentActions() []model.MonitorAutomationWorkbenchAction {
	actions := make([]model.MonitorAutomationWorkbenchAction, 0, 3)

	type notifyLogRow struct {
		RobotName string
		Status    string
		Title     string
		Summary   string
		CreatedAt string
	}
	var notifyLog notifyLogRow
	if err := common.GetDB().Raw(`
		SELECT
			COALESCE(robot_name, '') AS robot_name,
			COALESCE(status, '') AS status,
			COALESCE(alert_title, '') AS title,
			COALESCE(response_body, '') AS summary,
			COALESCE(DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), '') AS created_at
		FROM monitor_webhook_notify_log
		WHERE status <> 'success'
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`).Scan(&notifyLog).Error; err == nil && strings.TrimSpace(notifyLog.CreatedAt) != "" {
		title := "通知失败"
		if strings.TrimSpace(notifyLog.RobotName) != "" {
			title = title + ": " + strings.TrimSpace(notifyLog.RobotName)
		}
		actions = append(actions, model.MonitorAutomationWorkbenchAction{
			Title:   title,
			Status:  "failed",
			Summary: firstNonEmptyAutomation(notifyLog.Title, notifyLog.Summary, notifyLog.CreatedAt),
			Path:    "/monitor/alert-notify",
		})
	}

	type sslLogRow struct {
		Domain     string
		HostName   string
		ErrorMsg   string
		CreateTime string
	}
	var sslLog sslLogRow
	if err := common.GetDB().Raw(`
		SELECT
			COALESCE(domain, '') AS domain,
			COALESCE(host_name, '') AS host_name,
			COALESCE(error_msg, '') AS error_msg,
			COALESCE(DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s'), '') AS create_time
		FROM monitor_ssl_cert_deploy_log
		WHERE status = 3
		ORDER BY create_time DESC, id DESC
		LIMIT 1
	`).Scan(&sslLog).Error; err == nil && strings.TrimSpace(sslLog.CreateTime) != "" {
		actions = append(actions, model.MonitorAutomationWorkbenchAction{
			Title:   fmt.Sprintf("SSL 部署失败: %s", firstNonEmptyAutomation(sslLog.Domain, "-")),
			Status:  "failed",
			Summary: fmt.Sprintf("%s / %s / %s", firstNonEmptyAutomation(sslLog.HostName, "-"), firstNonEmptyAutomation(sslLog.ErrorMsg, "待人工排查"), sslLog.CreateTime),
			Path:    "/monitor/https?tab=ssl",
		})
	}

	type dbRow struct {
		DatabaseName string
		ErrorMsg     string
		UpdateTime   string
	}
	var dbIssue dbRow
	if err := common.GetDB().Raw(`
		SELECT
			COALESCE(database_name, '') AS database_name,
			COALESCE(error_msg, '') AS error_msg,
			COALESCE(DATE_FORMAT(update_time, '%Y-%m-%d %H:%i:%s'), '') AS update_time
		FROM monitor_db_health_snapshot
		WHERE available = 0
		ORDER BY update_time DESC, id DESC
		LIMIT 1
	`).Scan(&dbIssue).Error; err == nil && strings.TrimSpace(dbIssue.UpdateTime) != "" {
		actions = append(actions, model.MonitorAutomationWorkbenchAction{
			Title:   fmt.Sprintf("数据库异常: %s", firstNonEmptyAutomation(dbIssue.DatabaseName, "-")),
			Status:  "warning",
			Summary: fmt.Sprintf("%s / %s", firstNonEmptyAutomation(dbIssue.ErrorMsg, "连接异常"), dbIssue.UpdateTime),
			Path:    "/monitor/https?tab=db",
		})
	}

	return actions
}

func (s *MonitorAutomationServiceImpl) buildAutomationRiskTips(overview model.MonitorAutomationOverview) []string {
	tips := make([]string, 0, 6)
	if overview.OpenEventCount > 0 {
		tips = append(tips, fmt.Sprintf("当前仍有 %d 条开放事件，建议优先在告警中心完成闭环。", overview.OpenEventCount))
	}
	if overview.ExpiringDomainCount > 0 {
		tips = append(tips, fmt.Sprintf("有 %d 个域名证书即将过期，需要尽快复核 SSL 自动化计划。", overview.ExpiringDomainCount))
	}
	if overview.DatabaseTotalCount > overview.DatabaseHealthyCount {
		tips = append(tips, fmt.Sprintf("数据库健康快照异常 %d 个，建议复核连通性与延迟阈值。", overview.DatabaseTotalCount-overview.DatabaseHealthyCount))
	}
	if overview.FailedNotifyLogCount > 0 {
		tips = append(tips, fmt.Sprintf("通知失败日志累计 %d 条，建议检查飞书 / 钉钉 / 企业微信 / 邮件配置。", overview.FailedNotifyLogCount))
	}
	if len(tips) == 0 {
		tips = append(tips, "当前未发现明显高风险项，可继续按主机 / 数据库 / SSL 工作台执行例行巡检。")
	}
	return tips
}

func (s *MonitorAutomationServiceImpl) buildAutomationRecommendedActions(overview model.MonitorAutomationOverview) []model.MonitorAutomationWorkbenchAction {
	actions := make([]model.MonitorAutomationWorkbenchAction, 0, 4)
	actions = append(actions, model.MonitorAutomationWorkbenchAction{
		Title:   "告警中心",
		Status:  "primary",
		Summary: "统一处理开放事件、恢复状态和历史告警。",
		Path:    "/monitor/alert-center",
	})
	actions = append(actions, model.MonitorAutomationWorkbenchAction{
		Title:   "通知配置",
		Status:  "primary",
		Summary: fmt.Sprintf("当前启用机器人 %d 个，通知失败 %d 条。", overview.EnabledRobotCount, overview.FailedNotifyLogCount),
		Path:    "/monitor/alert-notify",
	})
	actions = append(actions, model.MonitorAutomationWorkbenchAction{
		Title:   "自动化工作台",
		Status:  "primary",
		Summary: "联动主机告警、数据库告警和 SSL 自动化规则。",
		Path:    "/monitor/https",
	})
	return actions
}

func (s *MonitorAutomationServiceImpl) GetHostAlertTemplates(c *gin.Context) {
	result.Success(c, builtInHostAlertTemplates())
}

func (s *MonitorAutomationServiceImpl) GetHostAlertRuleList(c *gin.Context) {
	list, err := monitorDao.GetMonitorHostAlertRuleList()
	if err != nil {
		result.Failed(c, 500, "获取主机告警规则失败: "+err.Error())
		return
	}

	hostIDs := make([]uint, 0, len(list))
	for _, item := range list {
		hostIDs = append(hostIDs, item.HostID)
	}
	hosts, _ := monitorDao.GetCmdbHostsByIDs(uniqueUintValues(hostIDs))
	hostMap := make(map[uint]cmdbModel.CmdbHost, len(hosts))
	for _, host := range hosts {
		hostMap[host.ID] = host
	}

	data := make([]gin.H, 0, len(list))
	for _, item := range list {
		host := hostMap[item.HostID]
		data = append(data, gin.H{
			"id":             item.ID,
			"name":           item.Name,
			"hostId":         item.HostID,
			"hostName":       hostDisplayName(host),
			"hostIp":         firstNonEmptyAutomation(host.SSHIP, host.PrivateIP, host.PublicIP),
			"metricKey":      item.MetricKey,
			"operator":       item.Operator,
			"thresholdValue": item.ThresholdValue,
			"severity":       item.Severity,
			"status":         item.Status,
			"notifyRobotIds": decodeUintList(item.NotifyRobotIDs),
			"remark":         item.Remark,
			"lastScanAt":     item.LastScanAt,
			"createTime":     item.CreateTime,
			"updateTime":     item.UpdateTime,
		})
	}
	result.Success(c, data)
}

func (s *MonitorAutomationServiceImpl) CreateHostAlertRule(c *gin.Context, req model.MonitorHostAlertRuleUpsertRequest) {
	entity, err := s.buildHostRuleEntity(0, req)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}
	if err = monitorDao.CreateMonitorHostAlertRule(entity); err != nil {
		result.Failed(c, 500, "创建主机告警规则失败: "+err.Error())
		return
	}
	result.Success(c, entity)
}

func (s *MonitorAutomationServiceImpl) UpdateHostAlertRule(c *gin.Context, id uint, req model.MonitorHostAlertRuleUpsertRequest) {
	if _, err := monitorDao.GetMonitorHostAlertRuleByID(id); err != nil {
		result.Failed(c, 404, "主机告警规则不存在")
		return
	}
	entity, err := s.buildHostRuleEntity(id, req)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}
	if err = monitorDao.UpdateMonitorHostAlertRule(entity); err != nil {
		result.Failed(c, 500, "更新主机告警规则失败: "+err.Error())
		return
	}
	result.Success(c, entity)
}

func (s *MonitorAutomationServiceImpl) UpdateHostAlertRuleStatus(c *gin.Context, id uint, req model.MonitorStatusUpdateRequest) {
	if _, err := monitorDao.GetMonitorHostAlertRuleByID(id); err != nil {
		result.Failed(c, 404, "主机告警规则不存在")
		return
	}
	if err := monitorDao.UpdateMonitorHostAlertRuleStatus(id, req.Status); err != nil {
		result.Failed(c, 500, "更新主机告警规则状态失败: "+err.Error())
		return
	}
	result.Success(c, gin.H{"id": id, "status": req.Status})
}

func (s *MonitorAutomationServiceImpl) DeleteHostAlertRule(c *gin.Context, id uint) {
	if _, err := monitorDao.GetMonitorHostAlertRuleByID(id); err != nil {
		result.Failed(c, 404, "主机告警规则不存在")
		return
	}
	if err := monitorDao.DeleteMonitorHostAlertRule(id); err != nil {
		result.Failed(c, 500, "删除主机告警规则失败: "+err.Error())
		return
	}
	result.Success(c, gin.H{"id": id})
}

func (s *MonitorAutomationServiceImpl) ScanHostAlerts(c *gin.Context) {
	if err := s.RunHostAlertScan(c.Request.Context()); err != nil {
		result.Failed(c, 500, "执行主机告警扫描失败: "+err.Error())
		return
	}
	result.Success(c, gin.H{"message": "主机告警扫描完成"})
}

func (s *MonitorAutomationServiceImpl) GetDBAlertRuleList(c *gin.Context) {
	list, err := monitorDao.GetMonitorDBAlertRuleList()
	if err != nil {
		result.Failed(c, 500, "获取数据库告警规则失败: "+err.Error())
		return
	}

	dbIDs := make([]uint, 0, len(list))
	for _, item := range list {
		dbIDs = append(dbIDs, item.DatabaseID)
	}
	dbMap := s.getDatabaseMap(uniqueUintValues(dbIDs))
	data := make([]gin.H, 0, len(list))
	for _, item := range list {
		target := dbMap[item.DatabaseID]
		data = append(data, gin.H{
			"id":             item.ID,
			"name":           item.Name,
			"databaseId":     item.DatabaseID,
			"databaseName":   target.Name,
			"databaseType":   target.Type,
			"metricKey":      item.MetricKey,
			"operator":       item.Operator,
			"thresholdValue": item.ThresholdValue,
			"severity":       item.Severity,
			"status":         item.Status,
			"notifyRobotIds": decodeUintList(item.NotifyRobotIDs),
			"remark":         item.Remark,
			"lastScanAt":     item.LastScanAt,
			"createTime":     item.CreateTime,
			"updateTime":     item.UpdateTime,
		})
	}
	result.Success(c, data)
}

func (s *MonitorAutomationServiceImpl) CreateDBAlertRule(c *gin.Context, req model.MonitorDBAlertRuleUpsertRequest) {
	entity, err := s.buildDBRuleEntity(0, req)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}
	if err = monitorDao.CreateMonitorDBAlertRule(entity); err != nil {
		result.Failed(c, 500, "创建数据库告警规则失败: "+err.Error())
		return
	}
	result.Success(c, entity)
}

func (s *MonitorAutomationServiceImpl) UpdateDBAlertRule(c *gin.Context, id uint, req model.MonitorDBAlertRuleUpsertRequest) {
	if _, err := monitorDao.GetMonitorDBAlertRuleByID(id); err != nil {
		result.Failed(c, 404, "数据库告警规则不存在")
		return
	}
	entity, err := s.buildDBRuleEntity(id, req)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}
	if err = monitorDao.UpdateMonitorDBAlertRule(entity); err != nil {
		result.Failed(c, 500, "更新数据库告警规则失败: "+err.Error())
		return
	}
	result.Success(c, entity)
}

func (s *MonitorAutomationServiceImpl) UpdateDBAlertRuleStatus(c *gin.Context, id uint, req model.MonitorStatusUpdateRequest) {
	if _, err := monitorDao.GetMonitorDBAlertRuleByID(id); err != nil {
		result.Failed(c, 404, "数据库告警规则不存在")
		return
	}
	if err := monitorDao.UpdateMonitorDBAlertRuleStatus(id, req.Status); err != nil {
		result.Failed(c, 500, "更新数据库告警规则状态失败: "+err.Error())
		return
	}
	result.Success(c, gin.H{"id": id, "status": req.Status})
}

func (s *MonitorAutomationServiceImpl) DeleteDBAlertRule(c *gin.Context, id uint) {
	if _, err := monitorDao.GetMonitorDBAlertRuleByID(id); err != nil {
		result.Failed(c, 404, "数据库告警规则不存在")
		return
	}
	if err := monitorDao.DeleteMonitorDBAlertRule(id); err != nil {
		result.Failed(c, 500, "删除数据库告警规则失败: "+err.Error())
		return
	}
	result.Success(c, gin.H{"id": id})
}

func (s *MonitorAutomationServiceImpl) GetDBHealthSnapshots(c *gin.Context) {
	list, err := monitorDao.GetMonitorDBHealthSnapshotList()
	if err != nil {
		result.Failed(c, 500, "获取数据库健康快照失败: "+err.Error())
		return
	}
	result.Success(c, list)
}

func (s *MonitorAutomationServiceImpl) ScanDBAlerts(c *gin.Context) {
	if err := s.RunDBAlertScan(c.Request.Context()); err != nil {
		result.Failed(c, 500, "执行数据库告警扫描失败: "+err.Error())
		return
	}
	result.Success(c, gin.H{"message": "数据库告警扫描完成"})
}

func (s *MonitorAutomationServiceImpl) GetAutomationEventList(c *gin.Context, query model.MonitorAutomationEventQuery) {
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	if query.PageNum < 1 {
		query.PageNum = 1
	}
	list, total, err := monitorDao.GetMonitorAutomationEventList(query)
	if err != nil {
		result.Failed(c, 500, "获取监控事件失败: "+err.Error())
		return
	}
	result.SuccessWithPage(c, list, total, query.PageNum, query.PageSize)
}

func (s *MonitorAutomationServiceImpl) GetSSLDomainList(c *gin.Context) {
	list, err := monitorDao.GetMonitorDomainList()
	if err != nil {
		result.Failed(c, 500, "获取域名巡检列表失败: "+err.Error())
		return
	}
	result.Success(c, list)
}

func (s *MonitorAutomationServiceImpl) CreateSSLDomain(c *gin.Context, req model.MonitorDomainUpsertRequest) {
	entity, err := buildDomainEntity(0, req)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}
	if err = monitorDao.CreateMonitorDomain(entity); err != nil {
		result.Failed(c, 500, "创建域名巡检对象失败: "+err.Error())
		return
	}
	result.Success(c, entity)
}

func (s *MonitorAutomationServiceImpl) UpdateSSLDomain(c *gin.Context, id uint, req model.MonitorDomainUpsertRequest) {
	if _, err := monitorDao.GetMonitorDomainByID(id); err != nil {
		result.Failed(c, 404, "域名巡检对象不存在")
		return
	}
	entity, err := buildDomainEntity(id, req)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}
	if err = monitorDao.UpdateMonitorDomain(entity); err != nil {
		result.Failed(c, 500, "更新域名巡检对象失败: "+err.Error())
		return
	}
	result.Success(c, entity)
}

func (s *MonitorAutomationServiceImpl) UpdateSSLDomainStatus(c *gin.Context, id uint, req model.MonitorStatusUpdateRequest) {
	if _, err := monitorDao.GetMonitorDomainByID(id); err != nil {
		result.Failed(c, 404, "域名巡检对象不存在")
		return
	}
	if err := monitorDao.UpdateMonitorDomainStatus(id, int64(req.Status)); err != nil {
		result.Failed(c, 500, "更新域名巡检状态失败: "+err.Error())
		return
	}
	result.Success(c, gin.H{"id": id, "status": req.Status})
}

func (s *MonitorAutomationServiceImpl) DeleteSSLDomain(c *gin.Context, id uint) {
	if _, err := monitorDao.GetMonitorDomainByID(id); err != nil {
		result.Failed(c, 404, "域名巡检对象不存在")
		return
	}
	if err := monitorDao.DeleteMonitorDomain(id); err != nil {
		result.Failed(c, 500, "删除域名巡检对象失败: "+err.Error())
		return
	}
	result.Success(c, gin.H{"id": id})
}

func (s *MonitorAutomationServiceImpl) GetSSLDomainScheduleList(c *gin.Context) {
	list, err := monitorDao.GetMonitorDomainScheduleList()
	if err != nil {
		result.Failed(c, 500, "获取域名巡检调度配置失败: "+err.Error())
		return
	}
	result.Success(c, list)
}

func (s *MonitorAutomationServiceImpl) UpdateSSLDomainSchedule(c *gin.Context, id uint, req model.MonitorDomainScheduleUpsertRequest) {
	item, err := monitorDao.GetMonitorDomainScheduleByID(id)
	if err != nil {
		result.Failed(c, 404, "域名巡检调度配置不存在")
		return
	}
	if err = applyDomainScheduleRequest(item, req); err != nil {
		result.Failed(c, 400, err.Error())
		return
	}
	if err = monitorDao.UpdateMonitorDomainSchedule(item); err != nil {
		result.Failed(c, 500, "更新域名巡检调度配置失败: "+err.Error())
		return
	}
	result.Success(c, item)
}

func (s *MonitorAutomationServiceImpl) ScanSSLDomains(c *gin.Context) {
	if err := s.RunSSLDomainScan(c.Request.Context()); err != nil {
		result.Failed(c, 500, "执行 SSL 域名扫描失败: "+err.Error())
		return
	}
	result.Success(c, gin.H{"message": "SSL 域名扫描完成"})
}

func (s *MonitorAutomationServiceImpl) GetSSLCertList(c *gin.Context) {
	list, err := monitorDao.GetMonitorSSLCertList()
	if err != nil {
		result.Failed(c, 500, "获取 SSL 证书列表失败: "+err.Error())
		return
	}
	result.Success(c, list)
}

func (s *MonitorAutomationServiceImpl) GetSSLCertDeployLogs(c *gin.Context) {
	list, err := monitorDao.GetMonitorSSLCertDeployLogList(50)
	if err != nil {
		result.Failed(c, 500, "获取 SSL 部署日志失败: "+err.Error())
		return
	}
	result.Success(c, list)
}

func (s *MonitorAutomationServiceImpl) DeploySSLCert(c *gin.Context, req model.MonitorSSLDeployRequest) {
	logEntry, err := s.deploySSLCertificate(c.Request.Context(), req)
	if err != nil {
		result.Failed(c, 500, "执行 SSL 证书部署失败: "+err.Error())
		return
	}
	result.Success(c, logEntry)
}

func (s *MonitorAutomationServiceImpl) RunHostAlertScan(ctx context.Context) error {
	rules, err := monitorDao.GetEnabledMonitorHostAlertRules()
	if err != nil {
		return err
	}
	if len(rules) == 0 {
		return nil
	}

	hostIDs := make([]uint, 0, len(rules))
	for _, rule := range rules {
		hostIDs = append(hostIDs, rule.HostID)
	}
	hosts, err := monitorDao.GetCmdbHostsByIDs(uniqueUintValues(hostIDs))
	if err != nil {
		return err
	}
	hostMap := make(map[uint]cmdbModel.CmdbHost, len(hosts))
	for _, host := range hosts {
		hostMap[host.ID] = host
	}

	scanAt := time.Now()
	for _, rule := range rules {
		host, ok := hostMap[rule.HostID]
		if !ok {
			continue
		}

		value, summary, detail, err := s.evaluateHostRule(ctx, host, rule)
		_ = monitorDao.UpdateMonitorHostAlertRuleLastScan(rule.ID, scanAt)
		triggered := err != nil
		if !triggered {
			triggered = compareThreshold(value, rule.Operator, rule.ThresholdValue)
		}

		labels := map[string]interface{}{
			"hostId":     host.ID,
			"hostName":   hostDisplayName(host),
			"hostIp":     firstNonEmptyAutomation(host.SSHIP, host.PrivateIP, host.PublicIP),
			"metricKey":  rule.MetricKey,
			"ruleName":   rule.Name,
			"sourceType": "host",
		}
		if triggered {
			if err != nil {
				detail = firstNonEmptyAutomation(detail, err.Error())
				summary = firstNonEmptyAutomation(summary, fmt.Sprintf("%s 扫描失败", rule.Name))
			}
			_, eventErr := s.triggerAutomationEvent(automationTriggerInput{
				RuleCategory:   model.MonitorRuleCategoryHost,
				RuleID:         rule.ID,
				ResourceType:   "host",
				ResourceID:     host.ID,
				ResourceName:   hostDisplayName(host),
				EventKey:       rule.MetricKey,
				Title:          fmt.Sprintf("主机告警: %s - %s", hostDisplayName(host), hostMetricLabel(rule.MetricKey)),
				Summary:        summary,
				Detail:         detail,
				Severity:       normalizeSeverity(rule.Severity),
				Operator:       rule.Operator,
				ThresholdValue: rule.ThresholdValue,
				CurrentValue:   value,
				NotifyRobotIDs: decodeUintList(rule.NotifyRobotIDs),
				Labels:         labels,
			})
			if eventErr != nil {
				return eventErr
			}
			continue
		}

		if err = s.resolveAutomationEvent(automationResolveInput{
			Fingerprint:    buildAutomationFingerprint("host", host.ID, rule.MetricKey),
			RecoveryValue:  value,
			Solution:       fmt.Sprintf("%s 已恢复到安全阈值内", hostMetricLabel(rule.MetricKey)),
			RecoveryRemark: detail,
			Labels:         labels,
		}); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	return nil
}

func (s *MonitorAutomationServiceImpl) RunDBAlertScan(ctx context.Context) error {
	rules, err := monitorDao.GetEnabledMonitorDBAlertRules()
	if err != nil {
		return err
	}
	if len(rules) == 0 {
		return nil
	}

	dbIDs := make([]uint, 0, len(rules))
	for _, rule := range rules {
		dbIDs = append(dbIDs, rule.DatabaseID)
	}
	dbMap := s.getDatabaseMap(uniqueUintValues(dbIDs))
	scanAt := time.Now()

	for _, rule := range rules {
		target, ok := dbMap[rule.DatabaseID]
		if !ok {
			continue
		}
		account, decrypted, err := s.resolveAccount(target.AccountID)
		if err != nil {
			return err
		}

		available, latencyMs, errMsg := checkDatabaseHealth(ctx, target, account, decrypted)
		_ = monitorDao.UpsertMonitorDBHealthSnapshot(&model.MonitorDBHealthSnapshotEntity{
			DatabaseID:   target.ID,
			DatabaseName: target.Name,
			DatabaseType: target.Type,
			Host:         account.Host,
			Port:         account.Port,
			Available:    boolToInt(available),
			LatencyMs:    latencyMs,
			ErrorMsg:     errMsg,
			LastCheckAt:  scanAt,
		})
		_ = monitorDao.UpdateMonitorDBAlertRuleLastScan(rule.ID, scanAt)

		value := float64(latencyMs)
		if rule.MetricKey == "connectivity" {
			if available {
				value = 0
			} else {
				value = 1
			}
		}
		triggered := !available && rule.MetricKey == "connectivity"
		if !triggered {
			triggered = compareThreshold(value, rule.Operator, rule.ThresholdValue)
		}
		summary := buildDatabaseAlertSummary(target, rule, available, latencyMs, errMsg)
		labels := map[string]interface{}{
			"databaseId":   target.ID,
			"databaseName": target.Name,
			"databaseType": databaseTypeText(target.Type),
			"metricKey":    rule.MetricKey,
			"ruleName":     rule.Name,
			"sourceType":   "database",
			"host":         account.Host,
			"port":         account.Port,
		}
		if triggered {
			_, eventErr := s.triggerAutomationEvent(automationTriggerInput{
				RuleCategory:   model.MonitorRuleCategoryDB,
				RuleID:         rule.ID,
				ResourceType:   "database",
				ResourceID:     target.ID,
				ResourceName:   target.Name,
				EventKey:       rule.MetricKey,
				Title:          fmt.Sprintf("数据库告警: %s - %s", target.Name, databaseMetricLabel(rule.MetricKey)),
				Summary:        summary,
				Detail:         summary,
				Severity:       normalizeSeverity(rule.Severity),
				Operator:       rule.Operator,
				ThresholdValue: rule.ThresholdValue,
				CurrentValue:   value,
				NotifyRobotIDs: decodeUintList(rule.NotifyRobotIDs),
				Labels:         labels,
			})
			if eventErr != nil {
				return eventErr
			}
			continue
		}

		if err = s.resolveAutomationEvent(automationResolveInput{
			Fingerprint:    buildAutomationFingerprint("database", target.ID, rule.MetricKey),
			RecoveryValue:  value,
			Solution:       fmt.Sprintf("%s 已恢复正常", databaseMetricLabel(rule.MetricKey)),
			RecoveryRemark: summary,
			Labels:         labels,
		}); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	return nil
}

func (s *MonitorAutomationServiceImpl) RunSSLDomainScan(ctx context.Context) error {
	domains, err := monitorDao.GetEnabledMonitorDomains()
	if err != nil {
		return err
	}
	if len(domains) == 0 {
		return nil
	}

	schedules, err := monitorDao.GetEnabledMonitorDomainSchedules()
	if err != nil {
		return err
	}
	schedule := firstEnabledSchedule(schedules)
	timeout := time.Duration(defaultDomainTimeout)
	expireDays := defaultExpireDays
	notifyRobotIDs := []uint{}
	if schedule != nil {
		if schedule.ScanTimeoutMs > 0 {
			timeout = time.Duration(schedule.ScanTimeoutMs) * time.Millisecond
		}
		if schedule.ExpireAlertDays > 0 {
			expireDays = schedule.ExpireAlertDays
		}
		if schedule.NotifyEnabled && schedule.NotifyRobotID > 0 {
			notifyRobotIDs = []uint{schedule.NotifyRobotID}
		}
		now := time.Now()
		var nextRunAt *time.Time
		if strings.TrimSpace(schedule.CronExpr) != "" {
			if next, calcErr := calculateNextCronRun(schedule.CronExpr, now); calcErr == nil {
				nextRunAt = &next
			}
		}
		_ = monitorDao.UpdateMonitorDomainScheduleRuntime(schedule.ID, "running", &now, nextRunAt)
	}

	httpClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		},
	}

	for _, domain := range domains {
		checked, checkErr := checkDomainHealth(ctx, domain.Domain, timeout, httpClient)
		updateEntity := domain
		now := time.Now()
		updateEntity.LastCheckAt = &now
		if checkErr != nil {
			updateEntity.IsAlive = 0
			updateEntity.StatusCode = 0
			updateEntity.ResponseTime = 0
			updateEntity.ErrorMsg = checkErr.Error()
			updateEntity.SSLExpireAt = nil
			updateEntity.SSLDaysLeft = 0
			updateEntity.SSLIssuer = ""
		} else {
			updateEntity.IsAlive = 1
			updateEntity.StatusCode = int64(checked.StatusCode)
			updateEntity.ResponseTime = checked.ResponseTime
			updateEntity.ErrorMsg = ""
			updateEntity.SSLExpireAt = checked.ExpireAt
			updateEntity.SSLDaysLeft = checked.DaysLeft
			updateEntity.SSLIssuer = checked.Issuer
		}
		if err = monitorDao.UpdateMonitorDomainCheckResult(&updateEntity); err != nil {
			return err
		}

		labels := map[string]interface{}{
			"domainId":    domain.ID,
			"domain":      domain.Domain,
			"sourceType":  "ssl",
			"responseMs":  updateEntity.ResponseTime,
			"statusCode":  updateEntity.StatusCode,
			"sslDaysLeft": updateEntity.SSLDaysLeft,
			"issuer":      updateEntity.SSLIssuer,
		}
		if updateEntity.IsAlive == 0 {
			_, eventErr := s.triggerAutomationEvent(automationTriggerInput{
				RuleCategory:   model.MonitorRuleCategorySSL,
				ResourceType:   "domain",
				ResourceID:     domain.ID,
				ResourceName:   domain.Domain,
				EventKey:       "domain_unreachable",
				Title:          fmt.Sprintf("SSL/域名告警: %s 不可用", domain.Domain),
				Summary:        fmt.Sprintf("域名 %s 巡检失败: %s", domain.Domain, updateEntity.ErrorMsg),
				Detail:         updateEntity.ErrorMsg,
				Severity:       "P2",
				Operator:       ">=",
				ThresholdValue: 1,
				CurrentValue:   1,
				NotifyRobotIDs: notifyRobotIDs,
				Labels:         labels,
			})
			if eventErr != nil {
				return eventErr
			}
		} else if err = s.resolveAutomationEvent(automationResolveInput{
			Fingerprint:    buildAutomationFingerprint("domain", domain.ID, "domain_unreachable"),
			RecoveryValue:  0,
			Solution:       fmt.Sprintf("%s 已恢复可访问", domain.Domain),
			RecoveryRemark: fmt.Sprintf("域名状态码=%d, 响应时间=%dms", updateEntity.StatusCode, updateEntity.ResponseTime),
			Labels:         labels,
		}); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if updateEntity.IsAlive == 1 && updateEntity.SSLDaysLeft > 0 && updateEntity.SSLDaysLeft <= expireDays {
			_, eventErr := s.triggerAutomationEvent(automationTriggerInput{
				RuleCategory:   model.MonitorRuleCategorySSL,
				ResourceType:   "domain",
				ResourceID:     domain.ID,
				ResourceName:   domain.Domain,
				EventKey:       "ssl_expiry",
				Title:          fmt.Sprintf("SSL 告警: %s 证书即将过期", domain.Domain),
				Summary:        fmt.Sprintf("%s SSL 证书剩余 %d 天，阈值 %d 天", domain.Domain, updateEntity.SSLDaysLeft, expireDays),
				Detail:         fmt.Sprintf("颁发者: %s, 到期时间: %v", updateEntity.SSLIssuer, updateEntity.SSLExpireAt),
				Severity:       "P2",
				Operator:       "<=",
				ThresholdValue: float64(expireDays),
				CurrentValue:   float64(updateEntity.SSLDaysLeft),
				NotifyRobotIDs: notifyRobotIDs,
				Labels:         labels,
			})
			if eventErr != nil {
				return eventErr
			}
		} else if err = s.resolveAutomationEvent(automationResolveInput{
			Fingerprint:    buildAutomationFingerprint("domain", domain.ID, "ssl_expiry"),
			RecoveryValue:  float64(updateEntity.SSLDaysLeft),
			Solution:       fmt.Sprintf("%s SSL 证书剩余天数恢复到阈值之外", domain.Domain),
			RecoveryRemark: fmt.Sprintf("当前剩余 %d 天", updateEntity.SSLDaysLeft),
			Labels:         labels,
		}); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	return nil
}

func (s *MonitorAutomationServiceImpl) buildHostRuleEntity(id uint, req model.MonitorHostAlertRuleUpsertRequest) (*model.MonitorHostAlertRuleEntity, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.MetricKey = strings.TrimSpace(req.MetricKey)
	req.Operator = strings.TrimSpace(req.Operator)
	req.Severity = normalizeSeverity(req.Severity)
	req.Remark = strings.TrimSpace(req.Remark)
	if req.Name == "" {
		return nil, errors.New("规则名称不能为空")
	}
	if req.HostID == 0 {
		return nil, errors.New("主机不能为空")
	}
	if _, err := s.hostDao.GetCmdbHostById(req.HostID); err != nil {
		return nil, errors.New("主机不存在")
	}
	if _, ok := hostMetricQueryBuilders()[req.MetricKey]; !ok {
		return nil, errors.New("主机告警指标不支持")
	}
	if !isSupportedComparison(req.Operator) {
		return nil, errors.New("比较符仅支持 > >= < <= == !=")
	}
	if req.Status != 0 && req.Status != 1 {
		req.Status = 1
	}
	return &model.MonitorHostAlertRuleEntity{
		ID:             id,
		Name:           req.Name,
		HostID:         req.HostID,
		MetricKey:      req.MetricKey,
		Operator:       req.Operator,
		ThresholdValue: req.ThresholdValue,
		Severity:       req.Severity,
		Status:         req.Status,
		NotifyRobotIDs: monitorDao.EncodeUintList(req.NotifyRobotIDs),
		Remark:         req.Remark,
	}, nil
}

func (s *MonitorAutomationServiceImpl) buildDBRuleEntity(id uint, req model.MonitorDBAlertRuleUpsertRequest) (*model.MonitorDBAlertRuleEntity, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.MetricKey = strings.TrimSpace(req.MetricKey)
	req.Operator = strings.TrimSpace(req.Operator)
	req.Severity = normalizeSeverity(req.Severity)
	req.Remark = strings.TrimSpace(req.Remark)
	if req.Name == "" {
		return nil, errors.New("规则名称不能为空")
	}
	if req.DatabaseID == 0 {
		return nil, errors.New("数据库不能为空")
	}
	if _, err := s.databaseDao.GetByID(req.DatabaseID); err != nil {
		return nil, errors.New("数据库资源不存在")
	}
	if req.MetricKey != "connectivity" && req.MetricKey != "latency_ms" {
		return nil, errors.New("数据库告警指标仅支持 connectivity / latency_ms")
	}
	if !isSupportedComparison(req.Operator) {
		return nil, errors.New("比较符仅支持 > >= < <= == !=")
	}
	if req.Status != 0 && req.Status != 1 {
		req.Status = 1
	}
	return &model.MonitorDBAlertRuleEntity{
		ID:             id,
		Name:           req.Name,
		DatabaseID:     req.DatabaseID,
		MetricKey:      req.MetricKey,
		Operator:       req.Operator,
		ThresholdValue: req.ThresholdValue,
		Severity:       req.Severity,
		Status:         req.Status,
		NotifyRobotIDs: monitorDao.EncodeUintList(req.NotifyRobotIDs),
		Remark:         req.Remark,
	}, nil
}

type automationTriggerInput struct {
	RuleCategory   string
	RuleID         uint
	ResourceType   string
	ResourceID     uint
	ResourceName   string
	EventKey       string
	Title          string
	Summary        string
	Detail         string
	Severity       string
	Operator       string
	ThresholdValue float64
	CurrentValue   float64
	NotifyRobotIDs []uint
	Labels         map[string]interface{}
}

type automationResolveInput struct {
	Fingerprint    string
	RecoveryValue  float64
	Solution       string
	RecoveryRemark string
	Labels         map[string]interface{}
}

type automationDispatchResult struct {
	WebhookLogID uint
	Status       string
	NotifyCount  int64
	SuccessCount int64
	FailedCount  int64
}

func (s *MonitorAutomationServiceImpl) triggerAutomationEvent(input automationTriggerInput) (*model.MonitorAlertEventEntity, error) {
	fingerprint := buildAutomationFingerprint(input.ResourceType, input.ResourceID, input.EventKey)
	existing, err := monitorDao.GetLatestOpenMonitorAutomationEvent(fingerprint)
	if err == nil && existing != nil {
		existing.CurrentValue = normalizeFloat(input.CurrentValue)
		existing.LastTriggeredAt = time.Now()
		existing.DedupCount++
		existing.Summary = input.Summary
		existing.Detail = input.Detail
		existing.Labels = encodeAutomationJSON(input.Labels)
		return existing, monitorDao.UpdateMonitorAutomationEvent(existing)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	now := time.Now()
	incident := &model.MonitorIncidentEntity{
		AlertTime:     now,
		BusinessLine:  input.ResourceType,
		Frequency:     "频繁",
		AlertDesc:     input.Summary,
		AlertLevel:    input.Severity,
		Department:    "运维部",
		IncidentCause: input.Detail,
		Status:        1,
		CreateTime:    now,
		UpdateTime:    now,
		Remark:        input.Detail,
	}
	if err = monitorDao.CreateMonitorIncident(incident); err != nil {
		return nil, err
	}

	dispatchResult, err := s.dispatchInternalAlert(model.MonitorWebhookReceiveRequest{
		Source:         "opsnexus-automation",
		Title:          input.Title,
		Content:        input.Summary,
		Level:          strings.ToLower(input.Severity),
		Tags:           input.Labels,
		Extra:          map[string]interface{}{"detail": input.Detail, "severity": input.Severity},
		NotifyRobotIDs: input.NotifyRobotIDs,
	})
	if err != nil {
		return nil, err
	}

	event := &model.MonitorAlertEventEntity{
		RuleCategory:     input.RuleCategory,
		RuleID:           input.RuleID,
		ResourceType:     input.ResourceType,
		ResourceID:       input.ResourceID,
		ResourceName:     input.ResourceName,
		Fingerprint:      fingerprint,
		EventKey:         input.EventKey,
		Title:            input.Title,
		Summary:          input.Summary,
		Detail:           input.Detail,
		Severity:         input.Severity,
		Status:           model.MonitorEventStatusOpen,
		Operator:         input.Operator,
		ThresholdValue:   normalizeFloat(input.ThresholdValue),
		CurrentValue:     normalizeFloat(input.CurrentValue),
		DedupCount:       1,
		NotifyRobotIDs:   monitorDao.EncodeUintList(input.NotifyRobotIDs),
		IncidentID:       incident.ID,
		WebhookLogID:     dispatchResult.WebhookLogID,
		Labels:           encodeAutomationJSON(input.Labels),
		FirstTriggeredAt: now,
		LastTriggeredAt:  now,
		CreateTime:       now,
		UpdateTime:       now,
	}
	if err = monitorDao.CreateMonitorAutomationEvent(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *MonitorAutomationServiceImpl) resolveAutomationEvent(input automationResolveInput) error {
	event, err := monitorDao.GetLatestOpenMonitorAutomationEvent(input.Fingerprint)
	if err != nil {
		return err
	}
	solution := firstNonEmptyAutomation(strings.TrimSpace(input.Solution), "告警已恢复")
	if event.IncidentID > 0 {
		if err = monitorDao.UpdateMonitorIncidentResolved(event.IncidentID, solution, strings.TrimSpace(input.RecoveryRemark)); err != nil {
			return err
		}
	}
	dispatchResult, dispatchErr := s.dispatchInternalAlert(model.MonitorWebhookReceiveRequest{
		Source: "opsnexus-automation-recovery",
		Title:  "[恢复] " + event.Title,
		Content: firstNonEmptyAutomation(
			strings.TrimSpace(input.RecoveryRemark),
			fmt.Sprintf("%s 已恢复", event.ResourceName),
		),
		Level:          "info",
		Tags:           input.Labels,
		Extra:          map[string]interface{}{"recovered": true},
		NotifyRobotIDs: decodeUintList(event.NotifyRobotIDs),
	})
	if dispatchErr != nil {
		return dispatchErr
	}

	now := time.Now()
	event.Status = model.MonitorEventStatusResolved
	event.RecoveredAt = &now
	event.RecoveryValue = normalizeFloat(input.RecoveryValue)
	event.RecoveryLogID = dispatchResult.WebhookLogID
	event.Detail = firstNonEmptyAutomation(strings.TrimSpace(input.RecoveryRemark), event.Detail)
	event.Summary = solution
	event.Labels = encodeAutomationJSON(input.Labels)
	return monitorDao.UpdateMonitorAutomationEvent(event)
}

func (s *MonitorAutomationServiceImpl) dispatchInternalAlert(req model.MonitorWebhookReceiveRequest) (*automationDispatchResult, error) {
	logEntry := &model.MonitorWebhookLogEntity{
		Source:         firstNonEmptyAutomation(req.Source, "opsnexus-automation"),
		Title:          req.Title,
		Content:        req.Content,
		Level:          req.Level,
		Tags:           encodeJSONValue(req.Tags),
		Extra:          encodeJSONValue(req.Extra),
		NotifyRobotIDs: "[]",
		Status:         "success",
		ErrorMsg:       "",
		NotifyCount:    0,
		SuccessCount:   0,
		FailedCount:    0,
	}
	if err := monitorDao.CreateMonitorWebhookLog(logEntry); err != nil {
		return nil, err
	}
	robots, err := monitorDao.GetEnabledMonitorNotifyRobots(req.NotifyRobotIDs)
	if err != nil {
		return nil, err
	}

	notifyIDs := make([]uint, 0, len(robots))
	errorMessages := make([]string, 0)
	successCount := int64(0)
	failedCount := int64(0)
	for _, robot := range robots {
		notifyIDs = append(notifyIDs, robot.ID)
		sendErr := sendAlertToRobot(&robot, req)
		status := "success"
		errorMsg := ""
		if sendErr != nil {
			status = "failed"
			errorMsg = truncateText(sendErr.Error(), 1000)
			failedCount++
			errorMessages = append(errorMessages, errorMsg)
		} else {
			successCount++
		}
		_ = monitorDao.CreateMonitorWebhookNotifyLog(&model.MonitorWebhookNotifyLogEntity{
			WebhookLogID: logEntry.ID,
			RobotID:      robot.ID,
			RobotName:    robot.Name,
			RobotType:    robot.Type,
			Status:       status,
			ErrorMsg:     errorMsg,
		})
	}

	notifyCount := int64(len(robots))
	finalStatus := calculateDispatchStatus(req, len(robots), successCount, failedCount)
	if err = monitorDao.UpdateMonitorWebhookLogResult(
		logEntry.ID,
		finalStatus,
		truncateText(strings.Join(errorMessages, " | "), 1000),
		encodeJSONValue(notifyIDs),
		notifyCount,
		successCount,
		failedCount,
	); err != nil {
		return nil, err
	}
	return &automationDispatchResult{
		WebhookLogID: logEntry.ID,
		Status:       finalStatus,
		NotifyCount:  notifyCount,
		SuccessCount: successCount,
		FailedCount:  failedCount,
	}, nil
}

func (s *MonitorAutomationServiceImpl) evaluateHostRule(ctx context.Context, host cmdbModel.CmdbHost, rule model.MonitorHostAlertRuleEntity) (float64, string, string, error) {
	instance := firstNonEmptyAutomation(host.Name, host.HostName, host.SSHIP)
	queryBuilder, ok := hostMetricQueryBuilders()[rule.MetricKey]
	if !ok {
		return 0, "", "", errors.New("unsupported host metric")
	}
	query := queryBuilder(instance, rule.ThresholdValue)
	value, err := queryPrometheusInstant(ctx, query)
	if err != nil {
		return 0, "", "", err
	}
	summary := fmt.Sprintf("%s %s 当前值 %.2f，阈值 %s %.2f",
		hostDisplayName(host),
		hostMetricLabel(rule.MetricKey),
		value,
		rule.Operator,
		rule.ThresholdValue,
	)
	if rule.MetricKey == "offline" {
		if value >= 1 {
			summary = fmt.Sprintf("%s 在最近 %.0f 分钟内没有采集到监控指标", hostDisplayName(host), rule.ThresholdValue)
		} else {
			summary = fmt.Sprintf("%s 监控指标已恢复正常采集", hostDisplayName(host))
		}
	}
	return value, summary, summary, nil
}

func (s *MonitorAutomationServiceImpl) getDatabaseMap(ids []uint) map[uint]cmdbModel.CmdbSQL {
	resultMap := make(map[uint]cmdbModel.CmdbSQL, len(ids))
	for _, id := range ids {
		item, err := s.databaseDao.GetByID(id)
		if err != nil || item == nil {
			continue
		}
		resultMap[id] = *item
	}
	return resultMap
}

func (s *MonitorAutomationServiceImpl) resolveAccount(accountID uint) (*configAccountView, string, error) {
	account, err := configService.NewAccountAuthService().GetByID(accountID)
	if err != nil {
		return nil, "", fmt.Errorf("获取数据库账号失败: %v", err)
	}
	password, err := configService.NewAccountAuthService().DecryptPassword(accountID)
	if err != nil {
		return nil, "", fmt.Errorf("解密数据库账号失败: %v", err)
	}
	return &configAccountView{
		ID:   account.ID,
		Host: account.Host,
		Port: account.Port,
		Name: account.Name,
		Type: account.Type,
	}, password, nil
}

type configAccountView struct {
	ID   uint
	Host string
	Port int
	Name string
	Type int
}

func builtInHostAlertTemplates() []model.MonitorAlertRuleTemplate {
	return []model.MonitorAlertRuleTemplate{
		{Name: "CPU 高使用率", MetricKey: "cpu_usage", Operator: ">=", ThresholdValue: 90, Severity: "P2", Description: "适用于主机 CPU 长时间打满预警"},
		{Name: "内存高使用率", MetricKey: "memory_usage", Operator: ">=", ThresholdValue: 85, Severity: "P2", Description: "适用于 Linux/Windows 统一内存占用预警"},
		{Name: "磁盘高使用率", MetricKey: "disk_usage", Operator: ">=", ThresholdValue: 85, Severity: "P2", Description: "使用实例维度最大磁盘使用率，兼容 Linux 与 Windows Exporter"},
		{Name: "主机离线", MetricKey: "offline", Operator: ">=", ThresholdValue: 5, Severity: "P1", Description: "最近 N 分钟无监控指标即触发离线告警"},
		{Name: "系统负载过高", MetricKey: "load1", Operator: ">=", ThresholdValue: 8, Severity: "P3", Description: "适用于高负载场景的快速预警"},
	}
}

func buildDomainEntity(id uint, req model.MonitorDomainUpsertRequest) (*model.MonitorDomainEntity, error) {
	req.Domain = normalizeDomain(req.Domain)
	if req.Domain == "" {
		return nil, errors.New("域名不能为空")
	}
	if req.Status != 0 && req.Status != 1 {
		req.Status = 1
	}
	return &model.MonitorDomainEntity{
		ID:     id,
		Domain: req.Domain,
		Tags:   strings.TrimSpace(req.Tags),
		Remark: strings.TrimSpace(req.Remark),
		Status: req.Status,
	}, nil
}

func applyDomainScheduleRequest(item *model.MonitorDomainScheduleEntity, req model.MonitorDomainScheduleUpsertRequest) error {
	item.Enabled = req.Enabled
	item.CronExpr = strings.TrimSpace(req.CronExpr)
	item.NotifyEnabled = req.NotifyEnabled
	item.NotifyRobotID = req.NotifyRobotID
	item.ExpireAlertDays = req.ExpireAlertDays
	item.ScanTimeoutMs = req.ScanTimeoutMs
	item.AutoRenewEnabled = req.AutoRenewEnabled
	item.AutoDeployEnabled = req.AutoDeployEnabled
	item.DeployHostID = req.DeployHostID
	item.DeployPath = strings.TrimSpace(req.DeployPath)
	item.ReloadCommand = strings.TrimSpace(req.ReloadCommand)
	if item.ExpireAlertDays <= 0 {
		item.ExpireAlertDays = defaultExpireDays
	}
	if item.ScanTimeoutMs <= 0 {
		item.ScanTimeoutMs = int64(defaultDomainTimeout / time.Millisecond)
	}
	if item.Enabled && item.CronExpr == "" {
		return errors.New("启用域名巡检调度时必须填写 Cron 表达式")
	}
	return nil
}

type domainCheckResult struct {
	StatusCode   int
	ResponseTime int64
	ExpireAt     *time.Time
	DaysLeft     int64
	Issuer       string
}

func checkDomainHealth(ctx context.Context, domain string, timeout time.Duration, client *http.Client) (*domainCheckResult, error) {
	domain = normalizeDomain(domain)
	if domain == "" {
		return nil, errors.New("empty domain")
	}
	start := time.Now()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+domain, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result := &domainCheckResult{StatusCode: resp.StatusCode, ResponseTime: time.Since(start).Milliseconds()}

	address := net.JoinHostPort(domain, "443")
	dialer := &net.Dialer{Timeout: timeout}
	conn, err := tls.DialWithDialer(dialer, "tcp", address, &tls.Config{
		ServerName:         domain,
		InsecureSkipVerify: true, //nolint:gosec
	})
	if err != nil {
		return result, nil
	}
	defer conn.Close()

	state := conn.ConnectionState()
	if len(state.PeerCertificates) > 0 {
		cert := state.PeerCertificates[0]
		expireAt := cert.NotAfter
		result.ExpireAt = &expireAt
		result.DaysLeft = int64(math.Ceil(expireAt.Sub(time.Now()).Hours() / 24))
		result.Issuer = cert.Issuer.CommonName
	}
	return result, nil
}

func hostMetricQueryBuilders() map[string]func(instance string, threshold float64) string {
	return map[string]func(instance string, threshold float64) string{
		"cpu_usage": func(instance string, threshold float64) string {
			return fmt.Sprintf(`max(system_cpu_usage_percent{instance="%s"})`, instance)
		},
		"memory_usage": func(instance string, threshold float64) string {
			return fmt.Sprintf(`max(system_memory_usage_percent{instance="%s"})`, instance)
		},
		"disk_usage": func(instance string, threshold float64) string {
			return fmt.Sprintf(`max(system_disk_usage_percent{instance="%s"})`, instance)
		},
		"load1": func(instance string, threshold float64) string {
			return fmt.Sprintf(`max(system_load_average{instance="%s",period="1min"})`, instance)
		},
		"offline": func(instance string, threshold float64) string {
			window := int64(threshold)
			if window <= 0 {
				window = 5
			}
			return fmt.Sprintf(`sum(absent_over_time(system_cpu_usage_percent{instance="%s"}[%dm]))`, instance, window)
		},
	}
}

func queryPrometheusInstant(ctx context.Context, query string) (float64, error) {
	prometheusURL := strings.TrimSpace(config.Config.Monitor.Prometheus.URL)
	if prometheusURL == "" {
		return 0, errors.New("prometheus URL not configured")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimSuffix(prometheusURL, "/")+"/api/v1/query", nil)
	if err != nil {
		return 0, err
	}
	q := req.URL.Query()
	q.Set("query", query)
	q.Set("time", strconv.FormatInt(time.Now().Unix(), 10))
	req.URL.RawQuery = q.Encode()
	resp, err := (&http.Client{Timeout: hostScanTimeout}).Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var payload model.PrometheusQueryResult
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, err
	}
	if payload.Status != "success" || len(payload.Data.Result) == 0 {
		return 0, nil
	}
	item := payload.Data.Result[0]
	if len(item.Value) >= 2 {
		if valueText, ok := item.Value[1].(string); ok {
			value, _ := strconv.ParseFloat(valueText, 64)
			return normalizeFloat(value), nil
		}
	}
	if len(item.Values) > 0 {
		last := item.Values[len(item.Values)-1]
		if len(last) >= 2 {
			if valueText, ok := last[1].(string); ok {
				value, _ := strconv.ParseFloat(valueText, 64)
				return normalizeFloat(value), nil
			}
		}
	}
	return 0, nil
}

func checkDatabaseHealth(ctx context.Context, target cmdbModel.CmdbSQL, account *configAccountView, password string) (bool, int64, string) {
	start := time.Now()
	switch target.Type {
	case 1:
		connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=10s&readTimeout=10s&writeTimeout=10s&parseTime=true", account.Name, password, account.Host, account.Port, target.Name)
		db, err := sql.Open("mysql", connStr)
		if err != nil {
			return false, 0, err.Error()
		}
		defer db.Close()
		pingCtx, cancel := context.WithTimeout(ctx, databaseScanTimeout)
		defer cancel()
		if err = db.PingContext(pingCtx); err != nil {
			return false, time.Since(start).Milliseconds(), err.Error()
		}
		return true, time.Since(start).Milliseconds(), ""
	case 3:
		client := redis.NewClient(&redis.Options{
			Addr:         net.JoinHostPort(account.Host, strconv.Itoa(account.Port)),
			Username:     account.Name,
			Password:     password,
			DialTimeout:  databaseScanTimeout,
			ReadTimeout:  databaseScanTimeout,
			WriteTimeout: databaseScanTimeout,
			DB:           0,
		})
		defer client.Close()
		pingCtx, cancel := context.WithTimeout(ctx, databaseScanTimeout)
		defer cancel()
		if err := client.Ping(pingCtx).Err(); err != nil {
			return false, time.Since(start).Milliseconds(), err.Error()
		}
		return true, time.Since(start).Milliseconds(), ""
	default:
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(account.Host, strconv.Itoa(account.Port)), databaseScanTimeout)
		if err != nil {
			return false, time.Since(start).Milliseconds(), err.Error()
		}
		_ = conn.Close()
		return true, time.Since(start).Milliseconds(), ""
	}
}

func buildDatabaseAlertSummary(target cmdbModel.CmdbSQL, rule model.MonitorDBAlertRuleEntity, available bool, latencyMs int64, errMsg string) string {
	switch rule.MetricKey {
	case "connectivity":
		if available {
			return fmt.Sprintf("%s 数据库连接正常", target.Name)
		}
		return fmt.Sprintf("%s 数据库连接失败: %s", target.Name, firstNonEmptyAutomation(errMsg, "unknown error"))
	case "latency_ms":
		return fmt.Sprintf("%s 数据库连接耗时 %dms，阈值 %s %.0fms", target.Name, latencyMs, rule.Operator, rule.ThresholdValue)
	default:
		return fmt.Sprintf("%s 数据库状态已更新", target.Name)
	}
}

func hostMetricLabel(metricKey string) string {
	labels := map[string]string{"cpu_usage": "CPU 使用率", "memory_usage": "内存使用率", "disk_usage": "磁盘使用率", "load1": "1 分钟负载", "offline": "主机离线"}
	return labels[metricKey]
}

func databaseMetricLabel(metricKey string) string {
	labels := map[string]string{"connectivity": "连通性", "latency_ms": "连接延迟"}
	return labels[metricKey]
}

func databaseTypeText(dbType int) string {
	switch dbType {
	case 1:
		return "MySQL"
	case 2:
		return "PostgreSQL"
	case 3:
		return "Redis"
	case 4:
		return "MongoDB"
	case 5:
		return "Elasticsearch"
	default:
		return "Unknown"
	}
}

func normalizeDomain(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "https://")
	raw = strings.TrimPrefix(raw, "http://")
	raw = strings.TrimSuffix(raw, "/")
	return raw
}

func isSupportedComparison(operator string) bool {
	switch operator {
	case ">", ">=", "<", "<=", "==", "!=":
		return true
	default:
		return false
	}
}

func compareThreshold(value float64, operator string, threshold float64) bool {
	switch operator {
	case ">":
		return value > threshold
	case ">=":
		return value >= threshold
	case "<":
		return value < threshold
	case "<=":
		return value <= threshold
	case "==":
		return value == threshold
	case "!=":
		return value != threshold
	default:
		return false
	}
}

func normalizeFloat(value float64) float64 {
	return math.Round(value*100) / 100
}

func normalizeSeverity(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	switch value {
	case "P1", "P2", "P3", "P4":
		return value
	default:
		return "P3"
	}
}

func hostDisplayName(host cmdbModel.CmdbHost) string {
	return firstNonEmptyAutomation(host.HostName, host.Name, host.SSHIP, fmt.Sprintf("host-%d", host.ID))
}

func buildAutomationFingerprint(resourceType string, resourceID uint, eventKey string) string {
	return fmt.Sprintf("%s:%d:%s", resourceType, resourceID, eventKey)
}

func uniqueUintValues(values []uint) []uint {
	if len(values) == 0 {
		return values
	}
	seen := make(map[uint]struct{}, len(values))
	result := make([]uint, 0, len(values))
	for _, value := range values {
		if value == 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

func decodeUintList(raw string) []uint {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []uint{}
	}
	var values []uint
	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		return []uint{}
	}
	return values
}

func encodeAutomationJSON(value interface{}) string {
	data, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func firstNonEmptyAutomation(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func firstEnabledSchedule(items []model.MonitorDomainScheduleEntity) *model.MonitorDomainScheduleEntity {
	for i := range items {
		if items[i].Enabled {
			return &items[i]
		}
	}
	return nil
}
