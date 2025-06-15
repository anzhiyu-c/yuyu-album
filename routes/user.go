package routes

import (
	"album-admin/controller"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	r.GET("/user/check-email", controller.CheckEmail)
}
