package controllers

import (
	"errors"
	"net/http"

	"cinema-ticket-booking/services"

	"github.com/gin-gonic/gin"
)

type SeatLockController struct {
	seatLockService *services.SeatLockService
	seatMapService  *services.SeatMapService
}

func NewSeatLockController(
	seatLockService *services.SeatLockService,
	seatMapService *services.SeatMapService,
) *SeatLockController {
	return &SeatLockController{
		seatLockService: seatLockService,
		seatMapService:  seatMapService,
	}
}

type LockSeatRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
}

func (c *SeatLockController) LockSeat(ctx *gin.Context) {

	showtimeID := ctx.Param("showtime_id")
	seatID := ctx.Param("seat_id")

	var req LockSeatRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := c.seatLockService.LockSeat(
		ctx.Request.Context(),
		showtimeID,
		seatID,
		req.UserID,
	)

	if err != nil {

		if errors.Is(err, services.ErrSeatLocked) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error":  "seat is already locked",
				"status": "LOCKED",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// แจ้งผู้ใช้ทุกคนที่กำลังดู Showtime นี้
	c.seatMapService.Broadcast(
		showtimeID,
		services.SeatMapMessage{
			Type:       "seat_status",
			ShowtimeID: showtimeID,
			SeatID:     seatID,
			Status:     "LOCKED",
			UserID:     req.UserID,
			Message:    "Seat locked",
		},
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "seat locked successfully",
		"seat_id":    seatID,
		"status":     "LOCKED",
		"expires_in": "15 minutes",
	})
}

func (c *SeatLockController) UnlockSeat(ctx *gin.Context) {

	showtimeID := ctx.Param("showtime_id")
	seatID := ctx.Param("seat_id")

	var req LockSeatRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := c.seatLockService.UnlockSeat(
		ctx.Request.Context(),
		showtimeID,
		seatID,
		req.UserID,
	)

	if err != nil {

		if errors.Is(err, services.ErrSeatLockOwner) {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "you do not own this seat lock",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.seatMapService.Broadcast(
		showtimeID,
		services.SeatMapMessage{
			Type:       "seat_status",
			ShowtimeID: showtimeID,
			SeatID:     seatID,
			Status:     "AVAILABLE",
			UserID:     req.UserID,
			Message:    "Seat available",
		},
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "seat unlocked successfully",
		"seat_id": seatID,
		"status":  "AVAILABLE",
	})
}
