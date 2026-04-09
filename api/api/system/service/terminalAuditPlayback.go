package service

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"time"

	"dodevops-api/api/system/model"
	"dodevops-api/common/util"
)

const (
	terminalAuditPlaybackStateReady       = "ready"
	terminalAuditPlaybackStateCommandOnly = "command_only"
	terminalAuditPlaybackStateMissing     = "missing"
	terminalAuditPlaybackStateEmpty       = "empty"
	terminalAuditPlaybackStateUnsupported = "unsupported"
	terminalAuditPlaybackStateReadError   = "read_error"
)

func enrichTerminalAuditSession(session *model.TerminalAuditSession) model.TerminalAuditRecordingHealth {
	session.SessionType = inferTerminalAuditSessionType(*session)
	health := buildTerminalAuditRecordingHealth(*session)
	applyTerminalAuditHealth(session, health)
	return health
}

func applyTerminalAuditHealth(session *model.TerminalAuditSession, health model.TerminalAuditRecordingHealth) {
	session.PlaybackAvailable = health.CanPlayback
	session.RecordingState = health.State
	session.RecordingStateText = health.StateText
	session.ActualFileSize = health.ActualFileSize
	session.RecordingWarning = strings.TrimSpace(health.Message)
}

func inferTerminalAuditSessionType(session model.TerminalAuditSession) string {
	hostName := strings.ToLower(strings.TrimSpace(session.HostName))
	sshUser := strings.ToLower(strings.TrimSpace(session.SSHUser))

	switch {
	case strings.Contains(hostName, "/kubectl") || sshUser == "kubectl":
		return "kubectl"
	case strings.Contains(hostName, "/pod/"):
		return "pod"
	case strings.TrimSpace(session.HostIP) != "" || strings.TrimSpace(session.HostName) != "" || strings.TrimSpace(session.SSHUser) != "":
		return "ssh"
	default:
		return "unknown"
	}
}

func buildTerminalAuditRecordingHealth(session model.TerminalAuditSession) model.TerminalAuditRecordingHealth {
	health := model.TerminalAuditRecordingHealth{
		RecordedFileSize: session.FileSize,
		ActualFileSize:   session.FileSize,
	}

	if session.DataSource != "recording" {
		health.State = terminalAuditPlaybackStateCommandOnly
		health.StateText = "仅命令聚合"
		health.Message = "当前会话仅保留命令审计聚合数据，暂无可回放录像文件。"
		return health
	}

	if session.StorageType != terminalAuditStorageLocal {
		health.State = terminalAuditPlaybackStateUnsupported
		health.StateText = "存储暂不支持"
		health.Message = "当前会话不是本地录像存储，暂不支持回放和下载。"
		return health
	}

	filePath := strings.TrimSpace(session.FilePath)
	if filePath == "" {
		health.State = terminalAuditPlaybackStateMissing
		health.StateText = "文件缺失"
		health.Message = "录像文件路径为空，无法执行回放。"
		return health
	}

	info, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			health.State = terminalAuditPlaybackStateMissing
			health.StateText = "文件缺失"
			health.Message = "录像文件不存在，可能已被清理或落盘失败。"
			return health
		}
		health.State = terminalAuditPlaybackStateReadError
		health.StateText = "读取失败"
		health.Message = "录像文件状态检查失败: " + err.Error()
		return health
	}

	health.FileExists = true
	health.FileReadable = true
	health.CanDownload = true
	health.ActualFileSize = info.Size()

	if info.Size() == 0 {
		health.State = terminalAuditPlaybackStateEmpty
		health.StateText = "空文件"
		health.IsEmpty = true
		health.Message = "录像文件已生成但内容为空，无法回放。"
		return health
	}

	health.State = terminalAuditPlaybackStateReady
	health.StateText = "可回放"
	health.CanPlayback = true
	if session.FileSize > 0 && session.FileSize != info.Size() {
		health.SizeMismatch = true
		health.Message = "录像文件大小与数据库记录不一致，已按实际文件大小回放。"
		return health
	}
	health.Message = "录像文件状态正常，可直接回放或下载。"
	return health
}

func buildTerminalAuditPlaybackResponse(session model.TerminalAuditSession, query model.TerminalAuditPlaybackQuery) model.TerminalAuditPlaybackResponse {
	session.SessionType = inferTerminalAuditSessionType(session)
	return model.TerminalAuditPlaybackResponse{
		Session:   session,
		PageSize:  query.PageSize,
		PageNum:   query.PageNum,
		Keyword:   strings.TrimSpace(query.Keyword),
		EventType: strings.TrimSpace(query.EventType),
		Events:    []model.TerminalAuditPlaybackEvent{},
	}
}

