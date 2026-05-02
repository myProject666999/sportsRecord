package controllers

import (
	"sportsRecord/database"
	"sportsRecord/models"
	"sportsRecord/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAnnouncementList(c *gin.Context) {
	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		page = 0
		for _, ch := range p {
			if ch >= '0' && ch <= '9' {
				page = page*10 + int(ch-'0')
			}
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		pageSize = 0
		for _, ch := range ps {
			if ch >= '0' && ch <= '9' {
				pageSize = pageSize*10 + int(ch-'0')
			}
		}
	}

	var total int64
	query := database.DB.Model(&models.Announcement{}).Where("status = ?", 1)
	query.Count(&total)

	var announcements []models.Announcement
	offset := (page - 1) * pageSize
	if result := query.Offset(offset).Limit(pageSize).Order("sort desc, created_at desc").Find(&announcements); result.Error != nil {
		utils.ServerError(c, "获取公告列表失败")
		return
	}

	utils.Success(c, gin.H{
		"list":      announcements,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetAnnouncementDetail(c *gin.Context) {
	id := c.Param("id")

	var announcement models.Announcement
	if result := database.DB.First(&announcement, id); result.Error != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	if announcement.Status != 1 {
		utils.NotFound(c, "公告不存在")
		return
	}

	database.DB.Model(&announcement).Update("view_count", announcement.ViewCount+1)

	utils.Success(c, announcement)
}

func GetAdminAnnouncementList(c *gin.Context) {
	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		page = 0
		for _, ch := range p {
			if ch >= '0' && ch <= '9' {
				page = page*10 + int(ch-'0')
			}
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		pageSize = 0
		for _, ch := range ps {
			if ch >= '0' && ch <= '9' {
				pageSize = pageSize*10 + int(ch-'0')
			}
		}
	}

	var total int64
	database.DB.Model(&models.Announcement{}).Count(&total)

	var announcements []models.Announcement
	offset := (page - 1) * pageSize
	if result := database.DB.Offset(offset).Limit(pageSize).Order("sort desc, created_at desc").Find(&announcements); result.Error != nil {
		utils.ServerError(c, "获取公告列表失败")
		return
	}

	utils.Success(c, gin.H{
		"list":      announcements,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateAnnouncement(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Cover   string `json:"cover"`
		Sort    int    `json:"sort"`
		Status  int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	now := time.Now()
	announcement := models.Announcement{
		Title:       req.Title,
		Content:     req.Content,
		Cover:       req.Cover,
		Sort:        req.Sort,
		Status:      req.Status,
		PublishTime: &now,
	}

	if result := database.DB.Create(&announcement); result.Error != nil {
		utils.ServerError(c, "创建失败: "+result.Error.Error())
		return
	}

	utils.Success(c, announcement)
}

func UpdateAnnouncement(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Cover   string `json:"cover"`
		Sort    int    `json:"sort"`
		Status  int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var announcement models.Announcement
	if result := database.DB.First(&announcement, id); result.Error != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Cover != "" {
		updates["cover"] = req.Cover
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.Status != 0 || announcement.Status != req.Status {
		updates["status"] = req.Status
	}

	if result := database.DB.Model(&announcement).Updates(updates); result.Error != nil {
		utils.ServerError(c, "更新失败: "+result.Error.Error())
		return
	}

	utils.Success(c, announcement)
}

func DeleteAnnouncement(c *gin.Context) {
	id := c.Param("id")

	var announcement models.Announcement
	if result := database.DB.First(&announcement, id); result.Error != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	if result := database.DB.Delete(&announcement); result.Error != nil {
		utils.ServerError(c, "删除失败: "+result.Error.Error())
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}
