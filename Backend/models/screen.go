package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Screen represents an individual auditorium inside a Theater.
//
// The actual foreign key is TheaterRefID (uint), pointing at
// Theater.ID — GORM's default belongs-to convention (a numeric FK
// referencing the parent's primary key). TheaterID (uuid string)
// is kept purely for API compatibility: it is a denormalized copy
// of the parent Theater's business UUID, written once at create
// time and never used as the relational key.
type Screen struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	ScreenID string `gorm:"type:uuid;uniqueIndex;not null" json:"screen_id"`

	TheaterRefID uint    `gorm:"not null;index" json:"-"`
	Theater      Theater `gorm:"foreignKey:TheaterRefID" json:"-" binding:"-"`
	TheaterID    string  `gorm:"type:uuid;not null;index" json:"theater_id" binding:"required,uuid"`

	Name string `gorm:"not null" json:"name" binding:"required"`

	TotalSeats int `json:"total_seats" binding:"required,gte=1"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Screen) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ScreenID == "" {
		s.ScreenID = uuid.NewString()
	}
	return
}
