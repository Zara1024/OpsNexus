import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/monitor/LoginLog.vue', import.meta.url), 'utf8')

function testBatchDeleteDisabledStateExplainsSelectionRequirement() {
  assert.match(
    source,
    /<el-tooltip[\s\S]*content="请先勾选要删除的登录日志"[\s\S]*<el-button[\s\S]*批量删除/s,
    'expected the disabled batch-delete button to explain that users must select login logs first'
  )
}

async function main() {
  testBatchDeleteDisabledStateExplainsSelectionRequirement()
  console.log('monitor loginlog usability tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
