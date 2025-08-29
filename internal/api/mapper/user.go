package mapper

import (
	"pubot/internal/api/dto"
	"pubot/internal/model"
)

func UserModelToRespUserDTO(user model.User) dto.RespUser {
	return dto.RespUser{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt.Time,
	}
}

func ReqUserDTOToUserModel(user dto.ReqCreateUser) model.User {
	return model.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
		Password: user.Password,
	}
}

func AuthUserDTOToUserModel(user dto.AuthUser) model.User {
	return model.User{
		Name:     user.Name,
		Email:    user.Emal,
		Password: user.Password,
	}
}

func UserModelsToRespUsersDTO(users []model.User) []dto.RespUser {
	dtoUsers := make([]dto.RespUser, 0)
	for _, user := range users {
		dtoUsers = append(dtoUsers, UserModelToRespUserDTO(user))
	}
	return dtoUsers
}
