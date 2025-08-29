package dao

import (
	"pubot/internal/model"

	"gorm.io/gorm"
)

type TaskDao struct {
	db *gorm.DB
}

func NewTaskDao(db *gorm.DB) *TaskDao {
	return &TaskDao{db: db}
}

func (taskDao *TaskDao) Create(modelTask model.Task) error {

	return nil
}

func (taskDao *TaskDao) Delete(id uint) error {

	return nil
}

func (taskDao *TaskDao) Update(modelTask model.Task) error {

	return nil
}

func (taskDao *TaskDao) Get(id uint) (*model.Task, error) {

	return nil, nil
}

func (taskDao *TaskDao) List() ([]model.Task, error) {

	return nil, nil
}

func (taskDao *TaskDao) State(id uint) (*model.Task, error) {

	return nil, nil
}
