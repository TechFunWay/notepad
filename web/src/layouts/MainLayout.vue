<template>
  <div class="app-shell" :class="{ 'immersive-mode': immersiveMode }">
    <a class="skip-link" href="#main-content">跳到主要内容</a>

    <aside class="app-sidebar" aria-label="主导航">
      <router-link to="/notes-list" class="brand">
        <span class="brand-mark" aria-hidden="true">
          <span class="brand-letter">N</span>
        </span>
        <span class="brand-copy">
          <strong>记事本</strong>
          <small>灵感与知识工作台</small>
        </span>
      </router-link>

      <nav class="primary-nav">
        <p class="nav-label">工作空间</p>
        <router-link
          v-for="item in mainNav"
          :key="item.to"
          :to="item.to"
          class="nav-item"
          :class="{ active: isActive(item) }"
        >
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.label }}</span>
        </router-link>

        <template v-if="user?.role === 'admin'">
          <p class="nav-label nav-label-admin">系统管理</p>
          <router-link
            v-for="item in adminNav"
            :key="item.to"
            :to="item.to"
            class="nav-item"
            :class="{ active: isActive(item) }"
          >
            <el-icon><component :is="item.icon" /></el-icon>
            <span>{{ item.label }}</span>
          </router-link>
        </template>
      </nav>

      <div class="sidebar-footer">
        <button
          class="theme-control"
          type="button"
          :aria-label="isDark ? '切换亮色模式' : '切换暗色模式'"
          @click="toggleTheme"
        >
          <el-icon><Sunny v-if="isDark" /><Moon v-else /></el-icon>
          <span>{{ isDark ? '亮色模式' : '暗色模式' }}</span>
        </button>

        <el-dropdown trigger="click" placement="top-start" class="account-dropdown">
          <button class="account-card" type="button" aria-label="打开账号菜单">
            <span class="user-avatar">{{ userInitial }}</span>
            <span class="account-copy">
              <strong>{{ user?.username || '用户' }}</strong>
              <small>{{ user?.role === 'admin' ? '管理员' : '个人账号' }}</small>
            </span>
            <el-icon class="account-more"><MoreFilled /></el-icon>
          </button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="router.push('/profile')">
                <el-icon><User /></el-icon>
                个人设置
              </el-dropdown-item>
              <el-dropdown-item divided @click="handleLogout">
                <el-icon><SwitchButton /></el-icon>
                退出登录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

        <span class="version">Notepad {{ version }}</span>
      </div>
    </aside>

    <header v-if="!immersiveMode" class="mobile-header">
      <router-link to="/notes-list" class="mobile-brand">
        <span class="brand-mark" aria-hidden="true">
          <span class="brand-letter">N</span>
        </span>
        <span>记事本</span>
      </router-link>
      <div class="mobile-actions">
        <button
          class="icon-control"
          type="button"
          :aria-label="isDark ? '切换亮色模式' : '切换暗色模式'"
          @click="toggleTheme"
        >
          <el-icon><Sunny v-if="isDark" /><Moon v-else /></el-icon>
        </button>
        <el-dropdown trigger="click" placement="bottom-end" class="mobile-account-dropdown">
          <button class="mobile-avatar" type="button" aria-label="打开账号菜单">
            {{ userInitial }}
          </button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="router.push('/profile')">
                <el-icon><User /></el-icon>
                个人设置
              </el-dropdown-item>
              <el-dropdown-item divided @click="handleLogout">
                <el-icon><SwitchButton /></el-icon>
                退出登录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </header>

    <main id="main-content" class="app-content" tabindex="-1">
      <!-- 赞赏横幅：仅管理员可见，未支持当前版本前常驻展示（关闭不记忆，
           刷新后再次出现，升级到新版本后也会再次出现） -->
      <div v-if="support.bannerVisible && !support.donateSupported" class="donate-banner" role="note">
        <span class="donate-banner-text">☕ 如果记事本对你有帮助，欢迎请作者喝杯咖啡——1 元也是心意，完全自愿。</span>
        <span class="donate-banner-actions">
          <button type="button" class="donate-banner-link" @click="support.open()">去赞赏</button>
          <button type="button" class="donate-banner-close" title="关闭" aria-label="关闭赞赏横幅" @click="support.dismissBanner()">×</button>
        </span>
      </div>
      <router-view v-slot="{ Component }">
        <transition name="page-fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- 支持按钮：常驻入口，始终可以打开赞赏弹窗 -->
    <button
      v-if="user?.role === 'admin'"
      type="button"
      class="support-fab"
      aria-label="支持作者"
      title="支持作者"
      @click="support.open()"
    >☕</button>

    <SupportModal />

    <nav v-if="!immersiveMode" class="mobile-nav" aria-label="移动端主导航">
      <router-link
        v-for="item in mobileNav"
        :key="item.to"
        :to="item.to"
        class="mobile-nav-item"
        :class="{ active: isActive(item) }"
      >
        <span class="mobile-nav-icon">
          <el-icon><component :is="item.icon" /></el-icon>
        </span>
        <span>{{ item.shortLabel || item.label }}</span>
      </router-link>
    </nav>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Collection,
  EditPen,
  MoreFilled,
  Moon,
  Setting,
  Sunny,
  SwitchButton,
  User,
  UserFilled
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { useTheme } from '@/composables/useTheme'
import { message } from '@/utils/message'
import api from '@/api/request'
import { useSupportStore } from '@/stores/support'
import SupportModal from '@/components/SupportModal.vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const { isDark, toggleTheme } = useTheme()
const support = useSupportStore()

