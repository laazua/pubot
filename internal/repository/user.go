package repository

import "pubot/internal/model"

type UserRepo interface {
	Create(user model.User) error
	Delete(user model.User) error
	Update(user model.User) error
	Get(modelUser model.User) (model.User, error)
	List() ([]model.User, error)
}
