# OpsNexus Production Readiness Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Move OpsNexus from a demo/acceptance deployment posture to a production-safe baseline by closing critical security exposure, fixing credential and password handling, and hardening the deployment/runtime configuration.

**Architecture:** Treat this as a phased hardening effort instead of one bulk refactor. First eliminate externally reachable unsafe defaults and rotate secrets, then repair core authentication/password logic and deployment defaults in code, and finally add production verification, release guardrails, and operational documentation so the system stays safe after the first fix.

**Tech Stack:** Go, Gin, Vue 3, MySQL, Redis, Nginx, systemd, Docker Compose, Helm, Linux host deployment.

---

## Priority Guide

### P0: Must finish before production exposure

- External API debug surface must be closed.
- All hardcoded or leaked secrets must be removed from repo and rotated on the remote host.
- Password storage and password-change logic must be corrected.
- JWT and auth middleware must stop using unsafe defaults and secret-leaking logs.
- Remote service must stop running in its current root + open-port + manual-overlay posture.

### P1: Strongly recommended before first production tenant

- Seed/demo data and default weak accounts must be removed from shipping paths.
- Deployment manifests must converge on one production path instead of multiple conflicting defaults.
- Security-sensitive endpoints must have regression tests and smoke checks.

### P2: Operational maturity after baseline hardening

- Add release checklist, rollback procedure, and production health verification.
- Add secret/bootstrap runbook and controlled admin initialization flow.
- Add CI checks that block future reintroduction of debug mode and hardcoded credentials.

---

### Task 1: Lock down externally reachable debug surfaces

**Files:**
- Modify: `api/config.yaml`
- Modify: `api/main.go`
- Modify: `api/router/router.go`
- Modify: `deploy/helm/opsnexus/values.yaml`
- Modify: `docker/api/config.yaml`
- Modify: `docker/docker-compose.yml`
- Modify remote: `/opt/opsnexus-remote-test/config.yaml`
- Modify remote: `/etc/systemd/system/opsnexus-api.service`

- [ ] **Step 1: Add a backend config test that asserts production-safe defaults**

Target behavior:
- `server.model` defaults to `release`
- `enableSwagger` defaults to `false`
- config/env overrides still work intentionally

Run: `cd api && go test ./common/config -run TestLoadConfig -v`
Expected: FAIL until defaults are tightened.

- [ ] **Step 2: Change code and sample configs to ship with production-safe defaults**

Implement:
- switch default server mode to `release`
- keep Swagger opt-in only
- remove dev-oriented startup log copy that claims development mode on server boot

- [ ] **Step 3: Remove direct public exposure of port `8000` on the remote host**

Implement on `10.0.0.200`:
- bind API to `127.0.0.1:8000` or firewall it
- keep external access through `8080` reverse proxy only

Verify:
- `curl http://10.0.0.200:8000/swagger/index.html` should fail externally
- `curl http://10.0.0.200:8080/api/v1/captcha` should still return `200`

- [ ] **Step 4: Re-run backend config tests**

Run: `cd api && go test ./common/config ./router -v`
Expected: PASS

### Task 2: Remove hardcoded secrets from code and rotate all exposed credentials

**Files:**
- Modify: `api/config.yaml`
- Modify: `docker/.env`
- Modify: `docker/docker-compose.yml`
- Modify: `deploy/helm/opsnexus/values.yaml`
- Modify: `api/common/config/config.go`
- Modify: `README.md`
- Modify or create: `deploy/systemd/opsnexus-api.env`
- Modify remote: `/opt/opsnexus-remote-test/config.yaml`
- Modify remote: `/etc/systemd/system/opsnexus-api.service`

- [ ] **Step 1: Add/extend config tests for env-only secret overrides**

Target secrets:
- DB password
- Redis password
- JWT secret
- AI API key
- monitor heartbeat token
- monitor webhook token

Run: `cd api && go test ./common/config -v`
Expected: FAIL until missing secret fields and overrides are implemented.

- [ ] **Step 2: Add explicit config support for `JWT_SECRET` and fail-safe secret loading**

Implement:
- do not hardcode production secrets in YAML
- prefer env or external secret file loading for all sensitive values
- keep example configs blank or placeholder-only

- [ ] **Step 3: Replace committed secrets with placeholders and update docs**

Implement:
- remove committed passwords from repo-owned config samples
- document required env vars for local/dev/prod
- document that production secrets must be injected by systemd env file, Docker env, or Kubernetes Secret

- [ ] **Step 4: Rotate every secret already exposed**

Remote actions on `10.0.0.200`:
- rotate MySQL password
- rotate Redis password
- rotate webhook token
- rotate agent heartbeat token
- rotate JWT secret
- revoke and replace the AI API key currently present in remote config

Verify:
- service restarts cleanly
- login still works
- monitoring webhooks and agent heartbeat still authenticate with new values

