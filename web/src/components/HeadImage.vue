<template>
  <div class="user-info">
    <div class="user-copy">
      <span class="user-label">Signed in</span>
      <span class="user-username">{{ sysAdmin.username }}</span>
    </div>

    <el-dropdown trigger="click" @command="handleCommand">
      <div class="avatar-wrap">
        <img
            :src="avatarSrc"
            alt="头像"
            class="user-avatar"
            @error="useDefaultAvatar"
        />
      </div>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="profile">个人信息</el-dropdown-item>
          <el-dropdown-item command="logout">退出登录</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script>
import storage from '@/utils/storage';
import { resolveAvatarUrl } from '@/utils/avatarUrl.mjs';

export default {
  name: "HeadImage",
  data() {
    return {
      sysAdmin: storage.getItem("sysAdmin") || {}
    };
  },
  computed: {
    avatarSrc() {
      return resolveAvatarUrl(this.sysAdmin.icon) || require('./../assets/image/touxiang.jpg');
    }
  },
  methods: {
    // 图片加载失败时使用默认头像
    useDefaultAvatar(e) {
      e.target.src = require('./../assets/image/touxiang.jpg');
    },
    // 菜单点击事件
    handleCommand(command) {
          if (command === 'logout') {
              // 调用新的 logout 方法
              this.logout();
            } else if (command === 'profile') {
              this.$router.push('/system/personal'); // 跳转到个人页面路由
            }
          },
          // 退出登录
          async logout() {
            const confirmResult = await this.$confirm('确定要退出登录吗, 是否继续?', '提示', {
              confirmButtonText: '确定',
              cancelButtonText: '取消',
              type: 'warning'
            }).catch(err => err)

            if (confirmResult !== 'confirm') {
              return this.$message.info('已取消退出')
            }
            // 清除本地存储并跳转到登录页
            this.$storage.clearAll()
            this.$router.push('/login')
            this.$message.success('退出成功')
          }
  }
};
</script>

<style lang="less" scoped>
.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.user-copy {
  display: grid;
  gap: 2px;
  text-align: right;
}

.user-label {
  font-size: 11px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--text-subtle, rgba(148, 163, 184, 0.64));
}

.user-username {
  font-size: 13px;
  font-weight: 700;
  color: var(--text-primary, #f8fbff);
  letter-spacing: 0.02em;
}

.avatar-wrap {
  position: relative;
  cursor: pointer;
  padding: 2px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--color-primary, #0ea5e9), var(--color-accent, #6366f1));
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 0 18px rgba(14, 165, 233, 0.3);
  }
}

.user-avatar {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  object-fit: cover;
  display: block;
  border: 2px solid rgba(4, 11, 23, 0.9);
}

@media (max-width: 768px) {
  .user-copy {
    display: none;
  }
}
</style>
