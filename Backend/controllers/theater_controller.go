package controllers

import (
	"net/http"
	"strings"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/services"

	"github.com/gin-gonic/gin"
)

var theaterService *services.TheaterService

// GetTheaters godoc
// @Summary Get all theaters
// @Description Retrieve all theaters
// @Tags Theaters
// @Produce json
// @Success 200 {array} models.Theater
// @Router /theaters [get]
func GetTheaters(c *gin.Context) {

	if theaterService == nil {
		theaterService = services.NewTheaterService()
	}

	theaters, err := theaterService.GetTheaters()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, theaters)
}

// CreateTheater godoc
// @Summary Create Theater
// @Description Create a new theater
// @Tags Theaters
// @Accept json
// @Produce json
// @Param theater body models.Theater true "Theater"
// @Success 201 {object} map[string]string
// @Router /theaters [post]
func CreateTheater(c *gin.Context) {

	if theaterService == nil {
		theaterService = services.NewTheaterService()
	}

	var theater models.Theater

	if err := c.ShouldBindJSON(&theater); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := theaterService.CreateTheater(&theater); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Theater created successfully",
	})
}

// UpdateTheater godoc
// @Summary Update theater
// @Description Update theater by ID
// @Tags Theaters
// @Accept json
// @Produce json
// @Param theater_id path string true "Theater ID"
// @Param theater body models.Theater true "Theater"
// @Success 200 {object} map[string]string
// @Router /theaters/{theater_id} [put]
func UpdateTheater(c *gin.Context) {

	if theaterService == nil {
		theaterService = services.NewTheaterService()
	}

	theaterID := c.Param("theater_id")

	var theater models.Theater

	if err := c.ShouldBindJSON(&theater); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := theaterService.UpdateTheater(theaterID, &theater); err != nil {

		if err.Error() == "theater not found" {
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
		"message": "Theater updated successfully",
	})
}

// DeleteTheater godoc
// @Summary Delete theater
// @Description Delete theater by ID
// @Tags Theaters
// @Produce json
// @Param theater_id path string true "Theater ID"
// @Success 200 {object} map[string]string
// @Router /theaters/{theater_id} [delete]
func DeleteTheater(c *gin.Context) {

	if theaterService == nil {
		theaterService = services.NewTheaterService()
	}

	theaterID := c.Param("theater_id")

	if err := theaterService.DeleteTheater(theaterID); err != nil {

		if err.Error() == "theater not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		if strings.Contains(err.Error(), "violates foreign key constraint") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "cannot delete theater because it is used by screens",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}
