package api

import (
	"net/http"

	"pubot/internal/api/dto"
	"pubot/internal/core/logx"
	"pubot/internal/core/token"
	"pubot/internal/core/web"
	"pubot/internal/service"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authservice *service.AuthService) *AuthController {
	return &AuthController{authService: authservice}
}

func (authController *AuthController) RegisterRouter(router *web.Router) {
	router.Post("/login", authController.login)
}

// 登录(ok)
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
	// 生成token
	tokenStr, err := token.Create(authUser.Name)
	if err != nil {
		logx.Error("生成token失败", logx.String("Err", err.Error()))
		web.Failure(w, web.Map{"code": 400, "message": err.Error()})
		return
	}
	web.Success(w, web.Map{"code": 200, "message": "user login success", "token": tokenStr})
}
