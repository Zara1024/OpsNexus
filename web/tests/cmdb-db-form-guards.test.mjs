import assert from 'node:assert/strict'
import fs from 'node:fs'

import { createDatabaseAssetFormModel } from '../src/utils/cmdbAssetPresentation.mjs'

const source = fs.readFileSync(new URL('../src/views/cmdb/cmdbDB.vue', import.meta.url), 'utf8')

function testDatabaseAssetFormModelStartsWithReadablePlatform() {
  const model = createDatabaseAssetFormModel()
  assert.equal(model.type, 1, 'expected the default database asset type to stay MySQL')
  assert.equal(model.platform, 'MySQL', 'expected the default database asset platform label to be prefilled for the default type')
}

function testDatabaseSubmitGuardSkipsValidationFailures() {
  assert.match(
    source,
    /isValidationFailure\s*\(error\)\s*\{[\s\S]*Array\.isArray/s,
    'expected the database asset page to define a validation-failure guard helper'
  )

  assert.match(
    source,
    /async submitForm\(\)\s*\{[\s\S]*catch \(error\) \{[\s\S]*if \(this\.isValidationFailure\(error\)\) \{\s*return\s*\}[\s\S]*this\.\$message\.error\(`数据库资产保存失败/s,
    'expected submitForm to exit quietly on validation failures before showing a save-failed toast'
  )
}

async function main() {
  testDatabaseAssetFormModelStartsWithReadablePlatform()
  testDatabaseSubmitGuardSkipsValidationFailures()
  console.log('cmdb db form guard tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
