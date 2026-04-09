<template>
  <div class="login-container">
    <canvas ref="snowCanvas" class="snow-canvas"></canvas>
    <div class="scene-grid"></div>
    <div class="scene-orb scene-orb--cyan"></div>
    <div class="scene-orb scene-orb--blue"></div>

    <div class="login-shell">
      <section class="brand-panel">
        <div class="brand-hero">
          <div class="brand-badge">AIOps Control Center</div>
          <div class="brand-row">
            <div class="brand-logo-wrap">
              <img :src="brand.loginLogo || brand.logo" :alt="brand.name" class="brand-logo">
            </div>
            <div class="brand-copy">
              <h1 class="brand-title">{{ brand.name }}</h1>
              <div class="brand-slogan">{{ brand.slogan }}</div>
            </div>
          </div>
        </div>
        <p class="brand-description">{{ brand.loginDescription || brand.description }}</p>

        <div class="brand-highlights">
          <div v-for="highlight in brand.loginHighlights || []" :key="highlight.name" class="highlight-card">
            <div class="highlight-name">{{ highlight.name }}</div>
            <div class="highlight-desc">{{ highlight.description }}</div>
            <div v-if="highlight.points && highlight.points.length" class="highlight-points">
              <span v-for="point in highlight.points" :key="point" class="highlight-point">{{ point }}</span>
            </div>
          </div>
        </div>
      </section>

      <section class="login-card">
        <div class="card-header">
          <div class="card-eyebrow">Secure Access</div>
          <h2 class="card-title">登录控制台</h2>
          <div class="card-subtitle">进入 {{ brand.name }} 运维工作台</div>
        </div>

        <el-form ref="loginFormRef" :rules="rules" :model="loginForm" class="login-form">
          <el-form-item prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="请输入账号"
              clearable
              class="dark-input"
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              placeholder="请输入密码"
              type="password"
              show-password
              clearable
              class="dark-input"
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item prop="image">
            <div class="captcha-row">
              <el-input
                v-model="loginForm.image"
                placeholder="请输入验证码"
                maxlength="6"
                clearable
                class="dark-input"
                @keyup.enter="loginBtn"
              >
                <template #prefix>
                  <el-icon><CircleCheck /></el-icon>
                </template>
              </el-input>
              <button type="button" class="captcha-box" title="点击刷新验证码" @click="refreshCaptcha">
                <img :src="image" class="captcha-img" alt="验证码">
                <span class="captcha-tip">点击刷新</span>
              </button>
            </div>
          </el-form-item>

          <div class="captcha-helper">验证码已优化为高对比度样式，看不清时可随时刷新。</div>

          <el-form-item class="action-item">
            <el-button class="login-btn" type="primary" @click="loginBtn">登 录</el-button>
            <el-button class="reset-btn" @click="resetLoginForm">重 置</el-button>
          </el-form-item>
        </el-form>
      </section>
    </div>

    <div class="login-footer">
      <span>© {{ brand.loginFooterYear || '2026' }} {{ brand.name }} · 智能运维控制台</span>
    </div>
  </div>
</template>

<script>
import { BRANDING } from '@/constants/branding'

