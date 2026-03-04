<template>
  <div class="layout">
    <aside class="sidebar">OpsNexus</aside>
    <div class="main">
      <header class="header">顶栏区域</header>
      <section class="content">
        <RouterView />
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
// 应用主布局骨架。
</script>

<style scoped lang="scss">
.layout {
  display: flex;
  width: 100%;
  height: 100%;
}

.sidebar {
  width: 220px;
  background: rgba(15, 20, 40, 0.8);
  backdrop-filter: blur(8px);
  padding: 20px;
}

.main {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.header {
  height: 64px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  display: flex;
  align-items: center;
  padding: 0 20px;
}

.content {
  flex: 1;
  padding: 20px;
}
</style>
