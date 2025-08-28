package api

import (
	"net/http"

	"pubot/internal/api/dto"
	"pubot/internal/core/logx"
	"pubot/internal/core/web"
	"pubot/internal/service"
)

type AuthController struct {
	authService *service.Auth
}

func NewAuthController(authservice *service.Auth) *AuthController {
	return &AuthController{authService: authservice}
}

func (authController *AuthController) RegisterRoute(route *web.Router) {
	route.Post("/login", authController.login)
	route.Get("/test", authController.testApi)
}

func (authController *AuthController) login(w http.ResponseWriter, r *http.Request) {
	var authUser dto.AuthUser
	if err := web.Bind(r, &authUser); err != nil {
		logx.Error("绑定请求体参数失败", logx.String("Err", err.Error()))
		web.Failure(w, web.Map{"code": 400, "message": err.Error()})
		return
	}
	if err := authController.authService.Login(authUser); err != nil {
		logx.Error("用户认证失败", logx.String("Err", err.Error()))
		web.Failure(w, web.Map{"code": 400, "message": err.Error()})
		return
	}
	web.Success(w, web.Map{"code": 200, "message": "user login success"})
}

func (authController *AuthController) testApi(w http.ResponseWriter, r *http.Request) {
	logx.Info("this is a test")
	web.Success(w, web.Map{"code": 200, "message": "this is a test api"})
}
