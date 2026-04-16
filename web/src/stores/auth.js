import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../api/request'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  function setAuth(t, u) {
    token.value = t
    user.value = u
    localStorage.setItem('token', t)
    localStorage.setItem('user', JSON.stringify(u))
    if (u?.id) {
      localStorage.setItem('userId', u.id)
    }
  }

  function clearAuth() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('userId')
  }

  async function login(username, password) {
    const { data } = await api.post('/auth/login', { username, password })
    setAuth(data.token, data.user)
    return data
  }

  async function register(payload) {
    const { data } = await api.post('/auth/register', payload)
    setAuth(data.token, data.user)
    return data
  }

  async function logout() {
    try {
      await api.post('/auth/logout')
    } catch (e) {}
    clearAuth()
  }

  return {
    token,
    user,
    isAuthenticated,
    isAdmin,
    login,
    register,
    logout,
    setAuth,
    clearAuth
  }
})
