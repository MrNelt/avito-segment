package user

import "segment/pkg/models"

type IStorage interface {
	CreateUser(ID int) error
	GetUserSegmentsByUserID(ID int) ([]models.Segment, error)
	DeleteAddSegmentsToUser(deleteSegmentsNames, addSegmentsNames []string, ID int) error
}
