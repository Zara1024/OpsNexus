<script>
export default {
  name: 'DashboardIndex'
}
</script>

<script setup>
import { computed, ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import { getDashboardStats, getBusinessDistribution } from '@/api/dashboard'
import { GetAllTools, CreateTool, UpdateTool, DeleteTool as DeleteToolAPI, UploadIcon } from '@/api/tool'
import { ElMessage, ElMessageBox } from 'element-plus'
import { BRANDING } from '@/constants/branding'
import { normalizeDashboardToolIconUrl } from '@/utils/dashboardToolIcon.mjs'
import EmptyState from '@/components/platform/EmptyState.vue'
import PageContainer from '@/components/platform/PageContainer.vue'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import SectionCard from '@/components/platform/SectionCard.vue'
import StatStrip from '@/components/platform/StatStrip.vue'

const router = useRouter()
const brand = BRANDING

// 响应式数据
const loading = ref(true)
const editDialogVisible = ref(false)
const editingTool = ref(null)
const editingIndex = ref(-1)

// 统计数据
const stats = reactive({
  assets: {
    title: '资产详情',
    items: [
      { label: '主机总数', value: 0 },
      { label: '数据库总数', value: 0 },
      { label: 'K8s集群数量', value: 0 }
    ]
  },
  services: {
    title: '服务详情',
    items: [
      { label: '应用总数', value: 0 },
      { label: '业务线总数', value: 0 }
    ]
  },
  deployment: {
    title: '发布详情',
    items: [
      { label: '应用发布', value: 0 },
      { label: '任务执行', value: 0 },
      { label: '成功率', value: 0, unit: '%' }
    ]
  },
  monitor: {
    title: '监控告警',
    items: [
      { label: '活跃告警', value: 0 },
      { label: '历史告警', value: 0 },
      { label: '告警同比', value: 0, unit: '%' }
    ]
  }
})

// 图表实例
let trendChart = null
let pieChart = null
let heatChart = null

// 发布统计时间维度
const deployTimeRange = ref('week') // week, month, year

const heroStats = computed(() => ([
  {
    label: '主机资产',
    value: stats.assets.items[0]?.value || 0,
    hint: '当前纳管主机规模',
    tone: 'primary'
  },
  {
    label: '应用规模',
    value: stats.services.items[0]?.value || 0,
    hint: '已接入的应用数量',
    tone: 'success'
  },
  {
    label: '活跃告警',
    value: stats.monitor.items[0]?.value || 0,
    hint: '优先查看告警中心',
    tone: 'danger'
  },
  {
    label: '发布成功率',
    value: stats.deployment.items[2]?.value || 0,
    unit: '%',
    hint: '当前发布成功水平',
    tone: 'warning'
  }
]))

const priorityCards = computed(() => ([
  {
    title: '监控风险',
    value: stats.monitor.items[0]?.value || 0,
    description: '先确认是否有未处理告警和待跟进事件。',
    path: '/monitor/alert-center',
    tone: 'danger'
  },
  {
    title: '待执行作业',
    value: stats.deployment.items[1]?.value || 0,
    description: '查看任务中心和工单中心的积压动作。',
    path: '/task/job',
    tone: 'warning'
  },
  {
    title: '集群覆盖',
    value: stats.assets.items[2]?.value || 0,
    description: '聚焦容器资源和关键业务集群状态。',
    path: '/k8s/list',
    tone: 'primary'
  }
]))

const recommendedActions = computed(() => ([
  {
    title: '发起 AI 巡检',
    description: '通过对话式助手快速发起主机巡检和上下文分析。',
    action: () => openAIAssistant('帮我开始一次主机巡检')
  },
  {
    title: '查看待处理工单',
    description: '统一跟进发布、脚本、服务和 SQL 工单。',
    action: () => router.push('/work/orders')
  },
  {
    title: '检查全局搜索',
    description: '跨模块定位主机、应用、告警和知识。',
    action: () => router.push('/search/global')
  }
]))

// 快捷导航工具数据
const quickTools = reactive([])

// 编辑工具表单
const toolForm = reactive({
  title: '',
  icon: '',
  link: '',
  sort: 0
})

const normalizeToolIcon = value => normalizeDashboardToolIconUrl(value)

// 打开编辑弹窗
const openEditDialog = (tool, index) => {
  editingIndex.value = index
  editingTool.value = tool
  Object.assign(toolForm, {
    title: tool.title,
    icon: normalizeToolIcon(tool.icon),
    link: tool.link,
    sort: tool.sort || 0
  })
  editDialogVisible.value = true
}

// 添加新工具
const addNewTool = () => {
  editingIndex.value = -1
  editingTool.value = null
  Object.assign(toolForm, {
    title: '',
    icon: '',
    link: '',
    sort: 0
  })
  editDialogVisible.value = true
}

// 保存编辑
const saveToolEdit = async () => {
  if (!toolForm.title.trim()) {
    ElMessage.warning('请输入导航标题')
    return
  }

  if (!toolForm.icon) {
    ElMessage.warning('请上传导航图标')
    return
  }

  if (!toolForm.link.trim()) {
    ElMessage.warning('请输入链接地址')
    return
  }

  // 校验链接地址必须包含 http:// 或 https://
  const link = toolForm.link.trim()
  if (!link.startsWith('http://') && !link.startsWith('https://')) {
    ElMessage.warning('链接地址必须以 http:// 或 https:// 开头')
    return
  }

  try {
    const icon = normalizeToolIcon(toolForm.icon)
    if (editingIndex.value >= 0) {
      // 编辑现有工具
      await UpdateTool({
        id: editingTool.value.id,
        title: toolForm.title,
        icon,
        link: toolForm.link,
        sort: toolForm.sort
      })
      ElMessage.success('更新成功')
    } else {
      // 添加新工具
      await CreateTool({
        title: toolForm.title,
        icon,
        link: toolForm.link,
        sort: toolForm.sort
      })
      ElMessage.success('添加成功')
    }

    editDialogVisible.value = false
    // 重新加载导航工具列表
    await loadTools()
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败，请稍后重试')
  }
}

// 删除工具
const deleteTool = (index) => {
  const tool = quickTools[index]
  ElMessageBox.confirm('确定要删除这个导航吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await DeleteToolAPI(tool.id)
      ElMessage.success('删除成功')
      // 重新加载导航工具列表
      await loadTools()
    } catch (error) {
      console.error('删除失败:', error)
      ElMessage.error('删除失败，请稍后重试')
    }
  }).catch(() => {})
}

