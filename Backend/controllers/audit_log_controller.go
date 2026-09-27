package controllers

import (
	"net/http"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/services"

	"github.com/gin-gonic/gin"
)

type AuditLogController struct {
	service *services.AuditLogService
}

func NewAuditLogController(
	service *services.AuditLogService,
) *AuditLogController {

	return &AuditLogController{
		service: service,
	}
}

// CreateAuditLog godoc
// @Summary Create Audit Log
// @Description Create a new audit log
// @Tags Audit Logs
// @Accept json
// @Produce json
// @Param audit_log body models.AuditLog true "Audit Log"
// @Success 201 {object} map[string]string
// @Router /audit_logs [post]
func (c *AuditLogController) CreateAuditLog(
	ctx *gin.Context,
) {

	var auditLog models.AuditLog

	if err := ctx.ShouldBindJSON(&auditLog); err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	err := c.service.CreateAuditLog(
		ctx,
		auditLog,
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
		http.StatusCreated,
		gin.H{
			"message": "Audit log created successfully",
		},
	)
}

// GetAuditLogs godoc
// @Summary Get all audit logs
// @Description Retrieve all audit logs
// @Tags Audit Logs
// @Produce json
// @Success 200 {array} models.AuditLog
// @Router /audit_logs [get]
func (c *AuditLogController) GetAuditLogs(
	ctx *gin.Context,
) {

	logs, err := c.service.GetAuditLogs(ctx)

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
		logs,
	)
}

// PUT /audit_logs/:audit_log_id

// UpdateAuditLog godoc
// @Summary Update audit log
// @Description Update audit log by ID
// @Tags Audit Logs
// @Accept json
// @Produce json
// @Param audit_log_id path string true "Audit Log ID"
// @Param audit_log body models.AuditLog true "Audit Log"
// @Success 200 {object} map[string]string
// @Router /audit_logs/{audit_log_id} [put]
func (c *AuditLogController) UpdateAuditLog(
	ctx *gin.Context,
) {

	auditLogID := ctx.Param("audit_log_id")

	var auditLog models.AuditLog

	if err := ctx.ShouldBindJSON(&auditLog); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := c.service.UpdateAuditLog(
		ctx,
		auditLogID,
		auditLog,
	); err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Audit log updated successfully",
	})
}

// DELETE /audit_logs/:audit_log_id

// DeleteAuditLog godoc
// @Summary Delete audit log
// @Description Delete audit log by ID
// @Tags Audit Logs
// @Produce json
// @Param audit_log_id path string true "Audit Log ID"
// @Success 200 {object} map[string]string
// @Router /audit_logs/{audit_log_id} [delete]
func (c *AuditLogController) DeleteAuditLog(
	ctx *gin.Context,
) {

	auditLogID := ctx.Param("audit_log_id")

	if err := c.service.DeleteAuditLog(
		ctx,
		auditLogID,
	); err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Audit log deleted successfully",
	})
}
