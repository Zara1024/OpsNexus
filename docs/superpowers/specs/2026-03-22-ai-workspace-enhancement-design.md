# AI Workspace Enhancement Design

**Date:** 2026-03-22

**Goal:** Enhance the existing OpsNexus AI area into a clearer, more operational workspace by adding an AI overview layer, richer scene organization, and stronger assistant/diagnosis guidance without breaking current assistant, diagnosis, history, inspection, or knowledge flows.

## Context

OpsNexus already has:

- AI assistant chat
- AI diagnosis by scene
- prompt templates
- knowledge suggestions
- inspection templates and reports
- OpenAI Responses runtime integration

The reference project KubePolaris does not provide a direct AI module to port, but it does provide a strong model for information architecture:

- clear task-oriented entry points
- detail-first operational layouts
- stronger context framing
- better drilling from overview to detail

## Scope

- Add one lightweight AI overview API for workspace summary and runtime status.
- Enrich `/ai/assistant` with operational overview, scenario/domain entry cards, and better runtime/status visibility.
- Enrich `/ai/diagnosis` with overview data and new high-value diagnosis scenes.
- Add at least two practical diagnosis scenes that reuse existing platform data:
  - alert analysis
  - deployment review

## Non-Goals

- No breaking changes to existing AI routes.
- No replacement of current assistant orchestration.
- No model-provider management console in this phase.
- No destructive AI execution changes.

## Design Direction

### 1. AI Overview Layer

Add `GET /api/v1/ai/overview` to return:

- runtime/provider/model/enabled/fallback posture
- counts for prompt templates, knowledge items, diagnosis sessions, assistant sessions, inspection templates, reports, pending confirmations
- supported diagnosis scenes
- suggested entry prompts / domains

This gives frontend a stable summary contract and avoids hardcoding too much product state in the UI.

### 2. Assistant Workspace Enhancement

Keep `/ai/assistant` as the main conversational surface, but add:

- runtime status card
- domain cards for host, k8s, alert, work order, deployment, inspection
- clearer "what can I ask" guidance
- richer workspace stats backed by overview API

The assistant remains chat-first, but becomes easier to start using.

### 3. Diagnosis Workspace Enhancement

Keep `/ai/diagnosis` as a structured analysis surface, but expand scenes so it covers more real platform workflows.

New scenes:

- `alert_analysis`
  Uses monitor alert summary + recent incidents + knowledge context to generate an AI report.

- `deployment_review`
  Uses quick deployment basic data + task details + knowledge context to generate release review and rollback-oriented guidance.

These scenes complement existing terminal audit, SQL work order, workload capacity, inspection report, and knowledge search analysis.

## Risk Control

- Reuse existing DAOs/models whenever possible.
- Keep assistant and diagnosis response shapes backward compatible.
- Add new scene support without changing existing scene behavior.
- Validate via backend tests where practical and fresh frontend production builds.
