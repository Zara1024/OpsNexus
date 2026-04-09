# K8s Pod Detail Lint Batch 8 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Clean the remaining lint debt in `k8s-pod.vue`, including the undefined owner-reference variables and search-loop lint issues.

**Architecture:** Use single-file `eslint` as the red-green loop. Fixes are limited to unused declarations, hoisting the owner reference variables to the correct scope, replacing `while(true)` loops with lint-safe equivalents, and removing one unreachable statement.

**Tech Stack:** Vue 3 SFC, Element Plus, ESLint, Vue CLI

---

### Task 1: Clean and fix `k8s-pod.vue`

**Files:**
- Modify: `web/src/views/K8s/pods/k8s-pod.vue`

- [ ] **Step 1: Run the failing lint check for `k8s-pod.vue`**
  Run: `npx eslint src/views/K8s/pods/k8s-pod.vue`
- [ ] **Step 2: Remove the unused imports and unused computed/helpers**
- [ ] **Step 3: Hoist `deploymentRef` and `replicaSetRef` so the later workload-pod lookup can use them safely**
- [ ] **Step 4: Replace the `while (true)` search loops with lint-safe loop constructs and remove unused local variables**
- [ ] **Step 5: Remove the unreachable statement and unused callback parameter**
- [ ] **Step 6: Re-run `npx eslint src/views/K8s/pods/k8s-pod.vue` and verify it passes**

### Task 2: Verify the pod-detail cleanup batch

**Files:**
- Modify: `web/src/views/K8s/pods/k8s-pod.vue`

- [ ] **Step 1: Run `npm run build`**
- [ ] **Step 2: Run full `npm run lint` and record the new remaining error count**

## Notes

- This batch intentionally keeps behavior unchanged; it only repairs lint issues and one real scope bug in the owner-reference lookup path.
- This workspace is not an initialized git worktree, so the plan is saved locally only.
