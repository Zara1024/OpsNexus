# Safe Global Lint Batch 10 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove the next set of low-risk lint errors and one known build warning without touching high-risk pages.

**Architecture:** Use targeted `eslint` and `npm run build` as the red-green loop. This batch only updates component names, removes unused imports/locals, and replaces the unsupported `Copy` icon in `k8s-pod.vue` with a supported icon export.

**Tech Stack:** Vue 3 SFCs, Element Plus, ESLint, Vue CLI

---

### Task 1: Fix the safe low-risk files

**Files:**
- Modify: `web/src/api/task.js`
- Modify: `web/src/main.js`
- Modify: `web/src/components/Menu.vue`
- Modify: `web/src/components/Tags.vue`
- Modify: `web/src/views/Home.vue`
- Modify: `web/src/views/Login.vue`
- Modify: `web/src/views/app/application.vue`
- Modify: `web/src/views/cmdb/Host/CreateCloud.vue`
- Modify: `web/src/views/K8s/pods/k8s-pod.vue`

- [ ] **Step 1: Run the failing lint check for the safe files**
  Run: `npx eslint src/api/task.js src/main.js src/components/Menu.vue src/components/Tags.vue src/views/Home.vue src/views/Login.vue src/views/app/application.vue src/views/cmdb/Host/CreateCloud.vue`
- [ ] **Step 2: Remove unused locals/imports in the JS files**
- [ ] **Step 3: Add multi-word component names for the Vue option components**
- [ ] **Step 4: Replace the unsupported `Copy` icon usage in `k8s-pod.vue` with `DocumentCopy`**
- [ ] **Step 5: Re-run the same `eslint` command and verify it passes**

### Task 2: Verify the tenth cleanup batch

**Files:**
- Modify: `web/src/api/task.js`
- Modify: `web/src/main.js`
- Modify: `web/src/components/Menu.vue`
- Modify: `web/src/components/Tags.vue`
- Modify: `web/src/views/Home.vue`
- Modify: `web/src/views/Login.vue`
- Modify: `web/src/views/app/application.vue`
- Modify: `web/src/views/cmdb/Host/CreateCloud.vue`
- Modify: `web/src/views/K8s/pods/k8s-pod.vue`

- [ ] **Step 1: Run `npm run build`**
- [ ] **Step 2: Run full `npm run lint` and record the new remaining error count**

## Notes

- This batch intentionally does not touch the remaining `k8s-pod.vue` lint debt beyond the icon warning fix.
- This workspace is not an initialized git worktree, so the plan is saved locally only.
