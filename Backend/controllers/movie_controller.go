package controllers

import (
	"errors"
	"net/http"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var movieService *services.MovieService

// GetMovies godoc
// @Summary Get Movies
// @Description Retrieve all movies
// @Tags Movies
// @Produce json
// @Success 200 {array} models.Movie
// @Router /movies [get]
func GetMovies(c *gin.Context) {

	// สร้าง Service หลังจาก ConnectDB() แล้ว
	if movieService == nil {
		movieService = services.NewMovieService()
	}

	movies, err := movieService.GetMovies()

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, movies)
}

// Create Movie
// @Summary Create movie
// @Description Create new movie
// @Tags Movies
// @Accept json
// @Produce json
// @Param movie body models.Movie true "Movie data"
// @Success 201 {object} models.Movie
// @Router /movies [post]
func CreateMovie(c *gin.Context) {

	// สร้าง Service หลังจาก ConnectDB() แล้ว
	if movieService == nil {
		movieService = services.NewMovieService()
	}

	var movie models.Movie

	if err := c.ShouldBindJSON(&movie); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err := movieService.CreateMovie(movie); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Movie created successfully",
	})
}

// UpdateMovie godoc
// @Summary Update movie
// @Description Update movie by ID
// @Tags Movies
// @Accept json
// @Produce json
// @Param movie_id path string true "Movie ID"
// @Param movie body models.Movie true "Movie"
// @Success 200 {object} map[string]string
// @Router /movies/{movie_id} [put]
func UpdateMovie(c *gin.Context) {

	if movieService == nil {
		movieService = services.NewMovieService()
	}

	movieID := c.Param("movie_id")

	var movie models.Movie

	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := movieService.UpdateMovie(movieID, movie); err != nil {

		if err.Error() == "movie not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Movie updated successfully",
	})
}

// DeleteMovie godoc
// @Summary Delete movie
// @Description Delete movie by ID
// @Tags Movies
// @Produce json
// @Param movie_id path string true "Movie ID"
// @Success 200 {object} map[string]string
// @Router /movies/{movie_id} [delete]
func DeleteMovie(c *gin.Context) {

	if movieService == nil {
		movieService = services.NewMovieService()
	}

	movieID := c.Param("movie_id")

	if err := movieService.DeleteMovie(movieID); err != nil {

		if err.Error() == "movie not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Movie deleted successfully",
	})
}

func GetMovieByID(c *gin.Context) {

	if movieService == nil {
		movieService = services.NewMovieService()
	}

	movieID := c.Param("movie_id")

	movie, err := movieService.GetMovieByID(movieID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "movie not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, movie)
}
