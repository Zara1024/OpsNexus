# CMDB Assets JumpServer Alignment Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rename and rebuild the `主机`、`网络设备`、`数据库` asset pages so they independently align with JumpServer’s asset UX, preserve existing OpsNexus host/database capabilities, and add a complete network-device CRUD and connection flow.

**Architecture:** Keep three independent asset routes and pages instead of a unified asset center, but standardize their list columns, create/edit forms, and action semantics around the JumpServer 4.10.16 asset model. Reuse existing host and database capabilities where possible, introduce a dedicated network-device backend/frontend slice, and sync verified frontend/backend artifacts to `10.0.0.200` after each finished change batch.

**Tech Stack:** Vue 3, Element Plus, Vue Router, Node `assert`-based frontend tests under `web/tests`, Go 1.24, Gin, GORM, PowerShell build/deploy commands, SCP/SSH deployment to `10.0.0.200`

---

## File Structure Map

- `web/src/router/cmdb.js`
  Route titles, route additions, and backward-compatible redirects for renamed asset pages.
- `web/src/views/Home.vue`
  Home-page copy updates so asset descriptions match the renamed menus.
- `web/src/utils/platformMenu.js`
  Menu normalization and injected child-menu naming for `主机`、`网络设备`、`数据库`.
- `web/src/api/cmdb.js`
  Frontend request layer for host, device, and database asset CRUD plus connectivity actions.
- `web/src/utils/cmdbPresentation.mjs`
  Existing CMDB presentation helpers. Extend only if database-specific helpers still belong here.
- `web/src/utils/cmdbAssetPresentation.mjs`
  New shared asset-page presentation helpers for list rows, addresses, account display, connectivity badges, and type-specific create/edit defaults.
- `web/src/views/cmdb/cmdbHost.vue`
  Host asset page shell, filters, batch actions, and detail entrypoints.
- `web/src/views/cmdb/Host/CmdbHostTable.vue`
  Host asset table layout and action-column presentation.
- `web/src/views/cmdb/Host/CreateHost.vue`
  Host create form aligned to JumpServer-style asset form sections.
- `web/src/views/cmdb/Host/EditHost.vue`
  Host edit form aligned to the same field structure.
- `web/src/views/cmdb/cmdbDevice.vue`
  New network-device asset page.
- `web/src/views/cmdb/Device/CmdbDeviceTable.vue`
  New network-device table component.
- `web/src/views/cmdb/Device/CreateDevice.vue`
  New network-device create form.
- `web/src/views/cmdb/Device/EditDevice.vue`
  New network-device edit form.
- `web/src/views/cmdb/cmdbDB.vue`
  Database asset page rebuilt from the current record-style page into an asset-style page.
- `web/src/views/cmdb/DBdetails.vue`
  Database detail page title/breadcrumb updates and asset-context entry handling.
- `web/src/views/cmdb/SQLWorkOrderCenter.vue`
  Database asset-context wording alignment where needed.
- `web/tests/cmdb-assets-route-meta.test.mjs`
  New frontend test for route/menu naming and asset-page exposure.
- `web/tests/cmdb-asset-presentation.test.mjs`
  New frontend test for shared asset presentation helpers.
- `web/tests/cmdb-host-asset-presentation.test.mjs`
  New frontend test for host asset row and form mapping.
- `web/tests/cmdb-device-presentation.test.mjs`
  New frontend test for network-device-specific list/form mapping.
- `web/tests/cmdb-database-asset-presentation.test.mjs`
  New frontend test for database asset row and form mapping.
- `api/router/cmdb/cmdb.go`
  CMDB route registration for device CRUD and expanded database asset endpoints.
- `api/api/cmdb/model/cmdbHost.go`
  Host model reference for list/display compatibility if response DTOs need extension.
- `api/api/cmdb/model/cmdbDevice.go`
  New network-device model and DTOs.
- `api/api/cmdb/dao/cmdbDevice.go`
  New network-device DAO.
- `api/api/cmdb/service/cmdbDevice.go`
  New network-device service.
- `api/api/cmdb/controller/cmdbDevice.go`
  New network-device controller.
- `api/api/cmdb/service/cmdbDevice_test.go`
  New backend test for network-device request normalization and connection selection logic.
- `api/api/cmdb/model/cmdbSQL.go`
  Database asset model expansion.
