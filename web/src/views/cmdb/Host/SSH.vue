<template>
  <div class="ssh-page">
    <div class="ssh-layout">
      <section class="surface-card ssh-sidebar">
        <header class="surface-card__header">
          <div>
            <div class="surface-card__eyebrow">Asset Groups</div>
            <h3 class="surface-card__title">资产分组</h3>
            <p class="surface-card__subtitle">选择主机后双击即可进入 SSH 终端。</p>
          </div>
          <div class="surface-card__badge">
            <strong>{{ groupList.length }}</strong>
            <span>Groups</span>
          </div>
        </header>

        <div class="ssh-sidebar__search">
          <el-input
            v-model="groupSearchText"
            placeholder="搜索分组或主机"
            clearable
            class="ssh-search-input"
            @input="handleGroupSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>

        <div class="ssh-sidebar__tree">
          <el-tree
            ref="groupTree"
            :data="groupList"
            :props="defaultProps"
            node-key="id"
            :expanded-keys="expandedKeys"
            :highlight-current="true"
            class="ssh-tree"
            @node-click="handleGroupClick"
            @node-expand="handleNodeExpand"
            @node-collapse="handleNodeCollapse"
          >
            <template v-slot="{ node, data }">
              <div
                :class="[
                  'ssh-tree-node',
                  {
                    'ssh-tree-node--root': !data.parentId,
                    'ssh-tree-node--host': data.isHost
                  }
                ]"
                @dblclick.stop="handleTreeNodeDblClick(data)"
              >
                <div class="ssh-tree-node__main">
                  <span class="ssh-tree-node__icon">
                    <template v-if="data.isHost">
                      <el-icon><Platform /></el-icon>
                    </template>
                    <template v-else-if="data.parentId">
                      <el-icon><FolderRemove /></el-icon>
                    </template>
                    <template v-else>
                      <el-icon v-if="expandedKeys.includes(node.key)"><FolderOpened /></el-icon>
                      <el-icon v-else><Folder /></el-icon>
                    </template>
                  </span>
                  <span class="ssh-tree-node__label">{{ node.label }}</span>
                </div>
                <span v-if="data.isHost" class="ssh-tree-node__tag">Host</span>
              </div>
            </template>
          </el-tree>
        </div>
      </section>

      <section class="surface-card terminal-preview">
        <header class="surface-card__header surface-card__header--wide">
          <div>
            <div class="surface-card__eyebrow">Terminal Access</div>
            <h3 class="surface-card__title">终端预览</h3>
            <p class="surface-card__subtitle">
              {{ currentHost ? '当前已选主机，点击下方按钮或双击左侧主机直接进入终端。' : '先从左侧资产分组中选择一台主机，再进入终端。' }}
            </p>
          </div>
          <div class="terminal-preview__summary">
            <div class="terminal-preview__summary-label">当前主机</div>
            <div class="terminal-preview__summary-value">
              {{ currentHost ? currentHost.hostName : '未选择' }}
            </div>
          </div>
        </header>

        <div class="terminal-preview__body">
          <div class="terminal-preview__screen">
            <div class="terminal-preview__screen-bar">
              <span class="terminal-preview__dot terminal-preview__dot--danger"></span>
              <span class="terminal-preview__dot terminal-preview__dot--warning"></span>
              <span class="terminal-preview__dot terminal-preview__dot--success"></span>
              <span class="terminal-preview__screen-title">OpsNexus SSH Preview</span>
            </div>

            <div class="terminal-preview__screen-content">
              <template v-if="currentHost">
                <div class="terminal-preview__welcome">已准备连接目标主机</div>
                <div class="terminal-preview__host">{{ currentHost.hostName }}</div>
                <div class="terminal-preview__address">
                  {{ currentHost.sshName }}@{{ currentHost.sshIp }}:{{ currentHost.sshPort }}
                </div>
                <div class="terminal-preview__prompt">
                  <span class="terminal-preview__prompt-user">{{ currentHost.sshName }}</span>
                  <span>@</span>
                  <span class="terminal-preview__prompt-host">{{ currentHost.hostName }}</span>
                  <span>:~$</span>
                  <span class="terminal-preview__cursor"></span>
                </div>
                <div class="terminal-preview__hint">双击左侧主机或点击“立即进入终端”开始连接。</div>
              </template>
              <template v-else>
                <div class="terminal-preview__welcome">欢迎访问 OpsNexus 业务终端系统</div>
                <div class="terminal-preview__address">从左侧分组树中选择主机后，即可发起 SSH 连接。</div>
                <div class="terminal-preview__prompt">
                  <span class="terminal-preview__prompt-user">ops</span>
                  <span>@</span>
                  <span class="terminal-preview__prompt-host">terminal</span>
                  <span>:~$</span>
                  <span class="terminal-preview__cursor"></span>
                </div>
                <div class="terminal-preview__hint">支持点击选择、搜索定位和双击直达终端。</div>
              </template>
            </div>
          </div>

          <div class="terminal-preview__actions">
            <el-button
              type="primary"
              :disabled="!currentHost"
              @click="openTerminal"
            >
              立即进入终端
            </el-button>
          </div>
        </div>
      </section>
    </div>

    <terminal ref="terminal" :current-host="currentHost" />
  </div>
