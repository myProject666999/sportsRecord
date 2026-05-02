package models

import (
	"time"

	"gorm.io/gorm"
)

type Type struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:50;not null"`
	Icon      string         `json:"icon" gorm:"size:255"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"` // 1:启用, 0:禁用
	Type      int            `json:"type" gorm:"default:1"`    // 1:健康科普, 2:教学视频
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Type) TableName() string {
	return "types"
}
