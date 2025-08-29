package service

import (
	"pubot/internal/api/dto"
	"pubot/internal/api/mapper"
	"pubot/internal/core/logx"
	"pubot/internal/repository"
)

type AuthService struct {
	authRepo repository.AuthRepo
}

func NewAuthService(authRepo repository.AuthRepo) *AuthService {
	return &AuthService{authRepo: authRepo}
}

func (authService *AuthService) Login(dtoAuth dto.AuthUser) error {
	modelUser := mapper.AuthUserDTOToUserModel(dtoAuth)
	logx.Info("auth service", logx.String("user", modelUser.Name))
	return authService.authRepo.Sign(modelUser)
}
