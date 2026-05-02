<template>
  <div class="manage-page">
    <van-nav-bar title="论坛管理" left-arrow @click-left="onClickLeft" fixed placeholder />
    
    <van-tabs v-model:active="activeType">
      <van-tab title="帖子">
        <van-pull-refresh v-model="refreshing" @refresh="onRefresh('posts')">
          <van-list
            v-model:loading="loadingPosts"
            :finished="finishedPosts"
            finished-text="没有更多了"
            @load="onLoad('posts')"
          >
            <van-cell-group inset>
              <div class="forum-item" v-for="item in postList" :key="item.id">
                <div class="item-header">
                  <van-avatar :src="item.user?.avatar || defaultAvatar" size="36" />
                  <div class="user-info">
                    <span class="username">{{ item.user?.nickname || item.user?.username }}</span>
                    <van-tag 
                      :type="item.status === 1 ? 'success' : 'danger'" 
                      size="mini"
                    >
                      {{ item.status === 1 ? '正常' : '已隐藏' }}
                    </van-tag>
                  </div>
                </div>
                <div class="item-content" @click="viewPost(item)">
                  <h4 class="title">{{ item.title }}</h4>
                  <p class="desc">{{ item.content?.substring(0, 80) || '' }}...</p>
                </div>
                <div class="item-stats">
                  <span><van-icon name="eye-o" /> {{ item.view_count }}</span>
                  <span><van-icon name="like-o" /> {{ item.like_count }}</span>
                  <span><van-icon name="chat-o" /> {{ item.comment_count }}</span>
                  <span class="time">{{ formatDate(item.created_at) }}</span>
                </div>
                <div class="item-actions">
                  <van-button 
                    size="mini" 
                    :type="item.status === 1 ? 'warning' : 'success'"
                    @click="toggleStatus(item)"
                  >
                    {{ item.status === 1 ? '隐藏' : '显示' }}
                  </van-button>
                  <van-button size="mini" type="danger" @click="onDeletePost(item)">删除</van-button>
                </div>
              </div>
            </van-cell-group>
            <van-empty v-if="!loadingPosts && postList.length === 0" description="暂无帖子" />
          </van-list>
        </van-pull-refresh>
      </van-tab>
      
      <van-tab title="评论">
        <van-pull-refresh v-model="refreshing" @refresh="onRefresh('comments')">
          <van-list
            v-model:loading="loadingComments"
            :finished="finishedComments"
            finished-text="没有更多了"
            @load="onLoad('comments')"
          >
            <van-cell-group inset>
              <div class="comment-item" v-for="item in commentList" :key="item.id">
                <div class="item-header">
                  <van-avatar :src="item.user?.avatar || defaultAvatar" size="32" />
                  <div class="user-info">
                    <span class="username">{{ item.user?.nickname || item.user?.username }}</span>
                    <span class="time">{{ formatDate(item.created_at) }}</span>
                  </div>
                </div>
                <div class="comment-content">
                  {{ item.content }}
                </div>
                <div class="item-actions">
                  <van-button size="mini" type="danger" @click="onDeleteComment(item)">删除</van-button>
                </div>
              </div>
            </van-cell-group>
            <van-empty v-if="!loadingComments && commentList.length === 0" description="暂无评论" />
          </van-list>
        </van-pull-refresh>
      </van-tab>
    </van-tabs>
    
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
import { getAdminForumPostList, updateForumPostStatus, deleteForumPost, getAdminForumComments, deleteForumComment } from '@/api'
import { showToast, showConfirmDialog } from 'vant'

const router = useRouter()

const activeTab = ref(2)
const activeType = ref(0)
const refreshing = ref(false)

const postList = ref([])
const loadingPosts = ref(false)
const finishedPosts = ref(false)
const postPage = ref(1)

const commentList = ref([])
const loadingComments = ref(false)
const finishedComments = ref(false)
const commentPage = ref(1)

const pageSize = 10

const defaultAvatar = computed(() => {
  return 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=user%20avatar%20default%20icon%20simple&image_size=square'
})

