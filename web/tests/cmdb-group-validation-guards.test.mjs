import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/cmdb/cmdbGroup.vue', import.meta.url), 'utf8')

function testCreateDialogDoesNotExposeRootSentinelValue() {
  assert.doesNotMatch(
    source,
    /addGroupForm:\s*\{\s*parentId:\s*0\s*\}/s,
    'expected the create-group dialog to avoid exposing root sentinel value 0 directly in the UI'
  )

  assert.match(
    source,
    /const payload = \{\s*\.\.\.this\.addGroupForm,\s*parentId:\s*this\.normalizeParentId\(this\.addGroupForm\.parentId\)\s*\}/s,
    'expected create-group submission to normalize an empty parent selection back to root id 0'
  )
}

function testCreateGroupValidationIsGuardedBeforeNetworkErrorHandling() {
  assert.match(
    source,
    /isValidationFailure\s*\(error\)\s*\{[\s\S]*Array\.isArray/s,
    'expected the group page to define a validation-failure guard helper'
  )

  assert.match(
    source,
    /async addGroup\(\)\s*\{[\s\S]*catch \(error\) \{[\s\S]*if \(this\.isValidationFailure\(error\)\) \{\s*return\s*\}[\s\S]*this\.\$message\.error\('新增分组失败/s,
    'expected addGroup to exit quietly on form validation failures before showing network-level errors'
  )
}

function testEditGroupValidationIsGuardedBeforeRetryLoopTreatsItAsTransportFailure() {
  assert.match(
    source,
    /async editGroup\(\)\s*\{[\s\S]*catch \(error\) \{[\s\S]*if \(this\.isValidationFailure\(error\)\) \{\s*return\s*\}[\s\S]*this\.\$message\.error\(`修改分组失败:/s,
    'expected editGroup to stop on form validation failures instead of retrying and reporting transport errors'
  )
}

async function main() {
  testCreateDialogDoesNotExposeRootSentinelValue()
  testCreateGroupValidationIsGuardedBeforeNetworkErrorHandling()
  testEditGroupValidationIsGuardedBeforeRetryLoopTreatsItAsTransportFailure()
  console.log('cmdb group validation guard tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
