<template>
  <div class="profile-page">
    <div class="profile-header">
      <van-nav-bar fixed placeholder>
        <template #right>
          <van-icon name="setting-o" size="20" @click="onSetting" />
        </template>
      </van-nav-bar>
      
      <div class="user-info" v-if="userStore.isLoggedIn">
        <van-avatar :src="userStore.userInfo.avatar || defaultAvatar" size="64" />
        <div class="info-text">
          <h3 class="nickname">{{ userStore.userInfo.nickname || userStore.userInfo.username || '用户' }}</h3>
          <p class="desc" v-if="userStore.userInfo.phone">
            <van-icon name="phone-o" /> {{ hidePhone(userStore.userInfo.phone) }}
          </p>
        </div>
        <van-icon name="arrow" class="edit-icon" @click="goEditProfile" />
      </div>
      
      <div class="user-info login-prompt" v-else @click="goLogin">
        <van-avatar :src="defaultAvatar" size="64" />
        <div class="info-text">
          <h3 class="nickname">点击登录</h3>
          <p class="desc">登录后享受更多功能</p>
        </div>
        <van-icon name="arrow" class="edit-icon" />
      </div>
    </div>
    
    <div class="profile-content">
      <van-cell-group inset class="menu-group">
        <van-cell title="步数记录" icon="orders-o" is-link @click="goStepRecord" />
        <van-cell title="我的收藏" icon="star-o" is-link @click="goCollection" />
        <van-cell title="浏览历史" icon="history" is-link @click="onHistory" />
      </van-cell-group>
      
      <van-cell-group inset class="menu-group">
        <van-cell title="编辑资料" icon="edit" is-link @click="goEditProfile" v-if="userStore.isLoggedIn" />
        <van-cell title="修改密码" icon="lock" is-link @click="goChangePassword" v-if="userStore.isLoggedIn" />
        <van-cell title="关于我们" icon="info-o" is-link @click="onAbout" />
        <van-cell title="意见反馈" icon="edit-pen" is-link @click="onFeedback" />
      </van-cell-group>
      
      <div class="logout-btn" v-if="userStore.isLoggedIn">
        <van-button round block type="danger" @click="onLogout">退出登录</van-button>
      </div>
    </div>
    
    <van-tabbar v-model="activeTab" route placeholder>
      <van-tabbar-item replace to="/home" icon="home-o">首页</van-tabbar-item>
      <van-tabbar-item replace to="/forum" icon="chat-o">论坛</van-tabbar-item>
      <van-tabbar-item replace to="/health-science" icon="medicines-o">科普</van-tabbar-item>
      <van-tabbar-item replace to="/profile" icon="user-o">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { showConfirmDialog, showToast } from 'vant'

const router = useRouter()
const userStore = useUserStore()

const activeTab = ref(3)

const defaultAvatar = computed(() => {
  return 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=user%20avatar%20default%20icon%20simple&image_size=square'
})

const hidePhone = (phone) => {
  if (!phone) return ''
  if (phone.length < 7) return phone
  return phone.substring(0, 3) + '****' + phone.substring(phone.length - 4)
}

const goLogin = () => {
  router.push('/login')
}

const goEditProfile = () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  router.push('/profile/edit')
}

const goStepRecord = () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  router.push('/step-record')
}

const goCollection = () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  router.push('/collection')
}

const goChangePassword = () => {
  router.push('/password/change')
}

const onSetting = () => {
  showToast('设置功能开发中')
}

const onHistory = () => {
  showToast('浏览历史功能开发中')
}

const onAbout = () => {
  showToast('关于我们功能开发中')
}

const onFeedback = () => {
  showToast('意见反馈功能开发中')
}

const onLogout = () => {
  showConfirmDialog({
    title: '提示',
    message: '确定要退出登录吗？',
  })
    .then(() => {
      userStore.logout()
      showToast('已退出登录')
      router.replace('/login')
    })
    .catch(() => {
      // on cancel
    })
}
</script>

<style lang="less" scoped>
.profile-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 50px;
}

.profile-header {
  background: linear-gradient(135deg, #4CAF50 0%, #81C784 100%);
  padding: 20px 16px 30px 16px;
  margin-bottom: -20px;
}

.user-info {
  display: flex;
  align-items: center;
  padding: 16px;
  background-color: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  margin-top: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  
  &.login-prompt {
    cursor: pointer;
  }
  
  .info-text {
    flex: 1;
    margin-left: 16px;
    
    .nickname {
      font-size: 18px;
      font-weight: bold;
      color: #333;
      margin: 0 0 4px 0;
    }
    
    .desc {
      font-size: 14px;
      color: #999;
      margin: 0;
      display: flex;
      align-items: center;
      
      .van-icon {
        margin-right: 4px;
      }
    }
  }
  
  .edit-icon {
    color: #999;
    font-size: 16px;
  }
}

.profile-content {
  padding: 20px 0;
  
  .menu-group {
    margin-bottom: 16px;
    
    :deep(.van-cell) {
      font-size: 15px;
    }
  }
  
  .logout-btn {
    padding: 16px;
    margin-top: 24px;
  }
}
</style>
