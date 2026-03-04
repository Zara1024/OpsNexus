import { defineStore } from 'pinia'
import { getUserInfoApi, loginApi, logoutApi, type LoginPayload, type MenuNode, type TokenPair, type UserInfo } from '@/api/auth'
import { usePermissionStore } from '@/stores/permission'

interface UserState {
  accessToken: string
  refreshToken: string
  userInfo: UserInfo | null
  permissions: string[]
  menus: MenuNode[]
}

const ACCESS_TOKEN_KEY = 'opsnexus_access_token'
const REFRESH_TOKEN_KEY = 'opsnexus_refresh_token'

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    accessToken: localStorage.getItem(ACCESS_TOKEN_KEY) || '',
    refreshToken: localStorage.getItem(REFRESH_TOKEN_KEY) || '',
    userInfo: null,
    permissions: [],
    menus: [],
  }),
  getters: {
    isLogin: (state) => !!state.accessToken,
  },
  actions: {
    setToken(token: TokenPair) {
      this.accessToken = token.access_token
      this.refreshToken = token.refresh_token
      localStorage.setItem(ACCESS_TOKEN_KEY, token.access_token)
      localStorage.setItem(REFRESH_TOKEN_KEY, token.refresh_token)
    },
    syncTokenFromStorage() {
      this.accessToken = localStorage.getItem(ACCESS_TOKEN_KEY) || ''
      this.refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY) || ''
    },
    clearToken() {
      this.accessToken = ''
      this.refreshToken = ''
      localStorage.removeItem(ACCESS_TOKEN_KEY)
      localStorage.removeItem(REFRESH_TOKEN_KEY)
    },
    async login(payload: LoginPayload) {
      const res = await loginApi(payload)
      this.setToken(res.token)
      await this.fetchUserInfo()
    },
    async fetchUserInfo() {
      this.syncTokenFromStorage()
      const res = await getUserInfoApi()
      this.userInfo = res
      this.permissions = res.permissions || []
      this.menus = res.menus || []
      usePermissionStore().buildByMenus(this.menus)
    },
    async logout() {
      if (this.accessToken) {
        await logoutApi()
      }
      this.reset()
    },
    reset() {
      this.clearToken()
      this.userInfo = null
      this.permissions = []
      this.menus = []
      usePermissionStore().reset()
    },
  },
})
