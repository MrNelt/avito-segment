package models

import "gorm.io/gorm"

type Segment struct {
	gorm.Model

	ID    string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Users []User `gorm:"foreignKey:UserID"`
}
