package db

import (
	"strings"
	"testing"
)

func TestBuildMySQLDSNForcesUTF8MB4WhenLegacyCharsetConfigured(t *testing.T) {
	dsn := buildMySQLDSN("root", "secret", "127.0.0.1", 3306, "autoops", "utf8")

	if !strings.Contains(dsn, "charset=utf8mb4") {
		t.Fatalf("expected dsn to use utf8mb4 charset, got %q", dsn)
	}
}

func TestBuildMySQLDSNDefaultsToUTF8MB4WhenCharsetMissing(t *testing.T) {
	dsn := buildMySQLDSN("root", "secret", "127.0.0.1", 3306, "autoops", "")

	if !strings.Contains(dsn, "charset=utf8mb4") {
		t.Fatalf("expected dsn to default to utf8mb4 charset, got %q", dsn)
	}
}