- `api/api/cmdb/dao/cmdbSQL.go`
  Database DAO query/field updates.
- `api/api/cmdb/service/cmdbSQL.go`
  Database service updates for asset-oriented CRUD.
- `api/api/cmdb/controller/cmdbSQL.go`
  Database controller updates for expanded payloads and list filtering.
- `api/api/cmdb/service/cmdbSQL_asset_test.go`
  New backend test for database asset normalization and compatibility behavior.

**Execution Notes**

- The current working directory is not the Git root, so commit steps below must be run from the actual repository root if/when it is available.
- The remote deployment shape is already documented in `README.md`:
  - backend binary target: `/opt/opsnexus-remote-test/opsnexus-api-linux-amd64`
  - frontend dist target: `/opt/opsnexus-remote-test/web-dist`
  - restart commands: `systemctl restart opsnexus-api.service` and `docker restart opsnexus-web`

### Task 1: Rename Menus And Align CMDB Route Metadata

**Files:**
- Modify: `web/src/router/cmdb.js`
- Modify: `web/src/views/Home.vue`
- Modify: `web/src/utils/platformMenu.js`
- Create: `web/src/views/cmdb/cmdbDevice.vue`
- Test: `web/tests/cmdb-assets-route-meta.test.mjs`

- [ ] **Step 1: Write the failing route/menu metadata test**

```js
import assert from 'node:assert/strict'
import cmdbRoutes from '../src/router/cmdb.js'
import { normalizePlatformMenu } from '../src/utils/platformMenu.js'

function getRoute(path) {
  return cmdbRoutes.find((item) => item.path === path)
}

const hostRoute = getRoute('/cmdb/ecs')
const deviceRoute = getRoute('/cmdb/device')
const databaseRoute = getRoute('/cmdb/db')

assert.equal(hostRoute.meta.tTitle, '主机')
assert.equal(databaseRoute.meta.tTitle, '数据库')
assert.equal(deviceRoute.meta.tTitle, '网络设备')

const normalized = normalizePlatformMenu([
  { id: 1, menuName: '资产管理', menuSvoList: [] }
])

const assetChildren = normalized.find((item) => item.menuName === '资产管理').menuSvoList
assert.ok(assetChildren.some((item) => item.menuName === '主机'))
assert.ok(assetChildren.some((item) => item.menuName === '网络设备'))
assert.ok(assetChildren.some((item) => item.menuName === '数据库'))

console.log('cmdb asset route/meta tests passed')
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `node web/tests/cmdb-assets-route-meta.test.mjs`

Expected: FAIL because `/cmdb/device` does not exist yet and current titles still resolve to `主机管理` / `数据管理`.

- [ ] **Step 3: Implement the route and menu metadata changes**

```js
// web/src/router/cmdb.js
import Device from '@/views/cmdb/cmdbDevice.vue'

const routes = [
  {
    path: '/cmdb/ecs',
    component: Host,
    meta: { sTitle: '资产管理', tTitle: '主机' }
  },
  {
    path: '/cmdb/device',
    component: Device,
    meta: { sTitle: '资产管理', tTitle: '网络设备' }
  },
  {
    path: '/cmdb/db',
    component: Db,
    meta: { sTitle: '资产管理', tTitle: '数据库' }
  }
]
```

```js
// web/src/utils/platformMenu.js
{
  match: menu => menu.menuName === '资产管理',
  items: [
    { id: 99981, menuName: '主机', url: 'cmdb/ecs', icon: 'Monitor' },
    { id: 99982, menuName: '网络设备', url: 'cmdb/device', icon: 'Connection' },
    { id: 99983, menuName: '数据库', url: 'cmdb/db', icon: 'Coin' },
    { id: 99993, menuName: 'SQL工单', url: 'cmdb/sql-work-orders', icon: 'Tickets' }
  ]
}
```

```vue
<!-- web/src/views/cmdb/cmdbDevice.vue -->
<template>
  <div class="cmdb-device-management">
    <el-empty description="网络设备页面建设中" />
  </div>
