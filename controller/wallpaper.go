/*
 * @Description: 壁纸控制器
 * @Author: 安知鱼
 * @Date: 2025-04-11 15:17:13
 * @LastEditTime: 2025-04-13 01:16:09
 * @LastEditors: 安知鱼
 */

package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"wallpaper-admin/database"
	"wallpaper-admin/model"

	"github.com/gin-gonic/gin"
)

// 壁纸列表查询（支持分页、标签和创建时间筛选）
func GetWallpapers(c *gin.Context) {

	type WallpaperResponse struct {
		ID             uint      `json:"id"`
		ImageUrl       string    `json:"imageUrl"`
		BigImageUrl    string    `json:"bigImageUrl"`
		DownloadUrl    string    `json:"downloadUrl"`
		ThumbParam     string    `json:"thumbParam"`
		BigParam       string    `json:"bigParam"`
		Tags           string    `json:"tags"`
		ViewCount      int       `json:"viewCount"`
		DownloadCount  int       `json:"downloadCount"`
		FileSize       int64     `json:"fileSize"`
		Format         string    `json:"format"`
		AspectRatio    string    `json:"aspectRatio"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
		Width          int       `json:"width"`
		Height         int       `json:"height"`
		WidthAndHeight string    `json:"widthAndHeight"`
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	tag := c.Query("tag")
	startStr := c.Query("createdAt[0]")
	endStr := c.Query("createdAt[1]")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	offset := (page - 1) * pageSize

	const layout = "2006/01/02 15:04:05" // 匹配前端 value-format

	var rawWallpapers []model.Wallpaper
	var wallpapers []WallpaperResponse
	query := database.DB

	if tag != "" {
		query = query.Where("FIND_IN_SET(?, tags)", tag)
	}

	if startStr != "" && endStr != "" {
		startTime, err1 := time.ParseInLocation(layout, startStr, time.Local)
		endTime, err2 := time.ParseInLocation(layout, endStr, time.Local)

		if err1 == nil && err2 == nil {
			query = query.Where("created_at BETWEEN ? AND ?", startTime, endTime)
		} else {
			// 输出错误日志便于调试
			if err1 != nil {
				fmt.Printf("start time parse error: %v\n", err1)
			}
			if err2 != nil {
				fmt.Printf("end time parse error: %v\n", err2)
			}
		}
	}

	var total int64
	query.Model(&model.Wallpaper{}).Count(&total)

	query.Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&rawWallpapers)

	for _, w := range rawWallpapers {
		wallpapers = append(wallpapers, WallpaperResponse{
			ID:             w.ID,
			ImageUrl:       w.ImageUrl,
			BigImageUrl:    w.BigImageUrl,
			DownloadUrl:    w.DownloadUrl,
			ThumbParam:     w.ThumbParam,
			BigParam:       w.BigParam,
			Tags:           w.Tags,
			ViewCount:      w.ViewCount,
			DownloadCount:  w.DownloadCount,
			CreatedAt:      w.CreatedAt,
			UpdatedAt:      w.UpdatedAt,
			FileSize:       w.FileSize,
			Format:         w.Format,
			AspectRatio:    w.AspectRatio,
			Width:          w.Width,
			Height:         w.Height,
			WidthAndHeight: fmt.Sprintf("%dx%d", w.Width, w.Height),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"list":     wallpapers,
			"total":    total,
			"pageNum":  page,
			"pageSize": pageSize,
		},
	})
}

// 新增壁纸
func AddWallpaper(c *gin.Context) {
	// 临时接收结构体
	var req struct {
		ImageUrl      string   `json:"imageUrl"`
		BigImageUrl   string   `json:"bigImageUrl"`
		DownloadUrl   string   `json:"downloadUrl"`
		ThumbParam    string   `json:"thumbParam"`
		BigParam      string   `json:"bigParam"`
		Tags          []string `json:"tags"`
		ViewCount     int      `json:"viewCount"`
		DownloadCount int      `json:"downloadCount"`
		Width         int      `json:"width"`
		Height        int      `json:"height"`
		FileSize      int64    `json:"fileSize"`
		Format        string   `json:"format"`
		AspectRatio   string   `json:"aspectRatio"`
		FileHash      string   `json:"fileHash"` // 图片的哈希值，用来判断是否重复
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "参数错误"})
		return
	}

	// 检查是否存在相同哈希值的图片，并返回重复图片的信息
	var existingWallpaper model.Wallpaper
	if err := database.DB.Where("file_hash = ?", req.FileHash).First(&existingWallpaper).Error; err == nil {
		// 如果数据库中有相同的文件哈希值，说明是重复的图片
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "这张图片已存在，id是" + fmt.Sprint(existingWallpaper.ID) + "，请勿重复添加",
			"duplicate": gin.H{
				"id":          existingWallpaper.ID,          // 重复图片的ID
				"imageUrl":    existingWallpaper.ImageUrl,    // 重复图片的URL
				"fileSize":    existingWallpaper.FileSize,    // 重复图片的大小
				"format":      existingWallpaper.Format,      // 重复图片的格式
				"aspectRatio": existingWallpaper.AspectRatio, // 重复图片的宽高比
				"fileHash":    existingWallpaper.FileHash,    // 重复图片的哈希值
			},
		})
		return
	}

	width := req.Width
	height := req.Height

	if width <= 0 || height <= 0 {
		width, height = 0, 0
	}

	// 默认填充
	if req.BigImageUrl == "" {
		req.BigImageUrl = req.ImageUrl
	}
	if req.DownloadUrl == "" {
		req.DownloadUrl = req.ImageUrl
	}
	if req.ThumbParam == "" {
		req.ThumbParam = "x-oss-process=image//resize,h_600/quality,q_100/auto-orient,0/interlace,1/format,avif"
	}
	if req.BigParam == "" {
		req.BigParam = "x-oss-process=image//resize,s_2000/quality,q_100/auto-orient,0/interlace,1/format,avif"
	}

	// 将标签数组转为英文逗号分隔的字符串
	tagsStr := ""
	if len(req.Tags) > 0 {
		tagsStr = joinTags(req.Tags)
	}

	// 保存壁纸
	wallpaper := model.Wallpaper{
		ImageUrl:      req.ImageUrl,
		BigImageUrl:   req.BigImageUrl,
		DownloadUrl:   req.DownloadUrl,
		ThumbParam:    req.ThumbParam,
		BigParam:      req.BigParam,
		Tags:          tagsStr,
		ViewCount:     req.ViewCount,
		DownloadCount: req.DownloadCount,
		Width:         req.Width,
		Height:        req.Height,
		FileSize:      req.FileSize,
		Format:        req.Format,
		AspectRatio:   req.AspectRatio,
		FileHash:      req.FileHash, // 将文件哈希值保存到数据库
	}

	if err := database.DB.Create(&wallpaper).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "添加失败", "code": 500})
		return
	}

	// 插入新标签（避免重复）
	for _, tag := range req.Tags {
		var count int64
		database.DB.Model(&model.Tag{}).Where("name = ?", tag).Count(&count)
		if count == 0 {
			database.DB.Create(&model.Tag{Name: tag})
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "添加成功", "code": 200})
}

// 删除壁纸
func DeleteWallpaper(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID非法", "code": 400})
		return
	}

	if err := database.DB.Delete(&model.Wallpaper{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败", "code": 500})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "删除成功", "code": 200})
}

func UpdateWallpaper(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID非法"})
		return
	}

	var req struct {
		ImageUrl      string   `json:"imageUrl"`
		BigImageUrl   string   `json:"bigImageUrl"`
		DownloadUrl   string   `json:"downloadUrl"`
		ThumbParam    string   `json:"thumbParam"`
		BigParam      string   `json:"bigParam"`
		Tags          []string `json:"tags"`
		ViewCount     int      `json:"viewCount"`
		DownloadCount int      `json:"downloadCount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "参数错误"})
		return
	}

	// 自动填充空字段
	if req.BigImageUrl == "" {
		req.BigImageUrl = req.ImageUrl
	}
	if req.DownloadUrl == "" {
		req.DownloadUrl = req.ImageUrl
	}

	tagsStr := joinTags(req.Tags)

	updateData := map[string]interface{}{
		"image_url":      req.ImageUrl,
		"big_image_url":  req.BigImageUrl,
		"download_url":   req.DownloadUrl,
		"thumb_param":    req.ThumbParam,
		"big_param":      req.BigParam,
		"tags":           tagsStr,
		"view_count":     req.ViewCount,
		"download_count": req.DownloadCount,
	}

	if err := database.DB.Model(&model.Wallpaper{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败", "code": 500})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "更新成功", "code": 200})
}
