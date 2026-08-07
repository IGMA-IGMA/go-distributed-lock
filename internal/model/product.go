package model

import "time"

type Product struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:255;not null"`
	Quantity  int       `gorm:"not null;default:0"`
	Version   int       `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}