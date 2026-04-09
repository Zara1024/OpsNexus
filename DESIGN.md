# OpsNexus DESIGN.md

## Product Context

OpsNexus is a Go + Vue3 enterprise operations platform that combines CMDB, configuration management, task execution, Kubernetes management, monitoring, alerts, audit, knowledge, and AI-assisted operations into one workbench.

The UI should feel like a serious operations product, not a generic admin template and not a flashy demo screen.

## Design Positioning

Use this visual direction:

- Primary direction: mature industrial product surface
- Secondary accent: localized blue tech language
- Product feeling: enterprise workbench + selected command-center moments

In inspiration terms:

- HashiCorp: 70%
- Sentry: 20%
- ClickHouse: 10%

Do not copy any reference literally. Use them only to guide tone, hierarchy, and component discipline.

## Core Rules

1. Preserve all business logic, routes, permissions, APIs, and workflows.
2. Do not redesign for novelty. Redesign for clarity, hierarchy, and product consistency.
3. Use blue tech accents only on flagship operational pages, not everywhere.
4. Keep dense tables, filters, forms, and dialogs highly usable.
5. Avoid generic AI-style dashboards, purple-on-white defaults, and overdecorated sci-fi UI.

## Global Visual System

### Color

- Base pages should use cool-neutral product tones.
- Main brand accent should be a stable blue.
- Brighter cyan-blue should be used sparingly as a tech highlight.
- Status colors should be restrained and operational:
  - success: cool green
  - warning: amber
  - danger: governance red

### Background

- Default lists, forms, and system pages should use light or lightly cool-tinted product surfaces.
- Login, dashboard, monitoring, K8s, logs, and terminal pages may use deeper blue or mixed dark surfaces.
- Do not make the whole product heavy dark mode by default.

### Cards and Panels

Use three panel types:

- `data card`: KPI and metrics
- `module card`: title + action + body
- `container card`: table, filter, form, neutral content

Prefer crisp borders and layered surfaces over large shadows.

### State Expression

State must never rely on color alone. Combine:

- text
- weak background
- border or icon cue
- consistent placement

Use unified categories such as:

- healthy / normal / enabled
- pending / attention / at risk
- alert / abnormal / blocking

## Page Guidance

### Login

- Make it feel credible, premium, and focused.
- Use a stronger brand and atmosphere treatment here than on ordinary pages.
- Layout recommendation:
  - left side: platform value, key capabilities, lightweight visual context
  - right side: clean login card

### Dashboard

- Treat dashboard as an operational overview and quick-action center.
- Strengthen KPI hierarchy, section framing, and chart containers.
- Allow more blue-steel tech language here, but keep it suitable for daily use.

### Monitoring / Alert Center

- Emphasize severity, trend, recent events, and next actions.
- Unify charts, event streams, detail panels, and action surfaces into one visual grammar.

### K8s

- Use a professional control-plane feel.
- Standardize toolbar, filters, status, detail shells, and object hierarchy across clusters, nodes, workloads, network, storage, and monitoring pages.

### Lists

- Keep these pages mature, clean, and predictable.
- Standard structure:
  - page header
  - filter region
  - optional stat strip
  - table container
  - pagination / batch actions

### Forms

- Optimize for long enterprise configuration workflows.
- Use grouping, sectioning, clear help text, and clear primary actions.
- Isolate destructive actions visually.

### Details

- Start with summary header + key status + key metrics + tabs.
- Help the user answer:
  - what is this
  - what state is it in
  - what risk exists
  - what can I do next

### Logs / Terminal / Audit

- These pages may be deeper and more tool-like.
- Prioritize typography clarity, output readability, and strong anomaly emphasis.

## Component Guidance

### Page Headers

Every main page should support:

- title
- short description
- key status or count
- primary actions

### Buttons

Use a clear hierarchy:

- primary
- secondary
- tertiary / ghost
- danger

### Tables

Treat tables as first-class product components. Optimize:

- scanability
- row density
- header clarity
- state rendering
- action rhythm

### Filters

Use one stable pattern:

- common filters first
- advanced filters collapsible
- consistent button positions

### Forms

Keep labels, helper text, validation, and spacing consistent across modules.

### Dialogs and Drawers

- dialogs: short actions, confirmation, light edit
- drawers: side details, longer edit, less disruptive workflows

Avoid modal nesting.

### Motion

Use only restrained motion:

- hover feedback
- selected state
- drawer or tab transitions

