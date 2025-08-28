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

// Config pubot程序配置
type Config struct {
	Address      string        `yaml:"address" default:"localhost"`         // pubot服务监听地址
	ReadTimeout  time.Duration `yaml:"readTimeout" default:"1h"`            // http读超时
	WriteTimeout time.Duration `yaml:"writeTimeout" default:"1h"`           // http写超时
	PwdSalt      string        `yaml:"pwdSalt" default:"1df15haM"`          // 密码salt
	TokenExpired time.Duration `yaml:"tokenExpired" default:"1h"`           // token过期时间
	SecretKey    string        `yaml:"secretKey" default:"nJ2mx&2nd12da2A"` // pubot密钥
	LogPath      string        `yaml:"logPath" default:"./logs"`            // 程序日志路径
	LogFormat    string        `yaml:"logFormat" default:"text"`            // 程序日志格式
	LogConsole   bool          `yaml:"logConsole" default:"true"`           // 程序日志是否输出到控制台
	LogLevel     string        `yaml:"logLevel" default:"debug"`            // 日志等级: debug|info|warn|error
	PgHost       string        `yaml:"pgHost" default:"localhost"`          // PG数据库地址
	PgPort       int           `yaml:"pgPort" default:"5432"`               // pg数据库端口
	PgName       string        `yaml:"pgName" default:"pubot"`              // pg数据库名称
	PgUser       string        `yaml:"pgUser" default:"zhangsan"`           // pg数据库用户
	PgPass       string        `yaml:"pgPass" default:"123456abc"`          // pg数据库密码
	PgPool       int           `ymal:"pgPool" default:"pgPool"`             // pg数据库连接池
	PgMaxIdle    int           `yaml:"pgMaxIdle" default:"20"`              // pg数据库Maxidle
	PgLifeTime   time.Duration `ymal:"pgLifeTime" default:"600s"`           // pg数据库LifeTime
}

func Init() error {
	fileName := filepath.Join(".", "config.yml")
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	// 从环境变量设置日志默认值
	config = &Config{}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return err
	}
	return nil
}

// Get 通过实例化Cofnig, 获取配置项
func Get() *Config {
	if config == nil {
		once.Do(func() {
			if err := Init(); err != nil {
				panic(err)
			}
		})
	}
	return config
}
