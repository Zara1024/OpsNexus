package service

import "testing"

func TestResolveDatabaseExecutionType(t *testing.T) {
	if got := resolveDatabaseExecutionType(1); got != databaseExecutionTypeMySQL {
		t.Fatalf("expected mysql execution type, got %q", got)
	}
	if got := resolveDatabaseExecutionType(2); got != databaseExecutionTypePostgreSQL {
		t.Fatalf("expected postgresql execution type, got %q", got)
	}
	if got := resolveDatabaseExecutionType(3); got != databaseExecutionTypeRedis {
		t.Fatalf("expected redis execution type, got %q", got)
	}
}

func TestClassifyRedisCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    redisCommandKind
		blocked bool
	}{
		{name: "read", command: "GET demo", want: redisCommandKindRead, blocked: false},
		{name: "write", command: "SET demo value", want: redisCommandKindWrite, blocked: false},
		{name: "blocked flushall", command: "FLUSHALL", want: redisCommandKindBlocked, blocked: true},
		{name: "blocked config set", command: "CONFIG SET requirepass secret", want: redisCommandKindBlocked, blocked: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, blocked := classifyRedisCommand(tt.command)
			if got != tt.want {
				t.Fatalf("expected kind %q, got %q", tt.want, got)
			}
			if blocked != tt.blocked {
				t.Fatalf("expected blocked=%v, got %v", tt.blocked, blocked)
			}
		})
	}
}
