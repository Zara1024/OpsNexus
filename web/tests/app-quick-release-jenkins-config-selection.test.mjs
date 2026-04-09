import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/app/app_quick_release.vue', import.meta.url), 'utf8')

function testQuickReleaseSelectedAppsRetainJenkinsConfigs() {
  assert.match(
    source,
    /createForm\.applications\.push\(\{[\s\S]*jenkins_envs:\s*app\.jenkins_envs\s*\|\|\s*\[\][\s\S]*\}\)/s,
    "expected confirmAppSelection to retain each selected app's jenkins_envs so the parameter-config dialog can resolve Jenkins jobs"
  )
}

async function main() {
  testQuickReleaseSelectedAppsRetainJenkinsConfigs()
  console.log('app quick release jenkins config selection tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
