import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/app/app_quick_release.vue', import.meta.url), 'utf8')

function testQuickReleaseSubmitGuardsValidationFailures() {
  assert.match(
    source,
    /isValidationFailure\s*=\s*\(error\)\s*=>[\s\S]*Array\.isArray/s,
    'expected the quick-release page to define a validation-failure guard helper'
  )

  assert.match(
    source,
    /const handleCreateSubmit = async \(\) => \{[\s\S]*catch \(error\) \{[\s\S]*if \(isValidationFailure\(error\)\) \{\s*return\s*\}[\s\S]*ElMessage\.error\('创建发布失败'\)/s,
    'expected handleCreateSubmit to stop on validation failures before showing the generic create-failed toast'
  )
}

async function main() {
  testQuickReleaseSubmitGuardsValidationFailures()
  console.log('app quick release validation guard tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
