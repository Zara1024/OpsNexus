# CMDB Assets JumpServer Alignment Design

**Goal:** Rename and reshape the OpsNexus asset menus and pages so `主机`、`网络设备`、`数据库` align with JumpServer 4.10.16 behavior and structure, while keeping the three asset pages independent and preserving existing OpsNexus capabilities such as host terminal, monitoring, database operation, and SQL work orders.

**Reference Baseline**

- Reference source: local code at `C:\zq\平台开发\Github参考项目\jumpserver-4.10.16`
- Runtime verification source: deployed JumpServer on `http://10.0.0.201`
- Verified JumpServer asset-list behavior:
  - top asset type tabs include `全部 / 主机 / 网络设备 / 数据库`
  - common list columns are `名称 / 地址 / 账号 / 平台 / 连接性 / 操作`
  - create forms share a common structure
  - `主机` uses `IP/主机`
  - `网络设备` uses `地址`
  - `数据库` adds `默认数据库`

**Current Problem**

- OpsNexus currently exposes `主机管理` and `数据管理`, but the naming and page structure do not match JumpServer.
- `主机` and `数据库` already exist, but they are shaped as separate legacy management pages rather than JumpServer-like asset pages.
- `网络设备` is missing as a complete asset type.
- The current database model is still closer to a record-management form than an asset-management form.

**Chosen Approach**

- Keep three independent asset pages, one each for `主机`、`网络设备`、`数据库`.
- Align all three pages to a shared asset-page visual and interaction model derived from JumpServer.
- Reuse existing OpsNexus host and database capabilities instead of collapsing everything into one shared asset center page.
- Add a new dedicated network-device backend and frontend flow.
- Extend the existing database asset model and responses so the database page becomes a true asset page.

**Why this approach**

- It satisfies the explicit requirement to keep the three pages independent.
- It still aligns closely with JumpServer at the UX and data-field level.
- It preserves existing host and database integrations with less regression risk than replacing them with a new unified asset center.
- It lets `网络设备` be added cleanly without forcing a larger CMDB rewrite in this round.

**Information Architecture**

- Under `资产管理`:
  - rename `主机管理` to `主机`
  - add `网络设备`
  - rename `数据管理` to `数据库`
- Keep existing `终端登录` available.
- Keep `数据库操作` and `SQL工单`, but treat them as database-asset actions and preserve route compatibility.
- Preserve legacy routes where needed by redirecting to the renamed routes so existing internal links do not break.

**Page Design**

- `主机` page:
  - keep the current left group tree and right content structure
  - reshape the main table toward JumpServer semantics:
    - `名称`
    - `地址`
    - `账号`
    - `平台`
    - `连接性`
    - `操作`
  - keep OpsNexus-specific actions in the detail drawer or action area:
    - terminal
    - monitoring
    - audit
    - sync
- `网络设备` page:
  - use the same overall layout as the host page
  - list columns align with JumpServer:
    - `名称`
    - `地址`
    - `账号`
    - `平台`
    - `连接性`
    - `操作`
  - operations initially include:
    - SSH login
    - Telnet login
    - Web entry
    - connectivity test
- `数据库` page:
  - redesign the current database list into a database asset page
  - list stays close to JumpServer’s common columns
  - database-specific information such as `默认数据库` is shown in form fields and detail/actions
  - preserve:
    - database operation page
    - SQL work order center

**Form Design**

- Shared form sections:
  - `基本设置`
  - `其它设置`
- `主机` required shape:
  - `名称`
  - `IP/主机`
  - `平台`
  - `分组/节点`
  - `协议组` or equivalent protocol settings
  - `账号`
  - `标签`
  - `激活中`
  - `备注`
- `网络设备` required shape:
  - `名称`
  - `地址`
  - `平台`
  - `分组/节点`
  - `协议组`
  - `账号`
  - `标签`
  - `激活中`
  - `备注`
- `数据库` required shape:
  - `名称`
  - `地址`
  - `平台`
  - `分组/节点`
  - `默认数据库`
  - `协议组`
  - `账号`
  - `标签`
  - `激活中`
  - `备注`

