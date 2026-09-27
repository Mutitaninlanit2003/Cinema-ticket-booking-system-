package repository

import (
	"context"

	"cinema-ticket-booking/models"

	"gorm.io/gorm"
)

type AdminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(
	db *gorm.DB,
) *AdminRepository {
	return &AdminRepository{
		DB: db,
	}
}

func (r *AdminRepository) FindBookings(
	ctx context.Context,
	userID *string,
	status *string,
) ([]*models.Booking, error) {

	var bookings []*models.Booking

	query := r.DB.WithContext(ctx).
		Model(&models.Booking{})

	if userID != nil && *userID != "" {
		query = query.Where("user_id = ?", *userID)
	}

	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	if err := query.
		Order("created_at DESC").
		Find(&bookings).Error; err != nil {
		return nil, err
	}

	return bookings, nil
}

// GetMe gets an admin by user_id.
func (r *AdminRepository) GetMe(
	ctx context.Context,
	userID string,
) (*models.User, error) {

	var user models.User

	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("role = ?", "admin").
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateMe updates the currently logged-in admin profile.
func (r *AdminRepository) UpdateMe(
	ctx context.Context,
	userID string,
	name string,
	email string,
) (*models.User, error) {

	var user models.User

	// Find the currently logged-in admin.
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("role = ?", "admin").
		First(&user).Error

	if err != nil {
		return nil, err
	}

	// Update profile information.
	user.Name = name
	user.Email = email

	if err := r.DB.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
