package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	UserID         int `gorm:"primary_key"`
	ExpirationDate time.Time
	Segments       []Segment `gorm:"many2many:segment_users"`
}
