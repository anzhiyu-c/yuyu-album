/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-06-15 13:06:01
 * @LastEditTime: 2025-06-15 13:06:06
 * @LastEditors: 安知鱼
 */
package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 对密码进行哈希处理
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 验证密码哈希
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
