<template>
  <div class="ops-shell">
    <div
      v-if="isMobile && mobileNavVisible"
      class="ops-shell__overlay"
      @click="mobileNavVisible = false"
    />

    <aside
      :class="[
        'ops-shell__aside',
        { 'is-collapsed': isCollapse, 'is-mobile-open': mobileNavVisible }
      ]"
      >
      <div class="ops-shell__brand">
        <div class="ops-shell__brand-icon">
          <img :src="brand.logo" :alt="brand.name" class="ops-shell__brand-logo" />
        </div>
      </div>

      <div v-show="!isCollapse || isMobile" class="ops-shell__nav-summary">
        <div class="ops-shell__nav-label">当前工作区</div>
        <div class="ops-shell__nav-title">{{ currentSection }}</div>
        <div class="ops-shell__nav-hint">{{ currentPageSubtitle }}</div>
      </div>

      <div class="ops-shell__menu-wrap">
        <el-menu
          router
          unique-opened
          background-color="transparent"
          :default-active="$route.path"
          :collapse="isCollapse && !isMobile"
          :collapse-transition="false"
          class="ops-shell__menu"
        >
          <el-menu-item
            v-for="item in leafMenus"
            :key="menuKey(item)"
            :index="menuPath(item.url)"
            @click="handleMenuClick(item.url)"
          >
            <el-icon><component :is="item.icon || 'Grid'" /></el-icon>
            <template #title>
              <span>{{ item.menuName }}</span>
            </template>
          </el-menu-item>

          <el-sub-menu
            v-for="item in groupedMenus"
            :key="menuKey(item)"
            :index="String(item.id || item.menuName)"
          >
            <template #title>
              <el-icon><component :is="item.icon || 'Grid'" /></el-icon>
              <span>{{ item.menuName }}</span>
            </template>
            <el-menu-item
              v-for="subItem in item.menuSvoList"
              :key="menuKey(subItem)"
              :index="menuPath(subItem.url)"
              @click="handleMenuClick(subItem.url)"
            >
              <el-icon><component :is="subItem.icon || 'Grid'" /></el-icon>
              <template #title>
                <span>{{ subItem.menuName }}</span>
              </template>
            </el-menu-item>
          </el-sub-menu>
        </el-menu>
      </div>
    </aside>

    <section class="ops-shell__content">
      <header class="ops-shell__topbar">
        <div class="ops-shell__topbar-main">
          <el-button link class="ops-shell__toggle" @click="toggleCollapse">
            <el-icon size="20"><component :is="toggleIcon" /></el-icon>
          </el-button>

          <div class="ops-shell__route">
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/dashboard' }">仪表盘</el-breadcrumb-item>
              <el-breadcrumb-item v-if="$route.meta?.sTitle">{{ $route.meta.sTitle }}</el-breadcrumb-item>
              <el-breadcrumb-item v-if="$route.meta?.tTitle">{{ $route.meta.tTitle }}</el-breadcrumb-item>
            </el-breadcrumb>
            <div class="ops-shell__route-copy">
              <div class="ops-shell__route-title">{{ currentPageTitle }}</div>
              <div class="ops-shell__route-subtitle">{{ currentPageSubtitle }}</div>
            </div>
          </div>
        </div>

        <div class="ops-shell__topbar-tools">
          <div class="ops-shell__search">
            <el-input
              v-model="globalSearchKeyword"
              clearable
              placeholder="搜索主机、集群、应用、告警、AI 会话、知识"
              @keyup.enter="goGlobalSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-button type="primary" @click="goGlobalSearch">全局搜索</el-button>
          </div>
          <HeadImage />
        </div>
      </header>

      <Tags class="ops-shell__tags" />

      <main class="ops-shell__main">
        <div :class="['ops-shell__page-shell', routeShellClass]">
          <router-view />
        </div>
      </main>
    </section>
  </div>
</template>

