package controllers

import (
	"sportsRecord/database"
	"sportsRecord/models"
	"sportsRecord/utils"

	"github.com/gin-gonic/gin"
)

func GetUserCollections(c *gin.Context) {
	userID := c.GetUint("user_id")
	page := 1
	pageSize := 20
	resourceType := c.Query("resource_type")

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
	query := database.DB.Model(&models.Collection{}).Where("user_id = ?", userID)
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	query.Count(&total)

	var collections []models.Collection
	offset := (page - 1) * pageSize
	if result := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&collections); result.Error != nil {
		utils.ServerError(c, "获取收藏列表失败")
		return
	}

	utils.Success(c, gin.H{
		"list":      collections,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func ToggleCollection(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		ResourceID   uint   `json:"resource_id" binding:"required"`
		ResourceType int    `json:"resource_type" binding:"required"`
		Title        string `json:"title"`
		Cover        string `json:"cover"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var existingCollection models.Collection
	result := database.DB.Where("user_id = ? AND resource_id = ? AND resource_type = ?", 
		userID, req.ResourceID, req.ResourceType).First(&existingCollection)

	if result.Error == nil {
		if delResult := database.DB.Delete(&existingCollection); delResult.Error != nil {
			utils.ServerError(c, "取消收藏失败: "+delResult.Error.Error())
			return
		}
		utils.Success(c, gin.H{
			"is_collected": false,
			"message":      "取消收藏成功",
		})
		return
	}

	collection := models.Collection{
		UserID:       userID,
		ResourceID:   req.ResourceID,
		ResourceType: req.ResourceType,
		Title:        req.Title,
		Cover:        req.Cover,
	}

	if createResult := database.DB.Create(&collection); createResult.Error != nil {
		utils.ServerError(c, "收藏失败: "+createResult.Error.Error())
		return
	}

	utils.Success(c, gin.H{
		"is_collected": true,
		"message":      "收藏成功",
	})
}

func CheckCollectionStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	resourceIDStr := c.Query("resource_id")
	resourceTypeStr := c.Query("resource_type")

	if resourceIDStr == "" || resourceTypeStr == "" {
		utils.BadRequest(c, "参数缺失")
		return
	}

	var resourceID uint
	for _, ch := range resourceIDStr {
		if ch >= '0' && ch <= '9' {
			resourceID = resourceID*10 + uint(ch-'0')
		}
	}

	var resourceType int
	for _, ch := range resourceTypeStr {
		if ch >= '0' && ch <= '9' {
			resourceType = resourceType*10 + int(ch-'0')
		}
	}

	var count int64
	database.DB.Model(&models.Collection{}).
		Where("user_id = ? AND resource_id = ? AND resource_type = ?", 
			userID, resourceID, resourceType).Count(&count)

	utils.Success(c, gin.H{
		"is_collected": count > 0,
	})
}
