import assert from 'node:assert/strict'

import {
  createDatabaseAssetFormModel,
  mapDatabaseRowToAssetRow
} from '../src/utils/cmdbAssetPresentation.mjs'

function testMapDatabaseRowToAssetRowUsesSuccessConnectivityText() {
  const row = mapDatabaseRowToAssetRow({
    id: 7,
    name: 'orders-primary',
    address: '10.20.30.40:3306',
    accountName: 'dbadmin',
    platform: 'MySQL 8.0',
    defaultDatabase: 'orders',
    reachable: true,
    type: 1
  })

  assert.equal(row.name, 'orders-primary')
  assert.equal(row.address, '10.20.30.40:3306')
  assert.equal(row.account, 'dbadmin')
  assert.equal(row.platform, 'MySQL 8.0')
  assert.equal(row.defaultDatabase, 'orders')
  assert.deepEqual(row.connectivity, {
    text: '成功',
    type: 'success'
  })
}

function testMapDatabaseRowToAssetRowUsesFailureConnectivityText() {
  const row = mapDatabaseRowToAssetRow({
    id: 9,
    name: 'orders-replica',
    address: '10.20.30.42:3306',
    accountName: 'readonly',
    reachable: false
  })

  assert.deepEqual(row.connectivity, {
    text: '失败',
    type: 'danger'
  })
}

function testMapDatabaseRowToAssetRowUsesUnknownConnectivityTextWithoutSignal() {
  const row = mapDatabaseRowToAssetRow({
    id: 8,
    name: 'orders-shadow',
    address: '10.20.30.41:3306',
    accountName: 'readonly'
  })

  assert.deepEqual(row.connectivity, {
    text: '未知',
    type: 'info'
  })
}

function testMapDatabaseRowToAssetRowFallsBackPlatformFromType() {
  const row = mapDatabaseRowToAssetRow({
    id: 10,
    name: 'orders-analytics',
    type: 2
  })

  assert.equal(row.platform, 'PostgreSQL')
}

function testMapDatabaseRowToAssetRowFallsBackAddressFromHostAndPort() {
  const row = mapDatabaseRowToAssetRow({
    id: 11,
    name: 'orders-cache',
    host: '10.20.30.50',
    port: 6379,
    type: 3
  })

  assert.equal(row.address, '10.20.30.50:6379')
}

function testCreateDatabaseAssetFormModelDefaults() {
  const form = createDatabaseAssetFormModel()

  assert.equal(form.id, '')
  assert.equal(form.name, '')
  assert.equal(form.address, '')
  assert.equal(form.platform, 'MySQL')
  assert.equal(form.groupId, '')
  assert.equal(form.defaultDatabase, '')
  assert.equal(form.protocolGroup, 'default')
  assert.equal(form.accountId, '')
  assert.equal(form.tags, '')
  assert.equal(form.isActive, true)
  assert.equal(form.remark, '')
  assert.equal(form.type, 1)
}

async function main() {
  testMapDatabaseRowToAssetRowUsesSuccessConnectivityText()
  testMapDatabaseRowToAssetRowUsesFailureConnectivityText()
  testMapDatabaseRowToAssetRowUsesUnknownConnectivityTextWithoutSignal()
  testMapDatabaseRowToAssetRowFallsBackPlatformFromType()
  testMapDatabaseRowToAssetRowFallsBackAddressFromHostAndPort()
  testCreateDatabaseAssetFormModelDefaults()
  console.log('cmdb database asset presentation tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
