package services

import (
	"context"
	"errors"
	"time"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/repository"
)

type PaymentService struct {
	repo    *repository.PaymentRepository
	booking *BookingService
}

func NewPaymentService(
	repo *repository.PaymentRepository,
	booking *BookingService,
) *PaymentService {

	return &PaymentService{
		repo:    repo,
		booking: booking,
	}
}

// Create Payment
func (s *PaymentService) CreatePayment(
	ctx context.Context,
	payment *models.Payment,
) error {

	if payment == nil {
		return errors.New("payment data is required")
	}

	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()

	// สร้าง Payment
	if err := s.repo.Create(ctx, payment); err != nil {
		return err
	}

	// ถ้าชำระเงินสำเร็จ
	if payment.Status == "paid" {

		if s.booking == nil {
			return errors.New("booking service is not initialized")
		}

		// เปลี่ยน Booking เป็น confirmed
		if err := s.booking.ConfirmBooking(
			ctx,
			payment.BookingID,
		); err != nil {
			return err
		}
	}

	return nil
}

// Get All Payments
func (s *PaymentService) GetPayments(
	ctx context.Context,
) ([]*models.Payment, error) {

	return s.repo.FindAll(ctx)
}

func (s *PaymentService) UpdatePayment(
	ctx context.Context,
	paymentID string,
	payment *models.Payment,
) error {

	if payment == nil {
		return errors.New("payment data is required")
	}

	payment.UpdatedAt = time.Now()

	return s.repo.Update(
		ctx,
		paymentID,
		payment,
	)
}

// Delete Payment
func (s *PaymentService) DeletePayment(
	ctx context.Context,
	paymentID string,
) error {

	return s.repo.Delete(
		ctx,
		paymentID,
	)
}
