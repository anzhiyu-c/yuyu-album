package controller

import (
	"net/http"
	"time"
	"wallpaper-admin/config"
	"wallpaper-admin/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Avatar       string   `json:"avatar"`
		Username     string   `json:"username"`
		Nickname     string   `json:"nickname"`
		Roles        []string `json:"roles"`
		Permissions  []string `json:"permissions"`
		AccessToken  string   `json:"accessToken"`
		RefreshToken string   `json:"refreshToken"`
		Expires      string   `json:"expires"`
	} `json:"data"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	adminUser := config.Conf.GetString("ADMIN_USERNAME")
	adminPass := config.Conf.GetString("ADMIN_PASSWORD")

	if req.Username != adminUser || req.Password != adminPass {
		c.JSON(http.StatusOK, gin.H{"message": "账号或密码错误", "code": 401})
		return
	}

	// 生成 AccessToken 和 RefreshToken
	accessToken, err := utils.GenerateToken(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成Token失败"})
		return
	}

	refreshToken, err := utils.GenerateToken(req.Username) // 刷新 Token 用不同方式生成或有效期长一点
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成RefreshToken失败"})
		return
	}

	// 设置过期时间（例如：expires in 30 days）
	expires := time.Now().Add(time.Hour * 720).Format("2006/01/02 15:04:05")

	// 返回
	response := LoginResponse{
		Success: true,
	}

	appLoginImage := config.Conf.GetString("USER_AVATAR")
	appName := config.Conf.GetString("APP_NAME")

	response.Data.Avatar = appLoginImage // 从配置文件获取头像 URL
	response.Data.Username = "admin"
	response.Data.Nickname = appName // 从配置文件获取昵称
	response.Data.Roles = []string{"admin"}
	response.Data.Permissions = []string{"admin"}
	response.Data.AccessToken = accessToken
	response.Data.RefreshToken = refreshToken
	response.Data.Expires = expires

	c.JSON(http.StatusOK, response)
}

func RefreshToken(c *gin.Context) {
	// 获取 refreshToken
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供RefreshToken"})
		return
	}

	// 解析 RefreshToken
	token, err := utils.ParseToken(refreshToken)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效RefreshToken"})
		return
	}

	// 生成新的 AccessToken
	username := token.Claims.(jwt.MapClaims)["username"].(string)
	accessToken, err := utils.GenerateToken(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成新Token失败"})
		return
	}

	// 设置过期时间
	expires := time.Now().Add(time.Hour * 720).Format("2006/01/02 15:04:05")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"accessToken": accessToken,
			"expires":     expires,
		},
	})
}
