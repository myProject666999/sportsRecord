<template>
  <div class="home-page">
    <van-nav-bar title="运动健康" fixed placeholder>
      <template #right>
        <van-icon name="search" size="20" @click="onSearch" />
      </template>
    </van-nav-bar>
    
    <div class="home-content">
      <van-swipe class="banner-swipe" :autoplay="3000" indicator-color="#fff" indicator-active-color="#4CAF50">
        <van-swipe-item v-for="item in carouselList" :key="item.id">
          <div class="banner-item">
            <img :src="item.image || defaultBanner" alt="banner" @click="onBannerClick(item)" />
          </div>
        </van-swipe-item>
      </van-swipe>
      
      <div class="quick-menu">
        <van-grid :column-num="4" :border="false">
          <van-grid-item icon="orders-o" text="步数记录" @click="router.push('/step-record')" />
          <van-grid-item icon="newspaper-o" text="公告信息" @click="router.push('/announcements')" />
          <van-grid-item icon="chat-o" text="论坛交流" @click="router.push('/forum')" />
          <van-grid-item icon="medicines-o" text="健康科普" @click="router.push('/health-science')" />
          <van-grid-item icon="video-o" text="教学视频" @click="router.push('/teaching-video')" />
          <van-grid-item icon="star-o" text="我的收藏" @click="router.push('/collection')" />
          <van-grid-item icon="user-o" text="个人中心" @click="router.push('/profile')" />
          <van-grid-item icon="setting-o" text="设置" @click="onSetting" />
        </van-grid>
      </div>
      
      <van-divider>最新公告</van-divider>
      
      <div class="announcement-list">
        <van-cell 
          v-for="item in announcementList" 
          :key="item.id"
          :title="item.title"
          :label="formatDate(item.created_at)"
          is-link
          @click="router.push(`/announcements/${item.id}`)"
        >
          <template #icon>
            <van-icon name="bell-o" color="#4CAF50" />
          </template>
        </van-cell>
      </div>
      
      <van-divider>今日运动</van-divider>
      
      <div class="today-steps">
        <van-cell-group inset>
          <van-cell title="今日步数" :value="todaySteps + ' 步'" center>
            <template #label>
              <span class="step-desc">距离: {{ todayDistance.toFixed(2) }} km | 消耗: {{ todayCalories.toFixed(0) }} 卡路里</span>
            </template>
            <template #right-icon>
              <van-button type="primary" size="small" @click="recordSteps">记录步数</van-button>
            </template>
          </van-cell>
        </van-cell-group>
      </div>
    </div>
    
    <van-tabbar v-model="activeTab" route placeholder>
      <van-tabbar-item replace to="/home" icon="home-o">首页</van-tabbar-item>
      <van-tabbar-item replace to="/forum" icon="chat-o">论坛</van-tabbar-item>
      <van-tabbar-item replace to="/health-science" icon="medicines-o">科普</van-tabbar-item>
      <van-tabbar-item replace to="/profile" icon="user-o">我的</van-tabbar-item>
    </van-tabbar>
    
    <van-popup v-model:show="showRecordPopup" round position="bottom">
      <div class="record-popup">
        <van-nav-bar title="记录步数" left-arrow @click-left="showRecordPopup = false" />
        <van-form @submit="onSubmitSteps">
          <van-cell-group inset>
            <van-field
              v-model="stepForm.steps"
              type="number"
              name="steps"
              label="步数"
              placeholder="请输入步数"
              :rules="[{ required: true, message: '请填写步数' }]"
            />
            <van-field
              v-model="stepForm.distance"
              type="number"
              name="distance"
              label="距离(公里)"
              placeholder="请输入距离（选填）"
            />
            <van-field
              v-model="stepForm.calories"
              type="number"
              name="calories"
              label="消耗卡路里"
              placeholder="请输入卡路里（选填）"
            />
            <van-field
              v-model="stepForm.duration"
              type="number"
              name="duration"
              label="运动时长(分钟)"
              placeholder="请输入时长（选填）"
            />
          </van-cell-group>
          <div style="margin: 16px;">
            <van-button round block type="primary" native-type="submit" :loading="recording">
              保存
            </van-button>
          </div>
        </van-form>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getCarouselList, getAnnouncementList, getTodayStepRecord, createStepRecord } from '@/api'
import { showToast } from 'vant'

const router = useRouter()
const userStore = useUserStore()

const activeTab = ref(0)
const carouselList = ref([])
const announcementList = ref([])
const todaySteps = ref(0)
const todayDistance = ref(0)
const todayCalories = ref(0)
const showRecordPopup = ref(false)
const recording = ref(false)

const defaultBanner = computed(() => {
  return 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fitness%20exercise%20banner%20sport%20healthy&image_size=landscape_16_9'
})

const stepForm = ref({
  steps: '',
  distance: '',
  calories: '',
  duration: ''
})

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
}

const loadCarousel = async () => {
  try {
    const res = await getCarouselList()
    carouselList.value = res.data || []
  } catch (error) {
    console.error('Load carousel failed:', error)
  }
}

const loadAnnouncements = async () => {
  try {
    const res = await getAnnouncementList({ page: 1, page_size: 5 })
    announcementList.value = res.data?.list || []
  } catch (error) {
    console.error('Load announcements failed:', error)
  }
}

const loadTodaySteps = async () => {
  if (!userStore.isLoggedIn) return
  
  try {
    const res = await getTodayStepRecord()
    const data = res.data
    todaySteps.value = data.steps || 0
    todayDistance.value = data.distance || 0
    todayCalories.value = data.calories || 0
  } catch (error) {
    console.error('Load today steps failed:', error)
  }
}

const onBannerClick = (item) => {
  if (item.link) {
    window.open(item.link, '_blank')
  }
}

const onSearch = () => {
  showToast('搜索功能开发中')
}

const onSetting = () => {
  showToast('设置功能开发中')
}

const recordSteps = () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  stepForm.value = { steps: '', distance: '', calories: '', duration: '' }
  showRecordPopup.value = true
}

const onSubmitSteps = async (values) => {
  recording.value = true
  try {
    await createStepRecord({
      steps: parseInt(values.steps) || 0,
      distance: parseFloat(values.distance) || 0,
      calories: parseFloat(values.calories) || 0,
      duration: parseInt(values.duration) || 0
    })
    showToast('记录成功')
    showRecordPopup.value = false
    loadTodaySteps()
  } catch (error) {
    console.error('Record steps failed:', error)
  } finally {
    recording.value = false
  }
}

onMounted(() => {
  loadCarousel()
  loadAnnouncements()
  if (userStore.isLoggedIn) {
    loadTodaySteps()
  }
})
</script>

<style lang="less" scoped>
.home-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 50px;
}

.home-content {
  padding-bottom: 60px;
}

.banner-swipe {
  height: 180px;
  
  .banner-item {
    height: 100%;
    
    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }
}

.quick-menu {
  background-color: #fff;
  padding: 16px 0;
  margin-bottom: 12px;
}

.announcement-list {
  background-color: #fff;
  margin: 0 12px;
  border-radius: 8px;
  overflow: hidden;
}

.today-steps {
  margin: 12px 0;
  
  .step-desc {
    font-size: 12px;
    color: #999;
  }
}

.record-popup {
  height: 60vh;
  overflow-y: auto;
  padding-bottom: 20px;
}
</style>
