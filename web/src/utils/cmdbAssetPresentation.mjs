import { getCmdbDatabasePlatformLabel } from './cmdbPresentation.mjs'

const normalizeDisplayValue = (value) => {
  if (typeof value === 'string' || typeof value === 'number') {
    return String(value).trim()
  }
  return ''
}

const firstNonEmpty = (...values) => {
  for (const value of values) {
    const normalized = normalizeDisplayValue(value)
    if (normalized) {
      return normalized
    }
  }
  return ''
}

const buildDatabaseHostPortAddress = (asset = {}) => {
  const host = normalizeDisplayValue(asset?.host)
  const port = normalizeDisplayValue(asset?.port)

  if (host && port) {
    return `${host}:${port}`
  }

  return host
}

const buildAddressByType = (asset = {}) => {
  const type = String(asset?.type || '').toLowerCase()
  if (type === 'host') {
    return firstNonEmpty(
      asset?.sshIp,
      asset?.privateIp,
      asset?.publicIp,
      asset?.address,
      asset?.ip
    )
  }
  if (type === 'device') {
    return firstNonEmpty(asset?.address, asset?.ip)
  }
  if (type === 'database') {
    return firstNonEmpty(
      asset?.address,
      asset?.endpoint,
      buildDatabaseHostPortAddress(asset)
    )
  }
  return firstNonEmpty(asset?.address, asset?.sshIp)
}

const CONNECTIVITY_TAGS = {
  success: { text: '成功', type: 'success' },
  failure: { text: '失败', type: 'danger' },
  unknown: { text: '未知', type: 'info' }
}

const parseBooleanLike = (value) => {
  if (typeof value === 'boolean') {
    return value
  }
  if (typeof value === 'number') {
    return value !== 0
  }
  if (typeof value === 'string') {
    const normalized = value.trim().toLowerCase()
    if (!normalized) {
      return false
    }
    return normalized === 'true' || normalized === '1'
  }
  return false
}

export function buildAssetAddress(asset = {}) {
  return buildAddressByType(asset)
}

export function buildAssetAccountLabel(asset = {}) {
  return firstNonEmpty(asset?.accountName, asset?.account?.name, asset?.owner)
}

export function buildAssetConnectivityTag(asset = {}) {
  const connected = [
    asset?.reachable,
    asset?.connected,
    asset?.online
  ].some(parseBooleanLike)
  return { ...(connected ? CONNECTIVITY_TAGS.success : CONNECTIVITY_TAGS.failure) }
}

const SUPPORTED_ASSET_TYPES = new Set(['host', 'device', 'database'])
const FORM_SECTION_BASE = { key: 'basic', title: '基本设置' }

export function createAssetFormSections(assetType = '') {
  const normalized = String(assetType || '').trim().toLowerCase()
  const type = SUPPORTED_ASSET_TYPES.has(normalized) ? normalized : 'host'
  return [{ ...FORM_SECTION_BASE, assetType: type }]
}

const HOST_ASSET_FORM_MODEL = {
  hostName: '',
  ip: '',
  port: 22,
  username: '',
  authId: '',
  groupId: '',
  remark: '',
  deviceType: 'linux',
  remoteDomain: '',
  rdpPort: 3389,
  rdpUsername: '',
  rdpPassword: '',
  rdpDomain: '',
  name: '',
  address: '',
  account: '',
  platform: ''
}

const DEVICE_ASSET_FORM_MODEL = {
  id: '',
  name: '',
  address: '',
  platform: '',
  groupId: '',
  accountId: '',
  protocolGroup: 'ssh',
  tags: '',
  isActive: true,
  remark: '',
  deviceType: 'network',
  sshPort: 22,
  telnetPort: 23,
  webUrl: ''
}

const DEVICE_PLATFORM_OPTIONS = Object.freeze([
  { label: 'H3C', value: 'H3C' },
  { label: 'Cisco', value: 'Cisco' },
  { label: 'Huawei', value: 'Huawei' },
  { label: 'Juniper', value: 'Juniper' },
  { label: 'Ruijie', value: 'Ruijie' },
  { label: '其他', value: '其他' }
])

const DATABASE_ASSET_FORM_MODEL = {
  id: '',
  name: '',
  address: '',
  platform: getCmdbDatabasePlatformLabel(1, ''),
  groupId: '',
  defaultDatabase: '',
  protocolGroup: 'default',
  accountId: '',
  tags: '',
  isActive: true,
  remark: '',
  description: '',
  type: 1
}

