<template>
  <div class="admin-login-page">
    <div class="login-container">
      <div class="login-header">
        <div class="logo">
          <van-icon name="manager-o" size="48" />
        </div>
        <h1 class="title">运动健康管理</h1>
        <p class="subtitle">管理员登录</p>
      </div>
      
      <van-form @submit="onSubmit">
        <van-field
          v-model="form.username"
          name="username"
          label="用户名"
          placeholder="请输入用户名"
          :rules="[{ required: true, message: '请输入用户名' }]"
        >
          <template #left-icon>
            <van-icon name="user-o" />
          </template>
        </van-field>
        
        <van-field
          v-model="form.password"
          name="password"
          label="密码"
          type="password"
          placeholder="请输入密码"
          :rules="[{ required: true, message: '请输入密码' }]"
        >
          <template #left-icon>
            <van-icon name="lock" />
          </template>
          <template #right-icon>
            <van-icon
              :name="showPassword ? 'eye' : 'eye-o'"
              @click="showPassword = !showPassword"
            />
          </template>
        </van-field>
        
        <div class="login-btn">
          <van-button
            type="primary"
            round
            block
            native-type="submit"
            :loading="loading"
            size="large"
          >
            登录
          </van-button>
        </div>
      </van-form>
      
      <div class="footer-links">
        <router-link to="/">返回用户端</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { adminLogin } from '@/api'
import { showToast } from 'vant'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  username: '',
  password: ''
})
const showPassword = ref(false)
const loading = ref(false)

const onSubmit = async (values) => {
  loading.value = true
  
  try {
    const res = await adminLogin(form.value)
    const data = res.data
    
    userStore.setToken(data.token)
    userStore.setUserInfo({
      id: data.id,
      username: data.username,
      isAdmin: true
    })
    
    showToast('登录成功')
    router.replace('/admin/dashboard')
  } catch (error) {
    console.error('Admin login failed:', error)
    showToast('登录失败，请检查用户名和密码')
  } finally {
    loading.value = false
  }
}
</script>

<style lang="less" scoped>
.admin-login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #4CAF50 0%, #2E7D32 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-container {
  width: 100%;
  max-width: 400px;
  background-color: #fff;
  border-radius: 16px;
  padding: 40px 24px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
}

.login-header {
  text-align: center;
  margin-bottom: 36px;
  
  .logo {
    width: 80px;
    height: 80px;
    background: linear-gradient(135deg, #4CAF50 0%, #81C784 100%);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 20px;
    color: #fff;
  }
  
  .title {
    font-size: 24px;
    font-weight: bold;
    color: #333;
    margin: 0 0 8px 0;
  }
  
  .subtitle {
    font-size: 14px;
    color: #999;
    margin: 0;
  }
}

.login-btn {
  margin-top: 32px;
}

.footer-links {
  text-align: center;
  margin-top: 24px;
  font-size: 14px;
  
  a {
    color: #4CAF50;
    text-decoration: none;
    
    &:hover {
      text-decoration: underline;
    }
  }
}

:deep(.van-cell) {
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
  
  &:last-child {
    border-bottom: none;
  }
}

:deep(.van-field__left-icon) {
  color: #999;
  margin-right: 12px;
}

:deep(.van-field__right-icon) {
  color: #999;
}
</style>
