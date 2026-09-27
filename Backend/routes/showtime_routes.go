package routes

import (
	"cinema-ticket-booking/controllers"

	"github.com/gin-gonic/gin"
)

func ShowtimeRoutes(
	router *gin.Engine,
	controller *controllers.ShowtimeController,
) {

	showtimes := router.Group("/showtimes")
	{

		showtimes.POST(
			"",
			controller.CreateShowtime,
		)

		showtimes.GET(
			"",
			controller.GetShowtimes,
		)
		showtimes.PUT(
			"/:showtime_id",
			controller.UpdateShowtime,
		)

		showtimes.DELETE(
			"/:showtime_id",
			controller.DeleteShowtime,
		)

	}

}
