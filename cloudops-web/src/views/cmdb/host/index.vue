<template>
  <el-card>
    <template #header>
      <div class="toolbar">
        <span>主机管理</span>
        <div class="actions">
          <el-button v-permission="'cmdb:host:create'" type="primary" @click="openCreate">新增主机</el-button>
          <el-button v-permission="'cmdb:host:batch'" type="warning" plain @click="handleBatch('enable')">批量启用</el-button>
          <el-button v-permission="'cmdb:host:batch'" type="info" plain @click="handleBatch('disable')">批量禁用</el-button>
          <el-button v-permission="'cmdb:host:batch'" type="danger" plain @click="handleBatch('delete')">批量删除</el-button>
        </div>
      </div>
    </template>

    <el-form :inline="true" class="search">
      <el-form-item>
        <el-input v-model="query.keyword" placeholder="主机名/IP" clearable />
      </el-form-item>
      <el-form-item>
        <el-select v-model="query.group_id" placeholder="分组" clearable style="width: 160px">
          <el-option v-for="g in groupOptions" :key="g.id" :label="g.group_name" :value="g.id" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-select v-model="query.status" placeholder="状态" clearable style="width: 120px">
          <el-option label="启用" :value="1" />
          <el-option label="禁用" :value="2" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="loadData">查询</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="list" v-loading="loading" row-key="id" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="50" />
      <el-table-column prop="hostname" label="主机名" min-width="130" />
      <el-table-column prop="inner_ip" label="内网IP" min-width="130" />
      <el-table-column prop="outer_ip" label="公网IP" min-width="130" />
      <el-table-column prop="os_type" label="系统" min-width="100" />
      <el-table-column prop="cpu_cores" label="CPU" width="80" />
      <el-table-column prop="memory_gb" label="内存(GB)" width="100" />
      <el-table-column prop="status" label="状态" width="80">
        <template #default="scope">{{ scope.row.status === 1 ? '启用' : '禁用' }}</template>
      </el-table-column>
      <el-table-column label="操作" min-width="330" fixed="right">
        <template #default="scope">
          <el-button v-permission="'cmdb:host:detail'" link type="primary" @click="openDetail(scope.row.id)">详情</el-button>
          <el-button v-permission="'cmdb:host:update'" link type="primary" @click="openEdit(scope.row)">编辑</el-button>
          <el-button v-permission="'cmdb:host:test'" link type="success" @click="handleTest(scope.row.id)">SSH测试</el-button>
          <el-button v-permission="'cmdb:host:terminal'" link type="warning" @click="openTerminal(scope.row.id)">终端</el-button>
          <el-button v-permission="'cmdb:host:delete'" link type="danger" @click="handleDelete(scope.row.id)">删除</el-button>
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

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑主机' : '新增主机'" width="760px">
      <el-form :model="form" label-width="100px">
        <el-row :gutter="12">
          <el-col :span="12"><el-form-item label="主机名"><el-input v-model="form.hostname" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="内网IP"><el-input v-model="form.inner_ip" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="公网IP"><el-input v-model="form.outer_ip" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="系统类型"><el-input v-model="form.os_type" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="系统版本"><el-input v-model="form.os_version" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="架构"><el-input v-model="form.arch" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="CPU核数"><el-input-number v-model="form.cpu_cores" :min="0" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="内存GB"><el-input-number v-model="form.memory_gb" :min="0" :precision="2" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="磁盘GB"><el-input-number v-model="form.disk_gb" :min="0" :precision="2" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="分组ID"><el-input-number v-model="form.group_id" :min="0" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="凭据ID"><el-input-number v-model="form.credential_id" :min="0" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="云厂商"><el-input v-model="form.cloud_provider" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="实例ID"><el-input v-model="form.instance_id" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="地域"><el-input v-model="form.region" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="Agent状态"><el-input v-model="form.agent_status" /></el-form-item></el-col>
          <el-col :span="24"><el-form-item label="标签JSON"><el-input v-model="form.labels" type="textarea" :rows="2" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="状态"><el-switch v-model="statusSwitch" /></el-form-item></el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailVisible" title="主机详情" width="620px">
      <el-descriptions :column="2" border v-if="detail">
        <el-descriptions-item label="主机名">{{ detail.hostname }}</el-descriptions-item>
        <el-descriptions-item label="内网IP">{{ detail.inner_ip }}</el-descriptions-item>
        <el-descriptions-item label="公网IP">{{ detail.outer_ip }}</el-descriptions-item>
        <el-descriptions-item label="系统">{{ detail.os_type }} {{ detail.os_version }}</el-descriptions-item>
        <el-descriptions-item label="CPU/内存">{{ detail.cpu_cores }} / {{ detail.memory_gb }}</el-descriptions-item>
        <el-descriptions-item label="磁盘">{{ detail.disk_gb }}</el-descriptions-item>
        <el-descriptions-item label="分组ID">{{ detail.group_id }}</el-descriptions-item>
        <el-descriptions-item label="凭据ID">{{ detail.credential_id }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  batchCMDBHostsApi,
  createCMDBHostApi,
  deleteCMDBHostApi,
  getCMDBHostApi,
  listCMDBGroupsApi,
  listCMDBHostsApi,
  testCMDBHostSSHApi,
  updateCMDBHostApi,
  type CMDBGroupNode,
  type CMDBHostItem,
} from '@/api/cmdb'