</template>
```

- [ ] **Step 4: Run the test to verify it passes**

Run: `node web/tests/cmdb-assets-route-meta.test.mjs`

Expected: PASS with `cmdb asset route/meta tests passed`

- [ ] **Step 5: Commit the metadata-only batch**

```bash
git add web/src/router/cmdb.js web/src/views/Home.vue web/src/utils/platformMenu.js web/src/views/cmdb/cmdbDevice.vue web/tests/cmdb-assets-route-meta.test.mjs
git commit -m "feat: align cmdb asset menu titles with jumpserver"
```

### Task 2: Add Shared CMDB Asset Presentation Helpers

**Files:**
- Create: `web/src/utils/cmdbAssetPresentation.mjs`
- Test: `web/tests/cmdb-asset-presentation.test.mjs`

- [ ] **Step 1: Write the failing helper test**

```js
import assert from 'node:assert/strict'
import {
  buildAssetAddress,
  buildAssetAccountLabel,
  buildAssetConnectivityTag,
  createAssetFormSections
} from '../src/utils/cmdbAssetPresentation.mjs'

assert.equal(buildAssetAddress({ type: 'host', sshIp: '10.0.0.8' }), '10.0.0.8')
assert.equal(buildAssetAddress({ type: 'device', address: '10.0.0.9' }), '10.0.0.9')
assert.equal(buildAssetAddress({ type: 'database', address: 'db.internal:3306' }), 'db.internal:3306')

assert.equal(buildAssetAccountLabel({ accountName: 'root' }), 'root')
assert.deepEqual(buildAssetConnectivityTag({ reachable: true }), { text: '成功', type: 'success' })
assert.equal(createAssetFormSections('database')[0].title, '基本设置')

console.log('cmdb asset presentation tests passed')
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `node web/tests/cmdb-asset-presentation.test.mjs`

Expected: FAIL with module-not-found because `web/src/utils/cmdbAssetPresentation.mjs` does not exist.

- [ ] **Step 3: Implement the shared helper module**

```js
export function buildAssetAddress(asset) {
  if (asset.type === 'host') return asset.sshIp || asset.publicIp || asset.privateIp || ''
  if (asset.type === 'device') return asset.address || ''
  return asset.address || asset.host || ''
}

export function buildAssetConnectivityTag(asset) {
  return asset.reachable
    ? { text: '成功', type: 'success' }
    : { text: '失败', type: 'danger' }
}

export function createAssetFormSections(assetType) {
  const base = [{ title: '基本设置' }, { title: '其它设置' }]
  return assetType === 'database' ? base : base
}
```

- [ ] **Step 4: Run the test to verify it passes**

Run: `node web/tests/cmdb-asset-presentation.test.mjs`

Expected: PASS with `cmdb asset presentation tests passed`

- [ ] **Step 5: Commit the helper foundation**

```bash
git add web/src/utils/cmdbAssetPresentation.mjs web/tests/cmdb-asset-presentation.test.mjs
git commit -m "feat: add shared cmdb asset presentation helpers"
```

### Task 3: Realign The Host Asset Page To JumpServer Semantics

**Files:**
- Modify: `web/src/api/cmdb.js`
- Modify: `web/src/views/cmdb/cmdbHost.vue`
- Modify: `web/src/views/cmdb/Host/CmdbHostTable.vue`
- Modify: `web/src/views/cmdb/Host/CreateHost.vue`
- Modify: `web/src/views/cmdb/Host/EditHost.vue`
- Test: `web/tests/cmdb-host-asset-presentation.test.mjs`

- [ ] **Step 1: Write the failing host asset test**

