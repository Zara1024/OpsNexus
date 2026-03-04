import { defineStore } from 'pinia'
import type { MenuNode } from '@/api/auth'

export const usePermissionStore = defineStore('permission', {
  state: () => ({
    menus: [] as MenuNode[],
  }),
  actions: {
    buildByMenus(menus: MenuNode[]) {
      this.menus = menus
    },
    reset() {
      this.menus = []
    },
  },
})
