<template>
  <el-card>
    <template #header>
      <div class="toolbar">
        <span>通知渠道</span>
        <el-button v-permission="'monitor:channel:create'" type="primary" @click="openCreate">新增渠道</el-button>
      </div>
    </template>

    <el-table :data="list" v-loading="loading">
      <el-table-column prop="name" label="名称" min-width="160" />
      <el-table-column prop="type" label="类型" width="120" />
      <el-table-column label="状态" width="100"><template #default="scope">{{ scope.row.enabled ? '启用' : '禁用' }}</template></el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="scope">
          <el-button v-permission="'monitor:channel:test'" link type="primary" @click="handleTest(scope.row.id)">测试</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" title="新增渠道" width="560px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.type" style="width: 100%">
            <el-option label="webhook" value="webhook" />
            <el-option label="email" value="email" />
          </el-select>
        </el-form-item>
        <el-form-item label="启用"><el-switch v-model="form.enabled" /></el-form-item>
        <el-form-item label="配置JSON"><el-input v-model="form.config" type="textarea" :rows="4" placeholder='例如 {"url":"https://example"}' /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { createChannelApi, listChannelsApi, testChannelApi, type ChannelItem } from '@/api/monitor'

const loading = ref(false)
const list = ref<ChannelItem[]>([])
const dialogVisible = ref(false)
const form = reactive<any>({ name: '', type: 'webhook', enabled: true, config: '{}' })

const loadData = async () => {
  loading.value = true
  try {
    const res = await listChannelsApi()
    list.value = res.list || []
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  Object.assign(form, { name: '', type: 'webhook', enabled: true, config: '{}' })
  dialogVisible.value = true
}

const submit = async () => {
  await createChannelApi(form)
  ElMessage.success('创建成功')
  dialogVisible.value = false
  await loadData()
}

const handleTest = async (id: number) => {
  await testChannelApi(id)
  ElMessage.success('测试成功')
}

onMounted(loadData)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
</style>
