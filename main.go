package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"pubot/internal/api"
	"pubot/internal/config"
	"pubot/internal/dao"
	"pubot/internal/service"
	"pubot/internal/utils"

	"github.com/gorilla/mux"
)

func main() {
	// 依赖注入
	userDao := dao.NewUserDao(dao.GetDb())
	userService := service.NewUserService(userDao)
	userApi := api.NewUserApi(userService)
	hub := utils.NewHub()
	taskDao := dao.NewTaskDao(dao.GetDb())
	taskService := service.NewTaskService(taskDao, hub)
	taskApi := api.NewTaskApi(taskService)

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()
	// 登录路由 - 不需要认证中间件
	apiRouter.HandleFunc("/login", userApi.Login).Methods("POST")
	// 用户路由分组
	userRouter := router.PathPrefix("/api").Subrouter()
	userRouter.Use(utils.AuthMw, utils.CorsMw)
	userApi.Register(userRouter)
	// 任务路由分组
	taskRouter := router.PathPrefix("/api").Subrouter()
	taskRouter.Use(utils.AuthMw, utils.CorsMw)
	taskApi.Register(taskRouter)
	wsTaskRouter := router.PathPrefix("/ws").Subrouter()
	wsTaskRouter.Use(utils.AuthWsMw) // 先 Use，再注册路由
	wsTaskRouter.HandleFunc("/task", hub.ServeWS)
	// 前端静态资源（放在最后，避免覆盖API路由）
	router.PathPrefix("/").Handler(api.WebHandler())
	server := http.Server{
		Handler: router,
		Addr:    config.Get().Listen,
	}
	start := make(chan error, 1)
	quit := make(chan os.Signal, 1)
	// 监听失败和退出信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func(start chan error) {
		slog.Info("pubot 程序启动...", slog.String("Listen", config.Get().Listen))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			start <- err
		}
	}(start)
	// 监听退出
	select {
	case err := <-start:
		slog.Error("pubot 启动失败", slog.String("Err", err.Error()))
	case ext := <-quit:
		slog.Info("pubot 程序关闭...", slog.Any("Shutdown", ext))
		dao.CloseDb()
	}
}
