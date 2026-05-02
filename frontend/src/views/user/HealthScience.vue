<template>
  <div class="health-science-page">
    <van-nav-bar title="健康科普" fixed placeholder />
    
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
            <div class="article-item" v-for="item in list" :key="item.id" @click="goDetail(item.id)">
              <div class="article-content">
                <h3 class="title">{{ item.title }}</h3>
                <p class="desc">{{ item.content?.substring(0, 60) || '' }}...</p>
                <div class="meta">
                  <span class="type">{{ item.type?.name || '科普' }}</span>
                  <span class="time">{{ formatDate(item.created_at) }}</span>
                  <span class="views"><van-icon name="eye-o" /> {{ item.view_count }}</span>
                </div>
              </div>
              <van-image
                :src="item.cover || defaultCover"
                width="100"
                height="70"
                fit="cover"
                round
              />
            </div>
            
            <van-empty v-if="!loading && list.length === 0" description="暂无科普文章" />
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
import { getHealthScienceList, getTypeList } from '@/api'

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
  return 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=health%20science%20article%20cover%20image%20green%20theme&image_size=square_hd'
})

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`
}

const goDetail = (id) => {
  router.push(`/health-science/${id}`)
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
    const res = await getHealthScienceList(params)
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
    console.error('Load health science failed:', error)
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
.health-science-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 50px;
}

.article-item {
  display: flex;
  align-items: flex-start;
  background-color: #fff;
  margin: 12px;
  padding: 12px;
  border-radius: 8px;
  
  .article-content {
    flex: 1;
    margin-right: 12px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    min-height: 70px;
    
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
      margin: 0 0 8px 0;
      line-height: 1.4;
      display: -webkit-box;
      -webkit-line-clamp: 1;
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
      
      .time, .views {
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