// 上传图标
const handleIconUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  // 验证文件类型
  if (!file.type.startsWith('image/')) {
    ElMessage.error('请上传图片文件')
    return
  }

  // 验证文件大小（限制为2MB）
  if (file.size > 2 * 1024 * 1024) {
    ElMessage.error('图片大小不能超过2MB')
    return
  }

  try {
    const formData = new FormData()
    formData.append('file', file)

    const response = await UploadIcon(formData)
    if (response.data && response.data.code === 200) {
      toolForm.icon = normalizeToolIcon(response.data.data)
      ElMessage.success('图标上传成功')
    } else {
      ElMessage.error(response.data?.message || '图标上传失败')
    }
  } catch (error) {
    console.error('上传图标失败:', error)
    ElMessage.error('图标上传失败，请稍后重试')
  }
}

// 触发文件选择
const triggerIconUpload = () => {
  document.getElementById('iconUpload').click()
}

// 点击导航项
const handleToolClick = (tool) => {
  if (!tool.link) return

  // 判断是外部链接还是内部路由
  if (tool.link.startsWith('http://') || tool.link.startsWith('https://')) {
    // 外部链接，新窗口打开
    window.open(tool.link, '_blank')
  } else {
    // 内部路由
    router.push(tool.link)
  }
}

// 获取发布数据（根据时间维度）
const openAIAssistant = (prompt = '') => {
  if (prompt) {
    router.push({ path: '/ai/assistant', query: { prompt } })
    return
  }
  router.push('/ai/assistant')
}

const openAIDiagnosis = () => {
  router.push('/ai/diagnosis')
}

const getDeploymentData = (timeRange) => {
  // 模拟数据，后续替换为真实API
  const mockData = {
    week: {
      labels: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
      production: [12, 15, 10, 18, 22, 8, 5],
      test: [25, 30, 28, 35, 40, 20, 15]
    },
    month: {
      labels: ['1日', '5日', '10日', '15日', '20日', '25日', '30日'],
      production: [45, 52, 48, 60, 55, 62, 58],
      test: [88, 95, 90, 102, 98, 105, 100]
    },
    year: {
      labels: ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月'],
      production: [180, 165, 195, 210, 205, 220, 215, 230, 225, 240, 235, 250],
      test: [320, 310, 340, 360, 355, 380, 375, 390, 385, 400, 395, 410]
    }
  }
  return mockData[timeRange]
}

// 初始化发布统计图
const initTrendChart = () => {
  const chartDom = document.getElementById('trendChart')
  if (!chartDom) return

  trendChart = echarts.init(chartDom, 'dark')
  updateTrendChart()
}

