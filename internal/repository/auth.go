package repository

import "pubot/internal/model"

type AuthRepo interface {
	Sign(user model.User) error
}
