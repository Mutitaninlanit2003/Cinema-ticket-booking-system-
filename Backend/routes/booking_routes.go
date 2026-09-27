package routes

import (
	"fmt"

	"cinema-ticket-booking/controllers"

	"github.com/gin-gonic/gin"
)

func BookingRoutes(
	router *gin.Engine,
	controller *controllers.BookingController,
) {
	fmt.Println("✅ BookingRoutes Loaded")

	bookings := router.Group("/bookings")

	bookings.POST("", controller.CreateBooking)
	bookings.GET("", controller.GetBookings)
	bookings.PUT("/:booking_id", controller.UpdateBooking)
	bookings.DELETE("/:booking_id", controller.DeleteBooking)
}

// func BookingRoutes(
// 	router *gin.Engine,
// 	controller *controllers.BookingController,
// ) {

// 	bookings := router.Group("/bookings")

// 	{

// 		bookings.POST(
// 			"",
// 			controller.CreateBooking,
// 		)

// 		bookings.GET(
// 			"",
// 			controller.GetBookings,
// 		)
// 	}
// }
