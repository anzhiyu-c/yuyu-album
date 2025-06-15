package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 是统一的API返回结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // Data可以是任意类型，omitempty表示如果为空则不显示
}

// Success 成功响应
func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