// 更新发布统计图
const updateTrendChart = () => {
  if (!trendChart) return

  const data = getDeploymentData(deployTimeRange.value)

  const option = {
    title: {
      text: '上线发布次数统计',
      left: 20,
      top: 10,
      textStyle: {
        fontSize: 16,
        fontWeight: 'normal',
          color: 'rgba(226, 232, 240, 0.92)'
      }
    },
    grid: {
      left: 60,
      right: 30,
      top: 60,
      bottom: 40
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: {
          backgroundColor: '#6a7985'
        }
      },
      formatter: (params) => {
        let result = params[0].name + '<br/>'
        params.forEach(item => {
          result += `${item.marker} ${item.seriesName}: ${item.value}次<br/>`
        })
        return result
      }
    },
    xAxis: {
      type: 'category',
      data: data.labels,
      axisLine: { lineStyle: { color: 'rgba(148, 163, 184, 0.22)' } },
      axisTick: { show: false },
      axisLabel: { color: 'rgba(148, 163, 184, 0.82)' }
    },
    yAxis: {
      type: 'value',
      name: '发布次数',
      nameTextStyle: { color: 'rgba(148, 163, 184, 0.82)', fontSize: 12 },
      splitLine: { lineStyle: { color: 'rgba(148, 163, 184, 0.12)' } },
      axisLabel: { color: 'rgba(148, 163, 184, 0.82)' }
    },
    legend: {
      data: ['生产环境', '测试环境'],
      top: 15,
      right: 120
    },
    series: [
      {
        name: '生产环境',
        type: 'line',
        smooth: true,
        data: data.production,
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(255, 107, 107, 0.4)' },
              { offset: 1, color: 'rgba(255, 107, 107, 0.05)' }
            ]
          }
        },
        lineStyle: { color: '#ff6b6b', width: 2 },
        itemStyle: { color: '#ff6b6b' }
      },
      {
        name: '测试环境',
        type: 'line',
        smooth: true,
        data: data.test,
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(78, 201, 176, 0.4)' },
              { offset: 1, color: 'rgba(78, 201, 176, 0.05)' }
            ]
          }
        },
        lineStyle: { color: '#4ecdc4', width: 2 },
        itemStyle: { color: '#4ecdc4' }
      }
    ]
  }
  trendChart.setOption(option)
}

// 切换时间维度
const changeTimeRange = (range) => {
  deployTimeRange.value = range
  updateTrendChart()
}

// 初始化环形图
const initPieChart = async () => {
  const chartDom = document.getElementById('pieChart')
  if (!chartDom) return

  pieChart = echarts.init(chartDom, 'dark')

  // 加载业务分布数据
  let businessData = []
  try {
    const response = await getBusinessDistribution()
    if (response.data && response.data.code === 200) {
      const data = response.data.data
      const colors = ['#5dade2', '#f8b739', '#48c9b0', '#9b59b6', '#ec7063', '#ff6b6b', '#4ecdc4', '#45b7d1']
      businessData = data.businessLines.map((line, index) => ({
        value: line.serviceCount,
        name: line.name,
        itemStyle: { color: colors[index % colors.length] }
      }))
    }
  } catch (error) {
    console.error('加载业务分布数据失败:', error)
    // 使用默认数据
    businessData = [
      { value: 10, name: '暂无数据', itemStyle: { color: 'rgba(148, 163, 184, 0.46)' } }
    ]
  }

  const option = {
    title: {
      text: '业务组应用分布',
      left: 20,
      top: 10,
      textStyle: {
        fontSize: 16,
        fontWeight: 'normal',
          color: 'rgba(226, 232, 240, 0.92)'
      }
    },
    tooltip: {
      trigger: 'item',
      formatter: '{b}<br/>应用数: {c}<br/>占比: {d}%'
    },
    legend: {
      orient: 'vertical',
      right: 30,
      top: 'center',
      itemWidth: 12,
      itemHeight: 12,
      textStyle: { fontSize: 12, color: 'rgba(148, 163, 184, 0.82)' }
    },
    series: [
      {
        type: 'pie',
        radius: ['50%', '70%'],
        center: ['35%', '50%'],
        avoidLabelOverlap: false,
        label: { show: false },
        labelLine: { show: false },
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.3)'
          }
        },
        data: businessData
      }
    ]
  }
  pieChart.setOption(option)
}

// 资源使用率类型
const resourceType = ref('cpu') // cpu, memory, disk

// 获取资源使用率数据
const getResourceData = (type) => {
  // 模拟数据，后续替换为真实API
  const mockData = {
    cpu: [
      { name: '服务器-01', value: 89.5 },
      { name: '服务器-02', value: 76.3 },
      { name: '服务器-03', value: 68.7 },
      { name: '服务器-04', value: 62.1 },
      { name: '服务器-05', value: 58.9 }
    ],
    memory: [
      { name: '服务器-03', value: 92.3 },
      { name: '服务器-07', value: 85.6 },
      { name: '服务器-01', value: 78.9 },
      { name: '服务器-12', value: 71.2 },
      { name: '服务器-05', value: 68.4 }
    ],
    disk: [
      { name: '服务器-05', value: 94.7 },
      { name: '服务器-08', value: 88.2 },
      { name: '服务器-11', value: 82.5 },
      { name: '服务器-03', value: 75.8 },
      { name: '服务器-09', value: 69.3 }
    ]
  }
  return mockData[type]
}

// 初始化资源使用率图表
const initHeatChart = () => {
  const chartDom = document.getElementById('heatChart')
  if (!chartDom) return

  heatChart = echarts.init(chartDom, 'dark')
  updateResourceChart()
}

