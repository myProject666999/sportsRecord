<template>
  <div class="manage-page">
    <van-nav-bar title="健康科普管理" left-arrow @click-left="onClickLeft" fixed placeholder>
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
          <div class="list-item" v-for="item in list" :key="item.id">
            <div class="item-content" @click="viewDetail(item)">
              <van-image
                :src="item.cover"
                width="100"
                height="70"
                fit="cover"
                round
              />
              <div class="info">
                <h4 class="title">{{ item.title }}</h4>
                <p class="desc">{{ item.content?.substring(0, 40) || '' }}...</p>
                <div class="meta">
                  <span class="type">{{ item.type?.name || '未分类' }}</span>
                  <span><van-icon name="eye-o" /> {{ item.view_count }}</span>
                </div>
              </div>
            </div>
            <div class="item-actions">
              <van-button size="mini" type="primary" @click="onEdit(item)">编辑</van-button>
              <van-button size="mini" type="danger" @click="onDelete(item)">删除</van-button>
            </div>
          </div>
        </van-cell-group>
        <van-empty v-if="!loading && list.length === 0" description="暂无科普文章" />
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
import { getHealthScienceList, deleteHealthScience } from '@/api'
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
  showToast('文章预览功能开发中')
}

const onAdd = () => {
  showToast('新增科普功能开发中')
}

const onEdit = (item) => {
  showToast('编辑科普功能开发中')
}

const onDelete = (item) => {
  showConfirmDialog({
    title: '提示',
    message: '确定要删除该文章吗？',
  })
    .then(async () => {
      try {
        await deleteHealthScience(item.id)
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
    const res = await getHealthScienceList({ page: page.value, page_size: pageSize })
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

.list-item {
  padding: 12px 16px;
  border-bottom: 1px solid #f5f5f5;
  
  &:last-child {
    border-bottom: none;
  }
  
  .item-content {
    display: flex;
    
    .info {
      flex: 1;
      margin-left: 12px;
      
      .title {
        font-size: 15px;
        font-weight: 500;
        color: #333;
        margin: 0 0 6px 0;
      }
      
      .desc {
        font-size: 13px;
        color: #666;
        margin: 0 0 8px 0;
        line-height: 1.5;
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
  }
  
  .item-actions {
    display: flex;
    gap: 8px;
    margin-top: 12px;
    justify-content: flex-end;
  }
}
</style>
