<template>
  <el-dialog
    v-model="visible"
    class="host-selector-dialog"
    title="选择执行主机"
    width="960px"
    destroy-on-close
  >
    <div class="host-selector-layout">
      <aside class="host-selector-sidebar">
        <div class="host-selector-sidebar__title">主机分组</div>
        <div class="host-selector-sidebar__subtitle">点击左侧分组，快速切换目标主机范围。</div>
        <el-tree
          :data="groupsWithHosts"
          :props="treeProps"
          node-key="id"
          highlight-current
          class="host-selector-tree"
          @node-click="handleGroupClick"
        >
          <template #default="{ node }">
            <div class="host-tree-node">
              <el-icon><Folder /></el-icon>
              <span>{{ node.label }}</span>
            </div>
          </template>
        </el-tree>
      </aside>

      <section class="host-selector-main">
        <div class="host-selector-main__header">
          <div>
            <div class="host-selector-main__title">主机列表</div>
            <div class="host-selector-main__subtitle">支持分页浏览并勾选多台主机，一次性加入当前任务。</div>
          </div>
          <span class="host-selector-main__meta">共 {{ pagination.total }} 台</span>
        </div>

        <el-table
          ref="hostTable"
          :data="currentGroupHosts"
          class="host-selector-table"
          @selection-change="handleHostSelectionChange"
        >
          <el-table-column type="selection" width="52" :selectable="() => true" />
          <el-table-column prop="name" label="主机名称" min-width="160" show-overflow-tooltip>
            <template #default="{ row }">
              <div class="host-cell host-cell--name">
                <el-icon><User /></el-icon>
                <span>{{ row.name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="privateIp" label="内网 IP" min-width="140">
            <template #default="{ row }">
              <div class="host-cell">
                <el-icon><HomeFilled /></el-icon>
                <span>{{ row.privateIp || '-' }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="publicIp" label="外网 IP" min-width="140">
            <template #default="{ row }">
              <div class="host-cell">
                <el-icon><Connection /></el-icon>
                <span>{{ row.publicIp || '-' }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="os" label="操作系统" min-width="180" show-overflow-tooltip />
          <el-table-column prop="remark" label="备注" min-width="180" show-overflow-tooltip />
        </el-table>

        <div class="pagination">
          <el-pagination
            :current-page="pagination.pageNum"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="pagination.pageSize"
            layout="total, sizes, prev, pager, next, jumper"
            :total="pagination.total"
            @current-change="handlePageChange"
            @size-change="handleSizeChange"
          />
        </div>
      </section>
    </div>
    <template #footer>
      <el-button @click="close">关闭</el-button>
      <el-button type="primary" @click="confirmSelection">添加选中主机</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import cmdbAPI from '@/api/cmdb'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  selectedHosts: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue', 'hosts-selected'])

const visible = ref(false)
const groupsWithHosts = ref([])
const currentGroupHosts = ref([])
const tempSelectedHosts = ref([])
const hostTable = ref(null)
const currentGroupId = ref(null)
const pagination = ref({
  pageNum: 1,
  pageSize: 10,
  total: 0
})

const treeProps = {
  label: 'name',
  children: 'children'
}

const normalizeHost = host => ({
  id: host.id,
  name: host.hostName || host.name,
  privateIp: host.privateIp,
  publicIp: host.publicIp,
  os: host.os,
  remark: host.remark
})

const fetchGroupsWithHosts = async () => {
  try {
    const response = await cmdbAPI.getGroupListWithHosts()
    if (response.data?.code === 200) {
      groupsWithHosts.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取分组和主机列表失败:', error)
  }
}

const getAllHosts = async () => {
  try {
    const response = await cmdbAPI.getCmdbHostList({
      page: pagination.value.pageNum,
      pageSize: pagination.value.pageSize
    })

    if (response.data?.code === 200) {
      currentGroupHosts.value = (response.data.data?.list || []).map(normalizeHost)
      pagination.value.total = response.data.data?.total || 0
    }
  } catch (error) {
    console.error('获取主机列表失败:', error)
    currentGroupHosts.value = []
  }
}

const handleGroupClick = async data => {
  currentGroupId.value = data.id
  try {
    const response = await cmdbAPI.getCmdbHostsByGroupId(data.id, {
      page: pagination.value.pageNum,
      pageSize: pagination.value.pageSize
    })

    if (response.data?.code === 200) {
      const hostList = response.data.data || []
      currentGroupHosts.value = hostList.map(normalizeHost)
      pagination.value.total = hostList.length
    }
  } catch (error) {
    console.error('获取主机列表失败:', error)
    currentGroupHosts.value = []
  }
}

const handlePageChange = page => {
  pagination.value.pageNum = page
  if (currentGroupId.value) {
    handleGroupClick({ id: currentGroupId.value })
    return
  }
  getAllHosts()
}

const handleSizeChange = size => {
  pagination.value.pageSize = size
  pagination.value.pageNum = 1
  if (currentGroupId.value) {
    handleGroupClick({ id: currentGroupId.value })
    return
  }
  getAllHosts()
}

const handleHostSelectionChange = selection => {
  tempSelectedHosts.value = selection
}

const confirmSelection = () => {
  emit('hosts-selected', tempSelectedHosts.value)
  tempSelectedHosts.value = []
}

const close = () => {
  visible.value = false
}

watch(
  () => props.modelValue,
  value => {
    visible.value = value
  }
)

watch(visible, value => {
  emit('update:modelValue', value)
})

onMounted(() => {
  fetchGroupsWithHosts()
  getAllHosts()
})
</script>

<style scoped>
.host-selector-layout {
  display: grid;
  grid-template-columns: 260px minmax(0, 1fr);
  gap: 16px;
  min-height: 520px;
}

.host-selector-sidebar,
.host-selector-main {
  min-width: 0;
  border-radius: 18px;
  border: 1px solid var(--border-subtle);
  background: rgba(255, 255, 255, 0.04);
}

.host-selector-sidebar {
  padding: 16px;
  display: grid;
  align-content: start;
  gap: 8px;
}

.host-selector-sidebar__title,
.host-selector-main__title {
  color: var(--text-primary);
  font-size: 16px;
  font-weight: 700;
}

.host-selector-sidebar__subtitle,
.host-selector-main__subtitle,
.host-selector-main__meta {
  color: var(--text-muted);
  font-size: 12px;
  line-height: 1.6;
}

.host-selector-tree {
  margin-top: 10px;
}

.host-tree-node {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.host-selector-main {
  padding: 16px;
  display: grid;
  gap: 14px;
}

.host-selector-main__header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: flex-start;
}

.host-cell {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.host-cell--name {
  color: var(--text-primary);
  font-weight: 600;
}

.pagination {
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 900px) {
  .host-selector-layout {
    grid-template-columns: 1fr;
  }

  .host-selector-main__header {
    flex-direction: column;
  }
}
</style>
