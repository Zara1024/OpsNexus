<template>
  <el-dialog
    title="新建设备"
    v-model="dialogVisible"
    width="46%"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <div class="asset-form-section">
        <div class="asset-form-section__title">基础信息</div>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="名称" prop="name">
              <el-input v-model="form.name" placeholder="请输入设备名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="平台" prop="platform">
              <el-radio-group v-model="form.platform" class="device-platform-group">
                <el-radio
                  v-for="option in platformOptions"
                  :key="option.value"
                  :label="option.value"
                >
                  {{ option.label }}
                </el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="地址" prop="address">
              <el-input v-model="form.address" placeholder="请输入管理地址或 IP" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="分组/节点" prop="groupId">
              <el-select
                v-model="form.groupId"
                filterable
                clearable
                style="width: 100%"
                placeholder="请选择分组/节点"
              >
                <el-option
                  v-for="group in selectableGroups"
                  :key="group.id"
                  :label="group.displayName"
                  :value="group.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </div>

      <div class="asset-form-section">
        <div class="asset-form-section__title">访问设置</div>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="协议组" prop="protocolGroup">
              <el-checkbox-group v-model="protocolGroupSelection">
                <el-checkbox label="ssh">SSH</el-checkbox>
                <el-checkbox label="telnet">Telnet</el-checkbox>
                <el-checkbox label="web">Web</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="账号" prop="accountId">
              <el-select
                v-model="form.accountId"
                filterable
                clearable
                style="width: 100%"
                placeholder="请选择账号"
              >
                <el-option
                  v-for="account in accountList"
                  :key="account.id"
                  :label="account.alias || account.name || String(account.id)"
                  :value="account.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="设备类型" prop="deviceType">
              <el-select v-model="form.deviceType" style="width: 100%">
                <el-option label="通用网络设备" value="network" />
                <el-option label="交换机" value="switch" />
                <el-option label="路由器" value="router" />
                <el-option label="防火墙" value="firewall" />
                <el-option label="负载均衡" value="loadbalancer" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="激活中" prop="isActive">
              <el-switch v-model="form.isActive" inline-prompt active-text="是" inactive-text="否" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="SSH 端口" prop="sshPort">
              <el-input-number v-model="form.sshPort" :min="1" :max="65535" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Telnet 端口" prop="telnetPort">
              <el-input-number v-model="form.telnetPort" :min="1" :max="65535" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="24">
            <el-form-item label="Web 地址" prop="webUrl">
              <el-input v-model="form.webUrl" placeholder="例如 https://10.0.0.1 或 10.0.0.1" />
            </el-form-item>
          </el-col>
        </el-row>
      </div>

      <div class="asset-form-section">
        <div class="asset-form-section__title">补充信息</div>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="标签" prop="tags">
              <el-input v-model="form.tags" placeholder="多个标签可用逗号分隔" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="备注" prop="remark">
              <el-input v-model="form.remark" placeholder="补充说明" />
            </el-form-item>
          </el-col>
        </el-row>
      </div>
    </el-form>

    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="submitForm">确定</el-button>
    </template>
  </el-dialog>
</template>

<script>
import {
  createDeviceAssetFormModel,
  getCmdbDevicePlatformOptions
} from '@/utils/cmdbAssetPresentation.mjs'

function flattenSelectableGroups(groupList = []) {
  const groups = []

  const walk = (items) => {
    items.forEach((group) => {
      if (group.isDefault || !group.parentId) {
        groups.push({
          id: group.id,
          displayName: group.name
        })
      }
      if (Array.isArray(group.children) && group.children.length) {
        group.children.forEach((child) => {
          if (!child.children || !child.children.length) {
            groups.push({
              id: child.id,
              displayName: child.name
            })
          }
        })
      }
    })
  }

  walk(groupList)
  return groups
}

export default {
  name: 'CreateDevice',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    groupList: {
      type: Array,
      required: true
    },
    accountList: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      dialogVisible: this.visible,
      form: createDeviceAssetFormModel(),
      protocolGroupSelection: ['ssh'],
      rules: {
        name: [{ required: true, message: '请输入设备名称', trigger: 'blur' }],
        platform: [{ required: true, message: '请选择平台', trigger: 'change' }],
        address: [{ required: true, message: '请输入设备地址', trigger: 'blur' }],
        groupId: [{ required: true, message: '请选择分组/节点', trigger: 'change' }],
        accountId: [{ required: true, message: '请选择账号', trigger: 'change' }],
        deviceType: [{ required: true, message: '请选择设备类型', trigger: 'change' }],
        sshPort: [{ type: 'number', required: true, message: '请输入 SSH 端口', trigger: 'change' }],
        telnetPort: [{ type: 'number', required: true, message: '请输入 Telnet 端口', trigger: 'change' }]
      }
    }
  },
  computed: {
    selectableGroups() {
      return flattenSelectableGroups(this.groupList)
    },
    platformOptions() {
      return getCmdbDevicePlatformOptions(this.form.platform)
    }
  },
  watch: {
    visible(newVal) {
      this.dialogVisible = newVal
      if (newVal) {
        this.resetForm()
      }
    }
  },
  methods: {
    resetForm() {
      this.form = createDeviceAssetFormModel()
      this.protocolGroupSelection = ['ssh']
      this.$nextTick(() => {
        this.$refs.formRef?.clearValidate()
      })
    },
    handleClose() {
      this.resetForm()
      this.$emit('close')
    },
    async submitForm() {
      try {
        await this.$refs.formRef.validate()

        const protocolGroup = this.protocolGroupSelection.length
          ? this.protocolGroupSelection.join(',')
          : 'ssh'

        this.$emit('submit', {
          ...this.form,
          groupId: Number(this.form.groupId),
          accountId: Number(this.form.accountId),
          protocolGroup,
          sshPort: Number(this.form.sshPort || 22),
          telnetPort: Number(this.form.telnetPort || 23),
          isActive: Boolean(this.form.isActive)
        })
      } catch (error) {
        console.error('设备创建表单校验失败:', error)
      }
    }
  }
}
</script>

<style scoped>
.asset-form-section {
  margin-bottom: 14px;
  padding: 12px 12px 4px;
  border-radius: 12px;
  border: 1px solid rgba(148, 163, 184, 0.18);
  background: rgba(15, 23, 42, 0.28);
}

.asset-form-section__title {
  margin-bottom: 10px;
  color: #e2e8f0;
  font-size: 13px;
  font-weight: 600;
}

.device-platform-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 16px;
}

.device-platform-group :deep(.el-radio) {
  margin-right: 0;
}
</style>
