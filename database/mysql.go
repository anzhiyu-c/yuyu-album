package database

import (
	"album-admin/config"
	"album-admin/model"
	"album-admin/utils"
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
			LogLevel:                  logger.Warn,
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

// InitSettings 检查并初始化默认配置项
func InitSettings(db *gorm.DB) {
	// 定义所有默认配置，这里使用 Map 更方便管理
	defaultSettings := map[string]struct {
		Value   string
		Comment string
	}{
		"JWT_SECRET":          {Value: "", Comment: "JWT密钥，首次启动自动生成"},
		"ABOUT_LINK":          {Value: "https://github.com/anzhiyu-c/yuyu-album", Comment: "关于我们链接"},
		"APP_NAME":            {Value: "鱼鱼相册", Comment: "应用名称"},
		"APP_VERSION":         {Value: "1.0.0", Comment: "应用版本"},
		"ICP_NUMBER":          {Value: "湘ICP备2023015794号-2", Comment: "ICP备案号"},
		"USER_AVATAR":         {Value: "https://npm.elemecdn.com/anzhiyu-blog-static@1.0.4/img/avatar.jpg", Comment: "用户默认头像URL"},
		"API_URL":             {Value: "https://album.anheyu.com/", Comment: "API地址"},
		"LOGO_URL":            {Value: "https://album.anheyu.com/logo.svg", Comment: "Logo图片URL"},
		"ICON_URL":            {Value: "https://album.anheyu.com/logo.svg", Comment: "Icon图标URL"},
		"DEFAULT_THUMB_PARAM": {Value: "x-oss-process=image//resize,h_600/quality,q_100/auto-orient,0/interlace,1/format,avif", Comment: "默认缩略图处理参数"},
		"DEFAULT_BIG_PARAM":   {Value: "x-oss-process=image//resize,s_2000/quality,q_100/auto-orient,0/interlace,1/format,avif", Comment: "默认大图处理参数"},
	}

	// 自动生成 JWT Secret
	jwtSecret, err := utils.GenerateRandomString(32)
	if err != nil {
		log.Fatalf("Failed to generate JWT Secret for initial settings: %v", err)
	}
	defaultSettings["JWT_SECRET"] = struct {
		Value   string
		Comment string
	}{Value: jwtSecret, Comment: "JWT密钥，首次启动自动生成"}

	for key, data := range defaultSettings {
		var setting model.Setting
		// 检查配置是否存在
		result := db.Where("key = ?", key).First(&setting)
		if result.Error == gorm.ErrRecordNotFound {
			// 如果不存在，则创建
			newSetting := model.Setting{
				Key:     key,
				Value:   data.Value,
				Comment: data.Comment,
			}
			if createErr := db.Create(&newSetting).Error; createErr != nil {
				log.Printf("Failed to create default setting '%s': %v", key, createErr)
			} else {
				log.Printf("Default setting '%s' created.", key)
			}
		} else if result.Error != nil {
			// 其他查询错误
			log.Printf("Error checking setting '%s': %v", key, result.Error)
		} else {
			// 配置已存在，不做处理
			// log.Printf("Setting '%s' already exists, skipping.", key)
		}
	}
	log.Println("站点配置初始化完成。")
}
