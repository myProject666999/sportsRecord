package models

import (
	"time"

	"gorm.io/gorm"
)

type TeachingVideo struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:200;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Cover       string         `json:"cover" gorm:"size:255"`
	VideoUrl    string         `json:"video_url" gorm:"size:500;not null"`
	Duration    int            `json:"duration" gorm:"default:0"` // 视频时长(秒)
	TypeID      uint           `json:"type_id" gorm:"index"`
	Type        *Type          `json:"type" gorm:"foreignKey:TypeID"`
	Sort        int            `json:"sort" gorm:"default:0"`
	Status      int            `json:"status" gorm:"default:1"` // 1:显示, 0:隐藏
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	LikeCount   int            `json:"like_count" gorm:"default:0"`
	PublishTime *time.Time     `json:"publish_time"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (TeachingVideo) TableName() string {
	return "teaching_videos"
}
