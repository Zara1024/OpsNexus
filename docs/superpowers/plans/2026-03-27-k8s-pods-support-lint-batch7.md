# K8s Pods Support Lint Batch 7 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove the next batch of pods-related lint debt from support components without changing active workflows.

**Architecture:** Use targeted `eslint` as the red-green loop. This batch focuses on dead code in terminal/monitor helper components and on stale methods in `k8s-operation-pod.vue` left behind after actions moved to emitted handlers.

**Tech Stack:** Vue 3 SFCs, Element Plus, ESLint, Vue CLI

---

### Task 1: Clean dead code from pods support components

**Files:**
- Modify: `web/src/views/K8s/pods/K8S-sterminal.vue`
- Modify: `web/src/views/K8s/pods/k8s-pod-monitor.vue`
- Modify: `web/src/views/K8s/pods/k8s-operation-pod.vue`

- [ ] **Step 1: Run the failing lint check for the three files**
  Run: `npx eslint src/views/K8s/pods/K8S-sterminal.vue src/views/K8s/pods/k8s-pod-monitor.vue src/views/K8s/pods/k8s-operation-pod.vue`
- [ ] **Step 2: Remove the unused `startHeartbeat` helper from `K8S-sterminal.vue`**
- [ ] **Step 3: Remove the unused `reactive` import from `k8s-pod-monitor.vue`**
- [ ] **Step 4: Remove the unused imports and stale helpers from `k8s-operation-pod.vue`**
- [ ] **Step 5: Remove the dead `defineExpose(viewPodList)` path and the unreachable statement in `editWorkload`**
- [ ] **Step 6: Re-run the same `eslint` command and verify it passes**

### Task 2: Verify the seventh cleanup batch

**Files:**
- Modify: `web/src/views/K8s/pods/K8S-sterminal.vue`
- Modify: `web/src/views/K8s/pods/k8s-pod-monitor.vue`
- Modify: `web/src/views/K8s/pods/k8s-operation-pod.vue`

- [ ] **Step 1: Run `npm run build`**
- [ ] **Step 2: Run full `npm run lint` and record the new remaining error count**

## Notes

- This batch intentionally does not touch `k8s-pod.vue`, which still has mixed unused-symbol and real undefined-symbol issues.
- This workspace is not an initialized git worktree, so the plan is saved locally only.
