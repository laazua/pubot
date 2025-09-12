package service

import (
	"errors"
	"fmt"
	"pubot/internal/dao"
	"pubot/internal/dto"
	"pubot/internal/model"
	"pubot/internal/utils"
)

type UserService struct {
	userDao *dao.UserDao
}

func NewUserService(userDao *dao.UserDao) *UserService {
	return &UserService{userDao: userDao}
}

func (us *UserService) Create(userDto dto.UserRequest) (*dto.UserRequest, error) {
	hashPwd, err := utils.Hash(userDto.Password)
	if err != nil {
		return nil, err
	}
	user := model.PbUser{
		Name:     userDto.Username,
		Password: hashPwd,
		Role:     userDto.Role,
	}
	if err := us.userDao.Create(&user); err != nil {
		return nil, err
	}
	userDto.Password = ""
	return &userDto, nil
}

func (us *UserService) Delete(id uint) error {
	// 先检查任务是否存在
	_, err := us.userDao.GetByID(id)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// 执行删除操作
	if err := us.userDao.Delete(id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	// 可以在这里添加一些额外的逻辑，比如：
	// - 记录删除日志
	// - 发送删除通知
	// - 清理相关资源等
	return nil
}

func (us *UserService) Update(id uint, userDto dto.UserRequest) (*model.PbUser, error) {
	// 1. 先查找现有任务
	existingUser, err := us.userDao.GetByID(id)
	if err != nil {
		return nil, err
	}
	// 3. 更新其他字段
	if userDto.Username != "" {
		existingUser.Name = userDto.Username
	}
	if existingUser.Role != "" {
		existingUser.Role = userDto.Role
	}
	if userDto.Password != "" {
		hashPwd, err := utils.Hash(userDto.Password)
		if err != nil {
			return nil, fmt.Errorf("更新密码失败: %w", err)
		}
		existingUser.Password = hashPwd
	}
	// 4. 执行更新
	if err := us.userDao.Update(existingUser); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	existingUser.Password = ""
	return existingUser, nil
}

func (us *UserService) GetById() {}

func (us *UserService) List() ([]dto.UserRequest, error) {
	dbUsers, err := us.userDao.GetAllUsers()
	if err != nil {
		return nil, err
	}
	var dtoUsers []dto.UserRequest
	for _, dbUser := range dbUsers {
		user := dto.UserRequest{
			Id:       dbUser.ID,
			Username: dbUser.Name,
			Role:     dbUser.Role,
		}
		dtoUsers = append(dtoUsers, user)
	}
	return dtoUsers, nil
}

func (us *UserService) Auth(userDto dto.LoginRequest) (*model.PbUser, error) {
	user := model.PbUser{
		Name:     userDto.Username,
		Password: userDto.Password,
	}
	dbUser, err := us.userDao.Auth(&user)
	if err != nil {
		return nil, err
	}
	ok, err := utils.Verify(dbUser.Password, userDto.Password)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("认证失败")

	}
	return dbUser, nil
}
