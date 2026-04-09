# OpsNexus UI Audit Phase 1 Remediation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Reduce severe table-layout issues identified in the shared-environment UI audit and verify each remediation step immediately after implementation.

**Architecture:** Use a route-level layout regression check to measure horizontal overflow before and after each change, then remediate the worst routes by compressing low-priority columns, consolidating visual density, and preserving core actions in the first viewport. Apply reusable patterns only when they reduce repeated page-specific fixes.

**Tech Stack:** Vue 3, Element Plus, local Vue dev server, Playwright-based smoke verification, existing audit artifacts

---

### Task 1: Add Route-Level Layout Verification

**Files:**
- Create: `scripts/layout-threshold-check.mjs`

- [ ] **Step 1: Write a failing route layout check**

Create a Playwright-based script that:
- loads a route on the local frontend
- reads the visible page layout
- calculates max horizontal overflow and max clipped-text overflow
- exits non-zero when the route exceeds passed thresholds

- [ ] **Step 2: Run the check against one known-bad route to verify it fails**

Run: `node scripts/layout-threshold-check.mjs --route /cmdb/ecs --base-url http://localhost:8080 --max-wide-overflow 300`
Expected: FAIL because `/cmdb/ecs` currently exceeds the threshold

- [ ] **Step 3: Keep the script as the per-change verification harness**

The same script will be used after each page remediation.

### Task 2: Fix CMDB Host Management Width Explosion

**Files:**
- Modify: `web/src/views/cmdb/Host/CmdbHostTable.vue`
- Modify: `web/src/views/cmdb/cmdbHost.vue`

- [ ] **Step 1: Run the route check and confirm RED**

Run: `node scripts/layout-threshold-check.mjs --route /cmdb/ecs --base-url http://localhost:8080 --max-wide-overflow 300`
Expected: FAIL

- [ ] **Step 2: Reduce first-screen width pressure**

Implement:
- consolidate CPU / memory / disk into one resource column
- compress action area into fewer primary controls plus overflow actions
- reduce duplicate icon-only columns that are better as secondary actions

- [ ] **Step 3: Run the route check and confirm GREEN**

Run the same command.
Expected: PASS

- [ ] **Step 4: Run a click-path smoke**

Verify host detail, edit, upload/command entry, and delete button visibility still work.

### Task 3: Fix K8s Workload List Width Explosion

**Files:**
- Modify: `web/src/views/K8s/k8s-workloads.vue`

- [ ] **Step 1: Run the route check and confirm RED**

Run: `node scripts/layout-threshold-check.mjs --route /k8s/workload --base-url http://localhost:8080 --max-wide-overflow 300`
Expected: FAIL

- [ ] **Step 2: Reduce low-priority columns and compress actions**

Implement:
- prioritize workload name, status, pod readiness, key resources, and primary operations
- move secondary metadata behind dialog/tooltip or compact summary blocks

- [ ] **Step 3: Run the route check and confirm GREEN**

Run the same command.
Expected: PASS

- [ ] **Step 4: Run a click-path smoke**

Verify name click, pod dialog, YAML dialog, scaling, and key operations still work.

### Task 4: Fix Application List Width Explosion

**Files:**
- Modify: `web/src/views/app/application.vue`

- [ ] **Step 1: Run the route check and confirm RED**

Run: `node scripts/layout-threshold-check.mjs --route /app/application --base-url http://localhost:8080 --max-wide-overflow 300`
Expected: FAIL

- [ ] **Step 2: Compress overview and table density**

Implement:
- keep app name, type/status, env coverage, and core actions visible
- demote low-priority detail to dialogs or secondary text

- [ ] **Step 3: Run the route check and confirm GREEN**

Run the same command.
Expected: PASS

- [ ] **Step 4: Run a click-path smoke**

Verify view, edit, delete, and Jenkins env management remain usable.

### Task 5: Fix Ansible Task List Width Explosion

**Files:**
- Modify: `web/src/views/task/TaskAnsible.vue`

- [ ] **Step 1: Run the route check and confirm RED**

Run: `node scripts/layout-threshold-check.mjs --route /task/ansible --base-url http://localhost:8080 --max-wide-overflow 300`
Expected: FAIL

- [ ] **Step 2: Compress table columns and preserve start/edit/delete path**

- [ ] **Step 3: Run the route check and confirm GREEN**

- [ ] **Step 4: Run a click-path smoke**

Verify start, edit, delete, and create dialog layout still work.

### Task 6: Fix Work Order List Width Explosion

**Files:**
- Modify: `web/src/views/work/Work-Order-list.vue`

- [ ] **Step 1: Run the route check and confirm RED**

Run: `node scripts/layout-threshold-check.mjs --route /work/orders --base-url http://localhost:8080 --max-wide-overflow 300`
Expected: FAIL

- [ ] **Step 2: Compress risk/detail/action presentation**

- [ ] **Step 3: Run the route check and confirm GREEN**

- [ ] **Step 4: Run a click-path smoke**

Verify detail, approve/reject/execute, and AI/knowledge links still work.

### Task 7: Fix Terminal Audit Width Explosion

**Files:**
- Modify: `web/src/views/monitor/Recording.vue`

- [ ] **Step 1: Run the route check and confirm RED**

Run: `node scripts/layout-threshold-check.mjs --route /monitor/recording --base-url http://localhost:8080 --max-wide-overflow 300`
Expected: FAIL

- [ ] **Step 2: Rebalance audit columns and replay/detail tabs**

- [ ] **Step 3: Run the route check and confirm GREEN**

- [ ] **Step 4: Run a click-path smoke**

Verify detail drawer, command tab, playback tab, and download button layout still work.

### Task 8: Fix DB Log Width Explosion

**Files:**
- Modify: `web/src/views/monitor/DBLog.vue`

- [ ] **Step 1: Run the route check and confirm RED**

Run: `node scripts/layout-threshold-check.mjs --route /monitor/dblog --base-url http://localhost:8080 --max-wide-overflow 300`
Expected: FAIL

- [ ] **Step 2: Compress table fields and batch-action bar**

- [ ] **Step 3: Run the route check and confirm GREEN**

- [ ] **Step 4: Run a click-path smoke**

Verify search, batch delete, clean log, and single delete remain usable.

### Task 9: Final Phase 1 Regression

**Files:**
- Reuse: `scripts/layout-threshold-check.mjs`

- [ ] **Step 1: Re-run all Phase 1 route checks**

Routes:
- `/cmdb/ecs`
- `/k8s/workload`
- `/app/application`
- `/task/ansible`
- `/work/orders`
- `/monitor/recording`
- `/monitor/dblog`

- [ ] **Step 2: Re-run focused browser smoke for changed pages**

Confirm each page still loads and primary interactions remain available.

- [ ] **Step 3: Update audit docs if needed**

If final status materially changes the severity distribution, refresh the generated docs.
