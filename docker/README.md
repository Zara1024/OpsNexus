# Docker 部署说明

本文档同步日期：`2026-03-20`。

当前 `docker/` 目录提供的是 OpsNexus 的 Docker Compose 最小部署方案，并已同步到 P1 状态。

## 目录结构

```text
docker/
├── docker-compose.yml
├── .env
├── api/
│   └── config.yaml
├── web/
│   └── devops.conf
├── mysql/
├── redis/
├── prometheus/
├── pushgateway/
└── alertmanager/
    └── alertmanager.yml
```

## 服务说明

默认包含：

- `mysql`
- `redis`
- `pushgateway`
- `prometheus`
- `devops-api`
- `devops-web`

可选 P1 告警组件：

- `alertmanager`

## 启动方式

基础启动：

```bash
cd docker
docker compose up -d --build
```

连同 P1 AlertManager 一起启动：

```bash
cd docker
docker compose --profile p1-alerting up -d --build
```

## 关键环境变量

见 `docker/.env`。提交到仓库的 `.env` 应仅保留示例值或占位符，真实敏感信息请在本地或目标环境中覆盖。当前重点变量包括：

- `WEB_PORT`
- `API_PORT`
- `MYSQL_PORT`
- `MYSQL_ROOT_HOST`
- `REDIS_PORT`
- `PROMETHEUS_PORT`
- `PUSHGATEWAY_PORT`
- `ALERTMANAGER_PORT`
- `IMAGE_HOST`
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `MYSQL_ROOT_PASSWORD`
- `REDIS_ADDR`
- `REDIS_PASSWORD`
- `MONITOR_PROMETHEUS_URL`
- `MONITOR_PUSHGATEWAY_URL`
- `MONITOR_AGENT_HEARTBEAT_SERVER_URL`
- `MONITOR_AGENT_HEARTBEAT_TOKEN`
- `MONITOR_WEBHOOK_TOKEN`
- `AI_ENABLED`
- `AI_PROVIDER`
- `AI_BASE_URL`
- `AI_API_KEY`
- `AI_MODEL`
- `AI_REASONING_EFFORT`
- `AI_TIMEOUT_SECONDS`

说明：

- `devops-api` 与 `devops-web` 现在默认从当前仓库源码构建，不再依赖仓库外预构建镜像版本。
- `MYSQL_ROOT_PASSWORD` 同时作为后端默认数据库连接密码使用；如需分离，请额外覆盖 `DB_PASSWORD`。
- `MYSQL_ROOT_HOST` 默认值为 `%`，确保后端容器可通过 `mysql` 服务名连接数据库。

## P1 对应说明

### 告警中心 / AlertManager

- `docker/alertmanager/alertmanager.yml` 提供最小可用的 AlertManager 配置。
- `docker/prometheus/prometheus.yml` 已补充 `alertmanager:9093` 目标。
- `docker/docker-compose.yml` 中 `alertmanager` 为可选 profile，不会强制影响原有部署。

### 配置覆盖

后端现已支持环境变量覆盖 `docker/api/config.yaml` 中的关键配置，因此 Compose 注入的 `DB_*`、`REDIS_*`、`MONITOR_*` 会真实生效。

## 验证建议

启动后至少验证：

1. `http://<host>:<WEB_PORT>` 可访问。
2. `http://<host>:<WEB_PORT>/api/v1/captcha` 返回 `200`。
3. `http://<host>:<WEB_PORT>/monitor/alert-center` 可访问前端路由。
4. `http://<host>:<API_PORT>/api/v1/tool/services` 在带登录态时应返回服务市场数据，不应报 `services.json` 缺失。
5. 如启用 `p1-alerting`，应能在 OpsNexus 中新增 `type=4` 的 AlertManager source 并成功获取状态。

## 注意

- 当前仓库仍未提供完整 Helm chart；Helm 侧最小合同见 [deploy/helm/values-alerting.example.yaml](../deploy/helm/values-alerting.example.yaml)。
- 如果旧文档与当前 Compose 文件冲突，以当前 `docker-compose.yml` 与 `.env` 为准。
