/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-04-11 06:02:54
 * @LastEditTime: 2025-04-11 14:06:43
 * @LastEditors: 安知鱼
 */
package database

import (
	"fmt"
	"wallpaper-admin/config"
	"wallpaper-admin/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMySQL() {
	host := config.Conf.GetString("DB_HOST")
	port := config.Conf.GetString("DB_PORT")
	user := config.Conf.GetString("DB_USER")
	pass := config.Conf.GetString("DB_PASS")
	name := config.Conf.GetString("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MySQL: %v", err))
	}

	// 自动迁移
	if err := DB.AutoMigrate(&model.Wallpaper{}, &model.Tag{}); err != nil {
		panic(fmt.Sprintf("AutoMigrate failed: %v", err))
	}
}
