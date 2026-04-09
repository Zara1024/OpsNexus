package controller

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"

	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	cmdbService "dodevops-api/api/cmdb/service"
	configModel "dodevops-api/api/configcenter/model"
	configService "dodevops-api/api/configcenter/service"
	sysModel "dodevops-api/api/system/model"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"dodevops-api/common/util"

	"github.com/gin-gonic/gin"
	redisv8 "github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

const (
	ParamError    = 1001
	DatabaseError = 1002
)

type CmdbSQLRecordController struct {
	recordService *cmdbService.CmdbSQLRecordService
	sqlService    *cmdbService.CmdbSQLService
}

var sqlRecordController *CmdbSQLRecordController

// SQLRequest matches the front-end SQL execution payload.
type SQLRequest struct {
	DatabaseID   uint   `json:"databaseId"`
	DatabaseName string `json:"databaseName"`
	SQL          string `json:"sql"`
}

func InitCmdbSQLRecordController(db *gorm.DB) {
	recordDao := dao.NewCmdbSQLRecordDao(db)
	recordService := cmdbService.NewCmdbSQLRecordService(recordDao)
	sqlService := cmdbService.NewCmdbSQLService(dao.NewCmdbSQLDao(db))
	sqlRecordController = &CmdbSQLRecordController{
		recordService: recordService,
		sqlService:    sqlService,
	}
}

func GetCmdbSQLRecordController() *CmdbSQLRecordController {
	return sqlRecordController
}

func NewCmdbSQLRecordController(db *gorm.DB) *CmdbSQLRecordController {
	recordDao := dao.NewCmdbSQLRecordDao(db)
	recordService := cmdbService.NewCmdbSQLRecordService(recordDao)
	sqlService := cmdbService.NewCmdbSQLService(dao.NewCmdbSQLDao(db))
	return &CmdbSQLRecordController{
		recordService: recordService,
		sqlService:    sqlService,
	}
}

func (c *CmdbSQLRecordController) getDBConnectionInfo(req SQLRequest) (*model.CmdbSQL, *configModel.AccountAuth, string, string, error) {
	dbInfo, schemaName, err := c.sqlService.ResolveDatabaseTarget(req.DatabaseID, req.DatabaseName)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to resolve database asset: %w", err)
	}

	account, err := configService.NewAccountAuthService().GetByID(dbInfo.AccountID)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to load database account: %w", err)
	}

	decrypted, err := configService.NewAccountAuthService().DecryptPassword(dbInfo.AccountID)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to decrypt database password: %w", err)
	}

	return dbInfo, account, decrypted, schemaName, nil
}

func validateSQLType(sqlText string, allowedTypes []string) bool {
	sqlText = strings.TrimSpace(strings.ToUpper(sqlText))
	for _, item := range allowedTypes {
		if strings.HasPrefix(sqlText, item) {
			return true
		}
	}
	if containsSQLType(allowedTypes, "SELECT") {
		for _, item := range []string{"SHOW", "DESC", "DESCRIBE", "EXPLAIN"} {
			if strings.HasPrefix(sqlText, item) {
				return true
			}
		}
	}
	return false
}

func containsSQLType(allowedTypes []string, target string) bool {
	for _, item := range allowedTypes {
		if item == target {
			return true
		}
	}
	return false
}

func getCurrentUsername(ctx *gin.Context) (string, error) {
	userObj, exists := ctx.Get(constant.ContextKeyUserObj)
	if !exists {
		return "", fmt.Errorf("failed to get current user")
	}

	switch value := userObj.(type) {
	case map[string]interface{}:
		if username, ok := value["username"].(string); ok && username != "" {
			return username, nil
		}
	case *sysModel.SysAdmin:
		if value.Username != "" {
			return value.Username, nil
		}
	case sysModel.SysAdmin:
		if value.Username != "" {
			return value.Username, nil
		}
	default:
		val := reflect.ValueOf(userObj)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() == reflect.Struct {
			field := val.FieldByName("Username")
			if field.IsValid() && field.Kind() == reflect.String {
				return field.String(), nil
			}
		}
	}

	return "", fmt.Errorf("failed to resolve current username")
}

