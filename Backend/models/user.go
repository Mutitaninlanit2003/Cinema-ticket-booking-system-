package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID string `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`

	Name string `gorm:"not null" json:"name" binding:"required"`

	Email string `gorm:"uniqueIndex;not null" json:"email" binding:"required,email"`

	// PasswordHash is used for normal email/password login.
	// Google OAuth users may not have a password.
	PasswordHash *string `json:"-"`

	Role string `gorm:"not null;default:user" json:"role"`

	ProfileImage *string `json:"profile_image,omitempty"`

	// Google OAuth
	//
	// GoogleID is the unique identifier returned by Google.
	// It is nullable because existing email/password users
	// do not necessarily have a Google account linked.
	GoogleID *string `gorm:"uniqueIndex" json:"-"`

	// Forgot Password
	ResetToken       *string    `json:"-"`
	ResetTokenExpiry *time.Time `json:"-"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UserID == "" {
		u.UserID = uuid.NewString()
	}

	if u.Role == "" {
		u.Role = "user"
	}

	return nil
}
