import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/cmdb/cmdbHost.vue', import.meta.url), 'utf8')

function testHostSyncDefinesBusinessStatusGuard() {
  assert.match(
    source,
    /isHostSyncAccepted\s*\(\s*responseData\s*\)\s*\{[\s\S]*responseData\?\.code[\s\S]*Number\(responseData\?\.data\?\.status[\s\S]*!==\s*3/s,
    'expected the host page to define a guard that rejects host sync responses whose business status is 3'
  )
}

function testBatchSyncUsesBusinessStatusGuard() {
  assert.match(
    source,
    /async handleBatchSync\(\)\s*\{[\s\S]*this\.isHostSyncAccepted\(item\.value\?\.data\)/s,
    'expected handleBatchSync to classify sync results with the host-sync business status guard'
  )
}

function testSingleHostSyncUsesBusinessStatusGuard() {
  assert.match(
    source,
    /async handleHostSync\(targetHost = null\)\s*\{[\s\S]*if\s*\(\s*this\.isHostSyncAccepted\(res\)\s*\)/s,
    'expected handleHostSync to verify the host-sync business status guard before showing success'
  )
}

async function main() {
  testHostSyncDefinesBusinessStatusGuard()
  testBatchSyncUsesBusinessStatusGuard()
  testSingleHostSyncUsesBusinessStatusGuard()
  console.log('cmdb host sync submit guard tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
