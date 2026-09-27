package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Seat represents a single physical seat inside a Screen. A Seat
// is unique per Screen by SeatNumber.
//
// The actual foreign key is ScreenRefID (uint), pointing at
// Screen.ID. ScreenID (uuid string) is a denormalized copy of the
// parent Screen's business UUID, kept for API compatibility only.
type Seat struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	SeatID string `gorm:"type:uuid;uniqueIndex;not null" json:"seat_id"`

	ScreenRefID uint   `gorm:"not null;index:idx_screen_seat_number,unique,priority:1" json:"-"`
	Screen      Screen `gorm:"foreignKey:ScreenRefID" json:"-" binding:"-"`
	ScreenID    string `gorm:"type:uuid;not null;index" json:"screen_id" binding:"required,uuid"`

	SeatNumber string    `gorm:"not null;index:idx_screen_seat_number,unique,priority:2" json:"seat_number" binding:"required"`
	SeatType   string    `gorm:"not null;default:standard" json:"seat_type" binding:"required,oneof=standard premium vip"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Status     string    `gorm:"not null;default:AVAILABLE" json:"status"`
}

func (s *Seat) BeforeCreate(tx *gorm.DB) (err error) {
	if s.SeatID == "" {
		s.SeatID = uuid.NewString()
	}

	return
}
