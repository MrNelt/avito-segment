package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	ID    string `gorm:"primary_key"`
	Users []User `gorm:"foreignKey:UserID"`
}
