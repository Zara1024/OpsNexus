package dao

import (
	"dodevops-api/api/system/model"
	. "dodevops-api/pkg/db"
	"strings"
)

const terminalAuditAggregationSQL = `
SELECT
    sca.recording_id,
    sca.session_id,
    COUNT(*) AS command_count,
    SUM(CASE WHEN sca.is_sensitive = 1 THEN 1 ELSE 0 END) AS sensitive_command_count,
    MAX(COALESCE(sca.risk_level, 0)) AS max_command_risk_level,
    MIN(sca.execute_time) AS first_command_time,
    MAX(sca.execute_time) AS last_command_time,
    SUBSTRING_INDEX(
        GROUP_CONCAT(sca.command ORDER BY sca.execute_time DESC, sca.id DESC SEPARATOR '\n'),
        '\n',
        1
    ) AS latest_command,
    SUBSTRING_INDEX(
        GROUP_CONCAT(COALESCE(NULLIF(sca.risk_reason, ''), '') ORDER BY sca.execute_time DESC, sca.id DESC SEPARATOR '\n'),
        '\n',
        1
    ) AS latest_risk_reason
FROM sys_command_audit sca
GROUP BY sca.recording_id, sca.session_id
`

func terminalAuditFromSQL() string {
	return `
FROM (` + terminalAuditAggregationSQL + `) agg
LEFT JOIN sys_session_recording ssr
  ON ssr.id = agg.recording_id
 AND ssr.delete_time IS NULL
`
}

func buildTerminalAuditWhere(query model.TerminalAuditQuery) (string, []interface{}) {
	var (
		parts []string
		args  []interface{}
	)

	if query.SessionID != "" {
		parts = append(parts, "agg.session_id LIKE ?")
		args = append(args, "%"+query.SessionID+"%")
	}
	if query.HostID > 0 {
		parts = append(parts, "ssr.host_id = ?")
		args = append(args, query.HostID)
	}
	if query.HostKeyword != "" {
		keyword := "%" + query.HostKeyword + "%"
		parts = append(parts, "(COALESCE(ssr.host_name, '') LIKE ? OR COALESCE(ssr.host_ip, '') LIKE ?)")
		args = append(args, keyword, keyword)
	}
	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		parts = append(parts, `(agg.session_id LIKE ? OR agg.latest_command LIKE ? OR EXISTS (
			SELECT 1
			FROM sys_command_audit sca2
			WHERE sca2.session_id = agg.session_id
			  AND sca2.command LIKE ?
		))`)
		args = append(args, keyword, keyword, keyword)
	}
	if query.RiskLevel >= 0 {
		parts = append(parts, "GREATEST(COALESCE(ssr.risk_level, 0), COALESCE(agg.max_command_risk_level, 0)) = ?")
		args = append(args, query.RiskLevel)
	}
	if query.SensitiveOnly {
		parts = append(parts, "agg.sensitive_command_count > 0")
	}
	if query.BeginTime != "" {
		parts = append(parts, "agg.last_command_time >= ?")
		args = append(args, query.BeginTime)
	}
	if query.EndTime != "" {
		parts = append(parts, "agg.first_command_time <= ?")
		args = append(args, query.EndTime)
	}

	if len(parts) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(parts, " AND "), args
}

// GetTerminalAuditSummary returns the top-level audit metrics.
func GetTerminalAuditSummary() (summary model.TerminalAuditSummary, err error) {
	sql := `
SELECT
    COUNT(DISTINCT sca.session_id) AS total_sessions,
    COUNT(DISTINCT CASE WHEN ssr.id IS NOT NULL THEN sca.session_id END) AS recorded_sessions,
    COUNT(DISTINCT CASE WHEN ssr.id IS NULL THEN sca.session_id END) AS command_only_sessions,
    COUNT(*) AS total_commands,
    SUM(CASE WHEN sca.is_sensitive = 1 THEN 1 ELSE 0 END) AS sensitive_commands,
    COUNT(DISTINCT CASE WHEN COALESCE(sca.risk_level, 0) > 0 THEN sca.session_id END) AS risky_sessions,
    MAX(sca.execute_time) AS latest_execute_time
FROM sys_command_audit sca
LEFT JOIN sys_session_recording ssr
  ON ssr.id = sca.recording_id
 AND ssr.delete_time IS NULL
`

	err = Db.Raw(sql).Scan(&summary).Error
	return summary, err
}

