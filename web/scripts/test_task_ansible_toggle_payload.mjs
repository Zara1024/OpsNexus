import assert from 'node:assert/strict'

import { buildRecurringTogglePayload } from '../src/utils/ansibleTaskPayload.mjs'

assert.deepEqual(
  buildRecurringTogglePayload({
    id: 56,
    is_recurring: 0,
    cron_expr: '0 0 2 * * ?'
  }),
  {
    id: 56,
    isRecurring: 0,
    cronExpr: '0 0 2 * * ?'
  }
)

assert.deepEqual(
  buildRecurringTogglePayload({
    id: 56,
    is_recurring: 1,
    cron_expr: ''
  }),
  {
    id: 56,
    isRecurring: 1,
    cronExpr: ''
  }
)

console.log('task ansible toggle payload tests passed')