<script>
import storage from '@/utils/storage'
import HeadImage from '@/components/HeadImage.vue'
import Tags from '@/components/Tags.vue'
import { BRANDING } from '@/constants/branding'
import { normalizePlatformMenu } from '@/utils/platformMenu'

export default {
  name: 'OpsHome',
  components: {
    HeadImage,
    Tags
  },
  data() {
    return {
      brand: BRANDING,
      leftMenuList: [],
      isCollapse: false,
      isMobile: false,
      mobileNavVisible: false,
      globalSearchKeyword: ''
    }
  },
  computed: {
    leafMenus() {
      return this.leftMenuList.filter(item => !Array.isArray(item.menuSvoList) || item.menuSvoList.length === 0)
    },
    groupedMenus() {
      return this.leftMenuList.filter(item => Array.isArray(item.menuSvoList) && item.menuSvoList.length > 0)
    },
    toggleIcon() {
      if (this.isMobile) {
        return this.mobileNavVisible ? 'Close' : 'Menu'
      }
      return this.isCollapse ? 'Expand' : 'Fold'
    },
    currentSection() {
      return this.$route.meta?.sTitle || '平台工作台'
    },
    currentPageTitle() {
      return this.$route.meta?.tTitle || '仪表盘'
    },
    currentPageSubtitle() {
      const moduleHints = {
        '仪表盘': '查看全局态势、风险与高频入口',
        '资产管理': '统一管理主机、网络设备、数据库与资源分组',
        '容器管理': '管理集群、节点与工作负载',
        '服务管理': '应用、发布与依赖关系集中治理',
        '任务中心': '统一承载脚本、流水线和作业执行',
        '运维工具': '工具市场、Agent 与自动化辅助能力',
        'AI 智能运维助手': '支持多模型接入、Agent 协作、知识检索与诊断巡检的统一工作台',
        '监控告警': '聚焦告警中心、推送链路、历史回溯与监控深化',
        '操作审计': '集中回看登录、操作、数据与终端审计记录',
        '系统管理': '统一维护组织、账号、角色、菜单与基础权限配置',
        '配置中心': '集中维护凭据、LDAP 和通用配置',
        '全局搜索': '跨模块检索资产、告警、用户和知识',
        '工单中心': '发布、变更、SQL 与协作工单总览',
        '知识库': '沉淀 SOP、FAQ、巡检和复盘经验'
      }
      return moduleHints[this.currentSection] || this.brand.description
    },
    routeShellClass() {
      const normalized = String(this.$route.path || '/dashboard')
        .replace(/^\/+|\/+$/g, '')
        .replace(/[^a-zA-Z0-9]+/g, '-')
        .replace(/^-+|-+$/g, '')
      return `route-shell--${normalized || 'dashboard'}`
    }
  },
  watch: {
    '$route.query.keyword': {
      immediate: true,
      handler(value) {
        this.globalSearchKeyword = typeof value === 'string' ? value : ''
      }
    },
    '$route.path'() {
      if (this.isMobile) {
        this.mobileNavVisible = false
      }
    }
  },
  methods: {
    loadMenuData() {
      const menuData = storage.getItem('leftMenuList')
      if (Array.isArray(menuData)) {
        this.leftMenuList = normalizePlatformMenu(menuData)
        storage.setItem('leftMenuList', this.leftMenuList)
        return
      }
      this.leftMenuList = normalizePlatformMenu([])
    },
    menuPath(url = '') {
      return `/${String(url || '').replace(/^\/+/, '')}`
    },
    menuKey(item) {
      return item.id || item.url || item.menuName
    },
    handleMenuClick(url) {
      storage.setItem('activePath', this.menuPath(url))
      if (this.isMobile) {
        this.mobileNavVisible = false
      }
    },
    toggleCollapse() {
      if (this.isMobile) {
        this.mobileNavVisible = !this.mobileNavVisible
        return
      }
      this.isCollapse = !this.isCollapse
    },
    goGlobalSearch() {
      const keyword = this.globalSearchKeyword.trim()
      this.$router.push({
        path: '/search/global',
        query: keyword ? { keyword } : {}
      })
    },
    syncViewport() {
      const mobile = window.innerWidth <= 960
      this.isMobile = mobile
      if (mobile) {
        this.mobileNavVisible = false
      }
      if (!mobile && this.isCollapse && window.innerWidth > 1400) {
        this.isCollapse = false
      }
    }
  },
  mounted() {
    this.loadMenuData()
    this.syncViewport()
    window.addEventListener('resize', this.syncViewport)
  },
  beforeUnmount() {
    window.removeEventListener('resize', this.syncViewport)
  }
}
</script>

