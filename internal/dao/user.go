package dao

import (
	"pubot/internal/model"

	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{db: db}
}

func (user *User) Create(modelUser model.User) error {
	return user.db.Create(&modelUser).Error
}

func (user *User) Delete(modelUser model.User) error {
	return user.db.Delete(&modelUser).Error
}

