package repository

import (
	"context"
	"errors"
	"log"

	"cinema-ticket-booking/models"

	"gorm.io/gorm"
)

type BookingRepository struct {
	DB *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{
		DB: db,
	}
}

// ------------------------------------------------------------
// Create Booking
// ------------------------------------------------------------

func (r *BookingRepository) Create(
	ctx context.Context,
	booking *models.Booking,
) error {

	if booking == nil {
		return errors.New("booking data is required")
	}

	if booking.UserID == nil || *booking.UserID == "" {
		return errors.New("user_id is missing")
	}

	if booking.ShowtimeID == nil || *booking.ShowtimeID == "" {
		return errors.New("showtime_id is missing")
	}

	return r.DB.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {

			// ------------------------------------------------
			// Find User
			// ------------------------------------------------

			var user models.User

			if err := tx.
				Where("user_id = ?", *booking.UserID).
				First(&user).Error; err != nil {

				if errors.Is(
					err,
					gorm.ErrRecordNotFound,
				) {
					return errors.New(
						"user not found for user_id: " +
							*booking.UserID,
					)
				}

				return err
			}

			booking.UserRefID = &user.ID

			// ------------------------------------------------
			// Find Showtime
			// ------------------------------------------------

			var showtime models.Showtime

			if err := tx.
				Where(
					"showtime_id = ?",
					*booking.ShowtimeID,
				).
				First(&showtime).Error; err != nil {

				if errors.Is(
					err,
					gorm.ErrRecordNotFound,
				) {
					return errors.New(
						"showtime not found for showtime_id: " +
							*booking.ShowtimeID,
					)
				}

				return err
			}

			booking.ShowtimeRefID = &showtime.ID

			// ------------------------------------------------
			// Create Booking
			// ------------------------------------------------

			if err := tx.Create(booking).Error; err != nil {
				return err
			}

			if booking.ID == nil {
				return errors.New(
					"booking id was not generated",
				)
			}

			if booking.BookingID == nil ||
				*booking.BookingID == "" {

				return errors.New(
					"booking uuid was not generated",
				)
			}

			// ------------------------------------------------
			// Check seats
			// ------------------------------------------------

			if booking.SeatNumbers == nil ||
				len(*booking.SeatNumbers) == 0 {

				return nil
			}

			if showtime.ScreenRefID == nil {
				return errors.New(
					"showtime does not have a screen assigned",
				)
			}

			// ------------------------------------------------
			// Find Seats
			// ------------------------------------------------

			var seats []models.Seat

			if err := tx.
				Where(
					"screen_ref_id = ? AND seat_number IN ?",
					*showtime.ScreenRefID,
					[]string(*booking.SeatNumbers),
				).
				Find(&seats).Error; err != nil {

				return err
			}

			if len(seats) != len(*booking.SeatNumbers) {
				return errors.New(
					"one or more seat numbers not found",
				)
			}

			// ------------------------------------------------
			// Create BookingSeat
			// ------------------------------------------------

			for _, seat := range seats {

				bookingSeat := models.BookingSeat{
					BookingRefID: *booking.ID,
					BookingID:    *booking.BookingID,
					SeatRefID:    seat.ID,
					SeatID:       seat.SeatID,
				}

				if err := tx.
					Create(&bookingSeat).
					Error; err != nil {

					return err
				}
			}

			return nil
		},
	)
}

// ------------------------------------------------------------
// Find All Bookings
// ------------------------------------------------------------

func (r *BookingRepository) FindAll(
	ctx context.Context,
) ([]*models.Booking, error) {

	var bookings []*models.Booking

	err := r.DB.
		WithContext(ctx).
		Find(&bookings).
		Error

	return bookings, err
}

// ------------------------------------------------------------
// Find Booking By Booking UUID
// ------------------------------------------------------------

func (r *BookingRepository) FindByBookingID(
	ctx context.Context,
	bookingID string,
) (*models.Booking, error) {

	if bookingID == "" {
		return nil, errors.New(
			"booking_id is required",
		)
	}

	var booking models.Booking

	err := r.DB.
		WithContext(ctx).
		Where(
			"booking_id = ?",
			bookingID,
		).
		First(&booking).
		Error

	if err != nil {

		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, errors.New(
				"booking not found",
			)
		}

		return nil, err
	}

	return &booking, nil
}

// ------------------------------------------------------------
// Update Booking
// ------------------------------------------------------------

func (r *BookingRepository) Update(
	ctx context.Context,
	bookingID string,
	booking *models.Booking,
) error {

	if bookingID == "" {
		return errors.New(
			"booking_id is required",
		)
	}

	if booking == nil {
		return errors.New(
			"booking data is required",
		)
	}

	result := r.DB.
		WithContext(ctx).
		Model(&models.Booking{}).
		Where(
			"booking_id = ?",
			bookingID,
		).
		Updates(booking)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New(
			"booking not found",
		)
	}

	return nil
}

// ------------------------------------------------------------
// Update Booking Status
// ------------------------------------------------------------

func (r *BookingRepository) UpdateStatus(
	ctx context.Context,
	bookingID string,
	status string,
) error {

	if bookingID == "" {
		return errors.New(
			"booking_id is required",
		)
	}

	if status == "" {
		return errors.New(
			"status is required",
		)
	}

	result := r.DB.
		WithContext(ctx).
		Model(&models.Booking{}).
		Where(
			"booking_id = ?",
			bookingID,
		).
		Update(
			"status",
			status,
		)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New(
			"booking not found",
		)
	}

	return nil
}

