package service

import (
	"pubot/internal/api/dto"
	"pubot/internal/api/mapper"
	"pubot/internal/repository"
	"pubot/internal/core/pwd"
)

type User struct {
	userRepo repository.UserRepo
}

func NewUser(userRepo repository.UserRepo) *User {
	return &User{userRepo: userRepo}
}

func (user *User) Create(dtoUser dto.ReqUser) (dto.RespUser, error) {
	modelUser := mapper.ReqUserDTOToUserModel(dtoUser)
	modelUser.Password = pwd.Hash(dtoUser.Password)
	err := user.userRepo.Create(modelUser)
	return mapper.UserModelToRespUserDTO(modelUser), err
}

func (user *User) Delete(id uint) error {
	dtoUser := dto.ReqUser{Id: id}
	modelUser := mapper.ReqUserDTOToUserModel(dtoUser)
	return user.userRepo.Delete(modelUser)
}
