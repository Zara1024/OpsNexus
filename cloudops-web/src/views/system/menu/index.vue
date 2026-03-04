<template>
  <el-card>
    <template #header>
      <div class="toolbar">
        <span>菜单管理</span>
        <el-button v-permission="'system:menu:create'" type="primary" @click="openCreate">新增菜单</el-button>
      </div>
    </template>

    <el-table :data="list" v-loading="loading" row-key="id" default-expand-all>
      <el-table-column prop="menu_name" label="菜单名称" min-width="160" />
      <el-table-column prop="route_path" label="路由" min-width="160" />
      <el-table-column prop="menu_type" label="类型" width="100" />
      <el-table-column prop="permission" label="权限标识" min-width="180" />
      <el-table-column label="操作" width="220">
        <template #default="scope">
          <el-button v-permission="'system:menu:update'" link type="primary" @click="openEdit(scope.row)">编辑</el-button>
          <el-button v-permission="'system:menu:delete'" link type="danger" @click="handleDelete(scope.row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑菜单' : '新增菜单'" width="560px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="上级ID"><el-input-number v-model="form.parent_id" :min="0" /></el-form-item>
        <el-form-item label="菜单名称"><el-input v-model="form.menu_name" /></el-form-item>
        <el-form-item label="路由地址"><el-input v-model="form.route_path" /></el-form-item>
        <el-form-item label="组件路径"><el-input v-model="form.component" /></el-form-item>
        <el-form-item label="菜单类型">
          <el-select v-model="form.menu_type"><el-option label="目录" value="dir" /><el-option label="菜单" value="menu" /><el-option label="按钮" value="button" /></el-select>
        </el-form-item>
        <el-form-item label="权限标识"><el-input v-model="form.permission" /></el-form-item>
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
import { createMenuApi, deleteMenuApi, menuTreeApi, updateMenuApi, type MenuItem } from '@/api/system'

const loading = ref(false)
const list = ref<MenuItem[]>([])
const dialogVisible = ref(false)
const form = reactive<any>({ id: 0, parent_id: 0, menu_name: '', route_path: '', component: '', icon: '', sort_order: 0, menu_type: 'menu', permission: '' })

const loadData = async () => {
  loading.value = true
  try {
    const res = await menuTreeApi()
    list.value = res.list || []
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  Object.assign(form, { id: 0, parent_id: 0, menu_name: '', route_path: '', component: '', icon: '', sort_order: 0, menu_type: 'menu', permission: '' })
  dialogVisible.value = true
}

const openEdit = (row: MenuItem) => {
  Object.assign(form, row)
  dialogVisible.value = true
}

const submit = async () => {
  if (form.id) {
    await updateMenuApi(form.id, form)
  } else {
    await createMenuApi(form)
  }
  ElMessage.success('保存成功')
  dialogVisible.value = false
  await loadData()
}

const handleDelete = async (id: number) => {
  await ElMessageBox.confirm('确认删除该菜单吗？', '提示', { type: 'warning' })
  await deleteMenuApi(id)
  ElMessage.success('删除成功')
  await loadData()
}

onMounted(loadData)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
</style>
