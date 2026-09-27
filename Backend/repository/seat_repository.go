package repository

import (
	"errors"
	"time"

	"cinema-ticket-booking/config"
	"cinema-ticket-booking/models"

	"gorm.io/gorm"
)

type SeatRepository struct {
	DB *gorm.DB
}

func NewSeatRepository() *SeatRepository {
	return &SeatRepository{
		DB: config.DB,
	}
}

func (r *SeatRepository) CreateSeat(
	seat models.Seat,
) error {

	return r.DB.Create(&seat).Error
}

func (r *SeatRepository) GetSeats() ([]models.Seat, error) {

	var seats []models.Seat

	err := r.DB.
		Preload("Screen").
		Find(&seats).Error

	return seats, err
}

func (r *SeatRepository) UpdateSeat(
	seatID string,
	seatNumber string,
	seatType string,
	updatedAt time.Time,
) error {

	result := r.DB.
		Model(&models.Seat{}).
		Where("seat_id = ?", seatID).
		Updates(map[string]interface{}{
			"seat_number": seatNumber,
			"seat_type":   seatType,
			"updated_at":  updatedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("seat not found")
	}

	return nil
}

func (r *SeatRepository) DeleteSeat(
	seatID string,
) error {

	result := r.DB.
		Where("seat_id = ?", seatID).
		Delete(&models.Seat{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("seat not found")
	}

	return nil
}

func (r *SeatRepository) GetScreenByUUID(
	screenID string,
) (*models.Screen, error) {

	var screen models.Screen

	err := r.DB.
		Where("screen_id = ?", screenID).
		First(&screen).Error

	if err != nil {
		return nil, err
	}

	return &screen, nil
}

// Create multiple seats in one transaction
func (r *SeatRepository) CreateSeatsBulk(
	seats []*models.Seat,
) error {

	if len(seats) == 0 {
		return errors.New("no seats to create")
	}

	return r.DB.Transaction(func(tx *gorm.DB) error {

		// Check whether any seat already exists
		for _, seat := range seats {

			var count int64

			err := tx.
				Model(&models.Seat{}).
				Where(
					"screen_ref_id = ? AND seat_number = ?",
					seat.ScreenRefID,
					seat.SeatNumber,
				).
				Count(&count).Error

			if err != nil {
				return err
			}

			if count > 0 {
				return errors.New(
					"seat " + seat.SeatNumber + " already exists",
				)
			}
		}

		// Create all seats
		if err := tx.Create(seats).Error; err != nil {
			return err
		}

		return nil
	})
}

// Delete all seats for a specific screen
func (r *SeatRepository) DeleteSeatsByScreenID(screenID string) error {
	var screen models.Screen

	// Find screen by UUID
	if err := r.DB.
		Where("screen_id = ?", screenID).
		First(&screen).Error; err != nil {
		return err
	}

	// Delete all seats belonging to this screen
	result := r.DB.
		Where("screen_ref_id = ?", screen.ID).
		Delete(&models.Seat{})

	return result.Error
}
