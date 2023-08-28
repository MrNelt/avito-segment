package segment

import "segment/pkg/models"

type IStorage interface {
	CreateSegment(name string) error
	GetSegmentByName(name string) (models.Segment, error)
	DeleteSegment(name string) error
}
