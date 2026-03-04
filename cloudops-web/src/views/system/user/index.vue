<template>
  <el-card>
    <template #header>
      <div class="toolbar">
        <span>用户管理</span>
        <el-button v-permission="'system:user:create'" type="primary" @click="openCreate">新增用户</el-button>
      </div>
    </template>

    <el-form :inline="true" class="search">
      <el-form-item>
        <el-input v-model="query.keyword" placeholder="用户名/昵称/邮箱" clearable />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="loadData">查询</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="list" v-loading="loading" row-key="id">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="username" label="用户名" min-width="140" />
      <el-table-column prop="nickname" label="昵称" min-width="120" />
      <el-table-column prop="email" label="邮箱" min-width="180" />
      <el-table-column prop="phone" label="手机号" width="140" />
      <el-table-column prop="status" label="状态" width="90">
        <template #default="scope">{{ scope.row.status === 1 ? '启用' : '禁用' }}</template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="scope">
          <el-button v-permission="'system:user:update'" link type="primary" @click="openEdit(scope.row)">编辑</el-button>
          <el-button v-permission="'system:user:delete'" link type="danger" @click="handleDelete(scope.row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pager">
      <el-pagination
        background
        layout="total, prev, pager, next"
        :total="total"
        :current-page="query.page"
        :page-size="query.page_size"
        @current-change="handlePageChange"
      />
    </div>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑用户' : '新增用户'" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="用户名"><el-input v-model="form.username" :disabled="!!form.id" /></el-form-item>
        <el-form-item v-if="!form.id" label="密码"><el-input v-model="form.password" show-password /></el-form-item>
        <el-form-item label="昵称"><el-input v-model="form.nickname" /></el-form-item>
        <el-form-item label="邮箱"><el-input v-model="form.email" /></el-form-item>
        <el-form-item label="手机号"><el-input v-model="form.phone" /></el-form-item>
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
import { createUserApi, deleteUserApi, listUsersApi, updateUserApi, type UserItem } from '@/api/system'

const loading = ref(false)
const list = ref<UserItem[]>([])
const total = ref(0)
const dialogVisible = ref(false)

const query = reactive({ page: 1, page_size: 10, keyword: '' })
const form = reactive<any>({ id: 0, username: '', password: '', nickname: '', email: '', phone: '', status: 1, department_id: 0 })

const loadData = async () => {
  loading.value = true
  try {
    const res = await listUsersApi(query)
    list.value = res.list || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number) => {
  query.page = page
  loadData()
}

const openCreate = () => {
  Object.assign(form, { id: 0, username: '', password: '', nickname: '', email: '', phone: '', status: 1, department_id: 0 })
  dialogVisible.value = true
}

const openEdit = (row: UserItem) => {
  Object.assign(form, row, { password: '' })
  dialogVisible.value = true
}

const submit = async () => {
  if (form.id) {
    await updateUserApi(form.id, form)
    ElMessage.success('更新成功')
  } else {
    await createUserApi(form)
    ElMessage.success('创建成功')
  }
  dialogVisible.value = false
  await loadData()
}

const handleDelete = async (id: number) => {
  await ElMessageBox.confirm('确认删除该用户吗？', '提示', { type: 'warning' })
  await deleteUserApi(id)
  ElMessage.success('删除成功')
  await loadData()
}

onMounted(loadData)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
.search { margin-bottom: 12px; }
.pager { margin-top: 12px; display: flex; justify-content: flex-end; }
</style>
