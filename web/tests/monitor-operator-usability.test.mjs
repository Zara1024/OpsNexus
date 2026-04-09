import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/monitor/Operator.vue', import.meta.url), 'utf8')

function testBatchDeleteDisabledStateExplainsSelectionRequirement() {
  assert.match(
    source,
    /<el-tooltip[\s\S]*content="请先勾选要删除的操作日志"[\s\S]*<el-button[\s\S]*批量删除/s,
    'expected the disabled batch-delete button to explain that users must select operation logs first'
  )
}

async function main() {
  testBatchDeleteDisabledStateExplainsSelectionRequirement()
  console.log('monitor operator usability tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
