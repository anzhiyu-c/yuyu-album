package controller

import (
	"album-admin/config" // 导入 config 包
	"album-admin/utils/response"

	"github.com/gin-gonic/gin"
)

func GetSiteConfig(c *gin.Context) {
	// 直接从 config 包中获取所有已加载的设置
	settings := make(map[string]string)
	for k, v := range config.SiteSettings {
		settings[k] = v
	}

	delete(settings, "JWT_SECRET")

	response.Success(c, settings, "获取站点配置成功")
}
