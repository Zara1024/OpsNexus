<template>
  <el-card>
    <template #header>
      <div class="toolbar">
        <span>部门管理</span>
        <el-button v-permission="'system:department:create'" type="primary" @click="openCreate">新增部门</el-button>
      </div>
    </template>

    <el-table :data="list" v-loading="loading" row-key="id" default-expand-all>
      <el-table-column prop="dept_name" label="部门名称" min-width="180" />
      <el-table-column prop="leader" label="负责人" width="120" />
      <el-table-column prop="phone" label="电话" width="140" />
      <el-table-column prop="email" label="邮箱" min-width="200" />
      <el-table-column label="操作" width="220">
        <template #default="scope">
          <el-button v-permission="'system:department:update'" link type="primary" @click="openEdit(scope.row)">编辑</el-button>
          <el-button v-permission="'system:department:delete'" link type="danger" @click="handleDelete(scope.row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑部门' : '新增部门'" width="560px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="上级ID"><el-input-number v-model="form.parent_id" :min="0" /></el-form-item>
        <el-form-item label="部门名称"><el-input v-model="form.dept_name" /></el-form-item>
        <el-form-item label="负责人"><el-input v-model="form.leader" /></el-form-item>
        <el-form-item label="电话"><el-input v-model="form.phone" /></el-form-item>
        <el-form-item label="邮箱"><el-input v-model="form.email" /></el-form-item>
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
import { createDepartmentApi, deleteDepartmentApi, listDepartmentsApi, updateDepartmentApi, type DepartmentItem } from '@/api/system'

const loading = ref(false)
const list = ref<DepartmentItem[]>([])
const dialogVisible = ref(false)
const form = reactive<any>({ id: 0, parent_id: 0, dept_name: '', sort_order: 0, leader: '', phone: '', email: '', status: 1 })

const toTree = (rows: DepartmentItem[]) => {
  const map = new Map<number, any>()
  const roots: any[] = []
  rows.forEach((item) => map.set(item.id, { ...item, children: [] }))
  map.forEach((node) => {
    if (node.parent_id && map.has(node.parent_id)) map.get(node.parent_id).children.push(node)
    else roots.push(node)
  })
  return roots
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await listDepartmentsApi()
    list.value = toTree(res.list || [])
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  Object.assign(form, { id: 0, parent_id: 0, dept_name: '', sort_order: 0, leader: '', phone: '', email: '', status: 1 })
  dialogVisible.value = true
}

const openEdit = (row: DepartmentItem) => {
  Object.assign(form, row)
  dialogVisible.value = true
}

const submit = async () => {
  if (form.id) {
    await updateDepartmentApi(form.id, form)
  } else {
    await createDepartmentApi(form)
  }
  ElMessage.success('保存成功')
  dialogVisible.value = false
  await loadData()
}

const handleDelete = async (id: number) => {
  await ElMessageBox.confirm('确认删除该部门吗？', '提示', { type: 'warning' })
  await deleteDepartmentApi(id)
  ElMessage.success('删除成功')
  await loadData()
}

onMounted(loadData)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
</style>
