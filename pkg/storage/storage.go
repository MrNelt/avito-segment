package storage

import (
	"segment/pkg/storage/postgres"
	"segment/pkg/storage/segment"

	"gorm.io/gorm"
)

type IStorage interface {
	Connect() error
	Close() error
	Init() *gorm.DB
	MakeMigrations() error

	GetSegmentStorage() segment.IStorage
}

var _ IStorage = &postgres.Storage{}
