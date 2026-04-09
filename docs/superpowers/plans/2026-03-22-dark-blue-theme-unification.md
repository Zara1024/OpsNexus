# Dark Blue Theme Unification Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Bring legacy OpsNexus pages and dialogs onto the current deep blue dark visual system without changing feature behavior.

**Architecture:** Use `web/src/assets/css/global.css` as the primary compatibility layer for legacy page roots, cards, tables, dialogs, and forms. Apply minimal page-level overrides only where scoped component styles remain more specific than the global layer.

**Tech Stack:** Vue 3, Element Plus, scoped SFC CSS, global theme tokens in `web/src/assets/css/global.css`

---

### Task 1: Add Global Legacy Theme Compatibility Layer

**Files:**
- Modify: `web/src/assets/css/global.css`

- [ ] Add a dedicated section for legacy dark-theme reconciliation.
- [ ] Target legacy page roots and shared card classes with deep blue dark surfaces.
- [ ] Target legacy tables, dialogs, inputs, tabs, and labels with token-based colors.
- [ ] Keep selectors narrow to known legacy roots.
- [ ] Rebuild frontend and confirm CSS compiles.

### Task 2: Fix Namespace Page and Table Drift

**Files:**
- Modify: `web/src/views/K8s/k8s-namespace.vue`
- Modify: `web/src/views/K8s/namespaces/NamespacesTable.vue`

- [ ] Remove or neutralize any remaining page-local purple/light-card rules that global CSS cannot safely override.
- [ ] Align namespace table header, row text, quota badges, and action buttons to the dark system.
- [ ] Align namespace dialogs and detail sections to the same surface hierarchy.
- [ ] Rebuild frontend and verify `/k8s/namespace` in browser.

### Task 3: Cover Remaining High-Impact Legacy Modules

**Files:**
- Modify as needed:
  - `web/src/views/K8s/k8s-config.vue`
  - `web/src/views/K8s/k8s-network.vue`
  - `web/src/views/K8s/k8s-storage.vue`
  - `web/src/views/K8s/k8s-nodes.vue`
  - `web/src/views/K8s/k8s-workloads.vue`
  - `web/src/views/cmdb/cmdbDB.vue`
  - `web/src/views/cmdb/cmdbHost.vue`
  - `web/src/views/configcenter/KeyManage.vue`
  - `web/src/views/configcenter/ecs-key.vue`
  - `web/src/views/Tools/DeployManage.vue`
  - `web/src/views/system/Personal.vue`

- [ ] Let the global layer absorb as much as possible.
- [ ] Only patch page-local CSS when a specific element still resists the global layer.
- [ ] Avoid script/template changes unless purely cosmetic and necessary.
- [ ] Rebuild frontend after the batch.

### Task 4: Verify and Sync

**Files:**
- Modify: none expected

- [ ] Run `cd web && npm run build`
- [ ] Inspect representative remote pages in browser:
  - `/k8s/namespace`
  - `/k8s/config`
  - `/k8s/network`
  - `/cmdb/db`
  - `/config/keymanage` or equivalent route
- [ ] Sync rebuilt frontend assets to `10.0.0.200`
- [ ] Re-verify bundle hashes and representative page rendering

### Notes

- This repo does not currently expose a dedicated frontend visual regression test harness, so verification will rely on fresh production builds plus browser-based route inspection and computed-style evidence.
