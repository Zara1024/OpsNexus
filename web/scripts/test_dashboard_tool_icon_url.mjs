import assert from 'node:assert/strict'

import { normalizeDashboardToolIconUrl } from '../src/utils/dashboardToolIcon.mjs'

assert.equal(
  normalizeDashboardToolIconUrl('http://10.0.0.200/api/v1/upload/20260324/33644063.png', {
    currentOrigin: 'http://10.0.0.200:8080'
  }),
  '/api/v1/upload/20260324/33644063.png'
)

assert.equal(
  normalizeDashboardToolIconUrl('/api/v1/upload/20260324/33644063.png', {
    currentOrigin: 'http://10.0.0.200:8080'
  }),
  '/api/v1/upload/20260324/33644063.png'
)

assert.equal(
  normalizeDashboardToolIconUrl('https://cdn.example.com/icon.png', {
    currentOrigin: 'http://10.0.0.200:8080'
  }),
  'https://cdn.example.com/icon.png'
)

assert.equal(
  normalizeDashboardToolIconUrl('', {
    currentOrigin: 'http://10.0.0.200:8080'
  }),
  ''
)

console.log('dashboard tool icon url tests passed')
