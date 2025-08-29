package api

import (
	"net/http"
	"pubot/internal/core/web"
	"pubot/internal/service"
)

type TaskController struct {
	taskService *service.TaskService
}

func NewTaskController(taskService *service.TaskService) *TaskController {
	return &TaskController{taskService: taskService}
}

func (taskController *TaskController) RegisterRouter(router *web.Router) {
	router.Post("/task", taskController.create)
	router.Delete("/task/:id", taskController.delete)
	router.Put("/task/:id", taskController.update)
	router.Get("/task/:id", taskController.get)
	router.Get("/tasks", taskController.list)
	router.Get("/task/state/:id", taskController.state)
	router.Post("/task/:name", taskController.run)
}

func (taskController *TaskController) create(w http.ResponseWriter, r *http.Request) {}

func (taskController *TaskController) delete(w http.ResponseWriter, r *http.Request) {}

func (taskController *TaskController) update(w http.ResponseWriter, r *http.Request) {}

func (taskController *TaskController) get(w http.ResponseWriter, r *http.Request) {}

func (taskController *TaskController) list(w http.ResponseWriter, r *http.Request) {
	web.Success(w, web.Map{"code": 200, "message": "success"})
}

func (taskController *TaskController) state(w http.ResponseWriter, r *http.Request) {}

func (taskController *TaskController) run(w http.ResponseWriter, r *http.Request) {}
