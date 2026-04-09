# OpsNexus 共享环境 UI 巡检报告

## 范围
- 巡检日期: 2026-03-22
- 巡检环境: `http://10.0.0.200:8080`
- 巡检账号: `admin`
- 巡检视口: 1280x900
- 实际巡检路由数: 44
- 证据 JSON: [ui-audit-report.json](C:/zq/平台开发/OpsNexus/OpsNexus/tmp/opsnexus-browser-artifacts/full-ui-audit/crawl/ui-audit-report.json)

## 巡检方法
- 通过共享环境真实登录，读取管理员菜单树。
- 逐个访问所有可见功能模块和子功能模块，共覆盖 44 个唯一路由。
- 对每个页面采集首屏截图、HTML、表格宽度、文本裁切、滚动容器等指标。
- 证据目录: [crawl](C:/zq/平台开发/OpsNexus/OpsNexus/tmp/opsnexus-browser-artifacts/full-ui-audit/crawl)

## 结论摘要
- 严重问题路由: 7 个
- 中等问题路由: 12 个
- 轻微问题路由: 12 个
- 布局稳定路由: 13 个
- 共享环境附带噪音: 44/44 页面都出现 `net::ERR_EMPTY_RESPONSE` 控制台资源错误，属于环境侧共性问题，建议单独排查后端/静态资源链路，不作为页面布局唯一结论依据。

## 关键发现
- 当前最主要的 UI 风险不是整页塌陷，而是“列表页表格宽度策略不统一”，导致不少页面在 1280px 企业常用视口下，进入页面后首屏就需要横向滚动才能看到核心列或操作列。
- 新旧两套页面风格并存。老页面常用自定义 'el-card + search-section + table' 结构，缺少统一的栅格/断点治理；新页面虽然接入了 'TablePage'，但列宽策略仍然偏保守，字段太多时依然会把操作列挤到视口外。
- 仪表盘、AI 工作台、工单申请、知识库、全局搜索等页面首屏整体稳定，说明平台外壳和平台化卡片基线已经可用，问题主要集中在数据密集型列表页面。

## 严重问题
- 资产管理 · 主机管理（`/cmdb/ecs`）: 最大横向超出 718px，截图见 [01-资产管理-主机管理.png](C:/zq/平台开发/OpsNexus/OpsNexus/tmp/opsnexus-browser-artifacts/full-ui-audit/crawl/01-资产管理-主机管理.png)
- 容器管理 · 工作负载（`/k8s/workload`）: 最大横向超出 1108px，截图见 [08-容器管理-工作负载.png](C:/zq/平台开发/OpsNexus/OpsNexus/tmp/opsnexus-browser-artifacts/full-ui-audit/crawl/08-容器管理-工作负载.png)
- 服务管理 · 应用列表（`/app/application`）: 最大横向超出 910px，截图见 [13-服务管理-应用列表.png](C:/zq/平台开发/OpsNexus/OpsNexus/tmp/opsnexus-browser-artifacts/full-ui-audit/crawl/13-服务管理-应用列表.png)
- 任务中心 · Ansible任务（`/task/ansible`）: 最大横向超出 998px，截图见 [17-任务中心-Ansible任务.png](C:/zq/平台开发/OpsNexus/OpsNexus/tmp/opsnexus-browser-artifacts/full-ui-audit/crawl/17-任务中心-Ansible任务.png)
- 工单中心 · 工单列表（`/work/orders`）: 最大横向超出 818px，截图见 [23-工单中心-工单列表.png](C:/zq/平台开发/OpsNexus/OpsNexus/tmp/opsnexus-browser-artifacts/full-ui-audit/crawl/23-工单中心-工单列表.png)
- 监控告警 / 操作审计 · 终端审计（`/monitor/recording`）: 最大横向超出 788px，截图见 [28-监控告警-操作审计-终端审计.png](C:/zq/平台开发/OpsNexus/OpsNexus/tmp/opsnexus-browser-artifacts/full-ui-audit/crawl/28-监控告警-操作审计-终端审计.png)
- 操作审计 · 数据日志（`/monitor/dblog`）: 最大横向超出 706px，截图见 [33-操作审计-数据日志.png](C:/zq/平台开发/OpsNexus/OpsNexus/tmp/opsnexus-browser-artifacts/full-ui-audit/crawl/33-操作审计-数据日志.png)