// ------------------------------------------------------------
// Delete Booking
// ------------------------------------------------------------

func (r *BookingRepository) Delete(
	ctx context.Context,
	bookingID string,
) error {

	if bookingID == "" {
		return errors.New(
			"booking_id is required",
		)
	}

	return r.DB.
		WithContext(ctx).
		Transaction(
			func(tx *gorm.DB) error {

				var booking models.Booking

				if err := tx.
					Where(
						"booking_id = ?",
						bookingID,
					).
					First(&booking).
					Error; err != nil {

					if errors.Is(
						err,
						gorm.ErrRecordNotFound,
					) {
						return errors.New(
							"booking not found",
						)
					}

					return err
				}

				if booking.ID == nil {
					return errors.New(
						"booking id is missing",
					)
				}

				// --------------------------------------------
				// Delete BookingSeat first
				// --------------------------------------------

				if err := tx.
					Where(
						"booking_ref_id = ?",
						*booking.ID,
					).
					Delete(
						&models.BookingSeat{},
					).
					Error; err != nil {

					return err
				}

				// --------------------------------------------
				// Delete Booking
				// --------------------------------------------

				if err := tx.
					Delete(&booking).
					Error; err != nil {

					return err
				}

				return nil
			},
		)
}

// ------------------------------------------------------------
// Mark Seats As BOOKED
//
// Used when CONFIRMING a brand-new booking (right after payment),
// where the showtime was just used moments earlier to create the
// booking in the first place — so it should always still exist.
// Left strict on purpose: if the showtime is somehow missing here,
// that's a genuine data-integrity problem worth surfacing loudly,
// unlike the cancel path below.
// ------------------------------------------------------------

func (r *BookingRepository) MarkSeatsAsBooked(
	ctx context.Context,
	showtimeID string,
	seatNumbers []string,
) error {

	if showtimeID == "" {
		return errors.New(
			"showtime_id is required",
		)
	}

	if len(seatNumbers) == 0 {
		return errors.New(
			"seat numbers are required",
		)
	}

	// --------------------------------------------------------
	// Find Showtime
	// --------------------------------------------------------

	var showtime models.Showtime

	if err := r.DB.
		WithContext(ctx).
		Where(
			"showtime_id = ?",
			showtimeID,
		).
		First(&showtime).
		Error; err != nil {

		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return errors.New(
				"showtime not found",
			)
		}

		return err
	}

	if showtime.ScreenRefID == nil {
		return errors.New(
			"showtime does not have a screen assigned",
		)
	}

	// --------------------------------------------------------
	// Update Seat Status
	// --------------------------------------------------------

	result := r.DB.
		WithContext(ctx).
		Model(&models.Seat{}).
		Where(
			"screen_ref_id = ? AND seat_number IN ?",
			*showtime.ScreenRefID,
			seatNumbers,
		).
		Update(
			"status",
			"BOOKED",
		)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != int64(len(seatNumbers)) {
		return errors.New(
			"one or more seats could not be marked as booked",
		)
	}

	return nil
}

// ------------------------------------------------------------
// Mark Seats As AVAILABLE
//
// ใช้ตอน Cancel Booking
// BOOKED / LOCKED -> AVAILABLE
//
// If the showtime this booking pointed to has since been deleted
// (e.g. an admin removed old test showtimes), there is no screen to
// resolve the seats against — meaning there is nothing meaningful
// left to mark available anyway. Rather than hard-failing and
// blocking the booking's own status from ever changing to
// "cancelled" (leaving it stuck forever), this is now a tolerated
// no-op: log it and move on, instead of returning an error.
// ------------------------------------------------------------

func (r *BookingRepository) MarkSeatsAsAvailable(
	ctx context.Context,
	showtimeID string,
	seatNumbers []string,
) error {

	if showtimeID == "" {
		return errors.New(
			"showtime_id is required",
		)
	}

	if len(seatNumbers) == 0 {
		return errors.New(
			"seat numbers are required",
		)
	}

	// --------------------------------------------------------
	// Find Showtime
	// --------------------------------------------------------

	var showtime models.Showtime

	if err := r.DB.
		WithContext(ctx).
		Where(
			"showtime_id = ?",
			showtimeID,
		).
		First(&showtime).
		Error; err != nil {

		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			log.Printf(
				"MarkSeatsAsAvailable: showtime %s no longer exists — "+
					"skipping seat release, nothing to update",
				showtimeID,
			)

			return nil
		}

		return err
	}

	if showtime.ScreenRefID == nil {
		// Same reasoning: no screen to resolve seats against, so
		// there is nothing to release. Not an error.
		log.Printf(
			"MarkSeatsAsAvailable: showtime %s has no screen assigned — "+
				"skipping seat release, nothing to update",
			showtimeID,
		)

		return nil
	}

	// --------------------------------------------------------
	// Update Seat Status
	// --------------------------------------------------------

	result := r.DB.
		WithContext(ctx).
		Model(&models.Seat{}).
		Where(
			"screen_ref_id = ? AND seat_number IN ?",
			*showtime.ScreenRefID,
			seatNumbers,
		).
		Update(
			"status",
			"AVAILABLE",
		)

	if result.Error != nil {
		return result.Error
	}

	// NOTE: intentionally not erroring if RowsAffected is less than
	// len(seatNumbers) here (unlike MarkSeatsAsBooked) — some of
	// those seats may themselves have been deleted since, and that
	// alone shouldn't block cancelling the booking either.

	return nil
}
