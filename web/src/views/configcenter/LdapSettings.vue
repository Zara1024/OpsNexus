<template>
  <div class="ldap-settings-page">
    <el-card shadow="hover" class="ldap-card">
      <template #header>
        <div class="card-header">
          <div>
            <div class="title">LDAP 集成</div>
            <div class="subtitle">管理 LDAP 配置、连接测试、默认角色与 LDAP 组角色映射。</div>
          </div>
          <div class="header-actions">
            <el-button icon="Refresh" @click="loadConfig">刷新</el-button>
            <el-button type="primary" :loading="saving" @click="saveConfig">保存配置</el-button>
          </div>
        </div>
      </template>

      <el-alert
        title="保存后，登录接口会在本地密码校验失败时回退到 LDAP；LDAP 新用户首次登录会自动建档，并根据默认角色和组映射补齐角色。"
        type="info"
        :closable="false"
        class="page-alert"
      />

      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px" class="ldap-form">
        <el-row :gutter="16">
          <el-col :xs="24" :md="12">
            <el-form-item label="启用 LDAP">
              <el-switch v-model="form.enable" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="12">
            <el-form-item label="默认角色" prop="defaultRoleId">
              <el-select v-model="form.defaultRoleId" placeholder="请选择默认角色" style="width: 100%" clearable>
                <el-option
                  v-for="role in roleOptions"
                  :key="role.id"
                  :label="`${role.roleName || role.role_name || role.name} (#${role.id})`"
                  :value="role.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :xs="24" :md="8">
            <el-form-item label="LDAP Host" prop="host">
              <el-input v-model="form.host" placeholder="例如 10.0.0.200" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="8">
            <el-form-item label="LDAP Port" prop="port">
              <el-input-number v-model="form.port" :min="1" :max="65535" controls-position="right" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="8">
            <el-form-item label="认证过滤器" prop="authFilter">
              <el-input v-model="form.authFilter" placeholder="例如 (uid=%s)" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="Base DN" prop="baseDn">
          <el-input v-model="form.baseDn" placeholder="例如 ou=users,dc=opsnexus,dc=local" />
        </el-form-item>

        <el-row :gutter="16">
          <el-col :xs="24" :md="12">
            <el-form-item label="Bind 用户">
              <el-input v-model="form.bindUser" placeholder="例如 cn=admin,dc=opsnexus,dc=local" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="12">
            <el-form-item label="Bind 密码">
              <el-input v-model="form.bindPass" show-password placeholder="Bind 用户密码" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :xs="24" :md="4">
            <el-form-item label="昵称属性">
              <el-input v-model="form.attributes.nickname" placeholder="默认 cn" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="4">
            <el-form-item label="邮箱属性">
              <el-input v-model="form.attributes.email" placeholder="默认 mail" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="4">
            <el-form-item label="手机属性">
              <el-input v-model="form.attributes.phone" placeholder="默认 mobile" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="6">
            <el-form-item label="组过滤器">
              <el-input
                v-model="form.groupFilter"
                placeholder="例如 (memberUid=%s) 或 (member={{dn}})"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="6">
            <el-form-item label="组名属性">
              <el-input v-model="form.groupNameAttr" placeholder="默认 cn" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :xs="24" :md="8">
            <el-form-item label="TLS">
              <el-switch v-model="form.tls" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="8">
            <el-form-item label="StartTLS">
              <el-switch v-model="form.startTLS" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="8">
            <el-form-item label="覆盖属性">
              <el-switch v-model="form.coverAttributes" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="记录 LDAP 域、默认角色和使用说明" />
        </el-form-item>

        <el-form-item label="组角色映射">
          <div class="mapping-list">
            <div
              v-for="(mapping, index) in form.roleMappings"
              :key="index"
              class="mapping-row"
            >
              <el-input v-model="mapping.groupName" placeholder="LDAP 组名，例如 ops-admins" />
              <el-select v-model="mapping.roleId" placeholder="请选择角色" clearable>
                <el-option
                  v-for="role in roleOptions"
                  :key="role.id"
                  :label="`${role.roleName || role.role_name || role.name} (#${role.id})`"
                  :value="role.id"
                />
              </el-select>
              <el-button icon="Delete" @click="removeRoleMapping(index)" />
            </div>
            <el-button type="primary" link icon="Plus" @click="addRoleMapping">添加映射</el-button>
          </div>
        </el-form-item>
      </el-form>

      <div class="test-panel">
        <div>
          <div class="test-title">连接测试</div>
          <div class="test-subtitle">使用当前表单配置直接测试 LDAP 连接与 Bind。</div>
        </div>
        <el-button type="success" :loading="testing" @click="testConfig">测试连接</el-button>
      </div>

      <el-alert
        v-if="testResult"
        :title="testResult"
        :type="testSuccess ? 'success' : 'error'"
        :closable="false"
        class="page-alert"
      />
    </el-card>
  </div>
</template>

<script>
import { ElMessage } from 'element-plus'

const createEmptyForm = () => ({
  enable: false,
  host: '',
  port: 389,
  baseDn: '',
  bindUser: '',
  bindPass: '',
  authFilter: '(uid=%s)',
  groupFilter: '(memberUid=%s)',
  groupNameAttr: 'cn',
  coverAttributes: true,
  tls: false,
  startTLS: false,
  defaultRoles: [],
  defaultRoleId: null,
  roleMappings: [],
  attributes: {
    nickname: 'cn',
    phone: 'mobile',
    email: 'mail'
  },
  remark: ''
})

