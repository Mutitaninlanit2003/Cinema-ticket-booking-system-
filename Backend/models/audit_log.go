package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditLog records an auditable action taken by a user against a
// resource. UserID intentionally has no hard foreign-key
// constraint so that log history survives even if the referenced
// user is later removed.
type AuditLog struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	AuditLogID string `gorm:"type:uuid;uniqueIndex;not null" json:"audit_log_id"`

	UserID string `gorm:"type:uuid;index" json:"user_id" binding:"required,uuid"`

	Action     string `gorm:"not null" json:"action" binding:"required"`
	Resource   string `json:"resource" binding:"required"`
	ResourceID string `json:"resource_id" binding:"required"`
	Details    string `json:"details" binding:"required"`

	CreatedAt time.Time `json:"created_at"`
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) (err error) {
	if a.AuditLogID == "" {
		a.AuditLogID = uuid.NewString()
	}
	return
}
