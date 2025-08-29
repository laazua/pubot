package dao

import (
	"fmt"
	"pubot/internal/core/pwd"
	"pubot/internal/model"

	"gorm.io/gorm"
)

type AuthDao struct {
	db *gorm.DB
}

func NewAuthDao(db *gorm.DB) *AuthDao {
	return &AuthDao{db: db}
}

func (authDao *AuthDao) Sign(user model.User) error {
	var modelUser model.User
	err := authDao.db.Find(&modelUser).Where("name = ?", user.Name).Error
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
