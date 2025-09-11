package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type PbTask struct {
	ID         uint            `gorm:"primaryKey;autoIncrement"`
	Name       string          `gorm:"type:varchar(255);not null"`
	YAML       string          `gorm:"type:text"`
	YAMLParsed json.RawMessage `gorm:"type:jsonb"` // 存储解析后的 JSON
	Status     string          `gorm:"type:varchar(20);default:'stopped'"`
	Count      int             `gorm:"default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (PbTask) TableName() string {
	return "pb_task"
}
