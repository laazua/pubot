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
	middlewares []Middleware                // 全局中间件
	groups      []group                     // 路由分组
	prefix      string                      // 路由前缀（用于分组）
}

type route struct {
	pattern string       // 路由模式
	handler http.Handler // 处理器
	params  []string     // 路由参数名
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

// Handle 注册路由
func (r *Router) Handle(method, pattern string, handler http.Handler) {
	if _, exists := r.routes[method]; !exists {
		r.routes[method] = make(map[string]route)
	}

	// 拼接分组前缀
	fullPattern := r.prefix + pattern

	parts, params := parsePattern(fullPattern)
	r.routes[method][strings.Join(parts, "/")] = route{
		pattern: fullPattern,
		handler: handler,
		params:  params,
	}
	logx.Debug("注册的路由", logx.String("Method", method), logx.String("Prefix", fullPattern))
}

// Group 创建一个路由分组
func (r *Router) Group(prefix string, mws ...Middleware) *Router {
	return &Router{
		routes:      r.routes,
		middlewares: append(r.middlewares, mws...),
		groups:      append(r.groups, group{prefix: prefix, middlewares: mws}),
		prefix:      r.prefix + prefix, // 递归累加前缀
	}
}

// Get 注册 GET 请求路由
func (r *Router) Get(pattern string, handlerFun http.HandlerFunc) {
	r.Handle(http.MethodGet, pattern, handlerFun)
}

// Post 注册 POST 请求路由
func (r *Router) Post(pattern string, handlerFun http.HandlerFunc) {
	r.Handle(http.MethodPost, pattern, handlerFun)
}

// Delete 注册 DELETE 请求路由
func (r *Router) Delete(pattern string, handlerFun http.HandlerFunc) {
	r.Handle(http.MethodDelete, pattern, handlerFun)
}

// Put 注册 PUT 请求路由
func (r *Router) Put(pattern string, handlerFun http.HandlerFunc) {
	r.Handle(http.MethodPut, pattern, handlerFun)
}

// ServeHTTP 实现 http.Handler 接口
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path

	var handler http.Handler
	var params map[string]string

	// 匹配路由
	handler, params = r.matchRoute(method, path)

	// 如果没有找到匹配的 handler
	if handler == nil {
		http.NotFound(w, req)
		return
	}

	// 应用全局中间件
	for _, mw := range r.middlewares {
		handler = mw(handler)
	}

	if handler != nil {
		// 获取查询参数并将其添加到上下文
		req = addParamsToContext(req, params)
		handler.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}

// matchRoute 函数
func (r *Router) matchRoute(method, path string) (http.Handler, map[string]string) {
	if routes, ok := r.routes[method]; ok {
		parts := splitPath(path)
		for _, rt := range routes {
			if params := matchPattern(rt.pattern, parts); params != nil {
				return rt.handler, params
			}
		}
	}
	return nil, nil
}

// 辅助函数：解析路径模式
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

// 分割路径
func splitPath(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return []string{}
	}
	return strings.Split(trimmed, "/")
}

// matchPattern 函数
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
		if strings.HasPrefix(part, ":") { // 如果是动态参数
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
