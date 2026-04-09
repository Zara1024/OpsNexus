package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	configModel "dodevops-api/api/configcenter/model"
	configService "dodevops-api/api/configcenter/service"
	systemModel "dodevops-api/api/system/model"
	"dodevops-api/common"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	sqlWorkOrderStatusPending   = 1
	sqlWorkOrderStatusApproved  = 2
	sqlWorkOrderStatusRejected  = 3
	sqlWorkOrderStatusExecuting = 4
	sqlWorkOrderStatusSuccess   = 5
	sqlWorkOrderStatusFailed    = 6

	sqlRiskLow    = 0
	sqlRiskMedium = 1
	sqlRiskHigh   = 2
)

type SQLWorkOrderService struct {
	workDAO   *dao.SQLWorkOrderDao
	sqlDAO    *dao.CmdbSQLDao
	recordDAO *dao.CmdbSQLRecordDao
}

func NewSQLWorkOrderServiceWithDB() *SQLWorkOrderService {
	dbConn := common.GetDB()
	return &SQLWorkOrderService{
		workDAO:   dao.NewSQLWorkOrderDao(dbConn),
		sqlDAO:    dao.NewCmdbSQLDao(dbConn),
		recordDAO: dao.NewCmdbSQLRecordDao(dbConn),
	}
}

func (s *SQLWorkOrderService) GetSummary(c *gin.Context) {
	summary, err := s.workDAO.GetSummary()
	if err != nil {
		result.Failed(c, 500, "failed to get SQL work order summary: "+err.Error())
		return
	}
	result.Success(c, summary)
}

func (s *SQLWorkOrderService) List(c *gin.Context, query model.CmdbSQLWorkOrderQuery) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}

	list, total, err := s.workDAO.List(query)
	if err != nil {
		result.Failed(c, 500, "failed to get SQL work order list: "+err.Error())
		return
	}
	result.Success(c, model.CmdbSQLWorkOrderListResponse{
		List:     list,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	})
}

func (s *SQLWorkOrderService) Detail(c *gin.Context, id uint) {
	order, err := s.workDAO.GetByID(id)
	if err != nil {
		result.Failed(c, 404, "SQL work order not found")
		return
	}
	result.Success(c, order)
}