export default {
  name: 'UserLogin',
  data() {
    return {
      brand: BRANDING,
      image: '',
      snowAnimId: null,
      flakes: [],
      rules: {
        username: [{ required: true, message: '请输入账号', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
        image: [{ required: true, message: '请输入验证码', trigger: 'blur' }]
      },
      loginForm: {
        username: '',
        password: '',
        image: '',
        idKey: ''
      }
    }
  },
  methods: {
    async getCaptcha() {
      const { data: res } = await this.$api.captcha()
      if (res.code !== 200) {
        this.$message.error(res.message)
        return
      }
      this.image = res.data.image
      this.loginForm.idKey = res.data.idKey
    },
    refreshCaptcha() {
      this.loginForm.image = ''
      this.getCaptcha()
    },
    resizeSnowCanvas() {
      const canvas = this.$refs.snowCanvas
      if (!canvas) return
      canvas.width = window.innerWidth
      canvas.height = window.innerHeight
    },
    initSnow() {
      const canvas = this.$refs.snowCanvas
      if (!canvas) return

      const ctx = canvas.getContext('2d')
      this.resizeSnowCanvas()
      this.flakes = Array.from({ length: 160 }, () => ({
        x: Math.random() * canvas.width,
        y: Math.random() * canvas.height,
        r: Math.random() * 2.8 + 0.6,
        speed: Math.random() * 0.9 + 0.2,
        opacity: Math.random() * 0.55 + 0.18,
        drift: (Math.random() - 0.5) * 0.35
      }))

      const draw = () => {
        ctx.clearRect(0, 0, canvas.width, canvas.height)
        this.flakes.forEach((flake) => {
          ctx.beginPath()
          ctx.arc(flake.x, flake.y, flake.r, 0, Math.PI * 2)
          ctx.fillStyle = `rgba(255,255,255,${flake.opacity})`
          ctx.fill()

          flake.y += flake.speed
          flake.x += flake.drift

          if (flake.y > canvas.height + 8) {
            flake.y = -10
            flake.x = Math.random() * canvas.width
          }
          if (flake.x > canvas.width + 8) {
            flake.x = -8
          }
          if (flake.x < -8) {
            flake.x = canvas.width + 8
          }
        })

        this.snowAnimId = requestAnimationFrame(draw)
      }

      draw()
    },
    loginBtn() {
      this.$refs.loginFormRef.validate(async (valid) => {
        if (!valid) {
          return false
        }

        const { data: res } = await this.$api.login(this.loginForm)
        if (res.code !== 200) {
          this.$message.error(res.message)
          this.refreshCaptcha()
          return
        }

        this.$message.success('登录成功')
        this.$store.commit('saveSysAdmin', res.data.sysAdmin)
        this.$store.commit('saveToken', res.data.token)
        this.$store.commit('saveLeftMenuList', res.data.leftMenuList)
        this.$store.commit('savePermissionList', res.data.permissionList)
        await this.$router.push('/home')
      })
    },
    resetLoginForm() {
      this.$refs.loginFormRef.resetFields()
      this.refreshCaptcha()
    },
    handleResize() {
      this.resizeSnowCanvas()
    }
  },
  async mounted() {
    await this.getCaptcha()
    this.initSnow()
    window.addEventListener('resize', this.handleResize)
  },
  beforeUnmount() {
    cancelAnimationFrame(this.snowAnimId)
    window.removeEventListener('resize', this.handleResize)
  }
}
</script>

<style lang="less" scoped>
.login-container {
  --panel-padding: 34px;
  --panel-radius: 28px;
  --field-height: 48px;
  --panel-border: rgba(148, 163, 184, 0.16);
  --panel-bg: rgba(6, 15, 29, 0.74);
  --text-secondary: #94a3b8;
  position: relative;
  min-height: 100vh;
  display: grid;
  place-items: center;
  overflow: hidden;
  background:
    radial-gradient(circle at top left, rgba(56, 189, 248, 0.2), transparent 28%),
    radial-gradient(circle at right, rgba(99, 102, 241, 0.22), transparent 30%),
    linear-gradient(145deg, #040b17 0%, #081423 45%, #040914 100%);
}

.scene-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(125, 211, 252, 0.08) 1px, transparent 1px),
    linear-gradient(90deg, rgba(125, 211, 252, 0.08) 1px, transparent 1px);
  background-size: 72px 72px;
  mask-image: linear-gradient(180deg, rgba(255, 255, 255, 0.4), transparent 88%);
  pointer-events: none;
}

.scene-orb {
  position: absolute;
  width: 520px;
  height: 520px;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.45;
}

.scene-orb--cyan {
  top: -120px;
  left: -120px;
  background: rgba(14, 165, 233, 0.42);
}

.scene-orb--blue {
  right: -120px;
  bottom: -120px;
  background: rgba(99, 102, 241, 0.38);
}

.snow-canvas {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.login-shell {
  position: relative;
  z-index: 1;
  width: min(1180px, calc(100% - 48px));
  display: grid;
  grid-template-columns: minmax(0, 1.05fr) 420px;
  gap: 28px;
  align-items: stretch;
}

.brand-panel,
.login-card {
  border: 1px solid var(--panel-border);
  background: linear-gradient(180deg, rgba(8, 18, 34, 0.86), var(--panel-bg));
  backdrop-filter: blur(24px);
  box-shadow: 0 30px 80px rgba(2, 8, 23, 0.34);
}

.brand-panel {
  border-radius: var(--panel-radius);
  padding: var(--panel-padding);
  display: grid;
  grid-template-rows: auto auto 1fr;
  align-content: start;
  gap: 24px;
}

.brand-badge,
.card-eyebrow {
  display: inline-flex;
  align-items: center;
  align-self: start;
  width: fit-content;
  padding: 8px 14px;
  border-radius: 999px;
  background: rgba(14, 165, 233, 0.12);
  border: 1px solid rgba(125, 211, 252, 0.16);
  color: rgba(125, 211, 252, 0.92);
  font-size: 12px;
  line-height: 1;
  font-weight: 600;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.brand-hero,
.card-header {
  display: grid;
  align-content: start;
  gap: 12px;
}

.brand-hero {
  min-height: 126px;
}

.brand-row {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  align-items: start;
  gap: 18px;
  min-height: 80px;
}

.brand-logo-wrap {
  position: relative;
  width: 86px;
  height: 86px;
  border-radius: 28px;
  display: grid;
  place-items: center;
  overflow: hidden;
  background: linear-gradient(135deg, rgba(10, 26, 50, 0.94), rgba(17, 58, 99, 0.76));
  border: 1px solid rgba(125, 211, 252, 0.24);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.05),
    0 18px 36px rgba(8, 47, 73, 0.22);
}

.brand-logo-wrap::before {
  content: '';
  position: absolute;
  inset: 10px;
  border-radius: 22px;
  background: radial-gradient(circle at top, rgba(96, 165, 250, 0.22), transparent 68%);
  pointer-events: none;
}

.brand-logo {
  position: relative;
  z-index: 1;
  width: 58px;
  height: 58px;
  object-fit: contain;
  filter: drop-shadow(0 12px 18px rgba(37, 99, 235, 0.2));
}

.brand-copy {
  min-width: 0;
}

.brand-title,
.card-title {
  margin: 0;
  color: #f8fafc;
  font-weight: 700;
}

.brand-title {
  font-size: 42px;
  line-height: 1.05;
  letter-spacing: -0.02em;
}

.brand-slogan,
.card-subtitle,
.highlight-desc,
.captcha-helper {
  color: var(--text-secondary);
}

.brand-slogan {
  margin-top: 8px;
  font-size: 16px;
}

.brand-description {
  margin: 0;
  max-width: 640px;
  color: rgba(226, 232, 240, 0.86);
  font-size: 15px;
  line-height: 1.9;
}

.brand-highlights {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.highlight-card {
  min-height: 188px;
  padding: 22px;
  border-radius: 22px;
  border: 1px solid rgba(125, 211, 252, 0.12);
  background: linear-gradient(180deg, rgba(15, 23, 42, 0.62), rgba(9, 16, 32, 0.46));
  display: grid;
  grid-template-rows: auto auto 1fr;
  align-content: start;
  gap: 10px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.03);
}

.highlight-name {
  color: #f8fafc;
  font-size: 17px;
  font-weight: 700;
  margin-bottom: 0;
}

.highlight-desc {
  font-size: 13px;
  line-height: 1.75;
}

.highlight-points {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-self: end;
}

.highlight-point {
  display: inline-flex;
  align-items: center;
  padding: 6px 9px;
  border-radius: 999px;
  border: 1px solid rgba(96, 165, 250, 0.16);
  background: rgba(15, 23, 42, 0.66);
  color: #cbd5e1;
  font-size: 11px;
  line-height: 1.2;
  white-space: nowrap;
}

.login-card {
  border-radius: var(--panel-radius);
  padding: var(--panel-padding);
  display: flex;
  flex-direction: column;
}

.card-header {
  min-height: 126px;
  margin-bottom: 28px;
}

.card-title {
  font-size: 32px;
  line-height: 1.1;
  letter-spacing: -0.02em;
}

.card-subtitle {
  margin: 0;
  font-size: 14px;
  line-height: 1.7;
}

.dark-input {
  :deep(.el-input__wrapper) {
    min-height: var(--field-height);
    height: var(--field-height);
    border-radius: 16px;
    background: var(--bg-input, rgba(2, 6, 23, 0.82));
    border: 1px solid var(--border-subtle, rgba(148, 163, 184, 0.16));
    box-shadow: none !important;
    transition: border-color 0.25s ease;
  }

  :deep(.el-input__wrapper:hover) {
    border-color: var(--border-hover, rgba(56, 189, 248, 0.36));
  }

  :deep(.el-input__wrapper.is-focus) {
    border-color: var(--color-primary, #0ea5e9);
    box-shadow: 0 0 0 3px rgba(14, 165, 233, 0.12) !important;
  }

  :deep(.el-input__inner) {
    height: 100%;
    color: var(--text-primary, #f8fafc);

    &::placeholder {
      color: var(--text-muted, rgba(148, 163, 184, 0.72));
    }
  }

  :deep(.el-input__prefix) {
    color: var(--text-muted, rgba(148, 163, 184, 0.72));
    font-size: 18px;
  }
}

.captcha-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 160px;
  gap: 12px;
  align-items: stretch;
}

.captcha-box {
  position: relative;
  display: block;
  border: 1px solid var(--border-glow, rgba(125, 211, 252, 0.18));
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.82), rgba(3, 9, 23, 0.94));
  border-radius: 16px;
  cursor: pointer;
  min-height: var(--field-height);
  height: var(--field-height);
  overflow: hidden;
  padding: 0;
  transition: transform 0.25s ease, border-color 0.25s ease, box-shadow 0.25s ease;

  &::before,
  &::after {
    content: '';
    position: absolute;
    inset: 0;
    pointer-events: none;
  }

  &::before {
    z-index: 1;
    background: linear-gradient(135deg, rgba(2, 6, 23, 0.16), rgba(2, 6, 23, 0.34));
  }

  &::after {
    z-index: 2;
    border-radius: inherit;
    box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.04);
  }

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 14px 28px rgba(8, 47, 73, 0.24);
    border-color: var(--border-hover, rgba(56, 189, 248, 0.42));
  }
}

