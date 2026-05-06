# OpsNexus 项目架构图提示词

请基于当前仓库源码，生成一张 `OpsNexus` 的项目架构图，输出格式必须是可直接打开的纯 `SVG` 文件。

## 目标

- 图的重点不是单纯展示运行时组件，而是同时表达“代码结构 + 运行依赖 + 外部集成 + 部署路径”。
- 图要让第一次接手仓库的开发者能在 10 秒内看懂：入口在哪里、前后端如何分层、核心业务域有哪些、依赖哪些基础设施、如何部署。
- 图风格应偏企业级运维平台，克制、清晰、可嵌入 README，不要做成炫技大屏。

## 项目事实

- 仓库是前后端分离架构。
- 前端目录是 `web/`，技术栈是 `Vue 3 + Element Plus + Vue Router + Vuex`。
- 前端核心目录包括 `web/src/router`、`web/src/api`、`web/src/views`、`web/src/components`、`web/src/utils`、`web/src/permission`。
- 后端目录是 `api/`，技术栈是 `Go + Gin + GORM`。
- 后端入口是 `api/main.go`，顶层路由聚合在 `api/router/router.go`。
- 后端存在清晰的领域分层：`api/api/<domain>/{controller,service,dao,model}`。
- 后端公共能力位于 `api/common`、`api/middleware`、`api/pkg`、`api/scheduler`。
- 主要业务域包括 `system`、`cmdb`、`configCenter`、`app`、`work`、`task`、`k8s`、`monitor`、`knowledge`、`search`、`ai`、`tool`。
- 主要内部依赖包括 `MySQL`、`Redis`、`Prometheus`、`Pushgateway`、`Alertmanager`、上传静态文件目录与 SSH key 挂载。
- 主要外部集成包括 `Hosts / SSH / Ansible`、`Kubernetes Clusters`、`OpenAI-compatible AI Runtime`。
- 项目支持多种部署方式：`Docker Compose`、`Helm Chart`、原生 `Kubernetes manifests`、`systemd + Nginx`。

## 推荐布局

- 使用横向主流程，从左到右依次表现 `User / Browser -> Web UI -> API Platform -> Dependencies / External Systems`。
- 在主流程下方增加一块 `Business Domains` 区域，把业务模块分成 3 到 4 组，而不是把十几个模块平铺成碎片。
- 在左下角单独放 `Delivery Paths`，说明仓库如何被部署。
- 在图中直接写出关键目录路径，强调这是“项目架构图”而不是抽象产品图。

## 必须包含的结构块

- `Access Layer`
- `Web UI / web/`
- `API Platform / api/`
- `Business Domains`
- `Internal Dependencies`
- `Managed Targets & External Integrations`
- `Delivery Paths`

## Business Domains 分组建议

- `Governance & Assets`: `system / cmdb / configCenter`
- `Delivery & Jobs`: `app / work / task / tool`
- `Platform & Observability`: `k8s / monitor`
- `Knowledge & AI`: `knowledge / search / ai`

## 视觉要求

- 使用蓝灰主色，局部用青蓝、靛蓝、绿色、橙色区分层次。
- 整体背景轻，卡片白色或浅色，边框清晰，阴影很轻。
- 每个大卡片顶部有一条细色带，帮助快速区分层级。
- 连线尽量少但语义明确，可以给关键连线加简短标签，例如 `HTTPS`、`REST / SSE / WS`、`DB / cache / metrics`、`SSH / K8s API / LLM API`。

## SVG 约束

- 必须输出纯 SVG。
- 只能使用标准 SVG 元素，例如 `rect`、`path`、`line`、`text`、`tspan`、`defs`、`linearGradient`、`filter`、`marker`。
- 不要使用 `foreignObject`。
- 不要依赖外部 CSS、外链字体、脚本或位图资源。
- 文字必须使用真实的 SVG `text` 节点，保证在 README、浏览器和文档平台中稳定显示。
- 建议画布尺寸在 `1600 x 1000` 左右，适合桌面阅读。

## 输出要求

- 文件名使用 `assets/opsnexus-architecture.svg`。
- 标题使用 `OpsNexus 项目架构图`。
- 副标题体现 `Code structure + runtime dependencies`。
- 在页脚补一行简短说明，指出图基于 `README.md`、`api/router/router.go`、`web/src/router/router.js`、`docker/docker-compose.yml`、`deploy/helm/opsnexus/values.yaml` 提炼。
