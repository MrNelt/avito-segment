package segment

import (
	"errors"
	errorType "segment/pkg/errors"
	"segment/pkg/models"

	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateSegment(name string) error {
	db := s.db
	if err := db.Create(&models.Segment{Name: name}).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetSegmentByName(name string) (models.Segment, error) {
	db := s.db
	var segment models.Segment
	err := db.Where("name = ?", name).First(&segment).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return models.Segment{}, errorType.ErrSegmentNotFound
	case err != nil:
		return models.Segment{}, err
	}
	return segment, nil
}

func (s *Storage) DeleteSegmentByName(name string) error {
	db := s.db
	err := db.Where("name = ?", name).Delete(&models.Segment{}).Error
	if err != nil {
		return err
	}
	return nil
}
