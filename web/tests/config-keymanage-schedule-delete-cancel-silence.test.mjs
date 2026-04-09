import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/configcenter/KeyManage.vue', import.meta.url), 'utf8')

function testScheduleDeleteCancelDoesNotLogConsoleError() {
  assert.match(
    source,
    /handleDeleteSchedule[\s\S]*if\s*\(\s*error\s*===\s*'cancel'\s*\|\|\s*error\s*===\s*'close'\s*\)\s*return/,
    'expected handleDeleteSchedule to silently ignore Element Plus confirm cancellations'
  )
}

async function main() {
  testScheduleDeleteCancelDoesNotLogConsoleError()
  console.log('config keymanage schedule delete cancel silence tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
