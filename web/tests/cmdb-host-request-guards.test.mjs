import assert from 'node:assert/strict'

import {
  createHostMonitorBatchRequester,
  normalizeHostMonitorIds
} from '../src/utils/hostMonitorBatchRequester.mjs'
import { createMessageDeduper } from '../src/utils/requestFeedback.mjs'

async function testHostMonitorBatchRequesterDedupesInFlightRequests() {
  const calls = []
  let resolveRequest
  const requestFn = (ids) => {
    calls.push(ids)
    return new Promise((resolve) => {
      resolveRequest = () => resolve({ data: { code: 200, data: { ids } } })
    })
  }

  const requester = createHostMonitorBatchRequester(requestFn)
  const hosts = [{ id: 2 }, { id: 1 }, { id: 2 }]

  const firstPromise = requester.fetch(hosts)
  const secondPromise = requester.fetch([{ id: 1 }, { id: 2 }])

  assert.equal(calls.length, 1, 'same host batch should reuse the in-flight request')
  assert.deepEqual(calls, ['1,2'])

  resolveRequest()

  const [firstResult, secondResult] = await Promise.all([firstPromise, secondPromise])
  assert.deepEqual(firstResult, secondResult, 'reused requests should resolve with the same payload')

  await requester.fetch(hosts)
  assert.equal(calls.length, 2, 'same host batch should be allowed again after the previous request settles')
}

function testNormalizeHostMonitorIds() {
  assert.equal(normalizeHostMonitorIds([{ id: 3 }, { id: 1 }, { id: 3 }, { id: 2 }]), '1,2,3')
  assert.equal(normalizeHostMonitorIds(['2', 1, '2', null, undefined, '']), '1,2')
  assert.equal(normalizeHostMonitorIds([]), '')
}

function testMessageDeduper() {
  let now = 1000
  const deduper = createMessageDeduper({
    windowMs: 1500,
    now: () => now
  })

  assert.equal(deduper.shouldNotify('timeout'), true, 'first message should pass')
  assert.equal(deduper.shouldNotify('timeout'), false, 'duplicate message inside the guard window should be suppressed')
  assert.equal(deduper.shouldNotify('network'), true, 'different message keys should not block each other')

  now += 1501
  assert.equal(deduper.shouldNotify('timeout'), true, 'message should be allowed again after the guard window')
}

async function main() {
  testNormalizeHostMonitorIds()
  await testHostMonitorBatchRequesterDedupesInFlightRequests()
  testMessageDeduper()
  console.log('cmdb host request guard tests passed')
}

main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
