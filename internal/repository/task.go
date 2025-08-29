package repository

import (
	"pubot/internal/model"
)

type TaskRepo interface {
	Create(modelTask model.Task) error
	Delete(id uint) error
	Update(modelTask model.Task) error
	Get(id uint) (*model.Task, error)
	List() ([]model.Task, error)
	State(id uint) (*model.Task, error)
}
