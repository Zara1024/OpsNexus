import assert from 'node:assert/strict'
import fs from 'node:fs'

const hostPageSource = fs.readFileSync(new URL('../src/views/cmdb/cmdbHost.vue', import.meta.url), 'utf8')
const compactTableSource = fs.readFileSync(new URL('../src/views/cmdb/Host/CmdbHostTableCompact.vue', import.meta.url), 'utf8')
const createHostSource = fs.readFileSync(new URL('../src/views/cmdb/Host/CreateHost.vue', import.meta.url), 'utf8')

function testStatusFilterLabelsStayReadable() {
  assert.match(
    hostPageSource,
    /statusList:\s*\[[\s\S]*\{\s*value:\s*2,\s*label:\s*'未认证'\s*\}/s,
    'expected the host status filter to show 未认证 instead of placeholder text'
  )

  assert.doesNotMatch(
    hostPageSource,
    /statusList:\s*\[[\s\S]*label:\s*'\?\?\?'/s,
    'expected the host status filter to avoid placeholder ??? labels'
  )
}

function testBatchDeleteExplainsWhyItIsDisabled() {
  assert.match(
    hostPageSource,
    /<el-tooltip[\s\S]*content="请先勾选需要删除的主机"[\s\S]*<el-button[\s\S]*class="delete-host-btn"/s,
    'expected the disabled batch delete button to explain the missing selection prerequisite'
  )
}

function testCompactTablePreservesSelectionAcrossMonitorRefresh() {
  assert.match(
    compactTableSource,
    /<el-table[\s\S]*row-key="id"/s,
    'expected the compact host table to define row-key="id" so selection survives monitor refreshes'
  )

  assert.match(
    compactTableSource,
    /<el-table-column\s+type="selection"[\s\S]*:reserve-selection="true"/s,
    'expected the compact host table selection column to reserve selection when rows are refreshed'
  )
}

function testCreateHostValidationDoesNotPolluteConsoleErrors() {
  assert.doesNotMatch(
    createHostSource,
    /console\.error\('表单验证失败:'/,
    'expected CreateHost validation failures to avoid console.error noise'
  )
}

async function main() {
  testStatusFilterLabelsStayReadable()
  testBatchDeleteExplainsWhyItIsDisabled()
  testCompactTablePreservesSelectionAcrossMonitorRefresh()
  testCreateHostValidationDoesNotPolluteConsoleErrors()
  console.log('cmdb host ecs usability regression tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
