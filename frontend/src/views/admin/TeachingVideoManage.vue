<template>
  <div class="manage-page">
    <van-nav-bar title="教学视频管理" left-arrow @click-left="onClickLeft" fixed placeholder>
      <template #right>
        <van-button size="small" type="primary" @click="onAdd">新增</van-button>
      </template>
    </van-nav-bar>
    
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <van-cell-group inset>
          <div class="video-item" v-for="item in list" :key="item.id">
            <div class="video-cover" @click="viewDetail(item)">
              <van-image
                :src="item.cover"
                width="100%"
                height="180"
                fit="cover"
              />
              <div class="play-mask">
                <van-icon name="play-circle" size="48" />
              </div>
            </div>
            <div class="item-info">
              <h4 class="title">{{ item.title }}</h4>
              <div class="meta">
                <span class="type">{{ item.type?.name || '未分类' }}</span>
                <span><van-icon name="eye-o" /> {{ item.view_count }}</span>
              </div>
            </div>
            <div class="item-actions">
              <van-button size="mini" type="primary" @click="onEdit(item)">编辑</van-button>
              <van-button size="mini" type="danger" @click="onDelete(item)">删除</van-button>
            </div>
          </div>
        </van-cell-group>
        <van-empty v-if="!loading && list.length === 0" description="暂无教学视频" />
      </van-list>
    </van-pull-refresh>
    
    <van-tabbar v-model="activeTab" route placeholder>
      <van-tabbar-item replace to="/admin/dashboard" icon="home-o">首页</van-tabbar-item>
      <van-tabbar-item replace to="/admin/users" icon="user-o">用户</van-tabbar-item>
      <van-tabbar-item replace to="/admin/forum" icon="chat-o">论坛</van-tabbar-item>
      <van-tabbar-item replace to="/admin/step-record" icon="chart-trending-o">统计</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getTeachingVideoList, deleteTeachingVideo } from '@/api'
import { showToast, showConfirmDialog } from 'vant'

const router = useRouter()

const activeTab = ref(1)
const list = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const page = ref(1)
const pageSize = 10

const onClickLeft = () => {
  router.back()
}

const viewDetail = (item) => {
  showToast('视频预览功能开发中')
}

const onAdd = () => {
  showToast('新增视频功能开发中')
}

const onEdit = (item) => {
  showToast('编辑视频功能开发中')
}

const onDelete = (item) => {
  showConfirmDialog({
    title: '提示',
    message: '确定要删除该视频吗？',
  })
    .then(async () => {
      try {
        await deleteTeachingVideo(item.id)
        list.value = list.value.filter(i => i.id !== item.id)
        showToast('删除成功')
      } catch (error) {
        console.error('Delete failed:', error)
        showToast('删除失败')
      }
    })
    .catch(() => {})
}

const loadData = async () => {
  try {
    const res = await getTeachingVideoList({ page: page.value, page_size: pageSize })
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
    console.error('Load data failed:', error)
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
.manage-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 50px;
}

.video-item {
  padding: 12px 16px;
  border-bottom: 1px solid #f5f5f5;
  
  &:last-child {
    border-bottom: none;
  }
  
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
  }
  
  .item-info {
    padding: 12px 0;
    
    .title {
      font-size: 15px;
      font-weight: 500;
      color: #333;
      margin: 0 0 8px 0;
    }
    
    .meta {
      display: flex;
      align-items: center;
      gap: 16px;
      font-size: 12px;
      color: #999;
      
      .type {
        background-color: rgba(76, 175, 80, 0.1);
        color: #4CAF50;
        padding: 2px 6px;
        border-radius: 4px;
      }
      
      span {
        display: flex;
        align-items: center;
        
        .van-icon {
          margin-right: 4px;
        }
      }
    }
  }
  
  .item-actions {
    display: flex;
    gap: 8px;
    justify-content: flex-end;
  }
}
</style>
