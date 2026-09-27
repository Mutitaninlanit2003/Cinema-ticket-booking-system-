package services

import (
	"context"
	"errors"
	"time"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/repository"
)

type BookingService struct {
	repo     *repository.BookingRepository
	seatLock *SeatLockService
	mq       *MessageQueueService
}

func NewBookingService(
	repo *repository.BookingRepository,
	seatLock *SeatLockService,
	mq *MessageQueueService,
) *BookingService {
	return &BookingService{
		repo:     repo,
		seatLock: seatLock,
		mq:       mq,
	}
}

// ------------------------------------------------------------
// Create Booking
// ------------------------------------------------------------

func (s *BookingService) CreateBooking(
	ctx context.Context,
	booking *models.Booking,
) error {

	if booking == nil {
		return errors.New("booking data is required")
	}

	if s.seatLock == nil {
		return errors.New("seat lock service is not initialized")
	}

	if s.repo == nil {
		return errors.New("booking repository is not initialized")
	}

	if booking.UserID == nil || *booking.UserID == "" {
		return errors.New("user_id is required")
	}

	if booking.ShowtimeID == nil || *booking.ShowtimeID == "" {
		return errors.New("showtime_id is required")
	}

	if booking.SeatNumbers == nil || len(*booking.SeatNumbers) == 0 {
		return errors.New("at least one seat is required")
	}

	// Booking ใหม่ต้องเป็น pending
	status := "pending"
	booking.Status = &status

	// --------------------------------------------------------
	// Lock seats ใน Redis
	// --------------------------------------------------------

	lockedSeats := make([]string, 0, len(*booking.SeatNumbers))

	for _, seatNumber := range *booking.SeatNumbers {

		err := s.seatLock.LockSeat(
			ctx,
			*booking.ShowtimeID,
			seatNumber,
			*booking.UserID,
		)

		if err != nil {

			// ถ้า lock บางที่นั่งสำเร็จไปแล้ว
			// ให้ปลด lock ทั้งหมดก่อนหน้า
			for _, lockedSeat := range lockedSeats {
				_ = s.seatLock.UnlockSeat(
					ctx,
					*booking.ShowtimeID,
					lockedSeat,
					*booking.UserID,
				)
			}

			if errors.Is(err, ErrSeatLocked) {
				return errors.New(
					"seat " + seatNumber + " is already locked",
				)
			}

			return err
		}

		lockedSeats = append(lockedSeats, seatNumber)
	}

	// --------------------------------------------------------
	// Timestamp
	// --------------------------------------------------------

	now := time.Now()

	booking.CreatedAt = &now
	booking.UpdatedAt = &now

	// --------------------------------------------------------
	// Save booking
	// --------------------------------------------------------

	if err := s.repo.Create(ctx, booking); err != nil {

		// ถ้าบันทึก Booking ไม่สำเร็จ
		// ต้องปลด Redis lock
		for _, seatNumber := range lockedSeats {
			_ = s.seatLock.UnlockSeat(
				ctx,
				*booking.ShowtimeID,
				seatNumber,
				*booking.UserID,
			)
		}

		return err
	}

	return nil
}

// ------------------------------------------------------------
// Get Bookings
// ------------------------------------------------------------

func (s *BookingService) GetBookings(
	ctx context.Context,
) ([]*models.Booking, error) {

	if s.repo == nil {
		return nil, errors.New(
			"booking repository is not initialized",
		)
	}

	return s.repo.FindAll(ctx)
}

// ------------------------------------------------------------
// Update Booking
//
// ใช้สำหรับแก้ไขสถานะ Booking
// เช่น
//
// pending -> cancelled
//
// ------------------------------------------------------------

func (s *BookingService) UpdateBooking(
	ctx context.Context,
	bookingID string,
	booking *models.Booking,
) error {

	if booking == nil {
		return errors.New("booking data is required")
	}

	if bookingID == "" {
		return errors.New("booking_id is required")
	}

	if s.repo == nil {
		return errors.New(
			"booking repository is not initialized",
		)
	}

	// ========================================================
	// CANCEL BOOKING
	// ========================================================

	if booking.Status != nil &&
		*booking.Status == "cancelled" {

		currentBooking, err := s.repo.FindByBookingID(
			ctx,
			bookingID,
		)

		if err != nil {
			return err
		}

		if currentBooking == nil {
			return errors.New("booking not found")
		}

		// ----------------------------------------------------
		// ถ้ายกเลิกไปแล้ว
		// ----------------------------------------------------

		if currentBooking.Status != nil &&
			*currentBooking.Status == "cancelled" {

			return errors.New(
				"booking is already cancelled",
			)
		}

		// ----------------------------------------------------
		// คืนสถานะ Seat เป็น AVAILABLE
		// ----------------------------------------------------

		if currentBooking.ShowtimeID != nil &&
			currentBooking.SeatNumbers != nil &&
			len(*currentBooking.SeatNumbers) > 0 {

			if err := s.repo.MarkSeatsAsAvailable(
				ctx,
				*currentBooking.ShowtimeID,
				[]string(*currentBooking.SeatNumbers),
			); err != nil {
				return err
			}

			// ------------------------------------------------
			// ปลด Redis Seat Lock
			// ------------------------------------------------

			if s.seatLock != nil &&
				currentBooking.UserID != nil {

				for _, seatNumber := range *currentBooking.SeatNumbers {

					_ = s.seatLock.UnlockSeat(
						ctx,
						*currentBooking.ShowtimeID,
						seatNumber,
						*currentBooking.UserID,
					)
				}
			}
		}

		// ----------------------------------------------------
		// Update Booking -> cancelled
		// ----------------------------------------------------

		now := time.Now()
		booking.UpdatedAt = &now

		return s.repo.UpdateStatus(
			ctx,
			bookingID,
			"cancelled",
		)
	}

	// ========================================================
	// NORMAL UPDATE
	// ========================================================

	now := time.Now()
	booking.UpdatedAt = &now

	return s.repo.Update(
		ctx,
		bookingID,
		booking,
	)
}

