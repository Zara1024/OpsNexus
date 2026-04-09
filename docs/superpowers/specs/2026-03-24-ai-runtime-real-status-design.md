# AI Runtime Real Status Design

**Goal:** Make the AI runtime badge on `/ai/diagnosis` and `/ai/assistant` reflect real end-to-end availability instead of only checking whether an API key exists, and switch the remote runtime base URL from `https://uniquefox.top/v1` to `https://beecode.cc/v1`.

**Current Problem**

- The backend overview response marks runtime as `ready` whenever `runtimeClient.IsEnabled()` is true.
- `IsEnabled()` currently means “API key exists and provider is openai”.
- Real assistant and diagnosis requests can still fall back, as shown by `usedLlm=false`, `usedFallback=true`, and fallback errors such as `invalid character '<' looking for beginning of value`.
- This makes the green runtime badge misleading.

**Chosen Approach**

- Add a lightweight real probe to the OpenAI-compatible runtime client.
- Cache probe results briefly in-process so the overview page stays truthful without creating excessive traffic.
- Expose three runtime states:
  - `ready`: latest real probe succeeded
  - `degraded`: runtime is configured but latest real probe failed
  - `fallback`: runtime is disabled or missing required config
- Include the latest probe error/details in `overview.runtime`.
- Update diagnosis and assistant UIs to render the new state truthfully.

**Why this approach**

- It matches actual user experience better than config-only readiness.
- It keeps the existing fallback execution path intact.
- It avoids tying status to one specific page’s last request by using a shared probe result.

**Scope**

- Backend:
  - extend runtime client with cached probe support
  - extend AI overview snapshot/response with probe-derived fields
- Frontend:
  - show degraded/fallback/ready truthfully on diagnosis and assistant pages
- Environment:
  - replace remote `baseUrl` with `https://beecode.cc/v1`

**Non-goals**

- Reworking the AI orchestration flow itself
- Replacing fallback behavior
- Adding persistent runtime-health storage to the database

**Files**

- Modify: `api/api/ai/service/runtime_openai.go`
- Modify: `api/api/ai/service/overview.go`
- Modify: `api/api/ai/model/overview.go`
- Modify: `web/src/views/ai/AIDiagnosis.vue`
- Modify: `web/src/views/ai/AIAssistant.vue`
- Modify: `/opt/opsnexus-remote-test/config.yaml` on remote host
- Create: `api/api/ai/service/runtime_openai_test.go`

**Verification**

- `go test ./api/ai/service -run Test.*Runtime.* -v`
- `cd web && npm run build`
- remote `GET /api/v1/ai/overview` shows `ready` only after real probe success
- remote real AI request confirms whether `usedLlm` changes from `false` to `true`

**Constraint Note**

- The superpowers spec-review loop recommends a reviewer subagent, but subagent delegation is not user-authorized in this session, so execution proceeds inline with local verification instead.
