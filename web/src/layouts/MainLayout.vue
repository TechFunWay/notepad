<template>
  <div class="main-layout">
    <el-container class="layout-container">
      <el-header class="header">
        <div class="header-inner">
          <div class="header-left">
            <router-link to="/" class="logo">
              <div class="logo-icon">
                <el-icon><Document /></el-icon>
              </div>
              <span class="logo-text">记事本</span>
            </router-link>
            <span class="version">{{ version }}</span>
          </div>
          <div class="header-right">
            <el-dropdown trigger="click" class="user-dropdown">
              <div class="user-info">
                <div class="user-avatar">
                  {{ user?.username?.charAt(0)?.toUpperCase() || 'U' }}
                </div>
                <span class="user-name">{{ user?.username }}</span>
                <el-icon class="dropdown-icon"><ArrowDown /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <router-link to="/profile">
                    <el-dropdown-item>
                      <el-icon><User /></el-icon>
                      <span>个人中心</span>
                    </el-dropdown-item>
                  </router-link>
                  <template v-if="user?.role === 'admin'">
                    <el-dropdown-item divided>
                      <el-icon><Setting /></el-icon>
                      <span>管理后台</span>
                    </el-dropdown-item>
                    <router-link to="/admin/users">
                      <el-dropdown-item>
                        <el-icon><UserFilled /></el-icon>
                        <span>用户管理</span>
                      </el-dropdown-item>
                    </router-link>
                    <router-link to="/admin/configs">
                      <el-dropdown-item>
                        <el-icon><Tools /></el-icon>
                        <span>系统配置</span>
                      </el-dropdown-item>
                    </router-link>
                  </template>
                  <el-dropdown-item divided @click="handleLogout">
                    <el-icon><SwitchButton /></el-icon>
                    <span>退出登录</span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </el-header>
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, User, SwitchButton, Setting, UserFilled, Tools, ArrowDown } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/api/request'

const router = useRouter()
const auth = useAuthStore()

const user = computed(() => auth.user)
const version = ref('1.0.0')

async function fetchVersion() {
  try {
    const { data } = await api.get('/version')
    version.value = data.version || '1.0.0'
  } catch {
    version.value = '1.0.0'
  }
}

async function handleLogout() {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    auth.logout()
    ElMessage.success('已退出登录')
    router.push('/login')
  } catch {
  }
}

onMounted(() => {
  fetchVersion()
})

onUnmounted(() => {
})
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
  background: #f5f7fa;
}

.layout-container {
  min-height: 100vh;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  box-shadow: 0 2px 12px rgba(102, 126, 234, 0.15);
  padding: 0;
  height: 64px;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-inner {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  text-decoration: none;
  color: white;
  transition: transform 0.2s;
}

.logo:hover {
  transform: scale(1.02);
}

.logo-icon {
  width: 44px;
  height: 44px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-icon .el-icon {
  font-size: 24px;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  letter-spacing: -0.5px;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-dropdown {
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 16px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  transition: all 0.2s;
}

.user-info:hover {
  background: rgba(255, 255, 255, 0.25);
  transform: translateY(-1px);
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: linear-gradient(135deg, #fff 0%, #f0f0f0 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 700;
  color: #667eea;
}

.user-name {
  color: white;
  font-size: 14px;
  font-weight: 500;
}

.dropdown-icon {
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
}

.el-dropdown-menu :deep(.el-dropdown-menu__item a) {
  color: inherit;
  text-decoration: none;
  display: flex;
  align-items: center;
  gap: 8px;
}

.el-dropdown-menu :deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 8px;
}

.main-content {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
  min-height: calc(100vh - 64px);
}

@media (max-width: 768px) {
  .header {
    height: 56px;
  }

  .header-inner {
    padding: 0 16px;
  }

  .logo-icon {
    width: 36px;
    height: 36px;
    border-radius: 10px;
  }

  .logo-icon .el-icon {
    font-size: 20px;
  }

  .logo-text {
    font-size: 17px;
  }

  .user-info {
    padding: 6px 12px;
    border-radius: 10px;
  }

  .user-avatar {
    width: 30px;
    height: 30px;
    border-radius: 8px;
    font-size: 14px;
  }

  .user-name {
    max-width: 60px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 13px;
  }

  .main-content {
    padding: 8px;
  }
}

.version {
  color: rgba(255, 255, 255, 0.7);
  font-size: 12px;
  font-weight: 500;
  padding: 4px 10px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
}

@media (max-width: 768px) {
  .version {
    font-size: 10px;
    padding: 2px 8px;
  }
}
</style>
