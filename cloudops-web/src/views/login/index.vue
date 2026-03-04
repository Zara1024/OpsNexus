<template>
  <div class="login-page">
    <div class="login-card">
      <img class="logo" src="/logo.png" alt="OpsNexus Logo" />
      <h1>OpsNexus</h1>
      <p>智维云枢 · 智能运维管理平台</p>

      <el-form :model="form" :rules="rules" ref="formRef" @keyup.enter="submit">
        <el-form-item prop="username"><el-input v-model="form.username" placeholder="请输入用户名" /></el-form-item>
        <el-form-item prop="password"><el-input v-model="form.password" show-password placeholder="请输入密码" /></el-form-item>
        <el-form-item prop="captcha">
          <div class="captcha-row">
            <el-input v-model="form.captcha" placeholder="请输入验证码" />
            <el-button @click="loadCaptcha">{{ captchaText || '获取验证码' }}</el-button>
          </div>
        </el-form-item>
        <el-button :loading="loading" type="primary" class="w-full" @click="submit">登录</el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { getCaptchaApi } from '@/api/auth'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)
const captchaText = ref('')
const captchaId = ref('')

const form = reactive({ username: 'admin', password: 'OpsNexus@2026', captcha: '' })

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  captcha: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
}

const loadCaptcha = async () => {
  const res = await getCaptchaApi()
  captchaText.value = res.captcha
  captchaId.value = res.captcha_id
}

const submit = async () => {
  await formRef.value?.validate()
  loading.value = true
  try {
    await userStore.login({
      username: form.username,
      password: form.password,
      captcha_id: captchaId.value,
      captcha: form.captcha,
    })
    ElMessage.success('登录成功')
    router.replace('/dashboard')
  } catch (_err) {
    await loadCaptcha()
  } finally {
    loading.value = false
  }
}

onMounted(loadCaptcha)
</script>

<style scoped lang="scss">
.login-page { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; background: linear-gradient(135deg, #1a1d4d, #667eea, #764ba2); }
.login-card { width: 400px; padding: 28px; border-radius: 12px; background: rgba(255, 255, 255, 0.12); backdrop-filter: blur(12px); }
.logo { width: 48px; height: 48px; margin-bottom: 10px; }
h1 { margin-bottom: 8px; }
p { margin-bottom: 20px; color: #d8defc; }
.w-full { width: 100%; }
.captcha-row { display: flex; gap: 8px; width: 100%; }
</style>
