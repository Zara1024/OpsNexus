<template>
  <el-card>
    <template #header>
      <div class="toolbar">
        <span>分组管理</span>
        <el-button v-permission="'cmdb:group:create'" type="primary" @click="openCreate">新增分组</el-button>
      </div>
    </template>

    <el-table :data="list" row-key="id" default-expand-all>
      <el-table-column prop="group_name" label="分组名称" min-width="180" />
      <el-table-column prop="description" label="描述" min-width="220" />
      <el-table-column prop="sort_order" label="排序" width="100" />
      <el-table-column label="操作" width="220">
        <template #default="scope">
          <el-button v-permission="'cmdb:group:update'" link type="primary" @click="openEdit(scope.row)">编辑</el-button>
          <el-button v-permission="'cmdb:group:delete'" link type="danger" @click="handleDelete(scope.row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑分组' : '新增分组'" width="560px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="上级ID"><el-input-number v-model="form.parent_id" :min="0" /></el-form-item>
        <el-form-item label="分组名称"><el-input v-model="form.group_name" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort_order" :min="0" /></el-form-item>
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
import { createCMDBGroupApi, deleteCMDBGroupApi, listCMDBGroupsApi, updateCMDBGroupApi, type CMDBGroupNode } from '@/api/cmdb'

const list = ref<CMDBGroupNode[]>([])
const dialogVisible = ref(false)
const form = reactive<any>({ id: 0, parent_id: 0, group_name: '', description: '', sort_order: 0 })

const loadData = async () => {
  const res = await listCMDBGroupsApi()
  list.value = res.list || []
}

const openCreate = () => {
  Object.assign(form, { id: 0, parent_id: 0, group_name: '', description: '', sort_order: 0 })
  dialogVisible.value = true
}

const openEdit = (row: CMDBGroupNode) => {
  Object.assign(form, row)
  dialogVisible.value = true
}

const submit = async () => {
  if (form.id) {
    await updateCMDBGroupApi(form.id, form)
  } else {
    await createCMDBGroupApi(form)
  }
  ElMessage.success('保存成功')
  dialogVisible.value = false
  await loadData()
}

const handleDelete = async (id: number) => {
  await ElMessageBox.confirm('确认删除该分组吗？', '提示', { type: 'warning' })
  await deleteCMDBGroupApi(id)
  ElMessage.success('删除成功')
  await loadData()
}

onMounted(loadData)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
</style>
