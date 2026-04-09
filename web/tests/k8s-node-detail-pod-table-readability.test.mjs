import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/K8s/nodes/NodesMonitoring.vue', import.meta.url), 'utf8')

function testPodTableUsesDedicatedDarkSurfaceAndReadableText() {
  assert.match(
    source,
    /\.pod-table\s*\{[\s\S]*background:\s*rgba\(15,\s*23,\s*42,\s*0\.92\)/s,
    'expected the pod table container to use a darker dedicated surface'
  )

  assert.match(
    source,
    /\.pod-header\s*\{[\s\S]*color:\s*#f8fafc/s,
    'expected the pod table header text to use a bright readable color'
  )

  assert.match(
    source,
    /\.pod-col\s*\{[\s\S]*color:\s*#e2e8f0/s,
    'expected the pod table body text to use a readable foreground color'
  )
}

function testPodRowsUseClearerSeparationAndHoverState() {
  assert.match(
    source,
    /\.pod-row\s*\{[\s\S]*background:\s*rgba\(15,\s*23,\s*42,\s*0\.72\)/s,
    'expected each pod row to use a darker contrasting row surface'
  )

  assert.match(
    source,
    /\.pod-row:hover\s*\{[\s\S]*background:\s*rgba\(30,\s*41,\s*59,\s*0\.92\)/s,
    'expected the pod row hover state to visibly stand out'
  )
}

async function main() {
  testPodTableUsesDedicatedDarkSurfaceAndReadableText()
  testPodRowsUseClearerSeparationAndHoverState()
  console.log('k8s node detail pod table readability tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