export default {
  name: 'LdapSettings',
  data() {
    return {
      loading: false,
      saving: false,
      testing: false,
      roleOptions: [],
      testResult: '',
      testSuccess: false,
      form: createEmptyForm(),
      rules: {
        host: [{
          validator: (rule, value, callback) => {
            if (this.form.enable && !value) {
              callback(new Error('请输入 LDAP Host'))
              return
            }
            callback()
          },
          trigger: 'blur'
        }],
        port: [{
          validator: (rule, value, callback) => {
            if (this.form.enable && (!value || Number(value) <= 0)) {
              callback(new Error('请输入合法 LDAP 端口'))
              return
            }
            callback()
          },
          trigger: 'change'
        }],
        baseDn: [{
          validator: (rule, value, callback) => {
            if (this.form.enable && !value) {
              callback(new Error('请输入 Base DN'))
              return
            }
            callback()
          },
          trigger: 'blur'
        }],
        authFilter: [{
          validator: (rule, value, callback) => {
            if (this.form.enable && !value) {
              callback(new Error('请输入认证过滤器'))
              return
            }
            callback()
          },
          trigger: 'blur'
        }]
      }
    }
  },
  methods: {
    async loadRoles() {
      const { data: res } = await this.$api.queryRoleList({ pageNum: 1, pageSize: 200 })
      if (res.code !== 200) {
        throw new Error(res.message || '获取角色列表失败')
      }
      const payload = res.data || {}
      this.roleOptions = payload.list || payload || []
    },
    normalizeRoleMappings(list) {
      if (!Array.isArray(list)) return []
      return list
        .map(item => ({
          groupName: item.groupName || '',
          roleId: item.roleId || null
        }))
        .filter(item => item.groupName || item.roleId)
    },
    async loadConfig() {
      this.loading = true
      try {
        await this.loadRoles()
        const { data: res } = await this.$api.getLdapConfig()
        if (res.code !== 200) {
          throw new Error(res.message || '获取 LDAP 配置失败')
        }
        this.form = {
          ...createEmptyForm(),
          ...(res.data || {}),
          roleMappings: this.normalizeRoleMappings((res.data || {}).roleMappings),
          attributes: {
            ...createEmptyForm().attributes,
            ...((res.data || {}).attributes || {})
          }
        }
      } catch (error) {
        ElMessage.error(error.message || '获取 LDAP 配置失败')
      } finally {
        this.loading = false
      }
    },
    buildSubmitPayload() {
      return {
        ...this.form,
        roleMappings: this.normalizeRoleMappings(this.form.roleMappings)
      }
    },
    async saveConfig() {
      const valid = await this.$refs.formRef.validate().catch(() => false)
      if (!valid) return

      this.saving = true
      try {
        const { data: res } = await this.$api.updateLdapConfig(this.buildSubmitPayload())
        if (res.code !== 200) {
          throw new Error(res.message || '保存 LDAP 配置失败')
        }
        ElMessage.success('LDAP 配置保存成功')
        await this.loadConfig()
      } catch (error) {
        ElMessage.error(error?.response?.data?.message || error.message || '保存 LDAP 配置失败')
      } finally {
        this.saving = false
      }
    },
    async testConfig() {
      const valid = await this.$refs.formRef.validate().catch(() => false)
      if (!valid) return

      this.testing = true
      this.testResult = ''
      try {
        const { data: res } = await this.$api.testLdapConfig(this.buildSubmitPayload())
        if (res.code !== 200) {
          throw new Error(res.message || 'LDAP 测试失败')
        }
        this.testSuccess = true
        this.testResult = `连接成功：${res.data?.host || this.form.host}:${res.data?.port || this.form.port}`
      } catch (error) {
        this.testSuccess = false
        this.testResult = error?.response?.data?.message || error.message || 'LDAP 测试失败'
      } finally {
        this.testing = false
      }
    },
    addRoleMapping() {
      this.form.roleMappings.push({
        groupName: '',
        roleId: null
      })
    },
    removeRoleMapping(index) {
      this.form.roleMappings.splice(index, 1)
    }
  },
  mounted() {
    this.loadConfig()
  }
}
</script>

<style scoped>
.ldap-settings-page {
  padding: 20px;
  min-height: 100vh;
  background:
    radial-gradient(circle at top left, rgba(34, 197, 94, 0.12), transparent 28%),
    radial-gradient(circle at top right, rgba(37, 99, 235, 0.14), transparent 22%),
    linear-gradient(145deg, #f7fff9 0%, #f7fbff 46%, #fcfffd 100%);
}

.ldap-card {
  border-radius: 20px;
  border: 1px solid rgba(148, 163, 184, 0.16);
  box-shadow: 0 18px 46px rgba(15, 23, 42, 0.08);
}

.card-header,
.test-panel {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.title {
  font-size: 24px;
  font-weight: 700;
  color: #0f172a;
}

.subtitle,
.test-subtitle {
  margin-top: 6px;
  font-size: 13px;
  color: #64748b;
}

.page-alert {
  margin-bottom: 18px;
}

.ldap-form {
  margin-top: 20px;
}

.ldap-form :deep(.el-input-number) {
  width: 100%;
}

.mapping-list {
  width: 100%;
}

.mapping-row {
  display: grid;
  grid-template-columns: minmax(180px, 1fr) minmax(180px, 1fr) 48px;
  gap: 12px;
  margin-bottom: 10px;
}

.test-panel {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid rgba(148, 163, 184, 0.14);
}

.test-title {
  font-size: 16px;
  font-weight: 700;
  color: #111827;
}

@media (max-width: 768px) {
  .card-header,
  .test-panel {
    flex-direction: column;
    align-items: stretch;
  }

  .mapping-row {
    grid-template-columns: 1fr;
  }
}
</style>
