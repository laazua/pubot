package service

import (
	"pubot/internal/api/dto"
	"pubot/internal/api/mapper"
	"pubot/internal/core/logx"
	"pubot/internal/core/pwd"
	"pubot/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (userService *UserService) Create(dtoUser dto.ReqCreateUser) error {
	modelUser := mapper.ReqUserDTOToUserModel(dtoUser)
	logx.Info("create user service: ", logx.Any("isAdmin", modelUser.IsAdmin))
	modelUser.Password = pwd.Hash(dtoUser.Password)
	err := userService.userRepo.Create(modelUser)
	return err
}

func (userService *UserService) Delete(id uint) error {
	dtoUser := dto.ReqCreateUser{ID: id}
	modelUser := mapper.ReqUserDTOToUserModel(dtoUser)
	return userService.userRepo.Delete(modelUser)
}

func (userService *UserService) Update(dtoUser dto.ReqCreateUser) error {
	modelUser := mapper.ReqUserDTOToUserModel(dtoUser)
	if modelUser.Password != "" {
		modelUser.Password = pwd.Hash(modelUser.Password)
		logx.Info("service更新密码", logx.String("Pass", modelUser.Password))
	}
	return userService.userRepo.Update(modelUser)
}

func (userService *UserService) Get(id uint) (dto.RespUser, error) {
	dtoUser := dto.ReqCreateUser{ID: id}
	modelUser := mapper.ReqUserDTOToUserModel(dtoUser)
	user, err := userService.userRepo.Get(modelUser)
	if err != nil {
		return dto.RespUser{}, err
	}
	return mapper.UserModelToRespUserDTO(user), nil
}

func (userService *UserService) List() ([]dto.RespUser, error) {
	modelUsers, err := userService.userRepo.List()
	if err != nil {
		return nil, err
	}
	return mapper.UserModelsToRespUsersDTO(modelUsers), nil
}
