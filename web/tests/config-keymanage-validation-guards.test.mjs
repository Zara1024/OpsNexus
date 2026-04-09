import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/configcenter/KeyManage.vue', import.meta.url), 'utf8')

function testSubmitFormStopsBeforeCatchWhenValidationFails() {
  assert.match(
    source,
    /const\s+valid\s*=\s*await\s+this\.\$refs\.formRef\.validate\(\)\.catch\(\(\)\s*=>\s*false\)/,
    'expected submitForm to normalize key form validation failures into a false result'
  )

  assert.match(
    source,
    /if\s*\(!valid\)\s*return/,
    'expected submitForm to return early when the key form is invalid'
  )
}

async function main() {
  testSubmitFormStopsBeforeCatchWhenValidationFails()
  console.log('config keymanage validation guard tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
