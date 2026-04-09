import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/cmdb/Host/CmdbHostTableCompact.vue', import.meta.url), 'utf8')

function testStatusColumnTemplateIsWellFormed() {
  assert.match(
    source,
    /<el-table-column min-width="100">\s*<template #header>\s*<span>状态<\/span>\s*<\/template>\s*<template v-slot="scope">\s*<div class="status-stack">/s,
    'expected the compact host table status column to render a well-formed header and status stack'
  )
}

function testCompactTableKeepsConfigStatusAndOperationOrder() {
  assert.match(
    source,
    /<div class="config-cell">[\s\S]*?<div class="status-stack">[\s\S]*?<div class="table-operation">/s,
    'expected the compact host table to keep config, status, and operation templates in order'
  )
}

function testCompactTableDefinesRuntimeHelpersForMetricAndStatusColumns() {
  assert.match(
    source,
    /methods:\s*{[\s\S]*\bgetUsageColor\s*\(/s,
    'expected the compact host table to define getUsageColor for the resource overview column'
  )

  assert.match(
    source,
    /methods:\s*{[\s\S]*\bgetStatusTagType\s*\(/s,
    'expected the compact host table to define getStatusTagType for the status chip column'
  )

  assert.match(
    source,
    /methods:\s*{[\s\S]*\bgetStatusText\s*\(/s,
    'expected the compact host table to define getStatusText for the status chip label'
  )
}

function testActionColumnKeepsRowActionsAlignedWithHeaderStart() {
  assert.match(
    source,
    /\.table-operation\s*{[\s\S]*justify-content:\s*flex-start\s*;/s,
    'expected the compact host table action row to start at the same horizontal origin as the 操作 header'
  )
}

function testActionColumnRemovesTheMoreDropdown() {
  assert.doesNotMatch(
    source,
    /table-operation__more|<el-dropdown[\s\S]*?<template #dropdown>[\s\S]*?<\/template>\s*<\/el-dropdown>/s,
    'expected the compact host table action column to remove the 更多 dropdown entirely'
  )
}

function testActionColumnExposesInlineDeleteAction() {
  assert.match(
    source,
    /<el-button type="danger" link v-authority="\['cmdb:ecs:delete'\]" @click="\$emit\('delete-host', scope\.row\)">[\s\S]*?删除主机[\s\S]*?<\/el-button>/s,
    'expected the compact host table action column to expose 删除主机 inline'
  )
}

async function main() {
  testStatusColumnTemplateIsWellFormed()
  testCompactTableKeepsConfigStatusAndOperationOrder()
  testCompactTableDefinesRuntimeHelpersForMetricAndStatusColumns()
  testActionColumnKeepsRowActionsAlignedWithHeaderStart()
  testActionColumnRemovesTheMoreDropdown()
  testActionColumnExposesInlineDeleteAction()
  console.log('cmdb host compact table layout tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})

