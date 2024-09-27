package pkg

import (
	"awesomeProject1/config"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

import (
	"gorm.io/driver/mysql"
	"os"
)

// Connect 连接到数据库并返回 *gorm.DB 实例
func Connect(dbConfig *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := dbConfig.User + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + dbConfig.Port + ")/" + dbConfig.Name + "?charset=utf8mb4&parseTime=True&loc=Local"

	// 创建一个新的 GORM 实例
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 设置默认的事务超时时间
		NowFunc: func() time.Time { return time.Now() },
		// 设置日志记录器
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // 使用标准库的Logger
			logger.Config{
				LogLevel:      logger.Info,            // 设置日志级别
				SlowThreshold: 100 * time.Millisecond, // 设置慢查询阈值
				Colorful:      true,                   // 是否彩色输出
			},
		),
	})
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB, err := db.DB()

	if err != nil {
		log.Fatalf("failed to get raw database connection: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	return db, nil
}
