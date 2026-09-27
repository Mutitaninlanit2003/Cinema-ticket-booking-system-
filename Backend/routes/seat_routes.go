package routes

import (
	"fmt"

	"cinema-ticket-booking/controllers"

	"github.com/gin-gonic/gin"
)

func SeatRoutes(router *gin.Engine) {

	fmt.Println("✅ SeatRoutes Loaded")

	router.GET(
		"/seats",
		controllers.GetSeats,
	)

	router.POST(
		"/seats",
		controllers.CreateSeat,
	)

	// Generate multiple seats
	router.POST(
		"/seats/generate",
		controllers.GenerateSeats,
	)

	router.PUT(
		"/seats/:seat_id",
		controllers.UpdateSeat,
	)

	router.DELETE(
		"/seats/:seat_id",
		controllers.DeleteSeat,
	)
	// router.DELETE(
	// 	"/screens/:screen_id/seats",
	// 	controllers.DeleteSeatsByScreenID,
	// )
}
