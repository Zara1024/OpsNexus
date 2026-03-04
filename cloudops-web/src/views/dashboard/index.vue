<template>
  <el-card>
    <h2>仪表盘</h2>
    <p>欢迎回来，{{ userStore.userInfo?.nickname || userStore.userInfo?.username }}</p>
    <p>当前权限数量：{{ userStore.permissions.length }}</p>
  </el-card>
</template>

<script setup lang="ts">
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
</script>
