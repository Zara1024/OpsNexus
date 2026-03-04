import { useUserStore } from '@/stores/user';
function check(el, binding) {
    const userStore = useUserStore();
    const required = binding.value;
    const list = userStore.permissions || [];
    if (!required || list.includes('*:*:*') || list.includes('system:*:*') || list.includes('cmdb:*:*')) {
        return;
    }
    const has = Array.isArray(required) ? required.some((p) => list.includes(p)) : list.includes(required);
    if (!has) {
        el.parentNode?.removeChild(el);
    }
}
// setupPermissionDirective 注册按钮权限指令。
export function setupPermissionDirective(app) {
    app.directive('permission', {
        mounted(el, binding) {
            check(el, binding);
        },
        updated(el, binding) {
            check(el, binding);
        },
    });
}
