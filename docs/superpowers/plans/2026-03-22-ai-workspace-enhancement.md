# AI Workspace Enhancement Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Upgrade OpsNexus AI into a more complete operations workspace by adding an overview API, richer assistant context, and expanded diagnosis scenes while preserving existing behavior.

**Architecture:** Extend the current AI backend with one summary endpoint and additional diagnosis scene handlers, then update the existing Vue assistant and diagnosis pages to consume the summary data and expose clearer operational entry points. Keep all changes additive and backwards compatible.

**Tech Stack:** Go, Gin, GORM, Vue 3, Element Plus

---

### Task 1: Add AI Overview Contract

**Files:**
- Create: `api/api/ai/model/overview.go`
- Modify: `api/api/ai/controller/ai.go`
- Modify: `api/router/ai/ai.go`
- Modify: `api/api/ai/service/ai.go`
- Test: `api/api/ai/service/overview_test.go`

- [ ] Write a failing unit test for the overview summary helper.
- [ ] Run the focused Go test and confirm it fails for the missing helper/behavior.
- [ ] Add overview response structs and summary builder/helper logic.
- [ ] Add `ListOverview` service/controller/route wiring.
- [ ] Run the focused Go test and confirm it passes.

### Task 2: Expand Diagnosis Scenes

**Files:**
- Modify: `api/api/ai/service/ai.go`
- Test: `api/api/ai/service/overview_test.go` or a new focused AI service test file if needed

- [ ] Write failing tests for scene metadata/scene normalization helpers if extracted.
- [ ] Add additive diagnosis support for `alert_analysis` and `deployment_review`.
- [ ] Reuse existing monitor/app data models instead of introducing duplicate queries.
- [ ] Keep existing diagnosis scenes unchanged.
- [ ] Run the relevant Go tests and confirm they pass.

### Task 3: Wire Frontend API

**Files:**
- Modify: `web/src/api/ai.js`

- [ ] Add frontend API method for the new overview endpoint.
- [ ] Keep existing API calls unchanged.

### Task 4: Enhance Assistant Workspace

**Files:**
- Modify: `web/src/views/ai/AIAssistant.vue`

- [ ] Add overview-driven runtime/status content.
- [ ] Add domain entry cards and better quick-start prompts.
- [ ] Reuse existing message flow and history behavior.
- [ ] Keep route `prompt` autostart working.
- [ ] Run frontend build to confirm no regressions.

### Task 5: Enhance Diagnosis Workspace

**Files:**
- Modify: `web/src/views/ai/AIDiagnosis.vue`

- [ ] Add overview-driven status/context content.
- [ ] Add the new diagnosis scenes into UI scene configuration.
- [ ] Keep existing route-query autoload behavior intact.
- [ ] Run frontend build to confirm no regressions.

### Task 6: Verify End to End

**Files:**
- Modify: none expected

- [ ] Run `cd api && go test ./...`
- [ ] Run `cd web && npm run build`
- [ ] Spot-check `/ai/assistant` and `/ai/diagnosis`
- [ ] If verified locally, sync frontend/backend changes as needed for remote validation

### Notes

- The frontend currently has no dedicated automated UI test harness, so verification will rely on focused Go tests plus production build validation and targeted browser inspection.