// 更新资源使用率图表
const updateResourceChart = () => {
  if (!heatChart) return

  const data = getResourceData(resourceType.value)
  const titles = {
    cpu: 'CPU使用率 TOP5',
    memory: '内存使用率 TOP5',
    disk: '磁盘占用 TOP5'
  }
  const colors = {
    cpu: ['#ff6b6b', '#ee5a52', '#e74c3c', '#c0392b', '#a93226'],
    memory: ['#3498db', '#2980b9', '#2472a4', '#1f618d', '#1a5276'],
    disk: ['#f39c12', '#e67e22', '#d68910', '#ca6f1e', '#ba4a00']
  }

  const option = {
    title: {
      text: titles[resourceType.value],
      left: 20,
      top: 10,
      textStyle: {
        fontSize: 16,
        fontWeight: 'normal',
          color: 'rgba(226, 232, 240, 0.92)'
      }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      },
      formatter: (params) => {
        const item = params[0]
        return `${item.name}<br/>${item.marker} ${item.value}%`
      }
    },
    grid: {
      left: 100,
      right: 40,
      top: 60,
      bottom: 40
    },
    xAxis: {
      type: 'value',
      max: 100,
      axisLabel: {
        formatter: '{value}%',
          color: 'rgba(148, 163, 184, 0.82)'
      },
      splitLine: {
        lineStyle: { color: 'rgba(148, 163, 184, 0.12)' }
      }
    },
    yAxis: {
      type: 'category',
      data: data.map(item => item.name),
      axisLine: { lineStyle: { color: 'rgba(148, 163, 184, 0.22)' } },
      axisTick: { show: false },
      axisLabel: { color: 'rgba(148, 163, 184, 0.82)' }
    },
    series: [
      {
        type: 'bar',
        data: data.map((item, index) => ({
          value: item.value,
          itemStyle: {
            color: colors[resourceType.value][index]
          }
        })),
        barWidth: 20,
        label: {
          show: true,
          position: 'right',
          formatter: '{c}%',
          color: 'rgba(226, 232, 240, 0.88)',
          fontSize: 12
        }
      }
    ]
  }
  heatChart.setOption(option)
}

// 切换资源类型
const changeResourceType = (type) => {
  resourceType.value = type
  updateResourceChart()
}

// 加载导航工具列表
const loadTools = async () => {
  try {
    const response = await GetAllTools()
    if (response.data && response.data.code === 200) {
      // 清空现有数据
      quickTools.splice(0, quickTools.length)
      // 添加新数据
      quickTools.push(...response.data.data.map(tool => ({
        ...tool,
        icon: normalizeToolIcon(tool.icon)
      })))
    }
  } catch (error) {
    console.error('加载导航工具失败:', error)
  }
}

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const response = await getDashboardStats()
    if (response.data && response.data.code === 200) {
      const data = response.data.data

      // 更新资产详情
      stats.assets.items[0].value = data.hostStats?.total || 0
      stats.assets.items[1].value = data.databaseStats?.total || 0
      stats.assets.items[2].value = data.k8sClusterStats?.total || 0

      // 更新服务详情
      stats.services.items[0].value = data.serviceStats?.total || 0
      stats.services.items[1].value = data.serviceStats?.businessLines || 0

      // 更新发布详情
      stats.deployment.items[0].value = data.deploymentStats?.total || 0
      stats.deployment.items[1].value = data.taskStats?.total || 0
      stats.deployment.items[2].value = data.deploymentStats?.successRate || 0

      // 更新监控告警(模拟数据，需要根据实际API调整)
      stats.monitor.items[0].value = 0
      stats.monitor.items[1].value = 0
      stats.monitor.items[2].value = 0
    }
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.value = false
  }
}

const refreshDashboard = async () => {
  await Promise.all([loadData(), loadTools()])

  if (trendChart) {
    updateTrendChart()
  } else {
    initTrendChart()
  }

  if (pieChart) {
    pieChart.dispose()
    pieChart = null
  }
  await initPieChart()

  if (heatChart) {
    updateResourceChart()
  } else {
    initHeatChart()
  }
}

// 窗口大小改变时重绘图表
const handleResize = () => {
  trendChart?.resize()
  pieChart?.resize()
  heatChart?.resize()
}

// 生命周期
onMounted(async () => {
  await loadData()
  await loadTools()
  setTimeout(async () => {
    initTrendChart()
    await initPieChart()
    initHeatChart()
  }, 100)
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  trendChart?.dispose()
  pieChart?.dispose()
  heatChart?.dispose()
})
</script>

