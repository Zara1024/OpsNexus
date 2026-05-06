# OpsNexus

OpsNexus 是一个面向企业运维场景的统一管理平台，覆盖 CMDB、主机纳管、应用发布、任务调度、Kubernetes 管理、监控告警、审计、知识库和 AI 辅助诊断。本仓库已收口为 **仅支持 Docker Compose 部署**，生产环境不再维护 Helm、Kubernetes manifest、systemd 或裸机启动方式。

## 技术栈

| 层级 | 技术 |
| --- | --- |
| 前端 | Vue 3、Vue Router、Vuex、Element Plus、ECharts、Nginx |
| 后端 | Go、Gin、GORM、JWT、Swagger、Kubernetes Client |
| 数据 | MySQL 8、Redis 7 |
| 监控 | Prometheus、Pushgateway、可选 Alertmanager |
| 部署 | Docker、Docker Compose |

## 系统架构

![opsnexus-poster](https://bucketbucket1.oss-cn-beijing.aliyuncs.com/imag/opsnexus-poster.png)

## 目录结构

```text
OpsNexus/
├── api/                  # Go 后端源码、Swagger 产物、服务市场模板
├── web/                  # Vue 3 前端源码和前端测试
├── docker/               # 唯一生产部署入口
│   ├── docker-compose.yml
│   ├── .env.example
│   ├── api/config.yaml
│   ├── web/devops.conf
│   ├── mysql/
│   ├── redis/
│   ├── prometheus/
│   ├── pushgateway/
│   └── alertmanager/
├── docs/assets/          # README 图片资产
├── assets/               # 产品展示和设计资产
└── LICENSE
```

## Docker 部署

生产部署目录建议为 `/opt/opsnexus`。

```bash
cd /opt/opsnexus/docker
cp .env.example .env
vi .env
docker compose up -d --build
```

可选启用 Alertmanager：

```bash
docker compose --profile alerting up -d --build
```

默认访问地址：

- Web：`http://10.0.0.202:8080/`
- API 健康检查：`http://10.0.0.202:8080/api/v1/captcha`
- Prometheus：`http://10.0.0.202:9090/`
- Pushgateway：`http://10.0.0.202:9091/`
- Alertmanager：`http://10.0.0.202:9093/`，仅启用 `alerting` profile 后可用

## 环境变量

复制 `docker/.env.example` 为 `docker/.env` 后再修改。`docker/.env` 不应提交到仓库。

| 变量 | 说明 |
| --- | --- |
| `WEB_PORT` | Web 对外端口，默认 `8080` |
| `API_PORT` | API 对外端口，默认 `8000` |
| `IMAGE_HOST` | 上传文件浏览器访问地址，例如 `http://10.0.0.202:8080` |
| `MYSQL_PORT` | MySQL 映射到宿主机的端口，默认 `3307` |
| `MYSQL_DATABASE` | 初始化数据库名，默认 `devops` |
| `MYSQL_ROOT_PASSWORD` | MySQL root 密码，生产必须替换 |
| `DB_HOST` / `DB_PORT` / `DB_NAME` / `DB_USER` / `DB_PASSWORD` | 后端数据库连接配置 |
| `REDIS_PORT` / `REDIS_ADDR` / `REDIS_PASSWORD` | Redis 端口、容器内地址和密码 |
| `PROMETHEUS_PORT` / `PUSHGATEWAY_PORT` / `ALERTMANAGER_PORT` | 监控组件端口 |
| `MONITOR_AGENT_HEARTBEAT_TOKEN` | Agent 心跳 token |
| `MONITOR_WEBHOOK_TOKEN` | 告警 webhook token |
| `AI_ENABLED` / `AI_API_KEY` / `AI_MODEL` | 可选 AI Runtime 配置 |

## 首次部署

1. 登录服务器，确认 Docker 和 Compose 可用：

```bash
docker version
docker compose version
```

2. 准备部署目录：

```bash
mkdir -p /opt/opsnexus
```

3. 将仓库同步到 `/opt/opsnexus`。

4. 创建生产环境变量：

```bash
cd /opt/opsnexus/docker
cp .env.example .env
vi .env
```

5. 启动：

```bash
docker compose up -d --build
```

6. 验证：

```bash
docker compose ps
curl -I http://127.0.0.1:${WEB_PORT:-8080}/
curl -fsS http://127.0.0.1:${WEB_PORT:-8080}/api/v1/captcha
```

## 更新发布

```bash
cd /opt/opsnexus
git pull
cd docker
docker compose build devops-api devops-web
docker compose up -d
docker compose ps
```

如只改前端或后端，可只构建对应服务：

```bash
docker compose build devops-web
docker compose up -d devops-web
```

## 常用运维命令

```bash
cd /opt/opsnexus/docker
docker compose ps
docker compose logs -f devops-api
docker compose logs -f devops-web
docker compose restart devops-api
docker compose down
docker compose up -d
docker compose pull
docker system df
```

进入容器：

```bash
docker compose exec devops-api sh
docker compose exec mysql mysql -uroot -p
docker compose exec redis redis-cli -a "$REDIS_PASSWORD"
```

## 数据持久化

| 路径 | 内容 |
| --- | --- |
| `docker/mysql/data` | MySQL 数据 |
| `docker/redis/data` | Redis AOF/RDB 数据 |
| `docker/prometheus/data` | Prometheus TSDB |
| `docker/pushgateway/data` | Pushgateway 数据 |
| `docker/alertmanager/data` | Alertmanager 数据 |
| `docker/api/upload` | 用户上传文件 |
| `docker/api/logs` | 后端日志 |
| `docker/api/ssh_keys` | 只读 SSH key 挂载目录 |

建议定期备份 `docker/.env`、`docker/mysql/data`、`docker/api/upload` 和业务所需的 SSH key。备份文件不要提交到 Git。

## 故障排查

| 现象 | 排查命令 |
| --- | --- |
| Web 打不开 | `docker compose ps devops-web`、`docker compose logs devops-web` |
| API 不健康 | `docker compose logs devops-api`、`curl http://127.0.0.1:8080/api/v1/captcha` |
| 数据库连接失败 | 检查 `.env` 中 `DB_PASSWORD` 与 `MYSQL_ROOT_PASSWORD` 是否一致，查看 `docker compose logs mysql` |
| Redis 认证失败 | 检查 `.env` 中 `REDIS_PASSWORD`，查看 `docker compose logs redis` |
| 上传文件无法访问 | 检查 `IMAGE_HOST` 和 `docker/api/upload` 权限 |
| 告警 webhook 401 | 检查 `MONITOR_WEBHOOK_TOKEN` 与 Alertmanager profile 启动日志 |
| 前端静态资源旧 | 重新执行 `docker compose build devops-web && docker compose up -d devops-web` |

## 安全注意事项

- 不要提交 `docker/.env`、真实密码、Token、SSH key、数据库数据、上传文件或日志。
- 生产环境必须替换 `.env.example` 中所有 `replace-with-...` 占位值。
- MySQL、Redis、Prometheus、Pushgateway、Alertmanager 端口如无外部访问需求，建议在防火墙上限制来源。
- `docker/api/ssh_keys` 仅以只读方式挂载到 API 容器，目录权限应限制为 root 可读。
- 对公网开放时建议在服务器前置 HTTPS 网关或云负载均衡，并限制管理端口访问。

## 许可

见 [LICENSE](LICENSE)。
