package database

import (
	"album-admin/config" // 导入 config 包
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMySQL() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn, // 设置日志级别为 Warn 或 Info，以便看到更多信息
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	dbUser := config.Conf.GetString("DB_USER")
	dbPass := config.Conf.GetString("DB_PASS")
	dbHost := config.Conf.GetString("DB_HOST")
	dbPortStr := config.Conf.GetString("DB_PORT")
	dbName := config.Conf.GetString("DB_NAME")

	if dbUser == "" || dbPass == "" || dbHost == "" || dbPortStr == "" || dbName == "" {
		log.Fatalf("Missing one or more database connection parameters (DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME) in .env or environment variables.")
	}

	_, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %s. Must be a number.", dbPortStr)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPortStr, dbName)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL with DSN '%s': %v", dsn, err)
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("MySQL 数据库连接成功！")
}
