import assert from 'node:assert/strict'
import fs from 'node:fs'

const dbLogSource = fs.readFileSync(new URL('../src/views/monitor/DBLog.vue', import.meta.url), 'utf8')
const cmdbApiSource = fs.readFileSync(new URL('../src/api/cmdb.js', import.meta.url), 'utf8')

function testCmdbApiExposesBatchDeleteSqlLog() {
  assert.match(
    cmdbApiSource,
    /BatchDeleteCmdbSqlLog\(ids\)\s*{[\s\S]*url:\s*'cmdb\/sqlLog\/batch\/delete'[\s\S]*method:\s*'delete'/,
    'expected cmdb api to expose a dedicated batch delete SQL log request'
  )
}

function testDbLogBatchDeleteUsesBatchApi() {
  assert.match(
    dbLogSource,
    /await this\.\$api\.BatchDeleteCmdbSqlLog\(sqlLogIds\)/,
    'expected DBLog batchHandleDelete to call the dedicated batch delete API'
  )

  assert.doesNotMatch(
    dbLogSource,
    /await this\.\$api\.DeleteCmdbSqlLogById\(sqlLogIds\)/,
    'batchHandleDelete should not reuse the single-delete API with an ids array'
  )
}

async function main() {
  testCmdbApiExposesBatchDeleteSqlLog()
  testDbLogBatchDeleteUsesBatchApi()
  console.log('monitor dblog batch delete api tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
