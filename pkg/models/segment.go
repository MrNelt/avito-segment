package models

import "gorm.io/gorm"

type Segment struct {
	gorm.Model
	SegmentID uint `gorm:"primarykey"`
	Name      string
}
