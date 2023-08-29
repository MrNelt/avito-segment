package postgres

import (
	"fmt"
	"segment/pkg/models"
	pgSegment "segment/pkg/storage/postgres/segment"
	pgUser "segment/pkg/storage/postgres/user"
	"segment/pkg/storage/user"

	"segment/pkg/storage/segment"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	cfg Config
	db  *gorm.DB
}

func NewStorage(cfg Config) *Storage {
	return &Storage{cfg: cfg}
}

func (s *Storage) Connect() error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.cfg.Host, s.cfg.Port, s.cfg.User, s.cfg.Password, s.cfg.Database)

	database, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}

	s.db = database
	return nil
}

func (s *Storage) Close() error {
	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func (s *Storage) Init() *gorm.DB {
	return s.db
}

func (s *Storage) MakeMigrations() error {
	if err := s.db.AutoMigrate(&models.User{}, &models.Segment{}); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetSegmentStorage() segment.IStorage {
	return pgSegment.NewStorage(s.db)
}

func (s *Storage) GetUserStorage() user.IStorage {
	return pgUser.NewStorage(s.db)
}
