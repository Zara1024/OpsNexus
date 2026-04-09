# CMDB Database Operation PostgreSQL And Redis Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Extend the existing `数据库操作` and `SQL工单` flow so PostgreSQL and Redis assets can use the same unified execution and work-order path as MySQL, while Redis blocks high-risk commands by default.

**Architecture:** Keep the current routes and page structure, but split execution by database type underneath. Frontend stays on one `DBdetails.vue` page and one work-order center; backend dispatches MySQL, PostgreSQL, and Redis to separate execution helpers while preserving the current DAO and request contracts where possible.

**Tech Stack:** Vue 3, Element Plus, Go, Gin, `database/sql`, MySQL driver, new PostgreSQL driver, Redis client, node-based frontend assertions, Go unit tests.

---

### File Structure

**Frontend**
- Modify: `web/src/views/cmdb/DBdetails.vue`
- Modify: `web/src/api/cmdb.js` if request payloads need small type-aware additions
- Create: `web/tests/cmdb-dbdetails-type-support.test.mjs`

**Backend**
- Modify: `api/api/cmdb/controller/cmdbSQLRecord.go`
- Modify: `api/api/cmdb/service/sqlWorkOrder.go`
- Modify: `api/go.mod`
- Create: `api/api/cmdb/service/cmdbSQL_type_support.go`
- Create: `api/api/cmdb/service/cmdbSQL_type_support_test.go`
- Modify: `api/api/cmdb/service/sqlWorkOrder_test.go`

These files split cleanly by responsibility:
- frontend page behavior
- backend direct execution and connection helpers
- backend work-order dispatch and risk policy
- focused tests for each layer

### Task 1: Add failing frontend tests for PostgreSQL and Redis mode behavior

**Files:**
- Create: `web/tests/cmdb-dbdetails-type-support.test.mjs`
- Modify: `web/src/views/cmdb/DBdetails.vue`

- [ ] **Step 1: Write the failing test**

Add assertions that `DBdetails.vue`:
- contains PostgreSQL-aware execution support
- contains Redis-specific command-mode support
- does not hard-code the page into a MySQL-only warning path

Use source assertions like:

```js
assert.match(source, /isRedisDatabaseType|redisCommandMode|resolveExecutionMode/)
assert.match(source, /isPostgreSQLDatabaseType|PostgreSQL/)
assert.doesNotMatch(source, /only MySQL databases are currently supported/)
```

- [ ] **Step 2: Run test to verify it fails**

Run: `node .\tests\cmdb-dbdetails-type-support.test.mjs`
Expected: FAIL because the current page still assumes MySQL-style execution only.

- [ ] **Step 3: Write minimal implementation**

Update `DBdetails.vue` so the page:
- resolves the current database type from `dbInfo.type`
- renders PostgreSQL using the current SQL-style form
- renders Redis with command-mode options and Redis-specific placeholders

- [ ] **Step 4: Run test to verify it passes**

Run: `node .\tests\cmdb-dbdetails-type-support.test.mjs`
Expected: PASS

### Task 2: Add failing backend tests for type-dispatched direct execution helpers

**Files:**
- Create: `api/api/cmdb/service/cmdbSQL_type_support.go`
- Create: `api/api/cmdb/service/cmdbSQL_type_support_test.go`
- Modify: `api/api/cmdb/controller/cmdbSQLRecord.go`
- Modify: `api/go.mod`

- [ ] **Step 1: Write the failing test**

Add Go tests for:
- PostgreSQL operation type resolution
- Redis command classification into `read / write / blocked`
- blocked Redis commands such as `FLUSHALL`, `FLUSHDB`, `SHUTDOWN`, `CONFIG SET`

Example:

```go
func TestClassifyRedisCommand(t *testing.T) {
    kind, blocked := classifyRedisCommand("FLUSHALL")
    if kind != redisCommandBlocked || !blocked {
        t.Fatalf("expected FLUSHALL to be blocked")
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./api/cmdb/service -run "TestClassifyRedisCommand|TestResolveDatabaseExecutionType"`
Expected: FAIL because helper functions do not exist yet.

- [ ] **Step 3: Write minimal implementation**

Add `cmdbSQL_type_support.go` with:
- database type constants/helpers
- PostgreSQL connection builder signature
- Redis command classifier
- blocked-command policy

Also add the PostgreSQL driver dependency to `api/go.mod`.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./api/cmdb/service -run "TestClassifyRedisCommand|TestResolveDatabaseExecutionType"`
Expected: PASS

### Task 3: Refactor direct execution controller to dispatch by database type

**Files:**
- Modify: `api/api/cmdb/controller/cmdbSQLRecord.go`
- Test: `api/api/cmdb/service/cmdbSQL_type_support_test.go`

- [ ] **Step 1: Write the failing test**

Add tests covering:
- PostgreSQL requests are no longer rejected by `dbInfo.Type != 1`
- Redis blocked commands are rejected with policy errors

If controller-level unit tests are too heavy, add service/helper tests that prove dispatch output and blocked-policy decisions.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./api/cmdb/service -run "TestPostgreSQLIsNotRejected|TestBlockedRedisCommandsAreRejected"`
Expected: FAIL while MySQL-only logic still remains.

