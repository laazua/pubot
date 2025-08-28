package service

import (
	"pubot/internal/api/dto"
	"pubot/internal/api/mapper"
	"pubot/internal/repository"
)

type Auth struct {
	authRepo repository.AuthRepo
}

func NewAuth(authRepo repository.AuthRepo) *Auth {
	return &Auth{authRepo: authRepo}
}

func (auth *Auth) Login(dtoAuth dto.AuthUser) error {
	modelUser := mapper.AuthUserDTOToUserModel(dtoAuth)
	return auth.authRepo.Sign(modelUser)
}
