package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	Nickname  string         `json:"nickname" gorm:"size:50"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	Role      int            `json:"role" gorm:"default:1"` // 1:普通管理员, 2:超级管理员
	Status    int            `json:"status" gorm:"default:1"` // 1:正常, 0:禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Admin) TableName() string {
	return "admins"
}
