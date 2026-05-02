import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/',
    redirect: '/home'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/user/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/user/Register.vue'),
    meta: { title: '注册' }
  },
  {
    path: '/home',
    name: 'Home',
    component: () => import('@/views/user/Home.vue'),
    meta: { title: '首页', requiresAuth: false }
  },
  {
    path: '/announcements',
    name: 'AnnouncementList',
    component: () => import('@/views/user/AnnouncementList.vue'),
    meta: { title: '公告信息', requiresAuth: false }
  },
  {
    path: '/announcements/:id',
    name: 'AnnouncementDetail',
    component: () => import('@/views/user/AnnouncementDetail.vue'),
    meta: { title: '公告详情', requiresAuth: false }
  },
  {
    path: '/forum',
    name: 'Forum',
    component: () => import('@/views/user/Forum.vue'),
    meta: { title: '论坛交流', requiresAuth: false }
  },
  {
    path: '/forum/:id',
    name: 'ForumPostDetail',
    component: () => import('@/views/user/ForumPostDetail.vue'),
    meta: { title: '帖子详情', requiresAuth: false }
  },
  {
    path: '/forum/create',
    name: 'CreatePost',
    component: () => import('@/views/user/CreatePost.vue'),
    meta: { title: '发布帖子', requiresAuth: true }
  },
  {
    path: '/health-science',
    name: 'HealthScience',
    component: () => import('@/views/user/HealthScience.vue'),
    meta: { title: '健康科普', requiresAuth: false }
  },
  {
    path: '/health-science/:id',
    name: 'HealthScienceDetail',
    component: () => import('@/views/user/HealthScienceDetail.vue'),
    meta: { title: '科普详情', requiresAuth: false }
  },
  {
    path: '/teaching-video',
    name: 'TeachingVideo',
    component: () => import('@/views/user/TeachingVideo.vue'),
    meta: { title: '教学视频', requiresAuth: false }
  },
  {
    path: '/teaching-video/:id',
    name: 'TeachingVideoDetail',
    component: () => import('@/views/user/TeachingVideoDetail.vue'),
    meta: { title: '视频详情', requiresAuth: false }
  },
  {
    path: '/step-record',
    name: 'StepRecord',
    component: () => import('@/views/user/StepRecord.vue'),
    meta: { title: '步数记录', requiresAuth: true }
  },
  {
    path: '/collection',
    name: 'Collection',
    component: () => import('@/views/user/Collection.vue'),
    meta: { title: '我的收藏', requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/user/Profile.vue'),
    meta: { title: '个人中心', requiresAuth: true }
  },
  {
    path: '/profile/edit',
    name: 'EditProfile',
    component: () => import('@/views/user/EditProfile.vue'),
    meta: { title: '编辑资料', requiresAuth: true }
  },
  {
    path: '/password/change',
    name: 'ChangePassword',
    component: () => import('@/views/user/ChangePassword.vue'),
    meta: { title: '修改密码', requiresAuth: true }
  },
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('@/views/admin/Login.vue'),
    meta: { title: '管理员登录', isAdmin: true }
  },
  {
    path: '/admin',
    redirect: '/admin/dashboard'
  },
  {
    path: '/admin/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/admin/Dashboard.vue'),
    meta: { title: '仪表盘', requiresAuth: true, isAdmin: true }
  },
  {
    path: '/admin/users',
    name: 'UserManage',
    component: () => import('@/views/admin/UserManage.vue'),
    meta: { title: '用户管理', requiresAuth: true, isAdmin: true }
  },
  {
    path: '/admin/forum',
    name: 'ForumManage',
    component: () => import('@/views/admin/ForumManage.vue'),
    meta: { title: '论坛管理', requiresAuth: true, isAdmin: true }
  },
  {
    path: '/admin/carousel',
    name: 'CarouselManage',
    component: () => import('@/views/admin/CarouselManage.vue'),
    meta: { title: '轮播图管理', requiresAuth: true, isAdmin: true }
  },
  {
    path: '/admin/announcement',
    name: 'AnnouncementManage',
    component: () => import('@/views/admin/AnnouncementManage.vue'),
    meta: { title: '公告管理', requiresAuth: true, isAdmin: true }
  },
  {
    path: '/admin/health-science',
    name: 'HealthScienceManage',
    component: () => import('@/views/admin/HealthScienceManage.vue'),
    meta: { title: '健康科普管理', requiresAuth: true, isAdmin: true }
  },
  {
    path: '/admin/step-record',
    name: 'StepRecordManage',
    component: () => import('@/views/admin/StepRecordManage.vue'),
    meta: { title: '步数记录管理', requiresAuth: true, isAdmin: true }
  },
  {
    path: '/admin/teaching-video',
    name: 'TeachingVideoManage',
    component: () => import('@/views/admin/TeachingVideoManage.vue'),
    meta: { title: '教学视频管理', requiresAuth: true, isAdmin: true }
  },
  {
    path: '/admin/type',
    name: 'TypeManage',
    component: () => import('@/views/admin/TypeManage.vue'),
    meta: { title: '类型管理', requiresAuth: true, isAdmin: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  document.title = to.meta.title || '运动健康'
  
  const userStore = useUserStore()
  
  if (to.meta.isAdmin) {
    if (to.path === '/admin/login') {
      if (userStore.isLoggedIn && userStore.isAdmin) {
        next('/admin/dashboard')
      } else {
        next()
      }
    } else {
      if (userStore.isLoggedIn && userStore.isAdmin) {
        next()
      } else {
        next('/admin/login')
      }
    }
  } else {
    if (to.meta.requiresAuth && !userStore.isLoggedIn) {
      next('/login')
    } else {
      if ((to.path === '/login' || to.path === '/register') && userStore.isLoggedIn && !userStore.isAdmin) {
        next('/home')
      } else {
        next()
      }
    }
  }
})

export default router
