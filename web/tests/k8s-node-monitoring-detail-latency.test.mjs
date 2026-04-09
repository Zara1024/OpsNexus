import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/K8s/nodes/NodesMonitoring.vue', import.meta.url), 'utf8')

function testDetailDialogOpensBeforeAwaitingRemoteFetches() {
  const match = source.match(/const viewNodeDetail = async \(node\) => \{([\s\S]*?)\n\}/)
  assert.ok(match, 'expected NodesMonitoring.vue to define viewNodeDetail')

  const body = match[1]
  const openIndex = body.indexOf('nodeDetailVisible.value = true')
  const awaitIndex = body.indexOf('await ')

  assert.notEqual(openIndex, -1, 'expected viewNodeDetail to set nodeDetailVisible.value = true')
  assert.notEqual(awaitIndex, -1, 'expected viewNodeDetail to await detail loading work')
  assert.ok(
    openIndex < awaitIndex,
    'expected the detail dialog to open before awaiting remote metric/detail fetches'
  )
}

function testDetailDialogTracksItsOwnLoadingState() {
  assert.match(
    source,
    /const nodeDetailLoading = ref\(false\)/,
    'expected NodesMonitoring.vue to define dedicated loading state for the detail dialog'
  )
}

async function main() {
  testDetailDialogOpensBeforeAwaitingRemoteFetches()
  testDetailDialogTracksItsOwnLoadingState()
  console.log('k8s node monitoring detail latency tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
