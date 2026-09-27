package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BookingSeat is the join table that normalizes the many-to-many
// relationship between Bookings and Seats.
//
// The actual foreign keys are BookingRefID and SeatRefID (uint),
// pointing at the parents' numeric primary keys. BookingID and
// SeatID (uuid strings) are denormalized copies of the parents'
// business UUIDs, kept for API compatibility / easy lookups.
type BookingSeat struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	BookingSeatID string `gorm:"type:uuid;uniqueIndex;not null" json:"booking_seat_id"`

	BookingRefID uint    `gorm:"not null;index:idx_booking_seat,unique,priority:1" json:"-"`
	Booking      Booking `gorm:"foreignKey:BookingRefID" json:"-"`
	BookingID    string  `gorm:"type:uuid;not null;index" json:"booking_id"`

	SeatRefID uint   `gorm:"not null;index:idx_booking_seat,unique,priority:2" json:"-"`
	Seat      Seat   `gorm:"foreignKey:SeatRefID" json:"-"`
	SeatID    string `gorm:"type:uuid;not null;index" json:"seat_id"`

	CreatedAt time.Time `json:"created_at"`
}

func (bs *BookingSeat) BeforeCreate(tx *gorm.DB) (err error) {
	if bs.BookingSeatID == "" {
		bs.BookingSeatID = uuid.NewString()
	}
	return
}