### Task 3: Replace MD5 password storage and repair password flows

**Files:**
- Modify: `api/common/util/encryption.go`
- Modify: `api/api/system/service/sysAdmin.go`
- Modify: `api/api/system/dao/sysAdmin.go`
- Modify: `api/api/system/model/sysAdmin.go`
- Create: `api/api/system/service/sysAdmin_password_test.go`
- Create: `api/api/system/dao/sysAdmin_password_test.go`
- Modify: `web/src/views/system/Admin.vue`
- Modify: `web/src/views/system/Personal.vue`

- [ ] **Step 1: Write failing tests for login, reset-password, and personal password change**

Cover:
- stored password hash is not raw plaintext
- login succeeds with new hash format
- admin reset writes hash, not plaintext
- personal password update verifies old password and writes hash

Run: `cd api && go test ./api/system/... -v`
Expected: FAIL

- [ ] **Step 2: Replace `MD5` password handling with bcrypt**

Implement:
- add `HashPassword` / `ComparePassword`
- remove `EncryptionMd5` from auth-sensitive password paths
- keep any unrelated legacy encryption utilities isolated if still needed elsewhere

- [ ] **Step 3: Fix the self-service password change bug**

Implement:
- verify current password before allowing change
- require `newPassword == resetPassword`
- hash the new password before save

- [ ] **Step 4: Remove unsafe UI behavior that reveals reset passwords in clear text**

Implement:
- success toast should not display the password value
- admin reset flow should use one-time operator entry only

- [ ] **Step 5: Run targeted tests and frontend build**

Run:
- `cd api && go test ./api/system/... -v`
- `cd web && npm run build`

Expected: PASS

### Task 4: Harden JWT handling and auth middleware

**Files:**
- Modify: `api/go.mod`
- Modify: `api/pkg/jwt/jwt.go`
- Modify: `api/middleware/authMiddleware.go`
- Create: `api/pkg/jwt/jwt_test.go`
- Create: `api/middleware/authMiddleware_test.go`

- [ ] **Step 1: Write failing tests for JWT secret loading and middleware behavior**

Cover:
- token validation fails with wrong secret
- middleware never logs the secret
- non-WebSocket HTTP requests do not accept query-token fallback
- WebSocket auth behavior remains explicitly supported if still required

Run: `cd api && go test ./pkg/jwt ./middleware -v`
Expected: FAIL

- [ ] **Step 2: Migrate off `github.com/dgrijalva/jwt-go`**

Implement:
- move to maintained `github.com/golang-jwt/jwt/v5`
- load JWT secret from config/env
- support secret rotation via restart

- [ ] **Step 3: Remove secret-bearing debug logs and tighten token transport**

Implement:
- delete log lines that print secret values or token metadata
- restrict query-token support to explicitly required websocket/SSE paths only
- keep normal API auth on `Authorization: Bearer ...`

- [ ] **Step 4: Re-run auth tests**

Run: `cd api && go test ./pkg/jwt ./middleware -v`
Expected: PASS

### Task 5: Tighten CORS, webhook auth, and internal token defaults

**Files:**
- Modify: `api/middleware/cors.go`
- Modify: `api/api/monitor/controller/alert.go`
- Modify: `api/api/monitor/controller/agent.go`
- Modify: `api/common/agent/agent.go`
- Create: `api/api/monitor/controller/alert_auth_test.go`

- [ ] **Step 1: Write failing tests for CORS and token fallback behavior**

Cover:
- wildcard origin is not allowed with credentials
- webhook and heartbeat reject missing or default token usage in production mode

Run: `cd api && go test ./api/monitor/... ./middleware -v`
Expected: FAIL

- [ ] **Step 2: Replace permissive CORS with explicit origin allowlist support**

Implement:
- configurable allowed origins
- disable credentials unless actually required
- preserve local dev convenience through explicit dev config only

- [ ] **Step 3: Remove dangerous hardcoded fallback tokens**

Implement:
- heartbeat and webhook endpoints must require configured secret values
- no silent fallback to `agent-heartbeat-token-2024` or `webhook-notify-token-2024` in production

- [ ] **Step 4: Re-run tests**

Run: `cd api && go test ./api/monitor/... ./middleware -v`
Expected: PASS

### Task 6: Clean shipping data and weak default accounts from bootstrap paths

**Files:**
- Modify: `docker/mysql/devops001.sql`
- Modify: `docker/mysql/init.sql`
- Modify: `api/sql/autoops.sql`
- Modify: `README.md`
- Create or modify: `scripts/bootstrap-admin.*`

- [ ] **Step 1: Audit which SQL files are actually used for bootstrap**

Confirm:
- which SQL file seeds Docker
- which SQL file seeds local/manual environments
- whether any remote environment was initialized from those files

