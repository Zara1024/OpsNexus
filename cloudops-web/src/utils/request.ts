import axios from 'axios'
import { ElMessage } from 'element-plus'

interface ApiBody<T = unknown> {
  code: number
  message: string
  data: T
}

const ACCESS_TOKEN_KEY = 'opsnexus_access_token'
const REFRESH_TOKEN_KEY = 'opsnexus_refresh_token'

const request: any = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

request.interceptors.request.use((config: any) => {
  const token = localStorage.getItem(ACCESS_TOKEN_KEY)
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  (response: any) => {
    const newAccessToken = response.headers['x-new-access-token'] as string | undefined
    const newRefreshToken = response.headers['x-new-refresh-token'] as string | undefined
    if (newAccessToken && newRefreshToken) {
      localStorage.setItem(ACCESS_TOKEN_KEY, newAccessToken)
      localStorage.setItem(REFRESH_TOKEN_KEY, newRefreshToken)
    }

    const body = response.data as ApiBody
    if (body.code !== 200) {
      ElMessage.error(body.message || '请求失败')
      return Promise.reject(new Error(body.message || '请求失败'))
    }
    return body.data
  },
  (error: any) => {
    if (error?.response?.status === 401) {
      localStorage.removeItem(ACCESS_TOKEN_KEY)
      localStorage.removeItem(REFRESH_TOKEN_KEY)
      if (location.pathname !== '/login') {
        location.href = '/login'
      }
    }
    return Promise.reject(error)
  },
)

export default request