// ------------------------------------------------------------
// Delete Booking
// ------------------------------------------------------------

func (s *BookingService) DeleteBooking(
	ctx context.Context,
	bookingID string,
) error {

	if bookingID == "" {
		return errors.New("booking_id is required")
	}

	if s.repo == nil {
		return errors.New(
			"booking repository is not initialized",
		)
	}

	return s.repo.Delete(
		ctx,
		bookingID,
	)
}

// ------------------------------------------------------------
// Confirm Booking
//
// Payment สำเร็จ
// -> ตรวจสอบ Seat Lock
// -> Mark Seat BOOKED
// -> Booking confirmed
// -> ลบ Redis lock
//
// ------------------------------------------------------------

func (s *BookingService) ConfirmBooking(
	ctx context.Context,
	bookingID string,
) error {

	if s.repo == nil {
		return errors.New(
			"booking repository is not initialized",
		)
	}

	if s.seatLock == nil {
		return errors.New(
			"seat lock service is not initialized",
		)
	}

	if bookingID == "" {
		return errors.New("booking_id is required")
	}

	// --------------------------------------------------------
	// หา Booking
	// --------------------------------------------------------

	booking, err := s.repo.FindByBookingID(
		ctx,
		bookingID,
	)

	if err != nil {
		return err
	}

	if booking == nil {
		return errors.New("booking not found")
	}

	// --------------------------------------------------------
	// ตรวจสอบสถานะ
	// --------------------------------------------------------

	if booking.Status != nil &&
		*booking.Status == "confirmed" {

		return errors.New(
			"booking is already confirmed",
		)
	}

	if booking.Status != nil &&
		*booking.Status == "cancelled" {

		return errors.New(
			"booking is already cancelled",
		)
	}

	if booking.UserID == nil ||
		*booking.UserID == "" {

		return errors.New(
			"booking user_id is missing",
		)
	}

	if booking.ShowtimeID == nil ||
		*booking.ShowtimeID == "" {

		return errors.New(
			"booking showtime_id is missing",
		)
	}

	if booking.SeatNumbers == nil ||
		len(*booking.SeatNumbers) == 0 {

		return errors.New(
			"booking seat numbers are missing",
		)
	}

	// --------------------------------------------------------
	// ตรวจสอบ Redis Seat Lock
	//
	// ถ้า lock หมดอายุ
	// จะพยายาม lock ใหม่ให้ user เดิม
	// --------------------------------------------------------

	for _, seatNumber := range *booking.SeatNumbers {

		err := s.seatLock.VerifySeatLock(
			ctx,
			*booking.ShowtimeID,
			seatNumber,
			*booking.UserID,
		)

		if err != nil {

			// ------------------------------------------------
			// Lock หมดอายุ
			// ------------------------------------------------

			if errors.Is(
				err,
				ErrSeatLockExpired,
			) {

				relockErr := s.seatLock.LockSeat(
					ctx,
					*booking.ShowtimeID,
					seatNumber,
					*booking.UserID,
				)

				if relockErr == nil {
					continue
				}

				// ------------------------------------------------
				// มีคนอื่นเอาที่นั่งไปแล้ว
				// ------------------------------------------------

				if errors.Is(
					relockErr,
					ErrSeatLocked,
				) {

					return errors.New(
						"seat " +
							seatNumber +
							" has just been taken by someone else, please select the seat again",
					)
				}

				return relockErr
			}

			return err
		}
	}

	// --------------------------------------------------------
	// เปลี่ยน Seat -> BOOKED
	// --------------------------------------------------------

	seatNumbers := []string(
		*booking.SeatNumbers,
	)

	if err := s.repo.MarkSeatsAsBooked(
		ctx,
		*booking.ShowtimeID,
		seatNumbers,
	); err != nil {
		return err
	}

	// --------------------------------------------------------
	// เปลี่ยน Booking -> confirmed
	// --------------------------------------------------------

	confirmedStatus := "confirmed"

	booking.Status = &confirmedStatus

	now := time.Now()
	booking.UpdatedAt = &now

	if err := s.repo.UpdateStatus(
		ctx,
		bookingID,
		confirmedStatus,
	); err != nil {
		return err
	}

	// --------------------------------------------------------
	// ลบ Redis lock หลังจาก Booking confirmed
	// --------------------------------------------------------

	for _, seatNumber := range seatNumbers {

		_ = s.seatLock.UnlockSeat(
			ctx,
			*booking.ShowtimeID,
			seatNumber,
			*booking.UserID,
		)
	}

	return nil
}
