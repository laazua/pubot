package config

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	config *Config
	once   sync.Once
)

type Config struct {
	Listen      string        `yaml:"listen" default:"127.0.0.1:7777"` // 监听地址
	SecretKey   string        `yaml:"secretKey" default:"1adnfdjkfa"`  // token key
	ExpiredTime time.Duration `yaml:"expiredTime" default:"12h"`       // token过期时间
	PgHost      string        `yaml:"pgHost" default:"127.0.0.1"`      // pgdb主机
	WorkSpace   string        `yaml:"workSpace" default:"."`           // 工作目录
	PgPort      int           `yaml:"pgPort" default:"5432"`           // pgdb端口
	PgUser      string        `yaml:"pgUser"`                          // pgdb认证用户
	PgPass      string        `yaml:"pgPass"`                          // pgdb认证密码
	PgName      string        `yaml:"pgName"`                          // pgdb数据库名
	PgPool      int           `yaml:"pgPool" default:"20"`             // pgdb池大小
	PgMaxIdle   int           `yaml:"pgMaxIdle" default:"50"`          // pgdb idle大小
	PgLifeTime  time.Duration `yaml:"pgLifeTime" default:"1h30m"`      // pgdb lifetime时间
}

func initConfig() error {
	fileName := filepath.Join(".", "config.yaml")
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return err
	}
	return nil
}

func Get() *Config {
	if config == nil {
		once.Do(func() {
			if err := initConfig(); err != nil {
				panic(err)
			}
		})
	}
	return config

}