<template>
    <PageContainer wide class="dashboard-page">
      <PageHeader eyebrow="OpsNexus Command Center" :title="brand.name" :subtitle="brand.description">
        <template #intro>
          <PageIntro
            title="平台定位"
            text="统一资产、容器、告警、工单与 AI 智能运维助手工作台。首页只保留核心指标、风险待办和高频入口，让值班、诊断和处置在一个视图里完成切换。"
          >
            <div class="page-chip-row">
              <span class="platform-chip">资产治理</span>
              <span class="platform-chip">容器运维</span>
              <span class="platform-chip">告警响应</span>
              <span class="platform-chip">多模型 AI 协作</span>
            </div>
          </PageIntro>
        </template>
        <template #actions>
          <el-button :loading="loading" @click="refreshDashboard">刷新</el-button>
          <el-button type="primary" @click="openAIAssistant()">助手工作台</el-button>
          <el-button @click="openAIAssistant('帮我开始一次主机巡检')">智能巡检</el-button>
          <el-button @click="openAIDiagnosis">诊断分析台</el-button>
        </template>
      </PageHeader>

    <StatStrip :items="heroStats" />

    <div class="platform-data-grid platform-data-grid--two">
      <SectionCard title="风险与待办" subtitle="先确认现在是否出问题，再决定优先处理什么。">
        <div class="dashboard-priority-grid">
          <button
            v-for="item in priorityCards"
            :key="item.title"
            :class="['dashboard-priority-card', `dashboard-priority-card--${item.tone}`]"
            @click="router.push(item.path)"
          >
            <div class="dashboard-priority-card__value">{{ item.value }}</div>
            <div class="dashboard-priority-card__title">{{ item.title }}</div>
            <div class="dashboard-priority-card__desc">{{ item.description }}</div>
          </button>
        </div>
      </SectionCard>

      <SectionCard title="AI 智能运维助手工作台" subtitle="多模型接入、Agent 协作、知识检索与诊断巡检在统一入口收敛。">
        <div class="dashboard-ai-workbench">
          <div class="dashboard-ai-workbench__intro">
            <div class="dashboard-ai-workbench__eyebrow">{{ brand.sidebarBadge }}</div>
            <div class="dashboard-ai-workbench__title">{{ brand.slogan }}</div>
            <div class="dashboard-ai-workbench__desc">{{ brand.assistantPositioning }}</div>
          </div>

          <div class="page-chip-row dashboard-ai-workbench__models">
            <span v-for="item in brand.aiModels || []" :key="item" class="platform-chip">{{ item }}</span>
          </div>

          <div class="dashboard-ai-capability-grid">
            <div v-for="item in brand.aiCapabilities || []" :key="item" class="dashboard-ai-capability">
              {{ item }}
            </div>
          </div>

          <div class="dashboard-action-list">
            <button
              v-for="item in recommendedActions"
              :key="item.title"
              class="dashboard-action-item"
              @click="item.action()"
            >
              <div class="dashboard-action-item__title">{{ item.title }}</div>
              <div class="dashboard-action-item__desc">{{ item.description }}</div>
            </button>
          </div>
        </div>
      </SectionCard>
    </div>

    <!-- 图表区域 -->
    <div class="charts-row">
      <!-- 发布统计图 -->
      <div class="chart-card large">
        <div class="chart-header">
          <div class="time-range-tabs">
            <button
              :class="['tab-btn', { active: deployTimeRange === 'week' }]"
              @click="changeTimeRange('week')"
            >
              周
            </button>
            <button
              :class="['tab-btn', { active: deployTimeRange === 'month' }]"
              @click="changeTimeRange('month')"
            >
              月
            </button>
            <button
              :class="['tab-btn', { active: deployTimeRange === 'year' }]"
              @click="changeTimeRange('year')"
            >
              年
            </button>
          </div>
        </div>
        <div id="trendChart" style="width: 100%; height: calc(100% - 40px);"></div>
      </div>

      <!-- 环形图 -->
      <div class="chart-card">
        <div id="pieChart" style="width: 100%; height: 100%;"></div>
      </div>
    </div>

    <!-- 底部区域 -->
    <div class="bottom-row">
      <!-- 快捷导航工具 -->
      <div class="tools-card">
        <div class="tools-header">
          <div class="tools-title">快捷导航工具</div>
          <button class="add-tool-btn" @click="addNewTool">
            <svg viewBox="0 0 24 24" fill="currentColor" style="width: 16px; height: 16px;">
              <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
            </svg>
            添加
          </button>
        </div>
        <div v-if="quickTools.length" class="tools-grid">
          <div
            class="tool-item"
            v-for="(tool, index) in quickTools"
            :key="tool.id"
          >
            <div class="tool-actions">
              <button class="action-btn edit-btn" @click.stop="openEditDialog(tool, index)">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"/>
                </svg>
              </button>
              <button class="action-btn delete-btn" @click.stop="deleteTool(index)">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z"/>
                </svg>
              </button>
            </div>
            <div class="tool-content" @click="handleToolClick(tool)">
              <div class="tool-icon">
                <img v-if="tool.icon" :src="tool.icon" :alt="tool.title" />
                <div v-else class="icon-placeholder">?</div>
              </div>
              <div class="tool-info">
                <div class="tool-name">{{ tool.title }}</div>
              </div>
            </div>
          </div>
        </div>
        <EmptyState
          v-else
          title="还没有常用动作"
          description="添加高频导航后，首页会保留真正常用的运维入口，而不是堆满所有模块链接。"
        />
      </div>

      <!-- 编辑弹窗 -->
      <el-dialog
        v-model="editDialogVisible"
        :title="editingIndex >= 0 ? '编辑导航' : '添加导航'"
        width="550px"
        class="ops-dialog ops-overlay--md dashboard-tool-dialog"
        :close-on-click-modal="false"
      >
        <div class="edit-form">
          <div class="form-item">
            <label><span class="required">*</span> 导航标题</label>
            <input
              v-model="toolForm.title"
              type="text"
              placeholder="例如：百度"
              class="form-input"
            />
          </div>

          <div class="form-item">
            <label><span class="required">*</span> 导航图标</label>
            <div class="icon-upload-wrapper">
              <div class="icon-preview-box">
                <img v-if="toolForm.icon" :src="toolForm.icon" alt="图标预览" />
                <div v-else class="empty-icon">
                  <svg viewBox="0 0 24 24" fill="currentColor">
                    <path d="M21 19V5c0-1.1-.9-2-2-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2zM8.5 13.5l2.5 3.01L14.5 12l4.5 6H5l3.5-4.5z"/>
                  </svg>
                  <span>暂无图标</span>
                </div>
              </div>
              <button class="upload-btn" @click="triggerIconUpload">
                <svg viewBox="0 0 24 24" fill="currentColor" style="width: 16px; height: 16px;">
                  <path d="M9 16h6v-6h4l-7-7-7 7h4zm-4 2h14v2H5z"/>
                </svg>
                选择图标
              </button>
              <input
                id="iconUpload"
                type="file"
                accept="image/*"
                @change="handleIconUpload"
                style="display: none;"
              />
            </div>
            <div class="form-tip">支持 PNG、JPG、SVG 格式，大小不超过 2MB</div>
          </div>

          <div class="form-item">
            <label><span class="required">*</span> 链接地址</label>
            <input
              v-model="toolForm.link"
              type="text"
              placeholder="例如：https://www.baidu.com/"
              class="form-input"
            />
            <div class="form-tip">必须以 http:// 或 https:// 开头</div>
          </div>
        </div>

        <template #footer>
          <div class="dialog-footer">
            <button class="btn-cancel" @click="editDialogVisible = false">取消</button>
            <button class="btn-confirm" @click="saveToolEdit">保存</button>
          </div>
        </template>
      </el-dialog>

      <!-- 资源使用率 -->
      <div class="chart-card">
        <div class="chart-header">
          <div class="resource-tabs">
            <button
              :class="['tab-btn', { active: resourceType === 'cpu' }]"
              @click="changeResourceType('cpu')"
            >
              CPU
            </button>
            <button
              :class="['tab-btn', { active: resourceType === 'memory' }]"
              @click="changeResourceType('memory')"
            >
              内存
            </button>
            <button
              :class="['tab-btn', { active: resourceType === 'disk' }]"
              @click="changeResourceType('disk')"
            >
              磁盘
            </button>
          </div>
        </div>
        <div id="heatChart" style="width: 100%; height: calc(100% - 40px);"></div>
      </div>
    </div>
  </PageContainer>
