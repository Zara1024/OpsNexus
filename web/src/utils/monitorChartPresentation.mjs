function niceRoundUp(value) {
  if (!Number.isFinite(value) || value <= 0) {
    return 1
  }

  const magnitude = 10 ** Math.floor(Math.log10(value))
  const normalized = value / magnitude

  if (normalized <= 1) return magnitude
  if (normalized <= 2) return 2 * magnitude
  if (normalized <= 5) return 5 * magnitude
  return 10 * magnitude
}

export function centerSinglePointSeries(labels = [], values = []) {
  if (labels.length === 1 && values.length === 1) {
    return {
      labels: ['', labels[0], ''],
      values: [null, values[0], null],
      snapshotOnly: true,
      latestLabel: labels[0]
    }
  }

  return {
    labels,
    values,
    snapshotOnly: false,
    latestLabel: labels[labels.length - 1] || ''
  }
}

export function computeChartAxisMax(valueGroups = [], { cap } = {}) {
  const numericValues = valueGroups
    .flat()
    .filter(value => Number.isFinite(Number(value)))
    .map(value => Number(value))

  if (numericValues.length === 0) {
    return 1
  }

  const maxValue = Math.max(...numericValues)
  const padded = niceRoundUp(maxValue * 1.15)
  if (Number.isFinite(cap)) {
    return Math.min(cap, Math.max(1, padded))
  }
  return Math.max(0.1, padded)
}
