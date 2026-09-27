package routes

import (
	"fmt"

	"cinema-ticket-booking/controllers"

	"github.com/gin-gonic/gin"
)

func TheaterRoutes(router *gin.Engine) {

	fmt.Println("✅ TheaterRoutes Loaded")

	router.GET("/theaters", controllers.GetTheaters)
	router.POST("/theaters", controllers.CreateTheater)
	router.PUT("/theaters/:theater_id", controllers.UpdateTheater)
	router.DELETE("/theaters/:theater_id", controllers.DeleteTheater)
}
