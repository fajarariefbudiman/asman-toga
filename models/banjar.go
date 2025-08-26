package models

import "time"

type Banjar struct {
	ID        int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"size:100" json:"name"`
	Slug      string `gorm:"type:text" json:"slug"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
