//    Copyright [2025] laazua
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"pubot/internal/api"
	"pubot/internal/core/config"
	"pubot/internal/core/db"
	"pubot/internal/core/logx"
	"pubot/internal/core/web"
	"pubot/internal/dao"
	"pubot/internal/service"
)

var (
	start = make(chan error, 1)
	quit  = make(chan os.Signal, 1)
)

func main() {
	startPubot()
}

// startPubot 初始化并启动
func startPubot() {
	dB, err := db.Get()
	if err != nil {
		start <- err
	}
	// 依赖组装
	authDao := dao.NewAuth(dB)
	userDao := dao.NewUser(dB)
	authService := service.NewAuth(authDao)
	userService := service.NewUser(userDao)
	authController := api.NewAuthController(authService)
	userController := api.NewUserController(userService)
	// 路由注册
	route := web.NewRouter()
	v1 := route.Group("/v1")
	authController.RegisterRoute(v1)
	userController.RegisterRoute(v1)
	// 实例化http server
	server := http.Server{
		Handler:      route,
		Addr:         config.Get().Address,
		ReadTimeout:  config.Get().ReadTimeout,
		WriteTimeout: config.Get().WriteTimeout,
	}
	// 协程启动服务
	go func() {
		logx.Info("pubot启动成功,", logx.String("Address", config.Get().Address))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			start <- err
		}
	}()
	// 监听失败和退出信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case err := <-start:
		logx.Error("pubot启动失败", logx.String("Error", err.Error()))
	case sig := <-quit:
		logx.Info("pubot关闭,并清理资源", logx.String("Signal", sig.String()))
		func() { _ = db.Close(); _ = logx.Close() }()
	}
}
