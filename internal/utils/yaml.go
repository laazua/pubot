package utils

import (
	"pubot/internal/dto"

	"gopkg.in/yaml.v3"
)

func ParseTaskYAML(yamlText string) (*dto.TaskYAML, error) {
	var parsed dto.TaskYAML
	err := yaml.Unmarshal([]byte(yamlText), &parsed)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
