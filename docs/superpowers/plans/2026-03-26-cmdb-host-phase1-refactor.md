# CMDB Host Phase 1 Refactor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver the Phase 1 CMDB host page refactor from the user-provided prompt by adding a connection center, batch host actions, command risk pre-check, and host-to-terminal-audit linkage without rewriting the existing CMDB flow.

**Architecture:** Keep the current Vue Options API host page and existing Go handlers, but extract new host-phase logic into small reusable helpers so the page stays maintainable. Reuse current SSH, sync, agent, and terminal audit capabilities; add only a lightweight connectivity endpoint plus query/filter wiring where existing APIs are insufficient.

**Tech Stack:** Vue 3 + Element Plus, axios request wrapper, Go + Gin + Gorm, Node built-in test runner, Go test.

---

### Task 1: Lock the reusable host-phase helper contract with tests

**Files:**
- Create: `web/src/utils/cmdbHostPhase1.mjs`
- Modify: `web/src/utils/cmdbHostPhase1.test.mjs`

- [ ] **Step 1: Run the existing host-phase utility test and confirm it fails**

Run: `node --test web/src/utils/cmdbHostPhase1.test.mjs`

Expected: FAIL because `web/src/utils/cmdbHostPhase1.mjs` does not exist yet.

- [ ] **Step 2: Extend the failing test coverage for batch toolbar and command risk helper behavior**

Add cases for:
- connection entries disabled reasons on Windows
- audit deep-link query generation
- batch action summary/count helpers
- command risk copy/confirm-state helper normalization

- [ ] **Step 3: Implement the minimal helper module**

Implement pure helpers for:
- building host connection-center entries
- building terminal-audit route query params from a host
- summarizing batch connectivity results
- normalizing command risk response state for the command dialog

- [ ] **Step 4: Re-run the host-phase utility test**

Run: `node --test web/src/utils/cmdbHostPhase1.test.mjs`

Expected: PASS.

### Task 2: Add lightweight backend support for batch connectivity and stable command-risk responses

**Files:**
- Modify: `api/api/cmdb/model/cmdbHost.go`
- Modify: `api/api/cmdb/controller/cmdbHost.go`
- Modify: `api/api/cmdb/service/cmdbHost.go`
- Modify: `api/router/cmdb/cmdb.go`
- Modify: `api/api/cmdb/controller/cmdbHostSSH.go`
- Add: `api/api/cmdb/service/cmdbHostConnectivity_test.go`

- [ ] **Step 1: Write the failing Go test for connectivity target/rule selection**

Cover:
- Linux hosts prefer SSH IP/port
- Windows hosts prefer RDP address/port
- missing connection fields produce explicit failure reasons

- [ ] **Step 2: Run the focused Go test and confirm it fails**

Run: `go test ./api/cmdb/service -run TestHostConnectivity -v`

Expected: FAIL because the helper/types do not exist yet.

- [ ] **Step 3: Implement minimal connectivity DTOs, service logic, and route**

Add a lightweight batch endpoint that:
- accepts `hostIds`
- resolves each host's connectivity target from existing fields
- performs a timeout-bounded TCP dial
- returns per-host success/failure plus reason text

- [ ] **Step 4: Normalize command-risk conflict responses**

Keep the current `cmdb/hostssh/command/:id` contract but ensure conflict messages and returned risk payload are stable, readable, and reusable by the frontend confirmation flow.

- [ ] **Step 5: Re-run the focused Go test**

Run: `go test ./api/cmdb/service -run TestHostConnectivity -v`

Expected: PASS.

### Task 3: Wire frontend APIs and terminal-audit page filters for real host deep links

**Files:**
- Modify: `web/src/api/cmdb.js`
- Modify: `web/src/views/monitor/Recording.vue`
- Modify: `web/src/api/system.js`

- [ ] **Step 1: Write/extend failing utility tests for audit-route and connectivity summaries**

Reuse `web/src/utils/cmdbHostPhase1.test.mjs` to cover the route payload expected by the Recording page.

- [ ] **Step 2: Confirm the test fails for the new route/filter behavior if needed**

Run: `node --test web/src/utils/cmdbHostPhase1.test.mjs`

Expected: FAIL if new helper expectations were added.

- [ ] **Step 3: Implement frontend API wrappers**

Add wrappers for:
- batch connectivity test
- command execution with `riskAck` and `confirmedRiskLevel`

- [ ] **Step 4: Upgrade terminal audit page filter handling**

Add:
- `hostKeyword` filter input
- hidden/route-backed `hostId`
- route query hydration on create/route-change
- deep-link friendly risk labels for low/medium/high

- [ ] **Step 5: Re-run the utility test**

Run: `node --test web/src/utils/cmdbHostPhase1.test.mjs`

Expected: PASS.

### Task 4: Refactor the CMDB host page around the new phase-1 workflow

**Files:**
- Modify: `web/src/views/cmdb/cmdbHost.vue`
- Modify: `web/src/views/cmdb/Host/CmdbHostTable.vue`
- Modify: `web/src/views/cmdb/Host/CmdbHostTableCompact.vue`

- [ ] **Step 1: Add failing expectations where possible via helper tests before UI wiring**

Use helper tests to define:
- connection-center entry order/state
- audit deep-link query output
- connectivity summary text handling

- [ ] **Step 2: Add table multi-select and batch toolbar**

Implement:
- selection column
- selected-host state/events
- batch sync
- batch connectivity test
- batch deploy/uninstall
- optional batch move-group only if low-risk and fits existing APIs cleanly

- [ ] **Step 3: Add host detail connection center**

In the detail drawer show:
- Linux: Web SSH / execute command / upload / sync
- Windows: RDP info plus disabled Linux-only actions with reasons
- explicit disabled-state copy rather than hiding actions

- [ ] **Step 4: Add command risk pre-check UX**

Upgrade the command dialog to:
- preview risk before execution
- show low/medium/high cues
- require explicit second confirmation for medium/high
- pass `riskAck` and `confirmedRiskLevel` only after confirmation

- [ ] **Step 5: Add terminal-audit linkage in host detail**

Show:
- recent audit summary rows
- recent risky session entry point
- buttons that route to `/monitor/recording` with real host filters

- [ ] **Step 6: Smoke the page in lintable form**

Run: `cd web && npm run lint`

Expected: PASS, or only pre-existing unrelated failures.

### Task 5: Full verification and delivery evidence

**Files:**
- No code changes required unless verification exposes regressions

- [ ] **Step 1: Run frontend verification**

Run: `cd web && npm run lint`

Expected: PASS.

- [ ] **Step 2: Run targeted Node utility verification**

Run: `node --test web/src/utils/cmdbHostPhase1.test.mjs`

Expected: PASS.

- [ ] **Step 3: Run backend verification**

Run: `cd api && go test ./...`

Expected: PASS, or document any pre-existing unrelated failures with exact evidence.

- [ ] **Step 4: Prepare manual regression checklist**

Cover:
- `/cmdb/ecs` load
- multi-select and batch toolbar
- one successful and one partial-failure batch action
- Linux connection center actions
- Windows RDP info + disabled reasons
- low/medium/high command execution flow
- host detail to terminal-audit deep link
- no regression in edit/upload/sync/process/TCP monitor

- [ ] **Step 5: Sync the completed code to `10.0.0.200`**

Use the provided root credential and copy the changed project state after verification.
