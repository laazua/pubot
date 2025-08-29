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
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

// 用户路由注册
func (userController *UserController) RegisterRouter(router *web.Router) {
	router.Post("/user", userController.create)
	router.Delete("/user/:id", userController.delete)
	router.Put("/user", userController.update)
	router.Get("/user/:id", userController.get)
	router.Get("/users", userController.list)
}

// 创建用户(ok)
func (userController *UserController) create(w http.ResponseWriter, r *http.Request) {
	logx.Info("create user ...")
	var userReq dto.ReqCreateUser
	if err := web.Bind(r, &userReq); err != nil {
		logx.Error("绑定请求体参数失败", logx.String("Err", err.Error()))
		web.Failure(w, web.Map{"code": 400, "message": err.Error()})
		return
	}
	err := userController.userService.Create(userReq)
	if err != nil {
		logx.Error("创建用户失败", logx.String("Err", err.Error()))
		web.Failure(w, web.Map{"code": 400, "message": err.Error()})
		return
	}
	web.Success(w, web.Map{"code": 200, "message": "创建用户成功"})
}

// 删除用户(ok)
func (userController *UserController) delete(w http.ResponseWriter, r *http.Request) {
	idStr := web.Param(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		logx.Error("路径参数获取失败", logx.String("Err", err.Error()))
		web.Failure(w, web.Map{"code": 400, "message": err.Error()})
		return
	}
	if err := userController.userService.Delete(uint(id)); err != nil {
		logx.Error("删除用户失败", logx.String("Err", err.Error()))
		web.Failure(w, web.Map{"code": 400, "message": err.Error()})
		return
	}
	web.Success(w, web.Map{"code": 200, "message": "delete user success"})
}

// 更新用户(ok)
func (userController *UserController) update(w http.ResponseWriter, r *http.Request) {
	var dtoUser dto.ReqCreateUser
	if err := web.Bind(r, &dtoUser); err != nil {
		logx.Error("更新用户绑定请求体失败", logx.String("Err", err.Error()))
		web.Failure(w, web.Map{"code": 400, "message": err.Error()})
		return
	}
	if err := userController.userService.Update(dtoUser); err != nil {
		logx.Error("更新用户数据库失败", logx.String("Err", err.Error()))
		web.Failure(w, web.Map{"code": 400, "message": err.Error()})
		return
	}
	web.Success(w, web.Map{"code": 200, "message": "更新用户成功"})
}

// 查询用户(ok)
func (userController *UserController) get(w http.ResponseWriter, r *http.Request) {
	idStr := web.Param(r, "id")
	logx.Info("用户ID", logx.String("ID", idStr))
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		logx.Error("路径参数获取失败", logx.String("Err", err.Error()))
		web.Failure(w, web.Map{"code": 400, "message": err.Error()})
		return
	}
	user, err := userController.userService.Get(uint(id))
	if err != nil {
		logx.Error("查询用户失败", logx.String("Err", err.Error()))
		web.Failure(w, web.Map{"code": 400, "message": err.Error()})
		return
	}
	web.Success(w, web.Map{"code": 200, "message": "查询用户成功", "data": user})
}

func (userController *UserController) list(w http.ResponseWriter, r *http.Request) {
	users, err := userController.userService.List()
	if err != nil {
		logx.Error("查询用户列表失败", logx.String("Err", err.Error()))
		web.Failure(w, web.Map{"code": 400, "message": err.Error()})
		return
	}
	web.Success(w, web.Map{"code": 200, "message": "查询用户成功", "data": users})
}
