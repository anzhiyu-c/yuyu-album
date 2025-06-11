/*
 * @Description: 壁纸管理后台服务入口
 * @Author: 安知鱼
 * @Date: 2025-04-11 06:02:54
 * @LastEditTime: 2025-04-12 15:29:54
 * @LastEditors: 安知鱼
 */

package main

import (
	"embed"
	"net/http"
	"wallpaper-admin/config"
	"wallpaper-admin/database"
	"wallpaper-admin/middleware"
	"wallpaper-admin/routes"
	"wallpaper-admin/utils"

	"github.com/gin-gonic/gin"
)

//go:embed assets/dist/*
var content embed.FS

func main() {
	// 加载配置、初始化工具和数据库
	config.LoadConfig()
	utils.InitJWT()
	database.InitMySQL()
	database.InitRedis()

	r := gin.Default()

	// 添加跨域中间件
	r.Use(middleware.Cors())

	// 注册 API 路由组
	apiGroup := r.Group("/api")
	routes.SetupRouter(apiGroup)

	// 静态资源嵌入式服务，挂载在根路径
	r.Use(middleware.Serve("/", middleware.EmbedFolder(content, "assets/dist")))

	// 没有匹配到API或静态资源，统一返回前端 index.html (支持前端history路由)
	r.NoRoute(func(c *gin.Context) {
		// 对于其他请求，返回嵌入的index.html
		data, err := content.ReadFile("assets/dist/index.html")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	// 启动服务
	port := config.Conf.GetString("PORT")
	if port == "" {
		port = "8091"
	}
	r.Run(":" + port)
}
