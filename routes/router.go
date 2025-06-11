/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-04-11 06:02:54
 * @LastEditTime: 2025-06-11 12:22:54
 * @LastEditors: 安知鱼
 */
package routes

import (
	"wallpaper-admin/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(api *gin.RouterGroup) {
	// 私有接口
	api.POST("/login", controller.Login)
	api.POST("/refresh-token", controller.RefreshToken)
	RegisterWallpaperRoutes(api)

	// 公共接口
	public := api.Group("/public")
	{
		public.GET("/wallpapers", controller.GetPublicWallpapers)
		public.PUT("/stat/:id", controller.UpdateWallpaperStat)
		public.GET("/site-config", controller.GetSiteConfig) 
	}
}
