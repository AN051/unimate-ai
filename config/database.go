package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	dsn := "root:123456@tcp(localhost:3306)/unimate_ai?charset=utf8mb4&parseTime=True&loc=Local"
	
	log.Println("连接数据库...")
	
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})

	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("数据库连接成功")
	return nil
}

func AutoMigrate(models ...interface{}) error {
	log.Println("创建数据表...")
	if err := DB.AutoMigrate(models...); err != nil {
		return err
	}
	log.Println("数据表创建完成")
	return nil
}