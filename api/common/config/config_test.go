package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigOverridesFromEnv(t *testing.T) {
	t.Setenv("SERVER_ADDRESS", "0.0.0.0:18000")
	t.Setenv("DB_HOST", "mysql-service")
	t.Setenv("DB_PORT", "13306")
	t.Setenv("DB_PASSWORD", "override-db-password")
	t.Setenv("REDIS_ADDR", "redis-service:6379")
	t.Setenv("REDIS_PASSWORD", "override-redis-password")
	t.Setenv("IMAGE_HOST", "http://opsnexus.test")
	t.Setenv("MONITOR_PROMETHEUS_URL", "http://prometheus:9090")
	t.Setenv("MONITOR_AGENT_HEARTBEAT_TOKEN", "override-agent-token")
	t.Setenv("MONITOR_WEBHOOK_TOKEN", "override-webhook-token")
	t.Setenv("SERVER_ENABLE_SWAGGER", "false")

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  address: 127.0.0.1:8000
  model: debug
  enableSwagger: true
db:
  dialects: mysql
  host: 127.0.0.1
  port: 3306
  db: autoops
  username: root
  password: original-db-password
  charset: utf8mb4
  maxIdle: 10
  maxOpen: 20
redis:
  address: 127.0.0.1:6379
  password: original-redis-password
imageSettings:
  uploadDir: ./upload
  imageHost: http://localhost:8080
log:
  path: ./log
  name: sys
  model: console
monitor:
  prometheus:
    url: http://127.0.0.1:9090
  pushgateway:
    url: http://127.0.0.1:9091
  agent:
    heartbeat_server_url: http://127.0.0.1:8000/api/v1/monitor/agent/heartbeat
    heartbeat_token: original-agent-token
  webhook:
    token: original-webhook-token
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}

	Config = nil
	if err := LoadConfig(configPath); err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if Config.Server.Address != "0.0.0.0:18000" {
		t.Fatalf("expected server address override, got %q", Config.Server.Address)
	}
	if Config.Server.EnableSwagger {
		t.Fatalf("expected swagger override to false")
	}
	if Config.Db.Host != "mysql-service" {
		t.Fatalf("expected db host override, got %q", Config.Db.Host)
	}
	if Config.Db.Port != 13306 {
		t.Fatalf("expected db port override, got %d", Config.Db.Port)
	}
	if Config.Db.Password != "override-db-password" {
		t.Fatalf("expected db password override, got %q", Config.Db.Password)
	}
	if Config.Redis.Address != "redis-service:6379" {
		t.Fatalf("expected redis address override, got %q", Config.Redis.Address)
	}
	if Config.Redis.Password != "override-redis-password" {
		t.Fatalf("expected redis password override, got %q", Config.Redis.Password)
	}
	if Config.ImageSettings.ImageHost != "http://opsnexus.test" {
		t.Fatalf("expected image host override, got %q", Config.ImageSettings.ImageHost)
	}
	if Config.Monitor.Prometheus.URL != "http://prometheus:9090" {
		t.Fatalf("expected prometheus url override, got %q", Config.Monitor.Prometheus.URL)
	}
	if Config.Monitor.Agent.HeartbeatToken != "override-agent-token" {
		t.Fatalf("expected heartbeat token override, got %q", Config.Monitor.Agent.HeartbeatToken)
	}
	if Config.Monitor.Webhook.Token != "override-webhook-token" {
		t.Fatalf("expected webhook token override, got %q", Config.Monitor.Webhook.Token)
	}
}
