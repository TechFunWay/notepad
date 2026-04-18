import api from './request'

export function getSecurityQuestion(username) {
  return api.get('/auth/security-question', { params: { username } })
}

export function forgotPassword(data) {
  return api.post('/auth/forgot-password', data)
}

export function changePassword(data) {
  return api.post('/auth/change-password', data)
}

export function updateSecurityQuestion(data) {
  return api.put('/auth/security-question', data)
}
