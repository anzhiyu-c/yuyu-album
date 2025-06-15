/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-06-15 12:04:20
 * @LastEditTime: 2025-06-15 13:03:48
 * @LastEditors: 安知鱼
 */
package migrate

import (
	"album-admin/model"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// AutoMigrateTables 用于执行数据库表结构的自动迁移
// 所有的模型都需要在这里注册，以便GORM能够识别并创建/更新对应的表
func AutoMigrateTables(db *gorm.DB) {
	if db == nil {
		log.Fatal("数据库连接未初始化，无法执行迁移。")
	}

	// 在这里列出所有需要进行自动迁移的模型
	// GORM会根据这些结构体定义来创建或更新数据库表
	err := db.AutoMigrate(
		&model.Album{},
		&model.Tag{},
		&model.Setting{},
		&model.User{},
	)
	if err != nil {
		log.Fatalf("数据库自动迁移失败: %v", err)
	}
	fmt.Println("数据库迁移成功完成！")
}
