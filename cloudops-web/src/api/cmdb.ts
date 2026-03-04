import request from '@/utils/request'

export interface CMDBGroupNode {
  id: number
  parent_id: number
  group_name: string
  description: string
  sort_order: number
  children: CMDBGroupNode[]
}

export interface CMDBCredentialItem {
  id: number
  name: string
  type: 'password' | 'key'
  username: string
  created_at: string
}

export interface CMDBHostItem {
  id: number
  hostname: string
  inner_ip: string
  outer_ip: string
  os_type: string
  os_version: string
  arch: string
  cpu_cores: number
  memory_gb: number
  disk_gb: number
  group_id: number
  credential_id: number
  cloud_provider: string
  instance_id: string
  region: string
  agent_status: string
  labels: string
  status: number
  created_at: string
  updated_at: string
}

export interface CMDBSSHRecordItem {
  id: number
  user_id: number
  host_id: number
  session_id: string
  start_time: string
  end_time: string
  client_ip: string
  recording_path: string
}

export const listCMDBHostsApi = (params: Record<string, unknown>) =>
  request.get('/v1/cmdb/hosts', { params }) as Promise<{ list: CMDBHostItem[]; total: number }>

export const getCMDBHostApi = (id: number) => request.get(`/v1/cmdb/hosts/${id}`) as Promise<CMDBHostItem>
export const createCMDBHostApi = (data: Record<string, unknown>) => request.post('/v1/cmdb/hosts', data)
export const updateCMDBHostApi = (id: number, data: Record<string, unknown>) => request.put(`/v1/cmdb/hosts/${id}`, data)
export const deleteCMDBHostApi = (id: number) => request.delete(`/v1/cmdb/hosts/${id}`)
export const testCMDBHostSSHApi = (id: number) => request.post(`/v1/cmdb/hosts/${id}/test`)
export const batchCMDBHostsApi = (data: { action: 'delete' | 'enable' | 'disable'; ids: number[] }) => request.post('/v1/cmdb/hosts/batch', data)

export const listCMDBGroupsApi = () => request.get('/v1/cmdb/groups') as Promise<{ list: CMDBGroupNode[] }>
export const createCMDBGroupApi = (data: Record<string, unknown>) => request.post('/v1/cmdb/groups', data)
export const updateCMDBGroupApi = (id: number, data: Record<string, unknown>) => request.put(`/v1/cmdb/groups/${id}`, data)
export const deleteCMDBGroupApi = (id: number) => request.delete(`/v1/cmdb/groups/${id}`)

export const listCMDBCredentialsApi = () => request.get('/v1/cmdb/credentials') as Promise<{ list: CMDBCredentialItem[] }>
export const createCMDBCredentialApi = (data: Record<string, unknown>) => request.post('/v1/cmdb/credentials', data)
export const updateCMDBCredentialApi = (id: number, data: Record<string, unknown>) => request.put(`/v1/cmdb/credentials/${id}`, data)
export const deleteCMDBCredentialApi = (id: number) => request.delete(`/v1/cmdb/credentials/${id}`)

export const listCMDBSSHRecordsApi = (params: { page: number; page_size: number }) =>
  request.get('/v1/cmdb/ssh-records', { params }) as Promise<{ list: CMDBSSHRecordItem[]; total: number }>
