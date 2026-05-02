<template>
  <div class="detail-page">
    <van-nav-bar title="科普详情" left-arrow @click-left="onClickLeft" fixed placeholder>
      <template #right>
        <van-icon :name="isCollected ? 'star' : 'star-o'" size="20" @click="toggleCollection" />
      </template>
    </van-nav-bar>
    
    <div class="detail-content" v-if="detail">
      <div class="detail-header">
        <h1 class="title">{{ detail.title }}</h1>
        <div class="meta">
          <span class="type" v-if="detail.type?.name">{{ detail.type.name }}</span>
          <span class="time">{{ formatDate(detail.created_at) }}</span>
          <span class="views"><van-icon name="eye-o" /> {{ detail.view_count }}</span>
        </div>
        <div class="author" v-if="detail.author">
          <van-icon name="manager-o" /> 作者：{{ detail.author }}
        </div>
      </div>
      
      <van-image
        v-if="detail.cover"
        :src="detail.cover"
        fit="cover"
        class="cover-image"
      />
      
      <div class="detail-body" v-html="detail.content"></div>
    </div>
    
    <van-loading v-if="loading" class="loading-center" />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getHealthScienceDetail, addCollection, removeCollection, checkCollection } from '@/api'
import { showToast } from 'vant'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const detail = ref(null)
const loading = ref(true)
const isCollected = ref(false)

const collectionType = 1

const onClickLeft = () => {
  router.back()
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
}

const checkIsCollected = async () => {
  if (!userStore.isLoggedIn) return
  
  try {
    const res = await checkCollection(collectionType, route.params.id)
    isCollected.value = res.data?.is_collected || false
  } catch (error) {
    console.error('Check collection failed:', error)
  }
}

const toggleCollection = () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  
  if (isCollected.value) {
    removeCollection(route.params.id)
      .then(() => {
        isCollected.value = false
        showToast('已取消收藏')
      })
      .catch((error) => {
        console.error('Remove collection failed:', error)
        showToast('操作失败')
      })
  } else {
    addCollection({
      type: collectionType,
      relation_id: route.params.id
    })
      .then(() => {
        isCollected.value = true
        showToast('收藏成功')
      })
      .catch((error) => {
        console.error('Add collection failed:', error)
        showToast('收藏失败')
      })
  }
}

const loadDetail = async () => {
  loading.value = true
  try {
    const res = await getHealthScienceDetail(route.params.id)
    detail.value = res.data
    checkIsCollected()
  } catch (error) {
    console.error('Load detail failed:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDetail()
})
</script>

<style lang="less" scoped>
.detail-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.detail-content {
  background-color: #fff;
  padding: 16px;
  margin: 12px;
  border-radius: 8px;
}

.detail-header {
  margin-bottom: 16px;
  
  .title {
    font-size: 20px;
    font-weight: bold;
    color: #333;
    margin: 0 0 12px 0;
    line-height: 1.5;
  }
  
  .meta {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 8px;
    
    .type {
      font-size: 12px;
      color: #4CAF50;
      background-color: rgba(76, 175, 80, 0.1);
      padding: 2px 8px;
      border-radius: 4px;
    }
    
    .time, .views {
      font-size: 12px;
      color: #999;
      display: flex;
      align-items: center;
      
      .van-icon {
        margin-right: 4px;
      }
    }
  }
  
  .author {
    font-size: 13px;
    color: #666;
    display: flex;
    align-items: center;
    
    .van-icon {
      margin-right: 4px;
    }
  }
}

.cover-image {
  width: 100%;
  height: 200px;
  margin-bottom: 16px;
  border-radius: 8px;
  overflow: hidden;
}

.detail-body {
  font-size: 15px;
  line-height: 1.8;
  color: #333;
  
  :deep(img) {
    max-width: 100%;
    height: auto;
    display: block;
    margin: 12px auto;
    border-radius: 4px;
  }
  
  :deep(p) {
    margin: 0 0 12px 0;
    text-indent: 2em;
  }
  
  :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
    margin: 16px 0 8px 0;
    font-weight: bold;
    text-indent: 0;
  }
  
  :deep(ul), :deep(ol) {
    padding-left: 20px;
    margin: 8px 0;
    text-indent: 0;
  }
  
  :deep(li) {
    margin: 4px 0;
    text-indent: 0;
  }
  
  :deep(blockquote) {
    border-left: 4px solid #4CAF50;
    padding-left: 12px;
    margin: 12px 0;
    color: #666;
    text-indent: 0;
  }
}

.loading-center {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px;
}
</style>
