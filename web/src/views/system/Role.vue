<template>
  <TablePage
    class="system-role-page"
    section-title="角色管理"
    section-subtitle="统一维护角色标识、状态与菜单权限范围"
    wide
  >
    <template #header>
      <PageHeader
        eyebrow="Identity Governance"
        title="角色管理"
        subtitle="通过角色承载权限集合、菜单访问范围与业务治理职责"
      >
        <template #actions>
          <el-button @click="refreshAll">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </template>
        <template #intro>
          <PageIntro
            title="权限说明"
            text="角色用于承接菜单权限与功能入口，建议按岗位职责拆分能力边界并保持标识稳定。"
          />
        </template>
      </PageHeader>
    </template>

    <template #toolbar>
      <PageToolbar>
        <el-form :inline="true" :model="queryParams" class="filter-cluster">
          <el-form-item class="filter-field" label="角色名称">
            <el-input
              v-model="queryParams.roleName"
              placeholder="请输入角色名称"
              clearable
              @keyup.enter="handleQuery"
            />
          </el-form-item>
          <el-form-item class="filter-field" label="状态">
            <el-select v-model="queryParams.status" placeholder="全部状态" clearable>
              <el-option v-for="item in statusList" :key="item.value" :label="item.label" :value="item.value" />
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
          <el-button type="primary" plain @click="openAddDialog" v-authority="['base:role:add']">
            <el-icon><Plus /></el-icon>
            新增角色
          </el-button>
        </template>
      </PageToolbar>
    </template>

    <template #stats>
      <StatStrip :items="statItems" />
    </template>

    <el-table v-loading="loading" :data="roleList" stripe class="data-table" empty-text="暂无角色数据">
      <el-table-column label="角色名称" prop="roleName" min-width="180" show-overflow-tooltip />
      <el-table-column label="角色标识" prop="roleKey" min-width="180" show-overflow-tooltip />
      <el-table-column label="创建时间" prop="createTime" min-width="170" />
      <el-table-column label="状态" width="150">
        <template #default="{ row }">
          <el-switch
            v-model="row.status"
            :active-value="1"
            :inactive-value="2"
            inline-prompt
            active-text="启用"
            inactive-text="停用"
            class="status-switch"
            @change="roleUpdateStatus(row)"
          />
        </template>
      </el-table-column>
      <el-table-column label="说明" prop="description" min-width="220" show-overflow-tooltip />
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <div class="row-actions">
            <el-button type="primary" link @click="showEditRoleDialog(row.id)" v-authority="['base:role:edit']">
                编辑
            </el-button>
            <el-button type="warning" link @click="showSetMenuDialog(row)" v-authority="['base:role:assign']">
                分配权限
            </el-button>
            <el-button type="danger" link @click="handleRoleDelete(row)" v-authority="['base:role:delete']">
                删除
            </el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <template #footer>
      <el-pagination
        :current-page="queryParams.pageNum"
        :page-sizes="[10, 50, 100, 500, 1000]"
        :page-size="queryParams.pageSize"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </template>

    <el-dialog
      v-model="addRoleDialogVisible"
      title="新增角色"
      width="560px"
      @closed="addRoleDialogClosed"
    >
      <el-form ref="addRoleFormRefForm" :model="addRoleForm" :rules="addRoleFormRules" label-width="88px" class="dialog-form">
        <el-form-item label="角色名称" prop="roleName">
          <el-input v-model="addRoleForm.roleName" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色标识" prop="roleKey">
          <el-input v-model="addRoleForm.roleKey" placeholder="请输入角色标识" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="addRoleForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="2">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="说明" prop="description">
          <el-input v-model="addRoleForm.description" type="textarea" :rows="4" placeholder="请输入角色说明" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="addRoleDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="addRole">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog
      v-model="editRoleDialogVisible"
      title="编辑角色"
      width="560px"
      @closed="editRoleDialogClosed"
    >
      <el-form ref="editRoleFormRefForm" :model="roleInfo" :rules="editRoleFormRules" label-width="88px" class="dialog-form">
        <el-form-item label="角色名称" prop="roleName">
          <el-input v-model="roleInfo.roleName" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色标识" prop="roleKey">
          <el-input v-model="roleInfo.roleKey" placeholder="请输入角色标识" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="roleInfo.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="2">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="说明" prop="description">
          <el-input v-model="roleInfo.description" type="textarea" :rows="4" placeholder="请输入角色说明" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="editRoleDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="editRole">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog
      v-model="setMenuDialogVisible"
      :title="currentRoleName ? `分配权限 - ${currentRoleName}` : '分配权限'"
      width="460px"
      @closed="setRightDialogClosed"
    >
      <div v-loading="permissionLoading" class="permission-dialog">
        <el-alert
          title="勾选该角色可访问的菜单与功能按钮"
          type="info"
          :closable="false"
        />
        <div class="permission-dialog__tree">
          <el-tree
            ref="treeRef"
            :key="permissionTreeKey"
            :data="menuList"
            :props="treeProps"
            :default-checked-keys="defKeys"
            node-key="id"
            show-checkbox
            class="permission-tree"
            @check="handleTreeCheck"
          />
        </div>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="setMenuDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="allotMenus">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </TablePage>
