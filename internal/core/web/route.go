package web

import (
	"context"
	"net/http"
	"strings"

	"pubot/internal/core/logx"
)

// contextKey 用于避免 context key 冲突
type contextKey string

// Middleware 定义中间件类型
type Middleware func(http.Handler) http.Handler

// Router 路由器实现
type Router struct {
	routes      map[string]map[string]route // 方法 -> 路径 -> 路由信息
	middlewares []Middleware                // 当前 router 累积的中间件（全局 + 分组）
	groups      []group                     // 路由分组（仅用于管理，不影响逻辑）
	prefix      string                      // 路由前缀（用于分组）
}

type route struct {
	pattern     string       // 路由模式
	handler     http.Handler // 处理器
	params      []string     // 路由参数名
	middlewares []Middleware // 路由级别中间件（包含全局 + 分组 + 路由级别）
}

type group struct {
	prefix      string
	middlewares []Middleware
}

// NewRouter 创建一个新的路由器
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]route),
	}
}

// Use 注册全局中间件
func (r *Router) Use(mw Middleware) {
	r.middlewares = append(r.middlewares, mw)
}

// Handle 注册路由（支持路由级中间件 + 分组中间件）
func (r *Router) Handle(method, pattern string, handler http.Handler, mws ...Middleware) {
	if _, exists := r.routes[method]; !exists {
		r.routes[method] = make(map[string]route)
	}

	fullPattern := r.prefix + pattern
	parts, params := parsePattern(fullPattern)

	// 合并：全局 + 分组中间件（在 r.middlewares） + 路由级别
	allMws := append([]Middleware{}, r.middlewares...)
	allMws = append(allMws, mws...)

	r.routes[method][strings.Join(parts, "/")] = route{
		pattern:     fullPattern,
		handler:     handler,
		params:      params,
		middlewares: allMws,
	}
	logx.Debug("注册的路由", logx.String("Method", method), logx.String("Pattern", fullPattern))
}

// Group 创建一个路由分组
func (r *Router) Group(prefix string, mws ...Middleware) *Router {
	return &Router{
		routes:      r.routes,
		middlewares: append(r.middlewares, mws...), // 累积父级中间件
		groups:      append(r.groups, group{prefix: prefix, middlewares: mws}),
		prefix:      r.prefix + prefix, // 累加前缀
	}
}

// Get 注册 GET 请求路由
func (r *Router) Get(pattern string, handlerFun http.HandlerFunc, mws ...Middleware) {
	r.Handle(http.MethodGet, pattern, handlerFun, mws...)
}

// Post 注册 POST 请求路由
func (r *Router) Post(pattern string, handlerFun http.HandlerFunc, mws ...Middleware) {
	r.Handle(http.MethodPost, pattern, handlerFun, mws...)
}

// Delete 注册 DELETE 请求路由
func (r *Router) Delete(pattern string, handlerFun http.HandlerFunc, mws ...Middleware) {
	r.Handle(http.MethodDelete, pattern, handlerFun, mws...)
}

// Put 注册 PUT 请求路由
func (r *Router) Put(pattern string, handlerFun http.HandlerFunc, mws ...Middleware) {
	r.Handle(http.MethodPut, pattern, handlerFun, mws...)
}

// ServeHTTP 实现 http.Handler 接口
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path

	handler, params, routeMws := r.matchRoute(method, path)
	if handler == nil {
		http.NotFound(w, req)
		return
	}

	// 应用所有中间件（已经包含全局 + 分组 + 路由级别）
	for _, mw := range routeMws {
		handler = mw(handler)
	}

	req = addParamsToContext(req, params)
	handler.ServeHTTP(w, req)
}

// matchRoute 匹配路由
func (r *Router) matchRoute(method, path string) (http.Handler, map[string]string, []Middleware) {
	if routes, ok := r.routes[method]; ok {
		parts := splitPath(path)
		for _, rt := range routes {
			if params := matchPattern(rt.pattern, parts); params != nil {
				return rt.handler, params, rt.middlewares
			}
		}
	}
	return nil, nil, nil
}

// parsePattern 解析路径模式
func parsePattern(pattern string) ([]string, []string) {
	parts := splitPath(pattern)
	var params []string
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			params = append(params, part[1:])
			parts[i] = "*"
		}
	}
	return parts, params
}

// splitPath 分割路径
func splitPath(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return []string{}
	}
	return strings.Split(trimmed, "/")
}

// matchPattern 匹配路径
func matchPattern(pattern string, pathParts []string) map[string]string {
	patternParts := splitPath(pattern)
	if len(patternParts) != len(pathParts) {
		return nil
	}
	params := make(map[string]string)
	for i, part := range patternParts {
		if part == "*" {
			continue
		}
		if strings.HasPrefix(part, ":") {
			params[part[1:]] = pathParts[i]
		} else if part != pathParts[i] {
			return nil
		}
	}
	return params
}

// 将路由参数添加到请求上下文
func addParamsToContext(req *http.Request, params map[string]string) *http.Request {
	ctx := req.Context()
	for key, value := range params {
		ctx = context.WithValue(ctx, contextKey(key), value)
	}
	return req.WithContext(ctx)
}

// Param 从请求上下文获取路由参数
func Param(req *http.Request, key string) string {
	if val, ok := req.Context().Value(contextKey(key)).(string); ok {
		return val
	}
	return ""
}
