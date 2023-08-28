package user

import "segment/pkg/models"

type IStorage interface {
	GetUserSegmentsByUserID(ID string) []models.Segment
}
