import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/router/config.js', import.meta.url), 'utf8')

function testAccountAuthRouteKeepsLegacySystemConfigAlias() {
  assert.match(
    source,
    /path:\s*'\/config\/accountauth'[\s\S]*alias:\s*\[\s*'\/system\/config'\s*\]/s,
    'route /config/accountauth should keep /system/config as a legacy alias'
  )
}

async function main() {
  testAccountAuthRouteKeepsLegacySystemConfigAlias()
  console.log('config route alias tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