func (s *SQLWorkOrderService) Create(c *gin.Context, req model.CmdbSQLWorkOrderCreateRequest) {
	sqlText := strings.TrimSpace(req.SQL)
	if sqlText == "" {
		result.Failed(c, 400, "SQL statement cannot be empty")
		return
	}

	target, account, decrypted, schemaName, err := s.resolveConnectionInfo(req.DatabaseID, req.DatabaseName)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}

	opType := detectSQLWorkOrderOperationType(sqlText)
	switch target.Type {
	case 1:
		if opType == "" {
			result.Failed(c, 400, "unsupported SQL type; only SELECT/INSERT/UPDATE/DELETE/CREATE/ALTER/DROP/TRUNCATE are supported")
			return
		}
		if opType == "SELECT" {
			result.Failed(c, 400, "SELECT statements should be executed directly; SQL work orders are only for change statements")
			return
		}
		if err = ensureMySQLSchema(c.Request.Context(), account.Host, account.Port, account.Name, decrypted, schemaName); err != nil {
			result.Failed(c, 400, "database connection or schema validation failed: "+err.Error())
			return
		}
	case 2:
		if opType == "" {
			result.Failed(c, 400, "unsupported SQL type; only SELECT/INSERT/UPDATE/DELETE/CREATE/ALTER/DROP/TRUNCATE are supported")
			return
		}
		if opType == "SELECT" {
			result.Failed(c, 400, "SELECT statements should be executed directly; SQL work orders are only for change statements")
			return
		}
		sqlDB, dbErr := openPostgreSQLConnection(c.Request.Context(), account.Host, account.Port, account.Name, decrypted, schemaName)
		if dbErr != nil {
			result.Failed(c, 400, "database connection or schema validation failed: "+dbErr.Error())
			return
		}
		_ = sqlDB.Close()
	case 3:
		opType = detectRedisWorkOrderOperationType(sqlText)
		if opType == "" {
			result.Failed(c, 400, "unsupported Redis command type")
			return
		}
		kind, blocked := classifyRedisCommand(sqlText)
		if blocked {
			result.Failed(c, 400, "redis high-risk commands are blocked by policy")
			return
		}
		if kind == redisCommandKindRead {
			result.Failed(c, 400, "read-only Redis commands should be executed directly; work orders are only for change commands")
			return
		}
	default:
		result.Failed(c, 400, "unsupported database type")
		return
	}

	admin, username := currentSQLWorkOrderAdmin(c)
	riskLevel, riskSummary, requireApproval, affectedTables, rollbackSQL, rollbackHint := analyzeSQLWorkOrder(sqlText)
	if target.Type == 3 {
		riskLevel, riskSummary, requireApproval, affectedTables, rollbackSQL, rollbackHint = analyzeRedisWorkOrder(sqlText)
	}
	now := time.Now()
	order := &model.CmdbSQLWorkOrder{
		OrderNo:         generateSQLWorkOrderNo(now),
		Title:           buildSQLWorkOrderTitle(req.Title, opType, schemaName),
		Reason:          strings.TrimSpace(req.Reason),
		DatabaseID:      target.ID,
		DatabaseName:    schemaName,
		InstanceName:    target.Name,
		InstanceHost:    fmt.Sprintf("%s:%d", account.Host, account.Port),
		AccountID:       target.AccountID,
		OperationType:   opType,
		SQLContent:      sqlText,
		RiskLevel:       riskLevel,
		RiskSummary:     riskSummary,
		AffectedTables:  affectedTables,
		RollbackSQL:     rollbackSQL,
		RollbackHint:    rollbackHint,
		RequireApproval: requireApproval,
		ApplicantID:     admin.ID,
		ApplicantName:   username,
		ClientIP:        util.GetClientIP(c.Request),
		Status:          sqlWorkOrderStatusPending,
		ResultStatus:    "PENDING",
		CreateTime:      util.HTime{Time: now},
		UpdateTime:      util.HTime{Time: now},
	}
	if !requireApproval {
		order.Status = sqlWorkOrderStatusApproved
	}

	if err = s.workDAO.Create(order); err != nil {
		result.Failed(c, 500, "failed to create SQL work order: "+err.Error())
		return
	}
	result.Success(c, order)
}

func (s *SQLWorkOrderService) Approve(c *gin.Context, id uint, req model.CmdbSQLWorkOrderActionRequest) {
	order, err := s.workDAO.GetByID(id)
	if err != nil {
		result.Failed(c, 404, "SQL work order not found")
		return
	}
	if order.Status != sqlWorkOrderStatusPending {
		result.Failed(c, 400, "the current SQL work order status does not allow approval")
		return
	}

	admin, username := currentSQLWorkOrderAdmin(c)
	now := time.Now()
	order.Status = sqlWorkOrderStatusApproved
	order.ApproverID = admin.ID
	order.ApproverName = username
	order.ApprovalComment = strings.TrimSpace(req.Comment)
	order.ApprovalTime = &now
	order.UpdateTime = util.HTime{Time: now}
	if err = s.workDAO.Update(order); err != nil {
		result.Failed(c, 500, "failed to approve SQL work order: "+err.Error())
		return
	}
	result.Success(c, order)
}

func (s *SQLWorkOrderService) Reject(c *gin.Context, id uint, req model.CmdbSQLWorkOrderActionRequest) {
	order, err := s.workDAO.GetByID(id)
	if err != nil {
		result.Failed(c, 404, "SQL work order not found")
		return
	}
	if order.Status != sqlWorkOrderStatusPending {
		result.Failed(c, 400, "the current SQL work order status does not allow rejection")
		return
	}

	admin, username := currentSQLWorkOrderAdmin(c)
	now := time.Now()
	order.Status = sqlWorkOrderStatusRejected
	order.ApproverID = admin.ID
	order.ApproverName = username
	order.ApprovalComment = firstNonEmptySQL(strings.TrimSpace(req.Comment), "rejected")
	order.ApprovalTime = &now
	order.ResultStatus = "REJECTED"
	order.ResultMessage = order.ApprovalComment
	order.UpdateTime = util.HTime{Time: now}
	if err = s.workDAO.Update(order); err != nil {
		result.Failed(c, 500, "failed to reject SQL work order: "+err.Error())
		return
	}
	result.Success(c, order)
}

