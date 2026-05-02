package models

import (
	"time"

	"gorm.io/gorm"
)

type ForumComment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	PostID    uint           `json:"post_id" gorm:"index;not null"`
	UserID    uint           `json:"user_id" gorm:"index;not null"`
	User      *User          `json:"user" gorm:"foreignKey:UserID"`
	ParentID  uint           `json:"parent_id" gorm:"index;default:0"` // 回复的评论ID
	Content   string         `json:"content" gorm:"type:text;not null"`
	LikeCount int            `json:"like_count" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"` // 1:正常, 0:已删除
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (ForumComment) TableName() string {
	return "forum_comments"
}
