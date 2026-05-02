package models

import (
	"time"

	"gorm.io/gorm"
)

type ForumPost struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"index;not null"`
	User        *User          `json:"user" gorm:"foreignKey:UserID"`
	Title       string         `json:"title" gorm:"size:200;not null"`
	Content     string         `json:"content" gorm:"type:text;not null"`
	Images      string         `json:"images" gorm:"size:1000"` // 多张图片用逗号分隔
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	LikeCount   int            `json:"like_count" gorm:"default:0"`
	CommentCount int           `json:"comment_count" gorm:"default:0"`
	Status      int            `json:"status" gorm:"default:1"` // 1:正常, 0:已删除, 2:审核中
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (ForumPost) TableName() string {
	return "forum_posts"
}
