# Deep Blue Dark Theme Unification Design

**Date:** 2026-03-22

**Goal:** Unify legacy pages and dialogs that still use the old purple gradient + light card visual system into the current deep blue dark OpsNexus design language without changing feature behavior.

## Scope

- K8s list/detail pages with obvious style drift, especially namespace, config, network, storage, nodes, workloads, cluster-related dialogs, and related table components.
- CMDB pages that still render with light cards or low-contrast text.
- Config center, tools, quick release, and personal center pages that still use the old legacy visual template.
- Dialogs, drawers, and detail cards opened from those pages when they visibly break the unified dark style.

## Non-Goals

- No API changes.
- No data model changes.
- No route changes.
- No interaction or permission logic changes.
- No visual redesign of already-aligned pages.

## Design Baseline

- Use the existing deep blue dark theme tokens from `web/src/assets/css/global.css`.
- Prefer dark surfaces built from `--bg-shell`, `--bg-surface`, `--bg-elevated`, `--border-*`, and `--text-*`.
- Remove or neutralize legacy purple gradients `#667eea/#764ba2`, light cards `rgba(255,255,255,0.95)`, and low-contrast text colors like `#2c3e50`, `#606266`, `#909399` when used on dark surfaces.
- Keep blue as the primary accent; avoid introducing a second page-specific color system.

## Implementation Strategy

1. Add global compatibility overrides for legacy page roots, cards, tables, forms, and dialogs.
2. Add targeted overrides for page-specific elements that remain visually inconsistent after global coverage.
3. Keep overrides selector-scoped to known legacy page roots and dialog classes to avoid affecting already-modern modules.
4. Verify using build output and real browser inspection on representative pages.

## Risk Control

- Use CSS-only changes where possible.
- Prefer global overrides with scoped selectors over template or script changes.
- Limit component file edits to cases where scoped styles are too specific to override cleanly from global CSS.
- Validate via `npm run build` and browser-based inspection on remote environment pages before sync completion.