const user = computed(() => auth.user)
const userInitial = computed(() => user.value?.username?.charAt(0)?.toUpperCase() || 'U')
const version = ref('1.0.0')

const mainNav = [
  { to: '/', label: '写作工作台', shortLabel: '写作', icon: EditPen, exact: true },
  { to: '/notes-list', label: '全部笔记', shortLabel: '笔记', icon: Collection },
  { to: '/profile', label: '个人设置', shortLabel: '我的', icon: User }
]

const adminNav = [
  { to: '/admin/users', label: '用户管理', shortLabel: '用户', icon: UserFilled },
  { to: '/admin/configs', label: '系统配置', shortLabel: '设置', icon: Setting }
]

const mobileNav = computed(() => {
  const items = [mainNav[0], mainNav[1]]
  if (user.value?.role === 'admin') {
    items.push({
      to: '/admin/users',
      label: '管理',
      shortLabel: '管理',
      icon: Setting,
      matchPrefix: '/admin'
    })
  }
  items.push(mainNav[2])
  return items
})

function isActive(item) {
  if (item.matchPrefix) return route.path.startsWith(item.matchPrefix)
  return item.exact ? route.path === item.to : route.path.startsWith(item.to)
}

const immersiveMode = computed(() => route.path === '/' && Boolean(route.query.note_id))

async function fetchVersion() {
  try {
    const { data } = await api.get('/version')
    version.value = data.version || '1.0.0'
    // 版本信息顺带交给赞赏提示调度（仅管理员生效）
    await support.init(auth.isAdmin, data)
  } catch {
    version.value = '1.0.0'
  }
}

async function handleLogout() {
  const confirmed = await message.confirm('确定要退出登录吗？', {
    title: '退出登录',
    confirmButtonText: '退出',
    cancelButtonText: '取消',
    type: 'warning'
  })
  if (!confirmed) return
  await auth.logout()
  await router.replace('/login')
}

onMounted(fetchVersion)

watch(
  () => route.fullPath,
  async () => {
    await nextTick()
    document.getElementById('main-content')?.focus({ preventScroll: true })
  }
)
</script>

<style scoped>
.app-shell {
  min-height: 100dvh;
  background: var(--bg-secondary);
}

#main-content:focus {
  outline: none;
}

.skip-link {
  position: fixed;
  top: 12px;
  left: 12px;
  z-index: 300;
  padding: 10px 14px;
  color: #fff;
  background: var(--primary-color);
  border-radius: 10px;
  transform: translateY(-160%);
  transition: transform 180ms ease;
}

.skip-link:focus {
  transform: translateY(0);
}

.app-sidebar {
  position: fixed;
  inset: 0 auto 0 0;
  z-index: 100;
  width: var(--app-sidebar-width);
  display: flex;
  flex-direction: column;
  padding: 22px 16px 16px;
  color: var(--text-primary);
  background: var(--surface-elevated);
  background: color-mix(in srgb, var(--bg-primary) 92%, transparent);
  border-right: 1px solid var(--border-color);
  box-shadow: 10px 0 35px rgba(8, 51, 68, 0.035);
  backdrop-filter: blur(20px);
}

