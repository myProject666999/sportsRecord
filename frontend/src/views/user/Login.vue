<template>
  <div class="login-page">
    <van-nav-bar title="登录" />
    
    <div class="login-content">
      <div class="logo-area">
        <van-icon name="fire-o" size="60" color="#4CAF50" />
        <h2>运动健康</h2>
        <p class="slogan">记录每一步，健康每一天</p>
      </div>
      
      <van-form @submit="onSubmit">
        <van-cell-group inset>
          <van-field
            v-model="form.username"
            name="username"
            label="用户名"
            placeholder="请输入用户名"
            :rules="[{ required: true, message: '请填写用户名' }]"
          />
          <van-field
            v-model="form.password"
            type="password"
            name="password"
            label="密码"
            placeholder="请输入密码"
            :rules="[{ required: true, message: '请填写密码' }]"
          />
        </van-cell-group>
        
        <div style="margin: 16px;">
          <van-button round block type="primary" native-type="submit" :loading="loading">
            登录
          </van-button>
        </div>
      </van-form>
      
      <div class="footer-actions">
        <router-link to="/register">还没有账号？立即注册</router-link>
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
  password: ''
})

const loading = ref(false)

const onSubmit = async (values) => {
  loading.value = true
  try {
    await userStore.login(values.username, values.password)
    showToast('登录成功')
    router.push('/home')
  } catch (error) {
    console.error('Login failed:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style lang="less" scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
}

.login-content {
  padding: 32px 16px;
}

.logo-area {
  text-align: center;
  padding: 40px 0;
  
  h2 {
    margin-top: 16px;
    font-size: 24px;
    color: #333;
  }
  
  .slogan {
    margin-top: 8px;
    color: #666;
    font-size: 14px;
  }
}

.footer-actions {
  text-align: center;
  margin-top: 24px;
  
  a {
    color: #4CAF50;
    text-decoration: none;
    font-size: 14px;
  }
}
</style>
