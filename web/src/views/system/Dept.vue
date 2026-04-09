<template>
  <TablePage
    class="system-dept-page"
    section-title="部门信息"
    section-subtitle="统一维护组织架构、层级关系与部门启停状态"
    wide
  >
    <template #header>
      <PageHeader
        eyebrow="Organization Map"
        title="部门信息"
        subtitle="按照公司、中心、部门三级结构维护组织节点与可见范围"
      >
        <template #actions>
          <el-button @click="refreshAll">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </template>
        <template #intro>
          <PageIntro
            title="治理说明"
            text="部门结构会影响管理员归属、数据权限与组织视图，建议按实际组织层级维护并保持状态准确。"
          />
        </template>
      </PageHeader>
    </template>

    <template #toolbar>
      <PageToolbar>
        <el-form :inline="true" :model="queryParams" class="filter-cluster">
          <el-form-item class="filter-field" label="部门名称">
            <el-input
              v-model="queryParams.deptName"
              placeholder="请输入部门名称"
              clearable
              @keyup.enter="handleQuery"
            />
          </el-form-item>
          <el-form-item class="filter-field" label="状态">
            <el-select v-model="queryParams.deptStatus" placeholder="全部状态" clearable>
              <el-option v-for="item in deptStatusList" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
        </el-form>
        <template #actions>
          <el-button type="primary" @click="handleQuery">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="resetQuery">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
          <el-button type="primary" plain @click="openAddDialog" v-authority="['base:dept:add']">
            <el-icon><Plus /></el-icon>
            新增部门
          </el-button>
          <el-button @click="toggleExpandAll">
            <el-icon><Sort /></el-icon>
            {{ isExpandAll ? '收起全部' : '展开全部' }}
          </el-button>
        </template>
      </PageToolbar>
    </template>

    <template #stats>
      <StatStrip :items="statItems" />
    </template>

    <el-table
      v-if="refreshTable"
      v-loading="loading"
      :data="deptList"
      :default-expand-all="isExpandAll"
      :tree-props="{ children: 'children' }"
      row-key="id"
      stripe
      class="data-table"
      empty-text="暂无部门数据"
    >
      <el-table-column label="部门名称" prop="deptName" min-width="220" show-overflow-tooltip />
      <el-table-column label="类型" width="120">
        <template #default="{ row }">
          <el-tag :type="deptTypeTagType(row.deptType)" effect="dark">{{ deptTypeText(row.deptType) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag :type="deptStatusTagType(row.deptStatus)" effect="dark">{{ deptStatusText(row.deptStatus) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" prop="createTime" min-width="170" />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <div class="row-actions">
            <el-button type="primary" link @click="showEditDeptDialog(row.id)" v-authority="['base:dept:edit']">
              编辑
            </el-button>
            <el-button
              type="danger"
              link
              :disabled="String(row.deptType) === '1'"
              @click="handleDeptDelete(row)"
              v-authority="['base:dept:delete']"
            >
              删除
            </el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="addDeptDialogVisible" title="新增部门" width="560px" @closed="addDeptDialogClosed">
      <el-form ref="addDeptFormRefForm" :model="addDeptForm" :rules="addDeptFormRules" label-width="88px" class="dept-form">
        <el-form-item label="部门类型" prop="deptType">
          <el-radio-group v-model="addDeptForm.deptType">
            <el-radio :label="1">公司</el-radio>
            <el-radio :label="2">中心</el-radio>
            <el-radio :label="3">部门</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="addDeptForm.deptType != 1" label="上级部门" prop="parentId">
          <treeselect v-model="addDeptForm.parentId" :options="optionsDeptList" placeholder="请选择上级部门" />
        </el-form-item>
        <el-form-item label="部门名称" prop="deptName">
          <el-input v-model="addDeptForm.deptName" placeholder="请输入部门名称" />
        </el-form-item>
        <el-form-item label="状态" prop="deptStatus">
          <el-radio-group v-model="addDeptForm.deptStatus">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="2">停用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="addDeptDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="addDept">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog v-model="editDeptDialogVisible" title="编辑部门" width="560px" @closed="editDeptDialogClosed">
      <el-form ref="editDeptFormRefForm" :model="deptInfo" :rules="editDeptFormRules" label-width="88px" class="dept-form">
        <el-form-item label="部门类型" prop="deptType">
          <el-radio-group v-model="deptInfo.deptType">
            <el-radio :label="1">公司</el-radio>
            <el-radio :label="2">中心</el-radio>
            <el-radio :label="3">部门</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="deptInfo.deptType != 1" label="上级部门" prop="parentId">
          <treeselect v-model="deptInfo.parentId" :options="optionsDeptList" placeholder="请选择上级部门" />
        </el-form-item>
        <el-form-item label="部门名称" prop="deptName">
          <el-input v-model="deptInfo.deptName" placeholder="请输入部门名称" />
        </el-form-item>
        <el-form-item label="状态" prop="deptStatus">
          <el-radio-group v-model="deptInfo.deptStatus">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="2">停用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="editDeptDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="editDept">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </TablePage>
</template>

<script>
import { ElMessage, ElMessageBox } from 'element-plus'
import Treeselect from 'vue3-treeselect'
import 'vue3-treeselect/dist/vue3-treeselect.css'
import {
  Search,
  Refresh,
  Plus,
  Sort
} from '@element-plus/icons-vue'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import StatStrip from '@/components/platform/StatStrip.vue'
import TablePage from '@/components/platform/TablePage.vue'

const createQueryParams = () => ({
  deptName: '',
  deptStatus: ''
})

const createDeptForm = () => ({
  deptType: undefined,
  parentId: undefined,
  deptName: '',
  deptStatus: 1
})

const flattenTree = list => {
  const nodes = []
  const walk = items => {
    (items || []).forEach(item => {
      nodes.push(item)
      if (Array.isArray(item.children) && item.children.length > 0) {
        walk(item.children)
      }
    })
  }
  walk(list)
  return nodes
}

const isActionCancelled = error => ['cancel', 'close', 'cancelled'].includes(error)

export default {
  name: 'SystemDeptPage',
  components: {
    PageHeader,
    PageIntro,
    PageToolbar,
    StatStrip,
    TablePage,
    Treeselect,
    Search,
    Refresh,
    Plus,
    Sort
  },
  data() {
    return {
      deptStatusList: [
        { value: '2', label: '停用' },
        { value: '1', label: '启用' }
      ],
      queryParams: createQueryParams(),
      loading: false,
      deptList: [],
      refreshTable: true,
      isExpandAll: true,
      optionsDeptList: [],
      addDeptDialogVisible: false,
      addDeptFormRules: {
        deptType: [{ required: true, message: '请选择部门类型', trigger: 'change' }],
        deptName: [{ required: true, message: '请输入部门名称', trigger: 'blur' }]
      },
      addDeptForm: createDeptForm(),
      editDeptDialogVisible: false,
      deptInfo: createDeptForm(),
      editDeptFormRules: {
        deptType: [{ required: true, message: '请选择部门类型', trigger: 'change' }],
        deptName: [{ required: true, message: '请输入部门名称', trigger: 'blur' }]
      }
    }
  },
  computed: {
    flatDeptList() {
      return flattenTree(this.deptList)
    },
    statItems() {
      const companyCount = this.flatDeptList.filter(item => Number(item.deptType) === 1).length
      const centerCount = this.flatDeptList.filter(item => Number(item.deptType) === 2).length
      const deptCount = this.flatDeptList.filter(item => Number(item.deptType) === 3).length
      const activeCount = this.flatDeptList.filter(item => Number(item.deptStatus) === 1).length

      return [
        { label: '组织节点', value: this.flatDeptList.length, hint: '当前组织节点总数', tone: 'primary' },
        { label: '公司', value: companyCount, hint: '一级组织节点', tone: 'neutral' },
        { label: '中心', value: centerCount, hint: '二级组织节点', tone: 'warning' },
        { label: '部门', value: deptCount, hint: '末级组织节点', tone: 'success' },
        { label: '启用部门', value: activeCount, hint: '当前启用状态', tone: 'success' }
      ]
    }
  },
  methods: {
    normalizeQueryParams() {
      const params = {}
      if (this.queryParams.deptName) {
        params.deptName = this.queryParams.deptName
      }
      if (this.queryParams.deptStatus) {
        params.deptStatus = this.queryParams.deptStatus
      }
      return params
    },
    deptTypeText(type) {
      return {
        1: '公司',
        2: '中心',
        3: '部门'
      }[Number(type)] || '未知'
    },
    deptTypeTagType(type) {
      return {
        1: 'info',
        2: 'warning',
        3: 'primary'
      }[Number(type)] || 'info'
    },
    deptStatusText(status) {
      return Number(status) === 1 ? '启用' : '停用'
    },
    deptStatusTagType(status) {
      return Number(status) === 1 ? 'success' : 'danger'
    },
    async getList() {
      this.loading = true
      try {
        const { data: res } = await this.$api.queryDeptList(this.normalizeQueryParams())
        if (res.code !== 200) {
          throw new Error(res.message || '获取部门列表失败')
        }

        this.deptList = this.$handleTree.handleTree(res.data || [], 'id')
      } catch (error) {
        ElMessage.error(error.message || '获取部门列表失败')
      } finally {
        this.loading = false
      }
    },
    async getDeptVoList() {
      try {
        const { data: res } = await this.$api.querySysDeptVoList()
        if (res.code !== 200) {
          throw new Error(res.message || '获取部门树失败')
        }

        this.optionsDeptList = this.$handleTree.handleTree(res.data || [], 'id')
      } catch (error) {
        ElMessage.error(error.message || '获取部门树失败')
      }
    },
    refreshAll() {
      Promise.all([this.getList(), this.getDeptVoList()])
    },
    handleQuery() {
      this.getList()
    },
    resetQuery() {
      this.queryParams = createQueryParams()
      this.getList()
    },
    toggleExpandAll() {
      this.refreshTable = false
      this.isExpandAll = !this.isExpandAll
      this.$nextTick(() => {
        this.refreshTable = true
      })
    },
    openAddDialog() {
      this.addDeptForm = createDeptForm()
      this.addDeptDialogVisible = true
      this.$nextTick(() => {
        this.$refs.addDeptFormRefForm?.clearValidate()
      })
    },
    addDeptDialogClosed() {
      this.$refs.addDeptFormRefForm?.resetFields()
      this.addDeptForm = createDeptForm()
    },
    async addDept() {
      const valid = await this.$refs.addDeptFormRefForm?.validate().catch(() => false)
      if (!valid) {
        return
      }

      try {
        const { data: res } = await this.$api.addDept(this.addDeptForm)
        if (res.code !== 200) {
          throw new Error(res.message || '新增部门失败')
        }

        ElMessage.success('新增部门成功')
        this.addDeptDialogVisible = false
        this.refreshAll()
      } catch (error) {
        ElMessage.error(error.message || '新增部门失败')
      }
    },
    async showEditDeptDialog(id) {
      try {
        const { data: res } = await this.$api.deptInfo(id)
        if (res.code !== 200) {
          throw new Error(res.message || '获取部门详情失败')
        }

        this.deptInfo = {
          ...createDeptForm(),
          ...(res.data || {})
        }
        this.editDeptDialogVisible = true
        this.$nextTick(() => {
          this.$refs.editDeptFormRefForm?.clearValidate()
        })
      } catch (error) {
        ElMessage.error(error.message || '获取部门详情失败')
      }
    },
    editDeptDialogClosed() {
      this.$refs.editDeptFormRefForm?.resetFields()
      this.deptInfo = createDeptForm()
    },
    async editDept() {
      const valid = await this.$refs.editDeptFormRefForm?.validate().catch(() => false)
      if (!valid) {
        return
      }

      try {
        const { data: res } = await this.$api.deptUpdate(this.deptInfo)
        if (res.code !== 200) {
          throw new Error(res.message || '更新部门失败')
        }

        ElMessage.success('更新部门成功')
        this.editDeptDialogVisible = false
        this.refreshAll()
      } catch (error) {
        ElMessage.error(error.message || '更新部门失败')
      }
    },
    async handleDeptDelete(row) {
      try {
        await ElMessageBox.confirm(`确认删除部门 ${row.deptName} 吗`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })

        const { data: res } = await this.$api.deleteDept(row.id)
        if (res.code !== 200) {
          throw new Error(res.message || '删除部门失败')
        }

        ElMessage.success('删除部门成功')
        this.refreshAll()
      } catch (error) {
        if (!isActionCancelled(error)) {
          ElMessage.error(error.message || '删除部门失败')
        }
      }
    }
  },
  created() {
    this.refreshAll()
  }
}
</script>

<style scoped>
.system-dept-page :deep(.page-actions) {
  align-items: center;
}

.system-dept-page :deep(.section-card__body) {
  display: grid;
  gap: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

@media (max-width: 960px) {
  .dialog-footer {
    width: 100%;
    justify-content: stretch;
  }

  .dialog-footer :deep(.el-button) {
    flex: 1 1 0;
  }
}
</style>
