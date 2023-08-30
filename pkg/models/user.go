package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	UserID   int       `gorm:"primary_key"`
	Segments []Segment `gorm:"many2many:segment_users"`
}
