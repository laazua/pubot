package api

import (
	"net/http"
	"strconv"

	"pubot/internal/api/dto"
	"pubot/internal/core/logx"
	"pubot/internal/core/web"
	"pubot/internal/service"
)

// 用户控制器
type UserController struct {
	userService *service.User
}

func NewUserController(userService *service.User) *UserController {
	return &UserController{userService: userService}
}

// 用户路由注册
func (userController *UserController) RegisterRoute(route *web.Router) {
	route.Post("/api/user", userController.create)
	route.Delete("/api/user/:id", userController.delete)
	route.Put("/api/user/:id", userController.update)
	route.Get("/api/user/:id", userController.get)
	route.Get("/api/user", userController.list)

}

// 创建用户
func (userController *UserController) create(w http.ResponseWriter, r *http.Request) {
	var userReq dto.ReqUser
	if err := web.Bind(r, &userReq); err != nil {
		logx.Error("绑定请求体参数失败", logx.String("Err", err.Error()))
		web.Failure(w, web.Map{"code": 400, "message": err.Error()})
		return
	}
	userResp, err := userController.userService.Create(userReq)
	if err != nil {
		logx.Error("创建用户失败", logx.String("Err", err.Error()))
		web.Failure(w, web.Map{"code": 400, "message": err.Error()})
		return
	}
	web.Success(w, web.Map{"code": 200, "message": "create user success", "data": userResp})
}

// 删除用户
func (UserController *UserController) delete(w http.ResponseWriter, r *http.Request) {
	strId := web.Param(r, "id")
	id, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		logx.Error("路径参数获取失败", logx.String("Err", err.Error()))
		web.Failure(w, web.Map{"code": 400, "message": err.Error()})
		return
	}
	if err := UserController.userService.Delete(uint(id)); err != nil {
		logx.Error("删除用户失败", logx.String("Err", err.Error()))
		web.Failure(w, web.Map{"code": 400, "message": err.Error()})
		return
	}
	web.Success(w, web.Map{"code": 200, "message": "delete user success"})
}

func (userController *UserController) update(w http.ResponseWriter, r *http.Request) {}

func (userController *UserController) get(w http.ResponseWriter, r *http.Request) {}

func (userController *UserController) list(w http.ResponseWriter, r *http.Request) {}
