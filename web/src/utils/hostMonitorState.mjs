const DEFAULT_UNAVAILABLE_REASON = '监控链路未接入'

export function hasHostMonitorData(payload) {
  return payload?.dataAvailable === true
}

export function getHostMonitorUnavailableReason(payload) {
  if (payload?.unavailableReason) {
    return payload.unavailableReason
  }
  if (payload?.collectionStatus === 'offline') {
    return '主机离线，无法采集监控数据'
  }
  return DEFAULT_UNAVAILABLE_REASON
}

export function getHostMonitorStatusLabel(payload) {
  if (payload?.collectionStatus === 'offline') {
    return '离线'
  }
  if (payload?.collectionStatus === 'unavailable') {
    return '未接入'
  }
  return '正常'
}

export function formatHostMonitorMetricValue(value, { digits = 2, fallback = '--' } = {}) {
  const numericValue = Number(value)
  if (!Number.isFinite(numericValue)) {
    return fallback
  }
  return numericValue.toFixed(digits)
}

export function hasHostMonitorCollectionData(payload, datasetKeys = []) {
  if (payload?.dataAvailable === true) {
    return true
  }

  return datasetKeys.some((key) => {
    const value = payload?.[key]
    if (Array.isArray(value)) {
      return value.length > 0
    }
    if (Array.isArray(value?.data)) {
      return value.data.length > 0
    }
    return false
  })
}

export function getHostMonitorPanelState(payload, datasetKeys = [], fallbackReason = DEFAULT_UNAVAILABLE_REASON) {
  if (hasHostMonitorCollectionData(payload, datasetKeys)) {
    return {
      ready: true,
      statusLabel: getHostMonitorStatusLabel(payload),
      reason: ''
    }
  }

  return {
    ready: false,
    statusLabel: getHostMonitorStatusLabel(payload),
    reason: payload?.unavailableReason || fallbackReason
  }
}