const onClickLeft = () => {
  router.back()
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`
}

const viewPost = (item) => {
  showToast('帖子预览功能开发中')
}

const toggleStatus = (item) => {
  const newStatus = item.status === 1 ? 0 : 1
  const action = newStatus === 1 ? '显示' : '隐藏'
  
  showConfirmDialog({
    title: '提示',
    message: `确定要${action}该帖子吗？`,
  })
    .then(async () => {
      try {
        await updateForumPostStatus(item.id, { status: newStatus })
        item.status = newStatus
        showToast(`${action}成功`)
      } catch (error) {
        console.error('Update status failed:', error)
        showToast('操作失败')
      }
    })
    .catch(() => {})
}

const onDeletePost = (item) => {
  showConfirmDialog({
    title: '提示',
    message: '确定要删除该帖子吗？',
  })
    .then(async () => {
      try {
        await deleteForumPost(item.id)
        postList.value = postList.value.filter(p => p.id !== item.id)
        showToast('删除成功')
      } catch (error) {
        console.error('Delete post failed:', error)
        showToast('删除失败')
      }
    })
    .catch(() => {})
}

const onDeleteComment = (item) => {
  showConfirmDialog({
    title: '提示',
    message: '确定要删除该评论吗？',
  })
    .then(async () => {
      try {
        await deleteForumComment(item.id)
        commentList.value = commentList.value.filter(c => c.id !== item.id)
        showToast('删除成功')
      } catch (error) {
        console.error('Delete comment failed:', error)
        showToast('删除失败')
      }
    })
    .catch(() => {})
}

const loadPosts = async () => {
  try {
    const res = await getAdminForumPostList({ page: postPage.value, page_size: pageSize })
    const data = res.data?.list || []
    
    if (refreshing.value) {
      postList.value = []
    }
    
    postList.value.push(...data)
    
    if (data.length < pageSize) {
      finishedPosts.value = true
    } else {
      postPage.value++
    }
  } catch (error) {
    console.error('Load posts failed:', error)
  } finally {
    loadingPosts.value = false
    refreshing.value = false
  }
}

const loadComments = async () => {
  try {
    const res = await getAdminForumComments({ page: commentPage.value, page_size: pageSize })
    const data = res.data?.list || []
    
    if (refreshing.value) {
      commentList.value = []
    }
    
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
    refreshing.value = false
  }
}

const onLoad = (type) => {
  if (type === 'posts') {
    loadPosts()
  } else {
    loadComments()
  }
}

const onRefresh = (type) => {
  if (type === 'posts') {
    postPage.value = 1
    finishedPosts.value = false
    refreshing.value = true
    loadPosts()
  } else {
    commentPage.value = 1
    finishedComments.value = false
    refreshing.value = true
    loadComments()
  }
}

onMounted(() => {
  loadPosts()
})
</script>

<style lang="less" scoped>
.manage-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 50px;
}

.forum-item {
  padding: 12px 16px;
  border-bottom: 1px solid #f5f5f5;
  
  &:last-child {
    border-bottom: none;
  }
  
  .item-header {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
    
    .user-info {
      margin-left: 10px;
      display: flex;
      align-items: center;
      gap: 8px;
      
      .username {
        font-size: 14px;
        color: #333;
      }
    }
  }
  
  .item-content {
    .title {
      font-size: 15px;
      font-weight: 500;
      color: #333;
      margin: 0 0 6px 0;
    }
    
    .desc {
      font-size: 13px;
      color: #666;
      margin: 0;
      line-height: 1.5;
    }
  }
  
  .item-stats {
    display: flex;
    align-items: center;
    gap: 16px;
    margin: 8px 0;
    font-size: 12px;
    color: #999;
    
    span {
      display: flex;
      align-items: center;
      
      .van-icon {
        margin-right: 4px;
      }
    }
    
    .time {
      margin-left: auto;
    }
  }
  
  .item-actions {
    display: flex;
    gap: 8px;
  }
}

.comment-item {
  padding: 12px 16px;
  border-bottom: 1px solid #f5f5f5;
  
  &:last-child {
    border-bottom: none;
  }
  
  .item-header {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
    
    .user-info {
      margin-left: 10px;
      display: flex;
      flex-direction: column;
      
      .username {
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
    color: #333;
    line-height: 1.6;
    margin-bottom: 8px;
  }
  
  .item-actions {
    display: flex;
    justify-content: flex-end;
  }
}
</style>
