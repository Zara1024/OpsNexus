import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/configcenter/KeyManage.vue', import.meta.url), 'utf8')

function testKeyDeleteCancelDoesNotLogConsoleError() {
  const start = source.lastIndexOf('async handleDelete(row)')
  const end = source.indexOf('\n  }\n}', start)
  assert.notEqual(start, -1, 'expected KeyManage.vue to define handleDelete')
  assert.notEqual(end, -1, 'expected KeyManage.vue to contain the next method after handleDelete')
  const methodSource = source.slice(start, end)

  assert.match(
    methodSource,
    /if\s*\(\s*error\s*===\s*'cancel'\s*\|\|\s*error\s*===\s*'close'\s*\)\s*return/,
    'expected handleDelete to silently ignore Element Plus confirm cancellations'
  )
}

async function main() {
  testKeyDeleteCancelDoesNotLogConsoleError()
  console.log('config keymanage delete cancel silence tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
