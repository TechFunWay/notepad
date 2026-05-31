import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { getSetupStatus } from '@/api/auth'

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

router.beforeEach(async (to, from, next) => {
  const auth = useAuthStore()

  // 首次访问检查：如果没有用户，跳转到管理员注册
  if (to.path === '/login' || to.path === '/') {
    try {
      const { data } = await getSetupStatus()
      if (data.needs_setup && to.path !== '/register') {
        next('/register?setup=true')
        return
      }
    } catch (e) {
      // API 不可用时静默处理
    }
  }

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    next('/login')
  } else if (to.meta.requiresAdmin && !auth.isAdmin) {
    next('/')
  } else if ((to.path === '/login' || to.path === '/register') && auth.isAuthenticated) {
    next('/')
  } else {
    next()
  }
})

export default router
