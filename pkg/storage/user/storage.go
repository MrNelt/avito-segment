package user

import "segment/pkg/models"

type IStorage interface {
	CreateUser(ID uint) error
	GetUserSegmentsByUserID(ID uint) ([]models.Segment, error)
	DeleteAddSegmentsToUser(deleteSegmentsNames, addSegmentsNames []string, ID uint) error
}
