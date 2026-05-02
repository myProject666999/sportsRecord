<template>
  <div class="step-record-page">
    <van-nav-bar title="步数记录" left-arrow @click-left="onClickLeft" fixed placeholder />
    
    <div class="today-stats">
      <div class="stats-card">
        <div class="stats-item">
          <span class="stats-num">{{ todayStats.steps || 0 }}</span>
          <span class="stats-label">今日步数</span>
        </div>
        <div class="stats-item">
          <span class="stats-num">{{ todayStats.calories || 0 }}</span>
          <span class="stats-label">消耗卡路里</span>
        </div>
        <div class="stats-item">
          <span class="stats-num">{{ todayStats.distance || 0 }}</span>
          <span class="stats-label">距离(km)</span>
        </div>
      </div>
    </div>
    
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <van-cell-group inset class="record-group">
          <div class="record-item" v-for="item in list" :key="item.id">
            <div class="record-date">
              <van-icon name="calendar" />
              <span>{{ item.record_date }}</span>
            </div>
            <div class="record-details">
              <div class="detail-item">
                <span class="detail-num">{{ item.steps }}</span>
                <span class="detail-label">步数</span>
              </div>
              <div class="detail-item">
                <span class="detail-num">{{ item.calories }}</span>
                <span class="detail-label">卡路里</span>
              </div>
              <div class="detail-item">
                <span class="detail-num">{{ item.distance }}</span>
                <span class="detail-label">km</span>
              </div>
            </div>
          </div>
        </van-cell-group>
        
        <van-empty v-if="!loading && list.length === 0" description="暂无步数记录" />
      </van-list>
    </van-pull-refresh>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getStepRecordList, getStepRecordStats } from '@/api'

const router = useRouter()

const list = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const page = ref(1)
const pageSize = 10
const todayStats = ref({ steps: 0, calories: 0, distance: 0 })

const onClickLeft = () => {
  router.back()
}

const loadStats = async () => {
  try {
    const res = await getStepRecordStats()
    todayStats.value = {
      steps: res.data?.today_steps || 0,
      calories: res.data?.today_calories || 0,
      distance: res.data?.today_distance || 0
    }
  } catch (error) {
    console.error('Load stats failed:', error)
  }
}

const loadData = async () => {
  try {
    const res = await getStepRecordList({ page: page.value, page_size: pageSize })
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
    console.error('Load step records failed:', error)
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
  loadStats()
  loadData()
}

onMounted(() => {
  loadStats()
  loadData()
})
</script>

<style lang="less" scoped>
.step-record-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.today-stats {
  padding: 16px;
  
  .stats-card {
    background: linear-gradient(135deg, #4CAF50 0%, #81C784 100%);
    border-radius: 12px;
    padding: 24px 16px;
    display: flex;
    justify-content: space-around;
    
    .stats-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      color: #fff;
      
      .stats-num {
        font-size: 24px;
        font-weight: bold;
        margin-bottom: 4px;
      }
      
      .stats-label {
        font-size: 12px;
        opacity: 0.9;
      }
    }
  }
}

.record-group {
  margin: 16px;
  
  .record-item {
    padding: 16px;
    border-bottom: 1px solid #f5f5f5;
    
    &:last-child {
      border-bottom: none;
    }
    
    .record-date {
      display: flex;
      align-items: center;
      font-size: 14px;
      color: #666;
      margin-bottom: 12px;
      
      .van-icon {
        margin-right: 6px;
      }
    }
    
    .record-details {
      display: flex;
      justify-content: space-around;
      
      .detail-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        
        .detail-num {
          font-size: 18px;
          font-weight: bold;
          color: #333;
          margin-bottom: 4px;
        }
        
        .detail-label {
          font-size: 12px;
          color: #999;
        }
      }
    }
  }
}
</style>
