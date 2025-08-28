package dao

import (
	"fmt"
	"pubot/internal/model"
	"pubot/internal/core/pwd"

	"gorm.io/gorm"
)

type Auth struct {
	db *gorm.DB
}

func NewAuth(db *gorm.DB) *Auth {
	return &Auth{db: db}
}

func (auth *Auth) Sign(user model.User) error {
	var modelUser model.User
	err := auth.db.First(&modelUser, user.Name).Error
	if err != nil {
		return err
	}
	ok, err := pwd.Verify(modelUser.Password, user.Password)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("用户认证失败")
	}
	return nil
}
