# CMDB Device Platform Single-Select Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace manual device platform input with a fixed single-select choice in the CMDB network-device dialogs while preserving the existing `platform` string payload.

**Architecture:** A shared utility exports the canonical platform option list and appends a temporary legacy option when editing older records. The create and edit dialogs consume that shared helper through `el-radio-group`, while the rest of the device data flow stays unchanged.

**Tech Stack:** Vue 3, Element Plus, local Node-based `.mjs` tests

---

### Task 1: Add a failing utility test for device platform options

**Files:**
- Create: `web/tests/cmdb-device-platform-options.test.mjs`
- Modify: `web/src/utils/cmdbAssetPresentation.mjs`

- [ ] **Step 1: Write the failing test**
- [ ] **Step 2: Run `node web/tests/cmdb-device-platform-options.test.mjs` and verify it fails because the helper does not exist yet**
- [ ] **Step 3: Implement the minimal shared helper in `web/src/utils/cmdbAssetPresentation.mjs`**
- [ ] **Step 4: Run `node web/tests/cmdb-device-platform-options.test.mjs` and verify it passes**

### Task 2: Switch create and edit dialogs to single-select platform radios

**Files:**
- Modify: `web/src/views/cmdb/Device/CreateDevice.vue`
- Modify: `web/src/views/cmdb/Device/EditDevice.vue`

- [ ] **Step 1: Replace the free-text `platform` input with a shared radio-group rendering**
- [ ] **Step 2: Add required validation for `platform`**
- [ ] **Step 3: Make edit mode append a temporary legacy option when needed**
- [ ] **Step 4: Keep submit payloads unchanged except for using the selected radio value**

### Task 3: Verify the device form behavior and regression coverage

**Files:**
- Modify: `web/tests/cmdb-device-presentation.test.mjs`
- Test: `web/tests/cmdb-device-platform-options.test.mjs`

- [ ] **Step 1: Run `node web/tests/cmdb-device-platform-options.test.mjs`**
- [ ] **Step 2: Run `node web/tests/cmdb-device-presentation.test.mjs`**
- [ ] **Step 3: Run targeted lint on the changed frontend files**
- [ ] **Step 4: Summarize any remaining risk, especially around legacy platform values**
