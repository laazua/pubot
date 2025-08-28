package dto

type ReqUser struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Emal     string `json:"emal"`
	Password string `json:"password"`
}

type RespUser struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Emal string `json:"emal"`
}

type AuthUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
