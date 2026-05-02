<template>
  <div class="user-manage-page">
    <van-nav-bar title="用户管理" left-arrow @click-left="onClickLeft" fixed placeholder>
      <template #right>
        <van-icon name="search" size="20" @click="showSearch = !showSearch" />
      </template>
    </van-nav-bar>
    
    <van-search
      v-model="searchKeyword"
      placeholder="搜索用户名/昵称"
      show-action
      @search="onSearch"
      @cancel="onCancel"
      v-if="showSearch"
    />
    
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <van-cell-group inset>
          <div class="user-item" v-for="user in list" :key="user.id">
            <van-avatar :src="user.avatar || defaultAvatar" size="44" />
            <div class="user-info">
              <div class="user-name">
                <span class="username">{{ user.nickname || user.username }}</span>
                <van-tag 
                  :type="user.status === 1 ? 'success' : 'danger'" 
                  size="mini"
                >
                  {{ user.status === 1 ? '正常' : '禁用' }}
                </van-tag>
              </div>
              <div class="user-meta">
                <span>ID: {{ user.id }}</span>
                <span>{{ formatDate(user.created_at) }}</span>
              </div>
            </div>
            <van-dropdown-menu>
              <van-dropdown-item :options="actionOptions" @change="(value) => onAction(value, user)" />
            </van-dropdown-menu>
          </div>
        </van-cell-group>
        
        <van-empty v-if="!loading && list.length === 0" description="暂无用户数据" />
      </van-list>
    </van-pull-refresh>
    
    <van-dialog
      v-model:show="showStatusDialog"
      title="确认操作"
      :message="statusDialogMessage"
      show-cancel-button
      @confirm="confirmStatusChange"
    />
    
    <van-tabbar v-model="activeTab" route placeholder>
      <van-tabbar-item replace to="/admin/dashboard" icon="home-o">首页</van-tabbar-item>
      <van-tabbar-item replace to="/admin/users" icon="user-o">用户</van-tabbar-item>
      <van-tabbar-item replace to="/admin/forum" icon="chat-o">论坛</van-tabbar-item>
      <van-tabbar-item replace to="/admin/step-record" icon="chart-trending-o">统计</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getUserList, updateUserStatus } from '@/api'
import { showToast, showConfirmDialog } from 'vant'

const router = useRouter()

const activeTab = ref(1)
const showSearch = ref(false)
const searchKeyword = ref('')
const list = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const page = ref(1)
const pageSize = 10

const showStatusDialog = ref(false)
const statusDialogMessage = ref('')
const pendingUser = ref(null)
const pendingStatus = ref(1)

const actionOptions = ref([
  { text: '禁用用户', value: 'disable' },
  { text: '启用用户', value: 'enable' }
])

const defaultAvatar = computed(() => {
  return 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=user%20avatar%20default%20icon%20simple&image_size=square'
})

const onClickLeft = () => {
  router.back()
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`
}

const loadData = async () => {
  try {
    const params = {
      page: page.value,
      page_size: pageSize
    }
    
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }
    
    const res = await getUserList(params)
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
    console.error('Load users failed:', error)
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

const onSearch = () => {
  page.value = 1
  finished.value = false
  list.value = []
  loadData()
}

const onCancel = () => {
  searchKeyword.value = ''
  showSearch.value = false
  page.value = 1
  finished.value = false
  list.value = []
  loadData()
}

const onAction = (value, user) => {
  pendingUser.value = user
  
  if (value === 'disable' && user.status === 1) {
    pendingStatus.value = 0
    statusDialogMessage.value = `确定要禁用用户 "${user.nickname || user.username}" 吗？`
    showStatusDialog.value = true
  } else if (value === 'enable' && user.status === 0) {
    pendingStatus.value = 1
    statusDialogMessage.value = `确定要启用用户 "${user.nickname || user.username}" 吗？`
    showStatusDialog.value = true
  } else {
    showToast('用户状态已符合要求')
  }
}

const confirmStatusChange = async () => {
  if (!pendingUser.value) return
  
  try {
    await updateUserStatus(pendingUser.value.id, {
      status: pendingStatus.value
    })
    
    pendingUser.value.status = pendingStatus.value
    showToast(pendingStatus.value === 1 ? '用户已启用' : '用户已禁用')
  } catch (error) {
    console.error('Update user status failed:', error)
    showToast('操作失败')
  } finally {
    showStatusDialog.value = false
    pendingUser.value = null
  }
}

onMounted(() => {
  loadData()
})
</script>

<style lang="less" scoped>
.user-manage-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 50px;
}

.user-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #f5f5f5;
  
  &:last-child {
    border-bottom: none;
  }
  
  .user-info {
    flex: 1;
    margin-left: 12px;
    
    .user-name {
      display: flex;
      align-items: center;
      gap: 8px;
      
      .username {
        font-size: 15px;
        font-weight: 500;
        color: #333;
      }
    }
    
    .user-meta {
      display: flex;
      gap: 16px;
      margin-top: 4px;
      font-size: 12px;
      color: #999;
    }
  }
}

:deep(.van-dropdown-menu) {
  background-color: transparent;
  
  .van-dropdown-menu__item {
    justify-content: flex-end;
    padding: 8px 0;
    
    .van-icon {
      color: #999;
    }
  }
}
</style>
