import assert from 'node:assert/strict'

import {
  formatHostMonitorMetricValue,
  getHostMonitorPanelState,
  getHostMonitorStatusLabel,
  getHostMonitorUnavailableReason,
  hasHostMonitorCollectionData,
  hasHostMonitorData
} from '../src/utils/hostMonitorState.mjs'

assert.equal(hasHostMonitorData({ dataAvailable: true }), true)
assert.equal(hasHostMonitorData({ dataAvailable: false }), false)
assert.equal(hasHostMonitorData(null), false)

assert.equal(
  getHostMonitorStatusLabel({ collectionStatus: 'offline' }),
  '离线'
)

assert.equal(
  getHostMonitorStatusLabel({ collectionStatus: 'unavailable' }),
  '未接入'
)

assert.equal(
  getHostMonitorUnavailableReason({ collectionStatus: 'offline' }),
  '主机离线，无法采集监控数据'
)

assert.equal(
  getHostMonitorUnavailableReason({ collectionStatus: 'unavailable', unavailableReason: 'Prometheus 未接入或不可达' }),
  'Prometheus 未接入或不可达'
)

assert.equal(
  formatHostMonitorMetricValue(18.236),
  '18.24'
)

assert.equal(
  formatHostMonitorMetricValue(undefined),
  '--'
)

assert.equal(
  hasHostMonitorCollectionData({ topCPU: [{ name: 'sshd' }] }, ['topCPU', 'topMemory']),
  true
)

assert.deepEqual(
  getHostMonitorPanelState(
    {
      collectionStatus: 'unavailable',
      unavailableReason: 'Prometheus 未接入或不可达'
    },
    ['topCPU']
  ),
  {
    ready: false,
    statusLabel: '未接入',
    reason: 'Prometheus 未接入或不可达'
  }
)

console.log('host monitor state tests passed')
