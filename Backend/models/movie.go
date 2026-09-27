package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Movie represents a film that can be scheduled into Showtimes.
//
// Movie is a parent table. It intentionally does NOT declare a
// back-reference field to Showtime (e.g. "Showtimes []Showtime").
// Only the child (Showtime) declares the association, via a plain
// "belongs to" foreign key. Declaring the relationship from both
// sides with matching tag names is what caused GORM to get
// confused about which table owns the foreign key and try to
// create the constraint on "movies" before "showtimes" existed.
type Movie struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	MovieID     string `gorm:"type:uuid;uniqueIndex;not null" json:"movie_id"`
	Title       string `gorm:"not null;index" json:"title" binding:"required"`
	Description string `json:"description"`

	Genre       pq.StringArray `gorm:"type:text[]" json:"genre" binding:"required" swaggertype:"array,string"`
	Duration    int            `json:"duration" binding:"required,gte=1"`
	Language    string         `json:"language" binding:"required"`
	Subtitle    string         `json:"subtitle"`
	Rating      string         `json:"rating" binding:"required"`
	PosterURL   string         `json:"poster_url"`
	TrailerURL  string         `json:"trailer_url"`
	ReleaseDate time.Time      `json:"release_date" binding:"required"`
	Status      string         `gorm:"not null;default:active" json:"status" binding:"required,oneof=NOW_SHOWING COMING_SOON ENDED"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Movie) BeforeCreate(tx *gorm.DB) (err error) {
	if m.MovieID == "" {
		m.MovieID = uuid.NewString()
	}
	return
}
