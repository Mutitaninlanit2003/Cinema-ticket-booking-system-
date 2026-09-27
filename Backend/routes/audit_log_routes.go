package routes

import (
	"cinema-ticket-booking/controllers"
	"fmt"

	"github.com/gin-gonic/gin"
)

func AuditLogRoutes(
	router *gin.Engine,
	controller *controllers.AuditLogController,
) {

	auditLogs := router.Group("/audit_logs")

	{
		auditLogs.POST("", controller.CreateAuditLog)
		auditLogs.GET("", controller.GetAuditLogs)
		auditLogs.PUT("/:audit_log_id", controller.UpdateAuditLog)
		auditLogs.DELETE("/:audit_log_id", controller.DeleteAuditLog)
	}

	fmt.Println("✅ AuditLogRoutes Loaded")
}
