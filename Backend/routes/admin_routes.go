package routes

import (
	"fmt"

	"cinema-ticket-booking/controllers"
	"cinema-ticket-booking/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(
	router *gin.Engine,
	controller *controllers.AdminController,
) {
	fmt.Println("✅ AdminRoutes Loaded")

	admin := router.Group("/admin")

	admin.Use(
		middleware.AuthMiddleware(),
		middleware.RequireRole("admin"),
	)

	// Admin Booking
	admin.GET("/bookings", controller.GetBookings)

	// Admin Profile
	admin.GET("/me", controller.GetMe)
	admin.PUT("/me", controller.UpdateMe)
}

// package routes

// import (
// 	"fmt"

// 	"cinema-ticket-booking/controllers"
// 	"cinema-ticket-booking/middleware"

// 	"github.com/gin-gonic/gin"
// )

// func AdminRoutes(
// 	router *gin.Engine,
// 	controller *controllers.AdminController,
// ) {

// 	fmt.Println("✅ AdminRoutes Loaded")

// 	admin := router.Group("/admin")

// 	admin.Use(
// 		middleware.AuthMiddleware(),
// 		middleware.RequireRole("admin"),
// 	)

// 	admin.GET("/bookings", controller.GetBookings)
// }
