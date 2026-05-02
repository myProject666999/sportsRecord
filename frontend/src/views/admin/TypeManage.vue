<template>
  <div class="manage-page">
    <van-nav-bar title="分类管理" left-arrow @click-left="onClickLeft" fixed placeholder>
      <template #right>
        <van-button size="small" type="primary" @click="showAddPopup = true">
          新增
        </van-button>
      </template>
    </van-nav-bar>
    
    <van-search
      v-model="searchKeyword"
      placeholder="搜索分类名称"
      show-action
      @search="onSearch"
      @cancel="onCancel"
    />
    
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <van-cell-group inset>
          <div class="list-item" v-for="item in list" :key="item.id">
            <div class="item-content">
              <span class="item-name">{{ item.name }}</span>
              <span class="item-sort">排序: {{ item.sort || 0 }}</span>
            </div>
            <div class="item-actions">
              <van-button size="mini" type="primary" @click="onEdit(item)">编辑</van-button>
              <van-button size="mini" type="danger" @click="onDelete(item)">删除</van-button>
            </div>
          </div>
        </van-cell-group>
        
        <van-empty v-if="!loading && list.length === 0" description="暂无分类数据" />
      </van-list>
    </van-pull-refresh>
    
    <van-popup v-model:show="showAddPopup" round position="bottom" :style="{ height: '60%' }">
      <div class="popup-header">
        <span class="popup-title">{{ isEdit ? '编辑分类' : '新增分类' }}</span>
        <van-icon name="cross" size="20" @click="closePopup" />
      </div>
      <van-form @submit="onSubmit">
        <van-cell-group inset>
          <van-field
            v-model="form.name"
            name="name"
            label="分类名称"
            placeholder="请输入分类名称"
            :rules="[{ required: true, message: '请输入分类名称' }]"
          />
          <van-field
            v-model="form.description"
            name="description"
            label="描述"
            placeholder="请输入分类描述"
            type="textarea"
            :rows="2"
          />
          <van-field
            v-model="form.sort"
            name="sort"
            label="排序"
            placeholder="数字越小越靠前"
            type="number"
          />
        </van-cell-group>
        <div class="popup-footer">
          <van-button round block type="primary" native-type="submit" :loading="submitting">
            确定
          </van-button>
        </div>
      </van-form>
    </van-popup>
    
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
import { getTypeList, createType, updateType, deleteType } from '@/api'
import { showToast, showConfirmDialog } from 'vant'

const router = useRouter()

const activeTab = ref(1)
const searchKeyword = ref('')
const list = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const page = ref(1)
const pageSize = 20

const showAddPopup = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const form = ref({
  id: null,
  name: '',
  description: '',
  sort: 0
})

const onClickLeft = () => {
  router.back()
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
    
    const res = await getTypeList(params)
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
    console.error('Load types failed:', error)
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
  page.value = 1
  finished.value = false
  list.value = []
  loadData()
}

const resetForm = () => {
  form.value = {
    id: null,
    name: '',
    description: '',
    sort: 0
  }
  isEdit.value = false
}

const closePopup = () => {
  showAddPopup.value = false
  resetForm()
}

const onEdit = (item) => {
  isEdit.value = true
  form.value = {
    id: item.id,
    name: item.name,
    description: item.description || '',
    sort: item.sort || 0
  }
  showAddPopup.value = true
}

const onDelete = (item) => {
  showConfirmDialog({
    title: '提示',
    message: `确定要删除分类 "${item.name}" 吗？`,
  })
    .then(async () => {
      try {
        await deleteType(item.id)
        showToast('删除成功')
        page.value = 1
        finished.value = false
        list.value = []
        loadData()
      } catch (error) {
        console.error('Delete type failed:', error)
        showToast('删除失败')
      }
    })
    .catch(() => {
      // on cancel
    })
}

const onSubmit = async () => {
  if (!form.value.name.trim()) {
    showToast('请输入分类名称')
    return
  }
  
  submitting.value = true
  
  try {
    const data = {
      name: form.value.name.trim(),
      description: form.value.description,
      sort: parseInt(form.value.sort) || 0
    }
    
    if (isEdit.value && form.value.id) {
      await updateType(form.value.id, data)
      showToast('更新成功')
    } else {
      await createType(data)
      showToast('创建成功')
    }
    
    closePopup()
    page.value = 1
    finished.value = false
    list.value = []
    loadData()
  } catch (error) {
    console.error('Submit failed:', error)
    showToast('操作失败')
  } finally {
    submitting.value = false
  }
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
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #f5f5f5;
  
  &:last-child {
    border-bottom: none;
  }
  
  .item-content {
    display: flex;
    flex-direction: column;
    
    .item-name {
      font-size: 15px;
      font-weight: 500;
      color: #333;
    }
    
    .item-sort {
      font-size: 12px;
      color: #999;
      margin-top: 4px;
    }
  }
  
  .item-actions {
    display: flex;
    gap: 8px;
  }
}

.popup-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #f5f5f5;
  
  .popup-title {
    font-size: 16px;
    font-weight: 500;
    color: #333;
  }
}

.popup-footer {
  padding: 16px;
}
</style>
