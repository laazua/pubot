package dao

import (
	"pubot/internal/model"

	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (userDao *UserDao) Create(modelUser model.User) error {
	return userDao.db.Create(&modelUser).Error
}

func (userDao *UserDao) Delete(modelUser model.User) error {
	return userDao.db.Delete(&modelUser).Error
}

func (userDao *UserDao) Update(modelUser model.User) error {
	return userDao.db.Model(&model.User{}).Where("id = ?", modelUser.ID).Updates(modelUser).Error
}

func (userDao *UserDao) Get(modelUser model.User) (model.User, error) {

	var user model.User
	if err := userDao.db.Find(&user, modelUser.ID).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (userDao *UserDao) List() ([]model.User, error) {
	var users []model.User
	err := userDao.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
