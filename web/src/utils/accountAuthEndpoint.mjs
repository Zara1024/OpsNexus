export function parseAccountHostPortInput(value) {
  const raw = String(value || '').trim()
  if (!raw) {
    return { host: '', port: null, valid: false }
  }

  const separatorIndex = raw.lastIndexOf(':')
  if (separatorIndex <= 0 || separatorIndex === raw.length - 1) {
    return { host: raw, port: null, valid: false }
  }

  const host = raw.slice(0, separatorIndex).trim()
  const portText = raw.slice(separatorIndex + 1).trim()
  const port = Number(portText)
  const valid = host !== '' && Number.isInteger(port) && port > 0 && port <= 65535

  return {
    host,
    port: valid ? port : null,
    valid
  }
}

export function formatAccountHostPortInput(host, port) {
  const normalizedHost = String(host || '').trim()
  if (!normalizedHost) return ''
  if (!port) return normalizedHost
  return `${normalizedHost}:${port}`
}
