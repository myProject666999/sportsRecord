import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userLogin, userRegister, getUserInfo, adminLogin, getAdminInfo } from '@/api'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))
  const isAdmin = ref(localStorage.getItem('isAdmin') === 'true')

  const isLoggedIn = computed(() => !!token.value)

  async function login(username, password) {
    const res = await userLogin({ username, password })
    token.value = res.data.token
    userInfo.value = res.data.user
    isAdmin.value = false
    localStorage.setItem('token', res.data.token)
    localStorage.setItem('userInfo', JSON.stringify(res.data.user))
    localStorage.setItem('isAdmin', 'false')
    return res
  }

  async function adminLoginAction(username, password) {
    const res = await adminLogin({ username, password })
    token.value = res.data.token
    userInfo.value = res.data.admin
    isAdmin.value = true
    localStorage.setItem('token', res.data.token)
    localStorage.setItem('userInfo', JSON.stringify(res.data.admin))
    localStorage.setItem('isAdmin', 'true')
    return res
  }

  async function register(username, password, nickname) {
    const res = await userRegister({ username, password, nickname })
    token.value = res.data.token
    userInfo.value = res.data.user
    isAdmin.value = false
    localStorage.setItem('token', res.data.token)
    localStorage.setItem('userInfo', JSON.stringify(res.data.user))
    localStorage.setItem('isAdmin', 'false')
    return res
  }

  async function fetchUserInfo() {
    if (isAdmin.value) {
      const res = await getAdminInfo()
      userInfo.value = res.data
      localStorage.setItem('userInfo', JSON.stringify(res.data))
      return res
    } else {
      const res = await getUserInfo()
      userInfo.value = res.data
      localStorage.setItem('userInfo', JSON.stringify(res.data))
      return res
    }
  }

  function logout() {
    token.value = ''
    userInfo.value = {}
    isAdmin.value = false
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    localStorage.removeItem('isAdmin')
  }

  function updateUserInfo(info) {
    userInfo.value = { ...userInfo.value, ...info }
    localStorage.setItem('userInfo', JSON.stringify(userInfo.value))
  }

  return {
    token,
    userInfo,
    isAdmin,
    isLoggedIn,
    login,
    adminLoginAction,
    register,
    fetchUserInfo,
    logout,
    updateUserInfo
  }
})
