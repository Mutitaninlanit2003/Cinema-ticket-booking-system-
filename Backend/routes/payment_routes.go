package routes

import (
	"fmt"

	"cinema-ticket-booking/controllers"

	"github.com/gin-gonic/gin"
)

func PaymentRoutes(
	router *gin.Engine,
	controller *controllers.PaymentController,
) {

	fmt.Println("✅ PaymentRoutes Loaded")

	payments := router.Group("/payments")

	payments.POST("", controller.CreatePayment)
	payments.GET("", controller.GetPayments)
	payments.PUT("/:payment_id", controller.UpdatePayment)
	payments.DELETE("/:payment_id", controller.DeletePayment)
}