func (s *SQLWorkOrderService) Execute(c *gin.Context, id uint, req model.CmdbSQLWorkOrderActionRequest) {
	order, err := s.workDAO.GetByID(id)
	if err != nil {
		result.Failed(c, 404, "SQL work order not found")
		return
	}
	if !(order.Status == sqlWorkOrderStatusApproved || (!order.RequireApproval && order.Status == sqlWorkOrderStatusPending)) {
		result.Failed(c, 400, "the current SQL work order status does not allow execution")
		return
	}

	target, account, decrypted, schemaName, err := s.resolveConnectionInfo(order.DatabaseID, order.DatabaseName)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}

	executor, executorName := currentSQLWorkOrderAdmin(c)
	start := time.Now()
	order.Status = sqlWorkOrderStatusExecuting
	order.ExecutorID = executor.ID
	order.ExecutorName = executorName
	order.ExecutionStart = &start
	order.UpdateTime = util.HTime{Time: start}
	_ = s.workDAO.Update(order)

	var (
		executionTime int64
		affectedRows  int64
		execErr       error
	)

	switch target.Type {
	case 1, 2:
		sqlDB, connErr := openTypedSQLWorkOrderConnection(target.Type, account.Host, account.Port, account.Name, decrypted, schemaName)
		if connErr != nil {
			s.markExecutionFailed(order, start, "database connection failed: "+connErr.Error())
			result.Failed(c, 500, "SQL work order execution failed: database connection failed")
			return
		}
		defer sqlDB.Close()

		if backupPreview, backupRows, backupHint := buildSQLRollbackPreview(c.Request.Context(), sqlDB, order); backupPreview != "" {
			order.BackupPreview = backupPreview
			order.BackupRowCount = backupRows
			order.RollbackHint = firstNonEmptySQL(backupHint, order.RollbackHint)
		}

		execCtx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		runStart := time.Now()
		execResult, runErr := sqlDB.ExecContext(execCtx, order.SQLContent)
		executionTime = time.Since(runStart).Milliseconds()
		if execResult != nil {
			if rows, rowsErr := execResult.RowsAffected(); rowsErr == nil {
				affectedRows = rows
			}
		}
		execErr = runErr
	case 3:
		kind, blocked := classifyRedisCommand(order.SQLContent)
		if blocked {
			s.markExecutionFailed(order, start, "redis high-risk commands are blocked by policy")
			result.Failed(c, 400, "redis high-risk commands are blocked by policy")
			return
		}
		if kind == redisCommandKindRead {
			s.markExecutionFailed(order, start, "read-only Redis commands should be executed directly")
			result.Failed(c, 400, "read-only Redis commands should be executed directly")
			return
		}

		client := buildRedisClient(account.Host, account.Port, decrypted, parseRedisDatabaseIndex(schemaName))
		defer client.Close()

		execCtx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()
		if err = client.Ping(execCtx).Err(); err != nil {
			s.markExecutionFailed(order, start, "redis connection failed: "+err.Error())
			result.Failed(c, 500, "failed to execute SQL work order: redis connection failed")
			return
		}

		runStart := time.Now()
		args := strings.Fields(strings.TrimSpace(order.SQLContent))
		commandName := strings.ToUpper(args[0])
		commandArgs := make([]interface{}, 0, len(args))
		commandArgs = append(commandArgs, commandName)
		for _, arg := range args[1:] {
			commandArgs = append(commandArgs, arg)
		}
		_, execErr = client.Do(execCtx, commandArgs...).Result()
		executionTime = time.Since(runStart).Milliseconds()
		affectedRows = 1
	default:
		result.Failed(c, 400, "unsupported database type")
		return
	}

	if execErr != nil {
		s.markExecutionFailed(order, start, execErr.Error())
		_ = NewCmdbSQLRecordService(s.recordDAO).RecordSQLExecution(
			order.InstanceHost,
			order.DatabaseName,
			order.OperationType,
			order.SQLContent,
			executorName,
			util.GetClientIP(c.Request),
			0,
			affectedRows,
			executionTime,
			0,
			"FAILED",
		)
		result.Failed(c, 500, "failed to execute SQL work order: "+execErr.Error())
		return
	}

	end := time.Now()
	order.Status = sqlWorkOrderStatusSuccess
	order.ExecutionEnd = &end
	order.ExecutionTime = end.Sub(start).Milliseconds()
	order.AffectedRows = affectedRows
	order.ResultStatus = "SUCCESS"
	order.ResultMessage = firstNonEmptySQL(strings.TrimSpace(req.Comment), "execution succeeded")
	order.UpdateTime = util.HTime{Time: end}

	if err = s.workDAO.Update(order); err != nil {
		result.Failed(c, 500, "failed to update SQL work order execution result: "+err.Error())
		return
	}

	_ = NewCmdbSQLRecordService(s.recordDAO).RecordSQLExecution(
		order.InstanceHost,
		order.DatabaseName,
		order.OperationType,
		order.SQLContent,
		executorName,
		util.GetClientIP(c.Request),
		0,
		affectedRows,
		executionTime,
		0,
		"SUCCESS",
	)
	result.Success(c, order)
}

