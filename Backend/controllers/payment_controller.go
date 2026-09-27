package controllers

import (
	"net/http"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/services"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	service *services.PaymentService
}

func NewPaymentController(
	service *services.PaymentService,
) *PaymentController {
	return &PaymentController{
		service: service,
	}
}

// CreatePayment godoc
// @Summary Create Payment
// @Description Create a new payment
// @Tags Payments
// @Accept json
// @Produce json
// @Param payment body models.Payment true "Payment"
// @Success 201 {object} map[string]string
// @Router /payments [post]
func (c *PaymentController) CreatePayment(ctx *gin.Context) {
	var payment models.Payment

	if err := ctx.ShouldBindJSON(&payment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if c.service == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "payment service is not initialized",
		})
		return
	}

	if err := c.service.CreatePayment(ctx, &payment); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Payment created successfully",
	})
}

// GetPayments godoc
// @Summary Get all payments
// @Description Retrieve all payments
// @Tags Payments
// @Produce json
// @Success 200 {array} models.Payment
// @Router /payments [get]
func (c *PaymentController) GetPayments(ctx *gin.Context) {
	if c.service == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "payment service is not initialized",
		})
		return
	}

	payments, err := c.service.GetPayments(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, payments)
}

// UpdatePayment godoc
// @Summary Update payment
// @Description Update payment by ID
// @Tags Payments
// @Accept json
// @Produce json
// @Param payment_id path string true "Payment ID"
// @Param payment body models.Payment true "Payment"
// @Success 200 {object} map[string]string
// @Router /payments/{payment_id} [put]
func (c *PaymentController) UpdatePayment(ctx *gin.Context) {
	paymentID := ctx.Param("payment_id")

	var payment models.Payment

	if err := ctx.ShouldBindJSON(&payment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if c.service == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "payment service is not initialized",
		})
		return
	}

	if err := c.service.UpdatePayment(
		ctx,
		paymentID,
		&payment,
	); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Payment updated successfully",
	})
}

// DeletePayment godoc
// @Summary Delete payment
// @Description Delete payment by ID
// @Tags Payments
// @Produce json
// @Param payment_id path string true "Payment ID"
// @Success 200 {object} map[string]string
// @Router /payments/{payment_id} [delete]
func (c *PaymentController) DeletePayment(ctx *gin.Context) {
	paymentID := ctx.Param("payment_id")

	if c.service == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "payment service is not initialized",
		})
		return
	}

	if err := c.service.DeletePayment(
		ctx,
		paymentID,
	); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Payment deleted successfully",
	})
}