- [ ] **Step 2: Remove weak default passwords and demo users from shipping SQL**

Implement:
- no `admin/123456`
- no broad demo users in production bootstrap
- if bootstrap admin is required, generate it through controlled env/script path

- [ ] **Step 3: Add a safe initialization workflow**

Implement:
- initialize first admin through env-driven bootstrap script or one-time CLI
- force password change on first login if a bootstrap admin is created

- [ ] **Step 4: Verify fresh bootstrap behavior**

Run:
- bring up a fresh local stack
- confirm login works with the intentional bootstrap path only
- confirm removed weak accounts no longer exist

### Task 7: Harden systemd and host deployment shape

**Files:**
- Modify: `deploy/systemd/opsnexus-api.service`
- Create: `deploy/systemd/opsnexus-api.env.example`
- Create or modify: `scripts/deploy-remote.*`
- Modify remote: `/etc/systemd/system/opsnexus-api.service`
- Modify remote: `/opt/opsnexus-remote-test/`

- [ ] **Step 1: Replace root service execution with a dedicated service user**

Implement:
- create `opsnexus` system user/group
- chown runtime directories appropriately
- run service under least privilege

- [ ] **Step 2: Move secrets out of the unit file and config blob**

Implement:
- use `EnvironmentFile=` for secrets on host deployment
- keep binary and config directories separated from secret storage

- [ ] **Step 3: Clean remote deployment directory and stop uncontrolled backup sprawl**

Implement:
- keep one current release, one previous release, one logs directory
- move backups to bounded archival location
- remove world-writable permissions from web assets

- [ ] **Step 4: Verify restart, boot persistence, and file ownership**

Run on remote:
- `systemctl daemon-reload`
- `systemctl restart opsnexus-api.service`
- `systemctl status opsnexus-api.service --no-pager`
- `ss -lntp | egrep ':(8000|8080)\\b'`

Expected:
- service active
- `8000` not publicly exposed
- service not running as `root`

### Task 8: Converge deployment defaults across Docker, Helm, and host mode

**Files:**
- Modify: `docker/docker-compose.yml`
- Modify: `docker/README.md`
- Modify: `deploy/helm/opsnexus/values.yaml`
- Modify: `deploy/helm/README.md`
- Modify: `README.md`

- [ ] **Step 1: Choose and document one production baseline per deployment mode**

Modes:
- host + systemd
- Docker Compose
- Helm/Kubernetes

- [ ] **Step 2: Make production-safe defaults consistent**

Implement:
- no debug mode by default
- no committed passwords
- no contradictory Redis password defaults between files
- no public raw API port in sample production topology

- [ ] **Step 3: Add deployment-specific verification commands to docs**

Examples:
- host mode curl checks
- compose health checks
- Helm readiness and smoke checks

- [ ] **Step 4: Rebuild and validate artifacts**

Run:
- `cd api && go test ./...`
- `cd web && npm run build`
- if using Docker mode, `docker compose config`
- if using Helm mode, `helm template`

Expected: PASS

### Task 9: Add release gate and production verification checklist

**Files:**
- Create: `docs/manual-checklists/2026-04-07-production-go-live-checklist.md`
- Modify: `README.md`
- Create or modify: `scripts/verify-production-readiness.*`

- [ ] **Step 1: Capture a go-live checklist**

Include:
- secrets injected
- Swagger disabled
- API port not public
- JWT secret rotated
- admin bootstrap confirmed
- backup and rollback confirmed
- logs and monitoring checked

- [ ] **Step 2: Script the critical smoke verification**

Checks:
- login page reachable on `8080`
- captcha reachable
- authenticated dashboard reachable
- selected protected API requires auth
- swagger path unavailable
- remote API raw port unavailable externally

- [ ] **Step 3: Run final verification on `10.0.0.200`**

Minimum checks:
- `curl -I http://10.0.0.200:8080/login`
- `curl http://10.0.0.200:8080/api/v1/captcha`
- `curl http://10.0.0.200:8080/api/v1/config/accountauth/list` returns `403` without auth
- `curl http://10.0.0.200:8000/swagger/index.html` fails externally
- browser login with intended admin works

- [ ] **Step 4: Update release notes with exact hardening changes**

Document:
- rotated secrets
- changed auth behavior
- required operator actions after deploy

---

## Recommended Execution Order

1. Task 1
2. Task 2
3. Task 3
4. Task 4
5. Task 5
6. Task 7
7. Task 6
8. Task 8
9. Task 9

## Notes

- The current remote environment on `10.0.0.200` should be treated as already compromised from a secret-hygiene perspective until rotation is complete.
- Do not treat “site loads and can log in” as production readiness; the acceptance bar is the checklist in Task 9.
- If execution becomes too wide for one pass, split follow-up implementation into two branches:
  - branch A: auth + secret hardening
  - branch B: deployment + runbook hardening
