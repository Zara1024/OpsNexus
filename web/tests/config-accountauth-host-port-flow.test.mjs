import assert from 'node:assert/strict'
import fs from 'node:fs'

import {
  formatAccountHostPortInput,
  parseAccountHostPortInput
} from '../src/utils/accountAuthEndpoint.mjs'

const source = fs.readFileSync(new URL('../src/views/configcenter/accountauth.vue', import.meta.url), 'utf8')

function testParseAccountHostPortInput() {
  assert.deepEqual(
    parseAccountHostPortInput('10.0.0.1:18080'),
    { host: '10.0.0.1', port: 18080, valid: true },
    'expected host:port input to split into host and numeric port'
  )

  assert.deepEqual(
    parseAccountHostPortInput('jenkins.example.com'),
    { host: 'jenkins.example.com', port: null, valid: false },
    'expected host without port to be rejected'
  )
}

function testFormatAccountHostPortInput() {
  assert.equal(
    formatAccountHostPortInput('10.0.0.1', 18080),
    '10.0.0.1:18080',
    'expected edit dialog to rehydrate stored host and port back into a single input value'
  )
}

function testAccountAuthViewUsesParsedPortForCreateAndEdit() {
  assert.match(
    source,
    /formatAccountHostPortInput\(row\.host,\s*row\.port\)/,
    'expected edit dialog to include the stored port in the host input'
  )

  assert.match(
    source,
    /const\s+parsedEndpoint\s*=\s*parseAccountHostPortInput\(this\.formData\.host\)/,
    'expected submitForm to parse the combined host input before sending the request'
  )

  assert.match(
    source,
    /port:\s*parsedEndpoint\.port/,
    'expected create\/update requests to include the parsed port'
  )

  assert.match(
    source,
    /API\.updateAccountAuth\(\{\s*id:\s*updateData\.id,[\s\S]*port:\s*updateData\.port,[\s\S]*\}\)/,
    'expected updateAccountAuth payload to include the parsed port'
  )
}

async function main() {
  testParseAccountHostPortInput()
  testFormatAccountHostPortInput()
  testAccountAuthViewUsesParsedPortForCreateAndEdit()
  console.log('config accountauth host-port flow tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
