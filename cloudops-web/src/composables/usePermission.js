import { computed } from 'vue';
import { useUserStore } from '@/stores/user';
// usePermission 提供统一的权限判断能力。
export function usePermission() {
    const userStore = useUserStore();
    const permissions = computed(() => userStore.permissions || []);
    const hasPermission = (permission) => {
        const current = permissions.value;
        if (!permission || current.includes('*:*:*') || current.includes('system:*:*') || current.includes('cmdb:*:*')) {
            return true;
        }
        if (Array.isArray(permission)) {
            return permission.some((p) => current.includes(p));
        }
        return current.includes(permission);
    };
    return {
        permissions,
        hasPermission,
    };
}
