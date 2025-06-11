package utils

import (
	"time"
	"wallpaper-admin/config"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("")

func InitJWT() {
	jwtKey = []byte(config.Conf.GetString("JWT_SECRET"))
}

func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 720).Unix(), // 设置长时间的过期时间（例如：30天）
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
}
