package service

import (
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

func (us *UserService) Delete() {}

func (us *UserService) Update() {}

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
	_, err = utils.Verify(dbUser.Password, userDto.Password)
	if err != nil {
		return nil, err
	}
	return dbUser, nil
}
