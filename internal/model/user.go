package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Name      string `gorm:"name"`
	Email     string `gorm:"emal"`
	Password  string `gorm:"password"`
	IsAdmin   bool   `gorm:"is_admin"`
}

// 设置表名
func (User) TableName() string {
	return "pb_user"
}
