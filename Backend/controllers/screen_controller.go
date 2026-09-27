package controllers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/repository"
	"cinema-ticket-booking/services"

	"github.com/gin-gonic/gin"
)

var screenService *services.ScreenService

// GetScreens godoc
// @Summary Get all screens
// @Description Retrieve all screens
// @Tags Screens
// @Produce json
// @Success 200 {array} models.Screen
// @Router /screens [get]
func GetScreens(c *gin.Context) {

	if screenService == nil {
		screenService = services.NewScreenService()
	}

	screens, err := screenService.GetScreens()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, screens)
}

// CreateScreen godoc
// @Summary Create Screen
// @Description Create a new screen
// @Tags Screens
// @Accept json
// @Produce json
// @Param screen body models.Screen true "Screen"
// @Success 201 {object} map[string]string
// @Router /screens [post]
func CreateScreen(c *gin.Context) {

	if screenService == nil {
		screenService = services.NewScreenService()
	}

	var screen models.Screen

	if err := c.ShouldBindJSON(&screen); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := screenService.CreateScreen(screen); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Screen created successfully",
	})
}

// UpdateScreen godoc
// @Summary Update screen
// @Description Update screen by ID
// @Tags Screens
// @Accept json
// @Produce json
// @Param screen_id path string true "Screen ID"
// @Param screen body models.Screen true "Screen"
// @Success 200 {object} map[string]string
// @Router /screens/{screen_id} [put]
func UpdateScreen(c *gin.Context) {

	if screenService == nil {
		screenService = services.NewScreenService()
	}

	screenID := c.Param("screen_id")

	var screen models.Screen

	if err := c.ShouldBindJSON(&screen); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err := screenService.UpdateScreen(screenID, screen); err != nil {

		if errors.Is(err, repository.ErrScreenNotFound) {
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
		"message": "Screen updated successfully",
	})
}

// DeleteScreen godoc
// @Summary Delete screen
// @Description Delete screen by ID. Use ?force=true to also delete its showtimes, bookings and payments (dev/test data only).
// @Tags Screens
// @Produce json
// @Param screen_id path string true "Screen ID"
// @Param force query bool false "Force delete with all related data"
// @Success 200 {object} map[string]string
// @Router /screens/{screen_id} [delete]
func DeleteScreen(c *gin.Context) {

	if screenService == nil {
		screenService = services.NewScreenService()
	}

	screenID := c.Param("screen_id")

	var err error
	if c.Query("force") == "true" {
		err = screenService.ForceDeleteScreen(screenID)
	} else {
		err = screenService.DeleteScreen(screenID)
	}

	if err != nil {

		// screen ไม่พบ -> 404
		if errors.Is(err, repository.ErrScreenNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		// ยังมี showtimes ใช้อยู่ -> 409
		if errors.Is(err, repository.ErrScreenInUse) {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		// fallback: FK จากตารางอื่นที่ยังไม่ได้จัดการ -> 409
		// error ของ Postgres ระบุชื่อตารางที่ขวางอยู่ เช่น `on table "seats"`
		if strings.Contains(err.Error(), "violates foreign key constraint") {

			message := "cannot delete screen because it is still in use"

			switch {
			case strings.Contains(err.Error(), `on table "seats"`):
				message = "cannot delete screen because it still has seats configured — delete its seats first"

			case strings.Contains(err.Error(), `on table "showtimes"`):
				message = "cannot delete screen because it is used by showtimes"

			case strings.Contains(err.Error(), `on table "bookings"`):
				message = "cannot delete screen because it has bookings associated with it"
			}

			c.JSON(http.StatusConflict, gin.H{
				"error": message,
			})
			return
		}

		// อย่างอื่น -> 500 และ log สาเหตุจริงไว้ดูใน terminal
		log.Println("delete screen error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Screen deleted successfully",
	})
}
