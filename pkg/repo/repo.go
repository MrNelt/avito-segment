package repo

import (
	"segment/pkg/repo/segment"
	"segment/pkg/repo/user"
	"segment/pkg/storage"
)

type IRepository interface {
	GetSegmentRepository() segment.IRepository
	GetUserRepository() user.IRepository
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

func (r *Repository) GetUserRepository() user.IRepository {
	return user.NewRepository(r.storage.GetUserStorage())
}
