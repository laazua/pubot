package dto

type Deploy struct {
	Platform string   `yaml:"platform" json:"platform"`
	Run      []string `yaml:"run" json:"run"`
}

type TaskYAML struct {
	Name   string   `yaml:"name" json:"name"`
	Build  []string `yaml:"build" json:"build"`
	Deploy Deploy   `yaml:"deploy" json:"deploy"`
}

// TaskCreateRequest 创建任务DTO
type TaskCreateRequest struct {
	Name   string    `json:"name" binding:"required"`
	YAML   string    `json:"yaml" binding:"required"` // 原始 YAML 文本
	Parsed *TaskYAML `json:"parsed,omitempty"`        // 解析后的结构体，不从 JSON 接收
}

type TaskUpdateRequest struct {
	Name   string    `json:"name" binding:"required"`
	YAML   string    `json:"yaml" binding:"required"`
	Status string    `json:"status,omitempty"`
	Parsed *TaskYAML `json:"parsed,omitempty"`
}
