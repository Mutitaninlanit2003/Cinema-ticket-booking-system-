package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Theater represents a physical cinema location.
//
// Parent table. No back-reference fields to Screen or Showtime —
// those children declare "belongs to Theater" themselves.
type Theater struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TheaterID string    `gorm:"type:uuid;uniqueIndex;not null" json:"theater_id"`
	Name      string    `gorm:"not null" json:"name" binding:"required"`
	Address   string    `json:"address" binding:"required"`
	Status    string    `gorm:"not null;default:active" json:"status" binding:"required,oneof=active inactive"`
	Location  string    `json:"location" binding:"required"`
	Capacity  int       `json:"capacity" binding:"required,gte=1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Theater) BeforeCreate(tx *gorm.DB) (err error) {
	if t.TheaterID == "" {
		t.TheaterID = uuid.NewString()
	}
	return
}
