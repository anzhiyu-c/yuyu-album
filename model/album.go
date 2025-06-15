/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-04-11 06:02:54
 * @LastEditTime: 2025-06-15 12:43:59
 * @LastEditors: 安知鱼
 */
package model

import (
	"time"

	"gorm.io/gorm"
)

type Album struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	ImageUrl      string         `gorm:"type:varchar(255);not null;comment:图片URL" json:"imageUrl"`
	BigImageUrl   string         `gorm:"type:varchar(255);comment:大图URL" json:"bigImageUrl"`
	DownloadUrl   string         `gorm:"type:varchar(255);comment:下载URL" json:"downloadUrl"`
	ThumbParam    string         `gorm:"type:varchar(512);comment:缩略图处理参数" json:"thumbParam"`
	BigParam      string         `gorm:"type:varchar(512);comment:大图处理参数" json:"bigParam"`
	Tags          string         `gorm:"type:varchar(255);comment:标签，逗号分隔" json:"tags"`
	ViewCount     int            `gorm:"default:0;comment:查看次数" json:"viewCount"`
	DownloadCount int            `gorm:"default:0;comment:下载次数" json:"downloadCount"`
	Width         int            `gorm:"comment:图片宽度" json:"width"`
	Height        int            `gorm:"comment:图片高度" json:"height"`
	FileSize      int64          `gorm:"comment:文件大小（字节）" json:"fileSize"`
	Format        string         `gorm:"type:varchar(50);comment:图片格式" json:"format"`
	AspectRatio   string         `gorm:"type:varchar(50);comment:图片宽高比" json:"aspectRatio"`
	FileHash      string         `gorm:"type:varchar(64);unique;comment:文件哈希值" json:"fileHash"`
}
