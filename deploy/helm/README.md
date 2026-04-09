# Helm 部署说明

本文档同步日期：`2026-03-20`。

当前 `deploy/helm/opsnexus` 已提供 OpsNexus 的完整 Helm chart，覆盖：

- `api`
- `web`
- `mysql`
- `redis`
- `pushgateway`
- `prometheus`
- `alertmanager`（可选）
- `Ingress`
- `ServiceAccount`
- `PVC`
- `helm test` 烟雾测试

## 目录结构

```text
deploy/helm/
├── README.md
├── values-alerting.example.yaml
├── values-current-source.example.yaml
└── opsnexus/
    ├── Chart.yaml
    ├── values.yaml
    ├── files/
    └── templates/
```

## 安装方式

基础安装：

```bash
helm upgrade --install opsnexus ./deploy/helm/opsnexus \
  -n opsnexus \
  --create-namespace
```

启用 AlertManager / Prometheus 告警链路：

```bash
helm upgrade --install opsnexus ./deploy/helm/opsnexus \
  -n opsnexus \
  --create-namespace \
  -f ./deploy/helm/values-alerting.example.yaml
```

部署“当前仓库源码构建的镜像”：

```bash
cp ./deploy/helm/values-current-source.example.yaml ./values-current-source.yaml
# 编辑 values-current-source.yaml 中的 web.image.* / api.image.*

helm upgrade --install opsnexus ./deploy/helm/opsnexus \
  -n opsnexus \
  --create-namespace \
  -f ./values-current-source.yaml
```

## 常用参数

- `api.config.imageHost`：上传文件对外访问地址
- `ingress.enabled`：是否启用 Ingress
- `global.storageClass`：统一存储类
- `mysql.enabled` / `redis.enabled`：是否使用内置依赖
- `externalServices.mysql.*` / `externalServices.redis.*`：外置依赖接入
- `alertmanager.enabled`：是否启用 AlertManager
- `prometheus.alerting.enabled`：是否让 Prometheus 指向 AlertManager
- `web.image.*` / `api.image.*`：切换到当前源码构建镜像时最关键的覆盖项

当前源码镜像覆盖建议：

- 新机器优先从 `deploy/helm/values-current-source.example.yaml` 复制出本地覆盖文件
- 示例文件默认把 Web 暴露为 `NodePort`，并固定在 `30080`
- 如需额外暴露 API，可把 `api.service.type` 改成 `NodePort`，并显式设置 `api.service.nodePort`
- 如果继续使用默认 `deploy/helm/opsnexus/values.yaml`，仍可能落到历史 `deviops-*` 镜像

安全建议：

- `values.yaml` 应只保留占位符或空值
- 生产环境请通过 `-f values-prod.yaml`、外部 Secret、或 CI/CD 注入真实密码与 token
- 不要把真实 `rootPassword`、`redis password`、`webhook token`、`AI_API_KEY` 提交到仓库

## 验证建议

至少执行：

```bash
helm lint ./deploy/helm/opsnexus
helm template opsnexus ./deploy/helm/opsnexus > /tmp/opsnexus.yaml
helm test opsnexus -n opsnexus
```

## 已验证环境

`2026-03-20` 已在 `10.0.0.200` 的 Kubernetes 环境中完成：

- `helm lint`
- `helm template`
- `helm upgrade --install`
- `helm test`

验证 namespace：`opsnexus-helm-test`

`2026-04-08` 已在 `10.0.0.203` 的 `k3s` 环境中完成当前源码镜像覆盖验证：

- 基于当前仓库源码构建 `web` / `api` 镜像
- 使用 `values-current-source` 风格覆盖文件执行 `helm upgrade --install`
- `helm test opsnexus -n opsnexus` 成功
- `GET /api/v1/captcha` 返回 `200`
- 浏览器登录页与标签页标题均显示为 `OpsNexus`
