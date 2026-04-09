import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/cmdb/DBdetails.vue', import.meta.url), 'utf8')

function testDbDetailsDefinesTypeAwareDatabaseModeHelpers() {
  assert.match(
    source,
    /isPostgreSQLDatabaseType|isPostgresDatabaseType/,
    'expected DBdetails to define a PostgreSQL-aware database mode helper'
  )

  assert.match(
    source,
    /isRedisDatabaseType/,
    'expected DBdetails to define a Redis-aware database mode helper'
  )
}

function testDbDetailsDefinesRedisCommandModeSupport() {
  assert.match(
    source,
    /redisCommand|commandMode|redisMode/,
    'expected DBdetails to define Redis command-mode support'
  )
}

function testDbDetailsIncludesPostgreSqlMetadataQueryGuidance() {
  assert.match(
    source,
    /SELECT datname FROM pg_database|SELECT datname AS database_name FROM pg_database/,
    'expected DBdetails to include a PostgreSQL metadata query example instead of relying on MySQL-style SHOW DATABASE syntax'
  )
}

function testDbDetailsDoesNotKeepMysqlOnlyFrontendGuard() {
  assert.doesNotMatch(
    source,
    /only MySQL databases are currently supported/,
    'expected DBdetails to stop surfacing the MySQL-only support guard in the page logic'
  )
}

async function main() {
  testDbDetailsDefinesTypeAwareDatabaseModeHelpers()
  testDbDetailsDefinesRedisCommandModeSupport()
  testDbDetailsIncludesPostgreSqlMetadataQueryGuidance()
  testDbDetailsDoesNotKeepMysqlOnlyFrontendGuard()
  console.log('cmdb DBdetails type support tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
