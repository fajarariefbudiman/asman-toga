package models

import (
	"time"

	"github.com/google/uuid"
)

type ResetPassword struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Email     string    `gorm:"size:255;index" json:"email"`
	OTP       string    `gorm:"size:6" json:"otp"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time
}
