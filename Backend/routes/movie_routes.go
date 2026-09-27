package routes

import (
	"cinema-ticket-booking/controllers"

	"github.com/gin-gonic/gin"
)

func MovieRoutes(router *gin.Engine) {
	registerMovieRoutes(router.Group("/movies"))
	registerMovieRoutes(router.Group("/api/movies"))
}

func registerMovieRoutes(movie *gin.RouterGroup) {
	movie.GET("", controllers.GetMovies)
	movie.GET("/:movie_id", controllers.GetMovieByID)

	movie.POST("", controllers.CreateMovie)
	movie.PUT("/:movie_id", controllers.UpdateMovie)
	movie.DELETE("/:movie_id", controllers.DeleteMovie)
}
