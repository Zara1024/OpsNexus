# CMDB Host Delete Actions Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the host-page terminal entry with delete actions, remove the entire `更多` dropdown from the compact host table, and support both single-host and batch-host deletion.

**Architecture:** Keep the backend contract unchanged and reuse the existing single-delete API for both row-level and batch deletion. Limit changes to the host page container, the compact host table component, and focused regression tests that lock the new action layout and delete entry points.

**Tech Stack:** Vue 3 options API, Element Plus, existing CMDB API client, node-based source assertions for regression coverage.

---

### Task 1: Lock the new delete-entry UI with failing tests

**Files:**
- Modify: `web/tests/cmdb-host-compact-table-layout.test.mjs`

- [ ] **Step 1: Write the failing test**

Add assertions that:
- the compact table no longer renders `table-operation__more`
- the compact table renders a row-level `删除主机` action

- [ ] **Step 2: Run test to verify it fails**

Run: `node .\tests\cmdb-host-compact-table-layout.test.mjs`
Expected: FAIL because the compact table still contains the `更多` dropdown and does not expose inline `删除主机`.

- [ ] **Step 3: Write minimal implementation**

Remove the dropdown from the compact table template and expose the row-level delete button directly.

- [ ] **Step 4: Run test to verify it passes**

Run: `node .\tests\cmdb-host-compact-table-layout.test.mjs`
Expected: PASS

### Task 2: Replace host-page terminal entry with batch delete

**Files:**
- Modify: `web/src/views/cmdb/cmdbHost.vue`

- [ ] **Step 1: Write the failing test**

Add source assertions that:
- the host toolbar button label is `删除主机`
- the toolbar click handler no longer points to `handleHostSSH`
- the batch toolbar exposes a delete action hook

- [ ] **Step 2: Run test to verify it fails**

Run the targeted node-based test file for the host page.
Expected: FAIL because the toolbar still shows `终端`.

- [ ] **Step 3: Write minimal implementation**

Update the toolbar button to:
- use delete permission
- call a new batch delete handler
- disable itself when no hosts are selected

Add a batch delete method that:
- confirms the selected host count
- calls `deleteCmdbHost(id)` for each selected host
- reports success/failure summary
- refreshes the host list afterward

- [ ] **Step 4: Run test to verify it passes**

Run the targeted node-based test file for the host page.
Expected: PASS

### Task 3: Keep single delete and remove dead `更多` actions

**Files:**
- Modify: `web/src/views/cmdb/Host/CmdbHostTableCompact.vue`
- Modify: `web/src/views/cmdb/cmdbHost.vue`

- [ ] **Step 1: Write the failing test**

Add assertions that the compact table only exposes `编辑 / 监控 / 删除主机` in the action column and does not reference `process`, `port`, `upload`, `command`, or `delete` dropdown commands anymore.

- [ ] **Step 2: Run test to verify it fails**

Run: `node .\tests\cmdb-host-compact-table-layout.test.mjs`
Expected: FAIL while old dropdown commands still exist.

- [ ] **Step 3: Write minimal implementation**

Remove the stale dropdown command handling and keep row-level delete wired to the existing `delete-host` event.

- [ ] **Step 4: Run test to verify it passes**

Run: `node .\tests\cmdb-host-compact-table-layout.test.mjs`
Expected: PASS

### Task 4: Verify the new delete flow and deploy

**Files:**
- Modify: `web/dist` (generated)
- Create: `tmp/web-dist-cmdb-host-delete-actions-20260328.tar.gz`

- [ ] **Step 1: Run focused regression tests**

Run: `node .\tests\cmdb-host-compact-table-layout.test.mjs`
Expected: PASS

- [ ] **Step 2: Run production build**

Run: `npm run build`
Expected: successful Vue production build with exit code `0`

- [ ] **Step 3: Verify in browser**

Use the local or remote OpsNexus browser flow to confirm:
- toolbar shows `删除主机`
- row actions show `编辑 / 监控 / 删除主机`
- `更多` is gone

- [ ] **Step 4: Deploy the verified frontend bundle**

Run the established remote sync flow to `10.0.0.200`:
- archive `web/dist`
- upload to `/tmp`
- extract into `/opt/opsnexus-remote-test/web-dist`
- restart `opsnexus-web`

- [ ] **Step 5: Verify the remote page**

Inspect `http://10.0.0.200:8080/cmdb/ecs`
Expected: the remote page matches the local verified delete-action layout.
