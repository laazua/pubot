package utils

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"pubot/internal/config"
	"pubot/internal/model"

	"github.com/gorilla/websocket"
)

type contextKey string

const ContextUserKey = contextKey("user")

func AuthMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenStr := r.Header.Get("Authorization")
		const BearerSchema = "Bearer "
		if !strings.HasPrefix(tokenStr, BearerSchema) {
			slog.Error("token格式错误", slog.Any("Token", tokenStr))
			Failure(w, Map{"code": 403, "message": "token格式错误"})
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(tokenStr, BearerSchema))

		jwtAuth := NewJWTAuth(config.Get().SecretKey, config.Get().ExpiredTime)
		claims, err := jwtAuth.GetUserFromToken(tokenString)
		if err != nil {
			// token 无效
			Failure(w, Map{"code": 403, "message": "token无效"})
			return
		}

		// 将用户放入 context
		user := &model.PbUser{
			ID:   claims.UserID,
			Name: claims.Username,
			Role: claims.Role,
		}
		ctx := context.WithValue(r.Context(), ContextUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// 定义一个私有的类型，避免与其他包冲突
type ctxKeyToken struct{}

// 包内唯一的 key
var tokenKey = ctxKeyToken{}

// WithToken 将 token 存入 context
func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

// GetToken 从 context 中取出 token
func GetToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(tokenKey).(string)
	return token, ok
}

// AuthWsMw WebSocket 认证中间件
func AuthWsMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从 Sec-WebSocket-Protocol 取 token
		protocols := websocket.Subprotocols(r)
		if len(protocols) == 0 {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		tokenString := protocols[0]

		// 校验 token
		token := NewJWTAuth(config.Get().SecretKey, config.Get().ExpiredTime)
		if !token.ValidateToken(tokenString) {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		// slog.Info("校验ws token成功", slog.String("Token", tokenString))
		// 将 token 存入 context
		ctx := WithToken(r.Context(), tokenString)
		// 继续执行后续 handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CorsMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
		w.Header().Set("Content-Type", "application/json")
	})
}

// // InMw request incoming to do something
// func InMw(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("请求....", r.URL.Path)
// 		next.ServeHTTP(w, r)
// 	})
// }

// // OutMw Response outing to do something
// func OutMw(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		next.ServeHTTP(w, r)
// 		w.Header().Set("xxx", "0000")
// 		fmt.Println("响应...", w.Header())
// 	})
// }
