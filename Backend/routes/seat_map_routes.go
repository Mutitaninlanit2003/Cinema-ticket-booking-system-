package routes

import (
	"cinema-ticket-booking/controllers"

	"github.com/gin-gonic/gin"
)

func SeatMapRoutes(
	router *gin.Engine,
	controller *controllers.SeatMapController,
) {

	println("✅ SeatMapRoutes Loaded")

	router.GET(
		"/showtimes/:showtime_id/seat-map",
		controller.WebSocket,
	)
}
