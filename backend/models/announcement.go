package models

import (
	"time"

	"gorm.io/gorm"
)

type Announcement struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:200;not null"`
	Content     string         `json:"content" gorm:"type:text;not null"`
	Cover       string         `json:"cover" gorm:"size:255"`
	Sort        int            `json:"sort" gorm:"default:0"`
	Status      int            `json:"status" gorm:"default:1"` // 1:显示, 0:隐藏
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	PublishTime *time.Time     `json:"publish_time"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Announcement) TableName() string {
	return "announcements"
}
