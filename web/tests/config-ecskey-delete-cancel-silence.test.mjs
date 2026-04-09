import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/configcenter/ecs-key.vue', import.meta.url), 'utf8')

function testDeleteCancelDoesNotLogConsoleError() {
  assert.match(
    source,
    /if\s*\(\s*error\s*===\s*'cancel'\s*\|\|\s*error\s*===\s*'close'\s*\)\s*return/,
    'expected handleDelete to silently ignore Element Plus confirm cancellations'
  )
}

async function main() {
  testDeleteCancelDoesNotLogConsoleError()
  console.log('config ecskey delete cancel silence tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
