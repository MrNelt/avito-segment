package storage

import (
	"segment/pkg/storage/postgres"
	"segment/pkg/storage/segment"
	"segment/pkg/storage/user"

	"gorm.io/gorm"
)

type IStorage interface {
	Init() *gorm.DB
	Connect() error
	Close() error
	MakeMigrations() error

	GetSegmentStorage() segment.IStorage
	GetUserStorage() user.IStorage
}

var _ IStorage = &postgres.Storage{}
