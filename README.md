# OpsNexus（智维云枢）

OpsNexus 是面向中大型企业的一站式智能运维管理平台。本仓库为第一阶段骨架工程，包含后端基础框架、前端基础框架与数据库迁移脚本。

## 目录结构

- `cloudops-server/`：Go 后端服务
- `cloudops-web/`：Vue3 前端应用
- `migrations/`：数据库迁移脚本
- `docker-compose.dev.yml`：本地 PostgreSQL/Redis 开发依赖

## 快速启动

### 1）启动依赖

```bash
docker compose -f docker-compose.dev.yml up -d
```

### 2）启动后端

```bash
cd cloudops-server
go mod tidy
go run ./cmd/server/main.go
```

访问健康检查：`GET http://localhost:8080/health`

### 3）启动前端

```bash
cd cloudops-web
npm install
npm run dev
```

浏览器访问：`http://localhost:5173`