Avoid decorative persistent animation.

## Rollout Priority

Start with these page classes:

1. Login
2. Dashboard
3. Monitoring center
4. K8s workbench
5. Shared list / form skeleton

After these sample pages are aligned, extend the same rules to:

- CMDB
- config center
- work orders
- AI pages
- system management

## Prompt Template

Use this as the base prompt for AI-assisted UI optimization:

```text
你正在为 OpsNexus 做 UI 优化。请严格遵守以下规则：

1. 这是一个企业级运维平台，不是官网，也不是大屏展示系统。
2. 保留现有业务逻辑、接口、权限、路由、表单字段和交互流程，不要改动功能语义。
3. 整体视觉方向采用“成熟工业产品面”为主，参考 HashiCorp 70%、Sentry 20%、ClickHouse 10%。
4. 只在关键页面局部吸收蓝色科技语言：登录页、Dashboard、监控中心、K8s 页面、日志/终端/审计页。
5. 列表页、配置页、系统管理页、表单页以克制、专业、耐看的产品面为主，不要做成重度深色炫技页面。
6. 优先优化信息层级、间距节奏、状态表达、标题区、卡片边界、图表容器、表格可读性和表单结构。
7. 避免通用 AI 风格后台界面，避免紫色偏好，避免过度渐变、过度阴影、过度发光。
8. 组件层统一按钮、状态标签、筛选区、表格、抽屉、弹窗和页头骨架。
9. 设计必须适配桌面运营场景，重点保证 1280px 以上工作台体验。
10. 如果在现有设计系统基础上修改，请优先复用已有样式和组件，不要无意义重写。

输出时请先给出：
- 设计思路
- 页面结构调整点
- 视觉语言说明
- 需要修改的 Vue 文件和样式文件
- 最终实现代码
```

## Module Prompts

### Login Prompt

```text
请优化 OpsNexus 登录页 UI。

要求：
- 风格以成熟工业产品面为主，局部吸收蓝色科技语言
- 左侧可展示平台定位、能力摘要、轻量态势信息
- 右侧登录卡片必须干净、可信、聚焦
- 保留现有登录逻辑、验证码逻辑、字段和接口
- 提升标题层级、表单可读性、按钮质感和错误提示
- 不要做成营销官网，也不要做成过度炫酷大屏
```

### Dashboard Prompt

```text
请优化 OpsNexus Dashboard 页面 UI。

要求：
- 保留现有模块结构、数据逻辑和交互行为
- 页面定位是“运维总览 + 风险焦点 + 快速行动中心”
- 用更强的层级处理页头、KPI 卡片、风险区、AI 工作区、图表区、快捷操作区
- 风格以成熟平台产品面为底，局部加入蓝钢科技语言
- 图表容器要像仪表工作区，而不是普通后台卡片
- 不要为了好看牺牲信息密度和可读性
```

### Monitoring Prompt

```text
请优化 OpsNexus 监控中心 / 告警中心页面 UI。

要求：
- 页面重点是态势感知、风险识别和事件处置，不是展示型大屏
- 强化告警等级、趋势、最近事件、处理动作的视觉优先级
- 统一图表、事件流、详情面板、操作区域的视觉语言
- 可使用适度深色和蓝色科技强化，但必须保证长时间使用不疲劳
- 保留现有筛选、列表、图表、抽屉和联动逻辑
```

### List and Form Prompt

```text
请优化 OpsNexus 的列表页和表单页 UI 骨架。

要求：
- 适用于 CMDB、配置中心、系统管理、任务中心等常规业务页
- 默认采用成熟、克制、耐看的产品面
- 统一页头、筛选区、统计条、表格容器、分页、批量操作区
- 统一表单分组、标签、帮助文案、错误提示、危险操作区
- 重点提升密集信息场景下的扫描效率和操作效率
- 不要加入重度科技背景和无意义装饰
```

### K8s and Terminal Prompt

```text
请优化 OpsNexus 的 K8s 页面、日志页、终端页和审计页 UI。

要求：
- 这类页面允许比普通页面更强的蓝色科技语言，但仍需保持专业控制台气质
- 统一工具条、状态、资源卡、详情区、输出区、侧边面板的层级
- 日志和终端强调可读性、边界感、风险提示和长时间操作舒适度
- 保留原有业务逻辑、对象结构、操作流程和权限控制
- 不要做成纯展示型炫技战情室
```
