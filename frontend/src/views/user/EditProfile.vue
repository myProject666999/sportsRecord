<template>
  <div class="edit-profile-page">
    <van-nav-bar title="编辑资料" left-arrow @click-left="onClickLeft" fixed placeholder />
    
    <van-form @submit="onSubmit">
      <van-cell-group inset>
        <van-field
          v-model="form.username"
          label="用户名"
          placeholder="请输入用户名"
          :disabled="true"
        />
        <van-field
          v-model="form.nickname"
          label="昵称"
          placeholder="请输入昵称"
        />
        <van-field
          v-model="form.phone"
          label="手机号"
          placeholder="请输入手机号"
          type="tel"
        />
        <van-field
          v-model="form.email"
          label="邮箱"
          placeholder="请输入邮箱"
          type="email"
        />
        <van-field
          v-model="form.avatar"
          label="头像"
          placeholder="请输入头像URL"
        >
          <template #button>
            <van-image
              v-if="form.avatar"
              :src="form.avatar"
              width="40"
              height="40"
              fit="cover"
              round
            />
            <van-icon v-else name="add" color="#ccc" size="24" />
          </template>
        </van-field>
      </van-cell-group>
      
      <div class="submit-btn">
        <van-button type="primary" round block native-type="submit">
          保存修改
        </van-button>
      </div>
    </van-form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getUserInfo, updateUserInfo } from '@/api'
import { showToast } from 'vant'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  username: '',
  nickname: '',
  phone: '',
  email: '',
  avatar: ''
})

const onClickLeft = () => {
  router.back()
}

const loadUserInfo = async () => {
  try {
    const res = await getUserInfo()
    const user = res.data
    form.value = {
      username: user.username || '',
      nickname: user.nickname || '',
      phone: user.phone || '',
      email: user.email || '',
      avatar: user.avatar || ''
    }
  } catch (error) {
    console.error('Load user info failed:', error)
  }
}

const onSubmit = async () => {
  try {
    await updateUserInfo(form.value)
    showToast('保存成功')
    userStore.setUserInfo({
      ...userStore.userInfo,
      ...form.value
    })
    router.back()
  } catch (error) {
    console.error('Update profile failed:', error)
    showToast('保存失败，请重试')
  }
}

onMounted(() => {
  loadUserInfo()
})
</script>

<style lang="less" scoped>
.edit-profile-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 24px;
}

.submit-btn {
  padding: 24px 16px;
}
</style>
