package config

import (
	"album-admin/model"
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
