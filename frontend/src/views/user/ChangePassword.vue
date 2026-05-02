<template>
  <div class="change-password-page">
    <van-nav-bar title="修改密码" left-arrow @click-left="onClickLeft" fixed placeholder />
    
    <van-form @submit="onSubmit">
      <van-cell-group inset>
        <van-field
          v-model="form.old_password"
          name="old_password"
          label="原密码"
          placeholder="请输入原密码"
          type="password"
          :rules="[{ required: true, message: '请输入原密码' }]"
        />
        <van-field
          v-model="form.new_password"
          name="new_password"
          label="新密码"
          placeholder="请输入新密码（至少6位）"
          type="password"
          :rules="[{ required: true, message: '请输入新密码' }, { min: 6, message: '密码长度不能少于6位' }]"
        />
        <van-field
          v-model="form.confirm_password"
          name="confirm_password"
          label="确认密码"
          placeholder="请再次输入新密码"
          type="password"
          :rules="[{ required: true, message: '请再次输入新密码' }]"
        />
      </van-cell-group>
      
      <div class="submit-btn">
        <van-button type="primary" round block native-type="submit">
          确认修改
        </van-button>
      </div>
    </van-form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { changePassword } from '@/api'
import { showToast, showConfirmDialog } from 'vant'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const onClickLeft = () => {
  router.back()
}

const onSubmit = async (values) => {
  if (form.value.new_password !== form.value.confirm_password) {
    showToast('两次输入的密码不一致')
    return
  }
  
  try {
    await changePassword({
      old_password: form.value.old_password,
      new_password: form.value.new_password
    })
    showToast('密码修改成功')
    
    showConfirmDialog({
      title: '提示',
      message: '密码已修改，请重新登录',
      confirmButtonText: '去登录',
    })
      .then(() => {
        userStore.logout()
        router.replace('/login')
      })
  } catch (error) {
    console.error('Change password failed:', error)
    showToast('密码修改失败，请检查原密码是否正确')
  }
}
</script>

<style lang="less" scoped>
.change-password-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 24px;
}

.submit-btn {
  padding: 24px 16px;
}
</style>
