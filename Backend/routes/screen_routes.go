package routes

import (
	"fmt"

	"cinema-ticket-booking/controllers"

	"github.com/gin-gonic/gin"
)

func ScreenRoutes(router *gin.Engine) {

	fmt.Println("✅ ScreenRoutes Loaded")

	router.GET("/screens", controllers.GetScreens)
	router.POST("/screens", controllers.CreateScreen)
	router.PUT("/screens/:screen_id", controllers.UpdateScreen)
	router.DELETE("/screens/:screen_id", controllers.DeleteScreen)
}
