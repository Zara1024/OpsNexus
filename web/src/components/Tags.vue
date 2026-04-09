<template>
  <!-- 标签组件 -->
  <div class="tags">
    <el-tag
      v-for="(item, index) in tags"
      :key="item.path"
      class="tag"
      size="default"
      :effect="item.title == $route.meta.tTitle ? 'dark' : 'plain'"
      :closable="index > 0"
      @click="goTo(item.path)"
      @close="close(index)"
    >
      <i class="circular" v-show="item.title == $route.meta.tTitle"></i>
      {{ item.title }}
    </el-tag>
  </div>
</template>

<script>
export default {
  name: 'AppTags',
  data() {
    return {
      tags: [
        {
          path: '/dashboard',
          title: '仪表盘',
        },
      ],
    }
  },
  watch: {
    $route: {
      immediate: true,
      handler(val) {
        const exists = this.tags.find(item => val.path == item.path)
        if (!exists) {
          this.tags.push({
            title: val.meta.tTitle,
            path: val.path,
          })
        }
      },
    },
  },
  methods: {
    // 路由跳转到指定位置
    goTo(path) {
      this.$router.push(path)
    },
    // 点击关闭标签
    close(index) {
      this.tags.splice(index, 1)
    },
  },
}
</script>

<style lang="less" scoped>
.tags {
  padding: 8px 18px 10px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  background: rgba(7, 18, 36, 0.76);
  border-bottom: 1px solid var(--border-subtle, rgba(148, 163, 184, 0.14));
  backdrop-filter: blur(16px);
}

.tag {
  cursor: pointer;
  border-radius: 8px !important;
  transition: all 0.25s ease;
  font-size: 12px;
  font-weight: 500;

  /* 非激活状态（plain） */
  &.el-tag--plain {
    background: rgba(15, 23, 42, 0.6) !important;
    border: 1px solid var(--border-subtle, rgba(148, 163, 184, 0.14)) !important;
    color: var(--text-muted, rgba(148, 163, 184, 0.72)) !important;
  }

  &.el-tag--plain:hover {
    background: rgba(14, 165, 233, 0.1) !important;
    border-color: var(--border-hover, rgba(56, 189, 248, 0.36)) !important;
    color: var(--text-primary, #f8fafc) !important;
    transform: translateY(-1px);
  }

  /* 激活状态（dark） */
  &.el-tag--dark {
    background: linear-gradient(135deg, rgba(14, 165, 233, 0.22), rgba(99, 102, 241, 0.18)) !important;
    border: 1px solid var(--border-glow, rgba(125, 211, 252, 0.26)) !important;
    color: var(--text-primary, #f8fafc) !important;
    box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.03);
  }

  /* 关闭按钮 */
  :deep(.el-tag__close) {
    color: var(--text-muted, rgba(148, 163, 184, 0.72)) !important;
    transition: all 0.2s ease;
  }

  :deep(.el-tag__close:hover) {
    background: rgba(239, 68, 68, 0.2) !important;
    color: #ef4444 !important;
  }
}

.circular {
  width: 7px;
  height: 7px;
  margin-right: 5px;
  background: var(--color-primary, #0ea5e9);
  border-radius: 50%;
  display: inline-block;
  box-shadow: 0 0 6px rgba(14, 165, 233, 0.5);
}
</style>
