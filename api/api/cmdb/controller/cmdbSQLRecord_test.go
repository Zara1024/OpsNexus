package controller

import "testing"

func TestValidateSQLTypeAllowsReadOnlyQueriesWhenSelectWhitelisted(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want bool
	}{
		{name: "select", sql: "SELECT * FROM users", want: true},
		{name: "show tables", sql: "SHOW TABLES", want: true},
		{name: "show create table", sql: "SHOW CREATE TABLE users", want: true},
		{name: "desc", sql: "DESC users", want: true},
		{name: "describe", sql: "DESCRIBE users", want: true},
		{name: "explain", sql: "EXPLAIN SELECT * FROM users", want: true},
		{name: "delete", sql: "DELETE FROM users", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateSQLType(tt.sql, []string{"SELECT"})
			if got != tt.want {
				t.Fatalf("validateSQLType(%q, [SELECT]) = %v, want %v", tt.sql, got, tt.want)
			}
		})
	}
}

func TestDetectOperationTypeRecognizesReadOnlyStatements(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want string
	}{
		{name: "select", sql: "SELECT * FROM users", want: "SELECT"},
		{name: "show", sql: "SHOW TABLES", want: "SHOW"},
		{name: "desc", sql: "DESC users", want: "DESCRIBE"},
		{name: "describe", sql: "DESCRIBE users", want: "DESCRIBE"},
		{name: "explain", sql: "EXPLAIN SELECT * FROM users", want: "EXPLAIN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectOperationType(tt.sql)
			if got != tt.want {
				t.Fatalf("detectOperationType(%q) = %q, want %q", tt.sql, got, tt.want)
			}
		})
	}
}

func TestNormalizeReadOnlyQueryForPostgreSQL(t *testing.T) {
	got := normalizeReadOnlyQueryForDatabaseType(2, "show database")
	want := "SELECT datname AS database_name FROM pg_database WHERE datistemplate = false ORDER BY datname"
	if got != want {
		t.Fatalf("normalizeReadOnlyQueryForDatabaseType(postgresql, show database) = %q, want %q", got, want)
	}

	got = normalizeReadOnlyQueryForDatabaseType(2, "SHOW DATABASES")
	if got != want {
		t.Fatalf("normalizeReadOnlyQueryForDatabaseType(postgresql, SHOW DATABASES) = %q, want %q", got, want)
	}

	original := "SELECT current_database()"
	got = normalizeReadOnlyQueryForDatabaseType(2, original)
	if got != original {
		t.Fatalf("expected non-mysql-style PostgreSQL query to stay unchanged, got %q", got)
	}
}