```js
import assert from 'node:assert/strict'
import { mapHostRowToAssetRow, createHostAssetFormModel } from '../src/utils/cmdbAssetPresentation.mjs'

const row = mapHostRowToAssetRow({
  id: 7,
  hostName: 'prod-host-01',
  sshIp: '10.0.0.7',
  sshName: 'root',
  os: 'Linux',
  status: 1
})

assert.equal(row.name, 'prod-host-01')
assert.equal(row.address, '10.0.0.7')
assert.equal(row.account, 'root')
assert.equal(row.platform, 'Linux')
assert.equal(row.connectivity.text, '成功')

const form = createHostAssetFormModel()
assert.equal(form.hostName, '')
assert.equal(form.sshIp, '')

console.log('cmdb host asset presentation tests passed')
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `node web/tests/cmdb-host-asset-presentation.test.mjs`

Expected: FAIL because host-specific asset helpers are not implemented yet.

- [ ] **Step 3: Implement the host mapping helpers and update the host page**

```js
export function mapHostRowToAssetRow(host) {
  return {
    id: host.id,
    name: host.hostName || host.name || '',
    address: buildAssetAddress({ type: 'host', ...host }),
    account: host.sshName || '',
    platform: host.os || host.deviceType || 'Linux',
    connectivity: buildAssetConnectivityTag({ reachable: Number(host.status) === 1 }),
    raw: host
  }
}
```

```vue
<!-- CmdbHostTable.vue -->
<el-table-column prop="name" label="名称" min-width="180" />
<el-table-column prop="address" label="地址" min-width="160" />
<el-table-column prop="account" label="账号" min-width="140" />
<el-table-column prop="platform" label="平台" min-width="120" />
<el-table-column label="连接性" min-width="120">
  <template #default="{ row }">
    <el-tag :type="row.connectivity.type">{{ row.connectivity.text }}</el-tag>
  </template>
</el-table-column>
```

- [ ] **Step 4: Run the test to verify it passes**

Run: `node web/tests/cmdb-host-asset-presentation.test.mjs`

Expected: PASS with `cmdb host asset presentation tests passed`

- [ ] **Step 5: Run a host-page build verification**

Run: `cd web && npm run build`

Expected: PASS and the build completes without new CMDB host errors.

- [ ] **Step 6: Commit and sync the host/menu frontend batch**

```bash
git add web/src/api/cmdb.js web/src/views/cmdb/cmdbHost.vue web/src/views/cmdb/Host/CmdbHostTable.vue web/src/views/cmdb/Host/CreateHost.vue web/src/views/cmdb/Host/EditHost.vue web/src/utils/cmdbAssetPresentation.mjs web/tests/cmdb-host-asset-presentation.test.mjs
git commit -m "feat: align host asset page with jumpserver layout"
```

```powershell
tar -czf tmp\web-dist-cmdb-assets-host.tar.gz -C web dist
scp tmp\web-dist-cmdb-assets-host.tar.gz root@10.0.0.200:/tmp/web-dist-cmdb-assets-host.tar.gz
ssh root@10.0.0.200 "rm -rf /opt/opsnexus-remote-test/web-dist/* && tar -xzf /tmp/web-dist-cmdb-assets-host.tar.gz -C /opt/opsnexus-remote-test/web-dist --strip-components=1 && docker restart opsnexus-web"
```

Expected: remote `http://10.0.0.200:8080/` shows renamed asset menu and updated host page layout.

### Task 4: Add Network Device Backend CRUD And Connectivity Support

**Files:**
- Create: `api/api/cmdb/model/cmdbDevice.go`
- Create: `api/api/cmdb/dao/cmdbDevice.go`
- Create: `api/api/cmdb/service/cmdbDevice.go`
- Create: `api/api/cmdb/controller/cmdbDevice.go`
- Modify: `api/router/cmdb/cmdb.go`
- Test: `api/api/cmdb/service/cmdbDevice_test.go`

- [ ] **Step 1: Write the failing backend test**

