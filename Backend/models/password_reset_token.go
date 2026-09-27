package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PasswordResetToken struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Token string `gorm:"type:uuid;uniqueIndex;not null" json:"-"`

	UserID string `gorm:"type:uuid;not null;index" json:"-"`

	ExpiresAt time.Time `gorm:"not null" json:"-"`

	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate generates a reset token automatically
// if one was not provided.
func (p *PasswordResetToken) BeforeCreate(tx *gorm.DB) error {
	if p.Token == "" {
		p.Token = uuid.NewString()
	}

	return nil
}
