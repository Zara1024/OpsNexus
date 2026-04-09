<template>
  <div class="empty-state">
    <div class="empty-state__halo"></div>
    <div class="empty-state__icon">
      <slot name="icon">
        <el-icon size="22"><Promotion /></el-icon>
      </slot>
    </div>
    <div class="empty-state__title">{{ title }}</div>
    <div class="empty-state__description">{{ description }}</div>
    <div v-if="$slots.actions" class="empty-state__actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup>
import { Promotion } from '@element-plus/icons-vue'

defineProps({
  title: {
    type: String,
    default: '暂无数据'
  },
  description: {
    type: String,
    default: ''
  }
})
</script>
