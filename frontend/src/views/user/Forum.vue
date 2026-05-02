<template>
  <div class="forum-page">
    <van-nav-bar title="论坛交流" fixed placeholder>
      <template #right>
        <van-icon name="edit" size="20" @click="goCreate" />
      </template>
    </van-nav-bar>
    
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <div class="post-item" v-for="item in list" :key="item.id" @click="goDetail(item.id)">
          <div class="post-header">
            <van-avatar :src="item.user?.avatar || defaultAvatar" size="40" />
            <div class="user-info">
              <span class="nickname">{{ item.user?.nickname || item.user?.username || '匿名用户' }}</span>
              <span class="time">{{ formatDate(item.created_at) }}</span>
            </div>
          </div>
          <div class="post-content">
            <h3 class="title">{{ item.title }}</h3>
            <p class="desc" v-if="item.content">{{ item.content.length > 100 ? item.content.substring(0, 100) + '...' : item.content }}</p>
            <div class="post-images" v-if="item.images">
              <van-image
                v-for="(img, idx) in imageList(item.images)"
                :key="idx"
                :src="img"
                width="80"
                height="80"
                fit="cover"
                round
              />
            </div>
          </div>
          <div class="post-footer">
            <div class="stat-item">
              <van-icon name="eye-o" />
              <span>{{ item.view_count }}</span>
            </div>
            <div class="stat-item">
              <van-icon name="like-o" />
              <span>{{ item.like_count }}</span>
            </div>
            <div class="stat-item">
              <van-icon name="chat-o" />
              <span>{{ item.comment_count }}</span>
            </div>
          </div>
        </div>
        
        <van-empty v-if="!loading && list.length === 0" description="暂无帖子，快来发布第一篇吧！" />
      </van-list>
    </van-pull-refresh>
    
    <van-tabbar v-model="activeTab" route placeholder>
      <van-tabbar-item replace to="/home" icon="home-o">首页</van-tabbar-item>
      <van-tabbar-item replace to="/forum" icon="chat-o">论坛</van-tabbar-item>
      <van-tabbar-item replace to="/health-science" icon="medicines-o">科普</van-tabbar-item>
      <van-tabbar-item replace to="/profile" icon="user-o">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getForumPostList } from '@/api'
import { showToast } from 'vant'

const router = useRouter()
const userStore = useUserStore()

const list = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const activeTab = ref(1)
const page = ref(1)
const pageSize = 10

const defaultAvatar = computed(() => {
  return 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=user%20avatar%20default%20icon%20simple&image_size=square'
})

const imageList = (imagesStr) => {
  if (!imagesStr) return []
  return imagesStr.split(',').filter(img => img.trim())
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now - date
  const minuteDiff = Math.floor(diff / (1000 * 60))
  const hourDiff = Math.floor(diff / (1000 * 60 * 60))
  const dayDiff = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (minuteDiff < 1) {
    return '刚刚'
  } else if (minuteDiff < 60) {
    return `${minuteDiff}分钟前`
  } else if (hourDiff < 24) {
    return `${hourDiff}小时前`
  } else if (dayDiff < 7) {
    return `${dayDiff}天前`
  } else {
    return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`
  }
}

const goCreate = () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  router.push('/forum/create')
}

const goDetail = (id) => {
  router.push(`/forum/${id}`)
}

const loadData = async () => {
  try {
    const res = await getForumPostList({ page: page.value, page_size: pageSize })
    const data = res.data?.list || []
    
    if (refreshing.value) {
      list.value = []
    }
    
    list.value.push(...data)
    
    if (data.length < pageSize) {
      finished.value = true
    } else {
      page.value++
    }
  } catch (error) {
    console.error('Load forum posts failed:', error)
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

const onLoad = () => {
  loadData()
}

const onRefresh = () => {
  page.value = 1
  finished.value = false
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style lang="less" scoped>
.forum-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 50px;
}

.post-item {
  background-color: #fff;
  margin: 12px;
  padding: 16px;
  border-radius: 8px;
  
  .post-header {
    display: flex;
    align-items: center;
    margin-bottom: 12px;
    
    .user-info {
      margin-left: 12px;
      display: flex;
      flex-direction: column;
      
      .nickname {
        font-size: 14px;
        font-weight: 500;
        color: #333;
      }
      
      .time {
        font-size: 12px;
        color: #999;
        margin-top: 2px;
      }
    }
  }
  
  .post-content {
    .title {
      font-size: 16px;
      font-weight: bold;
      color: #333;
      margin: 0 0 8px 0;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    
    .desc {
      font-size: 14px;
      color: #666;
      line-height: 1.6;
      margin: 0 0 12px 0;
    }
    
    .post-images {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
    }
  }
  
  .post-footer {
    display: flex;
    margin-top: 12px;
    padding-top: 12px;
    border-top: 1px solid #f5f5f5;
    
    .stat-item {
      display: flex;
      align-items: center;
      margin-right: 24px;
      font-size: 12px;
      color: #999;
      
      .van-icon {
        margin-right: 4px;
      }
    }
  }
}
</style>
