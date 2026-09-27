package controllers

import (
	"errors"
	"net/http"
	"strings"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var seatService *services.SeatService

// GetSeats godoc
// @Summary Get all seats
// @Description Retrieve all seats
// @Tags Seats
// @Produce json
// @Success 200 {array} models.Seat
// @Router /seats [get]
func GetSeats(c *gin.Context) {

	if seatService == nil {
		seatService = services.NewSeatService()
	}

	seats, err := seatService.GetSeats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, seats)
}

// CreateSeat godoc
// @Summary Create Seat
// @Description Create a new seat
// @Tags Seats
// @Accept json
// @Produce json
// @Param seat body models.Seat true "Seat"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /seats [post]
func CreateSeat(c *gin.Context) {

	if seatService == nil {
		seatService = services.NewSeatService()
	}

	// IMPORTANT:
	// Service.CreateSeat expects models.Seat
	// not *models.Seat
	var seat models.Seat

	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := seatService.CreateSeat(seat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Seat created successfully",
	})
}

// UpdateSeat godoc
// @Summary Update seat
// @Description Update seat by ID
// @Tags Seats
// @Accept json
// @Produce json
// @Param seat_id path string true "Seat ID"
// @Param seat body services.UpdateSeatRequest true "Update seat request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /seats/{seat_id} [put]
func UpdateSeat(c *gin.Context) {

	if seatService == nil {
		seatService = services.NewSeatService()
	}

	seatID := c.Param("seat_id")

	// IMPORTANT:
	// Use value here.
	// We will pass &req to the service.
	var req services.UpdateSeatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := seatService.UpdateSeat(
		seatID,
		&req,
	); err != nil {

		// Screen not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "screen not found",
			})
			return
		}

		// Seat not found
		if err.Error() == "seat not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "seat not found",
			})
			return
		}

		// Validation error
		if strings.Contains(
			err.Error(),
			"invalid seat_type",
		) ||
			strings.Contains(
				err.Error(),
				"seat_number is required",
			) {

			c.JSON(http.StatusBadRequest, gin.H{
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
		"message": "Seat updated successfully",
	})
}

// DeleteSeat godoc
// @Summary Delete seat
// @Description Delete seat by ID
// @Tags Seats
// @Produce json
// @Param seat_id path string true "Seat ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /seats/{seat_id} [delete]
func DeleteSeat(c *gin.Context) {

	if seatService == nil {
		seatService = services.NewSeatService()
	}

	seatID := c.Param("seat_id")

	if err := seatService.DeleteSeat(seatID); err != nil {

		if err.Error() == "seat not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "seat not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Seat deleted successfully",
	})
}

// GenerateSeats godoc
// @Summary Generate multiple seats
// @Description Generate multiple seats for a screen
// @Tags Seats
// @Accept json
// @Produce json
// @Param request body services.GenerateSeatsRequest true "Generate seats request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /seats/generate [post]
func GenerateSeats(c *gin.Context) {

	if seatService == nil {
		seatService = services.NewSeatService()
	}

	// IMPORTANT:
	// Use value here and pass &req to service.
	var req services.GenerateSeatsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	count, err := seatService.GenerateSeats(&req)

	if err != nil {

		// Screen not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "screen not found",
			})
			return
		}

		// Duplicate seat
		if strings.Contains(
			err.Error(),
			"already exists",
		) {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Validation errors
		if strings.Contains(
			err.Error(),
			"start_row",
		) ||
			strings.Contains(
				err.Error(),
				"end_row",
			) ||
			strings.Contains(
				err.Error(),
				"row must be between A and Z",
			) {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Seats generated successfully",
		"count":   count,
	})
}
