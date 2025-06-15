// utils/jwt.go (修改)
package utils

import (
	"album-admin/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(username string) (string, error) {
	jwtSecret := config.GetSetting("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT Secret is not loaded from site settings")
	}

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 720).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("生成token失败: %w", err)
	}
	return tokenString, nil
}

func ParseToken(tokenStr string) (*jwt.Token, error) {
	jwtSecret := config.GetSetting("JWT_SECRET") // 从新的配置获取方式
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT Secret is not loaded from site settings")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析token失败: %w", err)
	}
	return token, nil
}
