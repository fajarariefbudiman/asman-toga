package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPlant struct {
	ID         uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	UserID     uuid.UUID  `gorm:"type:char(36)" json:"user_id"`
	PlantID    int        `gorm:"not null" json:"plant_id"`
	Address    string     `gorm:"size:255" json:"address"`
	Latitude   float64    `json:"latitude"`
	Longitude  float64    `json:"longitude"`
	Notes      string     `json:"notes"`
	Status     string     `gorm:"type:enum('pending','approved','rejected');default:'pending'" json:"status"`
	ApprovedBy *uuid.UUID `gorm:"type:char(36)" json:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	User       User  `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Plant      Plant `gorm:"foreignKey:PlantID;references:ID" json:"plant"`
}

func (up *UserPlant) BeforeCreate(tx *gorm.DB) (err error) {
	up.ID = uuid.New()
	return
}
