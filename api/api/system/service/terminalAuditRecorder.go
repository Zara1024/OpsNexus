package service

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	systemDao "dodevops-api/api/system/dao"
	systemModel "dodevops-api/api/system/model"
)

const (
	terminalAuditStatusRecording = 1
	terminalAuditStatusCompleted = 2
	terminalAuditStatusAborted   = 3
	terminalAuditStorageLocal    = 1
)

type TerminalAuditRecorderOptions struct {
	AdminID   uint
	Username  string
	HostID    uint
	HostName  string
	HostIP    string
	SSHUser   string
	ClientIP  string
	UserAgent string
	Width     int
	Height    int
	FileDir   string
}

type TerminalAuditRecorder interface {
	SessionID() string
	RecordingID() uint
	RecordInput(data []byte)
	RecordOutput(data []byte)
	RecordResize(cols, rows int)
	Close(status int64, errMsg string) error
}

type terminalAuditRecorder struct {
	mu              sync.Mutex
	session         *systemModel.SysSessionRecording
	file            *os.File
	startTime       time.Time
	commandBuffer   strings.Builder
	sequence        int64
	inputCount      int64
	outputCount     int64
	resizeCount     int64
	commandCount    int64
	maxRiskLevel    int64
	hasSensitiveCmd bool
	closed          bool
}

func NewTerminalAuditRecorder(options TerminalAuditRecorderOptions) (TerminalAuditRecorder, error) {
	now := time.Now()
	sessionID, err := newTerminalAuditSessionID(now)
	if err != nil {
		return nil, err
	}

	baseDir := strings.TrimSpace(options.FileDir)
	if baseDir == "" {
		baseDir = "./upload/terminal-audit"
	}

	dayDir := filepath.Join(baseDir, now.Format("20060102"))
	if err = os.MkdirAll(dayDir, 0o755); err != nil {
		return nil, err
	}

	filePath := filepath.Join(dayDir, sessionID+".log")
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	session := &systemModel.SysSessionRecording{
		SessionID:       sessionID,
		AdminID:         options.AdminID,
		Username:        firstNonEmptyRecorder(options.Username, "unknown"),
		HostID:          options.HostID,
		HostName:        firstNonEmptyRecorder(options.HostName, "-"),
		HostIP:          options.HostIP,
		SSHUser:         firstNonEmptyRecorder(options.SSHUser, "-"),
		StartTime:       now,
		TerminalWidth:   int64(maxInt(options.Width, 80)),
		TerminalHeight:  int64(maxInt(options.Height, 24)),
		FilePath:        absPath,
		FileSize:        0,
		StorageType:     terminalAuditStorageLocal,
		InputCount:      0,
		OutputCount:     0,
		ResizeCount:     0,
		CommandCount:    0,
		ClientIP:        options.ClientIP,
		UserAgent:       options.UserAgent,
		RiskLevel:       0,
		HasSensitiveCmd: false,
		Status:          terminalAuditStatusRecording,
		ErrorMsg:        "",
		CreateTime:      now,
		UpdateTime:      now,
	}

	if err = systemDao.CreateTerminalAuditSession(session); err != nil {
		_ = file.Close()
		return nil, err
	}

	recorder := &terminalAuditRecorder{
		session:   session,
		file:      file,
		startTime: now,
	}
	recorder.writeFileLineLocked("SESSION START "+sessionID, now)
	return recorder, nil
}

func (r *terminalAuditRecorder) SessionID() string {
	return r.session.SessionID
}

func (r *terminalAuditRecorder) RecordingID() uint {
	return r.session.ID
}

func (r *terminalAuditRecorder) RecordInput(data []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed || len(data) == 0 {
		return
	}

	r.inputCount++
	now := time.Now()
	r.writeFileLineLocked("IN  "+sanitizeTerminalData(string(data)), now)
	r.captureCommandsLocked(data, now)
}

func (r *terminalAuditRecorder) RecordOutput(data []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed || len(data) == 0 {
		return
	}

	r.outputCount++
	r.writeFileLineLocked("OUT "+sanitizeTerminalData(string(data)), time.Now())
}

func (r *terminalAuditRecorder) RecordResize(cols, rows int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return
	}

	r.resizeCount++
	if cols > 0 {
		r.session.TerminalWidth = int64(cols)
	}
	if rows > 0 {
		r.session.TerminalHeight = int64(rows)
	}
	r.writeFileLineLocked("RESIZE cols="+intToString(cols)+" rows="+intToString(rows), time.Now())
}

