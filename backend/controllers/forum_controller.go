package controllers

import (
	"sportsRecord/database"
	"sportsRecord/models"
	"sportsRecord/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetForumPostList(c *gin.Context) {
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
	query := database.DB.Model(&models.ForumPost{}).Where("status = ?", 1)
	query.Count(&total)

	var posts []models.ForumPost
	offset := (page - 1) * pageSize
	if result := query.Offset(offset).Limit(pageSize).Preload("User").Order("created_at desc").Find(&posts); result.Error != nil {
		utils.ServerError(c, "获取帖子列表失败")
		return
	}

	utils.Success(c, gin.H{
		"list":      posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetForumPostDetail(c *gin.Context) {
	id := c.Param("id")

	var post models.ForumPost
	if result := database.DB.Preload("User").First(&post, id); result.Error != nil {
		utils.NotFound(c, "帖子不存在")
		return
	}

	if post.Status != 1 {
		utils.NotFound(c, "帖子不存在")
		return
	}

	database.DB.Model(&post).Update("view_count", post.ViewCount+1)

	utils.Success(c, post)
}

func CreateForumPost(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Images  string `json:"images"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	post := models.ForumPost{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
		Images:  req.Images,
		Status:  1,
	}

	if result := database.DB.Create(&post); result.Error != nil {
		utils.ServerError(c, "创建失败: "+result.Error.Error())
		return
	}

	utils.Success(c, post)
}

func DeleteForumPost(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var post models.ForumPost
	if result := database.DB.First(&post, id); result.Error != nil {
		utils.NotFound(c, "帖子不存在")
		return
	}

	if post.UserID != userID {
		utils.Forbidden(c, "无权删除该帖子")
		return
	}

	database.DB.Model(&post).Update("status", 0)

	utils.Success(c, gin.H{"message": "删除成功"})
}

func GetForumComments(c *gin.Context) {
	postID := c.Param("id")
	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		page = 0
		for _, ch := range p {
			if ch >= '0' && ch <= '9' {
				page = page*10 + int(ch-'0')
			}
		}
	}

	var total int64
	query := database.DB.Model(&models.ForumComment{}).Where("post_id = ? AND status = ?", postID, 1)
	query.Count(&total)

	var comments []models.ForumComment
	offset := (page - 1) * pageSize
	if result := query.Offset(offset).Limit(pageSize).Preload("User").Order("created_at desc").Find(&comments); result.Error != nil {
		utils.ServerError(c, "获取评论列表失败")
		return
	}

	utils.Success(c, gin.H{
		"list":      comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateForumComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	postIDStr := c.Param("id")

	var req struct {
		Content  string `json:"content" binding:"required"`
		ParentID uint   `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var postID uint
	for _, ch := range postIDStr {
		if ch >= '0' && ch <= '9' {
			postID = postID*10 + uint(ch-'0')
		}
	}

	comment := models.ForumComment{
		PostID:   postID,
		UserID:   userID,
		ParentID: req.ParentID,
		Content:  req.Content,
		Status:   1,
	}

	if result := database.DB.Create(&comment); result.Error != nil {
		utils.ServerError(c, "创建失败: "+result.Error.Error())
		return
	}

	database.DB.Model(&models.ForumPost{}).Where("id = ?", postID).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))

	utils.Success(c, comment)
}

func GetAdminForumPostList(c *gin.Context) {
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
	database.DB.Model(&models.ForumPost{}).Count(&total)

	var posts []models.ForumPost
	offset := (page - 1) * pageSize
	if result := database.DB.Offset(offset).Limit(pageSize).Preload("User").Order("created_at desc").Find(&posts); result.Error != nil {
		utils.ServerError(c, "获取帖子列表失败")
		return
	}

	utils.Success(c, gin.H{
		"list":      posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func UpdateForumPostStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var post models.ForumPost
	if result := database.DB.First(&post, id); result.Error != nil {
		utils.NotFound(c, "帖子不存在")
		return
	}

	if result := database.DB.Model(&post).Update("status", req.Status); result.Error != nil {
		utils.ServerError(c, "更新失败: "+result.Error.Error())
		return
	}

	utils.Success(c, gin.H{"message": "状态更新成功"})
}

func GetAdminForumComments(c *gin.Context) {
	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		page = 0
		for _, ch := range p {
			if ch >= '0' && ch <= '9' {
				page = page*10 + int(ch-'0')
			}
		}
	}

	var total int64
	database.DB.Model(&models.ForumComment{}).Count(&total)

	var comments []models.ForumComment
	offset := (page - 1) * pageSize
	if result := database.DB.Offset(offset).Limit(pageSize).Preload("User").Order("created_at desc").Find(&comments); result.Error != nil {
		utils.ServerError(c, "获取评论列表失败")
		return
	}

	utils.Success(c, gin.H{
		"list":      comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func DeleteForumComment(c *gin.Context) {
	id := c.Param("id")

	var comment models.ForumComment
	if result := database.DB.First(&comment, id); result.Error != nil {
		utils.NotFound(c, "评论不存在")
		return
	}

	if result := database.DB.Delete(&comment); result.Error != nil {
		utils.ServerError(c, "删除失败: "+result.Error.Error())
		return
	}

	database.DB.Model(&models.ForumPost{}).Where("id = ?", comment.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))

	utils.Success(c, gin.H{"message": "删除成功"})
}
