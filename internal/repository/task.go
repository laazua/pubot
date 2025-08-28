package repository

import "pubot/internal/model"

type TaskRepo interface {
	Create(task model.Task) error
	Delete(task model.Task) error
	Update(task model.Task) error
	Get(task model.Task) (model.Task, error)
	List() ([]model.Task, error)
	State(task model.Task) (model.Task, error)
}
