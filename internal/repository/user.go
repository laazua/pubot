package repository

import "pubot/internal/model"

type UserRepo interface {
	Create(user model.User) error
	Delete(user model.User) error
}