</template>

<script>
import Terminal from './Terminal.vue'

export default {
  components: {
    Terminal
  },
  data() {
    return {
      groupSearchText: '',
      expandedKeys: [],
      groupList: [],
      defaultProps: {
        children: 'children',
        label: 'name'
      },
      currentHost: null,
      showTerminalDrawer: false
    }
  },
  created() {
    this.loadGroupList()
  },
  methods: {
    async loadGroupList() {
      try {
        const response = await this.$api.getGroupListWithHosts()
        if (response.data.code === 200) {
          this.groupList = response.data.data.map(group => ({
            ...group,
            children: group.children ? group.children.map(child => ({
              ...child,
              children: child.hosts ? child.hosts.map(host => ({
                id: host.id,
                name: host.hostName,
                isHost: true,
                hostData: host
              })) : []
            })) : []
          }))
          
          if (this.groupList.length > 0) {
            this.expandedKeys = [this.groupList[0].id]
          }
        }
      } catch (error) {
        console.error('加载分组列表失败:', error)
        this.$message.error('加载分组列表失败')
      }
    },
    
    handleGroupClick(node, element) {
      if (element.data.isHost) {
        this.currentHost = {
          id: element.data.hostData.id,
          hostName: element.data.hostData.hostName,
          sshName: element.data.hostData.sshName,
          sshIp: element.data.hostData.sshIp,
          sshPort: element.data.hostData.sshPort
        }
      }
    },
    
    handleGroupDblClick(node, element) {
      if (element.data.isHost) {
        this.currentHost = {
          id: element.data.hostData.id,
          hostName: element.data.hostData.hostName,
          sshName: element.data.hostData.sshName,
          sshIp: element.data.hostData.sshIp,
          sshPort: element.data.hostData.sshPort
        }

        this.openTerminal()
      }
    },

    handleTreeNodeDblClick(data) {
      if (!data?.isHost) {
        return
      }

      this.currentHost = {
        id: data.hostData.id,
        hostName: data.hostData.hostName,
        sshName: data.hostData.sshName,
        sshIp: data.hostData.sshIp,
        sshPort: data.hostData.sshPort
      }

      this.openTerminal()
    },

    openTerminal() {
      if (!this.currentHost || !this.$refs.terminal) {
        return
      }

      this.$refs.terminal.show()
    },

    handleGroupSearch() {
      if (!this.groupSearchText) {
        this.expandedKeys = []
        return
      }
      
      try {
        const findMatchingGroup = (groups, searchText) => {
          for (const group of groups) {
            if (group.name.includes(searchText)) {
              return group
            }
            if (group.children && group.children.length > 0) {
              const found = findMatchingGroup(group.children, searchText)
              if (found) return found
            }
          }
          return null
        }
        
        const matchingGroup = findMatchingGroup(this.groupList, this.groupSearchText)
        if (matchingGroup) {
          const findPath = (groups, targetId, path = []) => {
            for (const group of groups) {
              if (group.id === targetId) {
                return [...path, group.id]
              }
              if (group.children && group.children.length > 0) {
                const foundPath = findPath(group.children, targetId, [...path, group.id])
                if (foundPath) return foundPath
              }
            }
            return null
          }
          
          const expandPath = findPath(this.groupList, matchingGroup.id)
          if (expandPath) {
            this.expandedKeys = expandPath.slice(0, -1)
            this.$nextTick(() => {
              this.$refs.groupTree.setCurrentKey(matchingGroup.id)
            })
          }
        } else {
          this.$message.warning('未找到匹配的分组或主机')
        }
      } catch (error) {
        console.error('搜索分组失败:', error)
        this.$message.error('搜索分组失败')
      }
    },
    
    handleNodeExpand(data, node) {
      if (!this.expandedKeys.includes(node.key)) {
        this.expandedKeys.push(node.key)
      }
    },
    
    handleNodeCollapse(data, node) {
      this.expandedKeys = this.expandedKeys.filter(key => key !== node.key)
    }
  }
}
</script>

