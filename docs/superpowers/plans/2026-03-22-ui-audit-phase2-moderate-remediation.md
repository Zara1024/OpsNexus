# OpsNexus UI Audit Phase 2 Moderate Remediation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Reduce the remaining moderate table-layout issues identified in the shared-environment audit while preserving primary user flows and verifying each page immediately after change.

**Architecture:** Continue using route-level layout threshold checks as regression gates. For each moderate page, shrink or remove low-priority columns from the first viewport, preserve primary actions, and fall back to details/dialogs for secondary information.

**Tech Stack:** Vue 3, Element Plus, local Vue dev server, Playwright-based smoke verification, existing `layout-threshold-check.mjs`

---

### Task 1: Fix SQL Work Order Center

**Files:**
- Modify: `web/src/views/cmdb/SQLWorkOrderCenter.vue`

- [ ] **Step 1: Run RED verification**
Run: `node scripts/layout-threshold-check.mjs --route /cmdb/sql-work-orders --base-url http://localhost:8080 --max-wide-overflow 300 --max-text-overflow 80`
Expected: FAIL

- [ ] **Step 2: Reduce low-priority list columns**
Keep order number, title, database, type, risk, status, and primary actions visible. Demote applicant/approver/create time to detail drawer if needed.

- [ ] **Step 3: Run GREEN verification**

- [ ] **Step 4: Run click-path smoke**
Verify detail drawer and approve/reject/execute actions still work.

### Task 2: Fix K8s Node Management

**Files:**
- Modify: `web/src/views/K8s/k8s-nodes.vue`
- Modify: `web/src/views/K8s/nodes/NodesTable.vue`

- [ ] **Step 1: Run RED verification**
Run: `node scripts/layout-threshold-check.mjs --route /k8s/node --base-url http://localhost:8080 --max-wide-overflow 300 --max-text-overflow 80`
Expected: FAIL

- [ ] **Step 2: Reduce first-screen width pressure**
Keep node name, role, status, allocatable summary, and primary actions visible. Move secondary resource detail or metadata to dialogs.

- [ ] **Step 3: Run GREEN verification**

- [ ] **Step 4: Run click-path smoke**
Verify node detail, taint/label management, and cordon actions still render.

### Task 3: Fix K8s Network Management

**Files:**
- Modify: `web/src/views/K8s/k8s-network.vue`

- [ ] **Step 1: Run RED verification**
Run: `node scripts/layout-threshold-check.mjs --route /k8s/network --base-url http://localhost:8080 --max-wide-overflow 300 --max-text-overflow 80`
Expected: FAIL

- [ ] **Step 2: Reduce service/ingress table width**
Keep name, core endpoint/host summary, type/status, and primary actions visible. Collapse verbose backend detail into dialog/tooltip.

- [ ] **Step 3: Run GREEN verification**

- [ ] **Step 4: Run click-path smoke**
Verify detail, YAML, events, and create/edit actions still work.

### Task 4: Fix Task Job List

**Files:**
- Modify: `web/src/views/task/TaskJob.vue`

- [ ] **Step 1: Run RED verification**
Run: `node scripts/layout-threshold-check.mjs --route /task/job --base-url http://localhost:8080 --max-wide-overflow 300 --max-text-overflow 80`
Expected: FAIL

- [ ] **Step 2: Compress schedule and remark columns**

- [ ] **Step 3: Run GREEN verification**

- [ ] **Step 4: Run click-path smoke**
Verify create task flow and start/delete operations remain usable.

### Task 5: Fix Task Template List

**Files:**
- Modify: `web/src/views/task/TaskTemplate.vue`

- [ ] **Step 1: Run RED verification**
Run: `node scripts/layout-threshold-check.mjs --route /task/template --base-url http://localhost:8080 --max-wide-overflow 300 --max-text-overflow 80`
Expected: FAIL

- [ ] **Step 2: Compress preview and metadata columns**

- [ ] **Step 3: Run GREEN verification**

- [ ] **Step 4: Run click-path smoke**
Verify view, edit, and delete flows.

### Task 6: Fix Task Config Center

**Files:**
- Modify: `web/src/views/task/TaskConfig.vue`

- [ ] **Step 1: Run RED verification**
Run: `node scripts/layout-threshold-check.mjs --route /task/config --base-url http://localhost:8080 --max-wide-overflow 300 --max-text-overflow 80`
Expected: FAIL

