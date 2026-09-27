package services

import (
	"context"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/repository"
)

type AdminService struct {
	repo *repository.AdminRepository
}

func NewAdminService(
	repo *repository.AdminRepository,
) *AdminService {
	return &AdminService{
		repo: repo,
	}
}

func (s *AdminService) GetBookings(
	ctx context.Context,
	userID *string,
	status *string,
) ([]*models.Booking, error) {

	return s.repo.FindBookings(
		ctx,
		userID,
		status,
	)
}

// GetMe gets the currently logged-in admin profile.
func (s *AdminService) GetMe(
	ctx context.Context,
	userID string,
) (*models.User, error) {

	return s.repo.GetMe(ctx, userID)
}

// UpdateMe updates the currently logged-in admin profile.
func (s *AdminService) UpdateMe(
	ctx context.Context,
	userID string,
	name string,
	email string,
) (*models.User, error) {
	return s.repo.UpdateMe(
		ctx,
		userID,
		name,
		email,
	)
}
