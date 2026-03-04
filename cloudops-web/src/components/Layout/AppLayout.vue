<template>
  <div class="layout">
    <aside class="sidebar">
      <div class="logo">OpsNexus</div>
      <el-menu :default-active="activeMenu" router class="menu" background-color="transparent" text-color="#cfd8ff" active-text-color="#8ea2ff">
        <template v-for="menu in sideMenus" :key="menu.id">
          <el-sub-menu v-if="menu.children?.length" :index="menu.route_path || String(menu.id)">
            <template #title>{{ menu.menu_name }}</template>
            <el-menu-item v-for="child in menu.children" :key="child.id" :index="normalizePath(child.route_path)">
              {{ child.menu_name }}
            </el-menu-item>
          </el-sub-menu>
          <el-menu-item v-else :index="normalizePath(menu.route_path)">{{ menu.menu_name }}</el-menu-item>
        </template>
      </el-menu>
    </aside>

    <div class="main">
      <header class="header">
        <span>{{ title }}</span>
        <el-button link type="primary" @click="handleLogout">退出登录</el-button>
      </header>
      <section class="content">
        <RouterView />
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const title = computed(() => (route.meta.title as string) || 'OpsNexus')
const activeMenu = computed(() => route.path)

const sideMenus = computed(() => {
  const menus = userStore.menus || []
  return menus.filter((item) => item.menu_type !== 'button')
})

const normalizePath = (path: string) => {
  if (!path) return '/dashboard'
  return path.startsWith('/') ? path : `/${path}`
}

const handleLogout = async () => {
  await userStore.logout()
  ElMessage.success('已退出登录')
  router.replace('/login')
}
</script>

<style scoped lang="scss">
.layout { display: flex; width: 100%; height: 100%; }
.sidebar { width: 240px; background: rgba(15, 20, 40, 0.86); padding: 16px; }
.logo { font-size: 20px; font-weight: 600; margin-bottom: 16px; color: #fff; }
.menu { border-right: none; }
.main { flex: 1; display: flex; flex-direction: column; }
.header { height: 56px; display: flex; justify-content: space-between; align-items: center; padding: 0 16px; border-bottom: 1px solid rgba(255,255,255,.08); }
.content { flex: 1; padding: 16px; overflow: auto; }
</style>
