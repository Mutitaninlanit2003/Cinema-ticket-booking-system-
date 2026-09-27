package routes

import (
	"cinema-ticket-booking/controllers"

	"github.com/gin-gonic/gin"
)

func SeatLockRoutes(
	router *gin.Engine,
	controller *controllers.SeatLockController,
) {
	showtimes := router.Group("/showtimes")

	showtimes.POST(
		"/:showtime_id/seat-lock/:seat_id",
		controller.LockSeat,
	)

	showtimes.DELETE(
		"/:showtime_id/seat-lock/:seat_id",
		controller.UnlockSeat,
	)

	println("✅ SeatLockRoutes Loaded")
}
