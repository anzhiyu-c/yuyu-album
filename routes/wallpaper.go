/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-04-11 15:18:15
 * @LastEditTime: 2025-06-16 00:01:18
 * @LastEditors: 安知鱼
 */
package routes

import (
	"album-admin/controller"
	"album-admin/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAlbumRoutes(r *gin.RouterGroup) {
	r.GET("/albums", middleware.AdminAuthMiddleware(), controller.GetAlbums)
	r.POST("/albums", middleware.AdminAuthMiddleware(), controller.AddAlbum)
	r.PUT("/albums/:id", middleware.AdminAuthMiddleware(), controller.UpdateAlbum)
	r.DELETE("/albums/:id", middleware.AdminAuthMiddleware(), controller.DeleteAlbum)
}
