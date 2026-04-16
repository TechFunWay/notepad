import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api/request'

export const useConfigStore = defineStore('config', () => {
  const siteTitle = ref('记事本')
  const allowRegister = ref(true)

  async function fetchPublicConfig() {
    try {
      const { data } = await api.get('/public-config')
      siteTitle.value = data.site_title || '记事本'
      allowRegister.value = data.allow_register === 'true'
    } catch (e) {}
  }

  return { siteTitle, allowRegister, fetchPublicConfig }
})
