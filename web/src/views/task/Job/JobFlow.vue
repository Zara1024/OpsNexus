<template>
  <el-dialog
    v-model="dialogVisible"
    class="modern-dialog task-flow-dialog"
    title="任务流程"
    width="84%"
    top="5vh"
    destroy-on-close
    @closed="handleDialogClose"
  >
    <div class="flow-shell">
      <div class="flow-shell__intro">
        <div class="flow-shell__title">任务执行链路</div>
        <div class="flow-shell__subtitle">进入后可启动整条任务链、查看每个模板节点的脚本与日志，并实时刷新状态。</div>
      </div>
      <div class="flow-container">
        <TaskFlow ref="taskFlowRef" :steps="steps" :task-id="currentTaskId" />
      </div>
    </div>
    <template #footer>
      <el-button @click="dialogVisible = false">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, ref } from 'vue'
import { GetTaskTemplates } from '@/api/task'
import TaskFlow from './TaskFlowTemp.vue'

const dialogVisible = ref(false)
const taskFlowRef = ref(null)
const taskCount = ref(0)
const currentTaskId = ref(null)
const templatesData = ref([])

const steps = computed(() => {
  if (!currentTaskId.value) {
    return [{
      task_count: 1,
      task_status: 1,
      template_id: 1,
      template_name: '默认任务',
      template_remark: '等待启动',
      id: 'task-1',
      hasDownConnector: false,
      showHorizontalConnector: false
    }]
  }

  if (currentTaskId.value && templatesData.value.length > 0) {
    return templatesData.value.map((template, index) => ({
      ...template,
      id: `task-${index + 1}`,
      hasDownConnector: (index + 1) % 3 === 0 && index < templatesData.value.length - 1,
      showHorizontalConnector: (index + 1) % 3 !== 0 && index < templatesData.value.length - 1
    }))
  }

  return []
})

const showFlow = async taskId => {
  try {
    if (!taskId) {
      throw new Error('任务ID未定义')
    }

    currentTaskId.value = taskId
    templatesData.value = []

    const response = await GetTaskTemplates({ id: taskId })
    if (!response?.data?.data) {
      throw new Error('无效的API响应数据')
    }

    const apiData = JSON.parse(JSON.stringify(response.data.data))
    templatesData.value = apiData.map(step => ({
      ...step,
      task_status: step.task_status ?? 1
    }))

    const apiTaskId = templatesData.value[0]?.task_id
    if (!apiTaskId) {
      throw new Error('API响应中未包含task_id')
    }

    currentTaskId.value = apiTaskId
    taskCount.value = templatesData.value.length
    dialogVisible.value = true
    return apiTaskId
  } catch (error) {
    console.error('获取模板信息失败:', error)
    templatesData.value = []
    taskCount.value = 1
    dialogVisible.value = true
    throw error
  }
}

const handleDialogClose = () => {
  if (taskFlowRef.value?.stopAllPolling) {
    taskFlowRef.value.stopAllPolling()
  }
}

defineExpose({
  showFlow
})
</script>

<style scoped>
.flow-shell {
  display: grid;
  gap: 16px;
}

.flow-shell__intro {
  padding: 14px 16px;
  border-radius: 18px;
  border: 1px solid var(--border-subtle);
  background: rgba(255, 255, 255, 0.04);
  display: grid;
  gap: 4px;
}

.flow-shell__title {
  color: var(--text-primary);
  font-size: 16px;
  font-weight: 700;
}

.flow-shell__subtitle {
  color: var(--text-muted);
  line-height: 1.6;
}

.flow-container {
  min-height: 400px;
  padding: 18px;
  border-radius: 20px;
  border: 1px solid var(--border-subtle);
  background: rgba(4, 11, 23, 0.56);
}
</style>
