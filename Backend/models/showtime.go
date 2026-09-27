package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Showtime represents a single scheduled screening of a Movie in a
// Theater (optionally a specific Screen) at a given date/time.
//
// The actual foreign keys are the *RefID (uint) fields, pointing
// at the parent's numeric primary key. MovieID / TheaterID /
// ScreenID (uuid strings) are denormalized copies of the parents'
// business UUIDs, kept purely so the JSON request/response shape
// is unchanged from before this refactor. ScreenRefID/ScreenID are
// pointers because the Screen association is optional; a non-nil
// pointer is required so the zero value is SQL NULL rather than 0
// or "", which would otherwise violate (or wrongly satisfy) the FK
// constraint.
type Showtime struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	ShowtimeID string `gorm:"type:uuid;uniqueIndex;not null" json:"showtime_id"`

	MovieRefID uint   `gorm:"not null;index" json:"-"`
	Movie      Movie  `gorm:"foreignKey:MovieRefID" json:"-" binding:"-"`
	MovieID    string `gorm:"type:uuid;not null;index" json:"movie_id" binding:"required,uuid"`

	TheaterRefID uint    `gorm:"not null;index" json:"-"`
	Theater      Theater `gorm:"foreignKey:TheaterRefID" json:"-" binding:"-"`
	TheaterID    string  `gorm:"type:uuid;not null;index" json:"theater_id" binding:"required,uuid"`

	ScreenRefID *uint   `gorm:"index" json:"-"`
	Screen      Screen  `gorm:"foreignKey:ScreenRefID" json:"-" binding:"-"`
	ScreenID    *string `gorm:"type:uuid;index" json:"screen_id,omitempty" binding:"omitempty,uuid"`

	ShowDate  time.Time `json:"show_date" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`

	// Per-seat-type pricing for this showing, replacing the old single
	// flat Price field. Which price applies to a given seat is decided
	// by that seat's Seat.SeatType ("standard" / "premium" / "vip") —
	// see Seat in seat.go. Explicit `column:` tags are used because
	// GORM's default snake_case conversion of "PriceVIP" is not
	// guaranteed to land on "price_vip".
	PriceStandard float64 `gorm:"column:price_standard;not null" json:"price_standard" binding:"required,gt=0"`
	PricePremium  float64 `gorm:"column:price_premium;not null" json:"price_premium" binding:"required,gt=0"`
	PriceVIP      float64 `gorm:"column:price_vip;not null" json:"price_vip" binding:"required,gt=0"`

	Status string `gorm:"not null;default:scheduled" json:"status" binding:"required,oneof=scheduled showing finished cancelled"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Showtime) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ShowtimeID == "" {
		s.ShowtimeID = uuid.NewString()
	}
	return
}
