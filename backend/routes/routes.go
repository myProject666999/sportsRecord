package routes

import (
	"sportsRecord/controllers"
	"sportsRecord/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/user/register", controllers.UserRegister)
		api.POST("/user/login", controllers.UserLogin)

		api.POST("/admin/login", controllers.AdminLogin)

		api.GET("/carousels", controllers.GetCarouselList)
		api.GET("/announcements", controllers.GetAnnouncementList)
		api.GET("/announcements/:id", controllers.GetAnnouncementDetail)
		api.GET("/types", controllers.GetTypeList)
		api.GET("/health-sciences", controllers.GetHealthScienceList)
		api.GET("/health-sciences/:id", controllers.GetHealthScienceDetail)
		api.GET("/teaching-videos", controllers.GetTeachingVideoList)
		api.GET("/teaching-videos/:id", controllers.GetTeachingVideoDetail)
		api.GET("/forum/posts", controllers.GetForumPostList)
		api.GET("/forum/posts/:id", controllers.GetForumPostDetail)
		api.GET("/forum/posts/:id/comments", controllers.GetForumComments)

		user := api.Group("/user").Use(middleware.AuthMiddleware(), middleware.UserMiddleware())
		{
			user.GET("/info", controllers.GetUserInfo)
			user.PUT("/info", controllers.UpdateUserInfo)
			user.PUT("/password", controllers.ChangeUserPassword)

			user.GET("/step-records", controllers.GetUserStepRecords)
			user.POST("/step-records", controllers.CreateStepRecord)
			user.GET("/step-records/today", controllers.GetTodayStepRecord)
			user.GET("/step-stats", controllers.GetStepStats)

			user.GET("/collections", controllers.GetUserCollections)
			user.POST("/collections/toggle", controllers.ToggleCollection)
			user.GET("/collections/status", controllers.CheckCollectionStatus)

			user.POST("/forum/posts", controllers.CreateForumPost)
			user.DELETE("/forum/posts/:id", controllers.DeleteForumPost)
			user.POST("/forum/posts/:id/comments", controllers.CreateForumComment)
		}

		admin := api.Group("/admin").Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			admin.GET("/info", controllers.GetAdminInfo)
			admin.GET("/dashboard", controllers.GetDashboardStats)

			admin.GET("/users", controllers.GetUserList)
			admin.PUT("/users/:id/status", controllers.UpdateUserStatus)

			admin.GET("/carousels", controllers.GetAdminCarouselList)
			admin.POST("/carousels", controllers.CreateCarousel)
			admin.PUT("/carousels/:id", controllers.UpdateCarousel)
			admin.DELETE("/carousels/:id", controllers.DeleteCarousel)

			admin.GET("/announcements", controllers.GetAdminAnnouncementList)
			admin.POST("/announcements", controllers.CreateAnnouncement)
			admin.PUT("/announcements/:id", controllers.UpdateAnnouncement)
			admin.DELETE("/announcements/:id", controllers.DeleteAnnouncement)

			admin.GET("/types", controllers.GetAdminTypeList)
			admin.POST("/types", controllers.CreateType)
			admin.PUT("/types/:id", controllers.UpdateType)
			admin.DELETE("/types/:id", controllers.DeleteType)

			admin.GET("/health-sciences", controllers.GetAdminHealthScienceList)
			admin.POST("/health-sciences", controllers.CreateHealthScience)
			admin.PUT("/health-sciences/:id", controllers.UpdateHealthScience)
			admin.DELETE("/health-sciences/:id", controllers.DeleteHealthScience)

			admin.GET("/teaching-videos", controllers.GetAdminTeachingVideoList)
			admin.POST("/teaching-videos", controllers.CreateTeachingVideo)
			admin.PUT("/teaching-videos/:id", controllers.UpdateTeachingVideo)
			admin.DELETE("/teaching-videos/:id", controllers.DeleteTeachingVideo)

			admin.GET("/forum/posts", controllers.GetAdminForumPostList)
			admin.PUT("/forum/posts/:id/status", controllers.UpdateForumPostStatus)
			admin.GET("/forum/comments", controllers.GetAdminForumComments)
			admin.DELETE("/forum/comments/:id", controllers.DeleteForumComment)

			admin.GET("/step-records", controllers.GetAdminStepRecords)
			admin.GET("/step-statistics", controllers.GetStepStatistics)
		}
	}
}
