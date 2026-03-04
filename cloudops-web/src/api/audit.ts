import axios from 'axios'
import request from '@/utils/request'

const ACCESS_TOKEN_KEY = 'opsnexus_access_token'

export interface AuditLogItem {
  id: number
  user_id: number
  username: string
  ip: string
  method: string
  path: string
  response_code: number
  module: string
  action: string
  description: string
  duration_ms: number
  created_at: string
}

export interface AuditListQuery {
  page?: number
  page_size?: number
  keyword?: string
  module?: string
  action?: string
  start_time?: string
  end_time?: string
}

export const listAuditLogsApi = (params: AuditListQuery) => request.get('/v1/audit/logs', { params }) as Promise<{ list: AuditLogItem[]; total: number }>
export const listLoginLogsApi = (params: AuditListQuery) => request.get('/v1/audit/login-logs', { params }) as Promise<{ list: AuditLogItem[]; total: number }>

export const exportAuditLogsApi = async (params: AuditListQuery) => {
  const token = localStorage.getItem(ACCESS_TOKEN_KEY) || ''
  const res = await axios.get('/api/v1/audit/logs/export', {
    params,
    responseType: 'blob',
    headers: token ? { Authorization: `Bearer ${token}` } : {},
  })
  return res.data as Blob
}