func (s *SQLWorkOrderService) markExecutionFailed(order *model.CmdbSQLWorkOrder, startedAt time.Time, message string) {
	end := time.Now()
	order.Status = sqlWorkOrderStatusFailed
	order.ExecutionEnd = &end
	order.ExecutionTime = end.Sub(startedAt).Milliseconds()
	order.ResultStatus = "FAILED"
	order.ResultMessage = strings.TrimSpace(message)
	order.UpdateTime = util.HTime{Time: end}
	_ = s.workDAO.Update(order)
}

func (s *SQLWorkOrderService) resolveConnectionInfo(databaseID uint, databaseName string) (*model.CmdbSQL, *configModel.AccountAuth, string, string, error) {
	if databaseID == 0 && strings.TrimSpace(databaseName) == "" {
		return nil, nil, "", "", errors.New("database id or name is required")
	}
	target, schemaName, err := resolveCmdbSQLTarget(databaseID, databaseName, s.sqlDAO.GetByID, s.sqlDAO.GetByName)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to resolve database asset: %w", err)
	}
	account, err := configService.NewAccountAuthService().GetByID(target.AccountID)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to load database account: %w", err)
	}
	decrypted, err := configService.NewAccountAuthService().DecryptPassword(target.AccountID)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to decrypt database password: %w", err)
	}

	return target, account, decrypted, schemaName, nil
}

func currentSQLWorkOrderAdmin(c *gin.Context) (systemModel.SysAdmin, string) {
	if admin, err := jwt.GetAdmin(c); err == nil && admin != nil {
		return systemModel.SysAdmin{
			ID:       admin.ID,
			Username: admin.Username,
			Nickname: admin.Nickname,
			Email:    admin.Email,
			Phone:    admin.Phone,
			Note:     admin.Note,
		}, firstNonEmptySQL(admin.Username, admin.Nickname, "unknown")
	}
	if userObj, exists := c.Get(constant.ContextKeyUserObj); exists {
		if admin, ok := userObj.(systemModel.SysAdmin); ok {
			return admin, firstNonEmptySQL(admin.Username, admin.Nickname, "unknown")
		}
		if admin, ok := userObj.(*systemModel.SysAdmin); ok && admin != nil {
			return *admin, firstNonEmptySQL(admin.Username, admin.Nickname, "unknown")
		}
	}
	return systemModel.SysAdmin{}, "unknown"
}

func openSQLWorkOrderConnection(host string, port int, username, password, databaseName string) (*sql.DB, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 5*time.Second)
	if err != nil {
		return nil, err
	}
	_ = conn.Close()

	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?timeout=60s&readTimeout=60s&writeTimeout=60s&parseTime=true&interpolateParams=true",
		username, password, host, port, databaseName,
	)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)
	return db, db.Ping()
}

