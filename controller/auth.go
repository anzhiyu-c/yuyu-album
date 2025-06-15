package controller

import (
	"album-admin/config"
	"album-admin/database"
	"album-admin/model"
	"album-admin/utils"
	"album-admin/utils/jwtutil"
	"album-admin/utils/response"
	"net/http"
	"strings"
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

	var user model.User
	if err := database.DB.Where("username = ? AND status = ?", req.Username, 1).First(&user).Error; err != nil {
		response.Fail(c, http.StatusUnauthorized, "账号或密码错误")
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		response.Fail(c, http.StatusUnauthorized, "账号或密码错误")
		return
	}

	// 更新最后登录时间
	now := time.Now()
	database.DB.Model(&user).Update("LastLoginAt", &now)

	// 支持多个角色用逗号分隔，例如："admin,editor"
	roles := strings.Split(user.Role, ",")

	accessToken, err := jwtutil.GenerateToken(user.Username, roles)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "生成AccessToken失败")
		return
	}

	refreshToken, err := jwtutil.GenerateToken(user.Username, roles)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "生成RefreshToken失败")
		return
	}

	expires := time.Now().Add(time.Hour * 720).Format("2006/01/02 15:04:05")

	userAvatar := user.Avatar
	if userAvatar == "" {
		userAvatar = config.GetSetting("USER_AVATAR")
	}
	userNickname := user.Nickname
	if userNickname == "" {
		userNickname = config.GetSetting("APP_NAME")
	}

	permissions := []string{}
	for _, role := range roles {
		permissions = append(permissions, role+"_permission")
	}

	data := LoginResponseData{
		Avatar:       userAvatar,
		Username:     user.Username,
		Nickname:     userNickname,
		Roles:        roles,
		Permissions:  permissions,
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

	token, err := jwtutil.ParseToken(refreshToken)
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

	// 校验用户是否存在且状态正常
	var user model.User
	if err := database.DB.Where("username = ? AND status = ?", username, 1).First(&user).Error; err != nil {
		response.Fail(c, http.StatusUnauthorized, "用户不存在或已被禁用")
		return
	}

	roles := strings.Split(user.Role, ",")

	accessToken, err := jwtutil.GenerateToken(username, roles)
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
