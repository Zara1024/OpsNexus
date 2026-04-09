export function createMessageDeduper({ windowMs = 1500, now = () => Date.now() } = {}) {
  const lastShownAt = new Map()

  return {
    shouldNotify(key) {
      const normalizedKey = `${key || ''}`.trim()
      if (!normalizedKey) {
        return true
      }

      const currentTime = now()
      const previousTime = lastShownAt.get(normalizedKey)
      if (previousTime !== undefined && currentTime - previousTime < windowMs) {
        return false
      }

      lastShownAt.set(normalizedKey, currentTime)
      return true
    },
    reset() {
      lastShownAt.clear()
    }
  }
}
