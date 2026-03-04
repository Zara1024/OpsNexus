<template>
  <el-card>
    <template #header>
      <div class="toolbar">
        <span>告警规则</span>
        <el-button v-permission="'monitor:rule:create'" type="primary" @click="openCreate">新增规则</el-button>
      </div>
    </template>

    <el-table :data="list" v-loading="loading">
      <el-table-column prop="name" label="名称" min-width="160" />
      <el-table-column prop="severity" label="等级" width="100" />
      <el-table-column prop="expression" label="表达式" min-width="260" show-overflow-tooltip />
      <el-table-column prop="duration" label="持续时长" width="120" />
      <el-table-column label="状态" width="100"><template #default="scope">{{ scope.row.enabled ? '启用' : '禁用' }}</template></el-table-column>
      <el-table-column label="操作" min-width="260" fixed="right">
        <template #default="scope">
          <el-button v-permission="'monitor:rule:update'" link type="primary" @click="openEdit(scope.row)">编辑</el-button>
          <el-button v-permission="'monitor:rule:toggle'" link type="warning" @click="handleToggle(scope.row)">{{ scope.row.enabled ? '禁用' : '启用' }}</el-button>
          <el-button v-permission="'monitor:rule:delete'" link type="danger" @click="handleDelete(scope.row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑规则' : '新增规则'" width="720px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="等级">
          <el-select v-model="form.severity" style="width: 100%">
            <el-option label="P1" value="P1" /><el-option label="P2" value="P2" /><el-option label="P3" value="P3" /><el-option label="P4" value="P4" />
          </el-select>
        </el-form-item>
        <el-form-item label="表达式"><el-input v-model="form.expression" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="持续时长"><el-input v-model="form.duration" placeholder="如 5m" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" /></el-form-item>
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { createAlertRuleApi, deleteAlertRuleApi, listAlertRulesApi, toggleAlertRuleApi, updateAlertRuleApi, type AlertRuleItem } from '@/api/monitor'

const loading = ref(false)
const list = ref<AlertRuleItem[]>([])
const dialogVisible = ref(false)
const form = reactive<any>({ id: 0, name: '', severity: 'P3', expression: '', duration: '5m', description: '' })

const loadData = async () => {
  loading.value = true
  try {
    const res = await listAlertRulesApi()
    list.value = res.list || []
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  Object.assign(form, { id: 0, name: '', severity: 'P3', expression: '', duration: '5m', description: '' })
  dialogVisible.value = true
}
const openEdit = (row: AlertRuleItem) => {
  Object.assign(form, row)
  dialogVisible.value = true
}

const submit = async () => {
  if (form.id) {
    await updateAlertRuleApi(form.id, form)
  } else {
    await createAlertRuleApi(form)
  }
  ElMessage.success('保存成功')
  dialogVisible.value = false
  await loadData()
}

const handleToggle = async (row: AlertRuleItem) => {
  await toggleAlertRuleApi(row.id, !row.enabled)
  ElMessage.success('状态已更新')
  await loadData()
}

const handleDelete = async (id: number) => {
  await ElMessageBox.confirm('确认删除该规则吗？', '提示', { type: 'warning' })
  await deleteAlertRuleApi(id)
  ElMessage.success('删除成功')
  await loadData()
}

onMounted(loadData)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
</style>
