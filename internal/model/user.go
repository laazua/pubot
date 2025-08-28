package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"name"`
	Email    string `gorm:"emal"`
	Password string `gorm:"password"`
	IsAdmin  bool   `gorm:"is_admin"`
}

// 设置表名
func (User) TableName() string {
	return "pb_user"
}