## 中等问题
- 资产管理 · SQL工单（`/cmdb/sql-work-orders`）: 最大横向超出 548px
- 容器管理 · 节点管理（`/k8s/node`）: 最大横向超出 448px
- 容器管理 · 网络管理（`/k8s/network`）: 最大横向超出 328px
- 任务中心 · 任务作业（`/task/job`）: 最大横向超出 568px
- 任务中心 · 任务模版（`/task/template`）: 最大横向超出 428px
- 任务中心 · 配置管理（`/task/config`）: 最大横向超出 318px
- AI 智能运维助手 · agent列表（`/ops/agent`）: 最大横向超出 453px
- 监控告警 / 操作审计 · 登录日志（`/monitor/loginlog`）: 最大横向超出 370px
- 监控告警 · 告警推送（`/monitor/alert-notify`）: 最大横向超出 637px
- 系统管理 · 用户信息（`/system/admin`）: 最大横向超出 672px
- 系统管理 · 菜单信息（`/system/menu`）: 最大横向超出 420px
- 配置中心 · 通用凭据（`/config/accountauth`）: 最大横向超出 568px

## 轻微问题
- 资产管理 · 数据管理（`/cmdb/db`）: 轻微横向或文本裁切
- 容器管理 · 集群管理（`/k8s/list`）: 轻微横向或文本裁切
- 容器管理 · 命名空间（`/k8s/namespace`）: 轻微横向或文本裁切
- 容器管理 · 存储管理（`/k8s/storage`）: 轻微横向或文本裁切
- 容器管理 · 监控仪表盘（`/k8s/monitoring`）: 轻微横向或文本裁切
- 服务管理 · 快速发布（`/app/quick-release`）: 轻微横向或文本裁切
- 监控告警 / 操作审计 · 操作日志（`/monitor/operator`）: 轻微横向或文本裁切
- 监控告警 · 告警中心（`/monitor/alert-center`）: 轻微横向或文本裁切
- 监控告警 · 告警历史（`/monitor/alert-history`）: 轻微横向或文本裁切
- 系统管理 · 角色信息（`/system/role`）: 轻微横向或文本裁切
- 系统管理 · 岗位信息（`/system/post`）: 轻微横向或文本裁切
- 配置中心 · 密钥管理（`/config/keymanage`）: 轻微横向或文本裁切

## 稳定页面
- 资产管理 · 资产分组（`/cmdb/group`）
- 容器管理 · 配置管理（`/k8s/config`）
- AI 智能运维助手 · 工具列表（`/ops/tools`）
- AI 智能运维助手 · 诊断分析台（`/ai/diagnosis`）
- AI 智能运维助手 · 助手工作台（`/ai/assistant`）
- 工单中心 · 工单申请（`/work/apply`）
- 知识库 · 知识文章（`/knowledge/base`）
- 监控告警 · 监控深化（`/monitor/https`）
- 系统管理 · 部门信息（`/system/dept`）
- 配置中心 · 主机凭据（`/config/ecskey`）
- 配置中心 · LDAP 集成（`/config/ldap`）
- 全局搜索 · 全局搜索（`/search/global`）
- 仪表盘 · 仪表盘（`/dashboard`）

## 模块概览
### 资产管理
- 页面数: 4
- 严重: 1
- 中等: 1
- 轻微: 1
- 稳定: 1

### 容器管理
- 页面数: 8
- 严重: 1
- 中等: 2
- 轻微: 4
- 稳定: 1

### 服务管理
- 页面数: 2
- 严重: 1
- 中等: 0
- 轻微: 1
- 稳定: 0

### 任务中心
- 页面数: 4
- 严重: 1
- 中等: 3
- 轻微: 0
- 稳定: 0

### AI 智能运维助手
- 页面数: 4
- 严重: 0
- 中等: 1
- 轻微: 0
- 稳定: 3

### 工单中心
- 页面数: 2
- 严重: 1
- 中等: 0
- 轻微: 0
- 稳定: 1

### 知识库
- 页面数: 1
- 严重: 0
- 中等: 0
- 轻微: 0
- 稳定: 1

### 监控告警
- 页面数: 7
- 严重: 1
- 中等: 2
- 轻微: 3
- 稳定: 1

### 操作审计
- 页面数: 4
- 严重: 2
- 中等: 1
- 轻微: 1
- 稳定: 0

### 系统管理
- 页面数: 5
- 严重: 0
- 中等: 2
- 轻微: 2
- 稳定: 1

### 配置中心
- 页面数: 4
- 严重: 0
- 中等: 1
- 轻微: 1
- 稳定: 2

### 全局搜索
- 页面数: 1
- 严重: 0
- 中等: 0
- 轻微: 0
- 稳定: 1

