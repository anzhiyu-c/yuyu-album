/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-04-11 18:48:11
 * @LastEditTime: 2025-04-12 11:54:21
 * @LastEditors: 安知鱼
 */
package controller

import (
	"net/http"
	"strconv"
	"time"
	"wallpaper-admin/database"
	"wallpaper-admin/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPublicWallpapers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "12")
	tag := c.Query("tag")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	offset := (page - 1) * pageSize

	db := database.DB

	// 标签过滤
	if tag != "" {
		db = db.Where("FIND_IN_SET(?, tags)", tag)
	}

	// 创建时间范围过滤
	startStr := c.Query("createdAt[0]")
	endStr := c.Query("createdAt[1]")
	const layout = "2006/01/02 15:04:05"

	if startStr != "" && endStr != "" {
		startTime, err1 := time.ParseInLocation(layout, startStr, time.Local)
		endTime, err2 := time.ParseInLocation(layout, endStr, time.Local)

		if err1 == nil && err2 == nil {
			db = db.Where("created_at BETWEEN ? AND ?", startTime, endTime)
		}
	}

	var list []model.Wallpaper
	var total int64

	db.Model(&model.Wallpaper{}).Count(&total)

	db.Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&list)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"list":     list,
			"total":    total,
			"pageNum":  page,
			"pageSize": pageSize,
		},
		"code": 200,
	})
}

// 更新访问量或下载量
func UpdateWallpaperStat(c *gin.Context) {
	idStr := c.Param("id")
	statType := c.Query("type") // "view" 或 "download"

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的ID"})
		return
	}

	var wallpaper model.Wallpaper
	if err := database.DB.First(&wallpaper, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "壁纸未找到"})
		return
	}

	switch statType {
	case "view":
		database.DB.Model(&wallpaper).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
	case "download":
		database.DB.Model(&wallpaper).UpdateColumn("download_count", gorm.Expr("download_count + ?", 1))
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的type参数"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "更新成功"})
}