<style scoped>
.ssh-page {
  min-height: calc(100vh - 180px);
  padding: 0;
  background: transparent;
}

.ssh-layout {
  display: grid;
  grid-template-columns: 300px minmax(0, 1fr);
  gap: 18px;
  min-height: calc(100vh - 220px);
}

.surface-card {
  display: flex;
  flex-direction: column;
  min-height: 0;
  background: rgba(12, 24, 41, 0.92);
  backdrop-filter: blur(18px);
  border-radius: 24px;
  border: 1px solid var(--border-subtle);
  box-shadow: var(--shadow-card);
  overflow: hidden;
}

.surface-card__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  padding: 22px 22px 18px;
  border-bottom: 1px solid rgba(148, 163, 184, 0.12);
}

.surface-card__header--wide {
  align-items: center;
}

.surface-card__eyebrow {
  margin-bottom: 8px;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: #60a5fa;
}

.surface-card__title {
  margin: 0;
  font-size: 28px;
  font-weight: 800;
  color: #f8fafc;
}

.surface-card__subtitle {
  margin: 10px 0 0;
  max-width: 620px;
  color: #94a3b8;
  font-size: 14px;
  line-height: 1.7;
}

.surface-card__badge {
  display: grid;
  justify-items: center;
  gap: 2px;
  min-width: 86px;
  padding: 12px 14px;
  border-radius: 16px;
  border: 1px solid rgba(96, 165, 250, 0.16);
  background: rgba(15, 23, 42, 0.74);
}

.surface-card__badge strong {
  color: #f8fafc;
  font-size: 24px;
  line-height: 1;
}

.surface-card__badge span {
  color: #94a3b8;
  font-size: 11px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.ssh-sidebar__search {
  padding: 18px 22px 14px;
}

.ssh-search-input :deep(.el-input__wrapper) {
  min-height: 46px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-subtle);
  border-radius: 14px;
  box-shadow: none;
}

.ssh-search-input :deep(.el-input__wrapper:hover) {
  border-color: rgba(96, 165, 250, 0.26);
}

.ssh-search-input :deep(.el-input__wrapper.is-focus) {
  border-color: var(--border-strong);
  box-shadow: 0 0 0 3px rgba(53, 164, 255, 0.12);
}

.ssh-search-input :deep(.el-input__inner) {
  color: #e2e8f0;
}

.ssh-search-input :deep(.el-input__prefix),
.ssh-search-input :deep(.el-input__suffix) {
  color: #94a3b8;
}

.ssh-sidebar__tree {
  flex: 1;
  min-height: 0;
  padding: 0 14px 18px;
  overflow: auto;
}

.ssh-tree {
  background: transparent;
}

.ssh-tree :deep(.el-tree) {
  background: transparent;
  color: inherit;
}

.ssh-tree :deep(.el-tree-node__content) {
  min-height: 46px;
  padding: 0;
  border-radius: 14px;
  background: transparent;
}

.ssh-tree :deep(.el-tree-node__content:hover),
.ssh-tree :deep(.el-tree-node:focus > .el-tree-node__content),
.ssh-tree :deep(.el-tree-node.is-current > .el-tree-node__content) {
  background: rgba(59, 130, 246, 0.12);
}

.ssh-tree :deep(.el-tree-node__expand-icon) {
  color: #94a3b8;
}

.ssh-tree-node {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 10px 14px;
}

.ssh-tree-node__main {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 10px;
}

.ssh-tree-node__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border-radius: 10px;
  color: #93c5fd;
  background: rgba(59, 130, 246, 0.12);
  flex-shrink: 0;
}

