<template>
  <div class="detail-page">
    <van-nav-bar title="公告详情" left-arrow @click-left="onClickLeft" fixed placeholder />
    
    <div class="detail-content" v-if="detail">
      <div class="detail-header">
        <h1 class="title">{{ detail.title }}</h1>
        <div class="meta">
          <span class="time">发布时间: {{ formatDate(detail.created_at) }}</span>
          <span class="views"><van-icon name="eye-o" /> {{ detail.view_count }}</span>
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
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getAnnouncementDetail } from '@/api'

const router = useRouter()
const route = useRoute()

const detail = ref(null)
const loading = ref(true)

const onClickLeft = () => {
  router.back()
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
}

const loadDetail = async () => {
  loading.value = true
  try {
    const res = await getAnnouncementDetail(route.params.id)
    detail.value = res.data
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
    font-size: 18px;
    font-weight: bold;
    color: #333;
    margin: 0 0 12px 0;
    line-height: 1.5;
  }
  
  .meta {
    display: flex;
    justify-content: space-between;
    font-size: 12px;
    color: #999;
    
    .views {
      display: flex;
      align-items: center;
      
      .van-icon {
        margin-right: 4px;
      }
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
  }
  
  :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
    margin: 16px 0 8px 0;
    font-weight: bold;
  }
  
  :deep(ul), :deep(ol) {
    padding-left: 20px;
    margin: 8px 0;
  }
  
  :deep(li) {
    margin: 4px 0;
  }
  
  :deep(blockquote) {
    border-left: 4px solid #4CAF50;
    padding-left: 12px;
    margin: 12px 0;
    color: #666;
  }
}

.loading-center {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px;
}
</style>