```go
package service

import "testing"

func TestNormalizeCreateCmdbDeviceRequestAppliesProtocolDefaults(t *testing.T) {
	req := CreateCmdbDeviceDto{
		Name:       "core-switch-01",
		Address:    "10.0.0.21",
		Platform:   "Cisco IOS",
		DeviceType: "switch",
	}

	got := normalizeCreateDeviceRequest(req)

	if got.Address != "10.0.0.21" {
		t.Fatalf("expected address to be preserved, got %q", got.Address)
	}
	if got.SSHPort != 22 {
		t.Fatalf("expected default ssh port 22, got %d", got.SSHPort)
	}
}

func TestBuildDeviceConnectivityPrefersSSHThenTelnetThenWeb(t *testing.T) {
	device := CmdbDevice{Address: "10.0.0.22", SSHPort: 22, TelnetPort: 23, WebURL: "http://10.0.0.22"}
	target := buildDeviceConnectivityTarget(device)
	if target.Protocol != "ssh" {
		t.Fatalf("expected ssh target, got %s", target.Protocol)
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `cd api && go test ./api/cmdb/service -run Test.*CmdbDevice.* -v`

Expected: FAIL because the device model, service, and helper functions do not exist yet.

- [ ] **Step 3: Implement the new network-device backend slice**

```go
type CmdbDevice struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Platform    string `json:"platform"`
	GroupID     uint   `json:"groupId"`
	AccountID   uint   `json:"accountId"`
	ProtocolGrp string `json:"protocolGroup"`
	Tags        string `json:"tags"`
	IsActive    bool   `json:"isActive"`
	Remark      string `json:"remark"`
	DeviceType  string `json:"deviceType"`
	SSHPort     int    `json:"sshPort"`
	TelnetPort  int    `json:"telnetPort"`
	WebURL      string `json:"webUrl"`
}
```

```go
router.POST("/cmdb/device", controller.NewCmdbDeviceController().Create)
router.PUT("/cmdb/device", controller.NewCmdbDeviceController().Update)
router.DELETE("/cmdb/device", controller.NewCmdbDeviceController().Delete)
router.GET("/cmdb/devicelist", controller.NewCmdbDeviceController().List)
router.GET("/cmdb/device/info", controller.NewCmdbDeviceController().Detail)
router.POST("/cmdb/device/connectivity", controller.NewCmdbDeviceController().BatchConnectivity)
```

- [ ] **Step 4: Run the test to verify it passes**

Run: `cd api && go test ./api/cmdb/service -run Test.*CmdbDevice.* -v`

Expected: PASS with both device normalization and connectivity-target tests green.

- [ ] **Step 5: Run backend compilation for the new slice**

Run: `cd api && go test ./api/cmdb/... -v`

Expected: PASS without route-registration or type errors.

- [ ] **Step 6: Commit the device backend batch**

```bash
git add api/api/cmdb/model/cmdbDevice.go api/api/cmdb/dao/cmdbDevice.go api/api/cmdb/service/cmdbDevice.go api/api/cmdb/controller/cmdbDevice.go api/api/cmdb/service/cmdbDevice_test.go api/router/cmdb/cmdb.go
git commit -m "feat: add network device cmdb backend"
```

### Task 5: Build The Network Device Frontend Page

**Files:**
- Modify: `web/src/router/cmdb.js`
- Modify: `web/src/api/cmdb.js`
- Modify: `web/src/utils/cmdbAssetPresentation.mjs`
- Create: `web/src/views/cmdb/cmdbDevice.vue`
- Create: `web/src/views/cmdb/Device/CmdbDeviceTable.vue`
- Create: `web/src/views/cmdb/Device/CreateDevice.vue`
- Create: `web/src/views/cmdb/Device/EditDevice.vue`
- Test: `web/tests/cmdb-device-presentation.test.mjs`

- [ ] **Step 1: Write the failing network-device presentation test**

```js
import assert from 'node:assert/strict'
import { mapDeviceRowToAssetRow, createDeviceAssetFormModel } from '../src/utils/cmdbAssetPresentation.mjs'

const row = mapDeviceRowToAssetRow({
  id: 3,
  name: 'core-switch-01',
  address: '10.0.0.21',
  accountName: 'netops',
  platform: 'Cisco IOS',
  reachable: true
})

assert.equal(row.name, 'core-switch-01')
assert.equal(row.address, '10.0.0.21')
assert.equal(row.account, 'netops')
assert.equal(row.platform, 'Cisco IOS')
assert.equal(row.connectivity.text, '成功')

const form = createDeviceAssetFormModel()
assert.equal(form.name, '')
assert.equal(form.address, '')

console.log('cmdb device presentation tests passed')
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `node web/tests/cmdb-device-presentation.test.mjs`

Expected: FAIL because device mapping helpers and device page files do not exist yet.

- [ ] **Step 3: Implement the frontend request layer and asset page**

```js
// web/src/api/cmdb.js
createDevice(data) { return request({ url: 'cmdb/device', method: 'post', data }) }
updateDevice(data) { return request({ url: 'cmdb/device', method: 'put', data }) }
deleteDevice(id) { return request({ url: 'cmdb/device', method: 'delete', data: { id } }) }
listDevices(params) { return request({ url: 'cmdb/devicelist', method: 'get', params }) }
getDevice(id) { return request({ url: 'cmdb/device/info', method: 'get', params: { id } }) }
testDeviceConnectivity(ids) { return request({ url: 'cmdb/device/connectivity', method: 'post', data: { deviceIds: ids } }) }
```

