package database

import (
	"log"
	"sportsRecord/models"
)

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Type{},
		&models.Announcement{},
		&models.Carousel{},
		&models.HealthScience{},
		&models.TeachingVideo{},
		&models.ForumPost{},
		&models.ForumComment{},
		&models.StepRecord{},
		&models.Collection{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")
}
