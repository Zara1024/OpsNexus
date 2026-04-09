# OpsNexus Demo Full Regression Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Execute every test case under `C:\zq\平台开发\OpsNexus\demo` against `http://10.0.0.200:8080`, preserve test data, fix product defects before continuing, and sync CSV/XLSX execution tables after each completed case.

**Architecture:** Treat the Markdown case files in `demo` as the source of truth for execution order and expected behavior. Use the OpsNexus browser automation skill to authenticate against the shared remote environment, collect repeatable artifacts, and drive page-level regression. Use `scripts/update_execution_tables.py` as the only writer for execution results so CSV and XLSX stay consistent after every case. When a case reveals a product defect, stop the current slice, isolate the failing route and code path, apply the smallest safe fix, rerun the affected case(s), then continue the full pass.

**Tech Stack:** PowerShell, Python 3.11, Playwright-based browser scripts, Vue frontend under `web`, Go backend under `api`, demo Markdown/CSV/XLSX assets.

---

### Task 1: Prepare the execution toolchain

**Files:**
- Modify: `C:\zq\平台开发\OpsNexus\OpsNexus\docs\superpowers\plans\2026-03-24-opsnexus-demo-full-regression.md`
- Read: `C:\zq\平台开发\OpsNexus\demo\00-先看这个-怎么执行.md`
- Read: `C:\zq\平台开发\OpsNexus\demo\00-test-execution-guide.md`
- Read: `C:\zq\平台开发\OpsNexus\OpsNexus\scripts\update_execution_tables.py`
- Test: `C:\zq\平台开发\OpsNexus\OpsNexus\scripts\test_update_execution_tables.py`

- [ ] **Step 1: Confirm the execution assets and counts**

Run: `rg --files C:\zq\平台开发\OpsNexus\demo`
Expected: all Markdown case files plus the three CSV tables and one XLSX workbook are listed.

- [ ] **Step 2: Verify the table updater before touching live records**

Run: `python -m unittest C:\zq\平台开发\OpsNexus\OpsNexus\scripts\test_update_execution_tables.py`
Expected: `Ran 2 tests` and `OK`.

- [ ] **Step 3: Bootstrap the browser runtime**

Run: `powershell -ExecutionPolicy Bypass -File C:\Users\zq\.codex\skills\opsnexus-browser-ui-adjuster\scripts\bootstrap-runtime.ps1`
Expected: Playwright/browser runtime is ready without install errors.

### Task 2: Validate remote login and artifact capture

**Files:**
- Read: `C:\Users\zq\.codex\skills\opsnexus-browser-ui-adjuster\references\opsnexus-defaults.md`
- Create: `C:\zq\平台开发\OpsNexus\OpsNexus\tmp\opsnexus-browser-artifacts\`

- [ ] **Step 1: Probe the shared remote environment**

Run: `Invoke-WebRequest http://10.0.0.200:8080/api/v1/captcha -UseBasicParsing`
Expected: HTTP 200 and captcha payload or image response.

- [ ] **Step 2: Log in as the admin test user and save browser state**

Run: `node C:\Users\zq\.codex\skills\opsnexus-browser-ui-adjuster\scripts\opsnexus-browser.mjs login --base-url http://10.0.0.200:8080 --username admin --password <ADMIN_PASSWORD> --artifacts-dir C:\zq\平台开发\OpsNexus\OpsNexus\tmp\opsnexus-browser-artifacts`
Expected: login summary JSON and storage-state file are written; browser lands on `/home` or `/dashboard`.

- [ ] **Step 3: Inspect a known page using the authenticated session**

Run: `node C:\Users\zq\.codex\skills\opsnexus-browser-ui-adjuster\scripts\opsnexus-browser.mjs inspect --base-url http://10.0.0.200:8080 --state-file C:\zq\平台开发\OpsNexus\OpsNexus\tmp\opsnexus-browser-artifacts\opsnexus-storage-state.json --route /dashboard --selector .layout-container --wait-ms 1500`
Expected: screenshot, HTML snapshot, and JSON page summary are written successfully.

### Task 3: Execute all base-suite cases and record every result immediately

