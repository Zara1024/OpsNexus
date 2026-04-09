import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/app/app_quick_release.vue', import.meta.url), 'utf8')

function testQuickReleaseDescriptionTextareaUsesNumericRowsProp() {
  assert.doesNotMatch(
    source,
    /<el-input[\s\S]*type="textarea"[\s\S]*(?<!:)rows="3"/s,
    'expected the quick-release description textarea to avoid passing rows as a string prop'
  )

  assert.match(
    source,
    /<el-input[\s\S]*type="textarea"[\s\S]*:rows="3"/s,
    'expected the quick-release description textarea to bind rows as a numeric prop'
  )
}

async function main() {
  testQuickReleaseDescriptionTextareaUsesNumericRowsProp()
  console.log('app quick release dialog prop tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
