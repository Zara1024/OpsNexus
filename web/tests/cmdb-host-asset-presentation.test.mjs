import assert from 'node:assert/strict'

import {
  buildAssetConnectivityTag,
  createHostAssetFormModel,
  mapHostRowToAssetRow
} from '../src/utils/cmdbAssetPresentation.mjs'

function testMapHostRowToAssetRowUsesPrimaryHostFields() {
  const row = mapHostRowToAssetRow({
    hostName: 'ops-prod-01',
    sshIp: '10.20.30.40',
    sshName: 'root',
    os: 'Ubuntu 22.04 LTS',
    status: 1
  })

  assert.equal(row.name, 'ops-prod-01')
  assert.equal(row.address, '10.20.30.40')
  assert.equal(row.account, 'root')
  assert.equal(row.platform, 'Ubuntu 22.04 LTS')
  assert.deepEqual(row.connectivity, buildAssetConnectivityTag({ reachable: true }))
}

function testMapHostRowToAssetRowSupportsFallbackFields() {
  const row = mapHostRowToAssetRow({
    name: 'fallback-name',
    privateIp: '172.16.1.9',
    sshName: '',
    deviceType: 'windows',
    status: 3
  })

  assert.equal(row.name, 'fallback-name')
  assert.equal(row.address, '172.16.1.9')
  assert.equal(row.account, '')
  assert.equal(row.platform, 'windows')
  assert.deepEqual(row.connectivity, buildAssetConnectivityTag({ reachable: false }))
}

function testMapHostRowToAssetRowFallsBackToStatusWhenExplicitSignalMissing() {
  const row = mapHostRowToAssetRow({
    hostName: 'fallback-status',
    status: 1,
    isAlive: undefined,
    onlineStatus: undefined
  })

  assert.deepEqual(
    row.connectivity,
    buildAssetConnectivityTag({ reachable: true }),
    'status fallback should still work when explicit monitor signal is missing'
  )
}

function testMapHostRowToAssetRowPrefersExplicitConnectivitySignals() {
  const onlineByAlive = mapHostRowToAssetRow({
    hostName: 'monitor-source',
    status: 3,
    isAlive: true
  })
  assert.deepEqual(
    onlineByAlive.connectivity,
    buildAssetConnectivityTag({ reachable: true }),
    'explicit alive signal should override status fallback'
  )

  const offlineByAlive = mapHostRowToAssetRow({
    hostName: 'monitor-source-2',
    status: 1,
    isAlive: false
  })
  assert.deepEqual(
    offlineByAlive.connectivity,
    buildAssetConnectivityTag({ reachable: false }),
    'explicit offline signal should override status fallback'
  )
}

function testMapHostRowToAssetRowPreservesCpuAndMemoryShape() {
  const row = mapHostRowToAssetRow({
    hostName: 'shape-host',
    cpu: 16,
    memory: 64,
    status: 1
  })

  assert.equal(row.cpu, 16)
  assert.equal(row.memory, 64)
}

function testCreateHostAssetFormModelDefaults() {
  const form = createHostAssetFormModel()

  assert.equal(form.hostName, '')
  assert.equal(form.ip, '')
  assert.equal(form.port, 22)
  assert.equal(form.username, '')
  assert.equal(form.authId, '')
  assert.equal(form.groupId, '')
  assert.equal(form.remark, '')
  assert.equal(form.deviceType, 'linux')
  assert.equal(form.remoteDomain, '')
  assert.equal(form.rdpPort, 3389)
  assert.equal(form.rdpUsername, '')
  assert.equal(form.rdpPassword, '')
  assert.equal(form.rdpDomain, '')
}

async function main() {
  testMapHostRowToAssetRowUsesPrimaryHostFields()
  testMapHostRowToAssetRowSupportsFallbackFields()
  testMapHostRowToAssetRowFallsBackToStatusWhenExplicitSignalMissing()
  testMapHostRowToAssetRowPrefersExplicitConnectivitySignals()
  testMapHostRowToAssetRowPreservesCpuAndMemoryShape()
  testCreateHostAssetFormModelDefaults()
  console.log('cmdb host asset presentation tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
