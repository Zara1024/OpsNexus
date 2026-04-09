# AI Runtime Real Status Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make AI runtime status reflect real probe results instead of config-only readiness, and switch the remote runtime base URL to `https://beecode.cc/v1`.

**Architecture:** Add a cached real probe in the OpenAI-compatible runtime client, feed probe outcome into AI overview, then render the diagnosis and assistant badges from that probe-derived state. Keep fallback execution behavior unchanged. Update the remote runtime base URL and verify with fresh real requests.

**Tech Stack:** Go, Vue 3, OpenAI-compatible HTTP client, remote systemd deployment.

---

### Task 1: Add a failing backend regression test

**Files:**
- Create: `api/api/ai/service/runtime_openai_test.go`

- [ ] **Step 1: Write a test for probe-derived runtime status**
- [ ] **Step 2: Run only that test and confirm it fails**

### Task 2: Implement backend real runtime status

**Files:**
- Modify: `api/api/ai/service/runtime_openai.go`
- Modify: `api/api/ai/service/overview.go`
- Modify: `api/api/ai/model/overview.go`

- [ ] **Step 1: Add cached probe support to the runtime client**
- [ ] **Step 2: Return `ready` / `degraded` / `fallback` from overview based on probe result**
- [ ] **Step 3: Re-run backend test and confirm it passes**

### Task 3: Update frontend status rendering

**Files:**
- Modify: `web/src/views/ai/AIDiagnosis.vue`
- Modify: `web/src/views/ai/AIAssistant.vue`

- [ ] **Step 1: Render degraded state distinctly from ready**
- [ ] **Step 2: Show runtime detail/error text when probe failed**
- [ ] **Step 3: Build frontend and confirm no build regressions**

### Task 4: Update remote runtime endpoint and verify

**Files:**
- Modify remote: `/opt/opsnexus-remote-test/config.yaml`

- [ ] **Step 1: Replace `https://uniquefox.top/v1` with `https://beecode.cc/v1`**
- [ ] **Step 2: Restart `opsnexus-api.service`**
- [ ] **Step 3: Verify `/api/v1/ai/overview` plus one assistant request and one diagnosis request**
