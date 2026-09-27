package controllers

import (
	"net/http"

	"cinema-ticket-booking/services"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	service *services.AdminService
}

func NewAdminController(
	service *services.AdminService,
) *AdminController {
	return &AdminController{
		service: service,
	}
}

// GetBookings godoc
// @Summary Admin - Get bookings
// @Description Get all bookings with optional filters
// @Tags Admin
// @Produce json
// @Param user_id query string false "Filter by User ID"
// @Param status query string false "Filter by booking status"
// @Success 200 {array} models.Booking
// @Failure 500 {object} map[string]string
// @Router /admin/bookings [get]
func (c *AdminController) GetBookings(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	status := ctx.Query("status")

	bookings, err := c.service.GetBookings(
		ctx,
		&userID,
		&status,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": bookings,
	})
}

// GetMe godoc
// @Summary Admin - Get profile
// @Description Get the profile of the currently logged-in admin
// @Tags Admin
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/me [get]
func (c *AdminController) GetMe(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("user_id")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "user is not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(string)

	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})
		return
	}

	user, err := c.service.GetMe(ctx, userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "admin not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// UpdateMe godoc
// @Summary Admin - Update profile
// @Description Update the profile of the currently logged-in admin
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body UpdateAdminProfileRequest true "Admin profile"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/me [put]
type UpdateAdminProfileRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func (c *AdminController) UpdateMe(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("user_id")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "user is not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(string)

	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})
		return
	}

	var request UpdateAdminProfileRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := c.service.UpdateMe(
		ctx,
		userID,
		request.Name,
		request.Email,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "admin not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
