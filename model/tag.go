/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-06-15 12:00:07
 * @LastEditTime: 2025-06-15 12:44:12
 * @LastEditors: 安知鱼
 */
package model

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"type:varchar(100);unique;not null;comment:标签名称" json:"name"`
}
