package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Segments []Segment `gorm:"foreignKey:SegmentID"`
}
