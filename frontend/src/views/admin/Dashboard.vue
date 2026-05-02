<template>
  <div class="admin-dashboard-page">
    <van-nav-bar title="管理后台" fixed placeholder>
      <template #right>
        <van-button size="small" type="default" @click="onLogout">
          退出
        </van-button>
      </template>
    </van-nav-bar>
    
    <div class="dashboard-content">
      <div class="stats-cards">
        <div class="stats-card user-card" @click="goTo('users')">
          <div class="card-icon">
            <van-icon name="user-o" size="32" />
          </div>
          <div class="card-info">
            <span class="card-num">{{ stats.user_count || 0 }}</span>
            <span class="card-label">用户总数</span>
          </div>
        </div>
        
        <div class="stats-card step-card" @click="goTo('step-record')">
          <div class="card-icon">
            <van-icon name="orders-o" size="32" />
          </div>
          <div class="card-info">
            <span class="card-num">{{ stats.total_steps || 0 }}</span>
            <span class="card-label">总步数</span>
          </div>
        </div>
        
        <div class="stats-card post-card" @click="goTo('forum')">
          <div class="card-icon">
            <van-icon name="chat-o" size="32" />
          </div>
          <div class="card-info">
            <span class="card-num">{{ stats.post_count || 0 }}</span>
            <span class="card-label">论坛帖子</span>
          </div>
        </div>
        
        <div class="stats-card view-card">
          <div class="card-icon">
            <van-icon name="eye-o" size="32" />
          </div>
          <div class="card-info">
            <span class="card-num">{{ stats.total_views || 0 }}</span>
            <span class="card-label">总浏览量</span>
          </div>
        </div>
      </div>
      
      <div class="quick-actions">
        <h3 class="section-title">快捷操作</h3>
        <div class="action-grid">
          <div class="action-item" @click="goTo('users')">
            <van-icon name="user-o" size="28" color="#4CAF50" />
            <span>用户管理</span>
          </div>
          <div class="action-item" @click="goTo('forum')">
            <van-icon name="chat-o" size="28" color="#4CAF50" />
            <span>论坛管理</span>
          </div>
          <div class="action-item" @click="goTo('announcement')">
            <van-icon name="info-o" size="28" color="#4CAF50" />
            <span>公告管理</span>
          </div>
          <div class="action-item" @click="goTo('carousel')">
            <van-icon name="photo-o" size="28" color="#4CAF50" />
            <span>轮播图</span>
          </div>
          <div class="action-item" @click="goTo('health-science')">
            <van-icon name="medicines-o" size="28" color="#4CAF50" />
            <span>科普管理</span>
          </div>
          <div class="action-item" @click="goTo('teaching-video')">
            <van-icon name="play-circle-o" size="28" color="#4CAF50" />
            <span>视频管理</span>
          </div>
          <div class="action-item" @click="goTo('step-record')">
            <van-icon name="orders-o" size="28" color="#4CAF50" />
            <span>步数统计</span>
          </div>
          <div class="action-item" @click="goTo('type')">
            <van-icon name="apps-o" size="28" color="#4CAF50" />
            <span>分类管理</span>
          </div>
        </div>
      </div>
      
      <div class="recent-data">
        <h3 class="section-title">最近注册用户</h3>
        <van-loading v-if="loadingUsers" class="loading-small" />
        <div class="user-list" v-else>
          <div class="user-item" v-for="user in recentUsers" :key="user.id">
            <van-avatar :src="user.avatar || defaultAvatar" size="40" />
            <div class="user-info">
              <span class="username">{{ user.nickname || user.username }}</span>
              <span class="time">注册于 {{ formatDate(user.created_at) }}</span>
            </div>
          </div>
          <van-empty v-if="recentUsers.length === 0" description="暂无用户数据" />
        </div>
      </div>
    </div>
    
    <van-tabbar v-model="activeTab" route placeholder>
      <van-tabbar-item replace to="/admin/dashboard" icon="home-o">首页</van-tabbar-item>
      <van-tabbar-item replace to="/admin/users" icon="user-o">用户</van-tabbar-item>
      <van-tabbar-item replace to="/admin/forum" icon="chat-o">论坛</van-tabbar-item>
      <van-tabbar-item replace to="/admin/step-record" icon="chart-trending-o">统计</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getAdminDashboardStats, getUserList } from '@/api'
