package routes

import (
	"fmt"

	"cinema-ticket-booking/controllers"

	"github.com/gin-gonic/gin"

	"cinema-ticket-booking/middleware"
)

func UserRoutes(router *gin.Engine) {

	fmt.Println("✅ UserRoutes Loaded")

	router.GET("/users", controllers.GetUsers)
	router.POST("/users", controllers.CreateUser)
	router.PUT("/users/:user_id", controllers.UpdateUser)
	router.DELETE("/users/:user_id", controllers.DeleteUser)

	router.GET(
		"/users/me",
		middleware.AuthMiddleware(),
		controllers.GetMe,
	)
	router.PUT(
		"/users/me",
		middleware.AuthMiddleware(),
		controllers.UpdateMe,
	)
	router.PUT(
		"/users/me/password",
		middleware.AuthMiddleware(),
		controllers.ChangeMyPassword,
	)
	router.POST(
		"/users/me/profile-image",
		middleware.AuthMiddleware(),
		controllers.UploadProfileImage,
	)
}
