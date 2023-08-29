package repo

import (
	"segment/pkg/repo/segment"
	"segment/pkg/storage"
)

type IRepository interface {
	GetSegmentRepository() segment.IRepository
}

type Repository struct {
	storage storage.IStorage
}

func NewRepository(storage storage.IStorage) *Repository {
	return &Repository{storage: storage}
}

func (r *Repository) GetSegmentRepository() segment.IRepository {
	return segment.NewRepository(r.storage.GetSegmentStorage())
}
