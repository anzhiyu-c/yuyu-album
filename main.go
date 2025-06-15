package main

import (
	"album-admin/config"
	"album-admin/database"
	"album-admin/middleware"
	"album-admin/migrate" // 依然需要导入 model 包，但不是为了 AutoMigrateTables，而是为了其他地方可能用到模型定义
	"album-admin/routes"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed assets/dist/*
var content embed.FS

func main() {
	// 1. 加载 .env 或环境变量配置 (例如数据库连接信息、端口等)
	config.LoadConfig()

	// 2. 初始化 MySQL 数据库连接
	fmt.Println("正在初始化 MySQL 数据库连接...")
	database.InitMySQL() // database.DB 全局变量在这里被赋值
	fmt.Println("正在初始化 Redis 数据库连接...")
	database.InitRedis()

	// 3. 执行数据库迁移 (确保所有表都已准备好，包括新的 Setting 表)
	fmt.Println("正在执行数据库迁移...")
	// *** 修正：直接调用 AutoMigrateTables，不再传入模型参数 ***
	migrate.AutoMigrateTables(database.DB)

	// 4. 初始化站点配置 (将默认键值对配置插入数据库，如果不存在)
	// 在数据库连接成功且表结构迁移后执行
	database.InitSettings(database.DB) // 调用新的 InitSettings 函数，传递数据库实例

	// 5. 从数据库加载所有键值对配置到 config.SiteSettings
	if err := config.LoadSettingsFromDB(database.DB); err != nil {
		log.Fatalf("Failed to load site settings from database: %v", err)
	}

	// 6. 初始化 Gin 引擎
	r := gin.Default()

	// 7. 添加跨域中间件
	r.Use(middleware.Cors())

	// 8. 注册 API 路由组
	apiGroup := r.Group("/api")
	routes.SetupRouter(apiGroup)

	// 9. 静态资源嵌入式服务，挂载在根路径
	r.Use(middleware.Serve("/", middleware.EmbedFolder(content, "assets/dist")))

	// 10. 没有匹配到API或静态资源，统一返回前端 index.html (支持前端history路由)
	r.NoRoute(func(c *gin.Context) {
		data, err := content.ReadFile("assets/dist/index.html")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	// 11. 启动服务
	port := config.Conf.GetString("PORT")
	if port == "" {
		port = "8091" // 提供一个默认端口
	}
	fmt.Printf("应用程序启动成功，正在监听端口: %s\n", port)
	r.Run(":" + port)
}