// GetTerminalAuditSessionList returns a paginated aggregated session list.
func GetTerminalAuditSessionList(query model.TerminalAuditQuery) (sessions []model.TerminalAuditSession, count int64, err error) {
	fromSQL := terminalAuditFromSQL()
	whereSQL, args := buildTerminalAuditWhere(query)

	countSQL := "SELECT COUNT(*) " + fromSQL + whereSQL
	if err = Db.Raw(countSQL, args...).Scan(&count).Error; err != nil {
		return sessions, count, err
	}

	listSQL := `
SELECT
    agg.recording_id AS recording_id,
    COALESCE(ssr.host_id, 0) AS host_id,
    agg.session_id AS session_id,
    COALESCE(ssr.username, '') AS username,
    COALESCE(ssr.host_name, '') AS host_name,
    COALESCE(ssr.host_ip, '') AS host_ip,
    COALESCE(ssr.ssh_user, '') AS ssh_user,
    COALESCE(ssr.start_time, agg.first_command_time) AS start_time,
    COALESCE(ssr.end_time, agg.last_command_time) AS end_time,
    COALESCE(ssr.duration, TIMESTAMPDIFF(SECOND, agg.first_command_time, agg.last_command_time)) AS duration,
    COALESCE(ssr.status, 1) AS status,
    GREATEST(COALESCE(ssr.risk_level, 0), COALESCE(agg.max_command_risk_level, 0)) AS risk_level,
    agg.command_count AS command_count,
    agg.sensitive_command_count AS sensitive_command_count,
    COALESCE(NULLIF(agg.latest_risk_reason, ''), '') AS latest_risk_reason,
    COALESCE(agg.latest_command, '') AS latest_command,
    COALESCE(ssr.file_path, '') AS file_path,
    COALESCE(ssr.file_size, 0) AS file_size,
    COALESCE(ssr.storage_type, 0) AS storage_type,
    CASE WHEN ssr.id IS NULL THEN 'command' ELSE 'recording' END AS data_source
` + fromSQL + whereSQL + `
ORDER BY agg.last_command_time DESC
LIMIT ? OFFSET ?
`

	listArgs := append([]interface{}{}, args...)
	listArgs = append(listArgs, query.PageSize, (query.PageNum-1)*query.PageSize)
	err = Db.Raw(listSQL, listArgs...).Scan(&sessions).Error
	return sessions, count, err
}

// GetTerminalAuditSessionDetail returns the aggregated header and command list for a session.
func GetTerminalAuditSessionDetail(sessionID string) (detail model.TerminalAuditSessionDetail, err error) {
	fromSQL := terminalAuditFromSQL()
	whereSQL := " WHERE agg.session_id = ?"

	detailSQL := `
SELECT
    agg.recording_id AS recording_id,
    COALESCE(ssr.host_id, 0) AS host_id,
    agg.session_id AS session_id,
    COALESCE(ssr.username, '') AS username,
    COALESCE(ssr.host_name, '') AS host_name,
    COALESCE(ssr.host_ip, '') AS host_ip,
    COALESCE(ssr.ssh_user, '') AS ssh_user,
    COALESCE(ssr.start_time, agg.first_command_time) AS start_time,
    COALESCE(ssr.end_time, agg.last_command_time) AS end_time,
    COALESCE(ssr.duration, TIMESTAMPDIFF(SECOND, agg.first_command_time, agg.last_command_time)) AS duration,
    COALESCE(ssr.status, 1) AS status,
    GREATEST(COALESCE(ssr.risk_level, 0), COALESCE(agg.max_command_risk_level, 0)) AS risk_level,
    agg.command_count AS command_count,
    agg.sensitive_command_count AS sensitive_command_count,
    COALESCE(NULLIF(agg.latest_risk_reason, ''), '') AS latest_risk_reason,
    COALESCE(agg.latest_command, '') AS latest_command,
    COALESCE(ssr.file_path, '') AS file_path,
    COALESCE(ssr.file_size, 0) AS file_size,
    COALESCE(ssr.storage_type, 0) AS storage_type,
    CASE WHEN ssr.id IS NULL THEN 'command' ELSE 'recording' END AS data_source
` + fromSQL + whereSQL + `
LIMIT 1
`

	if err = Db.Raw(detailSQL, sessionID).Scan(&detail.Session).Error; err != nil {
		return detail, err
	}

	commandSQL := `
SELECT
    id,
    recording_id,
    session_id,
    command,
    timestamp AS elapsed_seconds,
    sequence,
    is_sensitive,
    COALESCE(risk_level, 0) AS risk_level,
    COALESCE(risk_reason, '') AS risk_reason,
    execute_time,
    create_time
FROM sys_command_audit
WHERE session_id = ?
ORDER BY sequence ASC, execute_time ASC, id ASC
`

	err = Db.Raw(commandSQL, sessionID).Scan(&detail.Commands).Error
	return detail, err
}