**Files:**
- Read: `C:\zq\平台开发\OpsNexus\demo\01-仪表盘\01-仪表盘-测试用例.md`
- Read: `C:\zq\平台开发\OpsNexus\demo\02-资产管理\*.md`
- Read: `C:\zq\平台开发\OpsNexus\demo\03-容器管理\*.md`
- Read: `C:\zq\平台开发\OpsNexus\demo\04-服务管理\*.md`
- Read: `C:\zq\平台开发\OpsNexus\demo\05-任务中心\*.md`
- Read: `C:\zq\平台开发\OpsNexus\demo\06-AI智能运维助手\*.md`
- Read: `C:\zq\平台开发\OpsNexus\demo\07-运维工具\*.md`
- Read: `C:\zq\平台开发\OpsNexus\demo\08-工单中心\*.md`
- Read: `C:\zq\平台开发\OpsNexus\demo\09-知识库\*.md`
- Read: `C:\zq\平台开发\OpsNexus\demo\10-监控告警\*.md`
- Read: `C:\zq\平台开发\OpsNexus\demo\11-操作审计\*.md`
- Read: `C:\zq\平台开发\OpsNexus\demo\12-系统管理\*.md`
- Read: `C:\zq\平台开发\OpsNexus\demo\13-配置中心\*.md`
- Read: `C:\zq\平台开发\OpsNexus\demo\14-全局搜索\*.md`
- Modify: `C:\zq\平台开发\OpsNexus\demo\opsnexus-base-test-execution.csv`
- Modify: `C:\zq\平台开发\OpsNexus\demo\opsnexus-all-test-execution.csv`
- Modify: `C:\zq\平台开发\OpsNexus\demo\opsnexus-test-execution.xlsx`

- [ ] **Step 1: Execute each base-case Markdown file in order without sampling**

Run: follow each Markdown file’s `前置条件 -> 操作步骤 -> 预期结果` sequence against `http://10.0.0.200:8080`.
Expected: every base-suite case reaches a concrete result of `通过`, `失败`, or `阻塞`.

- [ ] **Step 2: After each completed base case, update the shared execution tables**

Run: `python C:\zq\平台开发\OpsNexus\OpsNexus\scripts\update_execution_tables.py --demo-root C:\zq\平台开发\OpsNexus\demo --case-id <CASE_ID> --status <通过|失败|阻塞> --executor Codex --date 2026-03-24 --bug-id <BUG_OR_EMPTY> --remark "<facts only>"`
Expected: the matching row changes in `opsnexus-base-test-execution.csv`, `opsnexus-all-test-execution.csv`, and `opsnexus-test-execution.xlsx`.

- [ ] **Step 3: Preserve test data for destructive actions**

Run: only create/delete/edit rows that are explicitly test data and always verify by refresh or re-search.
Expected: no business data is damaged; destructive validations are limited to test records.

### Task 4: Execute all button-suite cases and record every result immediately

**Files:**
- Read: `C:\zq\平台开发\OpsNexus\demo\15-重点菜单按钮级测试用例\*.md`
- Modify: `C:\zq\平台开发\OpsNexus\demo\opsnexus-button-test-execution.csv`
- Modify: `C:\zq\平台开发\OpsNexus\demo\opsnexus-all-test-execution.csv`
- Modify: `C:\zq\平台开发\OpsNexus\demo\opsnexus-test-execution.xlsx`

- [ ] **Step 1: Execute every button-level case in file order**

Run: follow each button-case Markdown file exactly, including post-save refresh and post-delete re-search checks.
Expected: every button-suite case reaches a concrete result and has evidence in remarks/artifacts.

- [ ] **Step 2: Update the tables after every button-level case**

Run: `python C:\zq\平台开发\OpsNexus\OpsNexus\scripts\update_execution_tables.py --demo-root C:\zq\平台开发\OpsNexus\demo --case-id <CASE_ID> --status <通过|失败|阻塞> --executor Codex --date 2026-03-24 --bug-id <BUG_OR_EMPTY> --remark "<facts only>"`
Expected: the matching row changes in `opsnexus-button-test-execution.csv`, `opsnexus-all-test-execution.csv`, and `opsnexus-test-execution.xlsx`.

### Task 5: Defect loop for any failing case

**Files:**
- Modify: `C:\zq\平台开发\OpsNexus\OpsNexus\web\src\**\*`
- Modify: `C:\zq\平台开发\OpsNexus\OpsNexus\api\api\**\*`
- Modify: `C:\zq\平台开发\OpsNexus\demo\opsnexus-*.csv`
- Modify: `C:\zq\平台开发\OpsNexus\demo\opsnexus-test-execution.xlsx`

- [ ] **Step 1: Capture the failure before patching**

Run: rerun the failing route with browser artifacts enabled and save the exact UI/API symptom.
Expected: reproducible evidence exists before code changes.

- [ ] **Step 2: Use systematic debugging and minimal fixes**

Run: identify the smallest frontend/backend file set that explains the failure, patch only those files, and avoid reverting unrelated work.
Expected: the failing behavior is addressed without collateral regressions.

- [ ] **Step 3: Re-run the failing case and adjacent smoke coverage**

Run: re-execute the original case plus the smallest related smoke path.
Expected: original case now passes or is downgraded to a documented blocker with evidence.

- [ ] **Step 4: Sync the updated outcome back to the execution tables immediately**

Run: `python C:\zq\平台开发\OpsNexus\OpsNexus\scripts\update_execution_tables.py --demo-root C:\zq\平台开发\OpsNexus\demo --case-id <CASE_ID> --status <通过|失败|阻塞> --executor Codex --date 2026-03-24 --bug-id <BUG_OR_EMPTY> --remark "<fix + verification facts>"`
Expected: CSV and XLSX match the post-fix result.
