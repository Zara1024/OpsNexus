CREATE TABLE IF NOT EXISTS cmdb_hosts (
    id BIGSERIAL PRIMARY KEY,
    hostname VARCHAR(128) NOT NULL,
    inner_ip VARCHAR(64) NOT NULL DEFAULT '',
    outer_ip VARCHAR(64) NOT NULL DEFAULT '',
    os_type VARCHAR(64) NOT NULL DEFAULT '',
    os_version VARCHAR(128) NOT NULL DEFAULT '',
    arch VARCHAR(64) NOT NULL DEFAULT '',
    cpu_cores INT NOT NULL DEFAULT 0,
    memory_gb NUMERIC(10,2) NOT NULL DEFAULT 0,
    disk_gb NUMERIC(10,2) NOT NULL DEFAULT 0,
    group_id BIGINT NOT NULL DEFAULT 0,
    credential_id BIGINT NOT NULL DEFAULT 0,
    cloud_provider VARCHAR(64) NOT NULL DEFAULT '',
    instance_id VARCHAR(128) NOT NULL DEFAULT '',
    region VARCHAR(64) NOT NULL DEFAULT '',
    agent_status VARCHAR(32) NOT NULL DEFAULT 'unknown',
    labels JSONB NOT NULL DEFAULT '{}'::jsonb,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_cmdb_hosts_group_id ON cmdb_hosts(group_id);
CREATE INDEX IF NOT EXISTS idx_cmdb_hosts_credential_id ON cmdb_hosts(credential_id);
CREATE INDEX IF NOT EXISTS idx_cmdb_hosts_inner_ip ON cmdb_hosts(inner_ip);
CREATE INDEX IF NOT EXISTS idx_cmdb_hosts_outer_ip ON cmdb_hosts(outer_ip);
