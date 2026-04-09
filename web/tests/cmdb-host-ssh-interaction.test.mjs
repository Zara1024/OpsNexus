import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/cmdb/Host/SSH.vue', import.meta.url), 'utf8')

function testHostTreeBindsNativeDoubleClickAction() {
  assert.match(
    source,
    /@dblclick\.stop="handleTreeNodeDblClick\(data\)"/,
    'expected host tree nodes to bind a native double-click action'
  )
}

function testDoubleClickHandlerOpensTerminalForHostNodes() {
  assert.match(
    source,
    /handleTreeNodeDblClick\(data\)\s*\{[\s\S]*if\s*\(!data\?\.isHost\)\s*\{\s*return[\s\S]*this\.openTerminal\(\)/,
    'expected double-click handler to open terminal only for host nodes'
  )
}

async function main() {
  testHostTreeBindsNativeDoubleClickAction()
  testDoubleClickHandlerOpensTerminalForHostNodes()
  console.log('cmdb host ssh interaction tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
