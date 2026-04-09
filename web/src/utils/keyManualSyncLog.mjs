export function extractManualSyncLogs(rows) {
  const list = Array.isArray(rows) ? rows : []
  return list.filter(item => item?.url === '/api/v1/config/keymanage/sync')
}

export default extractManualSyncLogs