func openTypedSQLWorkOrderConnection(dbType int, host string, port int, username, password, databaseName string) (*sql.DB, error) {
	switch dbType {
	case 2:
		return openPostgreSQLConnection(context.Background(), host, port, username, password, databaseName)
	default:
		return openSQLWorkOrderConnection(host, port, username, password, databaseName)
	}
}

func ensureMySQLSchema(ctx context.Context, host string, port int, username, password, schemaName string) error {
	rootConnStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/?timeout=60s&readTimeout=60s&writeTimeout=60s&parseTime=true&interpolateParams=true",
		username, password, host, port,
	)
	rootDB, err := sql.Open("mysql", rootConnStr)
	if err != nil {
		return err
	}
	defer rootDB.Close()

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var count int
	err = rootDB.QueryRowContext(
		timeoutCtx,
		"SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?",
		schemaName,
	).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("schema %q does not exist", schemaName)
	}
	return nil
}

func analyzeSQLWorkOrder(sqlText string) (int, string, bool, string, string, string) {
	normalized := strings.TrimSpace(sqlText)
	opType := detectSQLWorkOrderOperationType(normalized)
	tables := extractSQLWorkOrderTables(normalized)
	sort.Strings(tables)
	tables = uniqueSQLStrings(tables)

	reasons := make([]string, 0)
	riskLevel := sqlRiskLow
	requireApproval := opType != "SELECT"

	statements := splitSQLWorkOrderStatements(normalized)
	if len(statements) > 1 {
		riskLevel = sqlRiskHigh
		reasons = append(reasons, "contains multiple SQL statements and requires extra review")
	}

	switch opType {
	case "INSERT":
		riskLevel = maxSQLRisk(riskLevel, sqlRiskMedium)
		reasons = append(reasons, "INSERT writes new data; confirm unique keys and rollback conditions first")
	case "UPDATE":
		riskLevel = maxSQLRisk(riskLevel, sqlRiskMedium)
		reasons = append(reasons, "UPDATE modifies existing data; verify the affected scope before execution")
		if !strings.Contains(strings.ToUpper(normalized), " WHERE ") {
			riskLevel = maxSQLRisk(riskLevel, sqlRiskHigh)
			reasons = append(reasons, "UPDATE does not include a WHERE clause and may affect the full table")
		}
	case "DELETE":
		riskLevel = maxSQLRisk(riskLevel, sqlRiskHigh)
		reasons = append(reasons, "DELETE removes data and is treated as high risk by default")
		if !strings.Contains(strings.ToUpper(normalized), " WHERE ") {
			reasons = append(reasons, "DELETE does not include a WHERE clause and may affect the full table")
		}
	case "DROP", "TRUNCATE", "ALTER":
		riskLevel = sqlRiskHigh
		reasons = append(reasons, "DDL changes may affect schema structure or clear data and are treated as high risk")
	case "CREATE":
		riskLevel = maxSQLRisk(riskLevel, sqlRiskMedium)
		reasons = append(reasons, "CREATE changes add new schema objects; verify naming and impact")
	default:
		riskLevel = maxSQLRisk(riskLevel, sqlRiskHigh)
		reasons = append(reasons, "unrecognized change type; treated as high risk by default")
	}

	if strings.Contains(strings.ToUpper(normalized), " LIMIT ") && riskLevel < sqlRiskHigh {
		reasons = append(reasons, "LIMIT clause detected, which may reduce impact scope")
	}

	affectedTables := strings.Join(tables, ", ")
	if affectedTables == "" {
		affectedTables = "-"
	}

	rollbackSQL, rollbackHint := buildSQLRollbackGuidance(opType, tables)
	return riskLevel, strings.Join(uniqueSQLStrings(reasons), "; "), requireApproval, affectedTables, rollbackSQL, rollbackHint
}

