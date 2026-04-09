import assert from 'node:assert/strict'

import {
  extractCommandOutput,
  getCmdbDatabaseTypeIconKey,
  getCmdbDatabaseTypeLabel,
  getCmdbDatabaseTypeOptions,
  getCmdbDatabaseTypeTag
} from '../src/utils/cmdbPresentation.mjs'

const options = getCmdbDatabaseTypeOptions()
assert.equal(options.length, 8)
assert.equal(getCmdbDatabaseTypeLabel(6), '达梦数据库')
assert.equal(getCmdbDatabaseTypeLabel(7), 'GaussDB')
assert.equal(getCmdbDatabaseTypeLabel(8), 'Oracle')
assert.equal(getCmdbDatabaseTypeTag(8), 'danger')
assert.equal(getCmdbDatabaseTypeIconKey(7), 'gaussdb')

assert.equal(
  extractCommandOutput({
    output: {
      command: 'ls -a',
      output: '.\\n..\\n.profile\\n'
    }
  }),
  '.\n..\n.profile\n'
)

assert.equal(
  extractCommandOutput('line1\\nline2'),
  'line1\nline2'
)

assert.equal(
  extractCommandOutput('.\\n .Destination}}{{end}}\\n.profile\\n'),
  '.\n.profile\n'
)

console.log('cmdb presentation tests passed')
