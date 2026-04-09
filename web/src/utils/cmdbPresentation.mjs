const DATABASE_TYPE_META = {
  1: { label: 'MySQL', tagType: 'success', iconKey: 'mysql' },
  2: { label: 'PostgreSQL', tagType: 'warning', iconKey: 'postgresql' },
  3: { label: 'Redis', tagType: 'danger', iconKey: 'redis' },
  4: { label: 'MongoDB', tagType: 'info', iconKey: 'mongodb' },
  5: { label: 'Elasticsearch', tagType: 'primary', iconKey: 'elasticsearch' },
  6: { label: '达梦数据库', tagType: 'warning', iconKey: 'dameng' },
  7: { label: 'GaussDB', tagType: 'success', iconKey: 'gaussdb' },
  8: { label: 'Oracle', tagType: 'danger', iconKey: 'oracle' }
}

export function getCmdbDatabaseTypeOptions() {
  return Object.entries(DATABASE_TYPE_META).map(([value, meta]) => ({
    value: Number(value),
    label: meta.label,
    tagType: meta.tagType,
    iconKey: meta.iconKey
  }))
}

export function getCmdbDatabaseTypeMeta(type) {
  return DATABASE_TYPE_META[Number(type)] || DATABASE_TYPE_META[1]
}

export function getCmdbDatabaseTypeLabel(type) {
  return getCmdbDatabaseTypeMeta(type).label
}

export function getCmdbDatabasePlatformLabel(type, platform = '') {
  const normalized = String(platform || '').trim()
  return normalized || getCmdbDatabaseTypeLabel(type)
}

export function getCmdbDatabaseTypeTag(type) {
  return getCmdbDatabaseTypeMeta(type).tagType
}

export function getCmdbDatabaseTypeIconKey(type) {
  return getCmdbDatabaseTypeMeta(type).iconKey
}

function normalizeTerminalOutput(text) {
  const normalized = String(text || '')
    .replace(/\r\n/g, '\n')
    .replace(/\r/g, '\n')
    .replace(/\\r\\n/g, '\n')
    .replace(/\\n/g, '\n')
    .replace(/\u001b\[[0-9;]*[A-Za-z]/g, '')
  const lines = normalized.split('\n').filter((line) => {
    const trimmed = line.trim()
    return !trimmed.includes('{{end}}') && !trimmed.includes('.Destination}}')
  })
  return lines.join('\n')
}

export function extractCommandOutput(payload) {
  if (!payload) {
    return ''
  }

  if (typeof payload === 'string') {
    return normalizeTerminalOutput(payload)
  }

  if (typeof payload.output === 'string') {
    return normalizeTerminalOutput(payload.output)
  }

  if (payload.output && typeof payload.output === 'object') {
    return extractCommandOutput(payload.output)
  }

  if (payload.stdout || payload.stderr) {
    return normalizeTerminalOutput([payload.stdout, payload.stderr].filter(Boolean).join('\n'))
  }

  return normalizeTerminalOutput(JSON.stringify(payload, null, 2))
}
