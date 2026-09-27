package routes

import (
	"cinema-ticket-booking/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(
	router *gin.Engine,
	controller *controllers.AuthController,
) {

	auth := router.Group("/auth")

	auth.POST("/register", controller.Register)
	auth.POST("/login", controller.Login)
	auth.POST("/forgot-password", controller.ForgotPassword)
	auth.POST("/reset-password", controller.ResetPassword)

	auth.GET("/google", controller.GoogleLogin)
	auth.GET("/google/callback", controller.GoogleCallback)
}
