/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-04-11 06:02:54
 * @LastEditTime: 2025-04-13 00:37:01
 * @LastEditors: 安知鱼
 */
package model

import (
	"time"

	"gorm.io/gorm"
)

type Wallpaper struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ImageUrl      string         `json:"imageUrl"`      // 壁纸图
	BigImageUrl   string         `json:"bigImageUrl"`   // 大图（可选）
	DownloadUrl   string         `json:"downloadUrl"`   // 下载地址（可选）
	ThumbParam    string         `json:"thumbParam"`    // 缩略参数（JSON 字符串）
	BigParam      string         `json:"bigParam"`      // 大图参数（JSON 字符串）
	Tags          string         `json:"tags"`          // 标签，多个以英文逗号分隔
	ViewCount     int            `json:"viewCount"`     // 查看次数
	DownloadCount int            `json:"downloadCount"` // 下载次数
	Width         int            `json:"width"`         // ✅ 添加宽度
	Height        int            `json:"height"`        // ✅ 添加高度
	FileSize      int64          `json:"fileSize"`      // 文件大小，单位字节
	Format        string         `json:"format"`        // 图片格式（如 jpg、png、webp、avif）
	AspectRatio   string         `json:"aspectRatio"`   // 宽高比，例如 "16:9"
	FileHash      string         `json:"fileHash"`      // 图片的哈希值，用来判断是否重复
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Tag struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;size:50"`
	CreatedAt time.Time
}
