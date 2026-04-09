package service

import (
	"encoding/json"
	"errors"
	"strings"

	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// MonitorAlertServiceInterface defines the alert center operations.
type MonitorAlertServiceInterface interface {
	GetAlertSummary(c *gin.Context)
	GetIncidentList(c *gin.Context, query model.MonitorIncidentQuery)
	GetWebhookLogList(c *gin.Context, query model.MonitorWebhookLogQuery)
	GetNotifyLogList(c *gin.Context, query model.MonitorNotifyLogQuery)
	GetNotifyRobotList(c *gin.Context)
	GetAlertSourceList(c *gin.Context)
	ReceiveWebhook(c *gin.Context, payload map[string]interface{})
	CreateNotifyRobot(c *gin.Context, req model.MonitorNotifyRobotUpsertRequest)
	UpdateNotifyRobot(c *gin.Context, id uint, req model.MonitorNotifyRobotUpsertRequest)
	UpdateNotifyRobotStatus(c *gin.Context, id uint, req model.MonitorStatusUpdateRequest)
	DeleteNotifyRobot(c *gin.Context, id uint)
	TestNotifyRobot(c *gin.Context, id uint, req model.MonitorNotifyRobotTestRequest)
	CreateAlertSource(c *gin.Context, req model.MonitorAlertSourceUpsertRequest)
	UpdateAlertSource(c *gin.Context, id uint, req model.MonitorAlertSourceUpsertRequest)
	UpdateAlertSourceStatus(c *gin.Context, id uint, req model.MonitorStatusUpdateRequest)
	DeleteAlertSource(c *gin.Context, id uint)
	GetAlertManagerStatus(c *gin.Context, query model.MonitorAlertManagerQuery)
	GetAlertManagerSilenceList(c *gin.Context, query model.MonitorAlertManagerQuery)
	CreateAlertManagerSilence(c *gin.Context, req model.MonitorAlertManagerSilenceCreateRequest)
	DeleteAlertManagerSilence(c *gin.Context, query model.MonitorAlertManagerQuery, silenceID string)
	GetAlertManagerReceiverList(c *gin.Context, query model.MonitorAlertManagerQuery)
}

type MonitorAlertServiceImpl struct{}

// NewMonitorAlertService creates an alert center service instance.
func NewMonitorAlertService() MonitorAlertServiceInterface {
	return &MonitorAlertServiceImpl{}
}

func (s *MonitorAlertServiceImpl) GetAlertSummary(c *gin.Context) {
	summary, err := dao.GetMonitorAlertSummary()
	if err != nil {
		result.Failed(c, 500, "获取告警中心摘要失败: "+err.Error())
		return
	}
	result.Success(c, summary)
}

func (s *MonitorAlertServiceImpl) GetIncidentList(c *gin.Context, query model.MonitorIncidentQuery) {
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	if query.PageNum < 1 {
		query.PageNum = 1
	}

	list, total, err := dao.GetMonitorIncidentList(query)
	if err != nil {
		result.Failed(c, 500, "获取事件告警列表失败: "+err.Error())
		return
	}
	result.SuccessWithPage(c, list, total, query.PageNum, query.PageSize)
}

func (s *MonitorAlertServiceImpl) GetWebhookLogList(c *gin.Context, query model.MonitorWebhookLogQuery) {
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	if query.PageNum < 1 {
		query.PageNum = 1
	}

	list, total, err := dao.GetMonitorWebhookLogList(query)
	if err != nil {
		result.Failed(c, 500, "获取告警历史失败: "+err.Error())
		return
	}
	result.SuccessWithPage(c, list, total, query.PageNum, query.PageSize)
}

func (s *MonitorAlertServiceImpl) GetNotifyLogList(c *gin.Context, query model.MonitorNotifyLogQuery) {
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	if query.PageNum < 1 {
		query.PageNum = 1
	}

	list, total, err := dao.GetMonitorWebhookNotifyLogList(query)
	if err != nil {
		result.Failed(c, 500, "获取告警推送日志失败: "+err.Error())
		return
	}
	result.SuccessWithPage(c, list, total, query.PageNum, query.PageSize)
}

func (s *MonitorAlertServiceImpl) GetNotifyRobotList(c *gin.Context) {
	list, err := dao.GetMonitorNotifyRobotList()
	if err != nil {
		result.Failed(c, 500, "获取通知机器人列表失败: "+err.Error())
		return
	}
	result.Success(c, list)
}

func (s *MonitorAlertServiceImpl) GetAlertSourceList(c *gin.Context) {
	list, err := dao.GetMonitorAlertSourceList()
	if err != nil {
		result.Failed(c, 500, "获取告警源列表失败: "+err.Error())
		return
	}
	result.Success(c, list)
}

func (s *MonitorAlertServiceImpl) CreateNotifyRobot(c *gin.Context, req model.MonitorNotifyRobotUpsertRequest) {
	robot, err := buildNotifyRobotEntity(0, req)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}
	if err = dao.CreateMonitorNotifyRobot(robot); err != nil {
		result.Failed(c, 500, "创建通知机器人失败: "+err.Error())
		return
	}
	result.Success(c, robot)
}