</template>

<style scoped lang="scss">
.dashboard {
  position: relative;
  padding: 20px;
  background:
    radial-gradient(circle at top left, rgba(56, 189, 248, 0.14), transparent 28%),
    radial-gradient(circle at right, rgba(99, 102, 241, 0.16), transparent 24%),
    linear-gradient(145deg, #040b18 0%, #07111f 44%, #040813 100%);
  min-height: calc(100vh - 60px);
  color: #e2e8f0;
  overflow: hidden;
}

.dashboard-hero {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 20px;
  padding: 24px 28px;
  border-radius: 20px;
  border: 1px solid rgba(125, 211, 252, 0.14);
  background: rgba(7, 16, 31, 0.72);
  box-shadow: 0 24px 60px rgba(2, 8, 23, 0.3);
  backdrop-filter: blur(18px);
}

.hero-brand {
  display: flex;
  align-items: center;
  gap: 18px;
}

.hero-logo-wrap {
  width: 72px;
  height: 72px;
  border-radius: 22px;
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, rgba(14, 165, 233, 0.18), rgba(99, 102, 241, 0.12));
  border: 1px solid rgba(125, 211, 252, 0.18);
}

.hero-logo {
  width: 40px;
  height: 40px;
}

.hero-eyebrow {
  color: rgba(125, 211, 252, 0.88);
  font-size: 12px;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.hero-title {
  margin: 8px 0 10px;
  color: #f8fafc;
  font-size: 32px;
  line-height: 1.1;
}

.hero-description {
  margin: 0;
  max-width: 620px;
  color: rgba(226, 232, 240, 0.78);
  line-height: 1.8;
}

.hero-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.hero-btn {
  height: 44px;
  padding: 0 18px;
  border-radius: 14px;
  border: 1px solid rgba(148, 163, 184, 0.18);
  background: rgba(15, 23, 42, 0.72);
  color: #e2e8f0;
  cursor: pointer;
  transition: all 0.25s ease;

  &:hover {
    transform: translateY(-1px);
    border-color: rgba(56, 189, 248, 0.42);
    box-shadow: 0 18px 36px rgba(2, 8, 23, 0.24);
  }
}

.hero-btn--primary {
  border: none;
  background: linear-gradient(135deg, #0ea5e9, #4f46e5);
  color: #fff;
}

// 顶部统计卡片
.stats-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 20px;
}

.stat-card {
  background: rgba(7, 16, 31, 0.72);
  border: 1px solid rgba(148, 163, 184, 0.14);
  border-radius: 18px;
  padding: 16px 18px;
  backdrop-filter: blur(16px);
  box-shadow: 0 18px 44px rgba(2, 8, 23, 0.22);
  transition: all 0.3s ease;

  &:hover {
    box-shadow: 0 22px 54px rgba(14, 165, 233, 0.14);
    transform: translateY(-2px);
    border-color: rgba(56, 189, 248, 0.34);
  }
}

.stat-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
  padding-bottom: 10px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.stat-title {
  font-size: 15px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
}

.stat-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  svg {
    width: 22px;
    height: 22px;
  }
}

