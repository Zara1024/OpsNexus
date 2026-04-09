# OpsNexus UI Design System Design

**Date:** 2026-04-07

**Goal:** Define a project-wide UI design baseline for OpsNexus so future UI optimization work can stay visually unified, operationally readable, and compatible with the platform's existing routes, workflows, permissions, and business logic.

## Current Problem

- The project already covers many modules, but the overall visual language is still fragmented between older admin-style pages, dark blue refinements, and page-local styling decisions.
- Some key pages are strong individually, but the platform does not yet read as one coherent product across login, dashboard, K8s, monitoring, CMDB, config center, work orders, and system management.
- The user wants to use `awesome-design-md` as a design reference input, but not as a literal template source.
- The chosen visual direction must balance two needs:
  - daily enterprise workbench usability
  - stronger OpsNexus platform identity on flagship pages

## Chosen Direction

- Use `B` as the primary direction: a mature industrial product surface similar in spirit to platforms like HashiCorp / Sentry.
- Absorb `A` selectively: inject blue-steel / blue-cyan tech accents only into high-value Ops pages such as login, dashboard, monitoring, K8s, logs, terminal, and audit surfaces.
- Avoid `C` as a system-wide baseline because it is better suited to demo surfaces or large-screen visualization than a daily operational platform.

## Why This Direction

- It gives OpsNexus a more professional and durable product identity instead of a generic admin template look.
- It avoids over-styling dense workflows such as configuration management, forms, tables, and governance pages.
- It still leaves room for a distinctive command-center tone on flagship pages where stronger atmosphere is appropriate.
- It is the lowest-risk way to unify a large existing platform without rewriting structure or behavior.

## Design Principles

1. **Product Surface First**
   - The default platform tone is stable, professional, and restrained.
   - Readability, workflow continuity, and operational efficiency take priority over decorative effects.

2. **Localized Tech Accent**
   - Blue tech treatment is a highlight tool, not the default background for every page.
   - Stronger atmosphere is reserved for login, dashboard, monitoring, K8s workbench, logs, terminal, and audit views.

3. **Unified Skeleton, Contextual Expression**
   - Shared tokens and shared component rules must be global.
   - Page classes may express different emphasis:
     - monitoring / K8s: more command-center
     - CMDB / config center: more industrial product surface
     - system pages: more governance-oriented and restrained

4. **Information Density Over Decoration**
   - OpsNexus is a production tool, not a marketing site.
   - Visual improvement should increase hierarchy clarity, state recognition, and decision speed.

## Scope

- Define a reusable design system for:
  - color and surface tokens
  - page shells and page headers
  - cards and panels
  - tables, filters, forms, dialogs, drawers
  - status expression
  - chart containers and monitoring surfaces
  - page-type guidance across major modules
- Provide prompt-ready guidance for future human or AI-driven UI optimization.

## Non-Goals

- No route changes.
- No API changes.
- No permission model changes.
- No business workflow changes.
- No forced full dark mode rollout.
- No visual cloning of external reference sites.

## Visual Baseline

### 1. Color Strategy

- Use a cool-neutral product base with restrained blue brand emphasis.
- Suggested roles:
  - `base neutral`: shell background, page background, text, dividers
  - `brand primary`: stable blue for selected actions, active state, links, focus
  - `tech accent`: brighter cyan-blue used only on flagship pages and key metrics
  - `success`: cool green
  - `warning`: amber
  - `danger`: governance red, not fluorescent red
- Avoid over-saturating large surfaces with pure blue or pure black.

### 2. Background Rules

- Default business pages use a light or lightly cool-tinted product surface.
- Key flagship pages may use deep blue or mixed dark-light shells, but within controlled boundaries.
- Recommended emphasis by page class:
  - `login`: strongest atmosphere
  - `dashboard / monitoring / K8s`: medium atmosphere
  - `lists / forms / system pages`: restrained product surface

### 3. Card System

- Standardize three panel types:
  - `data card`: KPI or metric emphasis
  - `module card`: title + action + content
  - `container card`: table, filter, form, and neutral content container
- Use clearer borders and layer contrast instead of heavy shadows.
- On tech-accent pages, add local title-bar glow, selected-edge emphasis, or subtle inner highlights rather than global page glow.

### 4. Border and Shadow Rules

- Favor crisp edge definition over thick drop shadows.
- Default pages:
  - fine border
  - light shadow
  - controlled corner radius
- Flagship pages:
  - fine outline
  - local highlight edge
  - minimal inner glow where useful

### 5. State Expression

- State should never rely on text color alone.
- Standardize state expression using:
  - label tone
  - optional icon
  - weak background
  - border or emphasis cue
  - ordering / placement in critical areas
- Priority states:
  - healthy / normal / enabled
  - attention / pending / at risk
  - alert / abnormal / blocking

### 6. Charts and Monitoring Surfaces

- Charts must optimize scanability and comparability.
- Keep chart palettes tighter and more intentional.
- Use container framing so title, filters, time-range switches, and charts read as one instrument unit.
- Monitoring pages can be more atmospheric, but chart readability remains the first constraint.

## Page-Level Guidance

### 1. Login

- Positioning:
  - enterprise product landing shell + focused login panel
- Recommended treatment:
  - stronger brand narrative
  - deeper blue background with restrained tech texture
  - clean, trustworthy login card
  - obvious captcha and feedback states
- Goal:
  - feel credible, modern, and operationally serious

### 2. Dashboard

- Positioning:
  - operational overview + risk focus + quick action hub
- Recommended treatment:
  - more obvious command-center flavor
  - stronger KPI hierarchy
  - better separation among risk, pending, AI/workbench, and chart zones
  - chart containers styled as instrument panels, not generic cards

### 3. Monitoring and Alert Center