.ssh-tree-node__label {
  min-width: 0;
  color: #dbe7ff;
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.ssh-tree-node--root .ssh-tree-node__label {
  color: #f8fafc;
  font-weight: 700;
}

.ssh-tree-node--host .ssh-tree-node__icon {
  color: #67e8f9;
  background: rgba(34, 211, 238, 0.14);
}

.ssh-tree-node__tag {
  flex-shrink: 0;
  padding: 4px 8px;
  border-radius: 999px;
  background: rgba(34, 211, 238, 0.12);
  color: #67e8f9;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.04em;
  text-transform: uppercase;
}

.terminal-preview__summary {
  display: grid;
  gap: 4px;
  min-width: 180px;
  padding: 14px 16px;
  border-radius: 18px;
  border: 1px solid rgba(148, 163, 184, 0.14);
  background: rgba(15, 23, 42, 0.72);
}

.terminal-preview__summary-label {
  color: #94a3b8;
  font-size: 12px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.terminal-preview__summary-value {
  color: #f8fafc;
  font-size: 16px;
  font-weight: 700;
}

.terminal-preview__body {
  flex: 1;
  min-height: 0;
  padding: 22px;
  display: grid;
  gap: 18px;
}

.terminal-preview__screen {
  flex: 1;
  min-height: 460px;
  display: flex;
  flex-direction: column;
  border-radius: 22px;
  overflow: hidden;
  border: 1px solid rgba(30, 41, 59, 0.92);
  background: linear-gradient(180deg, rgba(8, 15, 27, 0.98), rgba(11, 18, 32, 0.98));
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.02);
}

.terminal-preview__screen-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 18px;
  border-bottom: 1px solid rgba(148, 163, 184, 0.12);
  background: rgba(255, 255, 255, 0.02);
}

.terminal-preview__dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.terminal-preview__dot--danger {
  background: #fb7185;
}

.terminal-preview__dot--warning {
  background: #fbbf24;
}

.terminal-preview__dot--success {
  background: #4ade80;
}

.terminal-preview__screen-title {
  margin-left: 10px;
  color: #94a3b8;
  font-size: 12px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.terminal-preview__screen-content {
  flex: 1;
  display: grid;
  align-content: center;
  justify-items: center;
  gap: 14px;
  padding: 36px;
  text-align: center;
  background:
    radial-gradient(circle at top, rgba(37, 99, 235, 0.14), transparent 38%),
    linear-gradient(180deg, rgba(2, 6, 23, 0.82), rgba(15, 23, 42, 0.96));
  font-family: Consolas, Monaco, 'Courier New', monospace;
}

.terminal-preview__welcome {
  color: #e2e8f0;
  font-size: 30px;
  font-weight: 700;
  line-height: 1.4;
}

.terminal-preview__host {
  color: #67e8f9;
  font-size: 24px;
  font-weight: 700;
}

.terminal-preview__address {
  color: #94a3b8;
  font-size: 15px;
  line-height: 1.8;
}

.terminal-preview__prompt {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 14px 18px;
  border-radius: 14px;
  border: 1px solid rgba(34, 211, 238, 0.16);
  background: rgba(15, 23, 42, 0.72);
  color: #dbe7ff;
  font-size: 18px;
}

.terminal-preview__prompt-user {
  color: #4ade80;
}

.terminal-preview__prompt-host {
  color: #67e8f9;
}

.terminal-preview__cursor {
  width: 10px;
  height: 20px;
  margin-left: 4px;
  border-radius: 2px;
  background: #67e8f9;
  animation: cursor-blink 1s steps(2, start) infinite;
}

.terminal-preview__hint {
  max-width: 540px;
  color: #94a3b8;
  font-size: 14px;
  line-height: 1.8;
}

.terminal-preview__actions {
  display: flex;
  justify-content: flex-end;
}

@keyframes cursor-blink {
  0%, 49% {
    opacity: 1;
  }

  50%, 100% {
    opacity: 0;
  }
}

@media (max-width: 1100px) {
  .ssh-layout {
    grid-template-columns: 1fr;
  }

  .surface-card__header--wide {
    flex-direction: column;
    align-items: flex-start;
  }

  .terminal-preview__summary {
    min-width: 0;
    width: 100%;
  }
}

@media (max-width: 768px) {
  .surface-card__header,
  .terminal-preview__body {
    padding: 18px;
  }

  .surface-card__title {
    font-size: 24px;
  }

  .terminal-preview__screen {
    min-height: 400px;
  }

  .terminal-preview__screen-content {
    padding: 28px 20px;
  }

  .terminal-preview__welcome {
    font-size: 22px;
  }

  .terminal-preview__host {
    font-size: 18px;
  }

  .terminal-preview__prompt {
    font-size: 15px;
    flex-wrap: wrap;
    justify-content: center;
  }
}
</style>
