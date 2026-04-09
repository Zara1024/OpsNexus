import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/configcenter/accountauth.vue', import.meta.url), 'utf8')

function testFilterOptionsMatchStoredAccountTypes() {
  const filterBlockMatch = source.match(/<el-select v-model="queryParams\.type"[\s\S]*?<\/el-select>/)
  assert.ok(filterBlockMatch, 'expected to find the toolbar account type filter select')

  const filterBlock = filterBlockMatch[0]

  assert.match(
    filterBlock,
    /<el-option label="Zabbix" :value="5" \/>/,
    'expected toolbar filter dropdown to expose Zabbix accounts as type 5'
  )

  assert.match(
    filterBlock,
    /<el-option label="通用账号" :value="6" \/>/,
    'expected toolbar filter dropdown to expose generic accounts as type 6'
  )
}

async function main() {
  testFilterOptionsMatchStoredAccountTypes()
  console.log('config accountauth filter option tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
