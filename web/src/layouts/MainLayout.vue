<template>
  <div class="main-layout">
    <el-container class="layout-container">
      <el-header class="header">
        <div class="header-inner">
          <div class="header-left">
            <router-link to="/notes-list" class="logo">
              <div class="logo-icon">
                <el-icon><Document /></el-icon>
              </div>
              <span class="logo-text">记事本</span>
            </router-link>
            <span class="version">{{ version }}</span>
          </div>
          <div class="header-right">
            <el-tooltip :content="isDark ? '切换亮色模式' : '切换暗色模式'" placement="bottom">
              <button class="theme-toggle" @click="toggleTheme">
                <el-icon><template v-if="isDark"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg></template><template v-else><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg></template></el-icon>
              </button>
            </el-tooltip>
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
                    <router-link to="/admin/users">
                      <el-dropdown-item divided>
                        <el-icon><UserFilled /></el-icon>
                        <span>用户管理</span>
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
import { message } from '@/utils/message'
import { Document, User, SwitchButton, Setting, UserFilled, Tools, ArrowDown } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { useTheme } from '@/composables/useTheme'
const { isDark, toggleTheme } = useTheme()
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
  const confirmed = await message.confirm('确定要退出登录吗？', {
    title: '提示',
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
  if (!confirmed) return
  auth.logout()
  window.location.href = '/login'
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
  background: var(--bg-secondary);
}

.layout-container {
  min-height: 100vh;
}

.header {
  background: var(--header-gradient);
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
  gap: 8px;
}

.theme-toggle {
  width: 40px;
  height: 40px;
  border: none;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.theme-toggle:hover {
  background: rgba(255, 255, 255, 0.25);
  transform: translateY(-1px);
}

.theme-toggle .el-icon {
  font-size: 20px;
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
