package config

import (
	"album-admin/model"
	"album-admin/utils"
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	Conf *viper.Viper
	// SiteSettings 存储从数据库加载的键值对配置
	SiteSettings map[string]string
	// settingsLoaded 确保站点配置只从数据库加载一次
	settingsLoaded sync.Once
)

func LoadConfig() {
	Conf = viper.New()
	Conf.SetConfigFile(".env")
	Conf.AutomaticEnv()
	if err := Conf.ReadInConfig(); err != nil {
		log.Printf("Warning: Failed to read .env config, using environment variables or defaults: %v", err)
	}
}

// LoadSettingsFromDB 从数据库加载所有键值对配置到 SiteSettings
func LoadSettingsFromDB(db *gorm.DB) error {
	var err error
	settingsLoaded.Do(func() {
		if db == nil {
			err = fmt.Errorf("database DB instance is nil when loading settings")
			return
		}

		var settings []model.Setting
		if findErr := db.Find(&settings).Error; findErr != nil {
			err = fmt.Errorf("failed to load settings from database: %w", findErr)
			return
		}

		SiteSettings = make(map[string]string)
		for _, s := range settings {
			SiteSettings[s.ConfigKey] = s.Value
		}
		log.Println("All site settings loaded from database successfully.")
	})
	return err
}

// GetSetting 获取单个配置值
func GetSetting(key string) string {
	if SiteSettings == nil {
		log.Println("Warning: SiteSettings not loaded. Returning empty string for key:", key)
		return ""
	}
	return SiteSettings[key]
}

// InitSettings 检查并初始化默认配置项和第一个管理员用户
func InitSettings(db *gorm.DB) {
	log.Println("--- 开始初始化站点配置 (Setting 表) ---")
	defaultSettings := map[string]struct {
		Value   string
		Comment string
	}{
		"JWT_SECRET":            {Value: "", Comment: "JWT密钥，首次启动自动生成"},
		"ABOUT_LINK":            {Value: "https://github.com/anzhiyu-c/yuyu-album", Comment: "关于我们链接"},
		"APP_NAME":              {Value: "鱼鱼相册", Comment: "应用名称"},
		"APP_VERSION":           {Value: "1.0.0", Comment: "应用版本"},
		"ICP_NUMBER":            {Value: "湘ICP备2023015794号-2", Comment: "ICP备案号"},
		"USER_AVATAR":           {Value: "/static/img/avatar.jpg", Comment: "用户默认头像URL"},
		"API_URL":               {Value: "https://album.anheyu.com/", Comment: "API地址"},
		"LOGO_URL":              {Value: "/static/img/logo.svg", Comment: "Logo图片URL (通用)"},
		"LOGO_URL_192x192":      {Value: "/static/img/logo-192x192.png", Comment: "Logo图片URL (192x192)"},
		"LOGO_URL_512x512":      {Value: "/static/img/logo-512x512.svg", Comment: "Logo图片URL (512x512)"},
		"LOGO_HORIZONTAL_DAY":   {Value: "/static/img/logo-horizontal-day.png", Comment: "横向Logo (白天模式)"},
		"LOGO_HORIZONTAL_NIGHT": {Value: "/static/img/logo-horizontal-night.png", Comment: "横向Logo (暗色模式)"},
		"ICON_URL":              {Value: "/static/img/logo.svg", Comment: "Icon图标URL"},
		"SITE_KEYWORDS":         {Value: "鱼鱼相册,相册,图片管理", Comment: "站点关键词"},
		"SITE_DESCRIPTION":      {Value: "鱼鱼相册是一个简单易用的图片管理系统，支持多种图片格式和处理方式。", Comment: "站点描述"},
		"DEFAULT_THUMB_PARAM":   {Value: "", Comment: "默认缩略图处理参数"},
		"DEFAULT_BIG_PARAM":     {Value: "", Comment: "默认大图处理参数"},
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
		log.Printf("检查配置项: %s", key)
		result := db.Where("config_key = ?", key).First(&setting)
		if result.Error == gorm.ErrRecordNotFound {
			newSetting := model.Setting{
				ConfigKey: key,
				Value:     data.Value,
				Comment:   data.Comment,
			}
			if createErr := db.Create(&newSetting).Error; createErr != nil {
				log.Printf("⚠️ 失败: 创建默认配置项 '%s' 失败: %v", key, createErr)
			} else {
				log.Printf("✅ 成功: 默认配置项 '%s' 已创建。", key)
			}
		} else if result.Error != nil {
			log.Printf("❌ 错误: 查询配置项 '%s' 时发生错误: %v", key, result.Error)
		} else {
			log.Printf("ℹ️ 跳过: 配置项 '%s' 已存在，跳过创建。", key)
		}
	}
	log.Println("--- 站点配置 (Setting 表) 初始化完成。---")

	// --- 检查 User 表并记录状态 ---
	log.Println("--- 开始检查 User 表状态 ---")
	var userCount int64
	if err := db.Model(&model.User{}).Count(&userCount).Error; err != nil {
		log.Printf("❌ 错误: 查询 User 表记录数量失败: %v。请确认 'users' 表已正确迁移。", err)
	} else {
		log.Printf("当前 User 表中用户数量: %d", userCount)
		if userCount == 0 {
			log.Println("User 表为空，第一个注册的用户将成为管理员。")
		} else {
			log.Println("User 表中已存在用户，注册将创建普通用户。")
		}
	}
	log.Println("--- User 表状态检查完成。---")
}
