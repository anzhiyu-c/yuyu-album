/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-06-15 11:30:55
 * @LastEditTime: 2025-06-15 12:18:41
 * @LastEditors: 安知鱼
 */
package controller

import (
	"album-admin/database"
	"album-admin/model"
	"album-admin/utils/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetPublicAlbums 获取公开的相册列表
func GetPublicAlbums(c *gin.Context) {
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

	var list []model.Album
	var total int64

	db.Model(&model.Album{}).Count(&total)

	db.Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&list)

	// 使用统一返回体
	response.Success(c, gin.H{
		"list":     list,
		"total":    total,
		"pageNum":  page,
		"pageSize": pageSize,
	}, "获取相册列表成功")
}

// UpdateAlbumStat 更新访问量或下载量
func UpdateAlbumStat(c *gin.Context) {
	idStr := c.Param("id")
	statType := c.Query("type") // "view" 或 "download"

	id, err := strconv.Atoi(idStr)
	if err != nil {
		// 使用统一返回体
		response.Fail(c, http.StatusBadRequest, "无效的ID")
		return
	}

	var album model.Album
	if err := database.DB.First(&album, id).Error; err != nil {
		// 使用统一返回体
		response.Fail(c, http.StatusNotFound, "相册未找到")
		return
	}

	switch statType {
	case "view":
		database.DB.Model(&album).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
	case "download":
		database.DB.Model(&album).UpdateColumn("download_count", gorm.Expr("download_count + ?", 1))
	default:
		// 使用统一返回体
		response.Fail(c, http.StatusBadRequest, "无效的type参数")
		return
	}

	// 使用统一返回体
	response.Success(c, nil, "更新成功")
}