</template>

<script>
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Plus
} from '@element-plus/icons-vue'
import PageHeader from '@/components/platform/PageHeader.vue'
import PageIntro from '@/components/platform/PageIntro.vue'
import PageToolbar from '@/components/platform/PageToolbar.vue'
import StatStrip from '@/components/platform/StatStrip.vue'
import TablePage from '@/components/platform/TablePage.vue'

const createQueryParams = () => ({
  roleName: '',
  status: '',
  pageNum: 1,
  pageSize: 10
})

const createRoleForm = () => ({
  roleName: '',
  roleKey: '',
  description: '',
  status: 1
})

const isActionCancelled = error => ['cancel', 'close', 'cancelled'].includes(error)

export default {
  name: 'SystemRolePage',
  components: {
    PageHeader,
    PageIntro,
    PageToolbar,
    StatStrip,
    TablePage,
    Search,
    Refresh,
    Plus
  },
  data() {
    return {
      statusList: [
        { value: '1', label: '启用' },
        { value: '2', label: '停用' }
      ],
      queryParams: createQueryParams(),
      loading: false,
      roleList: [],
      total: 0,
      addRoleDialogVisible: false,
      addRoleForm: createRoleForm(),
      addRoleFormRules: {
        roleName: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
        roleKey: [{ required: true, message: '请输入角色标识', trigger: 'blur' }],
        status: [{ required: true, message: '请选择角色状态', trigger: 'change' }],
        description: [{ required: true, message: '请输入角色说明', trigger: 'blur' }]
      },
      editRoleDialogVisible: false,
      roleInfo: createRoleForm(),
      editRoleFormRules: {
        roleName: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
        roleKey: [{ required: true, message: '请输入角色标识', trigger: 'blur' }],
        status: [{ required: true, message: '请选择角色状态', trigger: 'change' }],
        description: [{ required: true, message: '请输入角色说明', trigger: 'blur' }]
      },
      setMenuDialogVisible: false,
      permissionLoading: false,
      menuList: [],
      treeProps: {
        label: 'label'
      },
      defKeys: [],
      permissionTreeKey: 0,
      currentRoleId: '',
      currentRoleName: ''
    }
  },
  computed: {
    statItems() {
      const enabledCount = this.roleList.filter(item => Number(item.status) === 1).length
      const disabledCount = this.roleList.filter(item => Number(item.status) === 2).length
      const describedCount = this.roleList.filter(item => item.description).length

      return [
        { label: '角色总数', value: this.total, hint: '当前角色覆盖数量', tone: 'primary' },
        { label: '启用角色', value: enabledCount, hint: '当前启用角色', tone: 'success' },
        { label: '停用角色', value: disabledCount, hint: '当前停用角色', tone: 'danger' },
        { label: '已写说明', value: describedCount, hint: '说明完善度', tone: 'warning' }
      ]
    }
  },
  methods: {
    normalizeQueryParams() {
      const params = {
        pageNum: this.queryParams.pageNum,
        pageSize: this.queryParams.pageSize
      }

      if (this.queryParams.roleName) {
        params.roleName = this.queryParams.roleName
      }
      if (this.queryParams.status) {
        params.status = this.queryParams.status
      }

      return params
    },
    filterLeafPermissions(menuTree, rolePermissions) {
      const leafPermissions = []

      const collectLeafNodes = nodes => {
        nodes.forEach(node => {
          if (!node.children || node.children.length === 0) {
            leafPermissions.push(node.id)
            return
          }
          collectLeafNodes(node.children)
        })
      }

      collectLeafNodes(menuTree || [])
      return (rolePermissions || []).filter(id => leafPermissions.includes(id))
    },
    sortMenusByOrder(menuList, menuWithSort) {
      if (!Array.isArray(menuList)) {
        return []
      }

      const sortMap = {}
      if (Array.isArray(menuWithSort)) {
        menuWithSort.forEach(menu => {
          sortMap[menu.id] = menu.sort || 0
        })
      }

      const sortedList = [...menuList].sort((a, b) => {
        const sortA = sortMap[a.id] || 999
        const sortB = sortMap[b.id] || 999
        return sortA - sortB
      })

      return sortedList.map(menu => {
        if (menu.children && menu.children.length > 0) {
          return {
            ...menu,
            children: this.sortMenusByOrder(menu.children, menuWithSort)
          }
        }
        return menu
      })
    },
    async getRoleList() {
      this.loading = true
      try {
        const { data: res } = await this.$api.queryRoleList(this.normalizeQueryParams())
        if (res.code !== 200) {
          throw new Error(res.message || '获取角色列表失败')
        }

        const payload = res.data || {}
        this.roleList = payload.list || []
        this.total = payload.total || 0
      } catch (error) {
        ElMessage.error(error.message || '获取角色列表失败')
      } finally {
        this.loading = false
      }
    },
    refreshAll() {
      this.getRoleList()
    },
    handleQuery() {
      this.queryParams.pageNum = 1
      this.getRoleList()
    },
    resetQuery() {
      this.queryParams = createQueryParams()
      this.getRoleList()
    },
    handleSizeChange(newSize) {
      this.queryParams.pageSize = newSize
      this.queryParams.pageNum = 1
      this.getRoleList()
    },
    handleCurrentChange(newPage) {
      this.queryParams.pageNum = newPage
      this.getRoleList()
    },
    openAddDialog() {
      this.addRoleForm = createRoleForm()
      this.addRoleDialogVisible = true
      this.$nextTick(() => {
        this.$refs.addRoleFormRefForm?.clearValidate()
      })
    },
    async roleUpdateStatus(row) {
      const actionText = Number(row.status) === 2 ? '停用' : '启用'

      try {
        await ElMessageBox.confirm(`确认${actionText}角色 ${row.roleName} 吗`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })

        const { data: res } = await this.$api.updateRoleStatus(row.id, row.status)
        if (res.code !== 200) {
          throw new Error(res.message || `${actionText}失败`)
        }

        ElMessage.success(`${actionText}成功`)
      } catch (error) {
        if (!isActionCancelled(error)) {
          ElMessage.error(error.message || `${actionText}失败`)
        }
      } finally {
        this.getRoleList()
      }
    },
    addRoleDialogClosed() {
      this.$refs.addRoleFormRefForm?.resetFields()
      this.addRoleForm = createRoleForm()
    },
    async addRole() {
      const valid = await this.$refs.addRoleFormRefForm?.validate().catch(() => false)
      if (!valid) {
        return
      }

      try {
        const { data: res } = await this.$api.addRole(this.addRoleForm)
        if (res.code !== 200) {
          throw new Error(res.message || '新增角色失败')
        }

        ElMessage.success('新增角色成功')
        this.addRoleDialogVisible = false
        this.getRoleList()
      } catch (error) {
        ElMessage.error(error.message || '新增角色失败')
      }
    },
    editRoleDialogClosed() {
      this.$refs.editRoleFormRefForm?.resetFields()
      this.roleInfo = createRoleForm()
    },
    async showEditRoleDialog(id) {
      try {
        const { data: res } = await this.$api.roleInfo(id)
        if (res.code !== 200) {
          throw new Error(res.message || '获取角色详情失败')
        }

        this.roleInfo = {
          ...createRoleForm(),
          ...(res.data || {})
        }
        this.editRoleDialogVisible = true
        this.$nextTick(() => {
          this.$refs.editRoleFormRefForm?.clearValidate()
        })
      } catch (error) {
        ElMessage.error(error.message || '获取角色详情失败')
      }
    },
    async editRole() {
      const valid = await this.$refs.editRoleFormRefForm?.validate().catch(() => false)
      if (!valid) {
        return
      }

      try {
        const { data: res } = await this.$api.roleUpdate(this.roleInfo)
        if (res.code !== 200) {
          throw new Error(res.message || '更新角色失败')
        }

        ElMessage.success('更新角色成功')
        this.editRoleDialogVisible = false
        this.getRoleList()
      } catch (error) {
        ElMessage.error(error.message || '更新角色失败')
      }
    },
    async handleRoleDelete(row) {
      try {
        await ElMessageBox.confirm(`确认删除角色 ${row.roleName} 吗`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })

        const { data: res } = await this.$api.deleteRole(row.id)
        if (res.code !== 200) {
          throw new Error(res.message || '删除角色失败')
        }

        ElMessage.success('删除成功')
        this.getRoleList()
      } catch (error) {
        if (!isActionCancelled(error)) {
          ElMessage.error(error.message || '删除失败')
        }
      }
    },
    async showSetMenuDialog(role) {
      this.currentRoleId = role.id
      this.currentRoleName = role.roleName || ''
      this.setMenuDialogVisible = true
      this.permissionLoading = true
      this.permissionTreeKey += 1

      try {
        const [menuRes, roleMenuRes, menuWithSortRes] = await Promise.all([
          this.$api.querySysMenuVoList(),
          this.$api.QueryRoleMenuIdList(role.id),
          this.$api.queryMenuList({})
        ])

        if (menuRes.data.code !== 200) {
          throw new Error(menuRes.data.message || '获取菜单树失败')
        }
        if (roleMenuRes.data.code !== 200) {
          throw new Error(roleMenuRes.data.message || '获取角色权限失败')
        }
        if (menuWithSortRes.data.code !== 200) {
          throw new Error(menuWithSortRes.data.message || '获取菜单排序失败')
        }

        const allMenus = menuRes.data.data || []
        const rolePermissions = roleMenuRes.data.data || []
        const menuWithSort = menuWithSortRes.data.data || []
        const sortedMenus = this.sortMenusByOrder(allMenus, menuWithSort)

        this.menuList = this.$handleTree.handleTree(sortedMenus, 'id')
        this.defKeys = this.filterLeafPermissions(this.menuList, rolePermissions)
      } catch (error) {
        ElMessage.error(error.message || '加载权限数据失败')
      } finally {
        this.permissionLoading = false
      }
    },
    setRightDialogClosed() {
      this.defKeys = []
      this.menuList = []
      this.permissionTreeKey += 1
      this.currentRoleId = ''
      this.currentRoleName = ''
    },
    handleTreeCheck() {},
    async allotMenus() {
      if (!this.currentRoleId) {
        return
      }

      try {
        const checkedKeys = this.$refs.treeRef?.getCheckedKeys() || []
        const halfCheckedKeys = this.$refs.treeRef?.getHalfCheckedKeys() || []
        const allPermissionIds = [...new Set([...checkedKeys, ...halfCheckedKeys])]

        const { data: res } = await this.$api.AssignPermissions(this.currentRoleId, allPermissionIds)
        if (res.code !== 200) {
          throw new Error(res.message || '分配权限失败')
        }

        ElMessage.success('分配权限成功')
        this.setMenuDialogVisible = false
        this.getRoleList()
      } catch (error) {
        ElMessage.error(error.message || '分配权限失败')
      }
    }
  },
  created() {
    this.getRoleList()
  }
}
</script>

<style scoped>
.system-role-page :deep(.page-actions) {
  align-items: center;
}

.system-role-page :deep(.section-card__body) {
  display: grid;
  gap: 20px;
}

.status-switch {
  --el-switch-on-color: #2563eb;
  --el-switch-off-color: rgba(148, 163, 184, 0.3);
}

.dialog-form :deep(.el-textarea__inner) {
  min-height: 120px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.permission-dialog {
  display: grid;
  gap: 16px;
}

.permission-dialog__tree {
  max-height: 420px;
  overflow: auto;
  padding: 14px;
  border-radius: 18px;
  border: 1px solid rgba(148, 163, 184, 0.14);
  background: rgba(255, 255, 255, 0.03);
}

.permission-tree {
  background: transparent;
}

.permission-tree :deep(.el-tree-node__content) {
  height: 34px;
  border-radius: 10px;
}

.permission-tree :deep(.el-tree-node__content:hover),
.permission-tree :deep(.el-tree-node:focus > .el-tree-node__content) {
  background: rgba(59, 130, 246, 0.1);
}

.permission-tree :deep(.el-tree-node__label) {
  color: var(--text-secondary);
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
