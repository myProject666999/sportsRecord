<template>
  <div class="register-page">
    <van-nav-bar title="注册" left-arrow @click-left="onClickLeft" />
    
    <div class="register-content">
      <div class="logo-area">
        <van-icon name="fire-o" size="50" color="#4CAF50" />
        <h2>加入运动健康</h2>
      </div>
      
      <van-form @submit="onSubmit">
        <van-cell-group inset>
          <van-field
            v-model="form.username"
            name="username"
            label="用户名"
            placeholder="请输入用户名（3-50个字符）"
            :rules="[
              { required: true, message: '请填写用户名' },
              { min: 3, max: 50, message: '用户名长度为3-50个字符' }
            ]"
          />
          <van-field
            v-model="form.nickname"
            name="nickname"
            label="昵称"
            placeholder="请输入昵称（选填）"
          />
          <van-field
            v-model="form.password"
            type="password"
            name="password"
            label="密码"
            placeholder="请输入密码（至少6位）"
            :rules="[
              { required: true, message: '请填写密码' },
              { min: 6, message: '密码至少6位' }
            ]"
          />
          <van-field
            v-model="form.confirmPassword"
            type="password"
            name="confirmPassword"
            label="确认密码"
            placeholder="请再次输入密码"
            :rules="[
              { required: true, message: '请确认密码' },
              { validator: validateConfirmPassword, message: '两次输入的密码不一致' }
            ]"
          />
        </van-cell-group>
        
        <div style="margin: 16px;">
          <van-button round block type="primary" native-type="submit" :loading="loading">
            注册
          </van-button>
        </div>
      </van-form>
      
      <div class="footer-actions">
        <router-link to="/login">已有账号？立即登录</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { showToast } from 'vant'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  username: '',
  nickname: '',
  password: '',
  confirmPassword: ''
})

const loading = ref(false)

const validateConfirmPassword = (value) => {
  return value === form.value.password
}

const onClickLeft = () => {
  router.back()
}

const onSubmit = async (values) => {
  loading.value = true
  try {
    await userStore.register(values.username, values.password, values.nickname)
    showToast('注册成功')
    router.push('/home')
  } catch (error) {
    console.error('Register failed:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style lang="less" scoped>
.register-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
}

.register-content {
  padding: 16px;
}

.logo-area {
  text-align: center;
  padding: 24px 0;
  
  h2 {
    margin-top: 12px;
    font-size: 20px;
    color: #333;
  }
}

.footer-actions {
  text-align: center;
  margin-top: 16px;
  
  a {
    color: #4CAF50;
    text-decoration: none;
    font-size: 14px;
  }
}
</style>
