package storage

import (
	"segment/pkg/storage/postgres"

	"gorm.io/gorm"
)

type IStorage interface {
	Connect() error
	Close() error
	Init() *gorm.DB
	MakeMigrations() error
}

var _ IStorage = &postgres.Storage{}
