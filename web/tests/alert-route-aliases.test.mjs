import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/router/system.js', import.meta.url), 'utf8')

function testAlertCenterRouteKeepsLegacyIncidentAlias() {
  assert.match(
    source,
    /path:\s*'\/monitor\/alert-center'[\s\S]*alias:\s*\[\s*'\/monitor\/incident'\s*\]/s,
    'route /monitor/alert-center should keep /monitor/incident as a legacy alias'
  )
}

async function main() {
  testAlertCenterRouteKeepsLegacyIncidentAlias()
  console.log('alert route alias tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
