package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Payment records a payment transaction against a Booking.
//
// The actual foreign key is BookingRefID (uint), pointing at
// Booking.ID. BookingID (uuid string) is a denormalized copy of
// the parent Booking's business UUID, kept for API compatibility.
type Payment struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	PaymentID string `gorm:"type:uuid;uniqueIndex;not null" json:"payment_id"`

	BookingRefID uint    `gorm:"not null;index" json:"-"`
	Booking      Booking `gorm:"foreignKey:BookingRefID" json:"-" binding:"-"`
	BookingID    string  `gorm:"type:uuid;not null;index" json:"booking_id" binding:"required,uuid"`

	Amount float64 `gorm:"not null" json:"amount" binding:"required,gt=0"`

	Method string `json:"method" binding:"required,oneof=cash credit_card promptpay transfer"`

	Status string `gorm:"not null;default:pending" json:"status" binding:"required,oneof=pending paid failed refunded"`

	TransactionRef string `json:"transaction_ref" binding:"omitempty,alphanum,min=6,max=20"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if p.PaymentID == "" {
		p.PaymentID = uuid.NewString()
	}
	return
}
