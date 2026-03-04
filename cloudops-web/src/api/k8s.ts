import request from '@/utils/request'

export interface K8sClusterItem {
  id: number
  name: string
  api_server_url: string
  description: string
  version: string
  status: number
  node_count: number
  pod_count: number
  health?: string
}

export interface K8sNodeItem {
  name: string
  roles: string[]
  status: string
  kubelet_version: string
  os_image: string
  arch: string
  pod_count: number
  cpu_usage: number
  memory_usage: number
  labels?: Record<string, string>
}

export const listK8sClustersApi = () => request.get('/v1/k8s/clusters') as Promise<{ list: K8sClusterItem[] }>
export const createK8sClusterApi = (data: Record<string, unknown>) => request.post('/v1/k8s/clusters', data)
export const updateK8sClusterApi = (id: number, data: Record<string, unknown>) => request.put(`/v1/k8s/clusters/${id}`, data)
export const deleteK8sClusterApi = (id: number) => request.delete(`/v1/k8s/clusters/${id}`)
export const testK8sClusterApi = (id: number) => request.post(`/v1/k8s/clusters/${id}/test`)
export const getK8sClusterApi = (id: number) => request.get(`/v1/k8s/clusters/${id}`) as Promise<K8sClusterItem>

export const listK8sNamespacesApi = (clusterId: number) => request.get(`/v1/k8s/clusters/${clusterId}/namespaces`) as Promise<{ list: string[] }>
export const listK8sNodesApi = (clusterId: number) => request.get(`/v1/k8s/clusters/${clusterId}/nodes`) as Promise<{ list: K8sNodeItem[] }>
export const cordonK8sNodeApi = (clusterId: number, name: string) => request.put(`/v1/k8s/clusters/${clusterId}/nodes/${name}/cordon`)
export const uncordonK8sNodeApi = (clusterId: number, name: string) => request.put(`/v1/k8s/clusters/${clusterId}/nodes/${name}/uncordon`)
export const drainK8sNodeApi = (clusterId: number, name: string) => request.post(`/v1/k8s/clusters/${clusterId}/nodes/${name}/drain`)

export const listDeploymentsApi = (clusterId: number, ns: string) => request.get(`/v1/k8s/clusters/${clusterId}/ns/${ns}/deployments`) as Promise<{ list: any[] }>
export const getDeploymentApi = (clusterId: number, ns: string, name: string) => request.get(`/v1/k8s/clusters/${clusterId}/ns/${ns}/deployments/${name}`)
export const createDeploymentApi = (clusterId: number, ns: string, data: any) => request.post(`/v1/k8s/clusters/${clusterId}/ns/${ns}/deployments`, data)
export const updateDeploymentApi = (clusterId: number, ns: string, name: string, data: any) => request.put(`/v1/k8s/clusters/${clusterId}/ns/${ns}/deployments/${name}`, data)
export const deleteDeploymentApi = (clusterId: number, ns: string, name: string) => request.delete(`/v1/k8s/clusters/${clusterId}/ns/${ns}/deployments/${name}`)
export const scaleDeploymentApi = (clusterId: number, ns: string, name: string, replicas: number) => request.put(`/v1/k8s/clusters/${clusterId}/ns/${ns}/deployments/${name}/scale`, { replicas })
export const restartDeploymentApi = (clusterId: number, ns: string, name: string) => request.post(`/v1/k8s/clusters/${clusterId}/ns/${ns}/deployments/${name}/restart`)

export const listPodsApi = (clusterId: number, ns: string) => request.get(`/v1/k8s/clusters/${clusterId}/ns/${ns}/pods`) as Promise<{ list: any[] }>
export const getPodLogsApi = (clusterId: number, ns: string, pod: string, params?: Record<string, unknown>) =>
  request.get(`/v1/k8s/clusters/${clusterId}/ns/${ns}/pods/${pod}/logs`, { params }) as Promise<{ logs: string }>
export const deletePodApi = (clusterId: number, ns: string, pod: string) => request.delete(`/v1/k8s/clusters/${clusterId}/ns/${ns}/pods/${pod}`)

export const listServicesApi = (clusterId: number, ns: string) => request.get(`/v1/k8s/clusters/${clusterId}/ns/${ns}/services`) as Promise<{ list: any[] }>
export const listIngressesApi = (clusterId: number, ns: string) => request.get(`/v1/k8s/clusters/${clusterId}/ns/${ns}/ingresses`) as Promise<{ list: any[] }>

export const listConfigMapsApi = (clusterId: number, ns: string) => request.get(`/v1/k8s/clusters/${clusterId}/ns/${ns}/configmaps`) as Promise<{ list: any[] }>
export const saveConfigMapApi = (clusterId: number, ns: string, data: any) => request.post(`/v1/k8s/clusters/${clusterId}/ns/${ns}/configmaps`, data)
export const deleteConfigMapApi = (clusterId: number, ns: string, name: string) => request.delete(`/v1/k8s/clusters/${clusterId}/ns/${ns}/configmaps/${name}`)

export const listSecretsApi = (clusterId: number, ns: string, decode = false) =>
  request.get(`/v1/k8s/clusters/${clusterId}/ns/${ns}/secrets`, { params: { decode: decode ? 1 : 0 } }) as Promise<{ list: any[] }>
export const saveSecretApi = (clusterId: number, ns: string, data: any) => request.post(`/v1/k8s/clusters/${clusterId}/ns/${ns}/secrets`, data)
export const deleteSecretApi = (clusterId: number, ns: string, name: string) => request.delete(`/v1/k8s/clusters/${clusterId}/ns/${ns}/secrets/${name}`)
