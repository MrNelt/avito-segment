package user

import "segment/pkg/models"

type IStorage interface {
	CreateUser(ID string) error
	GetUserSegmentsByUserID(ID string) ([]models.Segment, error)
}