func (r *terminalAuditRecorder) Close(status int64, errMsg string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return nil
	}
	r.closed = true

	now := time.Now()
	r.flushPendingCommandLocked(now)
	r.writeFileLineLocked("SESSION END status="+int64ToString(status)+" err="+strings.TrimSpace(errMsg), now)

	if r.file != nil {
		_ = r.file.Sync()
		_ = r.file.Close()
	}

	fileSize := int64(0)
	if info, err := os.Stat(r.session.FilePath); err == nil {
		fileSize = info.Size()
	}

	duration := int64(now.Sub(r.startTime).Seconds())
	r.session.EndTime = &now
	r.session.Duration = &duration
	r.session.FileSize = fileSize
	r.session.InputCount = r.inputCount
	r.session.OutputCount = r.outputCount
	r.session.ResizeCount = r.resizeCount
	r.session.CommandCount = r.commandCount
	r.session.RiskLevel = r.maxRiskLevel
	r.session.HasSensitiveCmd = r.hasSensitiveCmd
	r.session.Status = status
	r.session.ErrorMsg = strings.TrimSpace(errMsg)
	r.session.UpdateTime = now

	return systemDao.UpdateTerminalAuditSession(r.session.ID, map[string]interface{}{
		"end_time":          r.session.EndTime,
		"duration":          duration,
		"terminal_width":    r.session.TerminalWidth,
		"terminal_height":   r.session.TerminalHeight,
		"file_size":         fileSize,
		"input_count":       r.inputCount,
		"output_count":      r.outputCount,
		"resize_count":      r.resizeCount,
		"command_count":     r.commandCount,
		"risk_level":        r.maxRiskLevel,
		"has_sensitive_cmd": r.hasSensitiveCmd,
		"status":            status,
		"error_msg":         strings.TrimSpace(errMsg),
		"update_time":       now,
	})
}

func (r *terminalAuditRecorder) captureCommandsLocked(data []byte, executeTime time.Time) {
	for _, b := range data {
		switch b {
		case '\r', '\n':
			r.flushPendingCommandLocked(executeTime)
		case 8, 127:
			current := r.commandBuffer.String()
			if len(current) == 0 {
				continue
			}
			r.commandBuffer.Reset()
			r.commandBuffer.WriteString(current[:len(current)-1])
		default:
			if b < 32 {
				continue
			}
			r.commandBuffer.WriteByte(b)
		}
	}
}

func (r *terminalAuditRecorder) flushPendingCommandLocked(executeTime time.Time) {
	command := strings.TrimSpace(sanitizeCommand(r.commandBuffer.String()))
	r.commandBuffer.Reset()
	if command == "" {
		return
	}

	r.sequence++
	r.commandCount++
	isSensitive, riskLevel, riskReason := classifyCommandRisk(command)
	if riskLevel > r.maxRiskLevel {
		r.maxRiskLevel = riskLevel
	}
	if isSensitive {
		r.hasSensitiveCmd = true
	}

	_ = systemDao.CreateTerminalAuditCommand(&systemModel.SysCommandAudit{
		RecordingID: r.session.ID,
		SessionID:   r.session.SessionID,
		Command:     command,
		Timestamp:   executeTime.Sub(r.startTime).Seconds(),
		Sequence:    r.sequence,
		IsSensitive: isSensitive,
		RiskLevel:   riskLevel,
		RiskReason:  riskReason,
		ExecuteTime: executeTime,
		CreateTime:  time.Now(),
	})
}

func (r *terminalAuditRecorder) writeFileLineLocked(content string, at time.Time) {
	if r.file == nil {
		return
	}
	_, _ = r.file.WriteString("[" + at.Format("2006-01-02 15:04:05") + "] " + content + "\n")
}

func newTerminalAuditSessionID(now time.Time) (string, error) {
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return int64ToString(now.Unix()) + "_" + hex.EncodeToString(randomBytes), nil
}

func classifyCommandRisk(command string) (bool, int64, string) {
	assessment := ClassifyTerminalCommandRisk(command)
	return assessment.IsSensitive, assessment.RiskLevel, assessment.Reason
}

func sanitizeCommand(command string) string {
	replacer := strings.NewReplacer(
		"\u001b", "",
		"\u0000", "",
	)
	return replacer.Replace(command)
}

func sanitizeTerminalData(text string) string {
	text = strings.ReplaceAll(text, "\r", "\\r")
	text = strings.ReplaceAll(text, "\n", "\\n")
	text = strings.ReplaceAll(text, "\u001b", "")
	text = strings.ReplaceAll(text, "\u0000", "")
	return text
}

func intToString(value int) string {
	return int64ToString(int64(value))
}

func int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func firstNonEmptyRecorder(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func maxInt(value, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}
