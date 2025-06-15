/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-06-15 11:30:55
 * @LastEditTime: 2025-06-15 23:45:48
 * @LastEditors: 安知鱼
 */

package routes

import (
	"album-admin/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(api *gin.RouterGroup) {
	RegisterAlbumRoutes(api)
	RegisterUserRoutes(api)

	// 公共接口
	public := api.Group("/public")
	{
		public.GET("/albums", controller.GetPublicAlbums)
		public.PUT("/stat/:id", controller.UpdateAlbumStat)
		public.GET("/site-config", controller.GetSiteConfig)
	}
}