.stat-items {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 0;
}

.item-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
}

.item-value {
  font-size: 15px;
  font-weight: 600;
  color: #00d4ff;
}

// 图表行
.charts-row {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}

.chart-card {
  background: rgba(7, 16, 31, 0.72);
  border: 1px solid rgba(148, 163, 184, 0.14);
  border-radius: 18px;
  padding: 20px;
  box-shadow: 0 18px 44px rgba(2, 8, 23, 0.22);
  height: 400px;

  &.large {
    height: 400px;
  }
}

.chart-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 10px;
  padding-bottom: 10px;
}

.time-range-tabs,
.resource-tabs {
  display: flex;
  gap: 8px;
  background: rgba(255, 255, 255, 0.08);
  padding: 4px;
  border-radius: 6px;
}

.tab-btn {
  min-height: var(--control-height-sm);
  padding: 0 16px;
  border: none;
  background: transparent;
  color: rgba(255, 255, 255, 0.6);
  font-size: 13px;
  cursor: pointer;
  border-radius: var(--control-radius-pill);
  transition: all 0.3s ease;

  &:hover {
    color: rgba(255, 255, 255, 0.9);
    background: rgba(255, 255, 255, 0.1);
  }

  &.active {
    background: rgba(99, 102, 241, 0.3);
    color: #a5b4fc;
    font-weight: 500;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
  }
}

// 底部行
.bottom-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

// 快捷工具
.tools-card {
  background: rgba(7, 16, 31, 0.72);
  border: 1px solid rgba(148, 163, 184, 0.14);
  border-radius: 18px;
  padding: 20px;
  box-shadow: 0 18px 44px rgba(2, 8, 23, 0.22);
}

.tools-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.tools-title {
  font-size: 16px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.9);
}

.add-tool-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  min-height: var(--control-height-sm);
  padding: 0 14px;
  background: linear-gradient(135deg, #0ea5e9, #4f46e5);
  color: white;
  border: none;
  border-radius: var(--control-radius-md);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    background: linear-gradient(135deg, #38bdf8, #6366f1);
    transform: translateY(-1px);
  }
}

.tools-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 12px;
}

.tool-item {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px 8px;
  border-radius: var(--action-card-radius);
  border: 1px solid rgba(148, 163, 184, 0.14);
  background: rgba(15, 23, 42, 0.42);
  transition: all 0.3s ease;

  &:hover {
    background: rgba(14, 165, 233, 0.14);
    border-color: rgba(56, 189, 248, 0.34);
    box-shadow: 0 18px 38px rgba(2, 8, 23, 0.2);

    .tool-actions {
      opacity: 1;
    }
  }
}

.tool-actions {
  position: absolute;
  top: 4px;
  right: 4px;
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.action-btn {
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;

  svg {
    width: 14px;
    height: 14px;
  }

  &.edit-btn {
    background: #e3f2fd;
    color: #2196f3;

    &:hover {
      background: #2196f3;
      color: white;
    }
  }

  &.delete-btn {
    background: #ffebee;
    color: #f44336;

    &:hover {
      background: #f44336;
      color: white;
    }
  }
}

.tool-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  width: 100%;
}

.tool-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: rgba(2, 6, 23, 0.72);
  border: 1px solid rgba(148, 163, 184, 0.18);

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .icon-placeholder {
    font-size: 24px;
    color: #c0c4cc;
  }
}

.tool-info {
  text-align: center;
  width: 100%;
}

.tool-name {
  font-size: 13px;
  color: rgba(226, 232, 240, 0.9);
  font-weight: 500;
  line-height: 1.35;
  overflow-wrap: anywhere;
}

// 编辑弹窗样式
.edit-form {
  padding: 10px 0;
}

.form-item {
  margin-bottom: 24px;

  label {
    display: block;
    font-size: 14px;
    color: #333;
    margin-bottom: 10px;
    font-weight: 500;

    .required {
      color: #f56c6c;
      margin-right: 4px;
    }
  }
}

.form-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.3s ease;
  box-sizing: border-box;

  &:focus {
    outline: none;
    border-color: #2196f3;
  }

  &::placeholder {
    color: #c0c4cc;
  }
}

.form-tip {
  margin-top: 6px;
  font-size: 12px;
  color: #909399;
}

.icon-upload-wrapper {
  display: flex;
  align-items: center;
  gap: 16px;
}

.icon-preview-box {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px dashed #dcdfe6;
  overflow: hidden;
  background: #fafafa;

  img {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }

  .empty-icon {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    color: #c0c4cc;

    svg {
      width: 32px;
      height: 32px;
    }

    span {
      font-size: 11px;
    }
  }
}

