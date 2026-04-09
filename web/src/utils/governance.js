const RISK_TAG_TYPE_MAP = Object.freeze({
  high: 'danger',
  medium: 'warning',
  low: 'success'
})

const RISK_LABEL_MAP = Object.freeze({
  high: '高',
  medium: '中',
  low: '低'
})

const PATH_KEY_MAP = Object.freeze({
  k8sWorkload: ['k8s_workload_path', 'k8sWorkloadPath'],
  alertCenter: ['alert_center_path', 'alertCenterPath'],
  aiDiagnosis: ['ai_diagnosis_path', 'aiDiagnosisPath']
})

export const getGovernanceSuggestion = (governance) => {
  if (!governance) {
    return null
  }
  return governance.latest_capacity_suggestion || governance.latestCapacitySuggestion || null
}

export const getGovernanceRiskLevel = (governance) => {
  const suggestion = getGovernanceSuggestion(governance)
  return String(suggestion?.risk_level || suggestion?.riskLevel || '').trim().toLowerCase()
}

export const getGovernanceRiskLabel = (governance) => RISK_LABEL_MAP[getGovernanceRiskLevel(governance)] || ''

export const getGovernanceTagType = (governance) => {
  if (!governance?.enabled) return 'info'
  if (governance?.blocking) return 'danger'
  return RISK_TAG_TYPE_MAP[getGovernanceRiskLevel(governance)] || 'info'
}

export const getGovernanceBlockingReason = (governance) =>
  governance?.blocking_reason || governance?.blockingReason || ''

export const getGovernanceSuggestionPath = (governance, target) => {
  const suggestion = getGovernanceSuggestion(governance)
  const keys = PATH_KEY_MAP[target] || []
  for (const key of keys) {
    const value = suggestion?.[key]
    if (typeof value === 'string' && value.trim()) {
      return value
    }
  }
  return ''
}

export const getGovernanceDisplayText = (governance, options = {}) => {
  const {
    disabledText = '未启用',
    unmappedText = '待映射',
    pendingText = '待预检',
    blockedFallbackText = '治理拦截',
    unknownRiskText = '风险:待确认'
  } = options

  if (!governance?.enabled) return disabledText
  if (!governance?.configured) return unmappedText
  if (!getGovernanceSuggestion(governance)) return pendingText
  if (governance?.blocking) {
    return `阻断:${getGovernanceBlockingReason(governance) || blockedFallbackText}`
  }
  const riskLabel = getGovernanceRiskLabel(governance)
  return riskLabel ? `风险:${riskLabel}` : unknownRiskText
}