.captcha-img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: fill;
  image-rendering: optimizeQuality;
  filter: brightness(0.82) contrast(1.18) saturate(0.85);
}

.captcha-tip {
  display: none;
  position: absolute;
  right: 8px;
  bottom: 6px;
  z-index: 3;
  color: #e2e8f0;
  font-size: 10px;
  line-height: 1;
  font-weight: 600;
  letter-spacing: 0.08em;
  padding: 4px 6px;
  border-radius: 999px;
  background: rgba(2, 6, 23, 0.58);
  backdrop-filter: blur(8px);
}

.captcha-helper {
  display: none;
  margin-top: 10px;
  margin-bottom: 24px;
  font-size: 12px;
  line-height: 1.75;
}

.action-item :deep(.el-form-item__content) {
  display: flex;
  gap: 12px;
  align-items: stretch;
}

.action-item :deep(.el-button + .el-button) {
  margin-left: 0;
}

.login-btn,
.reset-btn {
  height: var(--field-height);
  min-height: var(--field-height);
  border-radius: 16px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 0.12em;
  padding: 0 20px;
  transition: transform 0.25s ease, box-shadow 0.25s ease, border-color 0.25s ease;
}

.login-btn.el-button,
.reset-btn.el-button {
  border-radius: 16px !important;
  font-weight: 700 !important;
}

