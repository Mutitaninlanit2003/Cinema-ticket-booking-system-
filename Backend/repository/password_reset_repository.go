package repository

import (
	"errors"
	"time"

	"cinema-ticket-booking/config"
	"cinema-ticket-booking/models"

	"gorm.io/gorm"
)

type PasswordResetRepository struct {
	DB *gorm.DB
}

func NewPasswordResetRepository() *PasswordResetRepository {
	return &PasswordResetRepository{
		DB: config.DB,
	}
}

// --------------------------------------------------
// Create Reset Token
// --------------------------------------------------

// Create saves a new password reset token to the database.
func (r *PasswordResetRepository) Create(
	resetToken *models.PasswordResetToken,
) error {

	return r.DB.Create(resetToken).Error
}

// --------------------------------------------------
// Find Reset Token
// --------------------------------------------------

// FindByToken searches for a password reset token.
//
// The token must:
// 1. Exist in the database
// 2. Not be expired
func (r *PasswordResetRepository) FindByToken(
	token string,
) (*models.PasswordResetToken, error) {

	var resetToken models.PasswordResetToken

	err := r.DB.
		Where(
			"token = ? AND expires_at > ?",
			token,
			time.Now(),
		).
		First(&resetToken).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(
				"invalid or expired reset token",
			)
		}

		return nil, err
	}

	return &resetToken, nil
}

// --------------------------------------------------
// Delete Reset Token
// --------------------------------------------------

// DeleteByToken removes a reset token after it has been used.
func (r *PasswordResetRepository) DeleteByToken(
	token string,
) error {

	result := r.DB.
		Where("token = ?", token).
		Delete(&models.PasswordResetToken{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("reset token not found")
	}

	return nil
}

// --------------------------------------------------
// Delete User's Existing Reset Tokens
// --------------------------------------------------

// DeleteByUserID removes all existing reset tokens
// belonging to the specified user.
//
// This prevents multiple active reset links from
// remaining valid at the same time.
func (r *PasswordResetRepository) DeleteByUserID(
	userID string,
) error {

	return r.DB.
		Where("user_id = ?", userID).
		Delete(&models.PasswordResetToken{}).
		Error
}
