package controller

import (
	"album-admin/config"
	"album-admin/database"
	"album-admin/model"
	"album-admin/utils"
	"album-admin/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

// Register 处理用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误：用户名和密码不能为空")
		return
	}

	// 检查用户名是否已存在
	var existingUser model.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		response.Fail(c, http.StatusConflict, "该用户名已被注册")
		return
	}

	// 检查是否是第一个注册的用户
	var userCount int64
	database.DB.Model(&model.User{}).Count(&userCount)

	role := "normal" // 默认角色为普通用户
	if userCount == 0 {
		role = "admin" // 如果是第一个注册的用户，设置为管理员
		// 第一次注册时，如果昵称为空，可以默认设置为应用名称
		if req.Nickname == "" {
			req.Nickname = config.GetSetting("APP_NAME")
		}
	}

	// 对密码进行哈希处理
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "密码加密失败")
		return
	}

	// 如果昵称为空，提供一个默认值（例如：用户xxx）
	if req.Nickname == "" {
		req.Nickname = "用户" + req.Username
	}

	// 获取默认头像
	defaultAvatar := config.GetSetting("USER_AVATAR")

	user := model.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Nickname:     req.Nickname,
		Avatar:       defaultAvatar, // 注册时使用默认头像
		Email:        req.Email,
		Role:         role,
		Status:       1, // 新注册用户默认为正常状态
	}

	if err := database.DB.Create(&user).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "注册失败")
		return
	}

	response.Success(c, nil, "注册成功")
}

// UpdateUserPassword 用于用户修改自身密码的API
func UpdateUserPassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误：旧密码和新密码都不能为空")
		return
	}

	username := c.GetString("username")
	if username == "" {
		response.Fail(c, http.StatusUnauthorized, "未登录或无法获取当前用户信息")
		return
	}

	var user model.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取用户信息失败")
		return
	}

	if !utils.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
		response.Fail(c, http.StatusUnauthorized, "旧密码不正确")
		return
	}

	newHashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "密码哈希失败")
		return
	}

	if err := database.DB.Model(&user).Update("PasswordHash", newHashedPassword).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新密码失败")
		return
	}

	response.Success(c, nil, "密码修改成功")
}

// GetUserInfo 获取当前登录用户的信息
func GetUserInfo(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		response.Fail(c, http.StatusUnauthorized, "未登录或无法获取当前用户信息")
		return
	}

	var user model.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "用户不存在")
		return
	}

	// 为了安全，不返回密码哈希
	user.PasswordHash = ""

	response.Success(c, user, "获取用户信息成功")
}
