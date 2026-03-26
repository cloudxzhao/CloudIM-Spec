import axios from 'axios'
import { useUserStore } from '@/stores/user'

// 创建 axios 实例
const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 10000
})

// 请求拦截器 - 添加 Token
api.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    if (error.response?.status === 401) {
      // Token 过期或无效，跳转到登录页
      const userStore = useUserStore()
      userStore.logout()
      window.location.hash = '#/login'
    }
    return Promise.reject(error)
  }
)

export default api
