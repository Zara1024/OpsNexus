CREATE TABLE IF NOT EXISTS monitor_alert_events (
    id BIGSERIAL PRIMARY KEY,
    rule_id BIGINT NOT NULL,
    fingerprint VARCHAR(128) NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'firing',
    severity VARCHAR(8) NOT NULL DEFAULT 'P3',
    starts_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ends_at TIMESTAMPTZ,
    labels JSONB NOT NULL DEFAULT '{}'::jsonb,
    annotations JSONB NOT NULL DEFAULT '{}'::jsonb,
    acknowledged_by BIGINT NOT NULL DEFAULT 0,
    silenced_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_monitor_alert_events_rule FOREIGN KEY (rule_id) REFERENCES monitor_alert_rules(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_monitor_alert_events_fingerprint_active ON monitor_alert_events(fingerprint) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_monitor_alert_events_status ON monitor_alert_events(status);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_events_severity ON monitor_alert_events(severity);
CREATE INDEX IF NOT EXISTS idx_monitor_alert_events_starts_at ON monitor_alert_events(starts_at DESC);
