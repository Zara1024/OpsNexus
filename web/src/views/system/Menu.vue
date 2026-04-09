<template>
  <TablePage
    class="system-menu-page"
    section-title="菜单管理"
    section-subtitle="统一维护目录、菜单、按钮与平台导航结构"
    wide
  >
    <template #header>
      <PageHeader
        eyebrow="Navigation Schema"
        title="菜单管理"
        subtitle="按目录层级整理导航入口、权限标识与路由映射关系"
      >
        <template #actions>
          <el-button @click="refreshAll">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </template>
        <template #intro>
          <PageIntro
            title="结构说明"
            text="菜单配置决定侧边导航与权限树展示，建议优先维护目录结构，再补齐路由、标识与图标。"
          />
        </template>
      </PageHeader>
    </template>

    <template #toolbar>
      <PageToolbar>
        <el-form :inline="true" :model="queryParams" class="filter-cluster">
          <el-form-item class="filter-field" label="菜单名称">
            <el-input
              v-model="queryParams.menuName"
              placeholder="请输入菜单名称"
              clearable
              @keyup.enter="handleQuery"
            />
          </el-form-item>
          <el-form-item class="filter-field" label="状态">
            <el-select v-model="queryParams.menuStatus" placeholder="全部状态" clearable>
              <el-option v-for="item in menuStatusList" :key="item.value" :label="item.label" :value="item.value" />
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
          <el-button type="primary" plain @click="openAddDialog" v-authority="['base:menu:add']">
            <el-icon><Plus /></el-icon>
            新增菜单
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
      :data="menuList"
      :default-expand-all="isExpandAll"
      :tree-props="{ children: 'children' }"
      row-key="id"
      stripe
      class="data-table"
      empty-text="暂无菜单数据"
    >
      <el-table-column prop="menuName" label="菜单名称" min-width="180" show-overflow-tooltip />
      <el-table-column prop="sort" label="排序" width="80" />
      <el-table-column prop="value" label="权限标识" min-width="150" show-overflow-tooltip />
      <el-table-column prop="url" label="路由地址" min-width="180" show-overflow-tooltip />
      <el-table-column label="类型" width="110">
        <template #default="{ row }">
          <el-tag :type="menuTypeTagType(row.menuType)" effect="dark">{{ menuTypeText(row.menuType) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="110">
        <template #default="{ row }">
          <el-tag :type="menuStatusTagType(row.menuStatus)" effect="dark">{{ menuStatusText(row.menuStatus) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <div class="row-actions">
            <el-button type="primary" link @click="showEditMenuDialog(row.id)" v-authority="['base:menu:edit']">
                编辑
            </el-button>
            <el-button type="warning" link @click="handleCopyMenu(row)" v-authority="['base:menu:add']">
                复制
            </el-button>
            <el-button type="danger" link @click="handleMenuDelete(row)" v-authority="['base:menu:delete']">
                删除
            </el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      v-model="addMenuDialogVisible"
      title="新增菜单"
      width="720px"
      class="modern-dialog ops-dialog ops-overlay--md system-menu-dialog"
      @closed="addMenuDialogClosed"
    >
      <el-form ref="addMenuFormRefForm" :model="menuForm" :rules="addMenuFormRules" label-width="92px" class="menu-form">
        <el-row :gutter="16">
          <el-col :span="24">
            <el-form-item label="菜单类型" prop="menuType">
              <el-radio-group v-model="menuForm.menuType">
                <el-radio :label="1">目录</el-radio>
                <el-radio :label="2">菜单</el-radio>
                <el-radio :label="3">按钮</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col v-if="menuForm.menuType != 1" :span="24">
            <el-form-item label="上级菜单" prop="parentId">
              <treeselect v-model="menuForm.parentId" :options="treeList" placeholder="请选择上级菜单" />
            </el-form-item>
          </el-col>
          <el-col v-if="menuForm.menuType != 3" :span="24">
            <el-form-item label="图标" prop="icon">
              <el-select v-model="menuForm.icon" placeholder="请选择图标" clearable>
                <el-option v-for="item in iconList" :key="item.value" :label="item.label" :value="item.value">
                  <div class="icon-option">
                    <el-icon>
                      <component :is="item.value" />
                    </el-icon>
                    <span>{{ item.label }}</span>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="菜单名称" prop="menuName">
              <el-input v-model="menuForm.menuName" placeholder="请输入菜单名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序" prop="sort">
              <el-input-number v-model="menuForm.sort" controls-position="right" :min="0" />
            </el-form-item>
          </el-col>
          <el-col v-if="menuForm.menuType != 3" :span="24">
            <el-form-item label="路由地址" prop="url">
              <el-input v-model="menuForm.url" placeholder="请输入路由地址" />
            </el-form-item>
          </el-col>
          <el-col v-if="menuForm.menuType != 1" :span="24">
            <el-form-item label="权限标识" prop="value">
              <el-input v-model="menuForm.value" placeholder="请输入权限标识" maxlength="50" />
            </el-form-item>
          </el-col>
          <el-col v-if="menuForm.menuType != 3" :span="24">
            <el-form-item label="状态" prop="menuStatus">
              <el-radio-group v-model="menuForm.menuStatus">
                <el-radio :label="1">隐藏</el-radio>
                <el-radio :label="2">显示</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="addMenuDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="addMenu">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog
      v-model="editMenuDialogVisible"
      title="编辑菜单"
      width="720px"
      class="modern-dialog ops-dialog ops-overlay--md system-menu-dialog"
      @closed="editMenuDialogClosed"
    >
      <el-form ref="editMenuFormRefForm" :model="menuInfo" :rules="editMenuFormRules" label-width="92px" class="menu-form">
        <el-row :gutter="16">
          <el-col :span="24">
            <el-form-item label="菜单类型" prop="menuType">
              <el-radio-group v-model="menuInfo.menuType">
                <el-radio :label="1">目录</el-radio>
                <el-radio :label="2">菜单</el-radio>
                <el-radio :label="3">按钮</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col v-if="menuInfo.menuType != 1" :span="24">
            <el-form-item label="上级菜单" prop="parentId">
              <treeselect v-model="menuInfo.parentId" :options="treeList" placeholder="请选择上级菜单" />
            </el-form-item>
          </el-col>
          <el-col v-if="menuInfo.menuType != 3" :span="24">
            <el-form-item label="图标" prop="icon">
              <el-select v-model="menuInfo.icon" placeholder="请选择图标" clearable>
                <el-option v-for="item in iconList" :key="item.value" :label="item.label" :value="item.value">
                  <div class="icon-option">
                    <el-icon>
                      <component :is="item.value" />
                    </el-icon>
                    <span>{{ item.label }}</span>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="菜单名称" prop="menuName">
              <el-input v-model="menuInfo.menuName" placeholder="请输入菜单名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序" prop="sort">
              <el-input-number v-model="menuInfo.sort" controls-position="right" :min="0" />
            </el-form-item>
          </el-col>
          <el-col v-if="menuInfo.menuType != 3" :span="24">
            <el-form-item label="路由地址" prop="url">
              <el-input v-model="menuInfo.url" placeholder="请输入路由地址" />
            </el-form-item>
          </el-col>
          <el-col v-if="menuInfo.menuType != 1" :span="24">
            <el-form-item label="权限标识" prop="value">
              <el-input v-model="menuInfo.value" placeholder="请输入权限标识" maxlength="50" />
            </el-form-item>
          </el-col>
          <el-col v-if="menuInfo.menuType != 3" :span="24">
            <el-form-item label="状态" prop="menuStatus">
              <el-radio-group v-model="menuInfo.menuStatus">
                <el-radio :label="1">隐藏</el-radio>
                <el-radio :label="2">显示</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="editMenuDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="editMenu">确定</el-button>
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
  menuName: '',
  menuStatus: ''
})

const createMenuForm = () => ({
  parentId: undefined,
  menuName: '',
  icon: '',
  value: '',
  menuType: 1,
  url: '',
  sort: 0,
  menuStatus: 2
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
  name: 'SystemMenuPage',
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
      queryParams: createQueryParams(),
      menuStatusList: [
        { value: '2', label: '显示' },
        { value: '1', label: '隐藏' }
      ],
      loading: false,
      menuList: [],
      isExpandAll: false,
      refreshTable: true,
      iconList: [
        { value: 'HomeFilled', label: 'HomeFilled' },
        { value: 'UploadFilled', label: 'UploadFilled' },
        { value: 'Menu', label: 'Menu' },
        { value: 'Search', label: 'Search' },
        { value: 'Edit', label: 'Edit' },
        { value: 'Delete', label: 'Delete' },
        { value: 'More', label: 'More' },
        { value: 'Star', label: 'Star' },
        { value: 'StarFilled', label: 'StarFilled' },
        { value: 'Platform', label: 'Platform' },
        { value: 'TrendCharts', label: 'TrendCharts' },
        { value: 'Document', label: 'Document' },
        { value: 'Eleme', label: 'Eleme' },
        { value: 'Delete', label: 'Delete' },
        { value: 'Tools', label: 'Tools' },
        { value: 'Setting', label: 'Setting' },
        { value: 'User', label: 'User' },
        { value: 'Phone', label: 'Phone' },
        { value: 'Goods', label: 'Goods' },
        { value: 'Help', label: 'Help' },
        { value: 'Picture', label: 'Picture' },
        { value: 'Upload', label: 'Upload' },
        { value: 'Download', label: 'Download' },
        { value: 'Promotion', label: 'Promotion' },
        { value: 'Shop', label: 'Shop' },
        { value: 'menu', label: 'Menu' },
        { value: 'share', label: 'Share' },
        { value: 'bottom', label: 'Bottom' },
        { value: 'top', label: 'Top' },
        { value: 'key', label: 'Key' },
        { value: 'unlock', label: 'Unlock' },
        { value: 'shopping-cart-full', label: 'ShoppingCartFull' },
        { value: 'Coin', label: 'Coin' },
        { value: 'present', label: 'Present' },
        { value: 'box', label: 'Box' },
        { value: 'wallet', label: 'Wallet' },
        { value: 'discount', label: 'Discount' },
        { value: 'price-tag', label: 'PriceTag' },
        { value: 'guide', label: 'Guide' },
        { value: 'connection', label: 'Connection' },
        { value: 'chat-dot-round', label: 'ChatDotRound' }
      ],
      addMenuDialogVisible: false,
      menuForm: createMenuForm(),
      addMenuFormRules: {
        menuType: [{ required: true, message: '请选择菜单类型', trigger: 'change' }],
        menuName: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
        sort: [{ required: true, message: '请输入排序值', trigger: 'blur' }],
        value: [{ required: true, message: '请输入权限标识', trigger: 'blur' }]
      },
      treeList: [],
      editMenuDialogVisible: false,
      menuInfo: createMenuForm(),
      editMenuFormRules: {
        menuType: [{ required: true, message: '请选择菜单类型', trigger: 'change' }],
        menuName: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
        sort: [{ required: true, message: '请输入排序值', trigger: 'blur' }],
        value: [{ required: true, message: '请输入权限标识', trigger: 'blur' }]
      }
    }
  },
  computed: {
    flatMenuList() {
      return flattenTree(this.menuList)
    },
    statItems() {
      const directoryCount = this.flatMenuList.filter(item => Number(item.menuType) === 1).length
      const menuCount = this.flatMenuList.filter(item => Number(item.menuType) === 2).length
      const buttonCount = this.flatMenuList.filter(item => Number(item.menuType) === 3).length
      const enabledCount = this.flatMenuList.filter(item => Number(item.menuStatus) === 2).length

      return [
        { label: '菜单节点', value: this.flatMenuList.length, hint: '当前导航节点', tone: 'primary' },
        { label: '目录', value: directoryCount, hint: '一级导航目录', tone: 'neutral' },
        { label: '菜单', value: menuCount, hint: '可跳转页面', tone: 'success' },
        { label: '按钮', value: buttonCount, hint: '按钮权限项', tone: 'warning' },
        { label: '显示中', value: enabledCount, hint: '当前显示状态', tone: 'success' }
      ]
    }
  },
  methods: {
    normalizeQueryParams() {
      const params = {}
      if (this.queryParams.menuName) {
        params.menuName = this.queryParams.menuName
      }
      if (this.queryParams.menuStatus) {
        params.menuStatus = this.queryParams.menuStatus
      }
      return params
    },
    normalizeMenuPayload(payload) {
      const submitData = {
        ...payload
      }
      delete submitData.children
      delete submitData.createTime
      delete submitData.updateTime
      return submitData
    },
    menuTypeText(type) {
      return {
        1: '目录',
        2: '菜单',
        3: '按钮'
      }[Number(type)] || '未知'
    },
    menuTypeTagType(type) {
      return {
        1: 'info',
        2: 'primary',
        3: 'warning'
      }[Number(type)] || 'info'
    },
    menuStatusText(status) {
      return Number(status) === 2 ? '显示' : '隐藏'
    },
    menuStatusTagType(status) {
      return Number(status) === 2 ? 'success' : 'info'
    },
    async getMenuList() {
      this.loading = true
      try {
        const { data: res } = await this.$api.queryMenuList(this.normalizeQueryParams())
        if (res.code !== 200) {
          throw new Error(res.message || '获取菜单列表失败')
        }

        this.menuList = this.$handleTree.handleTree(res.data || [], 'id')
      } catch (error) {
        ElMessage.error(error.message || '获取菜单列表失败')
      } finally {
        this.loading = false
      }
    },
    async getMenuVoList() {
      try {
        const [menuVoRes, menuWithSortRes] = await Promise.all([
          this.$api.querySysMenuVoList(),
          this.$api.queryMenuList({})
        ])

        if (menuVoRes.data.code !== 200) {
          throw new Error(menuVoRes.data.message || '获取菜单树失败')
        }
        if (menuWithSortRes.data.code !== 200) {
          throw new Error(menuWithSortRes.data.message || '获取菜单排序失败')
        }

        const menuVoData = menuVoRes.data.data || []
        const menuWithSort = menuWithSortRes.data.data || []
        const sortMap = {}

        menuWithSort.forEach(menu => {
          sortMap[menu.id] = menu.sort || 0
        })

        const sortedMenus = [...menuVoData]
          .map(menu => ({
            ...menu,
            sort: sortMap[menu.id] || 999
          }))
          .sort((a, b) => a.sort - b.sort)

        this.treeList = this.$handleTree.handleTree(sortedMenus, 'id')
      } catch (error) {
        ElMessage.error(error.message || '获取菜单树失败')
      }
    },
    refreshAll() {
      Promise.all([this.getMenuList(), this.getMenuVoList()])
    },
    handleQuery() {
      this.getMenuList()
    },
    resetQuery() {
      this.queryParams = createQueryParams()
      this.getMenuList()
    },
    toggleExpandAll() {
      this.refreshTable = false
      this.isExpandAll = !this.isExpandAll
      this.$nextTick(() => {
        this.refreshTable = true
      })
    },
    openAddDialog() {
      this.menuForm = createMenuForm()
      this.addMenuDialogVisible = true
      this.$nextTick(() => {
        this.$refs.addMenuFormRefForm?.clearValidate()
      })
    },
    addMenuDialogClosed() {
      this.$refs.addMenuFormRefForm?.resetFields()
      this.menuForm = createMenuForm()
    },
    async addMenu() {
      const valid = await this.$refs.addMenuFormRefForm?.validate().catch(() => false)
      if (!valid) {
        return
      }

      try {
        const submitData = this.normalizeMenuPayload(this.menuForm)
        const { data: res } = await this.$api.addMenu(submitData)
        if (res.code !== 200) {
          throw new Error(res.message || '新增菜单失败')
        }

        ElMessage.success('新增菜单成功')
        this.addMenuDialogVisible = false
        this.refreshAll()
      } catch (error) {
        if (error.response?.status === 400) {
          ElMessage.warning('菜单可能已存在，已刷新列表')
          this.addMenuDialogVisible = false
          this.refreshAll()
          return
        }
        ElMessage.error(error.message || '新增菜单失败')
      }
    },
    editMenuDialogClosed() {
      this.$refs.editMenuFormRefForm?.resetFields()
      this.menuInfo = createMenuForm()
    },
    async showEditMenuDialog(id) {
      try {
        const { data: res } = await this.$api.menuInfo(id)
        if (res.code !== 200) {
          throw new Error(res.message || '获取菜单详情失败')
        }

        this.menuInfo = {
          ...createMenuForm(),
          ...(res.data || {})
        }
        this.editMenuDialogVisible = true
        this.$nextTick(() => {
          this.$refs.editMenuFormRefForm?.clearValidate()
        })
      } catch (error) {
        ElMessage.error(error.message || '获取菜单详情失败')
      }
    },
    async editMenu() {
      const valid = await this.$refs.editMenuFormRefForm?.validate().catch(() => false)
      if (!valid) {
        return
      }

      try {
        const { data: res } = await this.$api.menuUpdate(this.normalizeMenuPayload(this.menuInfo))
        if (res.code !== 200) {
          throw new Error(res.message || '更新菜单失败')
        }

        ElMessage.success('更新菜单成功')
        this.editMenuDialogVisible = false
        this.refreshAll()
      } catch (error) {
        ElMessage.error(error.message || '更新菜单失败')
      }
    },
    handleCopyMenu(menuData) {
      this.menuForm = {
        ...createMenuForm(),
        parentId: menuData.parentId,
        menuName: menuData.menuName ? `${menuData.menuName}_副本` : '',
        icon: menuData.icon || '',
        value: menuData.value ? `${menuData.value}_copy` : '',
        menuType: Number(menuData.menuType) || 1,
        url: menuData.url ? `${menuData.url}_copy` : '',
        sort: menuData.sort || 0,
        menuStatus: menuData.menuStatus || 2
      }
      this.addMenuDialogVisible = true
      this.$nextTick(() => {
        this.$refs.addMenuFormRefForm?.clearValidate()
      })
    },
    async handleMenuDelete(row) {
      try {
        await ElMessageBox.confirm(`确认删除菜单 ${row.menuName} 吗`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })

        const { data: res } = await this.$api.menuDelete(row.id)
        if (res.code !== 200) {
          throw new Error(res.message || '删除菜单失败')
        }

        ElMessage.success('删除成功')
        this.refreshAll()
      } catch (error) {
        if (!isActionCancelled(error)) {
          ElMessage.error(error.message || '删除失败')
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
.system-menu-page :deep(.page-actions) {
  align-items: center;
}

.system-menu-page :deep(.section-card__body) {
  display: grid;
  gap: 20px;
}

.icon-wrapper {
  width: 36px;
  height: 36px;
  display: grid;
  place-items: center;
  border-radius: 12px;
  border: 1px solid rgba(59, 130, 246, 0.18);
  background: rgba(59, 130, 246, 0.12);
  color: #bfdbfe;
}

.menu-icon {
  font-size: 18px;
}

.icon-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.menu-form :deep(.el-input-number),
.menu-form :deep(.el-select) {
  width: 100%;
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
