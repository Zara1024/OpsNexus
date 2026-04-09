import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/cmdb/cmdbHost.vue', import.meta.url), 'utf8')

function testToolbarReplacesTerminalWithDeleteButton() {
  assert.match(
    source,
    /<el-tooltip[\s\S]*content="请先勾选需要删除的主机"[\s\S]*?v-authority="\['cmdb:ecs:delete'\]"[\s\S]*?@click="handleBatchHostDelete"[\s\S]*?:disabled="!hasSelectedHosts"[\s\S]*?class="delete-host-btn"/s,
    'expected the host toolbar to replace the terminal entry with a batch delete button'
  )

  assert.doesNotMatch(
    source,
    /class="terminal-btn"[\s\S]*?@click="handleHostSSH"/s,
    'expected the old terminal toolbar button to be removed'
  )
}

function testHostPageDefinesBatchDeleteHandler() {
  assert.match(
    source,
    /async handleBatchHostDelete\(\)\s*\{[\s\S]*?this\.selectedHosts[\s\S]*?this\.\$api\.deleteCmdbHost\(host\.id\)/s,
    'expected the host page to define a batch delete handler that reuses deleteCmdbHost for selected hosts'
  )
}

async function main() {
  testToolbarReplacesTerminalWithDeleteButton()
  testHostPageDefinesBatchDeleteHandler()
  console.log('cmdb host delete action tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})

