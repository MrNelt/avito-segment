package user

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

func (s *Storage) GetUserByID(ID uint) (models.User, error) {
	db := s.db
	var user models.User
	err := db.Where("ID = ?", ID).First(&user).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return models.User{}, errorType.ErrUserNotFound
	case err != nil:
		return models.User{}, err
	}
	return user, nil
}

func (s *Storage) CreateUser(ID uint) error {
	db := s.db
	user := models.User{ID: ID}
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteAddSegmentsToUser(deleteSegmentsNames, addSegmentsNames []string, ID uint) error {
	db := s.db
	var deleteSegments []models.Segment
	var addSegments []models.Segment
	err := db.Where("name in (?)", deleteSegmentsNames).Find(&deleteSegments).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return errorType.ErrSegmentNotFound
	case err != nil:
		return err
	}
	err = db.Where("name in (?)", addSegmentsNames).Find(&addSegments).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return errorType.ErrSegmentNotFound
	case err != nil:
		return err
	}
	user, err := s.GetUserByID(ID)
	switch {
	case errors.Is(err, errorType.ErrUserNotFound):
		s.CreateUser(ID)
		user, err = s.GetUserByID(ID)
		if err != nil {
			return nil
		}
	case err != nil:
		return err
	}
	if err := db.Model(&user).Association("Segments").Append(&addSegments); err != nil {
		return err
	}

	if err := db.Model(&user).Association("Segments").Delete(&deleteSegments); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUserSegmentsByUserID(ID uint) ([]models.Segment, error) {
	db := s.db
	user, err := s.GetUserByID(ID)
	switch {
	case errors.Is(err, errorType.ErrUserNotFound):
		s.CreateUser(ID)
		user, err = s.GetUserByID(ID)
		if err != nil {
			return []models.Segment{}, nil
		}
	case err != nil:
		return []models.Segment{}, err
	}
	var segments []models.Segment
	if err := db.Model(&user).Association("Segments").Find(&segments); err != nil {
		return []models.Segment{}, err
	}
	return segments, nil
}
