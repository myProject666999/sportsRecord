package controllers

import (
	"sportsRecord/database"
	"sportsRecord/models"
	"sportsRecord/utils"

	"github.com/gin-gonic/gin"
)

func GetCarouselList(c *gin.Context) {
	var carousels []models.Carousel
	if result := database.DB.Where("status = ?", 1).Order("sort desc, created_at desc").Find(&carousels); result.Error != nil {
		utils.ServerError(c, "获取轮播图列表失败")
		return
	}

	utils.Success(c, carousels)
}

func GetAdminCarouselList(c *gin.Context) {
	var carousels []models.Carousel
	if result := database.DB.Order("sort desc, created_at desc").Find(&carousels); result.Error != nil {
		utils.ServerError(c, "获取轮播图列表失败")
		return
	}

	utils.Success(c, carousels)
}

func CreateCarousel(c *gin.Context) {
	var req struct {
		Title  string `json:"title"`
		Image  string `json:"image" binding:"required"`
		Link   string `json:"link"`
		Sort   int    `json:"sort"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	carousel := models.Carousel{
		Title:  req.Title,
		Image:  req.Image,
		Link:   req.Link,
		Sort:   req.Sort,
		Status: req.Status,
	}

	if result := database.DB.Create(&carousel); result.Error != nil {
		utils.ServerError(c, "创建失败: "+result.Error.Error())
		return
	}

	utils.Success(c, carousel)
}

func UpdateCarousel(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Title  string `json:"title"`
		Image  string `json:"image"`
		Link   string `json:"link"`
		Sort   int    `json:"sort"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var carousel models.Carousel
	if result := database.DB.First(&carousel, id); result.Error != nil {
		utils.NotFound(c, "轮播图不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Image != "" {
		updates["image"] = req.Image
	}
	if req.Link != "" {
		updates["link"] = req.Link
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.Status != 0 || carousel.Status != req.Status {
		updates["status"] = req.Status
	}

	if result := database.DB.Model(&carousel).Updates(updates); result.Error != nil {
		utils.ServerError(c, "更新失败: "+result.Error.Error())
		return
	}

	utils.Success(c, carousel)
}

func DeleteCarousel(c *gin.Context) {
	id := c.Param("id")

	var carousel models.Carousel
	if result := database.DB.First(&carousel, id); result.Error != nil {
		utils.NotFound(c, "轮播图不存在")
		return
	}

	if result := database.DB.Delete(&carousel); result.Error != nil {
		utils.ServerError(c, "删除失败: "+result.Error.Error())
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}