func openMySQLConnection(ctx *gin.Context, host string, port int, username, password, databaseName string) (*sql.DB, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database server %s:%d: %v", host, port, err)
	}
	conn.Close()

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

	ctxPing, cancelPing := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancelPing()
	if err = db.PingContext(ctxPing); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func scanRowsToMaps(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	resultsData := make([]map[string]interface{}, 0)
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err = rows.Scan(pointers...); err != nil {
			return nil, err
		}

		rowData := make(map[string]interface{}, len(columns))
		for i, col := range columns {
			if b, ok := values[i].([]byte); ok {
				rowData[col] = string(b)
			} else {
				rowData[col] = values[i]
			}
		}
		resultsData = append(resultsData, rowData)
	}
	return resultsData, nil
}

func parseRedisDatabaseIndex(name string) int {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return 0
	}
	index, err := strconv.Atoi(trimmed)
	if err != nil || index < 0 {
		return 0
	}
	return index
}

func normalizeReadOnlyQueryForDatabaseType(dbType int, sqlText string) string {
	trimmed := strings.TrimSpace(sqlText)
	if dbType != 2 {
		return trimmed
	}

	normalized := strings.ToUpper(trimmed)
	switch normalized {
	case "SHOW DATABASE", "SHOW DATABASES":
		return "SELECT datname AS database_name FROM pg_database WHERE datistemplate = false ORDER BY datname"
	default:
		return trimmed
	}
}

func executeRedisCommand(ctx *gin.Context, host string, port int, password, command string, dbIndex int) (interface{}, error) {
	client := cmdbService.BuildRedisClient(host, port, password, dbIndex)
	defer client.Close()

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	if err := client.Ping(timeoutCtx).Err(); err != nil {
		return nil, err
	}

	args := strings.Fields(strings.TrimSpace(command))
	if len(args) == 0 {
		return nil, fmt.Errorf("redis command cannot be empty")
	}

	commandName := strings.ToUpper(args[0])
	commandArgs := make([]interface{}, 0, len(args))
	commandArgs = append(commandArgs, commandName)
	for _, arg := range args[1:] {
		commandArgs = append(commandArgs, arg)
	}

	resultValue, err := client.Do(timeoutCtx, commandArgs...).Result()
	if err != nil {
		if err == redisv8.Nil {
			resultValue = nil
		} else {
			return nil, err
		}
	}

	return gin.H{
		"command":  commandName,
		"database": dbIndex,
		"result":   resultValue,
	}, nil
}

