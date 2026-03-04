import { defineStore } from 'pinia';
export const usePermissionStore = defineStore('permission', {
    state: () => ({
        menus: [],
    }),
    actions: {
        buildByMenus(menus) {
            this.menus = menus;
        },
        reset() {
            this.menus = [];
        },
    },
});
