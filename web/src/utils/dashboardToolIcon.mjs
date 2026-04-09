const ABSOLUTE_URL_PATTERN = /^[a-zA-Z][a-zA-Z\d+\-.]*:/
const UPLOAD_PATH_PREFIX = '/api/v1/upload/'

function getCurrentOrigin(currentOrigin) {
  if (currentOrigin) {
    return currentOrigin
  }
  if (typeof window !== 'undefined' && window.location?.origin) {
    return window.location.origin
  }
  return 'http://localhost'
}

export function normalizeDashboardToolIconUrl(value, options = {}) {
  const source = typeof value === 'string' ? value.trim() : ''
  if (!source) {
    return ''
  }

  let parsed
  try {
    parsed = new URL(source, getCurrentOrigin(options.currentOrigin))
  } catch (_error) {
    return ''
  }

  if (!['http:', 'https:'].includes(parsed.protocol)) {
    return ''
  }

  if (parsed.pathname.startsWith(UPLOAD_PATH_PREFIX)) {
    return `${parsed.pathname}${parsed.search}${parsed.hash}`
  }

  return ABSOLUTE_URL_PATTERN.test(source)
    ? parsed.toString()
    : `${parsed.pathname}${parsed.search}${parsed.hash}`
}

export default normalizeDashboardToolIconUrl
