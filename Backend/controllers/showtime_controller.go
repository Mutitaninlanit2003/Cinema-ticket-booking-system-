package controllers

import (
	"errors"
	"net/http"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/repository"
	"cinema-ticket-booking/services"

	"github.com/gin-gonic/gin"
)

type ShowtimeController struct {
	service *services.ShowtimeService
}

func NewShowtimeController(
	service *services.ShowtimeService,
) *ShowtimeController {

	return &ShowtimeController{
		service: service,
	}

}

// POST /showtimes

// CreateShowtime godoc
// @Summary Create Showtime
// @Description Create a new showtime
// @Tags Showtimes
// @Accept json
// @Produce json
// @Param showtime body models.Showtime true "Showtime"
// @Success 201 {object} map[string]string
// @Router /showtimes [post]
func (c *ShowtimeController) CreateShowtime(
	ctx *gin.Context,
) {

	var showtime models.Showtime

	if err := ctx.ShouldBindJSON(&showtime); err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	err := c.service.CreateShowtime(
		ctx,
		showtime,
	)

	if err != nil {

		ctx.JSON(
			500,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		201,
		gin.H{
			"message": "Showtime created successfully",
		},
	)

}

// GET /showtimes

// GetShowtimes godoc
// @Summary Get all showtimes
// @Description Retrieve all showtimes
// @Tags Showtimes
// @Produce json
// @Success 200 {array} models.Showtime
// @Router /showtimes [get]
func (c *ShowtimeController) GetShowtimes(
	ctx *gin.Context,
) {

	showtimes, err :=
		c.service.GetShowtimes(ctx)

	if err != nil {

		ctx.JSON(
			500,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		200,
		showtimes,
	)

}

// PUT /showtimes/:showtime_id

// UpdateShowtime godoc
// @Summary Update showtime
// @Description Update showtime by ID
// @Tags Showtimes
// @Accept json
// @Produce json
// @Param showtime_id path string true "Showtime ID"
// @Param showtime body models.Showtime true "Showtime"
// @Success 200 {object} map[string]string
// @Router /showtimes/{showtime_id} [put]
func (c *ShowtimeController) UpdateShowtime(
	ctx *gin.Context,
) {

	showtimeID := ctx.Param("showtime_id")

	var showtime models.Showtime

	if err := ctx.ShouldBindJSON(&showtime); err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	err := c.service.UpdateShowtime(
		ctx,
		showtimeID,
		showtime,
	)

	if err != nil {

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Showtime updated successfully",
		},
	)
}

// DELETE /showtimes/:showtime_id

// DeleteShowtime godoc
// @Summary Delete showtime
// @Description Delete showtime by ID
// @Tags Showtimes
// @Produce json
// @Param showtime_id path string true "Showtime ID"
// @Success 200 {object} map[string]string
// @Router /showtimes/{showtime_id} [delete]
func (c *ShowtimeController) DeleteShowtime(
	ctx *gin.Context,
) {

	showtimeID := ctx.Param("showtime_id")

	err := c.service.DeleteShowtime(
		ctx,
		showtimeID,
	)

	if err != nil {

		if errors.Is(err, repository.ErrShowtimeNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, repository.ErrShowtimeHasBookings) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Showtime deleted successfully",
		},
	)
}
