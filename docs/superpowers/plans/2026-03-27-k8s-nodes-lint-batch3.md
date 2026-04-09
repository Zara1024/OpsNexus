# K8s Nodes Lint Batch 3 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove the next batch of low-risk lint debt from the K8s nodes monitoring/detail pages without changing behavior.

**Architecture:** Continue using targeted `eslint` runs as the red-green loop. This batch is limited to unused imports, unused emits/helpers, and dead locals in the nodes detail, monitoring, and table views.

**Tech Stack:** Vue 3 SFCs, Element Plus, ESLint, Vue CLI

---

### Task 1: Clean dead code from `NodeDetails.vue`

**Files:**
- Modify: `web/src/views/K8s/nodes/NodeDetails.vue`

- [ ] **Step 1: Run the failing lint check for `NodeDetails.vue`**
  Run: `npx eslint src/views/K8s/nodes/NodeDetails.vue`
- [ ] **Step 2: Remove the unused Element Plus icon imports**
- [ ] **Step 3: Remove the unused `capacity` locals in the CPU and memory helper functions**
- [ ] **Step 4: Re-run `npx eslint src/views/K8s/nodes/NodeDetails.vue` and verify it passes**

### Task 2: Clean dead code from `NodesMonitoring.vue`

**Files:**
- Modify: `web/src/views/K8s/nodes/NodesMonitoring.vue`

- [ ] **Step 1: Run the failing lint check for `NodesMonitoring.vue`**
  Run: `npx eslint src/views/K8s/nodes/NodesMonitoring.vue`
- [ ] **Step 2: Remove the unused `emit` declaration**
- [ ] **Step 3: Remove the unused `getProgressColor` helper**
- [ ] **Step 4: Re-run `npx eslint src/views/K8s/nodes/NodesMonitoring.vue` and verify it passes**

### Task 3: Clean dead code from `NodesTable.vue`

**Files:**
- Modify: `web/src/views/K8s/nodes/NodesTable.vue`

- [ ] **Step 1: Run the failing lint check for `NodesTable.vue`**
  Run: `npx eslint src/views/K8s/nodes/NodesTable.vue`
- [ ] **Step 2: Remove the unused `getLabelCount` and `viewNodeLabels` helpers**
- [ ] **Step 3: Re-run `npx eslint src/views/K8s/nodes/NodesTable.vue` and verify it passes**

### Task 4: Clean dead code from `k8s-monitoring.vue`

**Files:**
- Modify: `web/src/views/K8s/nodes/k8s-monitoring.vue`

- [ ] **Step 1: Run the failing lint check for `k8s-monitoring.vue`**
  Run: `npx eslint src/views/K8s/nodes/k8s-monitoring.vue`
- [ ] **Step 2: Remove the unused Vue imports, icon imports, router handle, and formatting helpers**
- [ ] **Step 3: Re-run `npx eslint src/views/K8s/nodes/k8s-monitoring.vue` and verify it passes**

### Task 5: Verify the third K8s cleanup batch

**Files:**
- Modify: `web/src/views/K8s/nodes/NodeDetails.vue`
- Modify: `web/src/views/K8s/nodes/NodesMonitoring.vue`
- Modify: `web/src/views/K8s/nodes/NodesTable.vue`
- Modify: `web/src/views/K8s/nodes/k8s-monitoring.vue`

- [ ] **Step 1: Run the combined lint check for the four edited files**
  Run: `npx eslint src/views/K8s/nodes/NodeDetails.vue src/views/K8s/nodes/NodesMonitoring.vue src/views/K8s/nodes/NodesTable.vue src/views/K8s/nodes/k8s-monitoring.vue`
- [ ] **Step 2: Run `npm run build` to confirm the cleanup does not break the app**
- [ ] **Step 3: Re-run full `npm run lint` and record the new remaining error count**

## Notes

- This workspace is not an initialized git worktree, so the plan is saved locally only.
- Subagent review is skipped because delegation is not user-authorized in this session; review is handled inline.
