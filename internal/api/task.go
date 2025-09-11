package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"pubot/internal/dto"
	"pubot/internal/service"
	"pubot/internal/utils"

	"github.com/gorilla/mux"
)

type TaskApi struct {
	mu          sync.Mutex
	taskService *service.TaskService
}

func NewTaskApi(taskService *service.TaskService) *TaskApi {
	return &TaskApi{
		taskService: taskService,
	}
}

func (ta *TaskApi) Register(router *mux.Router) {
	router.HandleFunc("/task", ta.create).Methods("POST")
	router.HandleFunc("/task/{id:[0-9]+}", ta.delete).Methods("DELETE")
	router.HandleFunc("/task/{id:[0-9]+}", ta.update).Methods("PUT")
	router.HandleFunc("/task", ta.list).Methods("GET")
	router.HandleFunc("/task/{id:[0-9]+}", ta.get).Methods("GET")
	router.HandleFunc("/task/{id:[0-9]+}", ta.execute).Methods("POST")
}

// create 创建流水线任务模板(ok)
func (ta *TaskApi) create(w http.ResponseWriter, r *http.Request) {
	var req dto.TaskCreateRequest
	if err := utils.Bind(r, &req); err != nil {
		slog.Error("接口参数绑定失败", slog.Any("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 403, "message": "接口参数绑定失败"})
		return
	}
	fmt.Printf("%#v\n", req)
	task, err := ta.taskService.Create(req)
	if err != nil {
		slog.Error("添加流水线任务失败", slog.Any("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 503, "message": "添加流水线任务失败"})
		return
	}
	utils.Success(w, utils.Map{"code": 200, "message": "添加流水线任务成功", "data": task})
}

// delete 删除流水线任务模板
func (ta *TaskApi) delete(w http.ResponseWriter, r *http.Request) {
	taskIdStr := mux.Vars(r)["id"]
	// 如果需要数字类型，需要手动转换
	taskId, err := strconv.ParseUint(taskIdStr, 10, 0)
	if err != nil {
		slog.Error("无效的task ID", slog.Any("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 400, "message": "无效的 task ID"})
		return
	}
	err = ta.taskService.Delete(uint(taskId))
	if err != nil {
		slog.Error("删除任务失败", slog.Any("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 503, "message": "删除任务失败"})
		return
	}
	utils.Success(w, utils.Map{"code": 200, "message": "删除任务成功"})
}

// update 更新流水线任务模板(ok)
func (ta *TaskApi) update(w http.ResponseWriter, r *http.Request) {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	taskIdStr := mux.Vars(r)["id"]
	// 如果需要数字类型，需要手动转换
	taskId, err := strconv.ParseUint(taskIdStr, 10, 0)
	if err != nil {
		slog.Error("无效的task ID", slog.Any("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 400, "message": "无效的 task ID"})
		return
	}
	var req dto.TaskUpdateRequest
	if err := utils.Bind(r, &req); err != nil {
		slog.Error("绑定请求体参数失败", slog.Any("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 400, "message": "绑定请求体参数失败"})
		return
	}
	task, err := ta.taskService.Update(uint(taskId), req)
	if err != nil {
		slog.Error("更新流水线任务失败", slog.Any("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 500, "message": "更新流水线任务失败"})
		return
	}
	utils.Success(w, utils.Map{"code": 200, "message": "更新流水线任务成功", "data": task})
}

// get 获取单个流水线任务
func (ta *TaskApi) get(w http.ResponseWriter, r *http.Request) {
	taskIdStr := mux.Vars(r)["id"]
	// 如果需要数字类型，需要手动转换
	taskId, err := strconv.ParseUint(taskIdStr, 10, 0)
	if err != nil {
		slog.Error("无效的task ID", slog.Any("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 400, "message": "无效的 task ID"})
		return
	}
	task, err := ta.taskService.GetById(uint(taskId))
	if err != nil {
		slog.Error("获取任务失败", slog.Any("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 500, "message": "获取任务失败"})
		return
	}
	utils.Success(w, utils.Map{"code": 200, "message": "获取任务成功", "data": task})
}

// list 获取任务列表(ok)
func (ta *TaskApi) list(w http.ResponseWriter, r *http.Request) {
	tasks, err := ta.taskService.List()
	if err != nil {
		slog.Error("获取流水线任务列表失败", slog.Any("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": "503", "message": "获取流水线任务列表失败", "data": nil})
		return
	}
	utils.Success(w, utils.Map{"code": 200, "message": "获取任务列表成功", "data": tasks})
}

func (ta *TaskApi) execute(w http.ResponseWriter, r *http.Request) {
	taskIdStr := mux.Vars(r)["id"]
	// 如果需要数字类型，需要手动转换
	taskId, err := strconv.ParseUint(taskIdStr, 10, 0)
	if err != nil {
		slog.Error("无效的task ID", slog.Any("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 400, "message": "无效的 task ID"})
		return
	}
	err = ta.taskService.Execute(uint(taskId))
	if err != nil {

		return
	}

	utils.Success(w, utils.Map{"code": 200, "message": "执行任务操作成功"})
}
