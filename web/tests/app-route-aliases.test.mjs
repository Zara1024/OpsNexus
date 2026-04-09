import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/router/app.js', import.meta.url), 'utf8')

function testApplicationRouteKeepsLegacyViewAlias() {
  assert.match(
    source,
    /path:\s*'\/app\/application'[\s\S]*alias:\s*\[\s*'\/app\/view'\s*\]/s,
    'route /app/application should keep /app/view as a legacy alias'
  )
}

async function main() {
  testApplicationRouteKeepsLegacyViewAlias()
  console.log('app route alias tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
