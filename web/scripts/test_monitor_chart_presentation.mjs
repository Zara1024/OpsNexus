import assert from 'node:assert/strict'

import {
  centerSinglePointSeries,
  computeChartAxisMax
} from '../src/utils/monitorChartPresentation.mjs'

assert.deepEqual(
  centerSinglePointSeries(['20:47'], [47.68]),
  {
    labels: ['', '20:47', ''],
    values: [null, 47.68, null],
    snapshotOnly: true,
    latestLabel: '20:47'
  }
)

assert.deepEqual(
  centerSinglePointSeries(['20:40', '20:47'], [31.2, 47.68]),
  {
    labels: ['20:40', '20:47'],
    values: [31.2, 47.68],
    snapshotOnly: false,
    latestLabel: '20:47'
  }
)

assert.equal(
  computeChartAxisMax([[47.68]], { cap: 100 }),
  100
)

assert.equal(
  computeChartAxisMax([[8.33]], { cap: 100 }),
  10
)

assert.equal(
  computeChartAxisMax([[0.19], [0.17]]),
  0.5
)

console.log('monitor chart presentation tests passed')