import { showConfirmDialog, showToast } from 'vant'

const router = useRouter()
const userStore = useUserStore()

const activeTab = ref(0)
const stats = ref({
  user_count: 0,
  total_steps: 0,
  post_count: 0,
  total_views: 0
})
const recentUsers = ref([])
const loadingUsers = ref(false)

const defaultAvatar = computed(() => {
  return 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=user%20avatar%20default%20icon%20simple&image_size=square'
})

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`
}

const goTo = (path) => {
  router.push(`/admin/${path}`)
}

const loadStats = async () => {
  try {
    const res = await getAdminDashboardStats()
    stats.value = res.data || stats.value
  } catch (error) {
    console.error('Load stats failed:', error)
  }
}

const loadRecentUsers = async () => {
  loadingUsers.value = true
  try {
    const res = await getUserList({ page: 1, page_size: 5 })
    recentUsers.value = res.data?.list || []
  } catch (error) {
    console.error('Load recent users failed:', error)
  } finally {
    loadingUsers.value = false
  }
}

const onLogout = () => {
  showConfirmDialog({
    title: '提示',
    message: '确定要退出登录吗？',
  })
    .then(() => {
      userStore.logout()
      showToast('已退出登录')
      router.replace('/admin/login')
    })
    .catch(() => {
      // on cancel
    })
}

onMounted(() => {
  loadStats()
  loadRecentUsers()
})
</script>

<style lang="less" scoped>
.admin-dashboard-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 50px;
}

.dashboard-content {
  padding: 16px;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.stats-card {
  background-color: #fff;
  border-radius: 12px;
  padding: 16px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  
  .card-icon {
    width: 56px;
    height: 56px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 12px;
  }
  
  &.user-card .card-icon {
    background-color: rgba(76, 175, 80, 0.1);
    color: #4CAF50;
  }
  
  &.step-card .card-icon {
    background-color: rgba(33, 150, 243, 0.1);
    color: #2196F3;
  }
  
  &.post-card .card-icon {
    background-color: rgba(255, 152, 0, 0.1);
    color: #FF9800;
  }
  
  &.view-card .card-icon {
    background-color: rgba(156, 39, 176, 0.1);
    color: #9C27B0;
  }
  
  .card-info {
    display: flex;
    flex-direction: column;
    
    .card-num {
      font-size: 24px;
      font-weight: bold;
      color: #333;
    }
    
    .card-label {
      font-size: 12px;
      color: #999;
      margin-top: 4px;
    }
  }
}

.section-title {
  font-size: 16px;
  font-weight: bold;
  color: #333;
  margin: 0 0 16px 0;
}

.quick-actions {
  background-color: #fff;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 20px;
  
  .action-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 16px;
  }
  
  .action-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    
    span {
      font-size: 12px;
      color: #666;
      margin-top: 8px;
    }
  }
}

.recent-data {
  background-color: #fff;
  border-radius: 12px;
  padding: 16px;
  
  .user-list {
    .user-item {
      display: flex;
      align-items: center;
      padding: 12px 0;
      border-bottom: 1px solid #f5f5f5;
      
      &:last-child {
        border-bottom: none;
      }
      
      .user-info {
        display: flex;
        flex-direction: column;
        margin-left: 12px;
        
        .username {
          font-size: 14px;
          color: #333;
        }
        
        .time {
          font-size: 12px;
          color: #999;
          margin-top: 4px;
        }
      }
    }
  }
}

.loading-small {
  display: flex;
  justify-content: center;
  padding: 20px;
}
</style>
