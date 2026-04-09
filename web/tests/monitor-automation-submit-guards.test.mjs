import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('../src/views/monitor/https.vue', import.meta.url), 'utf8')

function testMonitorAutomationDefinesBusinessResponseGuard() {
  assert.match(
    source,
    /ensureApiSuccess\s*\(\s*response\s*,\s*fallbackMessage\s*\)\s*\{[\s\S]*response\?\.data\?\.code[\s\S]*throw new Error/s,
    'expected the monitor automation page to define a response-code guard for API calls'
  )
}

function testHostRuleSubmitValidatesRequiredFieldsAndChecksResponseCode() {
  assert.match(
    source,
    /validateHostRuleForm\s*\(\s*payload\s*\)\s*\{[\s\S]*payload\.name[\s\S]*payload\.hostId[\s\S]*payload\.metricKey[\s\S]*payload\.severity/s,
    'expected host-rule submission to validate required business fields before calling the API'
  )

  assert.match(
    source,
    /async submitHostRule\(\)\s*\{[\s\S]*this\.validateHostRuleForm\(payload\)[\s\S]*this\.ensureApiSuccess\(/s,
    'expected submitHostRule to validate the form and then verify the business response code'
  )
}

function testDatabaseRuleSubmitAlsoChecksResponseCode() {
  assert.match(
    source,
    /async submitDBRule\(\)\s*\{[\s\S]*this\.ensureApiSuccess\(/s,
    'expected submitDBRule to verify the business response code before showing success'
  )
}

async function main() {
  testMonitorAutomationDefinesBusinessResponseGuard()
  testHostRuleSubmitValidatesRequiredFieldsAndChecksResponseCode()
  testDatabaseRuleSubmitAlsoChecksResponseCode()
  console.log('monitor automation submit guard tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