```vue
<!-- web/src/views/cmdb/cmdbDevice.vue -->
<CmdbGroup ... />
<CmdbDeviceTable
  :device-list="deviceList"
  :loading="loading"
  @edit-device="showEditDialog"
  @delete-device="handleDelete"
  @connect-device="handleConnect"
/>
```

- [ ] **Step 4: Run the test to verify it passes**

Run: `node web/tests/cmdb-device-presentation.test.mjs`

Expected: PASS with `cmdb device presentation tests passed`

- [ ] **Step 5: Run the frontend build for device-page verification**

Run: `cd web && npm run build`

Expected: PASS and include the new `/cmdb/device` route in the built app.

- [ ] **Step 6: Commit and sync the network-device batch**

```bash
git add web/src/router/cmdb.js web/src/api/cmdb.js web/src/utils/cmdbAssetPresentation.mjs web/src/views/cmdb/cmdbDevice.vue web/src/views/cmdb/Device/CmdbDeviceTable.vue web/src/views/cmdb/Device/CreateDevice.vue web/src/views/cmdb/Device/EditDevice.vue web/tests/cmdb-device-presentation.test.mjs
git commit -m "feat: add network device asset page"
```

```powershell
cd api
$env:GOOS='linux'
$env:GOARCH='amd64'
go build -o ..\tmp\opsnexus-api-linux-amd64-cmdb-device .
cd ..
tar -czf tmp\web-dist-cmdb-device.tar.gz -C web dist
scp tmp\opsnexus-api-linux-amd64-cmdb-device root@10.0.0.200:/tmp/opsnexus-api-linux-amd64-cmdb-device
scp tmp\web-dist-cmdb-device.tar.gz root@10.0.0.200:/tmp/web-dist-cmdb-device.tar.gz
ssh root@10.0.0.200 "install -m 755 /tmp/opsnexus-api-linux-amd64-cmdb-device /opt/opsnexus-remote-test/opsnexus-api-linux-amd64 && rm -rf /opt/opsnexus-remote-test/web-dist/* && tar -xzf /tmp/web-dist-cmdb-device.tar.gz -C /opt/opsnexus-remote-test/web-dist --strip-components=1 && systemctl restart opsnexus-api.service && docker restart opsnexus-web"
```

Expected: remote `http://10.0.0.200:8080/` exposes the new `网络设备` page and CRUD/API requests succeed.

### Task 6: Expand The Database Backend Into A True Asset Model

**Files:**
- Modify: `api/api/cmdb/model/cmdbSQL.go`
- Modify: `api/api/cmdb/dao/cmdbSQL.go`
- Modify: `api/api/cmdb/service/cmdbSQL.go`
- Modify: `api/api/cmdb/controller/cmdbSQL.go`
- Modify: `api/api/cmdb/controller/cmdbSQLRecord.go`
- Modify: `api/api/cmdb/service/sqlWorkOrder.go`
- Test: `api/api/cmdb/service/cmdbSQL_asset_test.go`

- [ ] **Step 1: Write the failing database-asset test**

