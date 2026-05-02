import request from '@/utils/request'

export function userLogin(data) {
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}

export function userRegister(data) {
  return request({
    url: '/user/register',
    method: 'post',
    data
  })
}

export function getUserInfo() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}

export function updateUserInfo(data) {
  return request({
    url: '/user/info',
    method: 'put',
    data
  })
}

export function changePassword(data) {
  return request({
    url: '/user/password',
    method: 'put',
    data
  })
}

export function adminLogin(data) {
  return request({
    url: '/admin/login',
    method: 'post',
    data
  })
}

export function getAdminInfo() {
  return request({
    url: '/admin/info',
    method: 'get'
  })
}

export function getDashboardStats() {
  return request({
    url: '/admin/dashboard',
    method: 'get'
  })
}

export function getAdminDashboardStats() {
  return request({
    url: '/admin/dashboard',
    method: 'get'
  })
}

export function getCarouselList() {
  return request({
    url: '/carousels',
    method: 'get',
    showLoading: false
  })
}

export function getAdminCarouselList() {
  return request({
    url: '/admin/carousels',
    method: 'get'
  })
}

export function createCarousel(data) {
  return request({
    url: '/admin/carousels',
    method: 'post',
    data
  })
}

export function updateCarousel(id, data) {
  return request({
    url: `/admin/carousels/${id}`,
    method: 'put',
    data
  })
}

export function deleteCarousel(id) {
  return request({
    url: `/admin/carousels/${id}`,
    method: 'delete'
  })
}

export function getAnnouncementList(params) {
  return request({
    url: '/announcements',
    method: 'get',
    params,
    showLoading: false
  })
}

export function getAnnouncementDetail(id) {
  return request({
    url: `/announcements/${id}`,
    method: 'get'
  })
}

export function getAdminAnnouncementList(params) {
  return request({
    url: '/admin/announcements',
    method: 'get',
    params
  })
}

export function createAnnouncement(data) {
  return request({
    url: '/admin/announcements',
    method: 'post',
    data
  })
}

export function updateAnnouncement(id, data) {
  return request({
    url: `/admin/announcements/${id}`,
    method: 'put',
    data
  })
}

export function deleteAnnouncement(id) {
  return request({
    url: `/admin/announcements/${id}`,
    method: 'delete'
  })
}

export function getTypeList(params) {
  return request({
    url: '/types',
    method: 'get',
    params,
    showLoading: false
  })
}

export function getAdminTypeList() {
  return request({
    url: '/admin/types',
    method: 'get'
  })
}

export function createType(data) {
  return request({
    url: '/admin/types',
    method: 'post',
    data
  })
}

export function updateType(id, data) {
  return request({
    url: `/admin/types/${id}`,
    method: 'put',
    data
  })
}

export function deleteType(id) {
  return request({
    url: `/admin/types/${id}`,
    method: 'delete'
  })
}

export function getHealthScienceList(params) {
  return request({
    url: '/health-sciences',
    method: 'get',
    params,
    showLoading: false
  })
}

export function getHealthScienceDetail(id) {
  return request({
    url: `/health-sciences/${id}`,
    method: 'get'
  })
}

export function getAdminHealthScienceList(params) {
  return request({
    url: '/admin/health-sciences',
    method: 'get',
    params
  })
}

export function createHealthScience(data) {
  return request({
    url: '/admin/health-sciences',
    method: 'post',
    data
  })
}

export function updateHealthScience(id, data) {
  return request({
    url: `/admin/health-sciences/${id}`,
    method: 'put',
    data
  })
}

export function deleteHealthScience(id) {
  return request({
    url: `/admin/health-sciences/${id}`,
    method: 'delete'
  })
}

export function getTeachingVideoList(params) {
  return request({
    url: '/teaching-videos',
    method: 'get',
    params,
    showLoading: false
  })
}

export function getTeachingVideoDetail(id) {
  return request({
    url: `/teaching-videos/${id}`,
    method: 'get'
  })
}

export function getAdminTeachingVideoList(params) {
  return request({
    url: '/admin/teaching-videos',
    method: 'get',
    params
  })
}

