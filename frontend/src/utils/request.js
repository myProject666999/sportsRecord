import axios from 'axios'
import { showToast, showLoadingToast, closeToast } from 'vant'
import { useUserStore } from '@/stores/user'

const request = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

let loadingCount = 0

const showLoading = () => {
  if (loadingCount === 0) {
    showLoadingToast({
      message: '加载中...',
      forbidClick: true,
      duration: 0
    })
  }
  loadingCount++
}

const hideLoading = () => {
  loadingCount--
  if (loadingCount <= 0) {
    loadingCount = 0
    closeToast()
  }
}

request.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    
    if (config.showLoading !== false) {
      showLoading()
    }
    
    return config
  },
  (error) => {
    hideLoading()
    return Promise.reject(error)
  }
)

request.interceptors.response.use(
  (response) => {
    hideLoading()
    const res = response.data
    
    if (res.code === 200) {
      return res
    } else {
      showToast(res.message || '请求失败')
      return Promise.reject(new Error(res.message || '请求失败'))
    }
  },
  (error) => {
    hideLoading()
    
    if (error.response) {
      switch (error.response.status) {
        case 401:
          showToast('登录已过期，请重新登录')
          const userStore = useUserStore()
          userStore.logout()
          setTimeout(() => {
            window.location.href = '/login'
          }, 1000)
          break
        case 403:
          showToast('没有权限访问')
          break
        case 404:
          showToast('请求的资源不存在')
          break
        case 500:
          showToast('服务器错误')
          break
        default:
          showToast(error.message || '网络错误')
      }
    } else {
      showToast('网络连接失败')
    }
    
    return Promise.reject(error)
  }
)

export default request
