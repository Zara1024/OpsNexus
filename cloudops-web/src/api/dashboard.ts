import request from '@/utils/request'

export interface DashboardOverview {
  asset: {
    host_total: number
    host_online: number
    host_offline: number
    host_alert: number
    k8s_cluster_total: number
    db_instance_total: number
  }
  alert: {
    active_total: number
    today_handled: number
    severity_stats: Array<{ severity: string; count: number }>
  }
  resource: {
    cpu_avg: number
    memory_avg: number
    disk_avg: number
  }
  recent_ops: Array<{ id: number; username: string; module: string; action: string; path: string; created_at: string }>
  recent_alert_top5: Array<{ severity: string; count: number }>
}

export interface DashboardTrends {
  days: number
  alert_trend: Array<{ date: string; count: number }>
  resource_trend: Array<{ date: string; cpu: number; memory: number; disk: number }>
}

export const getDashboardOverviewApi = () => request.get('/v1/dashboard/overview') as Promise<DashboardOverview>
export const getDashboardTrendsApi = (days = 7) => request.get('/v1/dashboard/trends', { params: { days } }) as Promise<DashboardTrends>
