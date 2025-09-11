package dto

// LoginRequest 用户登录请求数据格式
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserRequest 用户操作请求数据格式
type UserRequest struct {
	Id       uint   `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