.upload-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  background: #f5f7fa;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
  color: #606266;
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    background: #2196f3;
    border-color: #2196f3;
    color: white;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 10px;
}

.btn-cancel,
.btn-confirm {
  min-height: var(--control-height-md);
  padding: 0 20px;
  border: none;
  border-radius: var(--control-radius-md);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-cancel {
  background: #f5f7fa;
  color: #606266;

  &:hover {
    background: #e0e3e8;
  }
}

.btn-confirm {
  background: #2196f3;
  color: white;

  &:hover {
    background: #1976d2;
  }
}

// 响应式设计
@media (max-width: 1400px) {
  .dashboard-hero {
    display: grid;
  }

  .stats-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  .charts-row {
    grid-template-columns: 1fr;
  }

  .bottom-row {
    grid-template-columns: 1fr;
  }

  .tools-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .dashboard {
    padding: 12px;
  }

  .dashboard-hero {
    padding: 20px;
  }

  .hero-title {
    font-size: 26px;
  }

  .stats-cards {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .charts-row,
  .bottom-row {
    gap: 12px;
  }

  .tools-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.dashboard-page {
  .charts-row,
  .bottom-row {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 16px;
  }

  .chart-card,
  .tools-card {
    border-radius: 24px;
    background: linear-gradient(180deg, rgba(13, 25, 43, 0.98) 0%, rgba(9, 18, 32, 0.92) 100%);
    border: 1px solid var(--border-subtle);
    box-shadow: var(--shadow-card);
  }

  .chart-header,
  .tools-header {
    align-items: center;
    margin-bottom: 16px;
    padding-bottom: 0;
  }

  .time-range-tabs,
  .resource-tabs {
    background: rgba(255, 255, 255, 0.04);
    border-radius: 999px;
    padding: 4px;
  }

  .tab-btn {
    border-radius: 999px;
    color: var(--text-muted);
  }

  .tab-btn.active {
    background: rgba(53, 164, 255, 0.16);
    color: var(--text-primary);
    box-shadow: none;
  }

  .dashboard-priority-grid,
  .dashboard-action-list {
    display: grid;
    gap: 14px;
  }

  .dashboard-priority-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .dashboard-priority-card,
  .dashboard-action-item {
    padding: 18px;
    border-radius: var(--action-card-radius);
    border: 1px solid var(--border-subtle);
    background: rgba(255, 255, 255, 0.03);
    text-align: left;
    cursor: pointer;
    transition: transform var(--transition-fast), border-color var(--transition-fast), background var(--transition-fast);
  }

  .dashboard-priority-card:hover,
  .dashboard-action-item:hover {
    transform: translateY(-2px);
    border-color: var(--border-strong);
    background: rgba(53, 164, 255, 0.06);
  }

  .dashboard-priority-card__value {
    font-size: 30px;
    font-weight: 800;
    color: var(--text-primary);
  }

  .dashboard-priority-card__title,
  .dashboard-action-item__title {
    margin-top: 10px;
    font-size: 16px;
    font-weight: 700;
    color: var(--text-primary);
  }

  .dashboard-priority-card__desc,
  .dashboard-action-item__desc {
    margin-top: 6px;
    font-size: 13px;
    line-height: 1.65;
    color: var(--text-muted);
  }

  .dashboard-ai-workbench {
    display: grid;
    gap: 16px;
  }

  .dashboard-ai-workbench__intro {
    display: grid;
    gap: 8px;
    padding: 18px;
    border-radius: 20px;
    border: 1px solid rgba(125, 211, 252, 0.2);
    background:
      radial-gradient(circle at top left, rgba(56, 189, 248, 0.18), transparent 38%),
      linear-gradient(135deg, rgba(9, 20, 38, 0.96), rgba(8, 16, 30, 0.88));
  }

  .dashboard-ai-workbench__eyebrow {
    font-size: 11px;
    font-weight: 700;
    letter-spacing: 0.14em;
    text-transform: uppercase;
    color: rgba(125, 211, 252, 0.92);
  }

  .dashboard-ai-workbench__title {
    font-size: 22px;
    font-weight: 800;
    color: var(--text-primary);
  }

  .dashboard-ai-workbench__desc {
    font-size: 13px;
    line-height: 1.7;
    color: var(--text-muted);
  }

  .dashboard-ai-workbench__models {
    gap: 10px;
  }

  .dashboard-ai-capability-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 12px;
  }

  .dashboard-ai-capability {
    padding: 14px 16px;
    border-radius: var(--action-card-radius);
    border: 1px solid var(--border-subtle);
    background: rgba(255, 255, 255, 0.03);
    color: var(--text-secondary);
    font-size: 13px;
    font-weight: 600;
  }
}

@media (max-width: 1200px) {
  .dashboard-page {
    .dashboard-priority-grid,
    .charts-row,
    .bottom-row {
      grid-template-columns: 1fr;
    }

    .dashboard-ai-capability-grid {
      grid-template-columns: 1fr;
    }
  }
}
</style>
