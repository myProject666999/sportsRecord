package controllers

import (
	"sportsRecord/database"
	"sportsRecord/models"
	"sportsRecord/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Nickname string `json:"nickname"`
}

type UpdateUserRequest struct {
	Nickname string  `json:"nickname"`
	Phone    string  `json:"phone"`
	Email    string  `json:"email"`
	Avatar   string  `json:"avatar"`
	Gender   int     `json:"gender"`
	Height   float64 `json:"height"`
	Weight   float64 `json:"weight"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func UserRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var existingUser models.User
	if result := database.DB.Where("username = ?", req.Username).First(&existingUser); result.Error == nil {
		utils.BadRequest(c, "用户名已存在")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ServerError(c, "密码加密失败")
		return
	}

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Status:   1,
	}

	if result := database.DB.Create(&user); result.Error != nil {
		utils.ServerError(c, "注册失败: "+result.Error.Error())
		return
	}

	token, err := utils.GenerateToken(user.ID, 0)
	if err != nil {
		utils.ServerError(c, "生成Token失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func UserLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if result := database.DB.Where("username = ?", req.Username).First(&user); result.Error != nil {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	if user.Status != 1 {
		utils.Forbidden(c, "账号已被禁用")
		return
	}

	token, err := utils.GenerateToken(user.ID, 0)
	if err != nil {
		utils.ServerError(c, "生成Token失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func GetUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, user)
}

func UpdateUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Gender != 0 {
		updates["gender"] = req.Gender
	}
	if req.Height > 0 {
		updates["height"] = req.Height
	}
	if req.Weight > 0 {
		updates["weight"] = req.Weight
	}

	if result := database.DB.Model(&user).Updates(updates); result.Error != nil {
		utils.ServerError(c, "更新失败: "+result.Error.Error())
		return
	}

	database.DB.First(&user, userID)
	utils.Success(c, user)
}

func ChangeUserPassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if !utils.CheckPassword(req.OldPassword, user.Password) {
		utils.BadRequest(c, "原密码错误")
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.ServerError(c, "密码加密失败")
		return
	}

	if result := database.DB.Model(&user).Update("password", hashedPassword); result.Error != nil {
		utils.ServerError(c, "修改密码失败: "+result.Error.Error())
		return
	}

	utils.Success(c, gin.H{"message": "密码修改成功"})
}
