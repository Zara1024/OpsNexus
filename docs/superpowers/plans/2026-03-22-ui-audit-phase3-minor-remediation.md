# OpsNexus UI Audit Phase 3 Minor Remediation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Eliminate the remaining minor layout issues from the UI audit, keeping first-screen readability high while preserving all existing workflows.

**Architecture:** Reuse the existing route-level layout threshold checker as the regression harness, but target smaller interventions: shave low-priority columns, reduce text overflow, and tighten action density only where the first viewport still shows avoidable overflow.

**Tech Stack:** Vue 3, Element Plus, local Vue dev server, Playwright-based layout smoke verification

---

### Task 1: Re-baseline all minor pages locally

**Files:**
- Reuse: `scripts/layout-threshold-check.mjs`

- [ ] **Step 1: Re-run all minor routes**

Routes:
- `/cmdb/db`
- `/k8s/list`
- `/k8s/namespace`
- `/k8s/storage`
- `/k8s/monitoring`
- `/app/quick-release`
- `/monitor/operator`
- `/monitor/alert-center`
- `/monitor/alert-history`
- `/system/role`
- `/system/post`
- `/config/keymanage`

- [ ] **Step 2: Split routes into “already green enough” and “still worth fixing”**

### Task 2: Fix Asset/Data Minor Pages

**Files:**
- Modify as needed:
  - `web/src/views/cmdb/cmdbDB.vue`
  - `web/src/views/system/Role.vue`
  - `web/src/views/system/Post.vue`
  - `web/src/views/configcenter/KeyManage.vue`

- [ ] **Step 1: Run RED verification for each page that still exceeds the target**

- [ ] **Step 2: Apply compact-table adjustments**

- [ ] **Step 3: Run GREEN verification after each page**

### Task 3: Fix K8s Minor Pages

**Files:**
- Modify as needed:
  - `web/src/views/K8s/k8s-clusters.vue`
  - `web/src/views/K8s/k8s-namespace.vue`
  - `web/src/views/K8s/k8s-storage.vue`
  - `web/src/views/K8s/nodes/k8s-monitoring.vue`

- [ ] **Step 1: Run RED verification per page**

- [ ] **Step 2: Remove low-priority columns from the first viewport**

- [ ] **Step 3: Run GREEN verification after each page**

### Task 4: Fix Service/Alert Minor Pages

**Files:**
- Modify as needed:
  - `web/src/views/app/app_quick_release.vue`
  - `web/src/views/monitor/Operator.vue`
  - `web/src/views/monitor/Alarm-rules.vue`
  - `web/src/views/monitor/alarm-history.vue`

- [ ] **Step 1: Run RED verification per page**

- [ ] **Step 2: Tighten columns and keep primary actions visible**

- [ ] **Step 3: Run GREEN verification after each page**

### Task 5: Final Minor Regression

**Files:**
- Reuse: `scripts/layout-threshold-check.mjs`

- [ ] **Step 1: Re-run all minor routes**

- [ ] **Step 2: Run a production build**

- [ ] **Step 3: Refresh audit docs if the local severity distribution has materially improved**
