# Tools Lint Batch 9 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove the next set of low-risk lint issues from the Tools pages without changing runtime behavior.

**Architecture:** Use targeted `eslint` checks as the red-green loop. This batch only fixes component names, unused icon imports, an unused `props` binding, and one unused slot parameter.

**Tech Stack:** Vue 3 SFCs, Element Plus, ESLint, Vue CLI

---

### Task 1: Clean low-risk Tools page lint issues

**Files:**
- Modify: `web/src/views/Tools/Agent.vue`
- Modify: `web/src/views/Tools/DeployManage.vue`
- Modify: `web/src/views/Tools/SelectDeployHost.vue`
- Modify: `web/src/views/Tools/Tools.vue`

- [ ] **Step 1: Run the failing lint check for the four files**
  Run: `npx eslint src/views/Tools/Agent.vue src/views/Tools/DeployManage.vue src/views/Tools/SelectDeployHost.vue src/views/Tools/Tools.vue`
- [ ] **Step 2: Add multi-word component names for `Agent.vue` and `Tools.vue`**
- [ ] **Step 3: Remove the unused icon imports from `Agent.vue`**
- [ ] **Step 4: Remove the unused `props` binding and dead icon imports from `DeployManage.vue`**
- [ ] **Step 5: Remove the unused `data` slot binding from `SelectDeployHost.vue`**
- [ ] **Step 6: Re-run the same `eslint` command and verify it passes**

### Task 2: Verify the ninth cleanup batch

**Files:**
- Modify: `web/src/views/Tools/Agent.vue`
- Modify: `web/src/views/Tools/DeployManage.vue`
- Modify: `web/src/views/Tools/SelectDeployHost.vue`
- Modify: `web/src/views/Tools/Tools.vue`

- [ ] **Step 1: Run `npm run build`**
- [ ] **Step 2: Run full `npm run lint` and record the new remaining error count**

## Notes

- This batch intentionally avoids `k8s-pod.vue` because that file now needs a separate, more careful repair pass.
- This workspace is not an initialized git worktree, so the plan is saved locally only.
