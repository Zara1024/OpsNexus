# K8s Cluster Lint Batch 1 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove the first batch of historical lint debt from the K8s cluster module without changing runtime behavior.

**Architecture:** Use targeted `eslint` runs as the red-green loop for each file. The cleanup is limited to stale imports, dead local helpers, and unused callback parameters left behind after the cluster module was split into table and dialog components.

**Tech Stack:** Vue 3 SFCs, Element Plus, ESLint, Vue CLI

---

### Task 1: Clean dead code from `ClusterTable.vue`

**Files:**
- Modify: `web/src/views/K8s/clusters/ClusterTable.vue`

- [ ] **Step 1: Run the failing lint check for `ClusterTable.vue`**
  Run: `npx eslint src/views/K8s/clusters/ClusterTable.vue`
- [ ] **Step 2: Remove the unused `Monitor` icon import**
- [ ] **Step 3: Remove the unused `handleKubectlTerminal` helper**
- [ ] **Step 4: Re-run `npx eslint src/views/K8s/clusters/ClusterTable.vue` and verify it passes**

### Task 2: Clean dead code from `K8sDetails.vue`

**Files:**
- Modify: `web/src/views/K8s/clusters/K8sDetails.vue`

- [ ] **Step 1: Run the failing lint check for `K8sDetails.vue`**
  Run: `npx eslint src/views/K8s/clusters/K8sDetails.vue`
- [ ] **Step 2: Remove the unused `computed`, `CpuFill`, and `MemoryStick` imports**
- [ ] **Step 3: Remove the unused `usedCores` and `usedMi` locals while keeping the existing CPU and memory calculations intact**
- [ ] **Step 4: Re-run `npx eslint src/views/K8s/clusters/K8sDetails.vue` and verify it passes**

### Task 3: Clean dead code from `k8s-clusters.vue`

**Files:**
- Modify: `web/src/views/K8s/k8s-clusters.vue`

- [ ] **Step 1: Run the failing lint check for `k8s-clusters.vue`**
  Run: `npx eslint src/views/K8s/k8s-clusters.vue`
- [ ] **Step 2: Remove the unused Vue, icon, API, component, and router imports**
- [ ] **Step 3: Remove the unused `row` parameters from the cluster refresh handlers**
- [ ] **Step 4: Re-run `npx eslint src/views/K8s/k8s-clusters.vue` and verify it passes**

### Task 4: Verify the first cluster-module cleanup batch

**Files:**
- Modify: `web/src/views/K8s/clusters/ClusterTable.vue`
- Modify: `web/src/views/K8s/clusters/K8sDetails.vue`
- Modify: `web/src/views/K8s/k8s-clusters.vue`

- [ ] **Step 1: Run the combined lint check for the three edited files**
  Run: `npx eslint src/views/K8s/clusters/ClusterTable.vue src/views/K8s/clusters/K8sDetails.vue src/views/K8s/k8s-clusters.vue`
- [ ] **Step 2: Run `npm run build` to verify the cleanup does not break the frontend build**
- [ ] **Step 3: Summarize the remaining lint debt outside this batch**

## Notes

- This workspace is not an initialized git worktree, so the plan can be saved locally but cannot be committed from here.
- Subagent review is skipped because delegation is not user-authorized in this session; review will be handled inline.
