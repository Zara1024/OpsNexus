<template>
  <PageContainer class="global-search-workspace">
    <PageHeader
      eyebrow="Unified Discovery"
      title="全局搜索"
      subtitle="跨主机、集群、应用、告警、用户和菜单统一检索，避免在多个模块之间反复切换。"
    >
      <template #intro>
        <PageIntro
          title="下一步建议"
          text="先输入高辨识度关键字，再按类型收敛结果范围；若暂无结果，可直接切换到知识库或 AI 诊断继续联动。"
        />
      </template>
      <template #actions>
        <el-button icon="Refresh" @click="executeSearch">刷新</el-button>
      </template>
    </PageHeader>

    <PageToolbar>
      <el-form :inline="true" class="filter-cluster">
        <el-form-item class="filter-field global-search-workspace__keyword">
          <el-input
            v-model="queryForm.keyword"
            placeholder="输入关键字，例如 admin、opsnexus-k3s-e2e、告警"
            clearable
            @keyup.enter="submitSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
      <template #actions>
        <el-button type="primary" @click="submitSearch">搜索</el-button>
        <el-button @click="resetSearch">重置</el-button>
      </template>
    </PageToolbar>

    <SectionCard title="搜索范围" subtitle="选择要参与检索的能力域。">
      <div class="page-chip-row">
        <el-check-tag
          v-for="item in typeOptions"
          :key="item.value"
          :checked="queryForm.types.includes(item.value)"
          @change="toggleType(item.value)"
        >
          {{ item.label }}
        </el-check-tag>
      </div>
    </SectionCard>

    <StatStrip :items="summaryItems" />

    <el-alert
      v-if="errorMessage"
      :title="errorMessage"
      type="error"
      :closable="false"
      class="page-alert"
    />

    <EmptyState
      v-else-if="!loading && !hasSearched"
      title="还没有开始搜索"
      description="输入主机名、业务关键词、告警等级或用户账号后，平台会按能力域给出聚合结果。"
    >
      <template #actions>
        <el-button type="primary" @click="queryForm.keyword = 'admin'; submitSearch()">搜索示例词</el-button>
      </template>
    </EmptyState>

    <EmptyState
      v-else-if="!loading && hasSearched && result.total === 0"
      title="没有找到匹配结果"
      description="可以尝试缩短关键字、减少筛选类型，或者转到知识库初始化 SOP 与经验文档。"
    >
      <template #actions>
        <el-button @click="$router.push('/knowledge/base')">查看知识库</el-button>
        <el-button type="primary" @click="$router.push('/ai/diagnosis')">去 AI 诊断</el-button>
      </template>
    </EmptyState>

    <div v-loading="loading" class="group-list">
      <SectionCard
        v-for="group in result.groups"
        :key="group.type"
        :title="group.typeLabel"
        :subtitle="`共 ${group.count} 条结果`"
        dense
      >
        <div class="result-grid">
          <article
            v-for="item in group.items"
            :key="`${item.type}-${item.id}-${item.route}`"
            class="result-card"
            @click="openResult(item)"
          >
            <div class="result-card__top">
              <el-tag size="small" :type="tagType(item.type)">{{ item.typeLabel }}</el-tag>
              <span class="result-card__status">{{ item.status || '可查看' }}</span>
            </div>
            <div class="result-card__title">{{ item.title }}</div>
            <div class="result-card__subtitle">{{ item.subtitle || '-' }}</div>
            <div class="result-card__desc">{{ item.description || '暂无描述' }}</div>
            <div v-if="item.tags && item.tags.length" class="result-card__tags">
              <el-tag v-for="tag in item.tags" :key="tag" size="small" effect="plain">{{ tag }}</el-tag>
            </div>
            <div class="result-card__route">{{ item.route }}</div>
          </article>
        </div>
      </SectionCard>
    </div>
  </PageContainer>
</template>

<script>
import { ElMessage } from 'element-plus'
import EmptyState from '@/components/platform/EmptyState.vue'
import PageContainer from '@/components/platform/PageContainer.vue'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import SectionCard from '@/components/platform/SectionCard.vue'
import StatStrip from '@/components/platform/StatStrip.vue'

const defaultTypes = ['cluster', 'host', 'application', 'alert', 'user', 'menu']

const createEmptyResult = () => ({
  keyword: '',
  total: 0,
  groups: [],
  results: []
})

