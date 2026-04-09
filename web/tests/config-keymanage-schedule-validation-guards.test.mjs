import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/configcenter/KeyManage.vue', import.meta.url), 'utf8')

function testSubmitScheduleFormStopsBeforeCatchWhenValidationFails() {
  assert.match(
    source,
    /const\s+valid\s*=\s*await\s+this\.\$refs\.scheduleFormRef\.validate\(\)\.catch\(\(\)\s*=>\s*false\)/,
    'expected submitScheduleForm to normalize schedule form validation failures into a false result'
  )

  assert.match(
    source,
    /if\s*\(!valid\)\s*return/,
    'expected submitScheduleForm to return early when the schedule form is invalid'
  )
}

async function main() {
  testSubmitScheduleFormStopsBeforeCatchWhenValidationFails()
  console.log('config keymanage schedule validation guard tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