### 仪表盘
- 页面数: 1
- 严重: 0
- 中等: 0
- 轻微: 0
- 稳定: 1

## 代码层根因
- 平台壳层在 [Home.vue](C:/zq/平台开发/OpsNexus/OpsNexus/web/src/views/Home.vue) 固定了较宽侧栏，1280px 视口下主内容实际可用宽度被压缩到约 922px 到 1004px；但很多列表页仍按 1400px 以上的列宽思路设计。
- 通用平台基座 [TablePage.vue](C:/zq/平台开发/OpsNexus/OpsNexus/web/src/components/platform/TablePage.vue) 与 [PageToolbar.vue](C:/zq/平台开发/OpsNexus/OpsNexus/web/src/components/platform/PageToolbar.vue) 已经提供统一容器，但缺少统一的“列优先级、操作列固定、次要列折叠、窄屏压缩”规范。
- 全局样式 [global.css](C:/zq/平台开发/OpsNexus/OpsNexus/web/src/assets/css/global.css) 已经统一了平台卡片和表单样式，但没有沉淀统一的数据表格布局策略，导致页面自己定义宽列、图标列、操作列时容易失控。
- 严重页面里仍有较多旧式实现，例如 [cmdbHost.vue](C:/zq/平台开发/OpsNexus/OpsNexus/web/src/views/cmdb/cmdbHost.vue)、[Recording.vue](C:/zq/平台开发/OpsNexus/OpsNexus/web/src/views/monitor/Recording.vue)；这些页面更容易出现左侧树、搜索区和表格互相争抢宽度。
- 中等问题页面虽然接入平台基座，但列仍然偏多，例如 [TaskJob.vue](C:/zq/平台开发/OpsNexus/OpsNexus/web/src/views/task/TaskJob.vue)、[TaskAnsible.vue](C:/zq/平台开发/OpsNexus/OpsNexus/web/src/views/task/TaskAnsible.vue)、[application.vue](C:/zq/平台开发/OpsNexus/OpsNexus/web/src/views/app/application.vue)、[accountauth.vue](C:/zq/平台开发/OpsNexus/OpsNexus/web/src/views/configcenter/accountauth.vue)、[Work-Order-list.vue](C:/zq/平台开发/OpsNexus/OpsNexus/web/src/views/work/Work-Order-list.vue)、[LoginLog.vue](C:/zq/平台开发/OpsNexus/OpsNexus/web/src/views/monitor/LoginLog.vue)。

## 生产可落地 UI 调整方案
### 第一阶段：本周可上线
- 先处理严重页面，目标是“1280px 首屏无需横向滚动即可完成核心查看和主要操作”。
- 对严重页面统一执行三件事:
  1. 固定右侧操作列，但将按钮数量压缩到 2 到 3 个主动作，次要动作移入更多菜单。
  2. 将次要信息列收敛为 tooltip、tag 或详情抽屉，不要求首屏全部展示。
  3. 让工具栏字段自动换行，不再要求单行塞满所有筛选条件。
- 首批优先修复路由:
  - `/cmdb/ecs`
  - `/k8s/workload`
  - `/app/application`
  - `/task/ansible`
  - `/work/orders`
  - `/monitor/recording`
  - `/monitor/dblog`

### 第二阶段：平台化治理
- 基于 [TablePage.vue](C:/zq/平台开发/OpsNexus/OpsNexus/web/src/components/platform/TablePage.vue) 补一层统一的数据表格约束，例如新增 'PlatformDataTable' 包装器或统一 mixin/composable。
- 建立列优先级规范:
  - 一级列: 名称/主键、状态、时间、主操作
  - 二级列: 统计信息、责任人、来源
  - 三级列: 长备注、扩展字段，默认隐藏到详情抽屉
- 在 [global.css](C:/zq/平台开发/OpsNexus/OpsNexus/web/src/assets/css/global.css) 中统一表格 header、cell、toolbar、pagination 的断点行为，避免各页面重复写一套“看起来很像但细节不同”的布局 CSS。

### 第三阶段：回归自动化
- 保留本次巡检脚本 [full-ui-audit.mjs](C:/zq/平台开发/OpsNexus/OpsNexus/tmp/full-ui-audit.mjs)，把它升级成固定的 UI smoke。
- 为关键页面补 Playwright 回归:
  - 首屏截图对比
  - 操作列可见性断言
  - 工具栏折行断言
  - 空态/弹窗稳定性断言
- 把 1280px 作为企业桌面基线，把 1440px 作为宽屏基线，再补一个 960px 的窄屏断点抽样。