func (c *CmdbSQLRecordController) executeSelectRequest(ctx *gin.Context, req SQLRequest) {
	dbInfo, account, decrypted, schemaName, err := c.getDBConnectionInfo(req)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, err.Error())
		return
	}
	var (
		resultsData   []map[string]interface{}
		returnedRows  int64
		executionTime int64
		operationType string
	)
	normalizedSQL := normalizeReadOnlyQueryForDatabaseType(dbInfo.Type, req.SQL)

	switch dbInfo.Type {
	case 1:
		rootConnStr := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/?timeout=60s&readTimeout=60s&writeTimeout=60s&parseTime=true&interpolateParams=true",
			account.Name, decrypted, account.Host, account.Port,
		)
		rootDB, err := sql.Open("mysql", rootConnStr)
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "failed to initialize database connection: "+err.Error())
			return
		}
		defer rootDB.Close()

		ctxCheckSchema, cancelCheckSchema := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
		defer cancelCheckSchema()

		var schemaExists int
		err = rootDB.QueryRowContext(
			ctxCheckSchema,
			"SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?",
			schemaName,
		).Scan(&schemaExists)
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "failed to check schema existence: "+err.Error())
			return
		}
		if schemaExists == 0 {
			result.FailedWithCode(ctx, DatabaseError, fmt.Sprintf("schema %q does not exist", schemaName))
			return
		}

		db, err := openMySQLConnection(ctx, account.Host, account.Port, account.Name, decrypted, schemaName)
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "failed to connect to database: "+err.Error())
			return
		}
		defer db.Close()

		queryCtx, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
		defer cancel()

		startTime := time.Now()
		rows, err := db.QueryContext(queryCtx, normalizedSQL)
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "query execution failed: "+err.Error())
			return
		}
		defer rows.Close()

		resultsData, err = scanRowsToMaps(rows)
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "failed to scan row data: "+err.Error())
			return
		}
		executionTime = time.Since(startTime).Milliseconds()
		returnedRows = int64(len(resultsData))
		operationType = detectOperationType(normalizedSQL)
	case 2:
		db, err := cmdbService.OpenPostgreSQLConnection(ctx.Request.Context(), account.Host, account.Port, account.Name, decrypted, schemaName)
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "failed to connect to database: "+err.Error())
			return
		}
		defer db.Close()

		queryCtx, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
		defer cancel()

		startTime := time.Now()
		rows, err := db.QueryContext(queryCtx, normalizedSQL)
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "query execution failed: "+err.Error())
			return
		}
		defer rows.Close()

		resultsData, err = scanRowsToMaps(rows)
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "failed to scan row data: "+err.Error())
			return
		}
		executionTime = time.Since(startTime).Milliseconds()
		returnedRows = int64(len(resultsData))
		operationType = detectOperationType(normalizedSQL)
	case 3:
		kind, blocked := cmdbService.ClassifyRedisCommand(normalizedSQL)
		if blocked {
			result.FailedWithCode(ctx, DatabaseError, "redis high-risk commands are blocked by policy")
			return
		}
		startTime := time.Now()
		payload, err := executeRedisCommand(ctx, account.Host, account.Port, decrypted, normalizedSQL, parseRedisDatabaseIndex(schemaName))
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "command execution failed: "+err.Error())
			return
		}
		executionTime = time.Since(startTime).Milliseconds()
		resultsData = []map[string]interface{}{{"payload": payload}}
		returnedRows = 1
		operationType = strings.ToUpper(kind)
	default:
		result.FailedWithCode(ctx, DatabaseError, "unsupported database type")
		return
	}

	username, err := getCurrentUsername(ctx)
	if err != nil {
		result.FailedWithCode(ctx, ParamError, err.Error())
		return
	}

	if operationType == "" {
		operationType = "SELECT"
	}

	err = c.recordService.RecordSQLExecution(
		account.Host,
		schemaName,
		operationType,
		normalizedSQL,
		username,
		util.GetClientIP(ctx.Request),
		0,
		0,
		executionTime,
		returnedRows,
		"SUCCESS",
	)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "database error: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"returnedRows":  returnedRows,
		"executionTime": executionTime,
		"results":       resultsData,
	})
}

func (c *CmdbSQLRecordController) executeMutationRequest(ctx *gin.Context, req SQLRequest, operationType string) {
	dbInfo, account, decrypted, schemaName, err := c.getDBConnectionInfo(req)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, err.Error())
		return
	}
	var (
		affectedRows  int64
		executionTime int64
	)

	switch dbInfo.Type {
	case 1:
		db, err := openMySQLConnection(ctx, account.Host, account.Port, account.Name, decrypted, schemaName)
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "failed to connect to database: "+err.Error())
			return
		}
		defer db.Close()

		queryCtx, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
		defer cancel()

		startTime := time.Now()
		execResult, err := db.ExecContext(queryCtx, req.SQL)
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, operationType+" execution failed: "+err.Error())
			return
		}

		affectedRows, err = execResult.RowsAffected()
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "failed to get affected rows: "+err.Error())
			return
		}
		executionTime = time.Since(startTime).Milliseconds()
	case 2:
		db, err := cmdbService.OpenPostgreSQLConnection(ctx.Request.Context(), account.Host, account.Port, account.Name, decrypted, schemaName)
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "failed to connect to database: "+err.Error())
			return
		}
		defer db.Close()

		queryCtx, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
		defer cancel()

		startTime := time.Now()
		execResult, err := db.ExecContext(queryCtx, req.SQL)
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, operationType+" execution failed: "+err.Error())
			return
		}

		affectedRows, err = execResult.RowsAffected()
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "failed to get affected rows: "+err.Error())
			return
		}
		executionTime = time.Since(startTime).Milliseconds()
	case 3:
		kind, blocked := cmdbService.ClassifyRedisCommand(req.SQL)
		if blocked {
			result.FailedWithCode(ctx, DatabaseError, "redis high-risk commands are blocked by policy")
			return
		}
		if kind == "read" {
			result.FailedWithCode(ctx, DatabaseError, "redis read commands should be executed through the direct execution path")
			return
		}
		startTime := time.Now()
		_, err := executeRedisCommand(ctx, account.Host, account.Port, decrypted, req.SQL, parseRedisDatabaseIndex(schemaName))
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, operationType+" execution failed: "+err.Error())
			return
		}
		affectedRows = 1
		executionTime = time.Since(startTime).Milliseconds()
	default:
		result.FailedWithCode(ctx, DatabaseError, "unsupported database type")
		return
	}

	username, err := getCurrentUsername(ctx)
	if err != nil {
		result.FailedWithCode(ctx, ParamError, err.Error())
		return
	}

	err = c.recordService.RecordSQLExecution(
		account.Host,
		schemaName,
		operationType,
		req.SQL,
		username,
		util.GetClientIP(ctx.Request),
		0,
		affectedRows,
		executionTime,
		0,
		"SUCCESS",
	)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, "database error: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"affectedRows":  affectedRows,
		"executionTime": executionTime,
		"operationType": operationType,
	})
}

