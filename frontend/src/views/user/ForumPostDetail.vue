<template>
  <div class="post-detail-page">
    <van-nav-bar title="帖子详情" left-arrow @click-left="onClickLeft" fixed placeholder />
    
    <div class="post-content" v-if="post">
      <div class="post-header">
        <van-avatar :src="post.user?.avatar || defaultAvatar" size="44" />
        <div class="user-info">
          <span class="nickname">{{ post.user?.nickname || post.user?.username || '匿名用户' }}</span>
          <span class="time">{{ formatDate(post.created_at) }}</span>
        </div>
      </div>
      
      <h1 class="post-title">{{ post.title }}</h1>
      
      <div class="post-body" v-html="post.content"></div>
      
      <div class="post-images" v-if="post.images">
        <van-image
          v-for="(img, idx) in imageList(post.images)"
          :key="idx"
          :src="img"
          width="100"
          height="100"
          fit="cover"
          round
          @click="previewImage(img)"
        />
      </div>
      
      <div class="post-stats">
        <div class="stat-item">
          <van-icon name="eye-o" />
          <span>{{ post.view_count }}</span>
        </div>
        <div class="stat-item">
          <van-icon name="like-o" />
          <span>{{ post.like_count }}</span>
        </div>
        <div class="stat-item">
          <van-icon name="chat-o" />
          <span>{{ post.comment_count }}</span>
        </div>
      </div>
    </div>
    
    <div class="comments-section">
      <div class="section-header">
        <van-icon name="chat-o" />
        <span class="title">评论 ({{ commentTotal }})</span>
      </div>
      
      <van-list
        v-model:loading="loadingComments"
        :finished="finishedComments"
        finished-text="没有更多评论了"
        @load="loadComments"
      >
        <div class="comment-item" v-for="item in commentList" :key="item.id">
          <div class="comment-header">
            <van-avatar :src="item.user?.avatar || defaultAvatar" size="32" />
            <div class="user-info">
              <span class="nickname">{{ item.user?.nickname || item.user?.username || '匿名用户' }}</span>
              <span class="time">{{ formatDate(item.created_at) }}</span>
            </div>
          </div>
          <div class="comment-content">{{ item.content }}</div>
        </div>
        
        <van-empty v-if="!loadingComments && commentList.length === 0" description="暂无评论，快来抢沙发！" />
      </van-list>
    </div>
    
    <div class="comment-input-fixed">
      <van-field
        v-model="commentContent"
        placeholder="发表评论..."
        @keyup.enter="submitComment"
      >
        <template #button>
          <van-button type="primary" size="small" :disabled="!commentContent.trim()" @click="submitComment">
            发送
          </van-button>
        </template>
      </van-field>
    </div>
    
    <van-loading v-if="loading" class="loading-center" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getForumPostDetail, getForumCommentList, createForumComment } from '@/api'
import { showToast, showImagePreview } from 'vant'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const post = ref(null)
const loading = ref(true)
const commentList = ref([])
const loadingComments = ref(false)
const finishedComments = ref(false)
const commentPage = ref(1)
const commentTotal = ref(0)
const commentContent = ref('')
const pageSize = 10

const defaultAvatar = computed(() => {
  return 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=user%20avatar%20default%20icon%20simple&image_size=square'
})

const imageList = (imagesStr) => {
  if (!imagesStr) return []
  return imagesStr.split(',').filter(img => img.trim())
}

const onClickLeft = () => {
  router.back()
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

const previewImage = (url) => {
  showImagePreview([url])
}

const loadPost = async () => {
  loading.value = true
  try {
    const res = await getForumPostDetail(route.params.id)
    post.value = res.data
  } catch (error) {
    console.error('Load post failed:', error)
  } finally {
    loading.value = false
  }
}

const loadComments = async () => {
  try {
    const res = await getForumCommentList(route.params.id, { 
      page: commentPage.value, 
      page_size: pageSize 
    })
    const data = res.data?.list || []
    commentTotal.value = res.data?.total || 0
    
    commentList.value.push(...data)
    
    if (data.length < pageSize) {
      finishedComments.value = true
    } else {
      commentPage.value++
    }
  } catch (error) {
    console.error('Load comments failed:', error)
  } finally {
    loadingComments.value = false
  }
}

const submitComment = async () => {
  if (!commentContent.value.trim()) {
    showToast('请输入评论内容')
    return
  }
  
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  
  try {
    await createForumComment({
      post_id: route.params.id,
      content: commentContent.value
    })
    
    showToast('评论成功')
    commentContent.value = ''
    
    commentPage.value = 1
    finishedComments.value = false
    commentList.value = []
    loadComments()
    
    if (post.value) {
      post.value.comment_count++
    }
  } catch (error) {
    console.error('Submit comment failed:', error)
    showToast('评论失败')
  }
}

onMounted(() => {
  loadPost()
  loadComments()
})
</script>

<style lang="less" scoped>
.post-detail-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 60px;
}

.post-content {
  background-color: #fff;
  margin: 12px;
  padding: 16px;
  border-radius: 8px;
  
  .post-header {
    display: flex;
    align-items: center;
    margin-bottom: 16px;
    
    .user-info {
      margin-left: 12px;
      display: flex;
      flex-direction: column;
      
      .nickname {
        font-size: 15px;
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
  
  .post-title {
    font-size: 18px;
    font-weight: bold;
    color: #333;
    margin: 0 0 12px 0;
    line-height: 1.5;
  }
  
  .post-body {
    font-size: 15px;
    line-height: 1.8;
    color: #333;
    margin-bottom: 12px;
    
    :deep(p) {
      margin: 0 0 12px 0;
    }
    
    :deep(img) {
      max-width: 100%;
      height: auto;
      display: block;
      margin: 12px auto;
      border-radius: 4px;
    }
  }
  
  .post-images {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
    margin-bottom: 16px;
  }
  
  .post-stats {
    display: flex;
    padding-top: 12px;
    border-top: 1px solid #f5f5f5;
    
    .stat-item {
      display: flex;
      align-items: center;
      margin-right: 24px;
      font-size: 13px;
      color: #999;
      
      .van-icon {
        margin-right: 4px;
      }
    }
  }
}

.comments-section {
  background-color: #fff;
  margin: 12px;
  padding: 16px;
  border-radius: 8px;
  
  .section-header {
    display: flex;
    align-items: center;
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 1px solid #f5f5f5;
    
    .van-icon {
      color: #4CAF50;
      margin-right: 6px;
    }
    
    .title {
      font-size: 15px;
      font-weight: 500;
      color: #333;
    }
  }
  
  .comment-item {
    padding: 12px 0;
    border-bottom: 1px solid #f5f5f5;
    
    &:last-child {
      border-bottom: none;
    }
    
    .comment-header {
      display: flex;
      align-items: center;
      margin-bottom: 8px;
      
      .user-info {
        margin-left: 10px;
        display: flex;
        flex-direction: column;
        
        .nickname {
          font-size: 14px;
          color: #333;
        }
        
        .time {
          font-size: 11px;
          color: #999;
          margin-top: 2px;
        }
      }
    }
    
    .comment-content {
      font-size: 14px;
      color: #666;
      line-height: 1.6;
      padding-left: 42px;
    }
  }
}

.comment-input-fixed {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: #fff;
  padding: 10px 16px;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
  z-index: 100;
}

.loading-center {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px;
}
</style>
