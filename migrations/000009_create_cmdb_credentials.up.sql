CREATE TABLE IF NOT EXISTS cmdb_credentials (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    type VARCHAR(16) NOT NULL,
    username VARCHAR(128) NOT NULL DEFAULT '',
    password_encrypted TEXT NOT NULL DEFAULT '',
    private_key_encrypted TEXT NOT NULL DEFAULT '',
    passphrase_encrypted TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
