<template>
  <div class="collection-page">
    <van-nav-bar title="我的收藏" left-arrow @click-left="onClickLeft" fixed placeholder />
    
    <van-tabs v-model:active="activeTab" type="card">
      <van-tab title="科普文章">
        <van-pull-refresh v-model="refreshing" @refresh="onRefresh('health')">
          <van-list
            v-model:loading="loadingHealth"
            :finished="finishedHealth"
            finished-text="没有更多了"
            @load="onLoad('health')"
          >
            <div class="collection-item" v-for="item in healthList" :key="item.id">
              <div class="item-content" @click="goHealthDetail(item.relation_id)">
                <van-image
                  :src="item.cover || defaultCover"
                  width="120"
                  height="80"
                  fit="cover"
                  round
                />
                <div class="item-info">
                  <h3 class="title">{{ item.title }}</h3>
                  <p class="desc">{{ item.content?.substring(0, 50) || '' }}...</p>
                </div>
              </div>
              <van-icon name="delete" class="delete-icon" @click="handleRemoveCollection(item.id)" />
            </div>
            <van-empty v-if="!loadingHealth && healthList.length === 0" description="暂无收藏" />
          </van-list>
        </van-pull-refresh>
      </van-tab>
      
      <van-tab title="教学视频">
        <van-pull-refresh v-model="refreshing" @refresh="onRefresh('video')">
          <van-list
            v-model:loading="loadingVideo"
            :finished="finishedVideo"
            finished-text="没有更多了"
            @load="onLoad('video')"
          >
            <div class="collection-item" v-for="item in videoList" :key="item.id">
              <div class="item-content" @click="goVideoDetail(item.relation_id)">
                <div class="video-cover">
                  <van-image
                    :src="item.cover || defaultCover"
                    width="120"
                    height="80"
                    fit="cover"
                    round
                  />
                  <van-icon name="play-circle" class="play-icon" size="32" />
                </div>
                <div class="item-info">
                  <h3 class="title">{{ item.title }}</h3>
                  <p class="desc">{{ item.content?.substring(0, 50) || '' }}...</p>
                </div>
              </div>
              <van-icon name="delete" class="delete-icon" @click="handleRemoveCollection(item.id)" />
            </div>
            <van-empty v-if="!loadingVideo && videoList.length === 0" description="暂无收藏" />
          </van-list>
        </van-pull-refresh>
      </van-tab>
    </van-tabs>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { getCollectionList, removeCollection as apiRemoveCollection } from '@/api'
import { showConfirmDialog, showToast } from 'vant'

const router = useRouter()

const activeTab = ref(0)
const refreshing = ref(false)

const healthList = ref([])
const loadingHealth = ref(false)
const finishedHealth = ref(false)
const healthPage = ref(1)

const videoList = ref([])
const loadingVideo = ref(false)
const finishedVideo = ref(false)
const videoPage = ref(1)

const pageSize = 10

const defaultCover = computed(() => {
  return 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=health%20science%20article%20cover%20image&image_size=square_hd'
})

const onClickLeft = () => {
  router.back()
}

const goHealthDetail = (id) => {
  router.push(`/health-science/${id}`)
}

const goVideoDetail = (id) => {
  router.push(`/teaching-video/${id}`)
}

const loadData = async (type) => {
  const collectionType = type === 'health' ? 1 : 2
  const page = type === 'health' ? healthPage : videoPage
  
  try {
    const res = await getCollectionList({ 
      page: page.value, 
      page_size: pageSize,
      type: collectionType
    })
    const data = res.data?.list || []
    
    if (type === 'health') {
      if (refreshing.value) {
        healthList.value = []
      }
      healthList.value.push(...data)
      if (data.length < pageSize) {
        finishedHealth.value = true
      } else {
        healthPage.value++
      }
      loadingHealth.value = false
    } else {
      if (refreshing.value) {
        videoList.value = []
      }
      videoList.value.push(...data)
      if (data.length < pageSize) {
        finishedVideo.value = true
      } else {
        videoPage.value++
      }
      loadingVideo.value = false
    }
  } catch (error) {
    console.error('Load collections failed:', error)
    if (type === 'health') loadingHealth.value = false
    else loadingVideo.value = false
  } finally {
    refreshing.value = false
  }
}

const onLoad = (type) => {
  if (type === 'health') {
    loadingHealth.value = true
  } else {
    loadingVideo.value = true
  }
  loadData(type)
}

const onRefresh = (type) => {
  if (type === 'health') {
    healthPage.value = 1
    finishedHealth.value = false
  } else {
    videoPage.value = 1
    finishedVideo.value = false
  }
  loadData(type)
}

const handleRemoveCollection = async (id) => {
  showConfirmDialog({
    title: '提示',
    message: '确定要取消收藏吗？',
  })
    .then(async () => {
      try {
        await apiRemoveCollection(id)
        showToast('已取消收藏')
        
        healthList.value = healthList.value.filter(item => item.id !== id)
        videoList.value = videoList.value.filter(item => item.id !== id)
      } catch (error) {
        console.error('Remove collection failed:', error)
        showToast('操作失败')
      }
    })
    .catch(() => {
      // on cancel
    })
}
</script>

<style lang="less" scoped>
.collection-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.collection-item {
  display: flex;
  align-items: flex-start;
  background-color: #fff;
  margin: 12px;
  padding: 12px;
  border-radius: 8px;
  
  .item-content {
    flex: 1;
    display: flex;
    align-items: flex-start;
  }
  
  .video-cover {
    position: relative;
    
    .play-icon {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      color: rgba(255, 255, 255, 0.9);
    }
  }
  
  .item-info {
    flex: 1;
    margin-left: 12px;
    display: flex;
    flex-direction: column;
    min-height: 80px;
    
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
    
    .desc {
      font-size: 12px;
      color: #999;
      margin: 0;
      line-height: 1.4;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
  }
  
  .delete-icon {
    color: #999;
    font-size: 18px;
    padding: 8px;
    margin-left: 8px;
  }
}

:deep(.van-tabs__wrap) {
  background-color: #fff;
}
</style>
