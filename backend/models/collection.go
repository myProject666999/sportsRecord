package models

import (
	"time"

	"gorm.io/gorm"
)

type Collection struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"index;not null"`
	User       *User          `json:"user" gorm:"foreignKey:UserID"`
	ResourceID uint           `json:"resource_id" gorm:"index;not null"`
	ResourceType int           `json:"resource_type" gorm:"index;not null"` // 1:健康科普, 2:教学视频, 3:论坛帖子
	Title      string         `json:"title" gorm:"size:200"`
	Cover      string         `json:"cover" gorm:"size:255"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Collection) TableName() string {
	return "collections"
}
