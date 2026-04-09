import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/task/TaskAnsible.vue', import.meta.url), 'utf8')

function testHandleEditFallsBackToCurrentRowSnapshot() {
  assert.match(
    source,
    /const handleEdit = async row => \{[\s\S]*const data = res\.data\.data\?\.task_info \|\| \{\}[\s\S]*description:\s*data\.Description \|\| row\.description \|\| ''[\s\S]*hostGroups:\s*parseObjectMaybe\(data\.HostGroups \|\| row\.hostGroups\)[\s\S]*variables:\s*data\.GlobalVars \|\| row\.variables \|\| ''[\s\S]*is_recurring:\s*Number\(data\.IsRecurring \?\? row\.is_recurring \?\? 0\)[\s\S]*cron_expr:\s*data\.CronExpr \|\| row\.cron_expr \|\| ''/s,
    'expected handleEdit to fall back to the current table row snapshot when the ansible detail endpoint only returns a minimal task_info payload'
  )
}

async function main() {
  testHandleEditFallsBackToCurrentRowSnapshot()
  console.log('task ansible edit row fallback tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
