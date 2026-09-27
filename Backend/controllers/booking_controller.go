package controllers

import (
	"errors"
	"net/http"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/services"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type BookingController struct {
	service *services.BookingService
}

// UpdateBookingInput ใช้เฉพาะตอนแก้ไข/ยกเลิก Booking
// ไม่มี field ไหนบังคับ required เพราะอาจแก้แค่ status อย่างเดียว
type UpdateBookingInput struct {
	UserID      *string         `json:"user_id"`
	ShowtimeID  *string         `json:"showtime_id"`
	SeatNumbers *pq.StringArray `json:"seat_numbers" swaggertype:"array,string"`
	TotalPrice  *float64        `json:"total_price"`
	Status      *string         `json:"status"`
}

func NewBookingController(
	service *services.BookingService,
) *BookingController {
	return &BookingController{
		service: service,
	}
}

// ------------------------------------------------------------
// Create Booking
// POST /bookings
// ------------------------------------------------------------

func (c *BookingController) CreateBooking(ctx *gin.Context) {

	var booking models.Booking

	if err := ctx.ShouldBindJSON(&booking); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ทุก booking ที่สร้างใหม่ต้องเริ่มเป็น pending
	status := "pending"
	booking.Status = &status

	if err := c.service.CreateBooking(
		ctx.Request.Context(),
		&booking,
	); err != nil {

		if errors.Is(err, services.ErrSeatLocked) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Booking created successfully",
	})
}

// ------------------------------------------------------------
// Get Bookings
// GET /bookings
// ------------------------------------------------------------

func (c *BookingController) GetBookings(ctx *gin.Context) {

	bookings, err := c.service.GetBookings(
		ctx.Request.Context(),
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, bookings)
}

// ------------------------------------------------------------
// Update Booking
// PUT /bookings/:booking_id
// ------------------------------------------------------------

func (c *BookingController) UpdateBooking(ctx *gin.Context) {

	bookingID := ctx.Param("booking_id")

	if bookingID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "booking_id is required",
		})
		return
	}

	var input UpdateBookingInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	booking := models.Booking{
		UserID:      input.UserID,
		ShowtimeID:  input.ShowtimeID,
		SeatNumbers: input.SeatNumbers,
		TotalPrice:  input.TotalPrice,
		Status:      input.Status,
	}

	if err := c.service.UpdateBooking(
		ctx.Request.Context(),
		bookingID,
		&booking,
	); err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Booking updated successfully",
	})
}

// ------------------------------------------------------------
// Delete Booking
// DELETE /bookings/:booking_id
// ------------------------------------------------------------

func (c *BookingController) DeleteBooking(ctx *gin.Context) {

	bookingID := ctx.Param("booking_id")

	if bookingID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "booking_id is required",
		})
		return
	}

	if err := c.service.DeleteBooking(
		ctx.Request.Context(),
		bookingID,
	); err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Booking deleted successfully",
	})
}
