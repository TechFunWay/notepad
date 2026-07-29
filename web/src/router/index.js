import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { checkSetupRequired } from '@/utils/setupStatus'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/RegisterView.vue')
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/views/ForgotPasswordView.vue')
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Notes',
        component: () => import('@/views/NotesView.vue')
      },
      {
        path: 'notes-list',
        name: 'NotesList',
        component: () => import('@/views/NotesListView.vue')
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/ProfileView.vue')
      },
      {
        path: 'admin/users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/UsersView.vue'),
        meta: { requiresAdmin: true }
      },
      {
        path: 'admin/configs',
        name: 'AdminConfigs',
        component: () => import('@/views/admin/ConfigView.vue'),
        meta: { requiresAdmin: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const shouldCheckSetup =
    requiresAuth ||
    to.path === '/login' ||
    to.path === '/register'

  if (shouldCheckSetup) {
    try {
      const needsSetup = await checkSetupRequired()
      if (needsSetup) {
        // 新数据库与浏览器中残留的旧令牌不能同时成立，否则会在 / 和
        // /register 之间循环跳转。初始化模式下以服务端状态为准。
        auth.clearAuth()
        if (to.path !== '/register' || to.query.setup !== 'true') {
          return {
            path: '/register',
            query: { setup: 'true' },
            replace: true
          }
        }
        return true
      }
    } catch (e) {
      // API 暂时不可用时交给目标页面展示，不制造重定向循环。
    }
  }

  if (requiresAuth && !auth.isAuthenticated) return '/login'
  if (to.meta.requiresAdmin && !auth.isAdmin) return '/'
  if ((to.path === '/login' || to.path === '/register') && auth.isAuthenticated) return '/'

  return true
})

export default router
