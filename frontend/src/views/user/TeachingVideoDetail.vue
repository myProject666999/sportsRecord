<template>
  <div class="detail-page">
    <van-nav-bar title="视频详情" left-arrow @click-left="onClickLeft" fixed placeholder>
      <template #right>
        <van-icon :name="isCollected ? 'star' : 'star-o'" size="20" @click="toggleCollection" />
      </template>
    </van-nav-bar>
    
    <div class="detail-content" v-if="detail">
      <div class="video-player" v-if="detail.video_url">
        <video
          ref="videoRef"
          :src="detail.video_url"
          poster="detail.cover"
          controls
          class="video-element"
        >
          您的浏览器不支持视频播放
        </video>
      </div>
      
      <div class="video-cover-only" v-else>
        <van-image
          :src="detail.cover || defaultCover"
          fit="cover"
          class="cover-image"
        />
        <div class="play-hint">
          <van-icon name="play-circle" size="48" />
          <span>视频链接暂未配置</span>
        </div>
      </div>
      
      <div class="video-info">
        <h1 class="title">{{ detail.title }}</h1>
        <div class="meta">
          <span class="type" v-if="detail.type?.name">{{ detail.type.name }}</span>
          <span class="duration" v-if="detail.duration">
            <van-icon name="clock-o" /> {{ formatDuration(detail.duration) }}
          </span>
          <span class="views"><van-icon name="eye-o" /> {{ detail.view_count }}</span>
        </div>
        <div class="author" v-if="detail.author">
          <van-icon name="manager-o" /> 教练：{{ detail.author }}
        </div>
      </div>
      
      <div class="video-description" v-if="detail.content">
        <h3 class="section-title">视频介绍</h3>
        <div class="desc-content" v-html="detail.content"></div>
      </div>
    </div>
    
    <van-loading v-if="loading" class="loading-center" />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getTeachingVideoDetail, addCollection, removeCollection, checkCollection } from '@/api'
import { showToast } from 'vant'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const videoRef = ref(null)
const detail = ref(null)
const loading = ref(true)
const isCollected = ref(false)

const collectionType = 2

const defaultCover = computed(() => {
  return 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=exercise%20workout%20video%20thumbnail%20fitness&image_size=landscape_16_9'
})

const onClickLeft = () => {
  router.back()
}

const formatDuration = (seconds) => {
  if (!seconds) return '00:00'
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`
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
    const res = await getTeachingVideoDetail(route.params.id)
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
}

.video-player {
  width: 100%;
  background-color: #000;
  
  .video-element {
    width: 100%;
    display: block;
  }
}

.video-cover-only {
  position: relative;
  width: 100%;
  height: 200px;
  
  .cover-image {
    width: 100%;
    height: 100%;
  }
  
  .play-hint {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    background-color: rgba(0, 0, 0, 0.4);
    color: rgba(255, 255, 255, 0.9);
    
    span {
      margin-top: 8px;
      font-size: 14px;
    }
  }
}

.video-info {
  padding: 16px;
  
  .title {
    font-size: 18px;
    font-weight: bold;
    color: #333;
    margin: 0 0 12px 0;
    line-height: 1.5;
  }
  
  .meta {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 12px;
    margin-bottom: 8px;
    
    .type {
      font-size: 12px;
      color: #4CAF50;
      background-color: rgba(76, 175, 80, 0.1);
      padding: 2px 8px;
      border-radius: 4px;
    }
    
    .duration, .views {
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

.video-description {
  padding: 16px;
  border-top: 8px solid #f5f5f5;
  
  .section-title {
    font-size: 16px;
    font-weight: bold;
    color: #333;
    margin: 0 0 12px 0;
  }
  
  .desc-content {
    font-size: 14px;
    line-height: 1.8;
    color: #666;
    
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
}

.loading-center {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px;
}
</style>
