package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	BanjarID  int       `gorm:"not null" json:"banjar_id"`
	Name      string    `gorm:"size:100" json:"name"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `gorm:"size:255" json:"-"`
	Role      string    `gorm:"type:enum('admin','user');default:'user'" json:"role"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Banjar    Banjar `gorm:"foreignKey:BanjarID;references:ID" json:"banjar"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
