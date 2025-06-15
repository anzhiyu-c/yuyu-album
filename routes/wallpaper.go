/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-04-11 15:18:15
 * @LastEditTime: 2025-04-11 17:27:24
 * @LastEditors: 安知鱼
 */
package routes

import (
	"album-admin/controller"

	"github.com/gin-gonic/gin"
)

func RegisterAlbumRoutes(r *gin.RouterGroup) {
	r.GET("/albums", controller.GetAlbums)
	r.POST("/albums", controller.AddAlbum)
	r.PUT("/albums/:id", controller.UpdateAlbum)
	r.DELETE("/albums/:id", controller.DeleteAlbum)
}
