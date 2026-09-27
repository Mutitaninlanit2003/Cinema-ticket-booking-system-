package repository

import (
	"cinema-ticket-booking/config"
	"cinema-ticket-booking/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

// Errors ที่ handler ใช้แยกแยะเพื่อตอบ HTTP status ให้ถูกต้อง
var (
	ErrScreenInUse    = errors.New("cannot delete screen because it is used by active showtimes (cancel them first)")
	ErrScreenNotFound = errors.New("screen not found")
)

type ScreenRepository struct {
	DB *gorm.DB
}

func NewScreenRepository() *ScreenRepository {
	return &ScreenRepository{
		DB: config.DB,
	}
}

func (r *ScreenRepository) CreateScreen(screen models.Screen) error {
	return r.DB.Create(&screen).Error
}

func (r *ScreenRepository) GetScreens() ([]models.Screen, error) {
	var screens []models.Screen

	err := r.DB.
		Preload("Theater").
		Find(&screens).Error

	return screens, err
}

func (r *ScreenRepository) UpdateScreen(
	screenID string,
	screen models.Screen,
) error {
	result := r.DB.
		Model(&models.Screen{}).
		Where("screen_id = ?", screenID).
		Updates(map[string]interface{}{
			"name":        screen.Name,
			"total_seats": screen.TotalSeats,
			"updated_at":  time.Now(), // แก้: เดิมใช้ screen.UpdatedAt ซึ่งอาจเป็นเวลาศูนย์
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrScreenNotFound
	}

	return nil
}

// DeleteScreen ลบ screen แบบปลอดภัย:
//   - ถ้ามี showtimes ใช้อยู่ -> คืน ErrScreenInUse (ต้องลบ showtimes ก่อน)
//   - ถ้าไม่มี -> ลบ seats ของ screen แล้วลบ screen (ใน transaction เดียว)
func (r *ScreenRepository) DeleteScreen(screenID string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {

		// 0) หา screen จาก UUID เพื่อเอา numeric ID (FK จริง)
		var screen models.Screen
		if err := tx.Where("screen_id = ?", screenID).First(&screen).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrScreenNotFound
			}
			return err
		}

		// 1) ห้ามลบถ้ายังมี showtimes ที่ active (ไม่ใช่ cancelled) ใช้อยู่
		//    หมายเหตุ: Delete ของ showtime เป็น soft delete (ตั้ง status = cancelled)
		var activeCount int64
		if err := tx.Model(&models.Showtime{}).
			Where("screen_ref_id = ? AND status <> ?", screen.ID, "cancelled").
			Count(&activeCount).Error; err != nil {
			return err
		}
		if activeCount > 0 {
			return ErrScreenInUse
		}

		// 2) showtimes ที่ cancelled แล้ว -> ตัดการเชื่อมกับ screen (screen เป็น optional)
		//    เก็บแถว showtime ไว้เพื่อรักษาประวัติ bookings/payments
		if err := tx.Model(&models.Showtime{}).
			Where("screen_ref_id = ?", screen.ID).
			Updates(map[string]interface{}{
				"screen_ref_id": nil,
				"screen_id":     nil,
			}).Error; err != nil {
			return err
		}

		// 3) ลบ booking_seats ที่อ้างอิงที่นั่งของ screen นี้ (ไม่งั้นชน fk_booking_seats_seat)
		//    หมายเหตุ: bookings ยังเก็บเลขที่นั่งไว้ใน seat_numbers จึงไม่สูญเสียประวัติทั้งหมด
		seatIDs := tx.Model(&models.Seat{}).
			Select("id").
			Where("screen_ref_id = ?", screen.ID)

		if err := tx.Where("seat_ref_id IN (?)", seatIDs).
			Delete(&models.BookingSeat{}).Error; err != nil {
			return err
		}

		// 4) ลบ seats ของ screen นี้
		if err := tx.
			Where("screen_ref_id = ?", screen.ID).
			Delete(&models.Seat{}).Error; err != nil {
			return err
		}

		// 5) ลบ screen
		return tx.Delete(&models.Screen{}, screen.ID).Error
	})
}

// ForceDeleteScreen ลบ screen พร้อมข้อมูลลูกทั้งหมด
// (booking_seats -> payments -> bookings -> showtimes -> seats -> screen)
//
// ⚠️ ใช้กับข้อมูลทดสอบ/dev เท่านั้น เพราะจะลบประวัติการจองและการชำระเงินด้วย
// ควรจำกัดสิทธิ์เฉพาะ admin และบันทึก Audit Log ที่ handler
func (r *ScreenRepository) ForceDeleteScreen(screenID string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {

		var screen models.Screen
		if err := tx.Where("screen_id = ?", screenID).First(&screen).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrScreenNotFound
			}
			return err
		}

		// showtime ทั้งหมดของ screen นี้
		showtimeIDs := tx.Model(&models.Showtime{}).
			Select("id").
			Where("screen_ref_id = ?", screen.ID)

		// booking ทั้งหมดของ showtime เหล่านั้น
		bookingIDs := tx.Model(&models.Booking{}).
			Select("id").
			Where("showtime_ref_id IN (?)", showtimeIDs)

		// seat ทั้งหมดของ screen นี้
		seatIDs := tx.Model(&models.Seat{}).
			Select("id").
			Where("screen_ref_id = ?", screen.ID)

		// 0) booking_seats (อ้างอิงทั้ง bookings และ seats)
		if err := tx.Where("booking_ref_id IN (?) OR seat_ref_id IN (?)", bookingIDs, seatIDs).
			Delete(&models.BookingSeat{}).Error; err != nil {
			return err
		}

		// 1) payments
		if err := tx.Where("booking_ref_id IN (?)", bookingIDs).
			Delete(&models.Payment{}).Error; err != nil {
			return err
		}

		// 2) bookings
		if err := tx.Where("showtime_ref_id IN (?)", showtimeIDs).
			Delete(&models.Booking{}).Error; err != nil {
			return err
		}

		// 3) showtimes
		if err := tx.Where("screen_ref_id = ?", screen.ID).
			Delete(&models.Showtime{}).Error; err != nil {
			return err
		}

		// 4) seats
		if err := tx.Where("screen_ref_id = ?", screen.ID).
			Delete(&models.Seat{}).Error; err != nil {
			return err
		}

		// 5) screen
		return tx.Delete(&models.Screen{}, screen.ID).Error
	})
}

func (r *ScreenRepository) GetTheaterByUUID(
	theaterID string,
) (*models.Theater, error) {
	var theater models.Theater

	err := r.DB.
		Where("theater_id = ?", theaterID).
		First(&theater).Error

	if err != nil {
		return nil, err
	}

	return &theater, nil
}
