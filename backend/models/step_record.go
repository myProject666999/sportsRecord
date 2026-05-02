package models

import (
	"time"

	"gorm.io/gorm"
)

type StepRecord struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"index;not null"`
	User       *User          `json:"user" gorm:"foreignKey:UserID"`
	Steps      int            `json:"steps" gorm:"not null"`
	Distance   float64        `json:"distance" gorm:"type:decimal(10,2)"` // 距离(公里)
	Calories   float64        `json:"calories" gorm:"type:decimal(10,2)"` // 消耗卡路里
	Duration   int            `json:"duration" gorm:"default:0"`            // 时长(分钟)
	RecordDate time.Time      `json:"record_date" gorm:"type:date;index"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (StepRecord) TableName() string {
	return "step_records"
}
