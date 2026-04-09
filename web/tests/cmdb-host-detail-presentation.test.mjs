import assert from 'node:assert/strict'

import {
  formatHostResourceValue,
  resolveHostSyncTarget
} from '../src/utils/cmdbHostDetailPresentation.mjs'

function testResolveHostSyncTargetFallsBackFromClickEvent() {
  const eventLike = {
    type: 'click',
    target: {}
  }
  const hostDetail = {
    id: 718,
    hostName: 'opsnexus-local-verify'
  }

  assert.deepEqual(resolveHostSyncTarget(eventLike, hostDetail), hostDetail)
}

function testResolveHostSyncTargetPrefersExplicitHost() {
  const row = {
    id: 719,
    hostName: 'linux-b'
  }
  const hostDetail = {
    id: 718,
    hostName: 'linux-a'
  }

  assert.deepEqual(resolveHostSyncTarget(row, hostDetail), row)
}

function testFormatHostResourceValueAppendsUnitToNumericValues() {
  assert.equal(formatHostResourceValue('16', '核'), '16核')
  assert.equal(formatHostResourceValue('8', 'G'), '8G')
  assert.equal(formatHostResourceValue('97', 'GB'), '97GB')
}

function testFormatHostResourceValueDoesNotExposeBareUnits() {
  assert.equal(formatHostResourceValue('', 'G'), '-')
  assert.equal(formatHostResourceValue(null, 'GB'), '-')
}

function testFormatHostResourceValueKeepsExistingUnits() {
  assert.equal(formatHostResourceValue('8G', 'G'), '8G')
  assert.equal(formatHostResourceValue('97GB', 'GB'), '97GB')
}

async function main() {
  testResolveHostSyncTargetFallsBackFromClickEvent()
  testResolveHostSyncTargetPrefersExplicitHost()
  testFormatHostResourceValueAppendsUnitToNumericValues()
  testFormatHostResourceValueDoesNotExposeBareUnits()
  testFormatHostResourceValueKeepsExistingUnits()
  console.log('cmdb host detail presentation tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
