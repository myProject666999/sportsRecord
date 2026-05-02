package controllers

import (
	"sportsRecord/database"
	"sportsRecord/models"
	"sportsRecord/utils"

	"github.com/gin-gonic/gin"
)

func GetTypeList(c *gin.Context) {
	typeParam := c.Query("type")
	
	var types []models.Type
	query := database.DB.Where("status = ?", 1)
	
	if typeParam != "" {
		query = query.Where("type = ?", typeParam)
	}
	
	if result := query.Order("sort desc, created_at desc").Find(&types); result.Error != nil {
		utils.ServerError(c, "获取类型列表失败")
		return
	}

	utils.Success(c, types)
}

func GetAdminTypeList(c *gin.Context) {
	var types []models.Type
	if result := database.DB.Order("sort desc, created_at desc").Find(&types); result.Error != nil {
		utils.ServerError(c, "获取类型列表失败")
		return
	}

	utils.Success(c, types)
}

func CreateType(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required"`
		Icon   string `json:"icon"`
		Sort   int    `json:"sort"`
		Type   int    `json:"type"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	_type := models.Type{
		Name:   req.Name,
		Icon:   req.Icon,
		Sort:   req.Sort,
		Type:   req.Type,
		Status: req.Status,
	}

	if result := database.DB.Create(&_type); result.Error != nil {
		utils.ServerError(c, "创建失败: "+result.Error.Error())
		return
	}

	utils.Success(c, _type)
}

func UpdateType(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name   string `json:"name"`
		Icon   string `json:"icon"`
		Sort   int    `json:"sort"`
		Type   int    `json:"type"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var _type models.Type
	if result := database.DB.First(&_type, id); result.Error != nil {
		utils.NotFound(c, "类型不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.Type != 0 {
		updates["type"] = req.Type
	}
	if req.Status != 0 || _type.Status != req.Status {
		updates["status"] = req.Status
	}

	if result := database.DB.Model(&_type).Updates(updates); result.Error != nil {
		utils.ServerError(c, "更新失败: "+result.Error.Error())
		return
	}

	utils.Success(c, _type)
}

func DeleteType(c *gin.Context) {
	id := c.Param("id")

	var _type models.Type
	if result := database.DB.First(&_type, id); result.Error != nil {
		utils.NotFound(c, "类型不存在")
		return
	}

	if result := database.DB.Delete(&_type); result.Error != nil {
		utils.ServerError(c, "删除失败: "+result.Error.Error())
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}