func parseTerminalAuditPlayback(session model.TerminalAuditSession, query model.TerminalAuditPlaybackQuery) (model.TerminalAuditPlaybackResponse, error) {
	response := buildTerminalAuditPlaybackResponse(session, query)
	health := buildTerminalAuditRecordingHealth(session)
	response.Health = health
	applyTerminalAuditHealth(&response.Session, health)
	if !health.CanPlayback {
		return response, nil
	}

	file, err := os.Open(strings.TrimSpace(session.FilePath))
	if err != nil {
		response.Health.State = terminalAuditPlaybackStateReadError
		response.Health.StateText = "读取失败"
		response.Health.Message = "打开录像文件失败: " + err.Error()
		response.Health.CanPlayback = false
		applyTerminalAuditHealth(&response.Session, response.Health)
		return response, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	events := make([]model.TerminalAuditPlaybackEvent, 0)
	stats := model.TerminalAuditPlaybackStats{}
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		event := parseTerminalAuditPlaybackLine(lineNo, scanner.Text(), session.StartTime.Time)
		if event.Line == 0 {
			continue
		}

		updateTerminalAuditPlaybackStats(&stats, event)
		event.Matched = matchTerminalAuditPlaybackKeyword(event, query.Keyword)
		if event.Matched {
			stats.MatchedEvents++
		}
		if !matchTerminalAuditPlaybackEvent(event, query) {
			continue
		}
		events = append(events, event)
	}

	if err = scanner.Err(); err != nil {
		response.Health.State = terminalAuditPlaybackStateReadError
		response.Health.StateText = "读取失败"
		response.Health.Message = "解析录像文件失败: " + err.Error()
		response.Health.CanPlayback = false
		applyTerminalAuditHealth(&response.Session, response.Health)
		return response, nil
	}

	response.Stats = stats
	response.Total = int64(len(events))
	response.Events = paginateTerminalAuditPlaybackEvents(events, query.PageNum, query.PageSize)
	return response, nil
}

func parseTerminalAuditPlaybackLine(lineNo int, line string, sessionStart time.Time) model.TerminalAuditPlaybackEvent {
	line = strings.TrimSpace(line)
	if line == "" {
		return model.TerminalAuditPlaybackEvent{}
	}

	event := model.TerminalAuditPlaybackEvent{Line: lineNo}
	content := line
	if strings.HasPrefix(line, "[") {
		if end := strings.Index(line, "] "); end > 1 {
			if at, err := time.ParseInLocation("2006-01-02 15:04:05", line[1:end], time.Local); err == nil {
				event.At = util.HTime{Time: at}
				if !sessionStart.IsZero() {
					event.RelativeSeconds = roundTerminalAuditSeconds(at.Sub(sessionStart).Seconds())
				}
			}
			content = line[end+2:]
		}
	}

	eventType, eventTypeText, payload := classifyTerminalAuditPlaybackEvent(content)
	event.EventType = eventType
	event.EventTypeText = eventTypeText
	event.Content = decodeTerminalAuditPlaybackContent(payload)
	return event
}

func classifyTerminalAuditPlaybackEvent(content string) (string, string, string) {
	switch {
	case strings.HasPrefix(content, "IN  "):
		return "input", "输入", strings.TrimPrefix(content, "IN  ")
	case strings.HasPrefix(content, "OUT "):
		return "output", "输出", strings.TrimPrefix(content, "OUT ")
	case strings.HasPrefix(content, "RESIZE "):
		return "resize", "窗口变化", strings.TrimPrefix(content, "RESIZE ")
	case strings.HasPrefix(content, "SESSION START "):
		return "system", "会话开始", strings.TrimPrefix(content, "SESSION START ")
	case strings.HasPrefix(content, "SESSION END "):
		return "system", "会话结束", strings.TrimPrefix(content, "SESSION END ")
	default:
		return "system", "系统事件", content
	}
}

func decodeTerminalAuditPlaybackContent(content string) string {
	replacer := strings.NewReplacer(
		"\\r\\n", "\r\n",
		"\\n", "\n",
		"\\r", "\r",
		"\\t", "\t",
	)
	return replacer.Replace(content)
}

func updateTerminalAuditPlaybackStats(stats *model.TerminalAuditPlaybackStats, event model.TerminalAuditPlaybackEvent) {
	stats.TotalEvents++
	switch event.EventType {
	case "input":
		stats.InputEvents++
	case "output":
		stats.OutputEvents++
	case "resize":
		stats.ResizeEvents++
	default:
		stats.SystemEvents++
	}
	if event.At.IsZero() {
		return
	}
	if stats.FirstEventTime.IsZero() || event.At.Time.Before(stats.FirstEventTime.Time) {
		stats.FirstEventTime = event.At
	}
	if stats.LastEventTime.IsZero() || event.At.Time.After(stats.LastEventTime.Time) {
		stats.LastEventTime = event.At
	}
}

func matchTerminalAuditPlaybackKeyword(event model.TerminalAuditPlaybackEvent, keyword string) bool {
	keyword = strings.TrimSpace(strings.ToLower(keyword))
	if keyword == "" {
		return false
	}
	return strings.Contains(strings.ToLower(event.Content), keyword) ||
		strings.Contains(strings.ToLower(event.EventTypeText), keyword)
}

func matchTerminalAuditPlaybackEvent(event model.TerminalAuditPlaybackEvent, query model.TerminalAuditPlaybackQuery) bool {
	eventType := strings.TrimSpace(strings.ToLower(query.EventType))
	if eventType != "" && eventType != event.EventType {
		return false
	}

	keyword := strings.TrimSpace(query.Keyword)
	if keyword == "" {
		return true
	}
	return event.Matched
}

func paginateTerminalAuditPlaybackEvents(events []model.TerminalAuditPlaybackEvent, pageNum, pageSize int) []model.TerminalAuditPlaybackEvent {
	if pageSize < 1 {
		pageSize = 200
	}
	if pageNum < 1 {
		pageNum = 1
	}

	start := (pageNum - 1) * pageSize
	if start >= len(events) {
		return []model.TerminalAuditPlaybackEvent{}
	}

	end := start + pageSize
	if end > len(events) {
		end = len(events)
	}
	return events[start:end]
}

func roundTerminalAuditSeconds(value float64) float64 {
	if value <= 0 {
		return 0
	}
	return float64(int(value*100+0.5)) / 100
}