func (s *MonitorAlertServiceImpl) UpdateNotifyRobot(c *gin.Context, id uint, req model.MonitorNotifyRobotUpsertRequest) {
	if err := dao.EnsureAlertResourcesExist("robot", id); err != nil {
		result.Failed(c, 404, err.Error())
		return
	}

	robot, err := buildNotifyRobotEntity(id, req)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}
	if err = dao.UpdateMonitorNotifyRobot(robot); err != nil {
		result.Failed(c, 500, "更新通知机器人失败: "+err.Error())
		return
	}
	result.Success(c, robot)
}

func (s *MonitorAlertServiceImpl) UpdateNotifyRobotStatus(c *gin.Context, id uint, req model.MonitorStatusUpdateRequest) {
	if err := dao.EnsureAlertResourcesExist("robot", id); err != nil {
		result.Failed(c, 404, err.Error())
		return
	}
	if err := dao.UpdateMonitorNotifyRobotStatus(id, req.Status); err != nil {
		result.Failed(c, 500, "更新通知机器人状态失败: "+err.Error())
		return
	}
	result.Success(c, gin.H{"id": id, "status": req.Status})
}

func (s *MonitorAlertServiceImpl) DeleteNotifyRobot(c *gin.Context, id uint) {
	if err := dao.EnsureAlertResourcesExist("robot", id); err != nil {
		result.Failed(c, 404, err.Error())
		return
	}
	if err := dao.DeleteMonitorNotifyRobot(id); err != nil {
		result.Failed(c, 500, "删除通知机器人失败: "+err.Error())
		return
	}
	result.Success(c, gin.H{"id": id})
}

func (s *MonitorAlertServiceImpl) CreateAlertSource(c *gin.Context, req model.MonitorAlertSourceUpsertRequest) {
	source, err := buildAlertSourceEntity(0, req)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}
	if err = dao.CreateMonitorAlertSource(source); err != nil {
		result.Failed(c, 500, "创建告警源失败: "+err.Error())
		return
	}
	result.Success(c, source)
}

func (s *MonitorAlertServiceImpl) UpdateAlertSource(c *gin.Context, id uint, req model.MonitorAlertSourceUpsertRequest) {
	if err := dao.EnsureAlertResourcesExist("source", id); err != nil {
		result.Failed(c, 404, err.Error())
		return
	}

	source, err := buildAlertSourceEntity(id, req)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}
	if err = dao.UpdateMonitorAlertSource(source); err != nil {
		result.Failed(c, 500, "更新告警源失败: "+err.Error())
		return
	}
	result.Success(c, source)
}

func (s *MonitorAlertServiceImpl) UpdateAlertSourceStatus(c *gin.Context, id uint, req model.MonitorStatusUpdateRequest) {
	if err := dao.EnsureAlertResourcesExist("source", id); err != nil {
		result.Failed(c, 404, err.Error())
		return
	}
	if err := dao.UpdateMonitorAlertSourceStatus(id, req.Status); err != nil {
		result.Failed(c, 500, "更新告警源状态失败: "+err.Error())
		return
	}
	result.Success(c, gin.H{"id": id, "status": req.Status})
}

