package user

import (
	"segment/pkg/models"
	"time"
)

type IStorage interface {
	CreateUser(ID int) error
	GetUserSegmentsByUserID(ID int) ([]models.Segment, error)
	DeleteAddSegmentsToUser(deleteSegmentsNames, addSegmentsNames []string, ID int, expirationDate time.Time) error
}
