package dao

import (
	"fmt"
	"log/slog"

	"sync"
	"time"

	"pubot/internal/config"
	"pubot/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	pgDb *gorm.DB
	once sync.Once
)

func initPg() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=disable TimeZone=Asia/Shanghai",
		config.Get().PgHost, config.Get().PgUser, config.Get().PgPass, config.Get().PgName, config.Get().PgPort)
	//pgDb, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
	//	Logger: logger.Default.LogMode(logger.Info),
	//})
	var err error
	pgDb, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local() // 使用本地时间
		},
	})
	if err != nil {
		slog.Error("打开数据库失败", slog.String("Err", err.Error()))
		return err
	}
	sqlDB, err := pgDb.DB()
	if err != nil {
		slog.Error("获取数据库句柄失败", slog.String("Err", err.Error()))
		return err
	}
	sqlDB.SetMaxIdleConns(config.Get().PgMaxIdle)
	sqlDB.SetMaxOpenConns(config.Get().PgPool)
	sqlDB.SetConnMaxLifetime(config.Get().PgLifeTime)
	if err := sqlDB.Ping(); err != nil {
		slog.Error("ping 数据库连接失败", slog.String("Err", err.Error()))
		return err
	}
	// 表迁移
	if err := pgDb.AutoMigrate(&model.PbUser{}, &model.PbTask{}); err != nil {
		slog.Error("数据库表迁移失败", slog.String("Err", err.Error()))
		return err
	}
	return nil
}

func GetDb() *gorm.DB {
	if pgDb == nil {
		once.Do(func() {
			if err := initPg(); err != nil {
				slog.Error("初始化数据库失败", slog.String("Err", err.Error()))
				panic(err)
			}
		})
	}
	return pgDb
}

func CloseDb() error {
	if pgDb != nil {
		sqlDB, err := pgDb.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
