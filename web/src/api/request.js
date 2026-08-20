import axios from 'axios'

const api = axios.create({
  baseURL: `${import.meta.env.BASE_URL}api`
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  res => res,
  err => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      const loginPath = `${import.meta.env.BASE_URL}login`
      if (window.location.pathname !== loginPath) {
        window.location.href = loginPath
      }
    }
    return Promise.reject(err)
  }
)

export default api
