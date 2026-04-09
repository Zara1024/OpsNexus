import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/api/app.js', import.meta.url), 'utf8')

function testGetJenkinsJobParametersUsesExtendedTimeout() {
  assert.match(
    source,
    /getJenkinsJobParameters\(serverId,\s*jobName\)\s*\{[\s\S]*url:\s*`\/jenkins\/\$\{serverId\}\/jobs\/\$\{jobName\}\/parameters`[\s\S]*timeout:\s*3[0-9]{4}[\s\S]*\}/s,
    'expected getJenkinsJobParameters to opt into a timeout longer than the shared 15s axios default so Jenkins parameter requests can surface backend responses instead of browser aborts'
  )
}

async function main() {
  testGetJenkinsJobParametersUsesExtendedTimeout()
  console.log('app api jenkins params timeout tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
