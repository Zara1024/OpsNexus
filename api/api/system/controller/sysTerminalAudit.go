package controller

import (
	"dodevops-api/api/system/model"
	"dodevops-api/api/system/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetTerminalAuditSummary returns top-level terminal audit metrics.
func GetTerminalAuditSummary(c *gin.Context) {
	service.SysTerminalAuditService().GetTerminalAuditSummary(c)
}

// GetTerminalAuditSessionList returns the aggregated session list.
func GetTerminalAuditSessionList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))

	riskLevel := -1
	if value := c.Query("riskLevel"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			riskLevel = parsed
		}
	}

	sensitiveOnly := false
	if value := strings.TrimSpace(strings.ToLower(c.Query("sensitiveOnly"))); value == "1" || value == "true" {
		sensitiveOnly = true
	}

	query := model.TerminalAuditQuery{
		SessionID:     strings.TrimSpace(c.Query("sessionId")),
		HostID:        uint(parseUintQuery(c.Query("hostId"))),
		HostKeyword:   strings.TrimSpace(c.Query("hostKeyword")),
		Keyword:       strings.TrimSpace(c.Query("keyword")),
		RiskLevel:     riskLevel,
		SensitiveOnly: sensitiveOnly,
		BeginTime:     strings.TrimSpace(c.Query("beginTime")),
		EndTime:       strings.TrimSpace(c.Query("endTime")),
		PageSize:      pageSize,
		PageNum:       pageNum,
	}

	service.SysTerminalAuditService().GetTerminalAuditSessionList(c, query)
}

func parseUintQuery(value string) uint64 {
	parsed, _ := strconv.ParseUint(strings.TrimSpace(value), 10, 32)
	return parsed
}

// GetTerminalAuditSessionDetail returns one session's commands and header.
func GetTerminalAuditSessionDetail(c *gin.Context) {
	sessionID := strings.TrimSpace(c.Param("sessionId"))
	service.SysTerminalAuditService().GetTerminalAuditSessionDetail(c, sessionID)
}

// GetTerminalAuditSessionPlayback returns parsed recording events when available.
func GetTerminalAuditSessionPlayback(c *gin.Context) {
	sessionID := strings.TrimSpace(c.Param("sessionId"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))

	query := model.TerminalAuditPlaybackQuery{
		Keyword:   strings.TrimSpace(c.Query("keyword")),
		EventType: strings.TrimSpace(c.Query("eventType")),
		PageSize:  pageSize,
		PageNum:   pageNum,
	}

	service.SysTerminalAuditService().GetTerminalAuditSessionPlayback(c, sessionID, query)
}

// DownloadTerminalAuditRecording downloads one recording file when available.
func DownloadTerminalAuditRecording(c *gin.Context) {
	sessionID := strings.TrimSpace(c.Param("sessionId"))
	service.SysTerminalAuditService().DownloadTerminalAuditRecording(c, sessionID)
}
