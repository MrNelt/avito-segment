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

func (s *Storage) GetUserByID(ID int) (models.User, error) {
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

func (s *Storage) CreateUser(ID int) error {
	db := s.db
	user := models.User{UserID: ID}
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteAddSegmentsToUser(deleteSegmentsNames, addSegmentsNames []string, ID int) error {
	if err := s.deleteSegmentsToUser(deleteSegmentsNames, ID); err != nil {
		return err
	}
	if err := s.addSegmentsToUser(addSegmentsNames, ID); err != nil {
		return err
	}
	return nil
}

func (s *Storage) addSegmentsToUser(addSegmentsNames []string, ID int) error {
	if len(addSegmentsNames) == 0 {
		return nil
	}
	var addSegments []models.Segment
	db := s.db
	err := db.Where("name in (?)", addSegmentsNames).Find(&addSegments).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return errorType.ErrSegmentNotFound
	case err != nil:
		return err
	}
	user := models.User{UserID: ID}
	for i := range addSegments {
		if err := db.Model(&addSegments[i]).Association("Users").Append(&user); err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) deleteSegmentsToUser(deleteSegmentsNames []string, ID int) error {
	if len(deleteSegmentsNames) == 0 {
		return nil
	}
	db := s.db
	var deleteSegments []models.Segment
	err := db.Where("name in (?)", deleteSegmentsNames).Find(&deleteSegments).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return errorType.ErrSegmentNotFound
	case err != nil:
		return err
	}
	user := models.User{UserID: ID}

	for i := range deleteSegments {
		if err := db.Model(&deleteSegments[i]).Association("Users").Delete(&user); err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) GetUserSegmentsByUserID(ID int) ([]models.Segment, error) {
	db := s.db
	var users []models.User
	if err := db.Preload("Segments").Where("user_id = ?", ID).Find(&users).Error; err != nil {
		return []models.Segment{}, err
	}
	var segments []models.Segment
	for _, user := range users {
		segments = append(segments, user.Segments...)
	}
	return segments, nil
}
