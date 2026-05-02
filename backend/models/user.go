package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password     string         `json:"-" gorm:"size:255;not null"`
	Email        string         `json:"email" gorm:"size:100"`
	Phone        string         `json:"phone" gorm:"size:20"`
	Avatar       string         `json:"avatar" gorm:"size:255"`
	Nickname     string         `json:"nickname" gorm:"size:50"`
	Gender       int            `json:"gender" gorm:"default:0"` // 0:未知, 1:男, 2:女
	Birthday     *time.Time     `json:"birthday"`
	Height       float64        `json:"height" gorm:"type:decimal(5,2)"`
	Weight       float64        `json:"weight" gorm:"type:decimal(5,2)"`
	Status       int            `json:"status" gorm:"default:1"` // 1:正常, 0:禁用
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
