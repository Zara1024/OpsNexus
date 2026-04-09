import yaml from 'js-yaml'

function ensureObject(value) {
  return value && typeof value === 'object' && !Array.isArray(value) ? value : {}
}

function normalizeStringRecord(value) {
  const record = ensureObject(value)
  return Object.fromEntries(
    Object.entries(record).map(([key, item]) => {
      if (typeof item === 'string') {
        return [key, item]
      }
      if (item === null || item === undefined) {
        return [key, '']
      }
      if (typeof item === 'object') {
        return [key, JSON.stringify(item, null, 2)]
      }
      return [key, String(item)]
    })
  )
}

function parseManifest(yamlContent, expectedKind) {
  const manifest = yaml.load(String(yamlContent || ''))
  if (!manifest || typeof manifest !== 'object' || Array.isArray(manifest)) {
    throw new Error('YAML 内容格式无效')
  }

  const metadata = ensureObject(manifest.metadata)
  const name = typeof metadata.name === 'string' ? metadata.name.trim() : ''
  if (!name) {
    throw new Error('YAML metadata.name 不能为空')
  }

  const kind = typeof manifest.kind === 'string' ? manifest.kind.trim() : ''
  if (expectedKind && kind && kind !== expectedKind) {
    throw new Error(`YAML kind 必须为 ${expectedKind}`)
  }

  return { manifest, metadata, name }
}

export function buildConfigMapCreatePayload(yamlContent) {
  const { manifest, metadata, name } = parseManifest(yamlContent, 'ConfigMap')
  return {
    name,
    labels: ensureObject(metadata.labels),
    data: normalizeStringRecord(manifest.data),
    binaryData: normalizeStringRecord(manifest.binaryData),
    immutable: Boolean(manifest.immutable)
  }
}

export function buildSecretCreatePayload(yamlContent) {
  const { manifest, metadata, name } = parseManifest(yamlContent, 'Secret')
  return {
    name,
    labels: ensureObject(metadata.labels),
    type: typeof manifest.type === 'string' && manifest.type.trim() ? manifest.type.trim() : 'Opaque',
    data: normalizeStringRecord(manifest.data),
    stringData: normalizeStringRecord(manifest.stringData),
    immutable: Boolean(manifest.immutable)
  }
}