- [ ] **Step 2: Reduce preview and metadata width**

- [ ] **Step 3: Run GREEN verification**

- [ ] **Step 4: Run click-path smoke**
Verify tab switch, view, edit, and delete actions.

### Task 7: Fix Agent Management

**Files:**
- Modify: `web/src/views/Tools/Agent.vue`

- [ ] **Step 1: Run RED verification**
Run: `node scripts/layout-threshold-check.mjs --route /ops/agent --base-url http://localhost:8080 --max-wide-overflow 300 --max-text-overflow 80`
Expected: FAIL

- [ ] **Step 2: Compress host, progress, and operation columns**

- [ ] **Step 3: Run GREEN verification**

- [ ] **Step 4: Run click-path smoke**
Verify deploy, restart, uninstall, and detail dialog.

### Task 8: Fix Login Log

**Files:**
- Modify: `web/src/views/monitor/LoginLog.vue`

- [ ] **Step 1: Run RED verification**
Run: `node scripts/layout-threshold-check.mjs --route /monitor/loginlog --base-url http://localhost:8080 --max-wide-overflow 300 --max-text-overflow 80`
Expected: FAIL

- [ ] **Step 2: Reduce browser/OS density and keep audit essentials**

- [ ] **Step 3: Run GREEN verification**

- [ ] **Step 4: Run click-path smoke**
Verify search, batch delete, and clean log operations.

### Task 9: Fix Alert Notify

**Files:**
- Modify: `web/src/views/monitor/Alarm-notify.vue`

- [ ] **Step 1: Run RED verification**
Run: `node scripts/layout-threshold-check.mjs --route /monitor/alert-notify --base-url http://localhost:8080 --max-wide-overflow 300 --max-text-overflow 80`
Expected: FAIL

- [ ] **Step 2: Reduce robot list and notify log table width**

- [ ] **Step 3: Run GREEN verification**

- [ ] **Step 4: Run click-path smoke**
Verify create/edit/test robot dialogs and notify log list.

### Task 10: Fix System Admin

**Files:**
- Modify: `web/src/views/system/Admin.vue`

- [ ] **Step 1: Run RED verification**
Run: `node scripts/layout-threshold-check.mjs --route /system/admin --base-url http://localhost:8080 --max-wide-overflow 300 --max-text-overflow 80`
Expected: FAIL

- [ ] **Step 2: Reduce table metadata and keep account governance essentials**

- [ ] **Step 3: Run GREEN verification**

- [ ] **Step 4: Run click-path smoke**
Verify add/edit/reset-password/delete flows and dept filtering.

### Task 11: Fix System Menu

**Files:**
- Modify: `web/src/views/system/Menu.vue`

- [ ] **Step 1: Run RED verification**
Run: `node scripts/layout-threshold-check.mjs --route /system/menu --base-url http://localhost:8080 --max-wide-overflow 300 --max-text-overflow 80`
Expected: FAIL

- [ ] **Step 2: Compress icon/time columns and keep hierarchy essentials**

- [ ] **Step 3: Run GREEN verification**

- [ ] **Step 4: Run click-path smoke**
Verify create/edit/copy/delete actions and tree expand state.

### Task 12: Fix Common Credential List

**Files:**
- Modify: `web/src/views/configcenter/accountauth.vue`

- [ ] **Step 1: Run RED verification**
Run: `node scripts/layout-threshold-check.mjs --route /config/accountauth --base-url http://localhost:8080 --max-wide-overflow 300 --max-text-overflow 80`
Expected: FAIL

- [ ] **Step 2: Reduce host/time/remark pressure and keep alias/type/actions**

- [ ] **Step 3: Run GREEN verification**

- [ ] **Step 4: Run click-path smoke**
Verify create/edit/decrypt/delete flows.

### Task 13: Final Moderate Regression

**Files:**
- Reuse: `scripts/layout-threshold-check.mjs`

- [ ] **Step 1: Re-run all moderate route checks**

Routes:
- `/cmdb/sql-work-orders`
- `/k8s/node`
- `/k8s/network`
- `/task/job`
- `/task/template`
- `/task/config`
- `/ops/agent`
- `/monitor/loginlog`
- `/monitor/alert-notify`
- `/system/admin`
- `/system/menu`
- `/config/accountauth`

- [ ] **Step 2: Re-run focused browser smoke for modified pages**

- [ ] **Step 3: Refresh UI audit docs if severity distribution changes materially**
