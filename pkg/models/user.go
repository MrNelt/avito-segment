package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	ID       string    `gorm:"primary_key"`
	Segments []Segment `gorm:"foreignKey:SegmentID"`
}
