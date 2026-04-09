import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/cmdb/Host/SSH.vue', import.meta.url), 'utf8')

function testSearchWarningMatchesPlaceholderScope() {
  assert.match(
    source,
    /placeholder="搜索分组或主机"/,
    'expected the SSH page search placeholder to mention both groups and hosts'
  )

  assert.match(
    source,
    /this\.\$message\.warning\('未找到匹配的分组或主机'\)/,
    'expected the no-match warning to mention both groups and hosts'
  )
}

function testProductionTreeSelectionAvoidsDebugConsoleLogs() {
  assert.doesNotMatch(
    source,
    /console\.log\('点击事件触发测试'\)|console\.log\('点击的是主机节点'\)/,
    'expected the SSH page to remove debug click logs from production code'
  )
}

async function main() {
  testSearchWarningMatchesPlaceholderScope()
  testProductionTreeSelectionAvoidsDebugConsoleLogs()
  console.log('cmdb ssh usability tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
