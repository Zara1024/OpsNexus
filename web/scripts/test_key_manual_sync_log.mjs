import assert from 'node:assert/strict'

import { extractManualSyncLogs } from '../src/utils/keyManualSyncLog.mjs'

const rows = [
  {
    id: 10,
    url: '/api/v1/config/keymanage/sync',
    description: '同步云主机',
    createTime: '2026-03-24 17:00:00'
  },
  {
    id: 11,
    url: '/api/v1/task/ansible',
    description: '创建Ansible任务',
    createTime: '2026-03-24 17:01:00'
  }
]

assert.deepEqual(extractManualSyncLogs(rows), [
  {
    id: 10,
    url: '/api/v1/config/keymanage/sync',
    description: '同步云主机',
    createTime: '2026-03-24 17:00:00'
  }
])

console.log('key manual sync log tests passed')
