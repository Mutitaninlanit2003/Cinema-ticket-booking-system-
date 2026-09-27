package repository

import (
	"context"

	"cinema-ticket-booking/models"

	"gorm.io/gorm"

	"errors"
)

type AuditLogRepository struct {
	DB *gorm.DB
}

func NewAuditLogRepository(
	db *gorm.DB,
) *AuditLogRepository {
	return &AuditLogRepository{
		DB: db,
	}
}

func (r *AuditLogRepository) Create(
	ctx context.Context,
	auditLog models.AuditLog,
) error {

	return r.DB.WithContext(ctx).Create(&auditLog).Error
}

func (r *AuditLogRepository) FindAll(
	ctx context.Context,
) ([]models.AuditLog, error) {

	var logs []models.AuditLog

	err := r.DB.WithContext(ctx).Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// Update Audit Log
func (r *AuditLogRepository) Update(
	ctx context.Context,
	auditLogID string,
	auditLog models.AuditLog,
) error {

	result := r.DB.WithContext(ctx).
		Model(&models.AuditLog{}).
		Where("audit_log_id = ?", auditLogID).
		Updates(auditLog)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("audit log not found")
	}

	return nil
}

// Delete Audit Log
func (r *AuditLogRepository) Delete(
	ctx context.Context,
	auditLogID string,
) error {

	result := r.DB.WithContext(ctx).
		Where("audit_log_id = ?", auditLogID).
		Delete(&models.AuditLog{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("audit log not found")
	}

	return nil
}
