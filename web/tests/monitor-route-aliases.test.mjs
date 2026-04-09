import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/router/system.js', import.meta.url), 'utf8')

function testMonitorWorkbenchRouteKeepsLegacyDomainAlias() {
  assert.match(
    source,
    /path:\s*'\/monitor\/https'[\s\S]*alias:\s*\[\s*'\/monitor\/domain'\s*\]/s,
    'route /monitor/https should keep /monitor/domain as a legacy alias'
  )
}

async function main() {
  testMonitorWorkbenchRouteKeepsLegacyDomainAlias()
  console.log('monitor route alias tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
