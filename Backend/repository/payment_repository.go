package repository

import (
	"context"
	"errors"

	"cinema-ticket-booking/models"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		DB: db,
	}
}

// Create Payment
func (r *PaymentRepository) Create(
	ctx context.Context,
	payment *models.Payment,
) error {

	if payment == nil {
		return errors.New("payment data is required")
	}

	var booking models.Booking

	if err := r.DB.WithContext(ctx).
		Where("booking_id = ?", payment.BookingID).
		First(&booking).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("booking not found")
		}

		return err
	}

	if booking.ID == nil {
		return errors.New("booking id is missing")
	}

	payment.BookingRefID = *booking.ID

	return r.DB.WithContext(ctx).Create(payment).Error
}

// Get All Payments
func (r *PaymentRepository) FindAll(
	ctx context.Context,
) ([]*models.Payment, error) {

	var payments []*models.Payment

	err := r.DB.WithContext(ctx).
		Find(&payments).Error

	return payments, err
}

// Update Payment
func (r *PaymentRepository) Update(
	ctx context.Context,
	paymentID string,
	payment *models.Payment,
) error {

	result := r.DB.WithContext(ctx).
		Model(&models.Payment{}).
		Where("payment_id = ?", paymentID).
		Updates(payment)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("payment not found")
	}

	return nil
}

// Delete Payment
func (r *PaymentRepository) Delete(
	ctx context.Context,
	paymentID string,
) error {

	result := r.DB.WithContext(ctx).
		Where("payment_id = ?", paymentID).
		Delete(&models.Payment{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("payment not found")
	}

	return nil
}
