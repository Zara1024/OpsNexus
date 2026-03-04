from pathlib import Path
import textwrap

files = {
    # backend core
    "cloudops-server/pkg/response/response.go": """
package response

import "github.com/gin-gonic/gin"

// Body 定义统一响应体。
type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 返回成功响应。
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Body{Code: 200, Message: "success", Data: data})
}

// Fail 返回业务失败响应。
func Fail(c *gin.Context, code int, message string) {
	c.JSON(200, Body{Code: code, Message: message, Data: gin.H{}})
}

// Error 返回HTTP错误响应。
func Error(c *gin.Context, httpCode int, code int, message string) {
	c.JSON(httpCode, Body{Code: code, Message: message, Data: gin.H{}})
}
""",
    "cloudops-server/pkg/errors/errors.go": """
package errors

// AppError 定义统一应用错误结构。
type AppError struct {
	Code     int    `json:"code"`
	HTTPCode int    `json:"-"`
	Message  string `json:"message"`
	Detail   string `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrUnauthorized = &AppError{Code: 10001, HTTPCode: 401, Message: "未登录或登录已过期"}
	ErrForbidden    = &AppError{Code: 10003, HTTPCode: 403, Message: "无权限访问"}
	ErrNotFound     = &AppError{Code: 10004, HTTPCode: 404, Message: "资源不存在"}
	ErrValidation   = &AppError{Code: 10005, HTTPCode: 400, Message: "参数校验失败"}
	ErrInternal     = &AppError{Code: 10006, HTTPCode: 500, Message: "服务器内部错误"}
)
""",
    "cloudops-server/pkg/middleware/cors.go": """
package middleware

import "github.com/gin-gonic/gin"

// CORS 提供跨域处理中间件。
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
""",
    "cloudops-server/pkg/middleware/recovery.go": """
package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

// Recovery 捕获panic并返回统一错误。
func Recovery(log *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Error("请求发生panic", "panic", recovered)
		response.Error(c, http.StatusInternalServerError, apperrors.ErrInternal.Code, apperrors.ErrInternal.Message)
	})
}
""",
    "cloudops-server/pkg/middleware/request_id.go": """
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDKey = "request_id"

// RequestID 为每个请求注入request_id。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set(requestIDKey, requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

// GetRequestID 从上下文获取request_id。
func GetRequestID(c *gin.Context) string {
	if v, ok := c.Get(requestIDKey); ok {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}
""",
    "cloudops-server/pkg/middleware/logger.go": """
package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger 记录HTTP请求访问日志。
func RequestLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		log.Info("HTTP请求",
			"request_id", GetRequestID(c),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start).String(),
			"client_ip", c.ClientIP(),
		)
	}
}
""",
    "cloudops-server/Makefile": """
APP_NAME=opsnexus-server

.PHONY: build run test lint tidy migrate-up migrate-down

build:
	go build -o bin/$(APP_NAME) ./cmd/server

run:
	go run ./cmd/server/main.go

test:
	go test ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

migrate-up:
	migrate -path ../migrations -database \"postgres://postgres:postgres@127.0.0.1:5432/opsnexus?sslmode=disable\" up

migrate-down:
	migrate -path ../migrations -database \"postgres://postgres:postgres@127.0.0.1:5432/opsnexus?sslmode=disable\" down
""",
    "cloudops-server/Dockerfile": """
# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /opsnexus-server ./cmd/server

FROM alpine:3.20
WORKDIR /app
RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /opsnexus-server /usr/local/bin/opsnexus-server
COPY config.yaml /app/config.yaml

USER app
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/opsnexus-server"]
""",
    "cloudops-server/.golangci.yml": """
run:
  timeout: 5m
  tests: true

linters:
  enable:
    - govet
    - revive
    - staticcheck
    - gofmt
    - goimports

issues:
  max-issues-per-linter: 0
  max-same-issues: 0
""",

    # frontend
    "cloudops-web/package.json": """
{
  "name": "cloudops-web",
  "private": true,
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc -b && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "@iconify/vue": "^4.1.2",
    "axios": "^1.8.4",
    "element-plus": "^2.10.0",
    "pinia": "^3.0.1",
    "vue": "^3.5.13",
    "vue-router": "^4.5.1"
  },
  "devDependencies": {
    "@types/node": "^22.13.11",
    "@vitejs/plugin-vue": "^5.2.1",
    "sass": "^1.85.1",
    "typescript": "~5.7.2",
    "unocss": "^66.0.0",
    "unplugin-auto-import": "^19.1.1",
    "unplugin-vue-components": "^28.4.1",
    "vite": "^6.2.0",
    "vue-tsc": "^2.2.8"
  }
}
""",
    "cloudops-web/index.html": """
<!doctype html>
<html lang="zh-CN">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>OpsNexus</title>
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.ts"></script>
  </body>
</html>
""",
    "cloudops-web/tsconfig.json": """
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "strict": true,
    "jsx": "preserve",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "esModuleInterop": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "types": ["node"],
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"]
    }
  },
  "include": ["src/**/*.ts", "src/**/*.d.ts", "src/**/*.tsx", "src/**/*.vue", "vite.config.ts", "uno.config.ts"]
}
""",
    "cloudops-web/vite.config.ts": """
import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import UnoCSS from 'unocss/vite'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

export default defineConfig({
  plugins: [
    vue(),
    UnoCSS(),
    AutoImport({
      resolvers: [ElementPlusResolver()],
    }),
    Components({
      resolvers: [ElementPlusResolver()],
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
""",
    "cloudops-web/uno.config.ts": """
import { defineConfig, presetAttributify, presetUno, presetIcons } from 'unocss'

export default defineConfig({
  presets: [presetUno(), presetAttributify(), presetIcons()],
})
""",
    "cloudops-web/src/vite-env.d.ts": """
/// <reference types=\"vite/client\" />
""",
    "cloudops-web/src/main.ts": """
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'uno.css'
import './styles/reset.scss'
import './styles/variables.scss'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(ElementPlus)

app.mount('#app')
""",
    "cloudops-web/src/App.vue": """
<template>
  <RouterView />
</template>
""",
    "cloudops-web/src/styles/variables.scss": """
:root {
  --color-bg-primary: #0b1020;
  --color-bg-secondary: #151b35;
  --color-accent-start: #1a1d4d;
  --color-accent-mid: #667eea;
  --color-accent-end: #764ba2;
  --color-text-primary: #f5f7ff;
  --color-text-secondary: #b9c1e8;
  --radius-md: 12px;
}

body {
  background: radial-gradient(circle at top right, var(--color-accent-end), var(--color-bg-primary) 45%);
  color: var(--color-text-primary);
  font-family: 'Source Han Sans SC', 'Inter', sans-serif;
}
""",
    "cloudops-web/src/styles/reset.scss": """
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body,
#app {
  width: 100%;
  height: 100%;
}
""",
    "cloudops-web/src/utils/request.ts": """
import axios from 'axios'

const request = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  (response) => response.data,
  (error) => {
    return Promise.reject(error)
  },
)

export default request
""",
    "cloudops-web/src/stores/app.ts": """
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    sidebarCollapsed: false,
    theme: 'dark',
  }),
  actions: {
    toggleSidebar() {
      this.sidebarCollapsed = !this.sidebarCollapsed
    },
  },
})
""",
    "cloudops-web/src/stores/user.ts": """
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: '',
    userInfo: null as Record<string, unknown> | null,
  }),
})
""",
    "cloudops-web/src/router/index.ts": """
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
  },
  {
    path: '/',
    component: () => import('@/components/Layout/AppLayout.vue'),
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
""",
    "cloudops-web/src/components/Layout/AppLayout.vue": """
<template>
  <div class="layout">
    <aside class="sidebar">OpsNexus</aside>
    <div class="main">
      <header class="header">顶栏区域</header>
      <section class="content">
        <RouterView />
      </section>
    </div>
  </div>
</template>

<style scoped lang="scss">
.layout {
  display: flex;
  width: 100%;
  height: 100%;
}

.sidebar {
  width: 220px;
  background: rgba(15, 20, 40, 0.8);
  backdrop-filter: blur(8px);
  padding: 20px;
}

.main {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.header {
  height: 64px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  display: flex;
  align-items: center;
  padding: 0 20px;
}

.content {
  flex: 1;
  padding: 20px;
}
</style>
""",
    "cloudops-web/src/views/login/index.vue": """
<template>
  <div class="login-page">
    <div class="login-card">
      <h1>OpsNexus</h1>
      <p>智维云枢 · 智能运维管理平台</p>
      <el-form>
        <el-form-item>
          <el-input placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item>
          <el-input type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-button type="primary" class="w-full">登录</el-button>
      </el-form>
    </div>
  </div>
</template>

<style scoped lang="scss">
.login-page {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-card {
  width: 380px;
  padding: 28px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(12px);
}

h1 {
  margin-bottom: 8px;
}

p {
  margin-bottom: 20px;
  color: #b9c1e8;
}

.w-full {
  width: 100%;
}
</style>
""",
    "cloudops-web/src/views/dashboard/index.vue": """
<template>
  <el-card>
    <h2>仪表盘</h2>
    <p>这里是仪表盘占位页面。</p>
  </el-card>
</template>
""",

    # migrations
    "migrations/000001_create_sys_users.up.sql": """
CREATE TABLE IF NOT EXISTS sys_users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(64) NOT NULL DEFAULT '',
    email VARCHAR(128) NOT NULL DEFAULT '',
    phone VARCHAR(32) NOT NULL DEFAULT '',
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
""",
    "migrations/000001_create_sys_users.down.sql": """
DROP TABLE IF EXISTS sys_users;
""",
    "migrations/000002_create_sys_roles.up.sql": """
CREATE TABLE IF NOT EXISTS sys_roles (
    id BIGSERIAL PRIMARY KEY,
    role_code VARCHAR(64) NOT NULL UNIQUE,
    role_name VARCHAR(64) NOT NULL,
    description VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
""",
    "migrations/000002_create_sys_roles.down.sql": """
DROP TABLE IF EXISTS sys_roles;
""",
    "migrations/000003_create_sys_menus.up.sql": """
CREATE TABLE IF NOT EXISTS sys_menus (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT NOT NULL DEFAULT 0,
    menu_name VARCHAR(64) NOT NULL,
    route_path VARCHAR(128) NOT NULL DEFAULT '',
    component VARCHAR(128) NOT NULL DEFAULT '',
    icon VARCHAR(64) NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
""",
    "migrations/000003_create_sys_menus.down.sql": """
DROP TABLE IF EXISTS sys_menus;
""",
    "migrations/000004_create_user_roles.up.sql": """
CREATE TABLE IF NOT EXISTS sys_user_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_sys_user_roles UNIQUE (user_id, role_id)
);
""",
    "migrations/000004_create_user_roles.down.sql": """
DROP TABLE IF EXISTS sys_user_roles;
""",
    "migrations/000005_create_role_menus.up.sql": """
CREATE TABLE IF NOT EXISTS sys_role_menus (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,
    menu_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_sys_role_menus UNIQUE (role_id, menu_id)
);
""",
    "migrations/000005_create_role_menus.down.sql": """
DROP TABLE IF EXISTS sys_role_menus;
""",

    # root
    "docker-compose.dev.yml": """
services:
  postgres:
    image: postgres:16
    container_name: opsnexus-postgres
    environment:
      POSTGRES_DB: opsnexus
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - pg_data:/var/lib/postgresql/data

  redis:
    image: redis:7
    container_name: opsnexus-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  pg_data:
  redis_data:
""",
    ".gitignore": """
# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
vendor/

# Node
node_modules/
dist/

# IDE
.vscode/
.idea/

# OS
.DS_Store
Thumbs.db

# Runtime
*.log
.env
""",
    "README.md": """
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
""",
}

# 生成后端模块空文件
modules = ["auth", "cmdb", "k8s", "app", "monitor", "task", "audit", "aiops"]
layers = ["handler", "service", "repository", "model"]
for module in modules:
    for layer in layers:
        pkg_name = layer
        title = layer.capitalize()
        path = f"cloudops-server/internal/{module}/{layer}/{layer}.go"
        files[path] = f"""
package {pkg_name}

// {title} 定义{module}模块{layer}层空结构。
type {title} struct{{}}
"""

for rel, content in files.items():
    p = Path(rel)
    p.parent.mkdir(parents=True, exist_ok=True)
    p.write_text(textwrap.dedent(content).lstrip("\n"), encoding="utf-8")

print("scaffold files written")
