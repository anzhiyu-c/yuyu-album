/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-04-11 15:18:15
 * @LastEditTime: 2025-04-11 17:27:24
 * @LastEditors: 安知鱼
 */
package routes

import (
	"wallpaper-admin/controller"

	"github.com/gin-gonic/gin"
)

func RegisterWallpaperRoutes(r *gin.RouterGroup) {
	r.GET("/wallpapers", controller.GetWallpapers)
	r.POST("/wallpapers", controller.AddWallpaper)
	r.PUT("/wallpapers/:id", controller.UpdateWallpaper)
	r.DELETE("/wallpapers/:id", controller.DeleteWallpaper)
}