- Positioning:
  - situational awareness + event handling
- Recommended treatment:
  - emphasize severity and next actions
  - unify trend panels, event lists, details, and action zones
  - allow medium-strength dark/blue atmosphere without becoming a demo wall

### 4. K8s Pages

- Positioning:
  - professional operations console
- Recommended treatment:
  - shared control-plane grammar across clusters, nodes, workloads, namespaces, network, storage, and monitoring pages
  - stronger hierarchy for state, resource usage, governance actions, and object relationships
  - stable toolbars, filters, tabs, and details layouts

### 5. List Pages

- Positioning:
  - mature enterprise workbench
- Recommended treatment:
  - stable page header
  - predictable filter region
  - optional stat strip where useful
  - strong table container
  - clean batch-action and pagination behavior
- Avoid:
  - heavy decorative backgrounds
  - inconsistent status tags
  - scattered action patterns

### 6. Form Pages

- Positioning:
  - low-cognitive-load configuration workflow
- Recommended treatment:
  - sectioned long forms
  - grouped fields
  - clear help text and validation rules
  - isolated warning / destructive zones
  - clear primary action ownership

### 7. Detail Pages

- Recommended shell:
  - summary header
  - key metrics and status
  - tabbed content body
- The top section should answer:
  - what is this object
  - what state is it in
  - what risk exists
  - what can the user do next

### 8. Logs / Terminal / Audit

- These pages can adopt a deeper and more tool-like surface.
- Priorities:
  - typography clarity
  - obvious boundary between controls and output
  - strong emphasis for anomalies, dangerous commands, and audit signals
  - fatigue-resistant long reading

## Component and Interaction Rules

### 1. Page Header

- Standard page header should support:
  - title
  - short context statement
  - key counts or status
  - primary actions
- Flagship pages may amplify the header visually, but the structure should remain consistent.

### 2. Button Hierarchy

- Standardize:
  - `primary`
  - `secondary`
  - `tertiary / ghost`
  - `danger`
- Not every action should look important.
- Primary CTA must remain visually dominant on dense surfaces.

### 3. Tables

- Treat tables as first-class platform components.
- Standardize:
  - row density
  - header contrast
  - hover feedback
  - state column rendering
  - action column rhythm
  - batch-operation placement
- Optimize for scanning speed and operator endurance.

### 4. Filter Regions

- Use one stable pattern:
  - common filters first
  - advanced filters foldable
  - query and reset positions stable
  - consistent control sizing

### 5. Forms

- Standardize:
  - label spacing
  - helper text placement
  - required-field markers
  - error presentation
  - grouped sections
- Reduce the chance of missed fields or accidental submissions.

### 6. Status Tags

- Centralize status vocabulary and status styling.
- Typical axes:
  - enabled / disabled
  - success / failed / running
  - healthy / at risk / abnormal
  - online / offline

### 7. Dialogs and Drawers

- Dialogs:
  - quick confirmation
  - short forms
  - lightweight edits
- Drawers:
  - side-context details
  - longer edits
  - workflows that should not disconnect the user from the page
- Avoid deep modal nesting.

### 8. Empty / Loading / Error States

- These states materially affect product quality and trust.
- Rules:
  - empty states should feel intentional, not missing
  - loading should not create excessive layout shift
  - errors should state the issue and a likely next step

### 9. Motion

- Use restrained motion only:
  - hover
  - selection
  - expand / collapse
  - tab / drawer transitions
- Avoid persistent glow, floating animation, or decorative movement.

### 10. Responsive Priority

- Optimize primarily for desktop operations, especially `1280px+` widths.
- Ensure medium-width laptops do not break layout hierarchy.
- Responsive support should protect usability rather than chase full mobile-admin parity.

## Rollout Strategy

### Recommended Approach

Use a **sample-page-first** rollout:

1. Establish design language on flagship pages:
   - `Login`
   - `Dashboard`
   - `Monitoring`
   - `K8s`
   - `List / Form skeleton`
2. Extract reusable tokens and shared patterns from those pages.
3. Extend the system to CMDB, config center, work orders, AI pages, and system management.

### Why This Approach

- It creates visible quality gains quickly.
- It avoids a long abstract design-system-only phase.
- It gives concrete UI references that later pages can follow.

## Reference Mapping

- Use `awesome-design-md` as a **reference input**, not a visual clone source.
- Recommended inspiration weighting:
  - `HashiCorp`: 70%
  - `Sentry`: 20%
  - `ClickHouse`: 10%
- Interpretation for OpsNexus:
  - HashiCorp provides mature product-surface discipline
  - Sentry provides clearer issue / event hierarchy
  - ClickHouse contributes sharper data-console framing where useful

## Prompt Alignment

Future AI-assisted UI changes should include these constraints:

- Preserve existing business logic, routes, permissions, APIs, and workflows.
- Improve layout hierarchy, visual consistency, and readability first.
- Avoid generic AI-looking admin design.
- Do not apply heavy sci-fi styling across every page.
- Use stronger blue tech accents only on flagship operational pages.
- Keep tables and forms highly usable for dense enterprise workflows.

## Acceptance Criteria

1. OpsNexus reads as one product instead of a collection of differently styled admin pages.
2. The platform default look is mature, durable, and professional.
3. Dashboard, monitoring, K8s, logs, and login have stronger identity without breaking daily usability.
4. List pages, forms, tables, and dialogs become more consistent and lower-friction.
5. Future UI optimization work can reference this document directly to stay aligned.

## Constraint Note

- The superpowers brainstorming flow recommends subagent-based spec review and git commit after writing the spec.
- In this session, subagent delegation has not been user-authorized, so the spec review is handled inline.
- The current working directory is not inside an initialized git repository, so this spec can be written locally but cannot be committed from the present location unless the actual git root is provided later.
