/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-06-15 13:02:05
 * @LastEditTime: 2025-06-15 13:02:11
 * @LastEditors: 安知鱼
 */
package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Username     string         `gorm:"type:varchar(50);unique;not null;comment:用户账号" json:"username"`
	PasswordHash string         `gorm:"type:varchar(255);not null;comment:密码哈希值" json:"-"` // 密码哈希，不对外暴露
	Nickname     string         `gorm:"type:varchar(50);comment:用户昵称" json:"nickname"`
	Avatar       string         `gorm:"type:varchar(255);comment:用户头像URL" json:"avatar"`
	Email        string         `gorm:"type:varchar(100);comment:用户邮箱" json:"email"`
	Role         string         `gorm:"type:varchar(20);default:'admin';comment:用户角色" json:"role"` // 例如：admin, editor, normal
	Status       int            `gorm:"type:tinyint;default:1;comment:用户状态 1:正常 0:禁用" json:"status"`
	LastLoginAt  *time.Time     `gorm:"comment:最后登录时间" json:"lastLoginAt"`
}
