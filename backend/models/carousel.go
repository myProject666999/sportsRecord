package models

import (
	"time"

	"gorm.io/gorm"
)

type Carousel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200"`
	Image     string         `json:"image" gorm:"size:255;not null"`
	Link      string         `json:"link" gorm:"size:255"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"` // 1:显示, 0:隐藏
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Carousel) TableName() string {
	return "carousels"
}