func detectOperationType(sqlText string) string {
	sqlText = strings.TrimSpace(strings.ToUpper(sqlText))
	switch {
	case strings.HasPrefix(sqlText, "SELECT"):
		return "SELECT"
	case strings.HasPrefix(sqlText, "SHOW"):
		return "SHOW"
	case strings.HasPrefix(sqlText, "DESC"), strings.HasPrefix(sqlText, "DESCRIBE"):
		return "DESCRIBE"
	case strings.HasPrefix(sqlText, "EXPLAIN"):
		return "EXPLAIN"
	case strings.HasPrefix(sqlText, "INSERT"):
		return "INSERT"
	case strings.HasPrefix(sqlText, "UPDATE"):
		return "UPDATE"
	case strings.HasPrefix(sqlText, "DELETE"):
		return "DELETE"
	case strings.HasPrefix(sqlText, "CREATE"):
		return "CREATE"
	case strings.HasPrefix(sqlText, "ALTER"):
		return "ALTER"
	case strings.HasPrefix(sqlText, "DROP"):
		return "DROP"
	default:
		return ""
	}
}

// ExecuteSelect handles read-only SQL execution.
func (c *CmdbSQLRecordController) ExecuteSelect(ctx *gin.Context) {
	var req SQLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.FailedWithCode(ctx, ParamError, "invalid params: "+err.Error())
		return
	}
	if !validateSQLType(req.SQL, []string{"SELECT"}) {
		result.FailedWithCode(ctx, ParamError, "only read-only query statements are allowed")
		return
	}
	c.executeSelectRequest(ctx, req)
}

// ExecuteInsert handles INSERT SQL execution.
func (c *CmdbSQLRecordController) ExecuteInsert(ctx *gin.Context) {
	var req SQLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.FailedWithCode(ctx, ParamError, "invalid params: "+err.Error())
		return
	}
	if !validateSQLType(req.SQL, []string{"INSERT"}) {
		result.FailedWithCode(ctx, ParamError, "only INSERT statements are allowed")
		return
	}
	c.executeMutationRequest(ctx, req, "INSERT")
}

// ExecuteUpdate handles UPDATE SQL execution.
func (c *CmdbSQLRecordController) ExecuteUpdate(ctx *gin.Context) {
	var req SQLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.FailedWithCode(ctx, ParamError, "invalid params: "+err.Error())
		return
	}
	if !validateSQLType(req.SQL, []string{"UPDATE"}) {
		result.FailedWithCode(ctx, ParamError, "only UPDATE statements are allowed")
		return
	}
	c.executeMutationRequest(ctx, req, "UPDATE")
}

