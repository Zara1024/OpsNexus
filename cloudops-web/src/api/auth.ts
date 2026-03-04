import request from '@/utils/request'

export interface CaptchaResp {
  captcha_id: string
  captcha: string
}

export interface LoginPayload {
  username: string
  password: string
  captcha_id: string
  captcha: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  expires_in: number
}

export interface LoginResp {
  token: TokenPair
}

export interface UserInfo {
  user_id: number
  username: string
  nickname: string
  role_codes: string[]
  permissions: string[]
  menus: MenuNode[]
}

export interface MenuNode {
  id: number
  parent_id: number
  menu_name: string
  route_path: string
  component: string
  icon: string
  sort_order: number
  menu_type: 'dir' | 'menu' | 'button'
  permission: string
  children: MenuNode[]
}

export const getCaptchaApi = () => request.get('/v1/auth/captcha') as Promise<CaptchaResp>
export const loginApi = (data: LoginPayload) => request.post('/v1/auth/login', data) as Promise<LoginResp>
export const logoutApi = () => request.post('/v1/auth/logout')
export const refreshApi = (refresh_token: string) => request.post('/v1/auth/refresh', { refresh_token }) as Promise<{ token: TokenPair }>
export const getUserInfoApi = () => request.get('/v1/auth/userinfo') as Promise<UserInfo>