- [ ] **Step 3: Write minimal implementation**

Refactor `cmdbSQLRecord.go` to:
- preserve MySQL behavior
- add PostgreSQL branches for database listing, query execution, and mutation execution
- add Redis execution for direct commands
- normalize Redis results into JSON-safe response payloads

Keep current route contracts:
- `/cmdb/sql/select`
- `/cmdb/sql`
- `/cmdb/sql/execute`
- `/cmdb/sql/databaselist`

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./api/cmdb/service -run "TestPostgreSQLIsNotRejected|TestBlockedRedisCommandsAreRejected"`
Expected: PASS

### Task 4: Refactor SQL work-order service into type-aware work orders

**Files:**
- Modify: `api/api/cmdb/service/sqlWorkOrder.go`
- Modify: `api/api/cmdb/service/sqlWorkOrder_test.go`

- [ ] **Step 1: Write the failing test**

Add tests for:
- PostgreSQL work-order creation allowed for change statements
- Redis work-order creation allowed for safe write commands
- Redis blocked commands rejected during creation and execution

Example:

```go
func TestRedisBlockedCommandIsRejected(t *testing.T) {
    if err := validateRedisWorkOrderCommand("FLUSHALL"); err == nil {
        t.Fatalf("expected blocked redis command to be rejected")
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./api/cmdb/service -run "TestRedisBlockedCommandIsRejected|TestAnalyzeSQLWorkOrder"`
Expected: FAIL while work-order service still rejects non-MySQL targets.

- [ ] **Step 3: Write minimal implementation**

Refactor `sqlWorkOrder.go` so:
- MySQL and PostgreSQL use SQL statement analysis
- Redis uses command classification
- blocked Redis commands fail fast
- rollback fields degrade gracefully for Redis with recovery guidance instead of rollback SQL

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./api/cmdb/service -run "TestRedisBlockedCommandIsRejected|TestAnalyzeSQLWorkOrder"`
Expected: PASS

### Task 5: Update the database operation page to use type-aware execution and work-order rules

**Files:**
- Modify: `web/src/views/cmdb/DBdetails.vue`
- Test: `web/tests/cmdb-dbdetails-type-support.test.mjs`

- [ ] **Step 1: Write the failing test**

Extend the new frontend test so it asserts:
- PostgreSQL still exposes SQL execution and work-order submission
- Redis uses command wording and type-aware placeholders
- Redis prevents select-only SQL assumptions from leaking into the UI

- [ ] **Step 2: Run test to verify it fails**

Run: `node .\tests\cmdb-dbdetails-type-support.test.mjs`
Expected: FAIL until the page reacts to `dbInfo.type`.

- [ ] **Step 3: Write minimal implementation**

Add computed helpers in `DBdetails.vue` such as:

```js
const isPostgreSQLDatabaseType = computed(() => Number(dbInfo.value.type) === 2)
const isRedisDatabaseType = computed(() => Number(dbInfo.value.type) === 3)
```

Update:
- selector options
- placeholder text
- direct execution routing
- work-order submit warnings
- execution result copy

- [ ] **Step 4: Run test to verify it passes**

Run: `node .\tests\cmdb-dbdetails-type-support.test.mjs`
Expected: PASS

### Task 6: Run full verification

**Files:**
- Generated: `web/dist`

- [ ] **Step 1: Run frontend regression tests**

Run:
- `node .\tests\cmdb-dbdetails-type-support.test.mjs`
- `node .\tests\cmdb-database-asset-presentation.test.mjs`

Expected: PASS

- [ ] **Step 2: Run backend service tests**

Run:
- `go test ./api/cmdb/service`

Expected: PASS

- [ ] **Step 3: Run frontend production build**

Run:
- `npm run build`

Expected: successful Vue production build with exit code `0`

- [ ] **Step 4: Manual verification**

Verify:
- MySQL asset still works
- PostgreSQL asset no longer hits MySQL-only rejection
- Redis asset supports command execution
- blocked Redis commands are rejected clearly

- [ ] **Step 5: Commit**

```bash
git add web/src/views/cmdb/DBdetails.vue web/tests/cmdb-dbdetails-type-support.test.mjs api/api/cmdb/controller/cmdbSQLRecord.go api/api/cmdb/service/sqlWorkOrder.go api/api/cmdb/service/cmdbSQL_type_support.go api/api/cmdb/service/cmdbSQL_type_support_test.go api/api/cmdb/service/sqlWorkOrder_test.go api/go.mod
git commit -m "feat: support postgres and redis db operations"
```