// ExecuteSQL handles raw SQL execution after operation type dispatch.
func (c *CmdbSQLRecordController) ExecuteSQL(ctx *gin.Context) {
	var req SQLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.FailedWithCode(ctx, ParamError, "invalid params: "+err.Error())
		return
	}

	dbInfo, _, _, _, err := c.getDBConnectionInfo(req)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, err.Error())
		return
	}

	if dbInfo.Type == 3 {
		kind, blocked := cmdbService.ClassifyRedisCommand(req.SQL)
		if blocked {
			result.FailedWithCode(ctx, DatabaseError, "redis high-risk commands are blocked by policy")
			return
		}
		if kind == "read" {
			c.executeSelectRequest(ctx, req)
			return
		}
		c.executeMutationRequest(ctx, req, "REDIS")
		return
	}

	allowedTypes := []string{"SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "ALTER", "DROP"}
	if !validateSQLType(req.SQL, allowedTypes) {
		result.FailedWithCode(ctx, ParamError, "this SQL statement type is not allowed")
		return
	}

	operationType := detectOperationType(req.SQL)
	switch operationType {
	case "SELECT", "SHOW", "DESCRIBE", "EXPLAIN":
		c.executeSelectRequest(ctx, req)
	case "INSERT", "UPDATE", "DELETE", "CREATE", "ALTER", "DROP":
		c.executeMutationRequest(ctx, req, operationType)
	default:
		result.FailedWithCode(ctx, ParamError, "this SQL statement type is not supported")
	}
}

// GetDatabaseList returns schema list from the target DB account.
func (c *CmdbSQLRecordController) GetDatabaseList(ctx *gin.Context) {
	var req SQLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.FailedWithCode(ctx, ParamError, "invalid params: "+err.Error())
		return
	}
	if req.DatabaseID == 0 {
		result.FailedWithCode(ctx, ParamError, "databaseId is required")
		return
	}

	dbInfo, account, decrypted, _, err := c.getDBConnectionInfo(req)
	if err != nil {
		result.FailedWithCode(ctx, DatabaseError, err.Error())
		return
	}
	var databases []string
	switch dbInfo.Type {
	case 1:
		connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/?timeout=60s&readTimeout=60s&writeTimeout=60s&parseTime=true&interpolateParams=true",
			account.Name, decrypted, account.Host, account.Port)
		db, err := sql.Open("mysql", connStr)
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "failed to initialize database connection: "+err.Error())
			return
		}
		defer db.Close()

		rows, err := db.Query("SHOW DATABASES")
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "failed to query database list: "+err.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			var dbName string
			if err = rows.Scan(&dbName); err != nil {
				result.FailedWithCode(ctx, DatabaseError, "failed to parse database list: "+err.Error())
				return
			}
			databases = append(databases, dbName)
		}
	case 2:
		db, err := cmdbService.OpenPostgreSQLConnection(ctx.Request.Context(), account.Host, account.Port, account.Name, decrypted, "postgres")
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "failed to initialize database connection: "+err.Error())
			return
		}
		defer db.Close()

		rows, err := db.Query(cmdbService.BuildPostgreSQLDatabaseListQuery())
		if err != nil {
			result.FailedWithCode(ctx, DatabaseError, "failed to query database list: "+err.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			var dbName string
			if err = rows.Scan(&dbName); err != nil {
				result.FailedWithCode(ctx, DatabaseError, "failed to parse database list: "+err.Error())
				return
			}
			databases = append(databases, dbName)
		}
	case 3:
		client := cmdbService.BuildRedisClient(account.Host, account.Port, decrypted, 0)
		defer client.Close()
		timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
		defer cancel()
		if err := client.Ping(timeoutCtx).Err(); err != nil {
			result.FailedWithCode(ctx, DatabaseError, "failed to initialize redis connection: "+err.Error())
			return
		}
		databases = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15"}
	default:
		result.FailedWithCode(ctx, DatabaseError, "unsupported database type")
		return
	}

	result.Success(ctx, gin.H{
		"databases": databases,
		"host":      account.Host,
		"port":      account.Port,
	})
}

func (c *CmdbSQLRecordController) ListDatabases(ctx *gin.Context) {
	c.GetDatabaseList(ctx)
}

// ExecuteDelete handles DELETE SQL execution.
func (c *CmdbSQLRecordController) ExecuteDelete(ctx *gin.Context) {
	var req SQLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.FailedWithCode(ctx, ParamError, "invalid params: "+err.Error())
		return
	}
	if !validateSQLType(req.SQL, []string{"DELETE"}) {
		result.FailedWithCode(ctx, ParamError, "only DELETE statements are allowed")
		return
	}
	c.executeMutationRequest(ctx, req, "DELETE")
}
