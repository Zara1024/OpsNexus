<template>
  <PageContainer wide class="service-market-page">
    <PageHeader
      eyebrow="Tooling Workspace"
      title="运维工具箱"
      subtitle="把服务市场和部署管理收敛到同一工作区，统一查看可用服务、版本与部署动作。"
    >
      <template #intro>
        <PageIntro
          title="使用方式"
          text="先在服务市场挑选标准能力，再进入部署管理跟踪安装、升级和运维动作，避免工具入口割裂。"
        />
      </template>
    </PageHeader>

    <PageToolbar>
      <div class="page-chip-row">
        <el-check-tag :checked="innerTab === 'market'" @change="innerTab = 'market'">服务市场</el-check-tag>
        <el-check-tag :checked="innerTab === 'deploy'" @change="innerTab = 'deploy'">部署管理</el-check-tag>
      </div>
      <template #actions>
        <div v-if="innerTab === 'market'" class="page-chip-row">
          <el-check-tag :checked="selectedCategory === ''" @change="selectedCategory = ''">全部</el-check-tag>
          <el-check-tag
            v-for="cat in categories"
            :key="cat.id"
            :checked="selectedCategory === cat.id"
            @change="selectedCategory = cat.id"
          >
            {{ cat.name }}
          </el-check-tag>
        </div>
      </template>
    </PageToolbar>

    <StatStrip :items="statItems" />

    <SectionCard v-if="innerTab === 'market'" title="服务市场" subtitle="按技术域浏览可部署服务，并快速发起部署。">
      <div v-loading="loading" class="service-grid">
        <button
          v-for="service in displayServices"
          :key="service.id"
          class="service-card"
          @click="openDeployDialog(service)"
        >
          <div class="service-icon">
            <img :src="getServiceIcon(service.id)" :alt="service.name" class="service-svg-icon" />
          </div>
          <div class="service-name">{{ service.name }}</div>
          <div class="service-description">{{ service.description }}</div>
          <div class="service-versions">
            <el-tag
              v-for="version in service.versions.slice(0, 2)"
              :key="version.id"
              :type="version.recommended ? 'success' : 'info'"
              size="small"
            >
              {{ version.name }}
            </el-tag>
          </div>
        </button>
      </div>

      <EmptyState
        v-if="!loading && !displayServices.length"
        title="当前分类没有可用服务"
        description="可以切换其他分类查看，也可以直接进入部署管理核对已有安装记录。"
      >
        <template #actions>
          <el-button @click="selectedCategory = ''">查看全部服务</el-button>
          <el-button type="primary" @click="innerTab = 'deploy'">进入部署管理</el-button>
        </template>
      </EmptyState>
    </SectionCard>

    <SectionCard v-else title="部署管理" subtitle="统一查看安装记录、运行状态与后续维护动作。">
      <DeployManage :embedded="true" />
    </SectionCard>

    <DeployDialog
      v-model="deployDialogVisible"
      :service="selectedService"
      @deployed="handleDeployed"
    />
  </PageContainer>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getServicesList } from '@/api/tool'
import DeployDialog from './DeployDialog.vue'
import DeployManage from './DeployManage.vue'
import EmptyState from '@/components/platform/EmptyState.vue'
import PageContainer from '@/components/platform/PageContainer.vue'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import SectionCard from '@/components/platform/SectionCard.vue'
import StatStrip from '@/components/platform/StatStrip.vue'

const services = ref([])
const categories = ref([])
const selectedCategory = ref('')
const deployDialogVisible = ref(false)
const selectedService = ref(null)
const loading = ref(false)
const innerTab = ref('market')

const getServiceIcon = (serviceId) => {
  try {
    const fileNameMap = {
      mysql: 'mysql',
      redis: 'redis',
      postgresql: 'PostgreSQL',
      jenkins: 'Jenkins',
      gitlab: 'gitlab',
      grafana: 'grafana',
      elasticsearch: 'Elasticsearch',
      loki: 'loki',
      prometheus: 'Prometheus',
      elk: 'ELK',
      n9e: 'n9e',
      jumpserver: 'jumpserver',
      nodejs: 'nodejs',
      java: 'java',
      golang: 'golang',
      mongodb: 'mongodb',
      fluentd: 'fluentd'
    }
    const fileName = fileNameMap[String(serviceId || '').toLowerCase()]
    if (fileName) {
      return require(`@/assets/image/${fileName}.svg`)
    }
  } catch (error) {
    console.warn(`加载图标失败: ${serviceId}`, error)
  }
  return require('@/assets/image/云主机服务器.svg')
}

const filteredServices = computed(() => {
  if (!selectedCategory.value) return services.value
  return services.value.filter(service => service.category === selectedCategory.value)
})

const displayServices = computed(() => {
  if (!selectedCategory.value) return filteredServices.value
  return filteredServices.value.slice(0, 10)
})

const statItems = computed(() => [
  { label: '服务总数', value: services.value.length, hint: '当前可浏览的服务能力', tone: 'primary' },
  { label: '分类数量', value: categories.value.length, hint: '按技术域组织', tone: 'success' },
  { label: '当前分类结果', value: filteredServices.value.length, hint: '适合当前筛选条件', tone: 'warning' },
  { label: '推荐版本', value: services.value.reduce((count, service) => count + service.versions.filter(item => item.recommended).length, 0), hint: '默认优先展示的版本', tone: 'neutral' }
])

const loadServices = async () => {
  loading.value = true
  try {
    const res = await getServicesList()
    if (res.data?.code === 200) {
      services.value = res.data.data?.services || []
      categories.value = res.data.data?.categories || []
    } else {
      throw new Error(res.data?.message || '获取服务列表失败')
    }
  } catch (error) {
    console.error('加载服务列表失败:', error)
    ElMessage.error(`加载服务列表失败: ${error.message}`)
  } finally {
    loading.value = false
  }
}

const openDeployDialog = (service) => {
  selectedService.value = service
  deployDialogVisible.value = true
}

const handleDeployed = () => {
  ElMessage.success('部署任务已创建')
}

onMounted(() => {
  loadServices()
})
</script>

<style scoped>
.service-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.service-card {
  display: grid;
  gap: 12px;
  text-align: left;
  padding: 18px;
  border-radius: 18px;
  border: 1px solid var(--border-subtle);
  background: rgba(255, 255, 255, 0.03);
  transition: transform var(--transition-fast), border-color var(--transition-fast), background var(--transition-fast);
  cursor: pointer;
}

.service-card:hover {
  transform: translateY(-2px);
  border-color: var(--border-strong);
  background: rgba(53, 164, 255, 0.06);
}

.service-icon {
  width: 56px;
  height: 56px;
  display: grid;
  place-items: center;
  border-radius: 16px;
  background: rgba(53, 164, 255, 0.12);
}

.service-svg-icon {
  width: 38px;
  height: 38px;
  object-fit: contain;
}

.service-name {
  font-size: 17px;
  font-weight: 700;
  color: var(--text-primary);
}

.service-description {
  min-height: 42px;
  font-size: 13px;
  line-height: 1.65;
  color: var(--text-muted);
}

.service-versions {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
</style>
