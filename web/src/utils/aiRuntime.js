const FRIENDLY_RUNTIME_ERROR_RULES = [
  {
    patterns: [
      'completed successfully but did not return any text payload',
      'returned empty text'
    ],
    message: '模型网关已响应，但当前返回里没有标准文本内容，系统先回退到内置逻辑。'
  },
  {
    patterns: [
      'context deadline exceeded',
      'client.timeout exceeded',
      'timeout awaiting response headers'
    ],
    message: '模型网关探测超时，系统先回退到内置逻辑，请稍后再试。'
  },
  {
    patterns: [
      'invalid character',
      'cannot unmarshal',
      'unexpected end of json input'
    ],
    message: '模型网关已返回结果，但返回格式和当前解析规则不兼容，系统先回退到内置逻辑。'
  }
]

export function normalizeAIRuntimeErrorMessage(value) {
  const raw = String(value || '').trim()
  if (!raw) {
    return ''
  }

  const lower = raw.toLowerCase()
  for (const rule of FRIENDLY_RUNTIME_ERROR_RULES) {
    if (rule.patterns.some((pattern) => lower.includes(pattern))) {
      return rule.message
    }
  }

  return raw
}

export function getAIRuntimeErrorLabel(status) {
  return status === 'degraded' ? '运行提示' : '最近错误'
}
