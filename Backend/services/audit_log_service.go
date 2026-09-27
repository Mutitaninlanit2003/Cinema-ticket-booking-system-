package services

import (
	"context"
	"time"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/repository"
)

type AuditLogService struct {
	repo *repository.AuditLogRepository
}

func NewAuditLogService(
	repo *repository.AuditLogRepository,
) *AuditLogService {

	return &AuditLogService{
		repo: repo,
	}
}

func (s *AuditLogService) CreateAuditLog(
	ctx context.Context,
	log models.AuditLog,
) error {

	log.CreatedAt = time.Now()

	return s.repo.Create(
		ctx,
		log,
	)
}

func (s *AuditLogService) GetAuditLogs(
	ctx context.Context,
) ([]models.AuditLog, error) {

	return s.repo.FindAll(ctx)
}
func (s *AuditLogService) UpdateAuditLog(
	ctx context.Context,
	auditLogID string,
	log models.AuditLog,
) error {

	return s.repo.Update(
		ctx,
		auditLogID,
		log,
	)
}

func (s *AuditLogService) DeleteAuditLog(
	ctx context.Context,
	auditLogID string,
) error {

	return s.repo.Delete(
		ctx,
		auditLogID,
	)
}
