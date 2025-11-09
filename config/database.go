package config

import (
	"log"
	"time"

	"github.com/glebarez/sqlite"  // ← 改用这个纯 Go 驱动
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接（SQLite 纯 Go 版本）
func InitDB() error {
	var err error
	
	// 使用 glebarez/sqlite（纯 Go 实现，不需要 CGO）
	DB, err = gorm.Open(sqlite.Open("unimate_ai.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})

	if err != nil {
		return err
	}

	log.Println("✅ 数据库连接成功（SQLite 本地文件 - 纯 Go 驱动）")
	log.Println("📁 数据文件：unimate_ai.db")
	
	return nil
}

// AutoMigrate 自动迁移数据表
func AutoMigrate(models ...interface{}) error {
	log.Println("🔄 开始创建/更新数据表...")
	if err := DB.AutoMigrate(models...); err != nil {
		return err
	}
	log.Println("✅ 数据表创建/更新完成")
	return nil
}