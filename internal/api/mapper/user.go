package mapper

import (
	"pubot/internal/api/dto"
	"pubot/internal/model"
)

func UserModelToRespUserDTO(user model.User) dto.RespUser {
	return dto.RespUser{
		Id:   user.ID,
		Name: user.Name,
		Emal: user.Email,
	}
}

func ReqUserDTOToUserModel(user dto.ReqUser) model.User {
	return model.User{
		Name:  user.Name,
		Email: user.Emal,
	}
}

func AuthUserDTOToUserModel(user dto.AuthUser) model.User {
	return model.User{
		Name:     user.Name,
		Password: user.Password,
	}
}
