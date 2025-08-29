package service

import (
	"pubot/internal/api/dto"
	"pubot/internal/repository"
)

type TaskService struct {
	taskRepo repository.TaskRepo
}

func NewTaskService(taskRepo repository.TaskRepo) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (taskService *TaskService) Create(dtoTask dto.Task) error {
	return nil
}
func (taskService *TaskService) Delete(id uint) error {
	return nil
}

func (taskService *TaskService) Update(dtoTask dto.Task) error {

	return nil
}

func (taskService *TaskService) Get(id uint) (*dto.Task, error) {

	return nil, nil
}

func (taskService *TaskService) List() ([]*dto.Task, error) {

	return nil, nil
}

func (taskService *TaskService) State(id uint) (*dto.Task, error) {

	return nil, nil
}