func buildSQLRollbackPreview(ctx context.Context, db *sql.DB, order *model.CmdbSQLWorkOrder) (string, int64, string) {
	tableName, whereClause, ok := extractSQLMutationTarget(order.SQLContent)
	if !ok {
		return "", 0, ""
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT 50", tableName, whereClause)
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := db.QueryContext(timeoutCtx, query)
	if err != nil {
		return "", 0, "failed to collect pre-execution snapshot; back up affected data manually before running SQL"
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return "", 0, "failed to collect pre-execution snapshot; back up affected data manually before running SQL"
	}

	list := make([]map[string]interface{}, 0)
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		if err = rows.Scan(pointers...); err != nil {
			break
		}
		item := make(map[string]interface{}, len(columns))
		for i, column := range columns {
			if bytes, ok := values[i].([]byte); ok {
				item[column] = string(bytes)
			} else {
				item[column] = values[i]
			}
		}
		list = append(list, item)
	}
	if len(list) == 0 {
		return "", 0, "no rollback snapshot rows were collected; verify the WHERE clause is correct"
	}
	preview, _ := json.MarshalIndent(list, "", "  ")
	return string(preview), int64(len(list)), "pre-execution snapshot collected; use it to generate rollback SQL or restore data manually"
}

func extractSQLMutationTarget(sqlText string) (string, string, bool) {
	updatePattern := regexp.MustCompile("(?is)^\\s*UPDATE\\s+([`a-zA-Z0-9_\\.]+)\\s+SET\\s+.+?\\s+WHERE\\s+(.+?)\\s*;?\\s*$")
	if matches := updatePattern.FindStringSubmatch(strings.TrimSpace(sqlText)); len(matches) == 3 {
		return matches[1], strings.TrimSpace(matches[2]), true
	}

	deletePattern := regexp.MustCompile("(?is)^\\s*DELETE\\s+FROM\\s+([`a-zA-Z0-9_\\.]+)\\s+WHERE\\s+(.+?)\\s*;?\\s*$")
	if matches := deletePattern.FindStringSubmatch(strings.TrimSpace(sqlText)); len(matches) == 3 {
		return matches[1], strings.TrimSpace(matches[2]), true
	}
	return "", "", false
}

func detectSQLWorkOrderOperationType(sqlText string) string {
	normalized := strings.ToUpper(strings.TrimSpace(sqlText))
	switch {
	case strings.HasPrefix(normalized, "SELECT"):
		return "SELECT"
	case strings.HasPrefix(normalized, "INSERT"):
		return "INSERT"
	case strings.HasPrefix(normalized, "UPDATE"):
		return "UPDATE"
	case strings.HasPrefix(normalized, "DELETE"):
		return "DELETE"
	case strings.HasPrefix(normalized, "CREATE"):
		return "CREATE"
	case strings.HasPrefix(normalized, "ALTER"):
		return "ALTER"
	case strings.HasPrefix(normalized, "DROP"):
		return "DROP"
	case strings.HasPrefix(normalized, "TRUNCATE"):
		return "TRUNCATE"
	default:
		return ""
	}
}

func detectRedisWorkOrderOperationType(command string) string {
	normalized := strings.TrimSpace(strings.ToUpper(command))
	if normalized == "" {
		return ""
	}
	tokens := strings.Fields(normalized)
	if len(tokens) == 0 {
		return ""
	}
	if len(tokens) >= 2 && tokens[0] == "CONFIG" && tokens[1] == "SET" {
		return "CONFIG SET"
	}
	return tokens[0]
}

func extractSQLWorkOrderTables(sqlText string) []string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile("(?is)\\bUPDATE\\s+([`a-zA-Z0-9_\\.]+)"),
		regexp.MustCompile("(?is)\\bINSERT\\s+INTO\\s+([`a-zA-Z0-9_\\.]+)"),
		regexp.MustCompile("(?is)\\bDELETE\\s+FROM\\s+([`a-zA-Z0-9_\\.]+)"),
		regexp.MustCompile("(?is)\\bALTER\\s+TABLE\\s+([`a-zA-Z0-9_\\.]+)"),
		regexp.MustCompile("(?is)\\bDROP\\s+TABLE\\s+([`a-zA-Z0-9_\\.]+)"),
		regexp.MustCompile("(?is)\\bTRUNCATE\\s+TABLE\\s+([`a-zA-Z0-9_\\.]+)"),
		regexp.MustCompile("(?is)\\bCREATE\\s+TABLE\\s+(?:IF\\s+NOT\\s+EXISTS\\s+)?([`a-zA-Z0-9_\\.]+)"),
	}

	tables := make([]string, 0)
	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(sqlText, -1)
		for _, match := range matches {
			if len(match) > 1 {
				tables = append(tables, strings.Trim(match[1], "`"))
			}
		}
	}
	return tables
}

