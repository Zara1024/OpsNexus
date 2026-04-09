# K8s Small Unused Lint Batch 5 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove a small set of remaining K8s unused-symbol lint errors from namespace and pod support components.

**Architecture:** Use targeted `eslint` checks as the verification loop. This batch is intentionally limited to declarations already proven unused by lint, with no behavior changes.

**Tech Stack:** Vue 3 SFCs, Element Plus, ESLint, Vue CLI

---

### Task 1: Clean simple unused declarations in namespace and pod support files

**Files:**
- Modify: `web/src/views/K8s/namespaces/NamespacesTable.vue`
- Modify: `web/src/views/K8s/pods/CreatePodDialog.vue`
- Modify: `web/src/views/K8s/pods/DialogManager.vue`
- Modify: `web/src/views/K8s/pods/PodListDialog.vue`
- Modify: `web/src/views/K8s/pods/k8s-container-pods.vue`

- [ ] **Step 1: Run the failing lint check for the five files**
  Run: `npx eslint src/views/K8s/namespaces/NamespacesTable.vue src/views/K8s/pods/CreatePodDialog.vue src/views/K8s/pods/DialogManager.vue src/views/K8s/pods/PodListDialog.vue src/views/K8s/pods/k8s-container-pods.vue`
- [ ] **Step 2: Remove the unused `props` assignment in `NamespacesTable.vue`**
- [ ] **Step 3: Remove the unused `reactive` import in `CreatePodDialog.vue`**
- [ ] **Step 4: Remove the unused `K8S_EVENTS` import in `DialogManager.vue`**
- [ ] **Step 5: Remove the unused `props` assignment in `PodListDialog.vue`**
- [ ] **Step 6: Remove the unused `computed` import and unused `emit` declaration in `k8s-container-pods.vue`**
- [ ] **Step 7: Re-run the same `eslint` command and verify it passes**

### Task 2: Verify the fifth cleanup batch

**Files:**
- Modify: `web/src/views/K8s/namespaces/NamespacesTable.vue`
- Modify: `web/src/views/K8s/pods/CreatePodDialog.vue`
- Modify: `web/src/views/K8s/pods/DialogManager.vue`
- Modify: `web/src/views/K8s/pods/PodListDialog.vue`
- Modify: `web/src/views/K8s/pods/k8s-container-pods.vue`

- [ ] **Step 1: Run `npm run build`**
- [ ] **Step 2: Run full `npm run lint` and record the new remaining error count**

## Notes

- This batch only removes unused declarations already flagged by lint.
- This workspace is not an initialized git worktree, so the plan is saved locally only.