```go
package service

import "testing"

func TestNormalizeDatabaseAssetPreservesCompatibility(t *testing.T) {
	req := DatabaseAssetPayload{
		Name:            "orders-db",
		Address:         "10.0.0.31:3306",
		Platform:        "MySQL",
		DefaultDatabase: "orders",
		AccountID:       9,
		GroupID:         2,
		IsActive:        true,
	}

	got := normalizeDatabaseAssetPayload(req)

	if got.Name != "orders-db" {
		t.Fatalf("expected name to survive normalization")
	}
	if got.DefaultDatabase != "orders" {
		t.Fatalf("expected default database to be preserved")
	}
	if !got.IsActive {
		t.Fatalf("expected active flag to be true")
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `cd api && go test ./api/cmdb/service -run TestNormalizeDatabaseAssetPreservesCompatibility -v`

Expected: FAIL because the expanded database asset payload/normalizer does not exist yet.

- [ ] **Step 3: Implement the expanded database asset model and compatibility mapping**

```go
type CmdbSQL struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Type            int    `json:"type"`
	Address         string `json:"address"`
	Platform        string `json:"platform"`
	DefaultDatabase string `json:"defaultDatabase"`
	AccountID       uint   `json:"accountId"`
	GroupID         uint   `json:"groupId"`
	ProtocolGroup   string `json:"protocolGroup"`
	Tags            string `json:"tags"`
	IsActive        bool   `json:"isActive"`
	Description     string `json:"description"`
}
```

```go
func normalizeDatabaseAssetPayload(req DatabaseAssetPayload) DatabaseAssetPayload {
	if strings.TrimSpace(req.Platform) == "" {
		req.Platform = databaseTypeText(req.Type)
	}
	if strings.TrimSpace(req.Address) == "" {
		req.Address = buildAddressFromAccount(req.AccountID)
	}
	return req
}
```

- [ ] **Step 4: Run the test to verify it passes**

Run: `cd api && go test ./api/cmdb/service -run TestNormalizeDatabaseAssetPreservesCompatibility -v`

Expected: PASS with the new asset normalization logic green.

- [ ] **Step 5: Run database-related backend regression tests**

Run: `cd api && go test ./api/cmdb/... -v`

Expected: PASS and existing SQL record / SQL work-order code still compiles.

- [ ] **Step 6: Commit the database backend batch**

```bash
git add api/api/cmdb/model/cmdbSQL.go api/api/cmdb/dao/cmdbSQL.go api/api/cmdb/service/cmdbSQL.go api/api/cmdb/controller/cmdbSQL.go api/api/cmdb/controller/cmdbSQLRecord.go api/api/cmdb/service/sqlWorkOrder.go api/api/cmdb/service/cmdbSQL_asset_test.go
git commit -m "feat: expand database assets for jumpserver alignment"
```

### Task 7: Rebuild The Database Frontend Page Around Asset Semantics

**Files:**
- Modify: `web/src/api/cmdb.js`
- Modify: `web/src/utils/cmdbPresentation.mjs`
- Modify: `web/src/utils/cmdbAssetPresentation.mjs`
- Modify: `web/src/views/cmdb/cmdbDB.vue`
- Modify: `web/src/views/cmdb/DBdetails.vue`
- Modify: `web/src/views/cmdb/SQLWorkOrderCenter.vue`
- Test: `web/tests/cmdb-database-asset-presentation.test.mjs`

- [ ] **Step 1: Write the failing database asset presentation test**

```js
import assert from 'node:assert/strict'
import { mapDatabaseRowToAssetRow, createDatabaseAssetFormModel } from '../src/utils/cmdbAssetPresentation.mjs'

const row = mapDatabaseRowToAssetRow({
  id: 11,
  name: 'orders-db',
  address: '10.0.0.31:3306',
  accountName: 'dba_admin',
  platform: 'MySQL',
  defaultDatabase: 'orders',
  isActive: true
})

assert.equal(row.name, 'orders-db')
assert.equal(row.address, '10.0.0.31:3306')
assert.equal(row.account, 'dba_admin')
assert.equal(row.platform, 'MySQL')
assert.equal(row.defaultDatabase, 'orders')

const form = createDatabaseAssetFormModel()
assert.equal(form.defaultDatabase, '')

console.log('cmdb database asset presentation tests passed')
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `node web/tests/cmdb-database-asset-presentation.test.mjs`

Expected: FAIL because the database page still uses the old record-style field shape.

- [ ] **Step 3: Implement the database asset UI alignment**

```vue
<!-- cmdbDB.vue -->
<el-table-column prop="name" label="名称" min-width="180" />
<el-table-column prop="address" label="地址" min-width="180" />
<el-table-column prop="account" label="账号" min-width="140" />
<el-table-column prop="platform" label="平台" min-width="120" />
<el-table-column label="连接性" min-width="120">
  <template #default="{ row }">
    <el-tag :type="row.connectivity.type">{{ row.connectivity.text }}</el-tag>
  </template>
</el-table-column>
<el-table-column label="操作" width="220">
  <!-- 编辑 / 数据库操作 / SQL工单 / 删除 -->
</el-table-column>
```

```js
// cmdbAssetPresentation.mjs
export function mapDatabaseRowToAssetRow(item) {
  return {
    id: item.id,
    name: item.name,
    address: item.address || '',
    account: item.accountAlias || item.accountName || '',
    platform: item.platform || getCmdbDatabaseTypeLabel(item.type),
    defaultDatabase: item.defaultDatabase || '',
    connectivity: buildAssetConnectivityTag({ reachable: item.isActive !== false })
  }
}
```

