package segment

import "segment/pkg/storage/segment"

type IRepository interface {
	CreateSegmentByName(name string) error
	DeleteSegmentByName(name string) error
}

type Repository struct {
	storage segment.IStorage
}

func NewRepository(storage segment.IStorage) *Repository {
	return &Repository{storage: storage}
}

func (r *Repository) CreateSegmentByName(name string) error {
	return r.storage.CreateSegment(name)
}

func (r *Repository) DeleteSegmentByName(name string) error {
	return r.storage.DeleteSegmentByName(name)
}

var _ IRepository = &Repository{}
