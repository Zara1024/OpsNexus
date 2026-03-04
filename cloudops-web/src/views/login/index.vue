<template>
  <div class="login-page">
    <div class="login-card">
      <img class="logo" src="/logo.png" alt="OpsNexus Logo" />
      <h1>OpsNexus</h1>
      <p>智维云枢 · 智能运维管理平台</p>
      <el-form>
        <el-form-item>
          <el-input placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item>
          <el-input type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-button type="primary" class="w-full">登录</el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
// 登录页面仅提供界面骨架，不包含业务逻辑。
</script>

<style scoped lang="scss">
.login-page {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-card {
  width: 380px;
  padding: 28px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(12px);
}

.logo {
  width: 48px;
  height: 48px;
  margin-bottom: 10px;
}

h1 {
  margin-bottom: 8px;
}

p {
  margin-bottom: 20px;
  color: #b9c1e8;
}

.w-full {
  width: 100%;
}
</style>