func buildSQLRollbackGuidance(operationType string, tables []string) (string, string) {
	target := firstNonEmptySQL(strings.Join(tables, ", "), "<table>")
	switch operationType {
	case "INSERT":
		return fmt.Sprintf("DELETE FROM %s WHERE <primary-key-or-unique-condition>;", target), "record the inserted primary or unique keys so the new rows can be deleted precisely if rollback is needed"
	case "UPDATE":
		return fmt.Sprintf("-- back up affected rows from %s before execution\nSELECT * FROM %s WHERE <original-WHERE-condition> LIMIT 50;", target, target), "back up rows affected by UPDATE first, then generate restorative UPDATE statements from the snapshot if needed"
	case "DELETE":
		return fmt.Sprintf("-- back up rows to be deleted from %s before execution\nSELECT * FROM %s WHERE <original-WHERE-condition> LIMIT 50;", target, target), "DELETE rollback depends on pre-execution backups; export rows first or generate INSERT restore statements"
	case "ALTER", "DROP", "TRUNCATE", "CREATE":
		return "-- use schema backup, schema diff, or backup database for rollback", "DDL changes should be backed up first and rolled back with backup snapshots or compensating migration scripts"
	default:
		return "-- prepare rollback scripts before execution", "define rollback steps and backup strategy before approval"
	}
}

func analyzeRedisWorkOrder(command string) (int, string, bool, string, string, string) {
	opType := detectRedisWorkOrderOperationType(command)
	kind, blocked := classifyRedisCommand(command)
	reasons := make([]string, 0, 2)

	if blocked {
		reasons = append(reasons, "high-risk Redis command is blocked by policy")
		return sqlRiskHigh, strings.Join(reasons, "; "), true, "-", "-- blocked by policy", "blocked Redis commands cannot be executed; use an alternative recovery plan"
	}

	switch kind {
	case redisCommandKindRead:
		reasons = append(reasons, "read-only Redis command should be executed directly")
		return sqlRiskLow, strings.Join(reasons, "; "), false, "-", "-- direct execution only", "read-only Redis commands do not require a work order"
	case redisCommandKindWrite:
		reasons = append(reasons, "Redis write command changes live cache data and requires approval")
		return sqlRiskMedium, strings.Join(reasons, "; "), true, firstNonEmptySQL(opType, "-"), "-- prepare compensating Redis commands before execution", "record the original key values first so they can be restored manually if needed"
	default:
		reasons = append(reasons, "unrecognized Redis command; treated as high risk")
		return sqlRiskHigh, strings.Join(reasons, "; "), true, "-", "-- review command manually before execution", "unknown Redis commands require manual review and rollback planning"
	}
}

func splitSQLWorkOrderStatements(sqlText string) []string {
	parts := strings.Split(sqlText, ";")
	list := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			list = append(list, trimmed)
		}
	}
	return list
}

func maxSQLRisk(current, next int) int {
	if next > current {
		return next
	}
	return current
}

func uniqueSQLStrings(values []string) []string {
	if len(values) == 0 {
		return values
	}
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func buildSQLWorkOrderTitle(title, operationType, schemaName string) string {
	if strings.TrimSpace(title) != "" {
		return strings.TrimSpace(title)
	}
	return fmt.Sprintf("%s 变更工单 - %s", operationType, schemaName)
}

func generateSQLWorkOrderNo(now time.Time) string {
	return "SQL" + now.Format("20060102150405") + strconv.FormatInt(now.UnixNano()%100000, 10)
}

func firstNonEmptySQL(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
