# CMDB Database Operation PostgreSQL And Redis Support Design

**Goal:** Extend the existing `数据库操作` and `SQL工单` workflow so PostgreSQL and Redis assets can use the same unified entry flow as MySQL, while Redis blocks high-risk commands by default.

## Current State

- Frontend database operation page lives at `web/src/views/cmdb/DBdetails.vue`.
- The page currently assumes a MySQL-style mental model:
  - SQL type selector is fixed to `select / insert / update / delete / raw`
  - direct execution always calls MySQL-oriented APIs
  - work-order submission treats the payload as SQL-only
- Backend execution is hard-coded to MySQL in:
  - `api/api/cmdb/controller/cmdbSQLRecord.go`
  - `api/api/cmdb/service/sqlWorkOrder.go`
- Both execution and work-order flows currently reject non-MySQL assets with `only MySQL databases are currently supported`.
- CMDB database assets already classify types as:
  - `1 = MySQL`
  - `2 = PostgreSQL`
  - `3 = Redis`

## Approved Direction

Use one unified page and one unified work-order center, but dispatch by database type underneath:

- MySQL:
  - keep current behavior
- PostgreSQL:
  - support direct execution
  - support work-order creation, approval, and execution
- Redis:
  - support direct command execution
  - support work-order creation, approval, and execution through the same work-order center
  - block high-risk commands by default, even if submitted through a work order

This keeps the product experience unified while allowing type-specific execution rules.

## Frontend Design

### Database Operation Page

File: `web/src/views/cmdb/DBdetails.vue`

The page remains the single entry point for database operations.

Behavior by type:

- MySQL:
  - keep existing SQL type selector and SQL text area
- PostgreSQL:
  - reuse the SQL type selector and SQL text area
  - update placeholder text and operation hints to be PostgreSQL-friendly where needed
- Redis:
  - keep the page layout and buttons unified
  - replace the MySQL-style SQL type selector with a Redis command mode selector
  - command mode categories:
    - `read`
    - `write`
    - `raw`
  - placeholder examples should become Redis command examples such as:
    - `GET key`
    - `SET key value`
    - `HGETALL hash_key`

### Submission Rules

- MySQL / PostgreSQL:
  - `select` stays direct-execution only
  - change statements can be submitted to work order
- Redis:
  - read commands can be executed directly
  - write commands can be submitted to work order
  - blocked commands must be rejected before execution and before work-order execution

### Result Presentation

The right-side result panel remains unified.

- SQL databases return structured query/mutation results
- Redis returns normalized command results so the existing `executionResult` panel can still display JSON text safely

## Backend Design

## Execution Controller Split

File: `api/api/cmdb/controller/cmdbSQLRecord.go`

Refactor the current MySQL-only logic into type-dispatched execution helpers:

- MySQL executor:
  - keep existing MySQL paths
- PostgreSQL executor:
  - add PostgreSQL connection builder
  - add database/schema discovery logic
  - add query execution
  - add mutation execution
- Redis executor:
  - add Redis client/session creation
  - add command parsing
  - add command execution
  - normalize returned values into frontend-friendly JSON

Proposed shape:

- `executeSelectRequest` becomes type-aware or delegates to per-type helpers
- `executeMutationRequest` becomes type-aware or delegates to per-type helpers
- `GetDatabaseList` becomes type-aware:
  - MySQL: `SHOW DATABASES`
  - PostgreSQL: database or schema list according to current asset semantics
  - Redis: return a synthetic list or single logical target if no schema concept exists

## Work Order Service Split

File: `api/api/cmdb/service/sqlWorkOrder.go`

Refactor work-order creation and execution into type-aware flows:

- MySQL work order:
  - keep existing behavior
- PostgreSQL work order:
  - reuse SQL classification and execution model with PostgreSQL-specific validation
- Redis work order:
  - treat work orders as command work orders inside the same center
  - keep the same list/detail/approval/execution lifecycle
  - use Redis command classification instead of SQL parsing for risk and validation

## Redis Risk Model

Redis commands must be classified into:

- safe read commands
  - examples: `GET`, `MGET`, `HGET`, `HGETALL`, `LRANGE`, `SMEMBERS`, `ZRANGE`, `TTL`, `EXISTS`
- normal write commands
  - examples: `SET`, `DEL`, `HSET`, `LPUSH`, `RPUSH`, `SADD`, `ZADD`, `EXPIRE`
- blocked high-risk commands
  - examples:
    - `FLUSHALL`
    - `FLUSHDB`
    - `SHUTDOWN`
    - `CONFIG SET`

Blocked commands:

- must be rejected during direct execution
- must be rejected during work-order execution
- should return explicit error messages indicating that the command is disallowed by policy

## Shared Work Order Model

Keep the current SQL work-order center instead of building a separate Redis center.

Compatibility approach:

- preserve existing storage model where possible
- keep `SQLContent` as the raw command text field even for Redis commands, to avoid a schema migration unless absolutely necessary
- broaden user-facing semantics from "SQL" to "database change command" where appropriate in logic and messages
- allow rollback-related fields to degrade gracefully:
  - MySQL / PostgreSQL can still generate rollback SQL or rollback hints
  - Redis uses human-readable recovery guidance instead of rollback SQL where exact rollback is not derivable

## Validation Rules

### MySQL

- unchanged

### PostgreSQL

- allow:
  - `SELECT`
  - `INSERT`
  - `UPDATE`
  - `DELETE`
  - `CREATE`
  - `ALTER`
  - `DROP`
  - `TRUNCATE`
- use PostgreSQL connection and database existence validation

### Redis

- parse the first command token
- map to `read / write / blocked`
- reject blocked commands with explicit policy errors
- work-order approval remains available for write commands

## Testing Plan

### Frontend

- add tests for `DBdetails.vue`:
  - PostgreSQL mode renders executable workflow without MySQL-only rejection
  - Redis mode renders command-oriented placeholders/options
  - submit button state and warning messages change by type

### Backend

- add controller/service tests for:
  - PostgreSQL direct execution
  - PostgreSQL work-order creation
  - PostgreSQL work-order execution
  - Redis read command execution
  - Redis write command work-order creation/execution
  - Redis blocked command rejection

### Regression

- MySQL direct execution still passes
- MySQL work-order flow still passes
- existing `数据库操作` route remains compatible

## Out Of Scope

- MongoDB support
- Elasticsearch support
- redesigning the SQL work-order center into a new database-command center UI
- introducing a brand new Redis shell or terminal page

## Risks

- PostgreSQL schema selection may differ from current MySQL-oriented `databaseName` semantics and needs a clear mapping.
- Redis does not naturally fit SQL-only terminology, so some field names remain implementation-compatible even if product wording becomes more generic.
- If current backend account credentials are not stored in a way compatible with PostgreSQL or Redis connection requirements, credential normalization may be needed during implementation.
