package repository

import (
	"context"
	"errors"

	// "strings"

	"cinema-ticket-booking/models"

	"gorm.io/gorm"
)

type ShowtimeRepository struct {
	DB *gorm.DB
}

func NewShowtimeRepository(
	db *gorm.DB,
) *ShowtimeRepository {

	return &ShowtimeRepository{
		DB: db,
	}
}

// Create resolves the incoming business UUIDs (MovieID, TheaterID,
// and the optional ScreenID) to their parents' numeric primary
// keys, then persists the Showtime using those numeric foreign
// keys. Resolution and insert happen inside a single transaction
// so a missing/invalid reference rolls back cleanly.
func (r *ShowtimeRepository) Create(
	ctx context.Context,
	showtime models.Showtime,
) error {

	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var movie models.Movie
		if err := tx.Where("movie_id = ?", showtime.MovieID).First(&movie).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("movie not found for movie_id: " + showtime.MovieID)
			}
			return err
		}
		showtime.MovieRefID = movie.ID

		var theater models.Theater
		if err := tx.Where("theater_id = ?", showtime.TheaterID).First(&theater).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("theater not found for theater_id: " + showtime.TheaterID)
			}
			return err
		}
		showtime.TheaterRefID = theater.ID

		if showtime.ScreenID != nil && *showtime.ScreenID != "" {
			var screen models.Screen
			if err := tx.Where("screen_id = ?", *showtime.ScreenID).First(&screen).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("screen not found for screen_id: " + *showtime.ScreenID)
				}
				return err
			}
			showtime.ScreenRefID = &screen.ID
		}

		return tx.Create(&showtime).Error
	})
}

// Get All Showtime
func (r *ShowtimeRepository) FindAll(
	ctx context.Context,
) ([]models.Showtime, error) {

	var showtimes []models.Showtime

	err := r.DB.WithContext(ctx).Find(&showtimes).Error
	if err != nil {
		return nil, err
	}

	return showtimes, nil
}

// Update Showtime
func (r *ShowtimeRepository) Update(
	ctx context.Context,
	showtimeID string,
	showtime models.Showtime,
) error {

	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var movie models.Movie
		if err := tx.Where("movie_id = ?", showtime.MovieID).First(&movie).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("movie not found")
			}

			return err
		}
		showtime.MovieRefID = movie.ID

		var theater models.Theater
		if err := tx.Where("theater_id = ?", showtime.TheaterID).First(&theater).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("theater not found")
			}

			return err
		}
		showtime.TheaterRefID = theater.ID

		// Screen เป็น Optional
		if showtime.ScreenID != nil && *showtime.ScreenID != "" {

			var screen models.Screen

			if err := tx.Where("screen_id = ?", *showtime.ScreenID).First(&screen).Error; err != nil {

				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("screen not found")
				}

				return err
			}

			showtime.ScreenRefID = &screen.ID

		} else {

			showtime.ScreenRefID = nil
			showtime.ScreenID = nil

		}
		result := tx.Model(&models.Showtime{}).
			Where("showtime_id = ?", showtimeID).
			Updates(map[string]interface{}{
				"movie_ref_id":   showtime.MovieRefID,
				"movie_id":       showtime.MovieID,
				"theater_ref_id": showtime.TheaterRefID,
				"theater_id":     showtime.TheaterID,
				"screen_ref_id":  showtime.ScreenRefID,
				"screen_id":      showtime.ScreenID,
				"show_date":      showtime.ShowDate,
				"start_time":     showtime.StartTime,
				"end_time":       showtime.EndTime,
				"price_standard": showtime.PriceStandard,
				"price_premium":  showtime.PricePremium,
				"price_vip":      showtime.PriceVIP,
				"status":         showtime.Status,
				"updated_at":     showtime.UpdatedAt,
			})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.New("showtime not found")
		}

		return nil
	})
}

// Errors ที่ controller ใช้แยกแยะเพื่อตอบ HTTP status
var (
	ErrShowtimeNotFound    = errors.New("showtime not found")
	ErrShowtimeHasBookings = errors.New("cannot delete showtime because it has bookings (kept for history)")
)

// Delete Showtime
//
//   - ไม่มี bookings ผูกอยู่  -> ลบแถวจริง (hard delete)
//   - มี bookings และยังไม่ cancelled -> เปลี่ยนเป็น "cancelled" (soft) เพื่อรักษาประวัติ
//   - มี bookings และ cancelled แล้ว -> คืน ErrShowtimeHasBookings (409)
func (r *ShowtimeRepository) Delete(
	ctx context.Context,
	showtimeID string,
) error {

	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var showtime models.Showtime
		if err := tx.Where("showtime_id = ?", showtimeID).First(&showtime).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrShowtimeNotFound
			}
			return err
		}

		var bookingCount int64
		if err := tx.Model(&models.Booking{}).
			Where("showtime_ref_id = ?", showtime.ID).
			Count(&bookingCount).Error; err != nil {
			return err
		}

		// ไม่มี bookings -> ลบจริง
		if bookingCount == 0 {
			return tx.Delete(&models.Showtime{}, showtime.ID).Error
		}

		// มี bookings อยู่แล้วและ cancelled แล้ว -> ลบไม่ได้
		if showtime.Status == "cancelled" {
			return ErrShowtimeHasBookings
		}

		// มี bookings -> soft cancel
		return tx.Model(&models.Showtime{}).
			Where("id = ?", showtime.ID).
			Updates(map[string]interface{}{
				"status":     "cancelled",
				"updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
			}).Error
	})
}
