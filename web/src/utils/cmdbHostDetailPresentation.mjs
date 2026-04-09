const hasOwnId = (value) => (
  value
  && typeof value === 'object'
  && 'id' in value
  && value.id !== undefined
  && value.id !== null
  && String(value.id).trim() !== ''
)

export function resolveHostSyncTarget(candidate, fallback = null) {
  if (hasOwnId(candidate)) {
    return candidate
  }

  if (hasOwnId(fallback)) {
    return fallback
  }

  return null
}

export function formatHostResourceValue(value, unit = '') {
  if (value === undefined || value === null) {
    return '-'
  }

  const normalized = String(value).trim()
  if (!normalized) {
    return '-'
  }

  if (!unit) {
    return normalized
  }

  const normalizedUnit = String(unit).trim()
  if (!normalizedUnit) {
    return normalized
  }

  if (normalized.toLowerCase().endsWith(normalizedUnit.toLowerCase())) {
    return normalized
  }

  if (/[a-zA-Z\u4e00-\u9fa5]+$/.test(normalized)) {
    return normalized
  }

  return `${normalized}${normalizedUnit}`
}