- [ ] **Step 4: Run the test to verify it passes**

Run: `node web/tests/cmdb-database-asset-presentation.test.mjs`

Expected: PASS with `cmdb database asset presentation tests passed`

- [ ] **Step 5: Run the full frontend build**

Run: `cd web && npm run build`

Expected: PASS and the rebuilt database asset page compiles with `DBdetails` and `SQLWorkOrderCenter`.

- [ ] **Step 6: Commit and sync the database batch**

```bash
git add web/src/api/cmdb.js web/src/utils/cmdbPresentation.mjs web/src/utils/cmdbAssetPresentation.mjs web/src/views/cmdb/cmdbDB.vue web/src/views/cmdb/DBdetails.vue web/src/views/cmdb/SQLWorkOrderCenter.vue web/tests/cmdb-database-asset-presentation.test.mjs
git commit -m "feat: align database assets with jumpserver"
```

```powershell
cd api
$env:GOOS='linux'
$env:GOARCH='amd64'
go build -o ..\tmp\opsnexus-api-linux-amd64-cmdb-database .
cd ..
tar -czf tmp\web-dist-cmdb-database.tar.gz -C web dist
scp tmp\opsnexus-api-linux-amd64-cmdb-database root@10.0.0.200:/tmp/opsnexus-api-linux-amd64-cmdb-database
scp tmp\web-dist-cmdb-database.tar.gz root@10.0.0.200:/tmp/web-dist-cmdb-database.tar.gz
ssh root@10.0.0.200 "install -m 755 /tmp/opsnexus-api-linux-amd64-cmdb-database /opt/opsnexus-remote-test/opsnexus-api-linux-amd64 && rm -rf /opt/opsnexus-remote-test/web-dist/* && tar -xzf /tmp/web-dist-cmdb-database.tar.gz -C /opt/opsnexus-remote-test/web-dist --strip-components=1 && systemctl restart opsnexus-api.service && docker restart opsnexus-web"
```

Expected: remote `http://10.0.0.200:8080/` shows the rebuilt `数据库` page and both database-detail and SQL-work-order entrypoints remain usable.

### Task 8: Run Integrated Verification And Final Smoke Checks

**Files:**
- Modify: none unless verification uncovers regressions
- Test: existing and newly added tests

- [ ] **Step 1: Run all targeted frontend tests**

Run:

```bash
node web/tests/cmdb-assets-route-meta.test.mjs
node web/tests/cmdb-asset-presentation.test.mjs
node web/tests/cmdb-host-asset-presentation.test.mjs
node web/tests/cmdb-device-presentation.test.mjs
node web/tests/cmdb-database-asset-presentation.test.mjs
```

Expected: every test prints its `... tests passed` message and exits `0`.

- [ ] **Step 2: Run all targeted backend tests**

Run:

```bash
cd api
go test ./api/cmdb/... -v
```

Expected: PASS across CMDB packages with the new device and database asset tests included.

- [ ] **Step 3: Run full backend and frontend verification**

Run:

```bash
cd api && go test ./...
cd web && npm run build
```

Expected:
- backend: PASS
- frontend: PASS

- [ ] **Step 4: Perform manual functional smoke checks**

Check:

```text
1. 资产管理菜单中显示 主机 / 网络设备 / 数据库
2. 主机页显示 名称 / 地址 / 账号 / 平台 / 连接性 / 操作
3. 网络设备页可完成新增、编辑、删除、连通性测试
4. 数据库页可完成新增、编辑、删除，并能进入数据库操作与 SQL工单
5. JumpServer 10.0.0.201 对照字段差异已收敛
```

Expected: all five checks pass locally and on `10.0.0.200`.

- [ ] **Step 5: Final commit**

```bash
git add web api docs/superpowers/specs/2026-03-26-cmdb-assets-jumpserver-alignment-design.md docs/superpowers/plans/2026-03-26-cmdb-assets-jumpserver-alignment.md
git commit -m "feat: align cmdb assets with jumpserver"
```

## Review Note

- The writing-plans skill recommends a reviewer subagent after the plan is written.
- Subagent delegation is not user-authorized in this session, so this plan must be reviewed inline by rereading the plan document before execution.
