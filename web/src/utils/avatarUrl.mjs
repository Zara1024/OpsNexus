const ABSOLUTE_URL_PATTERN = /^[a-zA-Z][a-zA-Z\d+\-.]*:/
const UPLOAD_PATH_PREFIX = '/api/v1/upload/'

function getCurrentOrigin(currentOrigin) {
  if (currentOrigin) {
    return currentOrigin
  }
  if (typeof window !== 'undefined' && window.location?.origin) {
    return window.location.origin
  }
  return ''
}

export function resolveAvatarUrl(value, options = {}) {
  const source = typeof value === 'string' ? value.trim() : ''
  if (!source) {
    return ''
  }

  const currentOrigin = getCurrentOrigin(options.currentOrigin)
  if (!currentOrigin) {
    return source
  }

  let parsed
  try {
    parsed = new URL(source, currentOrigin)
  } catch (error) {
    return ''
  }

  if (!['http:', 'https:'].includes(parsed.protocol)) {
    return ''
  }

  const isAbsolute = ABSOLUTE_URL_PATTERN.test(source)
  const isSameOrigin = parsed.origin === currentOrigin
  const isUploadAsset = parsed.pathname.startsWith(UPLOAD_PATH_PREFIX)

  if (!isAbsolute) {
    return `${parsed.pathname}${parsed.search}${parsed.hash}`
  }

  if (isSameOrigin) {
    return parsed.toString()
  }

  if (isUploadAsset) {
    return ''
  }

  return parsed.toString()
}

export default resolveAvatarUrl
