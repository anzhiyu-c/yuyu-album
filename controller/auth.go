package controller

import (
	"album-admin/config"
	"album-admin/database"
	"album-admin/model"
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

	var user model.User
	// 从数据库查询用户
	// 同时检查用户状态，只有状态为1（正常）的用户才能登录
	if err := database.DB.Where("username = ? AND status = ?", req.Username, 1).First(&user).Error; err != nil {
		// 如果用户不存在、被禁用或查询失败，统一返回“账号或密码错误”以避免信息泄露
		response.Fail(c, http.StatusUnauthorized, "账号或密码错误")
		return
	}

	// 验证密码
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		response.Fail(c, http.StatusUnauthorized, "账号或密码错误")
		return
	}

	// 更新最后登录时间
	now := time.Now()
	database.DB.Model(&user).Update("LastLoginAt", &now)

	accessToken, err := utils.GenerateToken(user.Username)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "生成AccessToken失败")
		return
	}

	refreshToken, err := utils.GenerateToken(user.Username)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "生成RefreshToken失败")
		return
	}

	expires := time.Now().Add(time.Hour * 720).Format("2006/01/02 15:04:05")

	// 从用户模型中获取头像和昵称，如果用户模型中没有，则使用站点配置的默认值
	userAvatar := user.Avatar
	if userAvatar == "" {
		userAvatar = config.GetSetting("USER_AVATAR")
	}
	userNickname := user.Nickname
	if userNickname == "" {
		userNickname = config.GetSetting("APP_NAME") // 可以用应用名作为默认昵称
	}

	data := LoginResponseData{
		Avatar:       userAvatar,
		Username:     user.Username,
		Nickname:     userNickname,
		Roles:        []string{user.Role},
		Permissions:  []string{user.Role + "_permission"},
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

	// 校验用户是否存在且状态正常
	var user model.User
	if err := database.DB.Where("username = ? AND status = ?", username, 1).First(&user).Error; err != nil {
		response.Fail(c, http.StatusUnauthorized, "用户不存在或已被禁用")
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