func (s *MonitorAlertServiceImpl) DeleteAlertSource(c *gin.Context, id uint) {
	if err := dao.EnsureAlertResourcesExist("source", id); err != nil {
		result.Failed(c, 404, err.Error())
		return
	}
	if err := dao.DeleteMonitorAlertSource(id); err != nil {
		result.Failed(c, 500, "删除告警源失败: "+err.Error())
		return
	}
	result.Success(c, gin.H{"id": id})
}

func buildNotifyRobotEntity(id uint, req model.MonitorNotifyRobotUpsertRequest) (*model.MonitorNotifyRobotEntity, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Type = strings.ToLower(strings.TrimSpace(req.Type))
	req.Webhook = strings.TrimSpace(req.Webhook)
	req.Secret = strings.TrimSpace(req.Secret)
	req.Remark = strings.TrimSpace(req.Remark)
	req.Server = strings.TrimSpace(req.Server)
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.Nickname = strings.TrimSpace(req.Nickname)
	req.Headers = strings.TrimSpace(req.Headers)
	req.Method = strings.ToUpper(strings.TrimSpace(req.Method))
	req.Template = strings.TrimSpace(req.Template)

	if req.Name == "" {
		return nil, errors.New("通知机器人名称不能为空")
	}

	validTypes := map[string]struct{}{
		"feishu":   {},
		"dingtalk": {},
		"wechat":   {},
		"email":    {},
		"webhook":  {},
		"teams":    {},
	}
	if _, ok := validTypes[req.Type]; !ok {
		return nil, errors.New("通知机器人类型不支持")
	}

	if req.Status != 0 && req.Status != 1 {
		return nil, errors.New("通知机器人状态只能是 0 或 1")
	}

	if req.Method == "" {
		req.Method = "POST"
	}
	validMethods := map[string]struct{}{
		"GET":   {},
		"POST":  {},
		"PUT":   {},
		"PATCH": {},
	}
	if _, ok := validMethods[req.Method]; !ok {
		return nil, errors.New("通知机器人请求方法仅支持 GET/POST/PUT/PATCH")
	}

	if req.Type == "email" {
		if req.Server == "" || req.Port <= 0 || req.Username == "" || req.Password == "" {
			return nil, errors.New("邮件机器人需填写 SMTP 服务地址、端口、用户名和密码")
		}
	} else if req.Webhook == "" {
		return nil, errors.New("Webhook 地址不能为空")
	}

	if req.Headers != "" {
		var headers map[string]interface{}
		if err := json.Unmarshal([]byte(req.Headers), &headers); err != nil {
			return nil, errors.New("自定义请求头必须是合法 JSON")
		}
	}

	return &model.MonitorNotifyRobotEntity{
		ID:       id,
		Name:     req.Name,
		Type:     req.Type,
		Webhook:  req.Webhook,
		Secret:   req.Secret,
		Status:   req.Status,
		Remark:   req.Remark,
		Server:   req.Server,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Headers:  req.Headers,
		Method:   req.Method,
		Template: req.Template,
	}, nil
}

func buildAlertSourceEntity(id uint, req model.MonitorAlertSourceUpsertRequest) (*model.MonitorAlertSourceEntity, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.AppKey = strings.TrimSpace(req.AppKey)
	req.APIBaseURL = strings.TrimSpace(req.APIBaseURL)
	req.Remark = strings.TrimSpace(req.Remark)

	if req.Name == "" {
		return nil, errors.New("告警源名称不能为空")
	}

	validTypes := map[int]struct{}{
		1: {},
		2: {},
		3: {},
		4: {},
	}
	if _, ok := validTypes[req.Type]; !ok {
		return nil, errors.New("告警源类型不支持")
	}
	if req.Status != 0 && req.Status != 1 {
		return nil, errors.New("告警源状态只能是 0 或 1")
	}

	return &model.MonitorAlertSourceEntity{
		ID:         id,
		Name:       req.Name,
		Type:       req.Type,
		AppKey:     req.AppKey,
		APIBaseURL: req.APIBaseURL,
		Status:     req.Status,
		Remark:     req.Remark,
		KeyID:      req.KeyID,
		HostID:     req.HostID,
	}, nil
}
