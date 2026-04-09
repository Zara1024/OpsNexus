import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/monitor/https.vue', import.meta.url), 'utf8')

function testSslDomainTableExposesDomainStatusAndRowActions() {
  const sslDomainTableMatch = source.match(
    /<el-table\s+:data="domains"[\s\S]*?<\/el-table>/
  )

  assert.ok(sslDomainTableMatch, 'expected the SSL domain table to exist in monitor automation page')

  const sslDomainTable = sslDomainTableMatch[0]

  assert.match(
    sslDomainTable,
    /<el-table-column label="状态"[\s\S]*toggleDomain\(row\)/s,
    'expected the SSL domain table to render a status column that can toggle a domain'
  )

  assert.match(
    sslDomainTable,
    /<el-table-column label="操作"[\s\S]*openDomainDialog\(row\)[\s\S]*deleteDomain\(row\)/s,
    'expected the SSL domain table to render edit and delete actions for each domain row'
  )
}

async function main() {
  testSslDomainTableExposesDomainStatusAndRowActions()
  console.log('monitor automation SSL domain action tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
