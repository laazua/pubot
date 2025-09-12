package dao

import (
	"errors"
	"pubot/internal/model"

	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (td *UserDao) Create(dbUser *model.PbUser) error {
	var pipelineExists model.PbUser
	if td.db.Where("name = ?", dbUser.Name).First(&pipelineExists).Error == nil {
		return errors.New("任务已经存在")
	}
	return td.db.Create(&dbUser).Error
}

func (td *UserDao) Delete(id uint) error {
	return td.db.Where("id = ?", id).Delete(&model.PbUser{}).Error
}

func (td *UserDao) GetByID(id uint) (*model.PbUser, error) {
	var modelUser model.PbUser
	err := td.db.First(&modelUser, id).Error
	if err != nil {
		return nil, err
	}
	return &modelUser, nil
}

func (td *UserDao) Update(dbUser *model.PbUser) error {
	result := td.db.Save(dbUser)
	return result.Error
}

func (td *UserDao) GetAllUsers() ([]model.PbUser, error) {
	var modelUsers []model.PbUser
	err := td.db.Find(&modelUsers).Error
	if err != nil {
		return nil, err
	}
	return modelUsers, nil
}

func (td *UserDao) Save(dbUser *model.PbUser) error {
	return td.db.Save(&dbUser).Error
}

func (td *UserDao) Auth(dbUser *model.PbUser) (*model.PbUser, error) {
	var user model.PbUser
	if err := td.db.Where("name = ?", dbUser.Name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