**Backend Design**

- `主机`
  - continue using `cmdb_host`
  - continue using existing host SSH/RDP and monitor/audit flows
  - extend list responses if needed so UI can directly render JumpServer-style list columns
- `网络设备`
  - add a dedicated table such as `cmdb_device`
  - core fields:
    - `id`
    - `name`
    - `address`
    - `platform`
    - `group_id`
    - `account_id`
    - `protocol_group`
    - `tags`
    - `is_active`
    - `remark`
  - OpsNexus-specific connection fields:
    - `device_type`
    - `ssh_port`
    - `telnet_port`
    - `web_url`
    - `vendor`
    - optional `snmp_community`
- `数据库`
  - extend `cmdb_sql` beyond the current minimal fields
  - add or map asset-oriented fields:
    - `address`
    - `platform`
    - `default_database`
    - `account_id`
    - `group_id`
    - `protocol_group`
    - `tags`
    - `is_active`
    - `remark`
  - preserve compatibility with SQL execution and SQL work-order services

**API Design**

- Keep separate API families:
  - host: existing `/cmdb/host*`
  - device: new `/cmdb/device*`
  - database: existing `/cmdb/database*` expanded
- Return structures should be normalized enough that all three asset pages can share common display helpers for:
  - address
  - account display
  - platform display
  - connectivity status
  - action enablement

**Implementation Scope**

- Frontend:
  - rename asset menus and route titles
  - reshape host page
  - create network-device page
  - reshape database page
  - preserve host and database action routes
- Backend:
  - add network-device model, DAO, service, controller, router
  - extend database model and CRUD behavior
  - normalize asset list responses
- Deployment:
  - after each verified change batch, sync to `10.0.0.200`

**Non-goals**

- Replacing the three pages with one shared unified asset-center route
- Full network-device configuration management or command orchestration
- Rebuilding host monitoring, terminal, or audit internals
- Replacing the existing SQL work-order flow

**Files**

- Modify: `web/src/router/cmdb.js`
- Modify: `web/src/views/Home.vue`
- Modify: `web/src/utils/platformMenu.js`
- Modify: `web/src/api/cmdb.js`
- Modify: `web/src/views/cmdb/cmdbHost.vue`
- Modify: `web/src/views/cmdb/cmdbDB.vue`
- Create: `web/src/views/cmdb/cmdbDevice.vue`
- Create: supporting `web/src/views/cmdb/Device/*` components if needed
- Modify: `api/router/cmdb/cmdb.go`
- Modify: `api/api/cmdb/model/cmdbSQL.go`
- Modify: `api/api/cmdb/dao/cmdbSQL.go`
- Modify: `api/api/cmdb/service/cmdbSQL.go`
- Modify: `api/api/cmdb/controller/cmdbSQL.go`
- Create: `api/api/cmdb/model/cmdbDevice.go`
- Create: `api/api/cmdb/dao/cmdbDevice.go`
- Create: `api/api/cmdb/service/cmdbDevice.go`
- Create: `api/api/cmdb/controller/cmdbDevice.go`

**Verification**

- Frontend:
  - `cd web && npm run build`
  - targeted UI smoke checks for `主机`、`网络设备`、`数据库`
- Backend:
  - run targeted Go tests for CMDB modules if present
  - add targeted tests for new network-device CRUD and database-field changes where practical
- Functional:
  - menus show `主机 / 网络设备 / 数据库`
  - host page keeps existing core actions working
  - network-device CRUD works end to end
  - database asset CRUD works end to end
  - database operation page and SQL work-order page still work from database assets
- Reference verification:
  - compare final labels and create-form fields against JumpServer on `10.0.0.201`
- Sync verification:
  - after each verified change batch, sync to `10.0.0.200`

**Constraints**

- The current workspace root is not a Git repository, so this spec can be written locally but cannot be committed from the current directory unless you later point me at the actual Git root.
- The superpowers review loop recommends a reviewer subagent, but subagent delegation is not user-authorized in this session, so review remains inline.