export default {
  name: 'GlobalSearch',
  components: {
    EmptyState,
    PageContainer,
    PageHeader,
    PageIntro,
    PageToolbar,
    SectionCard,
    StatStrip
  },
  data() {
    return {
      loading: false,
      hasSearched: false,
      errorMessage: '',
      queryForm: {
        keyword: '',
        types: [...defaultTypes]
      },
      typeOptions: [
        { value: 'cluster', label: '集群' },
        { value: 'host', label: '主机' },
        { value: 'application', label: '应用' },
        { value: 'alert', label: '告警' },
        { value: 'user', label: '用户' },
        { value: 'menu', label: '菜单' }
      ],
      result: createEmptyResult()
    }
  },
  computed: {
    summaryItems() {
      const groups = this.result.groups || []
      return [
        { label: '结果总数', value: this.result.total, hint: '当前查询命中的聚合结果', tone: 'primary' },
        ...groups.slice(0, 5).map(group => ({
          label: group.typeLabel,
          value: group.count,
          hint: `${group.typeLabel}结果`,
          tone: 'neutral'
        }))
      ]
    }
  },
  watch: {
    '$route.query.keyword': {
      immediate: true,
      handler() {
        this.syncFromRoute()
      }
    },
    '$route.query.types': {
      immediate: true,
      handler() {
        this.syncFromRoute()
      }
    }
  },
  methods: {
    syncFromRoute() {
      const keyword = typeof this.$route.query.keyword === 'string' ? this.$route.query.keyword : ''
      const routeTypes = typeof this.$route.query.types === 'string'
        ? this.$route.query.types.split(',').map(item => item.trim()).filter(Boolean)
        : []
      this.queryForm.keyword = keyword
      this.queryForm.types = routeTypes.length ? routeTypes : [...defaultTypes]
      if (keyword) {
        this.executeSearch()
      } else {
        this.result = createEmptyResult()
        this.hasSearched = false
        this.errorMessage = ''
      }
    },
    toggleType(value) {
      const exists = this.queryForm.types.includes(value)
      if (exists) {
        if (this.queryForm.types.length === 1) {
          ElMessage.warning('至少保留一个搜索类型')
          return
        }
        this.queryForm.types = this.queryForm.types.filter(item => item !== value)
        return
      }
      this.queryForm.types.push(value)
    },
    submitSearch() {
      const keyword = this.queryForm.keyword.trim()
      this.$router.push({
        path: '/search/global',
        query: {
          ...(keyword ? { keyword } : {}),
          ...(this.queryForm.types.length && this.queryForm.types.length !== defaultTypes.length
            ? { types: this.queryForm.types.join(',') }
            : {})
        }
      })
    },
    resetSearch() {
      this.queryForm.keyword = ''
      this.queryForm.types = [...defaultTypes]
      this.result = createEmptyResult()
      this.hasSearched = false
      this.errorMessage = ''
      this.$router.push({ path: '/search/global' })
    },
    async executeSearch() {
      const keyword = this.queryForm.keyword.trim()
      if (!keyword) return
      this.loading = true
      this.errorMessage = ''
      try {
        const { data: res } = await this.$api.globalSearch({
          keyword,
          types: this.queryForm.types.join(','),
          limit: 8
        })
        if (res.code !== 200) {
          throw new Error(res.message || '全局搜索失败')
        }
        this.result = res.data || createEmptyResult()
        this.hasSearched = true
      } catch (error) {
        this.errorMessage = error?.response?.data?.message || error.message || '全局搜索失败'
      } finally {
        this.loading = false
      }
    },
    openResult(item) {
      if (!item.route) return
      this.$router.push(item.route)
    },
    tagType(type) {
      return {
        cluster: 'primary',
        host: 'success',
        application: 'warning',
        alert: 'danger',
        user: 'info',
        menu: ''
      }[type] || 'info'
    }
  }
}
</script>

<style scoped>
.global-search-workspace__keyword {
  flex: 1 1 540px;
}

.group-list {
  display: grid;
  gap: 16px;
}

.result-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 14px;
}

.result-card {
  display: grid;
  gap: 10px;
  min-height: 220px;
  padding: 16px;
  border-radius: 18px;
  border: 1px solid var(--border-subtle);
  background: rgba(255, 255, 255, 0.03);
  transition: transform var(--transition-fast), border-color var(--transition-fast), background var(--transition-fast);
  cursor: pointer;
}

.result-card:hover {
  transform: translateY(-2px);
  border-color: var(--border-strong);
  background: rgba(53, 164, 255, 0.06);
}

.result-card__top,
.result-card__tags {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.result-card__status,
.result-card__route,
.result-card__subtitle,
.result-card__desc {
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-muted);
}

.result-card__title {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
}

@media (max-width: 768px) {
  .result-grid {
    grid-template-columns: 1fr;
  }
}
</style>
