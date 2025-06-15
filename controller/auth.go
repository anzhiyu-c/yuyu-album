package controller

import (
	"album-admin/config"
	"album-admin/utils"
	"album-admin/utils/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseData struct {
	Avatar       string   `json:"avatar"`
	Username     string   `json:"username"`
	Nickname     string   `json:"nickname"`
	Roles        []string `json:"roles"`
	Permissions  []string `json:"permissions"`
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	Expires      string   `json:"expires"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	adminUser := config.Conf.GetString("ADMIN_USERNAME")
	adminPass := config.Conf.GetString("ADMIN_PASSWORD")

	if req.Username != adminUser || req.Password != adminPass {
		response.Fail(c, http.StatusUnauthorized, "账号或密码错误")
		return
	}

	accessToken, err := utils.GenerateToken(req.Username)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "生成AccessToken失败")
		return
	}

	refreshToken, err := utils.GenerateToken(req.Username)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "生成RefreshToken失败")
		return
	}

	expires := time.Now().Add(time.Hour * 720).Format("2006/01/02 15:04:05")

	appLoginImage := config.GetSetting("USER_AVATAR")
	appName := config.GetSetting("APP_NAME")

	data := LoginResponseData{
		Avatar:       appLoginImage,
		Username:     "admin",
		Nickname:     appName,
		Roles:        []string{"admin"},
		Permissions:  []string{"admin"},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expires:      expires,
	}

	response.Success(c, data, "登录成功")
}

func RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		response.Fail(c, http.StatusUnauthorized, "未提供RefreshToken")
		return
	}

	token, err := utils.ParseToken(refreshToken)
	if err != nil || !token.Valid {
		response.Fail(c, http.StatusUnauthorized, "无效RefreshToken")
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "无效RefreshToken Claims")
		return
	}
	username, ok := claims["username"].(string)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "RefreshToken中缺少用户名信息")
		return
	}

	accessToken, err := utils.GenerateToken(username)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "生成新Token失败")
		return
	}

	expires := time.Now().Add(time.Hour * 720).Format("2006/01/02 15:04:05")

	response.Success(c, gin.H{
		"accessToken": accessToken,
		"expires":     expires,
	}, "刷新Token成功")
}
