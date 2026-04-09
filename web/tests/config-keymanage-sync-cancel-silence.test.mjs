import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/configcenter/KeyManage.vue', import.meta.url), 'utf8')

function testSyncCancelDoesNotLogConsoleErrorOrToast() {
  const start = source.lastIndexOf('async handleSyncHosts(row)')
  const end = source.indexOf('\n  }\n}', start)
  assert.notEqual(start, -1, 'expected KeyManage.vue to define handleSyncHosts')
  assert.notEqual(end, -1, 'expected KeyManage.vue to contain the next method after handleSyncHosts')
  const methodSource = source.slice(start, end)
  assert.match(
    methodSource,
    /if\s*\(\s*error\s*===\s*'cancel'\s*\|\|\s*error\s*===\s*'close'\s*\)\s*return/,
    'expected handleSyncHosts to silently ignore Element Plus confirm cancellations'
  )
}

async function main() {
  testSyncCancelDoesNotLogConsoleErrorOrToast()
  console.log('config keymanage sync cancel silence tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
