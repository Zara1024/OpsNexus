import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/task/TaskAnsible.vue', import.meta.url), 'utf8')

function testAnsibleSubmitGuardSkipsValidationFailures() {
  assert.match(
    source,
    /isValidationFailure\s*=\s*\(error\)\s*=>[\s\S]*Array\.isArray/s,
    'expected the Ansible task page to define a validation-failure guard helper'
  )

  assert.match(
    source,
    /const handleSubmit = async \(\) => \{[\s\S]*catch \(error\) \{[\s\S]*if \(isValidationFailure\(error\)\) \{\s*return\s*\}[\s\S]*ElMessage\.error\(`保存任务失败:/s,
    'expected handleSubmit to stop on validation failures before showing the generic save-failed toast'
  )
}

async function main() {
  testAnsibleSubmitGuardSkipsValidationFailures()
  console.log('task ansible validation guard tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
