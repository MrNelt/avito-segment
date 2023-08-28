package segment

import (
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
	if err := db.Where("name = ?", name).First(&segment).Error; err != nil {
		return models.Segment{}, nil
	}
	return segment, nil
}

func (s *Storage) DeleteSegmentByName(name string) error {
	db := s.db
	if err := db.Where("name = ?", name).Delete(&models.Segment{}).Error; err != nil {
		return err
	}
	return nil
}
