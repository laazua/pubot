package dto

type Task struct {
	Id         uint     `json:"id"`
	Name       string   `json:"name"`
	BuildSetps []string `json:"buildSteps"`
	Deploy     Deploy   `json:"deploy"`
}

type Deploy struct {
	Platform string   `json:"platform"`
	Run      []string `json:"run"`
}
