package service

import (
	"fmt"
	"strings"

	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common/result"
	jwtutil "dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type alertDispatchSummary struct {
	NotifyIDs    []uint
	NotifyCount  int64
	SuccessCount int64
	FailedCount  int64
	Status       string
	ErrorMsg     string
}

func (s *MonitorAlertServiceImpl) TestNotifyRobot(c *gin.Context, id uint, req model.MonitorNotifyRobotTestRequest) {
	if err := dao.EnsureAlertResourcesExist("robot", id); err != nil {
		result.Failed(c, 404, err.Error())
		return
	}

	robot, err := dao.GetMonitorNotifyRobotByID(id)
	if err != nil {
		result.Failed(c, 500, "获取通知机器人失败: "+err.Error())
		return
	}

	dispatchReq := buildManualRobotTestPayload(c, robot, req)
	renderedMessage := renderAlertTemplate(robot.Template, dispatchReq)
	if strings.TrimSpace(renderedMessage) == "" {
		renderedMessage = buildAlertMessage(dispatchReq)
	}

	logEntry, summary, err := dispatchAlert(dispatchReq, []model.MonitorNotifyRobotEntity{*robot})
	if err != nil {
		result.Failed(c, 500, "发送测试通知失败: "+err.Error())
		return
	}

	result.Success(c, gin.H{
		"robotId":         robot.ID,
		"robotName":       robot.Name,
		"robotType":       robot.Type,
		"robotEnabled":    robot.Status == 1,
		"webhookLogId":    logEntry.ID,
		"status":          summary.Status,
		"notifyCount":     summary.NotifyCount,
		"successCount":    summary.SuccessCount,
		"failedCount":     summary.FailedCount,
		"renderedMessage": renderedMessage,
	})
}

func buildManualRobotTestPayload(c *gin.Context, robot *model.MonitorNotifyRobotEntity, req model.MonitorNotifyRobotTestRequest) model.MonitorWebhookReceiveRequest {
	source := strings.TrimSpace(req.Source)
	if source == "" {
		source = "manual-test"
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		title = fmt.Sprintf("OpsNexus 通知测试 - %s", robot.Name)
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		content = fmt.Sprintf(
			"这是一条来自 OpsNexus 告警推送中心的测试消息，用于验证机器人“%s”的连通性、模板渲染与日志落库。",
			robot.Name,
		)
	}

	level := strings.TrimSpace(req.Level)
	if level == "" {
		level = "info"
	}

	operatorID, _ := jwtutil.GetAdminId(c)
	operatorName, _ := jwtutil.GetAdminName(c)

	return model.MonitorWebhookReceiveRequest{
		Source:  source,
		Title:   title,
		Content: content,
		Level:   level,
		Tags: map[string]interface{}{
			"mode":      "manual-test",
			"robotId":   robot.ID,
			"robotName": robot.Name,
			"robotType": robot.Type,
		},
		Extra: map[string]interface{}{
			"mode":         "manual-test",
			"robotId":      robot.ID,
			"robotName":    robot.Name,
			"robotType":    robot.Type,
			"robotEnabled": robot.Status == 1,
			"operatorId":   operatorID,
			"operatorName": operatorName,
		},
		NotifyRobotIDs: []uint{robot.ID},
	}
}

func dispatchAlert(req model.MonitorWebhookReceiveRequest, robots []model.MonitorNotifyRobotEntity) (*model.MonitorWebhookLogEntity, *alertDispatchSummary, error) {
	logEntry := &model.MonitorWebhookLogEntity{
		Source:         req.Source,
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
	if err := dao.CreateMonitorWebhookLog(logEntry); err != nil {
		return nil, nil, err
	}

	notifyIDs := make([]uint, 0, len(robots))
	errorMessages := make([]string, 0, len(robots))
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
			errorMessages = append(errorMessages, errorMsg)
			failedCount++
		} else {
			successCount++
		}

		_ = dao.CreateMonitorWebhookNotifyLog(&model.MonitorWebhookNotifyLogEntity{
			WebhookLogID: logEntry.ID,
			RobotID:      robot.ID,
			RobotName:    robot.Name,
			RobotType:    robot.Type,
			Status:       status,
			ErrorMsg:     errorMsg,
		})
	}

	notifyCount := int64(len(robots))
	status := calculateDispatchStatus(req, len(robots), successCount, failedCount)
	errorMsg := truncateText(strings.Join(errorMessages, " | "), 1000)
	if err := dao.UpdateMonitorWebhookLogResult(
		logEntry.ID,
		status,
		errorMsg,
		encodeJSONValue(notifyIDs),
		notifyCount,
		successCount,
		failedCount,
	); err != nil {
		return nil, nil, err
	}

	return logEntry, &alertDispatchSummary{
		NotifyIDs:    notifyIDs,
		NotifyCount:  notifyCount,
		SuccessCount: successCount,
		FailedCount:  failedCount,
		Status:       status,
		ErrorMsg:     errorMsg,
	}, nil
}
