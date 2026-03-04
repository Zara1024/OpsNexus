<template>
  <el-card>
    <template #header>
      <div class="toolbar">
        <span>凭据管理</span>
        <el-button v-permission="'cmdb:credential:create'" type="primary" @click="openCreate">新增凭据</el-button>
      </div>
    </template>

    <el-table :data="list" row-key="id">
      <el-table-column prop="name" label="名称" min-width="180" />
      <el-table-column prop="type" label="类型" width="120">
        <template #default="scope">{{ scope.row.type === 'key' ? '密钥' : '密码' }}</template>
      </el-table-column>
      <el-table-column prop="username" label="用户名" min-width="140" />
      <el-table-column label="敏感信息" min-width="160">
        <template #default>******</template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" min-width="180" />
      <el-table-column label="操作" width="220">
        <template #default="scope">
          <el-button v-permission="'cmdb:credential:update'" link type="primary" @click="openEdit(scope.row)">编辑</el-button>
          <el-button v-permission="'cmdb:credential:delete'" link type="danger" @click="handleDelete(scope.row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑凭据' : '新增凭据'" width="620px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="类型">
          <el-radio-group v-model="form.type">
            <el-radio label="password">密码</el-radio>
            <el-radio label="key">密钥</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="用户名"><el-input v-model="form.username" /></el-form-item>
        <el-form-item v-if="form.type === 'password'" label="密码"><el-input v-model="form.password" show-password /></el-form-item>
        <template v-else>
          <el-form-item label="私钥"><el-input v-model="form.private_key" type="textarea" :rows="5" /></el-form-item>
          <el-form-item label="私钥口令"><el-input v-model="form.passphrase" show-password /></el-form-item>
        </template>
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
import {
  createCMDBCredentialApi,
  deleteCMDBCredentialApi,
  listCMDBCredentialsApi,
  updateCMDBCredentialApi,
  type CMDBCredentialItem,
} from '@/api/cmdb'

const list = ref<CMDBCredentialItem[]>([])
const dialogVisible = ref(false)
const form = reactive<any>({ id: 0, name: '', type: 'password', username: '', password: '', private_key: '', passphrase: '' })

const loadData = async () => {
  const res = await listCMDBCredentialsApi()
  list.value = res.list || []
}

const openCreate = () => {
  Object.assign(form, { id: 0, name: '', type: 'password', username: '', password: '', private_key: '', passphrase: '' })
  dialogVisible.value = true
}

const openEdit = (row: CMDBCredentialItem) => {
  Object.assign(form, { id: row.id, name: row.name, type: row.type, username: row.username, password: '', private_key: '', passphrase: '' })
  dialogVisible.value = true
}

const submit = async () => {
  if (form.id) {
    await updateCMDBCredentialApi(form.id, form)
  } else {
    await createCMDBCredentialApi(form)
  }
  ElMessage.success('保存成功')
  dialogVisible.value = false
  await loadData()
}

const handleDelete = async (id: number) => {
  await ElMessageBox.confirm('确认删除该凭据吗？', '提示', { type: 'warning' })
  await deleteCMDBCredentialApi(id)
  ElMessage.success('删除成功')
  await loadData()
}

onMounted(loadData)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
</style>
