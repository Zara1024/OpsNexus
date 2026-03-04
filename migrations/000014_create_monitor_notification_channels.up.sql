CREATE TABLE IF NOT EXISTS monitor_notification_channels (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    type VARCHAR(32) NOT NULL,
    config_encrypted TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_monitor_notification_channels_name_active ON monitor_notification_channels(name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_monitor_notification_channels_type ON monitor_notification_channels(type);