export function createTeachingVideo(data) {
  return request({
    url: '/admin/teaching-videos',
    method: 'post',
    data
  })
}

export function updateTeachingVideo(id, data) {
  return request({
    url: `/admin/teaching-videos/${id}`,
    method: 'put',
    data
  })
}

export function deleteTeachingVideo(id) {
  return request({
    url: `/admin/teaching-videos/${id}`,
    method: 'delete'
  })
}

export function getForumPostList(params) {
  return request({
    url: '/forum/posts',
    method: 'get',
    params,
    showLoading: false
  })
}

export function getForumPostDetail(id) {
  return request({
    url: `/forum/posts/${id}`,
    method: 'get'
  })
}

export function createForumPost(data) {
  return request({
    url: '/user/forum/posts',
    method: 'post',
    data
  })
}

export function deleteForumPost(id) {
  return request({
    url: `/user/forum/posts/${id}`,
    method: 'delete'
  })
}

export function getForumComments(id, params) {
  return request({
    url: `/forum/posts/${id}/comments`,
    method: 'get',
    params,
    showLoading: false
  })
}

export function getAdminForumPostList(params) {
  return request({
    url: '/admin/forum/posts',
    method: 'get',
    params
  })
}

export function updateForumPostStatus(id, data) {
  return request({
    url: `/admin/forum/posts/${id}/status`,
    method: 'put',
    data
  })
}

export function getAdminForumComments(params) {
  return request({
    url: '/admin/forum/comments',
    method: 'get',
    params
  })
}

export function deleteForumComment(id) {
  return request({
    url: `/admin/forum/comments/${id}`,
    method: 'delete'
  })
}

export function getUserStepRecords(params) {
  return request({
    url: '/user/step-records',
    method: 'get',
    params
  })
}

export function createStepRecord(data) {
  return request({
    url: '/user/step-records',
    method: 'post',
    data
  })
}

export function getTodayStepRecord() {
  return request({
    url: '/user/step-records/today',
    method: 'get'
  })
}

export function getStepStats(params) {
  return request({
    url: '/user/step-stats',
    method: 'get',
    params
  })
}

export function getAdminStepRecords(params) {
  return request({
    url: '/admin/step-records',
    method: 'get',
    params
  })
}

export function getStepStatistics(params) {
  return request({
    url: '/admin/step-statistics',
    method: 'get',
    params
  })
}

export function getUserCollections(params) {
  return request({
    url: '/user/collections',
    method: 'get',
    params
  })
}

export function toggleCollection(data) {
  return request({
    url: '/user/collections/toggle',
    method: 'post',
    data
  })
}

export function checkCollectionStatus(params) {
  return request({
    url: '/user/collections/status',
    method: 'get',
    params,
    showLoading: false
  })
}

export function getUserList(params) {
  return request({
    url: '/admin/users',
    method: 'get',
    params
  })
}

export function updateUserStatus(id, data) {
  return request({
    url: `/admin/users/${id}/status`,
    method: 'put',
    data
  })
}

export function getStepRecordList(params) {
  return request({
    url: '/user/step-records',
    method: 'get',
    params
  })
}

export function getStepRecordStats() {
  return request({
    url: '/user/step-stats',
    method: 'get'
  })
}

export function getCollectionList(params) {
  return request({
    url: '/user/collections',
    method: 'get',
    params
  })
}

export function addCollection(data) {
  return request({
    url: '/user/collections',
    method: 'post',
    data
  })
}

export function checkCollection(type, relationId) {
  return request({
    url: '/user/collections/status',
    method: 'get',
    params: { type, relation_id: relationId }
  })
}

export function removeCollection(id) {
  return request({
    url: `/user/collections/${id}`,
    method: 'delete'
  })
}

export function getForumCommentList(postId, params) {
  return request({
    url: `/forum/posts/${postId}/comments`,
    method: 'get',
    params
  })
}

export function createForumComment(data) {
  return request({
    url: `/user/forum/posts/${data.post_id}/comments`,
    method: 'post',
    data
  })
}
