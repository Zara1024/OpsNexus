CREATE TABLE IF NOT EXISTS k8s_clusters (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    api_server_url VARCHAR(512) NOT NULL,
    kubeconfig_encrypted TEXT NOT NULL,
    description VARCHAR(512) NOT NULL DEFAULT '',
    version VARCHAR(64) NOT NULL DEFAULT '',
    status SMALLINT NOT NULL DEFAULT 1,
    node_count INTEGER NOT NULL DEFAULT 0,
    pod_count INTEGER NOT NULL DEFAULT 0,
    created_by BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_k8s_clusters_name_active ON k8s_clusters(name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_k8s_clusters_status ON k8s_clusters(status);

CREATE TABLE IF NOT EXISTS k8s_cluster_metrics (
    id BIGSERIAL PRIMARY KEY,
    cluster_id BIGINT NOT NULL,
    cpu_usage NUMERIC(8,2) NOT NULL DEFAULT 0,
    mem_usage NUMERIC(8,2) NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_k8s_cluster_metrics_cluster FOREIGN KEY (cluster_id) REFERENCES k8s_clusters(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_k8s_cluster_metrics_cluster_active ON k8s_cluster_metrics(cluster_id) WHERE deleted_at IS NULL;
