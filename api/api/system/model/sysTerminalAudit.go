package model

import "dodevops-api/common/util"

// TerminalAuditQuery defines the supported filters for terminal audit sessions.
type TerminalAuditQuery struct {
	SessionID     string
	HostID        uint
	HostKeyword   string
	Keyword       string
	RiskLevel     int
	SensitiveOnly bool
	BeginTime     string
	EndTime       string
	PageSize      int
	PageNum       int
}

// TerminalAuditSummary represents the top-level metrics shown on the audit page.
type TerminalAuditSummary struct {
	TotalSessions       int64      `json:"totalSessions"`
	RecordedSessions    int64      `json:"recordedSessions"`
	CommandOnlySessions int64      `json:"commandOnlySessions"`
	TotalCommands       int64      `json:"totalCommands"`
	SensitiveCommands   int64      `json:"sensitiveCommands"`
	RiskySessions       int64      `json:"riskySessions"`
	LatestExecuteTime   util.HTime `json:"latestExecuteTime"`
}

// TerminalAuditSession represents one aggregated session row.
type TerminalAuditSession struct {
	RecordingID           uint       `json:"recordingId"`
	HostID                uint       `json:"hostId"`
	SessionID             string     `json:"sessionId"`
	Username              string     `json:"username"`
	HostName              string     `json:"hostName"`
	HostIP                string     `json:"hostIp"`
	SSHUser               string     `json:"sshUser"`
	StartTime             util.HTime `json:"startTime"`
	EndTime               util.HTime `json:"endTime"`
	Duration              int64      `json:"duration"`
	Status                int        `json:"status"`
	RiskLevel             int        `json:"riskLevel"`
	CommandCount          int64      `json:"commandCount"`
	SensitiveCommandCount int64      `json:"sensitiveCommandCount"`
	LatestRiskReason      string     `json:"latestRiskReason"`
	LatestCommand         string     `json:"latestCommand"`
	FilePath              string     `json:"filePath"`
	FileSize              int64      `json:"fileSize"`
	StorageType           int        `json:"storageType"`
	DataSource            string     `json:"dataSource"`
	SessionType           string     `json:"sessionType"`
	PlaybackAvailable     bool       `json:"playbackAvailable"`
	RecordingState        string     `json:"recordingState"`
	RecordingStateText    string     `json:"recordingStateText"`
	ActualFileSize        int64      `json:"actualFileSize"`
	RecordingWarning      string     `json:"recordingWarning"`
}

// TerminalAuditCommand represents a single command audit record.
type TerminalAuditCommand struct {
	ID             uint       `json:"id"`
	RecordingID    uint       `json:"recordingId"`
	SessionID      string     `json:"sessionId"`
	Command        string     `json:"command"`
	ElapsedSeconds float64    `json:"elapsedSeconds"`
	Sequence       int64      `json:"sequence"`
	IsSensitive    bool       `json:"isSensitive"`
	RiskLevel      int        `json:"riskLevel"`
	RiskReason     string     `json:"riskReason"`
	ExecuteTime    util.HTime `json:"executeTime"`
	CreateTime     util.HTime `json:"createTime"`
}

// TerminalAuditSessionDetail is the session header and command timeline payload.
type TerminalAuditSessionDetail struct {
	Session  TerminalAuditSession   `json:"session"`
	Commands []TerminalAuditCommand `json:"commands"`
}

// TerminalAuditPlaybackQuery defines the playback filters.
type TerminalAuditPlaybackQuery struct {
	Keyword   string
	EventType string
	PageSize  int
	PageNum   int
}

// TerminalAuditRecordingHealth describes one recording file's health.
type TerminalAuditRecordingHealth struct {
	State            string `json:"state"`
	StateText        string `json:"stateText"`
	Message          string `json:"message"`
	CanPlayback      bool   `json:"canPlayback"`
	CanDownload      bool   `json:"canDownload"`
	FileExists       bool   `json:"fileExists"`
	FileReadable     bool   `json:"fileReadable"`
	IsEmpty          bool   `json:"isEmpty"`
	RecordedFileSize int64  `json:"recordedFileSize"`
	ActualFileSize   int64  `json:"actualFileSize"`
	SizeMismatch     bool   `json:"sizeMismatch"`
}

// TerminalAuditPlaybackStats contains parsed playback metrics.
type TerminalAuditPlaybackStats struct {
	TotalEvents    int64      `json:"totalEvents"`
	InputEvents    int64      `json:"inputEvents"`
	OutputEvents   int64      `json:"outputEvents"`
	ResizeEvents   int64      `json:"resizeEvents"`
	SystemEvents   int64      `json:"systemEvents"`
	MatchedEvents  int64      `json:"matchedEvents"`
	FirstEventTime util.HTime `json:"firstEventTime"`
	LastEventTime  util.HTime `json:"lastEventTime"`
}

// TerminalAuditPlaybackEvent is one parsed recording line.
type TerminalAuditPlaybackEvent struct {
	Line            int        `json:"line"`
	EventType       string     `json:"eventType"`
	EventTypeText   string     `json:"eventTypeText"`
	At              util.HTime `json:"at"`
	RelativeSeconds float64    `json:"relativeSeconds"`
	Content         string     `json:"content"`
	Matched         bool       `json:"matched"`
}

// TerminalAuditPlaybackResponse returns paged playback data.
type TerminalAuditPlaybackResponse struct {
	Session   TerminalAuditSession         `json:"session"`
	Health    TerminalAuditRecordingHealth `json:"health"`
	Stats     TerminalAuditPlaybackStats   `json:"stats"`
	Events    []TerminalAuditPlaybackEvent `json:"events"`
	Total     int64                        `json:"total"`
	PageSize  int                          `json:"pageSize"`
	PageNum   int                          `json:"pageNum"`
	Keyword   string                       `json:"keyword"`
	EventType string                       `json:"eventType"`
}
