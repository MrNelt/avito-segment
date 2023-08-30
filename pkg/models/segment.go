package models

import "gorm.io/gorm"

type Segment struct {
	gorm.Model

	ID    int `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Name  string
	Users []User `gorm:"many2many:segment_users"`
}
