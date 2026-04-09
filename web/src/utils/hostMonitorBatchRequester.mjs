function toComparableId(value) {
  if (value === null || value === undefined || value === '') {
    return null
  }

  const numeric = Number(value)
  if (Number.isFinite(numeric) && `${numeric}` === `${value}`.trim()) {
    return `${numeric}`
  }

  return `${value}`.trim()
}

export function normalizeHostMonitorIds(hostsOrIds = []) {
  const ids = hostsOrIds
    .map((item) => (item && typeof item === 'object' ? item.id : item))
    .map(toComparableId)
    .filter(Boolean)

  return [...new Set(ids)].sort((left, right) => Number(left) - Number(right)).join(',')
}

export function createHostMonitorBatchRequester(requestFn) {
  const inFlightRequests = new Map()

  return {
    fetch(hostsOrIds = []) {
      const ids = normalizeHostMonitorIds(hostsOrIds)
      if (!ids) {
        return Promise.resolve(null)
      }

      if (inFlightRequests.has(ids)) {
        return inFlightRequests.get(ids)
      }

      const pendingRequest = Promise.resolve(requestFn(ids))
        .finally(() => {
          inFlightRequests.delete(ids)
        })

      inFlightRequests.set(ids, pendingRequest)
      return pendingRequest
    },
    clear() {
      inFlightRequests.clear()
    }
  }
}
