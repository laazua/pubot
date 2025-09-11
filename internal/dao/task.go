package dao

import (
	"errors"
	"pubot/internal/model"

	"gorm.io/gorm"
)

type TaskDao struct {
	db *gorm.DB
}

func NewTaskDao(db *gorm.DB) *TaskDao {
	return &TaskDao{db: db}
}

func (td *TaskDao) Create(dbTask *model.PbTask) error {
	var pipelineExists model.PbTask
	if td.db.Where("name = ?", dbTask.Name).First(&pipelineExists).Error == nil {
		return errors.New("任务已经存在")
	}
	return td.db.Create(&dbTask).Error
}

func (td *TaskDao) Delete(id uint) error {
	return td.db.Where("id = ?", id).Delete(&model.PbTask{}).Error
}

func (td *TaskDao) GetByID(id uint) (*model.PbTask, error) {
	var modelTask model.PbTask
	err := td.db.First(&modelTask, id).Error
	if err != nil {
		return nil, err
	}
	return &modelTask, nil
}

func (td *TaskDao) Update(task *model.PbTask) error {
	result := td.db.Save(task)
	return result.Error
}

func (td *TaskDao) GetAllTask() ([]model.PbTask, error) {
	var modelTasks []model.PbTask
	err := td.db.Find(&modelTasks).Error
	if err != nil {
		return nil, err
	}
	return modelTasks, nil
}

func (td *TaskDao) Save(task *model.PbTask) error {
	return td.db.Save(&task).Error
}
