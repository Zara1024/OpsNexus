import test from 'node:test'
import assert from 'node:assert/strict'

import {
  buildHostConnectionEntries,
  buildHostAuditRouteQuery,
  getCommandRiskPresentation,
  summarizeBatchConnectivity
} from './cmdbHostPhase1.mjs'

test('buildHostConnectionEntries returns Linux connection actions when SSH is ready', () => {
  const entries = buildHostConnectionEntries({
    id: 718,
    deviceType: 'linux',
    hostName: 'opsnexus-local-verify',
    sshIp: '10.0.0.200',
    sshName: 'root',
    sshKeyId: 33,
    supportsSsh: true
  })

  assert.equal(entries.length, 4)
  assert.deepEqual(
    entries.map((item) => item.key),
    ['ssh', 'command', 'upload', 'sync']
  )
  assert.ok(entries.every((item) => item.available))
})

test('buildHostConnectionEntries returns Windows RDP plus disabled Linux-only actions', () => {
  const entries = buildHostConnectionEntries({
    id: 719,
    deviceType: 'windows',
    hostName: 'opsnexus-win',
    sshIp: '10.0.0.201',
    remoteDomain: '10.0.0.201',
    rdpPort: 3389,
    rdpUsername: 'Administrator'
  })

  assert.deepEqual(
    entries.map((item) => item.key),
    ['rdp', 'ssh', 'command', 'upload', 'sync']
  )
  assert.equal(entries[0].available, true)
  assert.ok(entries.slice(1).every((item) => item.available === false))
  assert.ok(entries.slice(1).every((item) => item.reason))
})

test('buildHostAuditRouteQuery carries host filters for audit deep links', () => {
  const query = buildHostAuditRouteQuery({
    id: 718,
    hostName: 'opsnexus-local-verify',
    sshIp: '10.0.0.200'
  }, { riskLevel: 2 })

  assert.deepEqual(query, {
    hostId: '718',
    hostKeyword: 'opsnexus-local-verify',
    hostIp: '10.0.0.200',
    riskLevel: '2'
  })
})

test('summarizeBatchConnectivity groups connected and disconnected hosts', () => {
  const summary = summarizeBatchConnectivity(
    [
      { id: 1, hostName: 'linux-a' },
      { id: 2, hostName: 'linux-b' },
      { id: 3, hostName: 'linux-c' }
    ],
    {
      1: { onlineStatus: 0 },
      2: { onlineStatus: 1, monitorUnavailableReason: 'SSH dial timeout' },
      3: {}
    }
  )

  assert.equal(summary.total, 3)
  assert.equal(summary.connected, 1)
  assert.equal(summary.disconnected, 2)
  assert.equal(summary.items[1].reason, 'SSH dial timeout')
  assert.equal(summary.items[2].status, 'disconnected')
})

test('getCommandRiskPresentation normalizes low medium and high risk states', () => {
  assert.deepEqual(
    getCommandRiskPresentation({
      riskLevel: 0,
      riskLevelLabel: 'low',
      riskLevelText: '低风险',
      requiresConfirmation: false,
      reason: ''
    }),
    {
      level: 0,
      label: 'low',
      text: '低风险',
      tagType: 'success',
      requiresConfirmation: false,
      reason: ''
    }
  )

  assert.deepEqual(
    getCommandRiskPresentation({
      riskLevel: 1,
      riskLevelLabel: 'medium',
      riskLevelText: '中风险',
      requiresConfirmation: true,
      reason: '包含敏感操作关键字'
    }),
    {
      level: 1,
      label: 'medium',
      text: '中风险',
      tagType: 'warning',
      requiresConfirmation: true,
      reason: '包含敏感操作关键字'
    }
  )

  assert.deepEqual(
    getCommandRiskPresentation({
      riskLevel: 2,
      riskLevelLabel: 'high',
      riskLevelText: '高风险',
      requiresConfirmation: true,
      reason: '包含高风险操作关键字'
    }),
    {
      level: 2,
      label: 'high',
      text: '高风险',
      tagType: 'danger',
      requiresConfirmation: true,
      reason: '包含高风险操作关键字'
    }
  )
})
