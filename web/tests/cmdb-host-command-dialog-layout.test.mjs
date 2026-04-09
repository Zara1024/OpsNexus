import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/cmdb/cmdbHost.vue', import.meta.url), 'utf8')

function testCommandDialogUsesWideResponsiveWidth() {
  assert.match(
    source,
    /<el-dialog[\s\S]*title="[^"]*执行命令"[\s\S]*width="72vw"/,
    'expected the command dialog to use a wider responsive width'
  )
}

function testCommandOutputIsOutsideFormLayout() {
  assert.match(
    source,
    /<\/el-form>\s*<div v-if="commandDialog\.output" class="command-output"/,
    'expected the command output block to render outside the form layout'
  )
}

async function main() {
  testCommandDialogUsesWideResponsiveWidth()
  testCommandOutputIsOutsideFormLayout()
  console.log('cmdb host command dialog layout tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
