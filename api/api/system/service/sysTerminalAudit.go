package service

import (
	"dodevops-api/api/system/dao"
	"dodevops-api/api/system/model"
	"dodevops-api/common/result"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type ISysTerminalAuditService interface {
	GetTerminalAuditSummary(c *gin.Context)
	GetTerminalAuditSessionList(c *gin.Context, query model.TerminalAuditQuery)
	GetTerminalAuditSessionDetail(c *gin.Context, sessionID string)
	GetTerminalAuditSessionPlayback(c *gin.Context, sessionID string, query model.TerminalAuditPlaybackQuery)
	DownloadTerminalAuditRecording(c *gin.Context, sessionID string)
}

type SysTerminalAuditServiceImpl struct{}

func (s SysTerminalAuditServiceImpl) GetTerminalAuditSummary(c *gin.Context) {
	summary, err := dao.GetTerminalAuditSummary()
	if err != nil {
		result.Failed(c, 500, "获取终端审计汇总失败: "+err.Error())
		return
	}
	result.Success(c, summary)
}

func (s SysTerminalAuditServiceImpl) GetTerminalAuditSessionList(c *gin.Context, query model.TerminalAuditQuery) {
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	if query.PageNum < 1 {
		query.PageNum = 1
	}

	sessions, count, err := dao.GetTerminalAuditSessionList(query)
	if err != nil {
		result.Failed(c, 500, "获取终端审计会话列表失败: "+err.Error())
		return
	}
	for i := range sessions {
		enrichTerminalAuditSession(&sessions[i])
	}

	result.Success(c, map[string]interface{}{
		"total":    count,
		"pageSize": query.PageSize,
		"pageNum":  query.PageNum,
		"list":     sessions,
	})
}

func (s SysTerminalAuditServiceImpl) GetTerminalAuditSessionDetail(c *gin.Context, sessionID string) {
	detail, err := dao.GetTerminalAuditSessionDetail(sessionID)
	if err != nil {
		result.Failed(c, 500, "获取终端审计详情失败: "+err.Error())
		return
	}
	if detail.Session.SessionID == "" {
		result.Failed(c, 404, "未找到对应的终端审计会话")
		return
	}
	enrichTerminalAuditSession(&detail.Session)
	result.Success(c, detail)
}

func (s SysTerminalAuditServiceImpl) GetTerminalAuditSessionPlayback(c *gin.Context, sessionID string, query model.TerminalAuditPlaybackQuery) {
	if query.PageSize < 1 {
		query.PageSize = 200
	}
	if query.PageNum < 1 {
		query.PageNum = 1
	}

	detail, err := dao.GetTerminalAuditSessionDetail(sessionID)
	if err != nil {
		result.Failed(c, 500, "获取录制详情失败: "+err.Error())
		return
	}
	if detail.Session.SessionID == "" {
		result.Failed(c, 404, "未找到对应的录制会话")
		return
	}

	playback, err := parseTerminalAuditPlayback(detail.Session, query)
	if err != nil {
		result.Failed(c, 500, "解析录制文件失败: "+err.Error())
		return
	}
	result.Success(c, playback)
}

func (s SysTerminalAuditServiceImpl) DownloadTerminalAuditRecording(c *gin.Context, sessionID string) {
	detail, err := dao.GetTerminalAuditSessionDetail(sessionID)
	if err != nil {
		result.Failed(c, 500, "获取录制详情失败: "+err.Error())
		return
	}
	if detail.Session.SessionID == "" {
		result.Failed(c, 404, "未找到对应的录制会话")
		return
	}
	health := enrichTerminalAuditSession(&detail.Session)
	if detail.Session.DataSource != "recording" {
		result.Failed(c, 400, "当前会话暂无录制文件")
		return
	}
	if detail.Session.StorageType != 1 {
		result.Failed(c, 501, "当前仅支持本地存储录制文件下载")
		return
	}
	if strings.TrimSpace(detail.Session.FilePath) == "" {
		result.Failed(c, 404, "录制文件路径为空")
		return
	}
	if !health.FileExists {
		result.Failed(c, 404, firstNonEmptyRecorder(health.Message, "录制文件不存在"))
		return
	}

	c.FileAttachment(detail.Session.FilePath, filepath.Base(detail.Session.FilePath))
}

var sysTerminalAuditService = SysTerminalAuditServiceImpl{}

func SysTerminalAuditService() ISysTerminalAuditService {
	return &sysTerminalAuditService
}
