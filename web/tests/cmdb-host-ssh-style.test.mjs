import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/cmdb/Host/SSH.vue', import.meta.url), 'utf8')

function testSshPageUsesUnifiedPageShell() {
  assert.match(
    source,
    /class="ssh-page"/,
    'expected terminal login page to expose a unified page shell class'
  )
}

function testSshPageExposesTerminalPreviewSurface() {
  assert.match(
    source,
    /terminal-preview__screen/,
    'expected terminal login page to define a styled terminal preview surface'
  )
}

function testSshPageRemovesLegacyBackgroundOverrides() {
  assert.doesNotMatch(
    source,
    /setContainerBackground|#f5f5f5/,
    'expected terminal login page to remove the legacy light background override'
  )
}

async function main() {
  testSshPageUsesUnifiedPageShell()
  testSshPageExposesTerminalPreviewSurface()
  testSshPageRemovesLegacyBackgroundOverrides()
  console.log('cmdb host ssh style tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
