package mw

import (
	"net/http"
	"pubot/internal/core/logx"
	"pubot/internal/core/token"
	"pubot/internal/core/web"
)

// AuthMiddleWare token认证中间件
func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		payload, err := token.Parse(tokenStr)
		if err != nil {
			logx.Debug("token认证失败...")
			web.Failure(w, web.Map{"code": 404, "message": err.Error()})
			return
		}
		logx.Debug("token认证成功...", logx.String("Name", payload.Name))
		next.ServeHTTP(w, r)
	})
}
