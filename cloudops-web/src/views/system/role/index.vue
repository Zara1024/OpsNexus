<template>
  <el-card>
    <template #header>
      <div class="toolbar">
        <span>角色管理</span>
        <el-button v-permission="'system:role:create'" type="primary" @click="openCreate">新增角色</el-button>
      </div>
    </template>

    <el-table :data="list" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="role_code" label="角色编码" min-width="180" />
      <el-table-column prop="role_name" label="角色名称" min-width="140" />
      <el-table-column prop="description" label="描述" min-width="200" />
      <el-table-column label="操作" width="280">
        <template #default="scope">
          <el-button v-permission="'system:role:update'" link type="primary" @click="openEdit(scope.row)">编辑</el-button>
          <el-button v-permission="'system:role:assign_menus'" link type="primary" @click="openAssign(scope.row)">分配菜单</el-button>
          <el-button v-permission="'system:role:delete'" link type="danger" @click="handleDelete(scope.row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑角色' : '新增角色'" width="500px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="角色编码"><el-input v-model="form.role_code" :disabled="!!form.id" /></el-form-item>
        <el-form-item label="角色名称"><el-input v-model="form.role_name" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="assignVisible" title="分配菜单" width="520px">
      <el-tree ref="treeRef" :data="menuTree" node-key="id" show-checkbox :props="{ label: 'menu_name', children: 'children' }" />
      <template #footer>
        <el-button @click="assignVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAssign">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { assignRoleMenusApi, createRoleApi, deleteRoleApi, listRolesApi, menuTreeApi, updateRoleApi, type MenuItem, type RoleItem } from '@/api/system'

const loading = ref(false)
const list = ref<RoleItem[]>([])
const dialogVisible = ref(false)
const assignVisible = ref(false)
const menuTree = ref<MenuItem[]>([])
const treeRef = ref<any>()
const currentRoleId = ref(0)

const form = reactive<any>({ id: 0, role_code: '', role_name: '', description: '' })

const loadData = async () => {
  loading.value = true
  try {
    const res = await listRolesApi()
    list.value = res.list || []
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  Object.assign(form, { id: 0, role_code: '', role_name: '', description: '' })
  dialogVisible.value = true
}

const openEdit = (row: RoleItem) => {
  Object.assign(form, row)
  dialogVisible.value = true
}

const submit = async () => {
  if (form.id) {
    await updateRoleApi(form.id, form)
    ElMessage.success('更新成功')
  } else {
    await createRoleApi(form)
    ElMessage.success('创建成功')
  }
  dialogVisible.value = false
  await loadData()
}

const handleDelete = async (id: number) => {
  await ElMessageBox.confirm('确认删除该角色吗？', '提示', { type: 'warning' })
  await deleteRoleApi(id)
  ElMessage.success('删除成功')
  await loadData()
}

const openAssign = async (row: RoleItem) => {
  currentRoleId.value = row.id
  const res = await menuTreeApi()
  menuTree.value = res.list || []
  assignVisible.value = true
}

const submitAssign = async () => {
  const checked = treeRef.value?.getCheckedKeys?.() || []
  await assignRoleMenusApi(currentRoleId.value, checked)
  ElMessage.success('分配成功')
  assignVisible.value = false
}

onMounted(loadData)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
</style>
