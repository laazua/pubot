package db

import (
	"fmt"
	"sync"
	"time"

	"pubot/internal/core/config"
	"pubot/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db      *gorm.DB
	initErr error
	mu      sync.Mutex
)

func Init() error {
	mu.Lock()
	defer mu.Unlock()
	if db != nil && initErr == nil {
		return nil
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=disable TimeZone=Asia/Shanghai",
		config.Get().PgHost, config.Get().PgUser, config.Get().PgPass, config.Get().PgName, config.Get().PgPort)
	var err error
	//db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
	//	Logger: logger.Default.LogMode(logger.Info),
	//})
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local() // 使用本地时间
		},
	})
	if err != nil {
		initErr = fmt.Errorf("打开数据库失败: %w", err)
		return initErr
	}
	sqlDB, err := db.DB()
	if err != nil {
		initErr = fmt.Errorf("获取数据库句柄失败: %w", err)
		return initErr
	}
	sqlDB.SetMaxIdleConns(config.Get().PgMaxIdle)
	sqlDB.SetMaxOpenConns(config.Get().PgPool)
	sqlDB.SetConnMaxLifetime(config.Get().PgLifeTime)
	if err := sqlDB.Ping(); err != nil {
		initErr = fmt.Errorf("ping数据库连接失败: %w", err)
		return initErr
	}
	// 表迁移
	if err := db.AutoMigrate(&model.User{}, &model.Task{}); err != nil {
		initErr = fmt.Errorf("数据库表迁移失败: %w", err)
		return initErr
	}
	initErr = nil
	return nil
}

func Get() (*gorm.DB, error) {
	if err := Init(); err != nil {
		return nil, err
	}
	return db, nil
}

func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
