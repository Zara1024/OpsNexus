CREATE TABLE IF NOT EXISTS cmdb_ssh_records (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL DEFAULT 0,
    host_id BIGINT NOT NULL DEFAULT 0,
    session_id VARCHAR(128) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    client_ip VARCHAR(64) NOT NULL DEFAULT '',
    recording_path VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_cmdb_ssh_records_user_id ON cmdb_ssh_records(user_id);
CREATE INDEX IF NOT EXISTS idx_cmdb_ssh_records_host_id ON cmdb_ssh_records(host_id);
CREATE INDEX IF NOT EXISTS idx_cmdb_ssh_records_session_id ON cmdb_ssh_records(session_id);
