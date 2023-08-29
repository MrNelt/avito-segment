package storage

import (
	"segment/pkg/storage/postgres"
	"segment/pkg/storage/segment"
	"segment/pkg/storage/user"
)

type IStorage interface {
	Connect() error
	Close() error
	MakeMigrations() error

	GetSegmentStorage() segment.IStorage
	GetUserStorage() user.IStorage
}

var _ IStorage = &postgres.Storage{}