<style scoped lang="less">
.ops-shell {
  display: flex;
  height: 100%;
  min-height: 100%;
  overflow: hidden;
  background: transparent;

  &__overlay {
    position: fixed;
    inset: 0;
    z-index: 24;
    background: rgba(3, 9, 18, 0.58);
    backdrop-filter: blur(4px);
  }

  &__aside {
    position: relative;
    z-index: 26;
    width: 276px;
    flex-shrink: 0;
    display: flex;
    flex-direction: column;
    gap: 16px;
    padding: 18px 14px 14px;
    background: linear-gradient(180deg, rgba(6, 14, 26, 0.98), rgba(4, 9, 18, 0.98));
    border-right: 1px solid var(--border-subtle);
    transition: width var(--transition-normal), transform var(--transition-normal);
    box-shadow: 20px 0 42px rgba(2, 8, 23, 0.26);
  }

  &__aside.is-collapsed {
    width: 92px;
  }

  &__brand {
    display: flex;
    justify-content: center;
    padding: 6px 0 2px;
  }

  &__brand-icon {
    width: 72px;
    height: 72px;
    display: grid;
    place-items: center;
    border-radius: 22px;
    background:
      linear-gradient(180deg, rgba(7, 18, 36, 0.96), rgba(11, 23, 43, 0.92)),
      radial-gradient(circle at 50% 32%, rgba(125, 211, 252, 0.18), transparent 72%);
    border: 1px solid rgba(125, 211, 252, 0.26);
    box-shadow:
      inset 0 0 0 1px rgba(255, 255, 255, 0.04),
      0 14px 28px rgba(3, 8, 20, 0.34);
    flex-shrink: 0;
  }

  &__brand-logo {
    width: 46px;
    height: 46px;
    object-fit: contain;
  }

  &__nav-summary {
    display: grid;
    gap: 4px;
    padding: 0 8px;
  }

  &__nav-label {
    font-size: 11px;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: var(--text-subtle);
  }

  &__nav-title {
    font-size: 15px;
    font-weight: 700;
    color: var(--text-primary);
  }

  &__nav-hint {
    font-size: 12px;
    line-height: 1.55;
    color: var(--text-muted);
  }

  &__menu-wrap {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    padding-right: 4px;
  }

  &__menu {
    border-right: none !important;
    background: transparent !important;
    padding: 4px 0;
  }

  &__content {
    flex: 1;
    min-width: 0;
    min-height: 0;
    height: 100%;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  &__topbar {
    position: sticky;
    top: 0;
    z-index: 20;
    display: flex;
    justify-content: space-between;
    gap: 18px;
    align-items: center;
    min-height: 84px;
    padding: 16px 22px;
    background: rgba(7, 15, 27, 0.84);
    backdrop-filter: blur(18px);
    border-bottom: 1px solid var(--border-subtle);
  }

  &__topbar-main,
  &__topbar-tools {
    display: flex;
    align-items: center;
    gap: 16px;
    min-width: 0;
  }

  &__topbar-main {
    flex: 1;
  }

  &__topbar-tools {
    justify-content: flex-end;
  }

  &__toggle {
    width: 40px;
    height: 40px;
    display: grid;
    place-items: center;
    border-radius: 12px;
    color: rgba(191, 219, 254, 0.96);
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid var(--border-subtle);
  }

  &__route {
    min-width: 0;
    display: grid;
    gap: 8px;
  }

  &__route-copy {
    display: grid;
    gap: 2px;
  }

  &__route-title {
    font-size: 20px;
    font-weight: 700;
    color: var(--text-primary);
  }

  &__route-subtitle {
    font-size: 13px;
    color: var(--text-muted);
  }

  &__search {
    display: flex;
    align-items: center;
    gap: 10px;
    min-width: min(460px, 48vw);
  }

  &__tags {
    position: sticky;
    top: 84px;
    z-index: 18;
  }

  &__main {
    flex: 1;
    min-height: 0;
    overflow-x: hidden;
    overflow-y: auto;
    overscroll-behavior: contain;
  }

  &__page-shell {
    min-height: 100%;
    padding: 22px 22px 28px;
  }

  &__page-shell > :not(.platform-page) {
    width: 100%;
    max-width: var(--content-width-ultra);
    margin: 0 auto;
  }
}

:deep(.ops-shell__menu .el-menu-item),
:deep(.ops-shell__menu .el-sub-menu__title) {
  height: 46px !important;
  line-height: 46px !important;
  border-radius: 14px;
  margin: 0 6px 6px;
  color: rgba(226, 232, 240, 0.84) !important;
  transition: background var(--transition-fast), transform var(--transition-fast), color var(--transition-fast);
}

:deep(.ops-shell__menu .el-menu-item:hover),
:deep(.ops-shell__menu .el-sub-menu__title:hover) {
  background: rgba(255, 255, 255, 0.06) !important;
  color: var(--text-primary) !important;
}

:deep(.ops-shell__menu .el-sub-menu .el-menu-item) {
  margin-left: 14px;
  padding-left: 50px !important;
}

:deep(.ops-shell__menu > .el-menu-item.is-active),
:deep(.ops-shell__menu .el-sub-menu .el-menu-item.is-active) {
  color: #fff !important;
  background: linear-gradient(135deg, rgba(37, 99, 235, 0.96), rgba(59, 130, 246, 0.92)) !important;
  box-shadow: var(--shadow-glow);
}

:deep(.ops-shell__menu .el-sub-menu__title .el-icon),
:deep(.ops-shell__menu .el-menu-item .el-icon) {
  width: 18px;
  font-size: 18px;
  margin-right: 10px;
}

:deep(.ops-shell__menu .el-sub-menu__icon-arrow) {
  color: rgba(148, 163, 184, 0.8) !important;
}

@media (max-width: 1200px) {
  .ops-shell {
    &__topbar {
      padding: 14px 18px;
    }

    &__search {
      min-width: min(360px, 42vw);
    }
  }
}

@media (max-width: 960px) {
  .ops-shell {
    &__aside {
      position: fixed;
      inset: 0 auto 0 0;
      transform: translateX(-100%);
      width: 304px;
      max-width: calc(100vw - 56px);
    }

    &__aside.is-mobile-open {
      transform: translateX(0);
    }

    &__topbar {
      min-height: 72px;
      padding: 14px 16px;
    }

    &__topbar,
    &__topbar-main,
    &__topbar-tools {
      flex-wrap: wrap;
    }

    &__search {
      width: 100%;
      min-width: 0;
    }

    &__tags {
      top: 72px;
    }

    &__page-shell {
      padding: 16px 16px 22px;
    }
  }
}

@media (max-width: 640px) {
  .ops-shell {
    &__route-title {
      font-size: 17px;
    }

    &__route-subtitle {
      font-size: 12px;
    }
  }
}
</style>

<style lang="less">
.el-menu--popup,
.el-menu--popup-container {
  background: linear-gradient(180deg, rgba(10, 20, 35, 0.98), rgba(8, 15, 27, 0.98)) !important;
  border: 1px solid var(--border-medium) !important;
  border-radius: 16px !important;
  box-shadow: var(--shadow-elevated) !important;
  padding: 6px !important;
}

.el-menu--popup .el-menu-item {
  border-radius: 12px !important;
  margin: 4px 0 !important;
}
</style>
