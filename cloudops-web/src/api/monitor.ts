import request from '@/utils/request'

export interface AlertRuleItem {
  id: number
  name: string
  description: string
  severity: 'P1' | 'P2' | 'P3' | 'P4'
  expression: string
  duration: string
  labels: string
  annotations: string
  enabled: boolean
  notify_channels: string
}

export interface AlertEventItem {
  id: number
  rule_id: number
  status: string
  severity: string
  starts_at: string
  ends_at?: string
  labels: string
  annotations: string
  acknowledged_by: number
  silenced_until?: string
}

export interface ChannelItem {
  id: number
  name: string
  type: string
  enabled: boolean
  config: Record<string, string>
  created_at: string
}

export const listAlertRulesApi = () => request.get('/v1/monitor/alert-rules') as Promise<{ list: AlertRuleItem[] }>
export const createAlertRuleApi = (data: Record<string, unknown>) => request.post('/v1/monitor/alert-rules', data)
export const updateAlertRuleApi = (id: number, data: Record<string, unknown>) => request.put(`/v1/monitor/alert-rules/${id}`, data)
export const deleteAlertRuleApi = (id: number) => request.delete(`/v1/monitor/alert-rules/${id}`)
export const toggleAlertRuleApi = (id: number, enabled: boolean) => request.put(`/v1/monitor/alert-rules/${id}/toggle`, { enabled })

export const listActiveAlertsApi = () => request.get('/v1/monitor/alerts') as Promise<{ list: AlertEventItem[] }>
export const listAlertHistoryApi = () => request.get('/v1/monitor/alerts/history') as Promise<{ list: AlertEventItem[] }>
export const ackAlertApi = (id: number) => request.put(`/v1/monitor/alerts/${id}/ack`)
export const silenceAlertApi = (id: number, minutes: number) => request.put(`/v1/monitor/alerts/${id}/silence`, { minutes })

export const listChannelsApi = () => request.get('/v1/monitor/channels') as Promise<{ list: ChannelItem[] }>
export const createChannelApi = (data: Record<string, unknown>) => request.post('/v1/monitor/channels', data)
export const testChannelApi = (id: number) => request.post(`/v1/monitor/channels/${id}/test`)