.login-btn {
  flex: 1;
  border: 1px solid rgba(96, 165, 250, 0.28);
  background: linear-gradient(135deg, #1d4ed8 0%, #0f7ae5 56%, #0891b2 100%) !important;
  box-shadow: 0 14px 28px rgba(14, 116, 214, 0.22);

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 18px 34px rgba(14, 116, 214, 0.24);
  }
}

.reset-btn {
  width: 110px;
  border: 1px solid rgba(148, 163, 184, 0.18) !important;
  background: rgba(15, 23, 42, 0.72) !important;
  color: #e2e8f0 !important;
}

:deep(.el-form-item) {
  margin-bottom: 18px;

  .el-form-item__error {
    color: #fda4af;
    padding-top: 6px;
  }
}

@media (max-width: 1040px) {
  .login-shell {
    grid-template-columns: 1fr;
  }

  .brand-hero,
  .card-header {
    min-height: auto;
  }

  .brand-highlights {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .login-container {
    padding: 12px;
  }

  .login-shell {
    width: 100%;
  }

  .brand-panel,
  .login-card {
    padding: 24px;
  }

  .brand-row {
    align-items: flex-start;
  }

  .brand-title {
    font-size: 32px;
  }

  .login-card {
    padding: 24px 20px;
  }

  .captcha-row {
    grid-template-columns: 1fr;
  }

  .captcha-box {
    width: 100%;
  }
}

.login-footer {
  position: absolute;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  color: var(--text-hint, rgba(148, 163, 184, 0.52));
  font-size: 12px;
  letter-spacing: 0.04em;
  white-space: nowrap;
}
</style>
