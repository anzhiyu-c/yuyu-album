// album-admin/middleware/auth.go
package middleware

import (
	"album-admin/utils/jwtutil"
	"album-admin/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware 是一个用于验证 JWT Token 的中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			// 使用 response.Fail 返回错误
			response.Fail(c, http.StatusUnauthorized, "请求头中缺少Auth Token或格式不正确")
			c.Abort() // 阻止后续处理器执行
			return
		}

		tokenString = tokenString[7:] // 提取Bearer后面的Token

		// 解析Token
		token, err := jwtutil.ParseToken(tokenString)
		if err != nil || !token.Valid {
			response.Fail(c, http.StatusUnauthorized, "无效的Auth Token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Fail(c, http.StatusUnauthorized, "无效的Auth Token Claims")
			c.Abort()
			return
		}

		// 从 claims 中获取用户信息并存储到Context中
		// 你可能需要根据你的实际 JWT Payload 结构来获取这些信息
		// 例如，如果你的 claims 中有 "username" 和 "roles"
		username, userExists := claims["username"].(string)
		rolesSlice, rolesExists := claims["roles"].([]interface{}) // JWT中数组通常解析为 []interface{}

		if !userExists || !rolesExists {
			response.Fail(c, http.StatusUnauthorized, "Auth Token信息不完整")
			c.Abort()
			return
		}

		// 将 rolesSlice 转换为 []string
		var roles []string
		for _, role := range rolesSlice {
			if r, ok := role.(string); ok {
				roles = append(roles, r)
			}
		}

		c.Set("username", username)
		c.Set("roles", roles)

		c.Next() // 继续处理请求
	}
}

// AdminAuthMiddleware 是一个用于验证用户是否为管理员的中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Context中获取用户角色信息（由JWTAuthMiddleware设置）
		rolesAny, exists := c.Get("roles")
		if !exists {
			response.Fail(c, http.StatusForbidden, "权限不足：无法获取用户角色信息")
			c.Abort()
			return
		}

		userRoles, ok := rolesAny.([]string) // 确保类型转换正确
		if !ok {
			response.Fail(c, http.StatusForbidden, "权限不足：无效的用户角色信息")
			c.Abort()
			return
		}

		isAdmin := false
		for _, role := range userRoles {
			if role == "admin" {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			response.Fail(c, http.StatusForbidden, "权限不足：此操作需要管理员权限")
			c.Abort()
			return
		}

		c.Next() // 继续处理请求
	}
}
