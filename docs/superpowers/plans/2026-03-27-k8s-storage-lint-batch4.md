# K8s Storage Lint Batch 4 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Clean the `k8s-storage` page lint debt, including the stale namespace-loading code that now references an undefined local state.

**Architecture:** Use targeted `eslint` as the red-green loop. The cleanup removes dead imports, dead resource-name refs, the obsolete `fetchNamespaceList` flow left over after adopting `NamespaceSelector`, and a few unused helper functions that no longer affect rendering.

**Tech Stack:** Vue 3 SFCs, Element Plus, ESLint, Vue CLI

---

### Task 1: Clean dead code from `k8s-storage.vue`

**Files:**
- Modify: `web/src/views/K8s/k8s-storage.vue`

- [ ] **Step 1: Run the failing lint check for `k8s-storage.vue`**
  Run: `npx eslint src/views/K8s/k8s-storage.vue`
- [ ] **Step 2: Remove the unused icon imports**
- [ ] **Step 3: Remove the unused `currentResourceName` and `currentResourceType` refs**
- [ ] **Step 4: Remove the obsolete `fetchNamespaceList` function that references `namespaceList`**
- [ ] **Step 5: Remove the unused `resetSearch`, `getClusterStatusText`, and `getClusterStatusTag` helpers**
- [ ] **Step 6: Re-run `npx eslint src/views/K8s/k8s-storage.vue` and verify it passes**

### Task 2: Verify the storage cleanup batch

**Files:**
- Modify: `web/src/views/K8s/k8s-storage.vue`

- [ ] **Step 1: Run `npm run build` to confirm the cleanup does not break the app**
- [ ] **Step 2: Re-run full `npm run lint` and record the new remaining error count**

## Notes

- The undefined `namespaceList` symbol is part of dead code; the live page now uses `NamespaceSelector` instead of managing a local namespace list.
- This workspace is not an initialized git worktree, so the plan is saved locally only.
