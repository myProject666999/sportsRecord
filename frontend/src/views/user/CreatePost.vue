<template>
  <div class="create-post-page">
    <van-nav-bar title="发布帖子" left-arrow @click-left="onClickLeft" fixed placeholder>
      <template #right>
        <van-button type="primary" size="small" :loading="submitting" @click="onSubmit">
          发布
        </van-button>
      </template>
    </van-nav-bar>
    
    <van-form ref="formRef" @submit="onSubmit">
      <van-cell-group inset class="form-group">
        <van-field
          v-model="form.title"
          name="title"
          label="标题"
          placeholder="请输入帖子标题（1-50字）"
          :rules="[{ required: true, message: '请输入标题' }, { max: 50, message: '标题不能超过50字' }]"
        />
        
        <van-field
          v-model="form.content"
          name="content"
          label="内容"
          type="textarea"
          placeholder="请输入帖子内容（1-500字）"
          :rows="6"
          :rules="[{ required: true, message: '请输入内容' }, { max: 500, message: '内容不能超过500字' }]"
        />
        
        <van-field
          v-model="form.images"
          name="images"
          label="图片"
          placeholder="可上传多张图片（用逗号分隔URL）"
        />
      </van-cell-group>
      
      <div class="tips">
        <p>温馨提示：</p>
        <ul>
          <li>请遵守社区规范，发布文明健康的内容</li>
          <li>图片请使用网络图片URL，多张图片用英文逗号分隔</li>
        </ul>
      </div>
    </van-form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { createForumPost } from '@/api'
import { showToast } from 'vant'

const router = useRouter()

const formRef = ref(null)
const submitting = ref(false)
const form = ref({
  title: '',
  content: '',
  images: ''
})

const onClickLeft = () => {
  router.back()
}

const onSubmit = async () => {
  if (!form.value.title.trim()) {
    showToast('请输入标题')
    return
  }
  
  if (!form.value.content.trim()) {
    showToast('请输入内容')
    return
  }
  
  submitting.value = true
  
  try {
    const postData = {
      title: form.value.title.trim(),
      content: form.value.content.trim()
    }
    
    if (form.value.images.trim()) {
      postData.images = form.value.images.trim()
    }
    
    await createForumPost(postData)
    showToast('发布成功')
    router.replace('/forum')
  } catch (error) {
    console.error('Create post failed:', error)
    showToast('发布失败，请重试')
  } finally {
    submitting.value = false
  }
}
</script>

<style lang="less" scoped>
.create-post-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 24px;
}

.form-group {
  margin-top: 12px;
}

.tips {
  margin: 24px 16px;
  padding: 16px;
  background-color: rgba(76, 175, 80, 0.05);
  border-radius: 8px;
  border: 1px solid rgba(76, 175, 80, 0.2);
  
  p {
    font-size: 14px;
    font-weight: 500;
    color: #4CAF50;
    margin: 0 0 8px 0;
  }
  
  ul {
    margin: 0;
    padding-left: 18px;
    
    li {
      font-size: 12px;
      color: #666;
      margin: 4px 0;
      line-height: 1.6;
    }
  }
}
</style>
