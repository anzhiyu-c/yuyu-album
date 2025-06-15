/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-06-15 12:25:12
 * @LastEditTime: 2025-06-15 12:48:21
 * @LastEditors: 安知鱼
 */
package model

import (
	"time"

	"gorm.io/gorm"
)

// Setting 键值对配置模型
type Setting struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Key       string         `gorm:"type:varchar(100);unique;not null;comment:配置键名" json:"key"`
	Value     string         `gorm:"type:text;comment:配置值" json:"value"`            // 使用 text 类型存储，可以存储较长的值
	Comment   string         `gorm:"type:varchar(255);comment:配置说明" json:"comment"` // 可选：用于记录配置的用途
}
