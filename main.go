package main

import (
	"album-admin/config"
	"album-admin/database"
	"album-admin/middleware"
	"album-admin/migrate"
	"album-admin/routes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

// --- 关键修正：修改 //go:embed 指令 ---
// 嵌入整个 assets/dist 目录及其所有内容。
// 这将使 `content` 的根目录变为 `assets/dist/`
//
//go:embed assets/dist
var content embed.FS

var parsedTemplates *template.Template

type CustomHTMLRender struct {
	Templates *template.Template
}

func (r CustomHTMLRender) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r.Templates,
		Name:     name,
		Data:     data,
	}
}

func main() {
	config.LoadConfig()
	fmt.Println("正在初始化 MySQL 数据库连接...")
	database.InitMySQL()
	fmt.Println("正在初始化 Redis 数据库连接...")
	database.InitRedis()
	fmt.Println("正在执行数据库迁移...")
	migrate.AutoMigrateTables(database.DB)
	config.InitSettings(database.DB)
	if err := config.LoadSettingsFromDB(database.DB); err != nil {
		log.Fatalf("Failed to load site settings from database: %v", err)
	}

	var err error
	// --- 关键修正：调整 template.ParseFS 中的路径 ---
	// `content` 的根是 `assets/dist`。要解析 `index.html`，需要给出其在 `embed.FS` 中的完整路径。
	parsedTemplates, err = template.ParseFS(content, "assets/dist/index.html")
	if err != nil {
		log.Fatalf("Failed to parse HTML template from embedded file system: %v", err)
	}

	// --- 调试检查：确认 platform-config.json 已经正确嵌入 ---
	// 修正：路径是 "assets/dist/platform-config.json"
	_, fileErr := content.ReadFile("assets/dist/platform-config.json")
	if fileErr != nil {
		log.Fatalf("错误: 'platform-config.json' 未在嵌入式文件系统中找到或无法读取。请确保它位于 'assets/dist/' 并且正确嵌入: %v", fileErr)
	} else {
		fmt.Println("调试: 'platform-config.json' 已成功在嵌入式文件系统中找到。")
	}
	// --- 调试检查结束 ---

	r := gin.Default()
	r.HTMLRender = CustomHTMLRender{Templates: parsedTemplates}
	r.Use(middleware.Cors())

	apiGroup := r.Group("/api")
	routes.SetupRouter(apiGroup)

	// --- 关键修改：调整静态文件服务的路径以匹配 Vue 构建输出 ---
	// 假设 Vue 的 JS/CSS/图片等资源在 `assets/dist/static` 目录下，并以 `/static/` 路径访问
	// 修正：创建子文件系统时，路径是 "assets/dist/static"
	subStaticFs, err := fs.Sub(content, "assets/dist/static")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem for assets/dist/static: %v", err)
	}
	r.StaticFS("/static", http.FS(subStaticFs))

	// 根目录下的静态文件，如 /logo.svg, /favicon.ico, /manifest.json
	// 修正：路径是相对于 `content` 的，所以是 "assets/dist/logo.svg"
	r.StaticFileFS("/logo.svg", "assets/dist/logo.svg", http.FS(content))
	r.StaticFileFS("/favicon.ico", "assets/dist/favicon.ico", http.FS(content))
	r.StaticFileFS("/manifest.json", "assets/dist/manifest.json", http.FS(content))

	// --- 手动处理 platform-config.json 路由 ---
	// 修正：路径是 "assets/dist/platform-config.json"
	r.GET("/platform-config.json", func(c *gin.Context) {
		fileData, readErr := content.ReadFile("assets/dist/platform-config.json")
		if readErr != nil {
			log.Printf("Failed to read platform-config.json from embedded FS: %v", readErr)
			c.Status(http.StatusNotFound)
			return
		}
		c.Data(http.StatusOK, "application/json", fileData)
	})

	r.NoRoute(func(c *gin.Context) {
		requestPath := c.Request.URL.Path

		if strings.HasPrefix(requestPath, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "API route not found"})
			return
		}
		if strings.HasPrefix(requestPath, "/static/") {
			c.Status(http.StatusNotFound)
			return
		}

		dataToTemplate := gin.H{
			"siteName":     config.GetSetting("APP_NAME"),
			"favicon":      config.GetSetting("ICON_URL"),
			"keywords":     config.GetSetting("SITE_KEYWORDS"),
			"description":  config.GetSetting("SITE_DESCRIPTION"),
			"pwaSmallIcon": config.GetSetting("ICON_URL"),
			"siteScript":   "",
		}

		c.HTML(http.StatusOK, "index.html", dataToTemplate) // `c.HTML` 的模板名称是 `index.html`
	})

	port := config.Conf.GetString("PORT")
	if port == "" {
		port = "8091"
	}
	fmt.Printf("应用程序启动成功，正在监听端口: %s\n", port)
	r.Run(":" + port)
}
