package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"album-admin/config" // 导入 config 包
	"album-admin/database"
	"album-admin/model"
	"album-admin/utils"
	"album-admin/utils/response"

	"github.com/gin-gonic/gin"
)

// 图片列表查询（支持分页、标签和创建时间筛选）
func GetAlbums(c *gin.Context) {

	type AlbumResponse struct {
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

	const layout = "2006/01/02 15:04:05"

	var rawAlbums []model.Album
	var albums []AlbumResponse
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
			if err1 != nil {
				fmt.Printf("start time parse error: %v\n", err1)
			}
			if err2 != nil {
				fmt.Printf("end time parse error: %v\n", err2)
			}
		}
	}

	var total int64
	query.Model(&model.Album{}).Count(&total)

	query.Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&rawAlbums)

	for _, w := range rawAlbums {
		albums = append(albums, AlbumResponse{
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

	response.Success(c, gin.H{
		"list":     albums,
		"total":    total,
		"pageNum":  page,
		"pageSize": pageSize,
	}, "获取图片列表成功")
}

// 新增图片
func AddAlbum(c *gin.Context) {
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
		FileHash      string   `json:"fileHash"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	var existingAlbum model.Album
	if err := database.DB.Where("file_hash = ?", req.FileHash).First(&existingAlbum).Error; err == nil {
		response.Fail(c, http.StatusOK, "这张图片已存在，id是"+fmt.Sprint(existingAlbum.ID)+"，请勿重复添加")
		return
	}

	width := req.Width
	height := req.Height

	if width <= 0 || height <= 0 {
		width, height = 0, 0
	}

	if req.BigImageUrl == "" {
		req.BigImageUrl = req.ImageUrl
	}
	if req.DownloadUrl == "" {
		req.DownloadUrl = req.ImageUrl
	}
	if req.ThumbParam == "" {
		// *** 调整此处：使用 config.GetSetting("DEFAULT_THUMB_PARAM") ***
		req.ThumbParam = config.GetSetting("DEFAULT_THUMB_PARAM")
	}
	if req.BigParam == "" {
		// *** 调整此处：使用 config.GetSetting("DEFAULT_BIG_PARAM") ***
		req.BigParam = config.GetSetting("DEFAULT_BIG_PARAM")
	}

	tagsStr := ""
	if len(req.Tags) > 0 {
		tagsStr = utils.JoinTags(req.Tags)
	}

	album := model.Album{
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
		FileHash:      req.FileHash,
	}

	if err := database.DB.Create(&album).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "添加失败")
		return
	}

	for _, tag := range req.Tags {
		var count int64
		database.DB.Model(&model.Tag{}).Where("name = ?", tag).Count(&count)
		if count == 0 {
			database.DB.Create(&model.Tag{Name: tag})
		}
	}

	response.Success(c, nil, "添加成功")
}

// 删除图片
func DeleteAlbum(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "ID非法")
		return
	}

	if err := database.DB.Delete(&model.Album{}, id).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(c, nil, "删除成功")
}

// 更新图片
func UpdateAlbum(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "ID非法")
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
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	if req.BigImageUrl == "" {
		req.BigImageUrl = req.ImageUrl
	}
	if req.DownloadUrl == "" {
		req.DownloadUrl = req.ImageUrl
	}
	if req.ThumbParam == "" {
		// *** 调整此处：使用 config.GetSetting("DEFAULT_THUMB_PARAM") ***
		req.ThumbParam = config.GetSetting("DEFAULT_THUMB_PARAM")
	}
	if req.BigParam == "" {
		// *** 调整此处：使用 config.GetSetting("DEFAULT_BIG_PARAM") ***
		req.BigParam = config.GetSetting("DEFAULT_BIG_PARAM")
	}

	tagsStr := utils.JoinTags(req.Tags)

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

	if err := database.DB.Model(&model.Album{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(c, nil, "更新成功")
}
