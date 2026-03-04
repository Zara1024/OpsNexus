CREATE TABLE IF NOT EXISTS monitor_alert_rules (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    description VARCHAR(512) NOT NULL DEFAULT '',
    severity VARCHAR(8) NOT NULL DEFAULT 'P3',
    expression TEXT NOT NULL,
    duration VARCHAR(64) NOT NULL DEFAULT '5m',
    labels JSONB NOT NULL DEFAULT '{}'::jsonb,
    annotations JSONB NOT NULL DEFAULT '{}'::jsonb,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    notify_channels JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_monitor_alert_rules_name_active ON monitor_alert_rules(name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_monitor_alert_rules_severity ON monitor_alert_rules(severity);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_rules_enabled ON monitor_alert_rules(enabled);
