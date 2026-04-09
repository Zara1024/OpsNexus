const isWindowsHost = (host) => String(host?.deviceType || 'linux').toLowerCase() === 'windows'

const canUseSSH = (host) => Boolean(
  host?.supportsSsh || (host?.sshIp && host?.sshName && Number(host?.sshKeyId) > 0)
)

const firstNonEmpty = (...values) => {
  for (const value of values) {
    const normalized = String(value || '').trim()
    if (normalized) {
      return normalized
    }
  }
  return ''
}

const toQueryValue = (value) => {
  if (value === undefined || value === null || value === '') {
    return ''
  }
  return String(value)
}

export const buildHostConnectionEntries = (host = {}) => {
  const sshReady = canUseSSH(host)
  const sshTarget = firstNonEmpty(host?.sshIp)
  const sshUser = firstNonEmpty(host?.sshName, 'root')
  const sshPort = Number(host?.sshPort) > 0 ? Number(host.sshPort) : 22

  if (isWindowsHost(host)) {
    const rdpTarget = firstNonEmpty(host?.remoteDomain, host?.sshIp, host?.privateIp, host?.publicIp)
    const rdpPort = Number(host?.rdpPort) > 0 ? Number(host.rdpPort) : 3389
    const rdpUser = firstNonEmpty(host?.rdpUsername, '未配置用户')
    const rdpDomain = firstNonEmpty(host?.rdpDomain)
    const rdpIdentity = rdpDomain ? `${rdpDomain}\\${rdpUser}` : rdpUser
    const disabledReason = '当前主机为 Windows，一期仅开放 RDP 信息，其余能力仍复用 Linux SSH 链路。'

    return [
      {
        key: 'rdp',
        title: 'RDP 连接',
        available: Boolean(rdpTarget),
        reason: rdpTarget ? '' : '当前主机未配置 RDP 地址或管理地址。',
        detail: rdpTarget ? `${rdpIdentity} @ ${rdpTarget}:${rdpPort}` : ''
      },
      {
        key: 'ssh',
        title: 'Web SSH',
        available: false,
        reason: disabledReason
      },
      {
        key: 'command',
        title: '执行命令',
        available: false,
        reason: disabledReason
      },
      {
        key: 'upload',
        title: '文件上传',
        available: false,
        reason: disabledReason
      },
      {
        key: 'sync',
        title: '配置同步',
        available: false,
        reason: disabledReason
      }
    ]
  }

  const unavailableReason = sshReady
    ? ''
    : '当前主机未配置完整 SSH 连接信息，无法通过现有 SSH 链路执行该操作。'

  return [
    {
      key: 'ssh',
      title: 'Web SSH',
      available: sshReady,
      reason: unavailableReason,
      detail: sshReady ? `${sshUser}@${sshTarget}:${sshPort}` : ''
    },
    {
      key: 'command',
      title: '执行命令',
      available: sshReady,
      reason: unavailableReason
    },
    {
      key: 'upload',
      title: '文件上传',
      available: sshReady,
      reason: unavailableReason
    },
    {
      key: 'sync',
      title: '配置同步',
      available: sshReady,
      reason: unavailableReason
    }
  ]
}

export const buildHostAuditRouteQuery = (host = {}, extra = {}) => {
  const query = {}
  const hostId = toQueryValue(host?.id)
  const hostKeyword = toQueryValue(firstNonEmpty(host?.hostName, host?.name))
  const hostIp = toQueryValue(firstNonEmpty(host?.sshIp, host?.privateIp, host?.publicIp))

  if (hostId) query.hostId = hostId
  if (hostKeyword) query.hostKeyword = hostKeyword
  if (hostIp) query.hostIp = hostIp

  Object.entries(extra || {}).forEach(([key, value]) => {
    const normalized = toQueryValue(value)
    if (normalized) {
      query[key] = normalized
    }
  })

  return query
}

const isConnectivityConnected = (item = {}) => {
  if (typeof item.reachable === 'boolean') {
    return item.reachable
  }
  if (typeof item.connected === 'boolean') {
    return item.connected
  }
  if (typeof item.onlineStatus === 'number') {
    return item.onlineStatus === 0
  }
  return item.status === 'connected'
}

export const summarizeBatchConnectivity = (hosts = [], resultMap = {}) => {
  const items = hosts.map((host) => {
    const result = Array.isArray(resultMap)
      ? resultMap.find((item) => Number(item?.hostId || item?.id) === Number(host?.id)) || {}
      : resultMap?.[host?.id] || {}
    const connected = isConnectivityConnected(result)
    return {
      id: host?.id,
      hostName: host?.hostName || `主机-${host?.id || 'unknown'}`,
      status: connected ? 'connected' : 'disconnected',
      reachable: connected,
      reason: connected ? '' : firstNonEmpty(result?.reason, result?.message, result?.monitorUnavailableReason, '连通性检测失败'),
      protocol: result?.protocol || ''
    }
  })

  const connected = items.filter((item) => item.status === 'connected').length
  const disconnected = items.length - connected

  return {
    total: items.length,
    connected,
    disconnected,
    items
  }
}

export const getCommandRiskPresentation = (assessment = {}) => {
  const level = Number(assessment?.riskLevel || 0)
  const label = firstNonEmpty(
    assessment?.riskLevelLabel,
    level >= 2 ? 'high' : (level === 1 ? 'medium' : 'low')
  )
  const text = firstNonEmpty(
    assessment?.riskLevelText,
    level >= 2 ? '高风险' : (level === 1 ? '中风险' : '低风险')
  )

  return {
    level,
    label,
    text,
    tagType: level >= 2 ? 'danger' : (level === 1 ? 'warning' : 'success'),
    requiresConfirmation: Boolean(assessment?.requiresConfirmation),
    reason: firstNonEmpty(assessment?.reason)
  }
}
