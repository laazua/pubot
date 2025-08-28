package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name       string   `gorm:"column:name"`
	BuildSteps []string `gorm:"column:build;type:jsonb"`
	Deoply     Deploy   `gorm:"column:deploy;type:jsonb"`
	
}

type Deploy struct {
	Platform string   `gorm:"column:platform"`
	Run      []string `gorm:"column:run;type:jsonb"`
}

// 设置task表名
func (Task) TableName() string {
	return "pb_task"
}

// 实现 Valuer 和 Scanner 接口
func (d Deploy) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *Deploy) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, d)
}

// // 创建任务
// task := Task{
//     Name:       "测试任务",
//     BuildSetps: []string{"install", "build", "test"},
//     Deploy: Deploy{
//         Platform: "linux",
//         Run:      []string{"start", "monitor", "log"},
//     },
// }

// db.Create(&task)
