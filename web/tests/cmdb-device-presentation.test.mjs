import assert from 'node:assert/strict'

import {
  buildAssetConnectivityTag,
  createDeviceAssetFormModel,
  mapDeviceRowToAssetRow
} from '../src/utils/cmdbAssetPresentation.mjs'

function testMapDeviceRowToAssetRowUsesPrimaryFields() {
  const row = mapDeviceRowToAssetRow({
    id: 41,
    name: 'core-sw-01',
    address: '10.0.0.41',
    accountName: 'netops-admin',
    platform: 'Huawei',
    reachable: true,
    protocolGroup: 'ssh,web',
    sshPort: 22,
    webUrl: 'https://10.0.0.41'
  })

  assert.equal(row.name, 'core-sw-01')
  assert.equal(row.address, '10.0.0.41')
  assert.equal(row.account, 'netops-admin')
  assert.equal(row.platform, 'Huawei')
  assert.deepEqual(row.connectivity, buildAssetConnectivityTag({ reachable: true }))
}

function testMapDeviceRowToAssetRowSupportsFallbackFields() {
  const row = mapDeviceRowToAssetRow({
    name: 'edge-router-02',
    ip: '172.16.20.2',
    accountAlias: 'router-ops',
    deviceType: 'router',
    connected: false
  })

  assert.equal(row.name, 'edge-router-02')
  assert.equal(row.address, '172.16.20.2')
  assert.equal(row.account, 'router-ops')
  assert.equal(row.platform, 'router')
  assert.deepEqual(row.connectivity, buildAssetConnectivityTag({ reachable: false }))
}

function testMapDeviceRowToAssetRowDoesNotRenderObjectAccountFallback() {
  const row = mapDeviceRowToAssetRow({
    name: 'agg-switch-01',
    address: '10.10.10.10',
    account: { id: 88 },
    accountId: 88
  })

  assert.equal(row.account, '88')
  assert.notEqual(row.account, '[object Object]')
}

function testMapDeviceRowToAssetRowUsesUnknownConnectivityWithoutSignal() {
  const row = mapDeviceRowToAssetRow({
    name: 'unknown-signal-device',
    address: '10.10.20.20'
  })

  assert.deepEqual(row.connectivity, {
    text: '未知',
    type: 'info'
  })
}

function testCreateDeviceAssetFormModelDefaults() {
  const form = createDeviceAssetFormModel()

  assert.equal(form.id, '')
  assert.equal(form.name, '')
  assert.equal(form.address, '')
  assert.equal(form.platform, '')
  assert.equal(form.groupId, '')
  assert.equal(form.accountId, '')
  assert.equal(form.protocolGroup, 'ssh')
  assert.equal(form.tags, '')
  assert.equal(form.isActive, true)
  assert.equal(form.remark, '')
  assert.equal(form.deviceType, 'network')
  assert.equal(form.sshPort, 22)
  assert.equal(form.telnetPort, 23)
  assert.equal(form.webUrl, '')
}

async function main() {
  testMapDeviceRowToAssetRowUsesPrimaryFields()
  testMapDeviceRowToAssetRowSupportsFallbackFields()
  testMapDeviceRowToAssetRowDoesNotRenderObjectAccountFallback()
  testMapDeviceRowToAssetRowUsesUnknownConnectivityWithoutSignal()
  testCreateDeviceAssetFormModelDefaults()
  console.log('cmdb device presentation tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