.brand,
.mobile-brand {
  display: flex;
  align-items: center;
  color: var(--text-primary);
  text-decoration: none;
}

.brand {
  gap: 12px;
  padding: 2px 8px 24px;
}

.brand-mark {
  width: 42px;
  height: 42px;
  flex: 0 0 42px;
  display: grid;
  place-items: center;
  color: #fff;
  background: var(--gradient-primary);
  border-radius: 13px;
  box-shadow: 0 9px 20px rgba(8, 145, 178, 0.22);
}

.brand-mark .el-icon {
  font-size: 21px;
}

.brand-copy {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.brand-copy strong {
  font-size: 18px;
  letter-spacing: -0.02em;
}

.brand-copy small,
.account-copy small {
  color: var(--text-muted);
  font-size: 11px;
}

.primary-nav {
  flex: 1;
}

.nav-label {
  margin: 12px 12px 8px;
  color: var(--text-muted);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.12em;
}

.nav-label-admin {
  margin-top: 26px;
}

.nav-item,
.theme-control {
  min-height: 46px;
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 0 13px;
  color: var(--text-secondary);
  background: transparent;
  border: 0;
  border-radius: 12px;
  font: inherit;
  font-size: 14px;
  font-weight: 600;
  text-decoration: none;
  cursor: pointer;
  transition: color 180ms ease, background 180ms ease, transform 180ms ease;
}

.nav-item + .nav-item {
  margin-top: 4px;
}

.nav-item .el-icon,
.theme-control .el-icon {
  flex: 0 0 20px;
  font-size: 19px;
}

.nav-item:hover,
.theme-control:hover {
  color: var(--primary-color);
  background: var(--primary-light);
}

.nav-item.active {
  color: var(--primary-strong);
  background: var(--primary-light);
  box-shadow: inset 3px 0 0 var(--primary-color);
}

.sidebar-footer {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.account-dropdown {
  width: 100%;
}

.account-card {
  width: 100%;
  min-height: 62px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px;
  color: var(--text-primary);
  background: var(--bg-tertiary);
  border: 1px solid var(--border-light);
  border-radius: 14px;
  cursor: pointer;
  text-align: left;
  transition: border-color 180ms ease, background 180ms ease;
}

.account-card:hover {
  background: var(--bg-primary);
  border-color: var(--border-strong);
}

.user-avatar,
.mobile-avatar {
  display: grid;
  place-items: center;
  color: var(--primary-strong);
  background: var(--primary-light);
  border: 1px solid var(--border-color);
  border-color: color-mix(in srgb, var(--primary-color) 24%, transparent);
  font-size: 14px;
  font-weight: 800;
}

/* 手机端头像是下拉菜单的触发按钮，重置 button 默认样式以对齐原链接外观 */
.mobile-avatar {
  padding: 0;
  font-family: inherit;
  cursor: pointer;
}

.user-avatar {
  width: 38px;
  height: 38px;
  flex: 0 0 38px;
  border-radius: 11px;
}

.account-copy {
  min-width: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.account-copy strong {
  overflow: hidden;
  font-size: 13px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.account-more {
  color: var(--text-muted);
}

.version {
  padding: 2px 10px;
  color: var(--text-muted);
  font-size: 10px;
  text-align: center;
}

.app-content {
  min-width: 0;
  min-height: 100dvh;
  margin-left: var(--app-sidebar-width);
  padding: 20px;
}

/* 赞赏横幅：未支持当前版本前常驻展示，关闭不记忆 */
.donate-banner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
  padding: 10px 14px;
  color: #fff;
  background: linear-gradient(90deg, #f59e0b, #f97316);
  border-radius: 12px;
  font-size: 13px;
  box-shadow: 0 6px 18px rgba(249, 115, 22, 0.18);
}

.donate-banner-text {
  min-width: 0;
}

.donate-banner-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.donate-banner-link {
  padding: 0;
  color: #fff;
  background: none;
  border: 0;
  font: inherit;
  font-weight: 700;
  text-decoration: underline;
  text-underline-offset: 2px;
  cursor: pointer;
}

.donate-banner-link:hover {
  opacity: 0.85;
}

.donate-banner-close {
  padding: 0 2px;
  color: #fff;
  background: none;
  border: 0;
  font-size: 16px;
  line-height: 1;
  cursor: pointer;
}

.donate-banner-close:hover {
  opacity: 0.85;
}

/* 支持按钮：常驻入口，始终可以打开赞赏弹窗 */
.support-fab {
  position: fixed;
  right: 22px;
  bottom: 24px;
  z-index: 90;
  width: 46px;
  height: 46px;
  display: grid;
  place-items: center;
  color: #fff;
  background: linear-gradient(135deg, #f59e0b, #f97316);
  border: 0;
  border-radius: 999px;
  font-size: 20px;
  cursor: pointer;
  box-shadow: 0 8px 22px rgba(249, 115, 22, 0.3);
  transition: transform 160ms ease;
}

.support-fab:hover {
  transform: scale(1.06);
}

.mobile-header,
.mobile-nav {
  display: none;
}

.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 160ms ease, transform 160ms ease;
}

.page-fade-enter-from {
  opacity: 0;
  transform: translateY(4px);
}

.page-fade-leave-to {
  opacity: 0;
}

@media (max-width: 1120px) and (min-width: 769px) {
  .app-sidebar {
    align-items: center;
    padding-inline: 12px;
  }

  .brand {
    padding-inline: 0;
  }

  .brand-copy,
  .nav-label,
  .nav-item span,
  .theme-control span,
  .account-copy,
  .account-more,
  .version {
    display: none;
  }

  .nav-item,
  .theme-control {
    justify-content: center;
    padding: 0;
  }

  .account-card {
    min-height: 54px;
    justify-content: center;
    padding: 7px;
  }
}

@media (max-width: 768px) {
  .app-shell {
    padding: calc(60px + env(safe-area-inset-top, 0px)) 0 calc(70px + env(safe-area-inset-bottom, 0px));
  }

  .app-sidebar {
    display: none;
  }

  .mobile-header {
    position: fixed;
    inset: 0 0 auto;
    z-index: 100;
    height: calc(60px + env(safe-area-inset-top, 0px));
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: env(safe-area-inset-top, 0px) 16px 0;
    background: var(--surface-elevated);
    background: color-mix(in srgb, var(--bg-primary) 90%, transparent);
    border-bottom: 1px solid var(--border-color);
    backdrop-filter: blur(20px);
  }

  .mobile-brand {
    gap: 10px;
    font-size: 17px;
    font-weight: 800;
  }

  .mobile-brand .brand-mark {
    width: 34px;
    height: 34px;
    flex-basis: 34px;
    border-radius: 10px;
  }

  .mobile-brand .brand-mark .el-icon {
    font-size: 18px;
  }

  .mobile-actions {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .icon-control,
  .mobile-avatar {
    width: 42px;
    height: 42px;
    border-radius: 12px;
  }

  .icon-control {
    display: grid;
    place-items: center;
    color: var(--text-secondary);
    background: transparent;
    border: 0;
    cursor: pointer;
  }

  .mobile-avatar {
    text-decoration: none;
  }

  .app-content {
    min-height: calc(100dvh - 130px);
    margin-left: 0;
    padding: 10px;
  }

  .support-fab {
    right: 14px;
    bottom: calc(84px + env(safe-area-inset-bottom, 0px));
    width: 42px;
    height: 42px;
    font-size: 18px;
  }

  .mobile-nav {
    position: fixed;
    inset: auto 0 0;
    z-index: 100;
    height: calc(70px + env(safe-area-inset-bottom, 0px));
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(0, 1fr));
    padding: 6px 10px env(safe-area-inset-bottom, 0px);
    background: var(--surface-elevated);
    background: color-mix(in srgb, var(--bg-primary) 94%, transparent);
    border-top: 1px solid var(--border-color);
    box-shadow: 0 -10px 30px rgba(8, 51, 68, 0.06);
    backdrop-filter: blur(20px);
  }

  .mobile-nav-item {
    min-width: 0;
    min-height: 56px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 3px;
    color: var(--text-muted);
    border-radius: 12px;
    font-size: 11px;
    font-weight: 700;
    text-decoration: none;
  }

  .mobile-nav-item .el-icon {
    font-size: 20px;
  }

  .mobile-nav-item.active {
    color: var(--primary-strong);
    background: var(--primary-light);
  }
}
</style>
