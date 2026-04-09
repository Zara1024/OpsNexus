package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"dodevops-api/api/cmdb/model"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

type databaseExecutionType string

const (
	databaseExecutionTypeMySQL      databaseExecutionType = "mysql"
	databaseExecutionTypePostgreSQL databaseExecutionType = "postgresql"
	databaseExecutionTypeRedis      databaseExecutionType = "redis"
)

type redisCommandKind string

const (
	redisCommandKindUnknown redisCommandKind = "unknown"
	redisCommandKindRead    redisCommandKind = "read"
	redisCommandKindWrite   redisCommandKind = "write"
	redisCommandKindBlocked redisCommandKind = "blocked"
)

var (
	errUnsupportedDatabaseType = errors.New("unsupported database type")

	redisBlockedCommands = map[string]struct{}{
		"FLUSHALL":   {},
		"FLUSHDB":    {},
		"SHUTDOWN":   {},
		"CONFIG SET": {},
	}
	redisReadCommands = map[string]struct{}{
		"EXISTS":   {},
		"GET":      {},
		"HGET":     {},
		"HGETALL":  {},
		"KEYS":     {},
		"LRANGE":   {},
		"MGET":     {},
		"SCARD":    {},
		"SMEMBERS": {},
		"STRLEN":   {},
		"TTL":      {},
		"TYPE":     {},
		"ZRANGE":   {},
	}
)

func resolveDatabaseExecutionType(dbType int) databaseExecutionType {
	switch dbType {
	case 2:
		return databaseExecutionTypePostgreSQL
	case 3:
		return databaseExecutionTypeRedis
	default:
		return databaseExecutionTypeMySQL
	}
}

func ResolveDatabaseExecutionType(dbType int) string {
	return string(resolveDatabaseExecutionType(dbType))
}

func classifyRedisCommand(command string) (redisCommandKind, bool) {
	normalized := strings.ToUpper(strings.TrimSpace(command))
	if normalized == "" {
		return redisCommandKindUnknown, false
	}

	tokens := strings.Fields(normalized)
	if len(tokens) == 0 {
		return redisCommandKindUnknown, false
	}

	if len(tokens) >= 2 {
		composite := tokens[0] + " " + tokens[1]
		if _, exists := redisBlockedCommands[composite]; exists {
			return redisCommandKindBlocked, true
		}
	}

	if _, exists := redisBlockedCommands[tokens[0]]; exists {
		return redisCommandKindBlocked, true
	}
	if _, exists := redisReadCommands[tokens[0]]; exists {
		return redisCommandKindRead, false
	}
	return redisCommandKindWrite, false
}

func ClassifyRedisCommand(command string) (string, bool) {
	kind, blocked := classifyRedisCommand(command)
	return string(kind), blocked
}

func openPostgreSQLConnection(ctx context.Context, host string, port int, username, password, databaseName string) (*sql.DB, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database server %s:%d: %v", host, port, err)
	}
	_ = conn.Close()

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=10",
		host, port, username, password, databaseName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err = db.PingContext(timeoutCtx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func OpenPostgreSQLConnection(ctx context.Context, host string, port int, username, password, databaseName string) (*sql.DB, error) {
	return openPostgreSQLConnection(ctx, host, port, username, password, databaseName)
}

func buildPostgreSQLDatabaseListQuery() string {
	return "SELECT datname FROM pg_database WHERE datistemplate = false ORDER BY datname"
}

func BuildPostgreSQLDatabaseListQuery() string {
	return buildPostgreSQLDatabaseListQuery()
}

func buildRedisClient(host string, port int, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	})
}

func BuildRedisClient(host string, port int, password string, db int) *redis.Client {
	return buildRedisClient(host, port, password, db)
}

func getDatabaseSchemaTarget(target *model.CmdbSQL, requestedName string) string {
	if target == nil {
		return strings.TrimSpace(requestedName)
	}
	return ResolveCmdbSQLSchemaName(*target, requestedName)
}

func GetDatabaseSchemaTarget(target *model.CmdbSQL, requestedName string) string {
	return getDatabaseSchemaTarget(target, requestedName)
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

func ParseRedisDatabaseIndex(name string) int {
	return parseRedisDatabaseIndex(name)
}
