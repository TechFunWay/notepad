import api from './request'

export function getPublicConfig() {
  return api.get('/public-config')
}

export function getConfigs() {
  return api.get('/configs')
}

export function updateConfig(key, value) {
  return api.put(`/configs/${key}`, { value })
}
