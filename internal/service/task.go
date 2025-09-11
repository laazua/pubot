package service

import (
	"encoding/json"
	"fmt"
	"pubot/internal/dao"
	"pubot/internal/dto"
	"pubot/internal/model"
	"pubot/internal/utils"

	"gopkg.in/yaml.v3"
)

type TaskService struct {
	hub     *utils.Hub
	taskDao *dao.TaskDao
}

func NewTaskService(taskDao *dao.TaskDao, hub *utils.Hub) *TaskService {
	return &TaskService{taskDao: taskDao, hub: hub}
}

func (ts *TaskService) Create(taskDto dto.TaskCreateRequest) (*model.PbTask, error) {
	parsed, err := utils.ParseTaskYAML(taskDto.YAML)
	if err != nil {
		return nil, err
	}
	tparsedJSON, err := json.Marshal(parsed)
	if err != nil {
		return nil, err
	}

	task := model.PbTask{
		Name:       taskDto.Name,
		YAML:       taskDto.YAML,
		Status:     "stopped",
		YAMLParsed: tparsedJSON,
		Count:      0,
	}

	if err := ts.taskDao.Create(&task); err != nil { // ✅ 注意这里是 &task
		return nil, err
	}

	return &task, nil // ✅ 现在返回时 task.ID 已经是数据库里的真实 ID
}

func (ts *TaskService) Delete(id uint) error {
	// 先检查任务是否存在
	_, err := ts.taskDao.GetByID(id)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// 执行删除操作
	if err := ts.taskDao.Delete(id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	// 可以在这里添加一些额外的逻辑，比如：
	// - 记录删除日志
	// - 发送删除通知
	// - 清理相关资源等
	return nil
}

func (ts *TaskService) Update(id uint, dtoTask dto.TaskUpdateRequest) (*model.PbTask, error) {
	// 1. 先查找现有任务
	existingTask, err := ts.taskDao.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 2. 如果有 YAML 更新，需要重新解析
	if dtoTask.YAML != "" {
		parsed, err := utils.ParseTaskYAML(dtoTask.YAML)
		if err != nil {
			return nil, fmt.Errorf("invalid YAML: %w", err)
		}

		parsedJSON, err := json.Marshal(parsed)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal parsed YAML: %w", err)
		}

		existingTask.YAML = dtoTask.YAML
		existingTask.YAMLParsed = parsedJSON
	}

	// 3. 更新其他字段
	if dtoTask.Name != "" {
		existingTask.Name = dtoTask.Name
	}

	if dtoTask.Status != "" {
		existingTask.Status = dtoTask.Status
	}

	// 4. 执行更新
	if err := ts.taskDao.Update(existingTask); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return existingTask, nil
}

func (ts *TaskService) GetById(id uint) (*model.PbTask, error) {
	return ts.taskDao.GetByID(id)
}

func (ts *TaskService) List() ([]model.PbTask, error) {
	return ts.taskDao.GetAllTask()
}

func (ts *TaskService) Execute(id uint) error {
	task, err := ts.taskDao.GetByID(id)
	if err != nil {
		return err
	}

	go func(t *model.PbTask) {
		// 1️⃣ 开始执行任务：持久化 running 状态
		t.Status = string(utils.TaskRunning)
		if err := ts.taskDao.Save(t); err != nil {
			ts.hub.Broadcast(utils.TaskStatus{ID: t.ID, Status: utils.TaskError, Count: t.Count})
			return
		}
		ts.hub.Broadcast(utils.TaskStatus{ID: t.ID, Status: utils.TaskRunning, Count: t.Count})

		var parsed map[string]interface{}
		if err := yaml.Unmarshal([]byte(t.YAML), &parsed); err != nil {
			// YAML 解析失败 → error
			t.Status = string(utils.TaskError)
			_ = ts.taskDao.Save(t)
			ts.hub.Broadcast(utils.TaskStatus{ID: t.ID, Status: utils.TaskError, Count: t.Count})
			return
		}

		// 2️⃣ 执行 build 阶段
		if buildSteps, ok := parsed["build"].([]interface{}); ok {
			for _, cmd := range buildSteps {
				if inerr := utils.RunCmd(cmd.(string)); inerr != nil {
					t.Status = string(utils.TaskError)
					_ = ts.taskDao.Save(t)
					ts.hub.Broadcast(utils.TaskStatus{ID: t.ID, Status: utils.TaskError, Count: t.Count})
					return
				}
			}
		}

		// 3️⃣ 执行 deploy 阶段
		if deploy, ok := parsed["deploy"].(map[string]interface{}); ok {
			if runSteps, ok := deploy["run"].([]interface{}); ok {
				for _, cmd := range runSteps {
					if inerr := utils.RunCmd(cmd.(string)); inerr != nil {
						t.Status = string(utils.TaskError)
						_ = ts.taskDao.Save(t)
						ts.hub.Broadcast(utils.TaskStatus{ID: t.ID, Status: utils.TaskError, Count: t.Count})
						return
					}
				}
			}
		}

		// 4️⃣ 成功完成：Count +1 并保存状态 success
		t.Count++
		t.Status = string(utils.TaskSuccess)
		if err := ts.taskDao.Save(t); err != nil {
			// 保存失败 → error
			ts.hub.Broadcast(utils.TaskStatus{ID: t.ID, Status: utils.TaskError, Count: t.Count})
			return
		}

		// 广播成功状态
		ts.hub.Broadcast(utils.TaskStatus{ID: t.ID, Status: utils.TaskSuccess, Count: t.Count})

	}(task)

	return nil
}
