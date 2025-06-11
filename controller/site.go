/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-06-11 12:22:08
 * @LastEditTime: 2025-06-11 13:31:45
 * @LastEditors: 安知鱼
 */
// controller/site.go
package controller

import (
	"net/http"
	"wallpaper-admin/config"

	"github.com/gin-gonic/gin"
)

// GetSiteConfig 获取站点配置
func GetSiteConfig(c *gin.Context) {
	// 获取 .env 文件中的配置项
	siteConfig := gin.H{
		"APP_NAME":    config.Conf.GetString("APP_NAME"),
		"APP_VERSION": config.Conf.GetString("APP_VERSION"),
		"ICP_NUMBER":  config.Conf.GetString("ICP_NUMBER"),
		"USER_AVATAR": config.Conf.GetString("USER_AVATAR"),
		"ABOUT_LINK":  config.Conf.GetString("ABOUT_LINK"),
		"API_URL":     config.Conf.GetString("API_URL"),
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": siteConfig,
	})
}