const router = useRouter()
const loading = ref(false)
const list = ref<CMDBHostItem[]>([])
const total = ref(0)
const selectedIDs = ref<number[]>([])
const groupOptions = ref<CMDBGroupNode[]>([])

const dialogVisible = ref(false)
const detailVisible = ref(false)
const detail = ref<CMDBHostItem | null>(null)

const query = reactive<any>({ page: 1, page_size: 10, keyword: '', group_id: undefined, status: undefined })
const form = reactive<any>({
  id: 0, hostname: '', inner_ip: '', outer_ip: '', os_type: '', os_version: '', arch: '',
  cpu_cores: 0, memory_gb: 0, disk_gb: 0, group_id: 0, credential_id: 0,
  cloud_provider: '', instance_id: '', region: '', agent_status: 'unknown', labels: '{}', status: 1,
})
const statusSwitch = computed({
  get: () => form.status === 1,
  set: (v: boolean) => { form.status = v ? 1 : 2 },
})

const flattenGroups = (nodes: CMDBGroupNode[]): CMDBGroupNode[] => {
  const result: CMDBGroupNode[] = []
  const walk = (items: CMDBGroupNode[]) => {
    items.forEach((i) => {
      result.push(i)
      if (i.children?.length) walk(i.children)
    })
  }
  walk(nodes)
  return result
}

const loadGroups = async () => {
  const res = await listCMDBGroupsApi()
  groupOptions.value = flattenGroups(res.list || [])
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await listCMDBHostsApi(query)
    list.value = res.list || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

const handleSelectionChange = (rows: CMDBHostItem[]) => {
  selectedIDs.value = rows.map((r) => r.id)
}

const handlePageChange = (page: number) => {
  query.page = page
  loadData()
}

const openCreate = () => {
  Object.assign(form, {
    id: 0, hostname: '', inner_ip: '', outer_ip: '', os_type: '', os_version: '', arch: '',
    cpu_cores: 0, memory_gb: 0, disk_gb: 0, group_id: 0, credential_id: 0,
    cloud_provider: '', instance_id: '', region: '', agent_status: 'unknown', labels: '{}', status: 1,
  })
  dialogVisible.value = true
}

const openEdit = (row: CMDBHostItem) => {
  Object.assign(form, row)
  dialogVisible.value = true
}

const openDetail = async (id: number) => {
  detail.value = await getCMDBHostApi(id)
  detailVisible.value = true
}

const submit = async () => {
  if (form.id) {
    await updateCMDBHostApi(form.id, form)
  } else {
    await createCMDBHostApi(form)
  }
  ElMessage.success('保存成功')
  dialogVisible.value = false
  await loadData()
}

const handleDelete = async (id: number) => {
  await ElMessageBox.confirm('确认删除该主机吗？', '提示', { type: 'warning' })
  await deleteCMDBHostApi(id)
  ElMessage.success('删除成功')
  await loadData()
}

const handleBatch = async (action: 'delete' | 'enable' | 'disable') => {
  if (!selectedIDs.value.length) {
    ElMessage.warning('请先选择主机')
    return
  }
  await batchCMDBHostsApi({ action, ids: selectedIDs.value })
  ElMessage.success('批量操作成功')
  await loadData()
}

const handleTest = async (id: number) => {
  await testCMDBHostSSHApi(id)
  ElMessage.success('SSH连接成功')
}

const openTerminal = (id: number) => {
  router.push(`/cmdb/terminal?host_id=${id}`)
}

onMounted(async () => {
  await loadGroups()
  await loadData()
})
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
.actions { display: flex; gap: 8px; }
.search { margin-bottom: 12px; }
.pager { margin-top: 12px; display: flex; justify-content: flex-end; }
</style>
