package models

import (
	"time"
)

type Plant struct {
	ID        int    `gorm:"primaryKey;autoIncrement" json:"id"`
	PlantName string `gorm:"size:100" json:"plant_name"`
	Slug      string `gorm:"type:text" json:"slug"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
