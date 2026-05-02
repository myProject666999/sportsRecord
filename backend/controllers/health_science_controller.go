package controllers

import (
	"sportsRecord/database"
	"sportsRecord/models"
	"sportsRecord/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func GetHealthScienceList(c *gin.Context) {
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
	query := database.DB.Model(&models.HealthScience{}).Where("status = ?", 1)
	if typeID != "" {
		query = query.Where("type_id = ?", typeID)
	}
	query.Count(&total)

	var healthSciences []models.HealthScience
	offset := (page - 1) * pageSize
	if result := query.Offset(offset).Limit(pageSize).Preload("Type").Order("sort desc, created_at desc").Find(&healthSciences); result.Error != nil {
		utils.ServerError(c, "获取健康科普列表失败")
		return
	}

	utils.Success(c, gin.H{
		"list":      healthSciences,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetHealthScienceDetail(c *gin.Context) {
	id := c.Param("id")

	var healthScience models.HealthScience
	if result := database.DB.Preload("Type").First(&healthScience, id); result.Error != nil {
		utils.NotFound(c, "健康科普不存在")
		return
	}

	if healthScience.Status != 1 {
		utils.NotFound(c, "健康科普不存在")
		return
	}

	database.DB.Model(&healthScience).Update("view_count", healthScience.ViewCount+1)

	utils.Success(c, healthScience)
}

func GetAdminHealthScienceList(c *gin.Context) {
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
	database.DB.Model(&models.HealthScience{}).Count(&total)

	var healthSciences []models.HealthScience
	offset := (page - 1) * pageSize
	if result := database.DB.Offset(offset).Limit(pageSize).Preload("Type").Order("sort desc, created_at desc").Find(&healthSciences); result.Error != nil {
		utils.ServerError(c, "获取健康科普列表失败")
		return
	}

	utils.Success(c, gin.H{
		"list":      healthSciences,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateHealthScience(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Cover   string `json:"cover"`
		TypeID  uint   `json:"type_id"`
		Sort    int    `json:"sort"`
		Status  int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	now := time.Now()
	healthScience := models.HealthScience{
		Title:       req.Title,
		Content:     req.Content,
		Cover:       req.Cover,
		TypeID:      req.TypeID,
		Sort:        req.Sort,
		Status:      req.Status,
		PublishTime: &now,
	}

	if result := database.DB.Create(&healthScience); result.Error != nil {
		utils.ServerError(c, "创建失败: "+result.Error.Error())
		return
	}

	utils.Success(c, healthScience)
}

func UpdateHealthScience(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Cover   string `json:"cover"`
		TypeID  uint   `json:"type_id"`
		Sort    int    `json:"sort"`
		Status  int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var healthScience models.HealthScience
	if result := database.DB.First(&healthScience, id); result.Error != nil {
		utils.NotFound(c, "健康科普不存在")
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
	if req.TypeID > 0 {
		updates["type_id"] = req.TypeID
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.Status != 0 || healthScience.Status != req.Status {
		updates["status"] = req.Status
	}

	if result := database.DB.Model(&healthScience).Updates(updates); result.Error != nil {
		utils.ServerError(c, "更新失败: "+result.Error.Error())
		return
	}

	utils.Success(c, healthScience)
}

func DeleteHealthScience(c *gin.Context) {
	id := c.Param("id")

	var healthScience models.HealthScience
	if result := database.DB.First(&healthScience, id); result.Error != nil {
		utils.NotFound(c, "健康科普不存在")
		return
	}

	if result := database.DB.Delete(&healthScience); result.Error != nil {
		utils.ServerError(c, "删除失败: "+result.Error.Error())
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}
