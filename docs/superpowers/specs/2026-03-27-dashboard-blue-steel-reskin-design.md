# Dashboard Blue Steel Reskin Design

**Date:** 2026-03-27

**Goal:** Upgrade the `/dashboard` page from a generic dark admin surface to a more polished OpsNexus command-center look using a blue-steel, tech-outline visual language, without changing layout structure, routes, data behavior, or user workflows.

## Current Problem

- The current dashboard is functional, but the overall look feels flat and inconsistent with the intended OpsNexus "platform workbench / operations command center" positioning.
- Several sections still read like standard dark admin cards instead of a unified instrument panel.
- The page has strong modules already, but visual hierarchy between header, KPIs, center cards, charts, and quick tools is not distinct enough.
- The user explicitly prefers a full dashboard-wide improvement, but only as a reskin, not a structural redesign.

## Chosen Approach

- Use the approved `A1` direction: pure reskin, no structural rearrangement.
- Use the approved visual direction: `A` / "blue steel tech outline", biased slightly toward a command-center tone but kept restrained enough for daily use.
- Keep the existing dashboard structure exactly as-is:
  - page header
  - KPI stat strip
  - risk / pending section
  - AI workbench section
  - chart row
  - quick tools area
- Apply the stronger visual language mainly at the dashboard page level, while only lightly enhancing shared platform components so other pages do not accidentally inherit an overly dashboard-specific style.

## Why This Approach

- It directly addresses the user's complaint that the dashboard "doesn't look good" without creating feature risk.
- It gives a clearly visible quality jump faster than a deeper redesign.
- It avoids turning the page into a full-screen exhibition dashboard, which would hurt usability for a real admin/workbench product.
- It fits the available PSD material best: the strongest reusable references are the header, tab, chart-frame, numeric-card, and control-panel pieces rather than full-page templates.

## Visual Baseline

- Primary mood: deep blue shell + blue-cyan accent + crisp outline hierarchy.
- Surface style: cut-corner cards, thin luminous borders, inner shadows, restrained glows.
- Background style: subtle tech texture / gradient atmosphere, not a heavy full-page map or noisy sci-fi wallpaper.
- Accent behavior:
  - `primary`: blue-cyan
  - `success`: cool teal
  - `warning`: amber-gold
  - `danger`: alert red, but still compatible with dark-blue surroundings
- Motion: minimal; only gentle hover lift, edge glow, and state emphasis.
- Readability rule: decoration must never compete with titles, metrics, actions, or chart data.

## Scope

- Dashboard-only visual refresh for:
  - page shell and background treatment
  - header panel
  - KPI stat strip
  - central module cards
  - chart containers and chart palette tuning
  - quick tools panel
  - dashboard-local dialog styling where visual mismatch is obvious
- Light shared-component enhancement for:
  - `web/src/components/platform/PageHeader.vue`
  - `web/src/components/platform/SectionCard.vue`
  - `web/src/components/platform/StatStrip.vue`

## Non-Goals

- No route changes.
- No API changes.
- No data model changes.
- No chart type changes.
- No module reorder or layout redesign.
- No new dashboard features.
- No global redesign of unrelated pages.

## Design Details

### 1. Header Panel

- Keep the current title, subtitle, intro, chips, and action buttons in their existing positions.
- Turn the header into a command-panel style block with:
  - stronger outline treatment
  - a subtle top-edge highlight
  - tech-divider details around eyebrow/title
  - more deliberate button hierarchy
- Reference source: the PSD header assets from the left-aligned title/header set, plus related search/tag framing styles from the visualization bundle.

### 2. KPI Stat Strip

- Keep the same four KPI cards and the same values.
- Reskin cards into compact instrument tiles:
  - stronger tone separation by status
  - better number emphasis
  - tighter label/value/hint hierarchy
  - hover glow and edge response
- Make these cards the first obvious "quality jump" when entering the dashboard.

### 3. Risk / Pending and AI Workbench Cards

- Keep the same left-right structure and same content blocks.
- Upgrade `SectionCard` to feel like a module frame rather than a generic content card.
- Strengthen the title bar, outline, and internal spacing rhythm.
- Make risk cards feel more state-driven, especially danger/warning items.
- Make AI capabilities, chips, and recommended actions share the same panel grammar so the right side no longer feels visually fragmented.

### 4. Chart Containers

- Keep the same trend chart, pie/ring chart, and heat/other chart usage.
- Improve the enclosing frames:
  - chart title/header framing
  - time-range switch styling
  - outline + inset layer treatment
  - cleaner dark grid/background harmony
- Tune ECharts colors so chart strokes, fills, labels, and gridlines match the new blue-steel palette.
- The goal is "chart bay / instrument frame", not "ordinary chart card".

### 5. Quick Tools and Dashboard Dialog

- Keep the quick tool CRUD behavior unchanged.
- Reskin the tools panel into a more deliberate control surface:
  - icon pedestal feel
  - better hover/focus edges
  - unified edit/delete action buttons
- Align the tool edit dialog with the dashboard theme so opening it does not revert the page back to a generic admin style.

## File Strategy

- Primary dashboard page styling and chart tuning:
  - `web/src/views/dashboard/index.vue`
- Light shared component upgrades:
  - `web/src/components/platform/PageHeader.vue`
  - `web/src/components/platform/SectionCard.vue`
  - `web/src/components/platform/StatStrip.vue`
- Optional token adjustment only if strictly necessary:
  - `web/src/assets/css/global.css`

## Risk Control

- Preserve all template structure unless a wrapper/class is truly needed for styling.
- Prefer CSS and small template-class additions over script changes.
- Keep shared-component changes restrained and generic.
- Contain strong dashboard visuals within dashboard-local classes to avoid polluting other platform pages.
- Validate readability, button clarity, and hover feedback after styling changes.

## Acceptance Criteria

1. The dashboard reads as one coherent blue-steel command-center surface rather than a collection of generic dark cards.
2. Header, KPI strip, center cards, charts, and quick tools visibly belong to the same visual system.
3. Layout structure, navigation, data loading, and tool CRUD behavior remain unchanged.
4. Charts still render correctly and gain clearer container hierarchy.
5. The page looks noticeably improved at first glance without becoming noisy or hard to use.
6. The design remains suitable for daily operations use, not just demo display.

## Verification

- `cd web && npm run build`
- Open `/dashboard` locally or in the target environment and verify:
  - header panel renders correctly
  - KPI cards render with correct tone differentiation
  - risk cards and AI workbench maintain click behavior
  - chart tabs and chart rendering still work
  - quick tool add/edit/delete dialog still opens and styles correctly
- Manual visual check on desktop width to confirm no clipping, overlap, or unusable glow/contrast issues.

## Constraint Note

- The superpowers brainstorming flow recommends a reviewer subagent after spec writing, but subagent delegation is not user-authorized in this session, so spec review will be handled inline.
- The current workspace root is not an initialized git worktree, so the design doc can be written locally but cannot be committed from this location unless the actual git root is provided later.
