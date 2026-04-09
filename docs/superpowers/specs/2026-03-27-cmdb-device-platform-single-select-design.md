# CMDB Device Platform Single-Select Design

**Goal:** Change the CMDB network-device create/edit form so the `平台` field is selected from a fixed single-choice list instead of typed manually.

**Current Problem**

- The network-device create and edit dialogs currently render `platform` as a free-text input.
- Manual entry leads to inconsistent vendor names such as `H3C`, `h3c`, `华三`, or mixed custom text.
- The backend already stores `platform` as a single string, so a single-choice UI maps naturally to the current data shape.

**Chosen UX**

- Replace the `平台` input in both create and edit dialogs with a single-select radio group.
- The default option set will be:
  - `H3C`
  - `Cisco`
  - `Huawei`
  - `Juniper`
  - `Ruijie`
  - `其他`
- `平台` becomes a required field in both dialogs.

**Data Contract**

- Keep using the existing `platform` string field in frontend form state and backend payloads.
- The selected radio value is written directly to `form.platform`.
- No backend API or database change is required.

**Compatibility**

- Existing records with a known platform continue to round-trip normally.
- If an edited device has a legacy custom platform value that is not in the fixed option list, the frontend will append a temporary `历史值` option for that exact value so the dialog can still represent and preserve the existing record.

**Implementation Notes**

- Add a shared device-platform option helper in the CMDB asset presentation utility so create and edit dialogs stay aligned.
- Update create and edit dialog validation rules to require `platform`.
- Keep table rendering unchanged because it already reads the stored `platform` string.

**Testing**

- Add a focused frontend utility test that verifies:
  - the standard device-platform options are exposed in the expected order
  - a legacy platform value is appended as a temporary option
- Run the existing CMDB device presentation tests plus the new platform-option test after implementation.
