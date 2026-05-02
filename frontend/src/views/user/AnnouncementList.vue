<template>
  <div class="announcement-page">
    <van-nav-bar title="公告信息" left-arrow @click-left="onClickLeft" fixed placeholder />
    
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <div class="announcement-item" v-for="item in list" :key="item.id" @click="goDetail(item.id)">
          <div class="item-header">
            <span class="title">{{ item.title }}</span>
            <van-tag type="primary" v-if="item.view_count > 100">热门</van-tag>
          </div>
          <div class="item-content">
            <van-image
              v-if="item.cover"
              :src="item.cover"
              round
              width="80"
              height="60"
              fit="cover"
            />
            <div class="content-text" :class="{ 'no-cover': !item.cover }">
              <p class="desc">{{ stripHtml(item.content) }}</p>
              <div class="item-footer">
                <span class="time">{{ formatDate(item.created_at) }}</span>
                <span class="views"><van-icon name="eye-o" /> {{ item.view_count }}</span>
              </div>
            </div>
          </div>
        </div>
        
        <van-empty v-if="!loading && list.length === 0" description="暂无公告" />
      </van-list>
    </van-pull-refresh>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getAnnouncementList } from '@/api'

const router = useRouter()

const list = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const page = ref(1)
const pageSize = 10

const onClickLeft = () => {
  router.back()
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now - date
  const dayDiff = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (dayDiff === 0) {
    return `${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
  } else if (dayDiff === 1) {
    return '昨天'
  } else if (dayDiff < 7) {
    return `${dayDiff}天前`
  } else {
    return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`
  }
}

const stripHtml = (html) => {
  if (!html) return ''
  const doc = new DOMParser().parseFromString(html, 'text/html')
  const text = doc.body.textContent || ''
  return text.length > 80 ? text.substring(0, 80) + '...' : text
}

const goDetail = (id) => {
  router.push(`/announcements/${id}`)
}

const loadData = async () => {
  try {
    const res = await getAnnouncementList({ page: page.value, page_size: pageSize })
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
    console.error('Load announcements failed:', error)
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
.announcement-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.announcement-item {
  background-color: #fff;
  margin: 12px;
  padding: 16px;
  border-radius: 8px;
  
  .item-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 12px;
    
    .title {
      font-size: 16px;
      font-weight: bold;
      color: #333;
      flex: 1;
      margin-right: 8px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
  
  .item-content {
    display: flex;
    
    .content-text {
      flex: 1;
      margin-left: 12px;
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      
      &.no-cover {
        margin-left: 0;
      }
      
      .desc {
        font-size: 14px;
        color: #666;
        line-height: 1.6;
        margin: 0;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
      }
      
      .item-footer {
        display: flex;
        justify-content: space-between;
        font-size: 12px;
        color: #999;
        margin-top: 8px;
        
        .views {
          display: flex;
          align-items: center;
          
          .van-icon {
            margin-right: 4px;
          }
        }
      }
    }
  }
}
</style>
