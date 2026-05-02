package controllers

import (
	"sportsRecord/database"
	"sportsRecord/models"
	"sportsRecord/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func GetTeachingVideoList(c *gin.Context) {
	page := 1
	pageSize := 10
	typeID := c.Query("type_id")

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
	query := database.DB.Model(&models.TeachingVideo{}).Where("status = ?", 1)
	if typeID != "" {
		query = query.Where("type_id = ?", typeID)
	}
	query.Count(&total)

	var videos []models.TeachingVideo
	offset := (page - 1) * pageSize
	if result := query.Offset(offset).Limit(pageSize).Preload("Type").Order("sort desc, created_at desc").Find(&videos); result.Error != nil {
		utils.ServerError(c, "获取教学视频列表失败")
		return
	}

	utils.Success(c, gin.H{
		"list":      videos,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetTeachingVideoDetail(c *gin.Context) {
	id := c.Param("id")

	var video models.TeachingVideo
	if result := database.DB.Preload("Type").First(&video, id); result.Error != nil {
		utils.NotFound(c, "教学视频不存在")
		return
	}

	if video.Status != 1 {
		utils.NotFound(c, "教学视频不存在")
		return
	}

	database.DB.Model(&video).Update("view_count", video.ViewCount+1)

	utils.Success(c, video)
}

func GetAdminTeachingVideoList(c *gin.Context) {
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
	database.DB.Model(&models.TeachingVideo{}).Count(&total)

	var videos []models.TeachingVideo
	offset := (page - 1) * pageSize
	if result := database.DB.Offset(offset).Limit(pageSize).Preload("Type").Order("sort desc, created_at desc").Find(&videos); result.Error != nil {
		utils.ServerError(c, "获取教学视频列表失败")
		return
	}

	utils.Success(c, gin.H{
		"list":      videos,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateTeachingVideo(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Cover       string `json:"cover"`
		VideoUrl    string `json:"video_url" binding:"required"`
		Duration    int    `json:"duration"`
		TypeID      uint   `json:"type_id"`
		Sort        int    `json:"sort"`
		Status      int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	now := time.Now()
	video := models.TeachingVideo{
		Title:       req.Title,
		Description: req.Description,
		Cover:       req.Cover,
		VideoUrl:    req.VideoUrl,
		Duration:    req.Duration,
		TypeID:      req.TypeID,
		Sort:        req.Sort,
		Status:      req.Status,
		PublishTime: &now,
	}

	if result := database.DB.Create(&video); result.Error != nil {
		utils.ServerError(c, "创建失败: "+result.Error.Error())
		return
	}

	utils.Success(c, video)
}

func UpdateTeachingVideo(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Cover       string `json:"cover"`
		VideoUrl    string `json:"video_url"`
		Duration    int    `json:"duration"`
		TypeID      uint   `json:"type_id"`
		Sort        int    `json:"sort"`
		Status      int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var video models.TeachingVideo
	if result := database.DB.First(&video, id); result.Error != nil {
		utils.NotFound(c, "教学视频不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Cover != "" {
		updates["cover"] = req.Cover
	}
	if req.VideoUrl != "" {
		updates["video_url"] = req.VideoUrl
	}
	if req.Duration != 0 {
		updates["duration"] = req.Duration
	}
	if req.TypeID > 0 {
		updates["type_id"] = req.TypeID
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.Status != 0 || video.Status != req.Status {
		updates["status"] = req.Status
	}

	if result := database.DB.Model(&video).Updates(updates); result.Error != nil {
		utils.ServerError(c, "更新失败: "+result.Error.Error())
		return
	}

	utils.Success(c, video)
}

func DeleteTeachingVideo(c *gin.Context) {
	id := c.Param("id")

	var video models.TeachingVideo
	if result := database.DB.First(&video, id); result.Error != nil {
		utils.NotFound(c, "教学视频不存在")
		return
	}

	if result := database.DB.Delete(&video); result.Error != nil {
		utils.ServerError(c, "删除失败: "+result.Error.Error())
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}
