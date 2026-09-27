package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Booking represents a user's reservation of one or more seats for
// a Showtime.
type Booking struct {
	ID        *uint   `gorm:"primaryKey" json:"id"`
	BookingID *string `gorm:"type:uuid;uniqueIndex;not null" json:"booking_id"`
	UserRefID *uint   `gorm:"not null;index" json:"-"`
	User      *User   `gorm:"foreignKey:UserRefID" json:"-" binding:"-"`

	UserID *string `gorm:"type:uuid;not null;index" json:"user_id" binding:"required,uuid"`

	ShowtimeRefID *uint     `gorm:"not null;index" json:"-"`
	Showtime      *Showtime `gorm:"foreignKey:ShowtimeRefID" json:"-" binding:"-"`

	ShowtimeID *string `gorm:"type:uuid;not null;index" json:"showtime_id" binding:"required,uuid"`

	SeatNumbers *pq.StringArray `gorm:"type:text[]" json:"seat_numbers" binding:"required,min=1" swaggertype:"array,string"`

	TotalPrice *float64 `json:"total_price" binding:"required,gt=0"`

	Status *string `gorm:"not null;default:pending" json:"status"`

	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// BeforeCreate generates a booking UUID when BookingID is empty.
func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.BookingID == nil || *b.BookingID == "" {
		id := uuid.NewString()
		b.BookingID = &id
	}

	return nil
}
