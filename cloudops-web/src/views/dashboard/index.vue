<template>
  <div class="dashboard-page" v-loading="loading">
    <el-row :gutter="16">
      <el-col :span="6">
        <el-card><template #header>主机总数</template><div class="stat-value">{{ overview.asset.host_total }}</div></el-card>
      </el-col>
      <el-col :span="6">
        <el-card><template #header>在线主机</template><div class="stat-value success">{{ overview.asset.host_online }}</div></el-card>
      </el-col>
      <el-col :span="6">
        <el-card><template #header>离线主机</template><div class="stat-value danger">{{ overview.asset.host_offline }}</div></el-card>
      </el-col>
      <el-col :span="6">
        <el-card><template #header>活跃告警</template><div class="stat-value warn">{{ overview.alert.active_total }}</div></el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="mt16">
      <el-col :span="12">
        <el-card>
          <template #header>告警趋势（7天）</template>
          <el-table :data="trends.alert_trend" size="small">
            <el-table-column prop="date" label="日期" width="140" />
            <el-table-column prop="count" label="告警数" />
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>资源水位均值</template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="CPU">{{ overview.resource.cpu_avg.toFixed(2) }}%</el-descriptions-item>
            <el-descriptions-item label="内存">{{ overview.resource.memory_avg.toFixed(2) }}%</el-descriptions-item>
            <el-descriptions-item label="磁盘">{{ overview.resource.disk_avg.toFixed(2) }}%</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="mt16">
      <el-col :span="24">
        <el-card>
          <template #header>最近操作日志</template>
          <el-table :data="overview.recent_ops" size="small">
            <el-table-column prop="username" label="用户" width="140" />
            <el-table-column prop="module" label="模块" width="120" />
            <el-table-column prop="action" label="动作" width="120" />
            <el-table-column prop="path" label="路径" min-width="260" />
            <el-table-column prop="created_at" label="时间" width="180" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { getDashboardOverviewApi, getDashboardTrendsApi, type DashboardOverview, type DashboardTrends } from '@/api/dashboard'

const loading = ref(false)

const overview = reactive<DashboardOverview>({
  asset: { host_total: 0, host_online: 0, host_offline: 0, host_alert: 0, k8s_cluster_total: 0, db_instance_total: 0 },
  alert: { active_total: 0, today_handled: 0, severity_stats: [] },
  resource: { cpu_avg: 0, memory_avg: 0, disk_avg: 0 },
  recent_ops: [],
  recent_alert_top5: [],
})

const trends = reactive<DashboardTrends>({ days: 7, alert_trend: [], resource_trend: [] })

const loadData = async () => {
  loading.value = true
  try {
    Object.assign(overview, await getDashboardOverviewApi())
    Object.assign(trends, await getDashboardTrendsApi(7))
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.mt16 { margin-top: 16px; }
.stat-value { font-size: 28px; font-weight: 700; }
.success { color: #67c23a; }
.danger { color: #f56c6c; }
.warn { color: #e6a23c; }
</style>
