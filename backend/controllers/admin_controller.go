package controllers

import (
	"sportsRecord/database"
	"sportsRecord/models"
	"sportsRecord/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func AdminLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var admin models.Admin
	if result := database.DB.Where("username = ?", req.Username).First(&admin); result.Error != nil {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	if !utils.CheckPassword(req.Password, admin.Password) {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	if admin.Status != 1 {
		utils.Forbidden(c, "账号已被禁用")
		return
	}

	token, err := utils.GenerateToken(admin.ID, 1)
	if err != nil {
		utils.ServerError(c, "生成Token失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"admin": admin,
	})
}

func GetAdminInfo(c *gin.Context) {
	adminID := c.GetUint("user_id")

	var admin models.Admin
	if result := database.DB.First(&admin, adminID); result.Error != nil {
		utils.NotFound(c, "管理员不存在")
		return
	}

	utils.Success(c, admin)
}

func CreateAdmin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname"`
		Role     int    `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var existingAdmin models.Admin
	if result := database.DB.Where("username = ?", req.Username).First(&existingAdmin); result.Error == nil {
		utils.BadRequest(c, "用户名已存在")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ServerError(c, "密码加密失败")
		return
	}

	admin := models.Admin{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Role:     req.Role,
		Status:   1,
	}

	if result := database.DB.Create(&admin); result.Error != nil {
		utils.ServerError(c, "创建失败: "+result.Error.Error())
		return
	}

	utils.Success(c, admin)
}

func GetUserList(c *gin.Context) {
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
	database.DB.Model(&models.User{}).Count(&total)

	var users []models.User
	offset := (page - 1) * pageSize
	if result := database.DB.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&users); result.Error != nil {
		utils.ServerError(c, "获取用户列表失败")
		return
	}

	utils.Success(c, gin.H{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if result := database.DB.First(&user, id); result.Error != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if result := database.DB.Model(&user).Update("status", req.Status); result.Error != nil {
		utils.ServerError(c, "更新失败: "+result.Error.Error())
		return
	}

	utils.Success(c, gin.H{"message": "状态更新成功"})
}

func GetDashboardStats(c *gin.Context) {
	var userCount int64
	var todayStepCount int64
	var postCount int64
	var videoCount int64

	today := time.Now().Format("2006-01-02")

	database.DB.Model(&models.User{}).Count(&userCount)
	database.DB.Model(&models.ForumPost{}).Count(&postCount)
	database.DB.Model(&models.TeachingVideo{}).Count(&videoCount)
	database.DB.Model(&models.StepRecord{}).Where("DATE(record_date) = ?", today).Count(&todayStepCount)

	var totalSteps int64
	database.DB.Model(&models.StepRecord{}).Select("COALESCE(SUM(steps), 0)").Scan(&totalSteps)

	utils.Success(c, gin.H{
		"user_count":        userCount,
		"today_step_count": todayStepCount,
		"total_steps":      totalSteps,
		"post_count":       postCount,
		"video_count":      videoCount,
	})
}
