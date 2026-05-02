package controllers

import (
	"sportsRecord/database"
	"sportsRecord/models"
	"sportsRecord/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUserStepRecords(c *gin.Context) {
	userID := c.GetUint("user_id")
	page := 1
	pageSize := 30

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
	query := database.DB.Model(&models.StepRecord{}).Where("user_id = ?", userID)
	query.Count(&total)

	var records []models.StepRecord
	offset := (page - 1) * pageSize
	if result := query.Offset(offset).Limit(pageSize).Order("record_date desc").Find(&records); result.Error != nil {
		utils.ServerError(c, "获取步数记录失败")
		return
	}

	utils.Success(c, gin.H{
		"list":      records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateStepRecord(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Steps    int     `json:"steps" binding:"required"`
		Distance float64 `json:"distance"`
		Calories float64 `json:"calories"`
		Duration int     `json:"duration"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	today := time.Now().Format("2006-01-02")

	var existingRecord models.StepRecord
	if result := database.DB.Where("user_id = ? AND DATE(record_date) = ?", userID, today).First(&existingRecord); result.Error == nil {
		existingRecord.Steps += req.Steps
		existingRecord.Distance += req.Distance
		existingRecord.Calories += req.Calories
		existingRecord.Duration += req.Duration

		if result := database.DB.Save(&existingRecord); result.Error != nil {
			utils.ServerError(c, "更新步数记录失败: "+result.Error.Error())
			return
		}

		utils.Success(c, existingRecord)
		return
	}

	record := models.StepRecord{
		UserID:     userID,
		Steps:      req.Steps,
		Distance:   req.Distance,
		Calories:   req.Calories,
		Duration:   req.Duration,
		RecordDate: time.Now(),
	}

	if result := database.DB.Create(&record); result.Error != nil {
		utils.ServerError(c, "创建步数记录失败: "+result.Error.Error())
		return
	}

	utils.Success(c, record)
}

func GetTodayStepRecord(c *gin.Context) {
	userID := c.GetUint("user_id")
	today := time.Now().Format("2006-01-02")

	var record models.StepRecord
	if result := database.DB.Where("user_id = ? AND DATE(record_date) = ?", userID, today).First(&record); result.Error != nil {
		utils.Success(c, gin.H{
			"steps":    0,
			"distance": 0,
			"calories": 0,
			"duration": 0,
		})
		return
	}

	utils.Success(c, record)
}

func GetStepStats(c *gin.Context) {
	userID := c.GetUint("user_id")
	period := c.Query("period") // week, month

	var startTime time.Time
	now := time.Now()

	switch period {
	case "week":
		startTime = now.AddDate(0, 0, -7)
	case "month":
		startTime = now.AddDate(0, -1, 0)
	default:
		startTime = now.AddDate(0, 0, -7)
	}

	var totalSteps int64
	var totalDistance float64
	var totalCalories float64
	var totalDuration int64
	var recordCount int64

	database.DB.Model(&models.StepRecord{}).
		Where("user_id = ? AND record_date >= ?", userID, startTime).
		Select("COALESCE(SUM(steps), 0)").Scan(&totalSteps)

	database.DB.Model(&models.StepRecord{}).
		Where("user_id = ? AND record_date >= ?", userID, startTime).
		Select("COALESCE(SUM(distance), 0)").Scan(&totalDistance)

	database.DB.Model(&models.StepRecord{}).
		Where("user_id = ? AND record_date >= ?", userID, startTime).
		Select("COALESCE(SUM(calories), 0)").Scan(&totalCalories)

	database.DB.Model(&models.StepRecord{}).
		Where("user_id = ? AND record_date >= ?", userID, startTime).
		Select("COALESCE(SUM(duration), 0)").Scan(&totalDuration)

	database.DB.Model(&models.StepRecord{}).
		Where("user_id = ? AND record_date >= ?", userID, startTime).
		Count(&recordCount)

	utils.Success(c, gin.H{
		"total_steps":    totalSteps,
		"total_distance": totalDistance,
		"total_calories": totalCalories,
		"total_duration": totalDuration,
		"record_count":   recordCount,
		"period":         period,
	})
}

func GetAdminStepRecords(c *gin.Context) {
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

	if ps := c.Query("page_size"); ps != "" {
		pageSize = 0
		for _, ch := range ps {
			if ch >= '0' && ch <= '9' {
				pageSize = pageSize*10 + int(ch-'0')
			}
		}
	}

	userID := c.Query("user_id")

	var total int64
	query := database.DB.Model(&models.StepRecord{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	query.Count(&total)

	var records []models.StepRecord
	offset := (page - 1) * pageSize
	if result := query.Offset(offset).Limit(pageSize).Preload("User").Order("created_at desc").Find(&records); result.Error != nil {
		utils.ServerError(c, "获取步数记录失败")
		return
	}

	utils.Success(c, gin.H{
		"list":      records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetStepStatistics(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	type DailyStats struct {
		Date         string  `json:"date"`
		TotalSteps   int64   `json:"total_steps"`
		TotalDistance float64 `json:"total_distance"`
		TotalCalories float64 `json:"total_calories"`
		UserCount    int64   `json:"user_count"`
	}

	var stats []DailyStats

	database.DB.Model(&models.StepRecord{}).
		Select("DATE(record_date) as date, COALESCE(SUM(steps), 0) as total_steps, COALESCE(SUM(distance), 0) as total_distance, COALESCE(SUM(calories), 0) as total_calories, COUNT(DISTINCT user_id) as user_count").
		Where("DATE(record_date) BETWEEN ? AND ?", startDate, endDate).
		Group("DATE(record_date)").
		Order("date desc").
		Scan(&stats)

	var overallTotalSteps int64
	var overallTotalDistance float64
	var overallTotalCalories float64
	var overallUserCount int64

	database.DB.Model(&models.StepRecord{}).
		Where("DATE(record_date) BETWEEN ? AND ?", startDate, endDate).
		Select("COALESCE(SUM(steps), 0)").Scan(&overallTotalSteps)
	database.DB.Model(&models.StepRecord{}).
		Where("DATE(record_date) BETWEEN ? AND ?", startDate, endDate).
		Select("COALESCE(SUM(distance), 0)").Scan(&overallTotalDistance)
	database.DB.Model(&models.StepRecord{}).
		Where("DATE(record_date) BETWEEN ? AND ?", startDate, endDate).
		Select("COALESCE(SUM(calories), 0)").Scan(&overallTotalCalories)
	database.DB.Model(&models.StepRecord{}).
		Where("DATE(record_date) BETWEEN ? AND ?", startDate, endDate).
		Distinct("user_id").Count(&overallUserCount)

	utils.Success(c, gin.H{
		"daily_stats": stats,
		"overview": gin.H{
			"total_steps":    overallTotalSteps,
			"total_distance": overallTotalDistance,
			"total_calories": overallTotalCalories,
			"user_count":     overallUserCount,
		},
		"start_date": startDate,
		"end_date":   endDate,
	})
}
