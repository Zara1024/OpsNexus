# K8s Top-Level Lint Batch 2 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove another low-risk batch of unused-symbol lint debt from the top-level K8s pages without changing user-visible behavior.

**Architecture:** Keep using targeted `eslint` runs as the red-green loop. This batch only removes unused imports, refs, router handles, and helper functions from pages where the lint report shows dead code rather than broken references.

**Tech Stack:** Vue 3 SFCs, Element Plus, ESLint, Vue CLI

---

### Task 1: Clean dead code from `k8s-config.vue`

**Files:**
- Modify: `web/src/views/K8s/k8s-config.vue`

- [ ] **Step 1: Run the failing lint check for `k8s-config.vue`**
  Run: `npx eslint src/views/K8s/k8s-config.vue`
- [ ] **Step 2: Remove the unused icon imports**
- [ ] **Step 3: Remove the unused `currentResourceName`, `currentResourceType`, and `resetSearch` locals**
- [ ] **Step 4: Re-run `npx eslint src/views/K8s/k8s-config.vue` and verify it passes**

### Task 2: Clean dead code from `k8s-namespace.vue`

**Files:**
- Modify: `web/src/views/K8s/k8s-namespace.vue`

- [ ] **Step 1: Run the failing lint check for `k8s-namespace.vue`**
  Run: `npx eslint src/views/K8s/k8s-namespace.vue`
- [ ] **Step 2: Remove the unused `computed` import, icon imports, and router handle**
- [ ] **Step 3: Remove the unused `getLabelCount` helper**
- [ ] **Step 4: Re-run `npx eslint src/views/K8s/k8s-namespace.vue` and verify it passes**

### Task 3: Clean dead code from `k8s-network.vue`

**Files:**
- Modify: `web/src/views/K8s/k8s-network.vue`

- [ ] **Step 1: Run the failing lint check for `k8s-network.vue`**
  Run: `npx eslint src/views/K8s/k8s-network.vue`
- [ ] **Step 2: Remove the unused icon imports plus the unused router and route handles**
- [ ] **Step 3: Remove the unused `handleViewIngressEvents`, `resetServiceForm`, and `resetIngressForm` helpers**
- [ ] **Step 4: Re-run `npx eslint src/views/K8s/k8s-network.vue` and verify it passes**

### Task 4: Clean dead code from `k8s-nodes.vue`

**Files:**
- Modify: `web/src/views/K8s/k8s-nodes.vue`

- [ ] **Step 1: Run the failing lint check for `k8s-nodes.vue`**
  Run: `npx eslint src/views/K8s/k8s-nodes.vue`
- [ ] **Step 2: Remove the unused icon imports**
- [ ] **Step 3: Remove the unused `getNodeStatus` and `getNodeRole` helpers**
- [ ] **Step 4: Re-run `npx eslint src/views/K8s/k8s-nodes.vue` and verify it passes**

### Task 5: Verify the second K8s cleanup batch

**Files:**
- Modify: `web/src/views/K8s/k8s-config.vue`
- Modify: `web/src/views/K8s/k8s-namespace.vue`
- Modify: `web/src/views/K8s/k8s-network.vue`
- Modify: `web/src/views/K8s/k8s-nodes.vue`

- [ ] **Step 1: Run the combined lint check for the four edited files**
  Run: `npx eslint src/views/K8s/k8s-config.vue src/views/K8s/k8s-namespace.vue src/views/K8s/k8s-network.vue src/views/K8s/k8s-nodes.vue`
- [ ] **Step 2: Run `npm run build` to confirm the cleanup does not break the app**
- [ ] **Step 3: Re-run full `npm run lint` and record the new remaining error count**

## Notes

- This workspace is not an initialized git worktree, so the plan is saved locally only.
- Subagent review is skipped because delegation is not user-authorized in this session; review is handled inline.
