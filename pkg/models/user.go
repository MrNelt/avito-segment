package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint      `gorm:"primaryKey"`
	Segments []Segment `gorm:"foreignKey:SegmentID"`
}
