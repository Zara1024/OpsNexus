# K8s Workloads Lint Batch 6 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove the large cluster of unused-symbol lint debt from `k8s-workloads.vue` without changing active behavior.

**Architecture:** Use `eslint` as the red-green loop for the single file. The cleanup is limited to imports, dialog flags, helper functions, and placeholder functions already proven unused by lint after earlier component extraction work.

**Tech Stack:** Vue 3 SFCs, Element Plus, ESLint, Vue CLI

---

### Task 1: Clean unused declarations from `k8s-workloads.vue`

**Files:**
- Modify: `web/src/views/K8s/k8s-workloads.vue`

- [ ] **Step 1: Run the failing lint check for `k8s-workloads.vue`**
  Run: `npx eslint src/views/K8s/k8s-workloads.vue`
- [ ] **Step 2: Remove the unused icon and component imports**
- [ ] **Step 3: Remove the unused dialog refs and arrays such as `workloadTypeOptions` and `yamlDialogVisible`**
- [ ] **Step 4: Remove the unused helper functions flagged by lint**
- [ ] **Step 5: Re-run `npx eslint src/views/K8s/k8s-workloads.vue` and verify it passes**

### Task 2: Verify the workloads cleanup batch

**Files:**
- Modify: `web/src/views/K8s/k8s-workloads.vue`

- [ ] **Step 1: Run `npm run build`**
- [ ] **Step 2: Run full `npm run lint` and record the new remaining error count**

## Notes

- This batch only targets declarations already flagged as unused by lint.
- This workspace is not an initialized git worktree, so the plan is saved locally only.
