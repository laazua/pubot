package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"pubot/internal/config"
	"pubot/internal/dto"
	"pubot/internal/model"
	"pubot/internal/service"
	"pubot/internal/utils"

	"github.com/gorilla/mux"
)

type UserApi struct {
	mu          sync.Mutex
	userService *service.UserService
}

func NewUserApi(userService *service.UserService) *UserApi {
	return &UserApi{userService: userService}
}

func (ua *UserApi) Register(router *mux.Router) {
	router.HandleFunc("/user", ua.create).Methods("POST")
	router.HandleFunc("/user/{id:[0-9]+}", ua.delete).Methods("DELETE")
	router.HandleFunc("/user/{id:[0-9]+}", ua.update).Methods("PUT")
	router.HandleFunc("/user", ua.list).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}", ua.get).Methods("GET") // 限制id只能是数字
	router.HandleFunc("/user/info", ua.info).Methods("GET")
}

func (ua *UserApi) Login(w http.ResponseWriter, r *http.Request) {
	slog.Info("login api ...")
	var req dto.LoginRequest
	if err := utils.Bind(r, &req); err != nil {
		utils.Failure(w, utils.Map{"code": 403, "message": "解析参数失败"})
		return
	}
	// 校验用户名密码
	user, err := ua.userService.Auth(req)
	if err != nil {
		slog.Error("用户认证失败", slog.String("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 401, "message": "用户认证失败"})
		return
	}
	auth := utils.NewJWTAuth(config.Get().SecretKey, config.Get().ExpiredTime)
	tokenString, err := auth.GenerateToken(user.ID, user.Name, user.Role)
	if err != nil {
		utils.Failure(w, utils.Map{"code": 500, "message": "生成token失败"})
		return
	}
	utils.Success(w, utils.Map{"code": 200, "message": "登录成功", "token": tokenString})
}

func (ua *UserApi) create(w http.ResponseWriter, r *http.Request) {
	var req dto.UserRequest
	if err := utils.Bind(r, &req); err != nil {
		slog.Error("绑定请求体参数失败", slog.String("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 403, "message": "绑定请求参数失败"})
		return
	}
	user, err := ua.userService.Create(req)
	if err != nil {
		slog.Error("创建用户失败", slog.String("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 500, "message": "创建用户失败"})
		return
	}
	utils.Success(w, utils.Map{"code": 200, "message": "创建用户成功", "data": user})
}

func (ua *UserApi) delete(w http.ResponseWriter, r *http.Request) {
	userIdStr := mux.Vars(r)["id"]
	// 如果需要数字类型，需要手动转换
	userId, err := strconv.ParseUint(userIdStr, 10, 0)
	if err != nil {
		slog.Error("无效的task ID", slog.Any("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 400, "message": "无效的 task ID"})
		return
	}
	if err := ua.userService.Delete(uint(userId)); err != nil {
		slog.Error("删除用户失败", slog.String("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 503, "message": "删除用户失败"})
		return
	}
	utils.Success(w, utils.Map{"code": 200, "message": "删除用户成功"})
}

func (ua *UserApi) update(w http.ResponseWriter, r *http.Request) {
	ua.mu.Lock()
	defer ua.mu.Unlock()
	userIdStr := mux.Vars(r)["id"]
	// 如果需要数字类型，需要手动转换
	userId, err := strconv.ParseUint(userIdStr, 10, 0)
	if err != nil {
		slog.Error("无效的task ID", slog.Any("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 400, "message": "无效的 task ID"})
		return
	}
	var req dto.UserRequest
	if err := utils.Bind(r, &req); err != nil {
		slog.Error("绑定请求体参数失败", slog.String("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 200, "message": "绑定请求头参数失败"})
		return
	}
	user, err := ua.userService.Update(uint(userId), req)
	if err != nil {
		slog.Error("更新用户失败", slog.String("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 503, "message": "跟新用户失败"})
		return
	}
	utils.Success(w, utils.Map{"code": 200, "message": "更新用户成功", "data": user})
}

func (ua *UserApi) get(w http.ResponseWriter, r *http.Request) {}

func (ua *UserApi) list(w http.ResponseWriter, r *http.Request) {
	users, err := ua.userService.List()
	if err != nil {
		slog.Error("获取用户列表失败", slog.String("Err", err.Error()))
		utils.Failure(w, utils.Map{"code": 200, "message": "获取用户列表失败"})
		return
	}

	utils.Success(w, utils.Map{"code": 200, "message": "获取用户列表成功", "data": users})
}

func (ua *UserApi) info(w http.ResponseWriter, r *http.Request) {
	slog.Info("info api ...")
	user, ok := r.Context().Value(utils.ContextUserKey).(*model.PbUser)
	if !ok || user == nil {
		utils.Failure(w, utils.Map{"code": 401, "message": "用户未登录"})
		return
	}

	data := map[string]any{
		"id":       user.ID,
		"username": user.Name,
		"role":     user.Role,
	}

	utils.Success(w, utils.Map{"code": 200, "message": "获取用户信息成功", "data": data})
}