function buildHostAssetAddress(host = {}) {
  return buildAssetAddress({ ...host, type: 'host' })
}

function resolveHostConnectivity(host = {}) {
  const explicitValues = [
    host?.reachable,
    host?.connected,
    host?.online,
    host?.isAlive,
    host?.isOnline,
    host?.alive
  ]
  for (const value of explicitValues) {
    if (value !== undefined && value !== null && String(value).trim() !== '') {
      return parseBooleanLike(value)
    }
  }

  if (host?.onlineStatus !== undefined && host?.onlineStatus !== null) {
    return Number(host.onlineStatus) === 0
  }

  return Number(host?.status) === 1
}

function resolveDeviceConnectivity(device = {}) {
  const explicitValues = [
    device?.reachable,
    device?.connected,
    device?.online
  ]
  const hasExplicitSignal = explicitValues.some((value) => (
    value !== undefined
    && value !== null
    && String(value).trim() !== ''
  ))

  if (!hasExplicitSignal) {
    return { ...CONNECTIVITY_TAGS.unknown }
  }

  return buildAssetConnectivityTag({
    reachable: device?.reachable,
    connected: device?.connected,
    online: device?.online
  })
}

export function mapHostRowToAssetRow(host = {}) {
  const connectivity = buildAssetConnectivityTag({
    reachable: resolveHostConnectivity(host)
  })
  return {
    ...host,
    name: firstNonEmpty(host?.hostName, host?.name),
    address: buildHostAssetAddress(host),
    account: firstNonEmpty(host?.sshName),
    platform: firstNonEmpty(host?.os, host?.deviceType),
    connectivity
  }
}

export function createHostAssetFormModel() {
  return { ...HOST_ASSET_FORM_MODEL }
}

export function mapDeviceRowToAssetRow(device = {}) {
  return {
    ...device,
    name: firstNonEmpty(device?.name),
    address: buildAssetAddress({ ...device, type: 'device' }),
    account: firstNonEmpty(
      device?.accountName,
      device?.accountAlias,
      device?.account?.name,
      device?.accountId
    ),
    platform: firstNonEmpty(device?.platform, device?.deviceType),
    connectivity: resolveDeviceConnectivity(device)
  }
}

export function createDeviceAssetFormModel() {
  return { ...DEVICE_ASSET_FORM_MODEL }
}

export function getCmdbDevicePlatformOptions(selectedPlatform = '') {
  const options = DEVICE_PLATFORM_OPTIONS.map((option) => ({ ...option }))
  const normalizedSelectedPlatform = normalizeDisplayValue(selectedPlatform)

  if (
    normalizedSelectedPlatform
    && !options.some((option) => option.value === normalizedSelectedPlatform)
  ) {
    options.push({
      label: `历史值：${normalizedSelectedPlatform}`,
      value: normalizedSelectedPlatform,
      legacy: true
    })
  }

  return options
}

function resolveDatabaseConnectivity(database = {}) {
  const explicitValues = [
    database?.reachable,
    database?.connected,
    database?.online
  ]
  const hasExplicitSignal = explicitValues.some((value) => (
    value !== undefined
    && value !== null
    && String(value).trim() !== ''
  ))

  if (hasExplicitSignal) {
    return buildAssetConnectivityTag({
      reachable: database?.reachable,
      connected: database?.connected,
      online: database?.online
    })
  }

  if (database?.isActive !== undefined && database?.isActive !== null) {
    return buildAssetConnectivityTag({
      reachable: database?.isActive
    })
  }

  return { ...CONNECTIVITY_TAGS.unknown }
}

export function mapDatabaseRowToAssetRow(database = {}) {
  return {
    ...database,
    name: firstNonEmpty(database?.name),
    address: buildAssetAddress({ ...database, type: 'database' }),
    account: firstNonEmpty(
      database?.accountAlias,
      database?.accountName,
      database?.account?.name,
      database?.accountId
    ),
    platform: getCmdbDatabasePlatformLabel(database?.type, database?.platform),
    defaultDatabase: firstNonEmpty(database?.defaultDatabase),
    connectivity: resolveDatabaseConnectivity(database)
  }
}

export function createDatabaseAssetFormModel() {
  return { ...DATABASE_ASSET_FORM_MODEL }
}
