<template>
  <div class="teaching-video-page">
    <van-nav-bar title="教学视频" fixed placeholder />
    
    <van-tabs v-model:active="activeType" type="card" sticky>
      <van-tab
        v-for="type in typeList"
        :key="type.id"
        :title="type.name"
      >
        <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
          <van-list
            v-model:loading="loading"
            :finished="finished"
            finished-text="没有更多了"
            @load="onLoad"
          >
            <div class="video-item" v-for="item in list" :key="item.id" @click="goDetail(item.id)">
              <div class="video-cover">
                <van-image
                  :src="item.cover || defaultCover"
                  width="100%"
                  height="180"
                  fit="cover"
                />
                <div class="play-mask">
                  <van-icon name="play-circle" size="48" />
                </div>
                <div class="duration" v-if="item.duration">
                  {{ formatDuration(item.duration) }}
                </div>
              </div>
              <div class="video-info">
                <h3 class="title">{{ item.title }}</h3>
                <div class="meta">
                  <span class="type" v-if="item.type?.name">{{ item.type.name }}</span>
                  <span class="views"><van-icon name="eye-o" /> {{ item.view_count }}</span>
                  <span class="time">{{ formatDate(item.created_at) }}</span>
                </div>
              </div>
            </div>
            
            <van-empty v-if="!loading && list.length === 0" description="暂无教学视频" />
          </van-list>
        </van-pull-refresh>
      </van-tab>
    </van-tabs>
    
    <van-tabbar v-model="activeTab" route placeholder>
      <van-tabbar-item replace to="/home" icon="home-o">首页</van-tabbar-item>
      <van-tabbar-item replace to="/forum" icon="chat-o">论坛</van-tabbar-item>
      <van-tabbar-item replace to="/health-science" icon="medicines-o">科普</van-tabbar-item>
      <van-tabbar-item replace to="/profile" icon="user-o">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getTeachingVideoList, getTypeList } from '@/api'

const router = useRouter()

const activeTab = ref(2)
const activeType = ref(0)
const typeList = ref([{ id: 0, name: '全部' }])
const list = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const page = ref(1)
const pageSize = 10

const defaultCover = computed(() => {
  return 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=exercise%20workout%20video%20thumbnail%20fitness&image_size=landscape_16_9'
})

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`
}

const formatDuration = (seconds) => {
  if (!seconds) return '00:00'
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`
}

const goDetail = (id) => {
  router.push(`/teaching-video/${id}`)
}

const loadTypes = async () => {
  try {
    const res = await getTypeList({ page: 1, page_size: 100 })
    const types = res.data?.list || []
    typeList.value = [{ id: 0, name: '全部' }, ...types]
  } catch (error) {
    console.error('Load types failed:', error)
  }
}

const loadData = async () => {
  const currentTypeId = typeList.value[activeType.value]?.id || 0
  const params = {
    page: page.value,
    page_size: pageSize
  }
  
  if (currentTypeId > 0) {
    params.type_id = currentTypeId
  }
  
  try {
    const res = await getTeachingVideoList(params)
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
    console.error('Load teaching videos failed:', error)
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

watch(activeType, () => {
  page.value = 1
  finished.value = false
  list.value = []
  loadData()
})

onMounted(() => {
  loadTypes()
  loadData()
})
</script>

<style lang="less" scoped>
.teaching-video-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 50px;
}

.video-item {
  background-color: #fff;
  margin: 12px;
  border-radius: 8px;
  overflow: hidden;
  
  .video-cover {
    position: relative;
    width: 100%;
    height: 180px;
    
    .play-mask {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      display: flex;
      justify-content: center;
      align-items: center;
      background-color: rgba(0, 0, 0, 0.2);
      color: rgba(255, 255, 255, 0.9);
    }
    
    .duration {
      position: absolute;
      right: 8px;
      bottom: 8px;
      background-color: rgba(0, 0, 0, 0.7);
      color: #fff;
      font-size: 12px;
      padding: 2px 6px;
      border-radius: 4px;
    }
  }
  
  .video-info {
    padding: 12px;
    
    .title {
      font-size: 15px;
      font-weight: 500;
      color: #333;
      margin: 0 0 8px 0;
      line-height: 1.4;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
    
    .meta {
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      gap: 8px;
      
      .type {
        font-size: 11px;
        color: #4CAF50;
        background-color: rgba(76, 175, 80, 0.1);
        padding: 2px 6px;
        border-radius: 4px;
      }
      
      .views, .time {
        font-size: 11px;
        color: #999;
        display: flex;
        align-items: center;
        
        .van-icon {
          margin-right: 2px;
        }
      }
    }
  }
}

:deep(.van-tabs__wrap) {
  background-color: #fff;
}
</style>
