import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/configcenter/ecs-key.vue', import.meta.url), 'utf8')

function testCreateDialogClearsStaleValidationState() {
  assert.match(
    source,
    /showAddDialog\s*\(\)\s*{[\s\S]*this\.dialogVisible\s*=\s*true[\s\S]*this\.\$nextTick\(\(\)\s*=>\s*{[\s\S]*this\.\$refs\.formRef\??\.\s*clearValidate\(\)/,
    'expected showAddDialog to clear stale Element Plus validation state after reopening the dialog'
  )
}

function testEditDialogClearsStaleValidationState() {
  assert.match(
    source,
    /showEditDialog\s*\(row\)\s*{[\s\S]*this\.dialogVisible\s*=\s*true[\s\S]*this\.\$nextTick\(\(\)\s*=>\s*{[\s\S]*this\.\$refs\.formRef\??\.\s*clearValidate\(\)/,
    'expected showEditDialog to clear stale Element Plus validation state after reopening the dialog'
  )
}

function testDialogDestroysStaleFormStateOnClose() {
  assert.match(
    source,
    /<el-dialog[^>]*destroy-on-close/,
    'expected the ECS credential dialog to destroy its inner form state on close so stale validation does not survive reopen'
  )
}

async function main() {
  testCreateDialogClearsStaleValidationState()
  testEditDialogClearsStaleValidationState()
  testDialogDestroysStaleFormStateOnClose()
  console.log('config ecskey dialog validation reset tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
