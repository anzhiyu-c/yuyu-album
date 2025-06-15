/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-04-11 06:02:54
 * @LastEditTime: 2025-06-11 12:22:54
 * @LastEditors: 安知鱼
 */
package routes

import (
	"album-admin/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(api *gin.RouterGroup) {
	// 私有接口
	api.POST("/login", controller.Login)
	api.POST("/refresh-token", controller.RefreshToken)
	RegisterAlbumRoutes(api)

	// 公共接口
	public := api.Group("/public")
	{
		public.GET("/albums", controller.GetPublicAlbums)
		public.PUT("/stat/:id", controller.UpdateAlbumStat)
		public.GET("/site-config", controller.GetSiteConfig)
	}
}
