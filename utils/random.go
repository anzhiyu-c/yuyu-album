/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-06-15 12:25:50
 * @LastEditTime: 2025-06-15 12:25:56
 * @LastEditors: 安知鱼
 */
package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomString
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// 使用 Base64 URL 编码，避免特殊字符问题
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
