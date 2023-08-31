package user

import (
	"errors"
	errorType "segment/pkg/errors"
	"segment/pkg/models"
	"time"

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

func (s *Storage) DeleteAddSegmentsToUser(deleteSegmentsNames, addSegmentsNames []string, ID int, expirationTime time.Time) error {
	if err := s.deleteSegmentsToUser(deleteSegmentsNames, ID); err != nil {
		return err
	}
	if err := s.addSegmentsToUser(addSegmentsNames, ID, expirationTime); err != nil {
		return err
	}
	return nil
}

func (s *Storage) addSegmentsToUser(addSegmentsNames []string, ID int, expirationDate time.Time) error {
	if len(addSegmentsNames) == 0 {
		return nil
	}
	var addSegments []models.Segment
	db := s.db
	req := db.Where("name in (?)", addSegmentsNames).Find(&addSegments)
	if req.RowsAffected != int64(len(addSegmentsNames)) {
		return errorType.ErrSegmentNotFound
	}
	switch {
	case errors.Is(req.Error, gorm.ErrRecordNotFound):
		return errorType.ErrSegmentNotFound
	case req.Error != nil:
		return req.Error
	}
	user := models.User{UserID: ID, ExpirationDate: expirationDate}
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

	deleteSegmentsNameSet := make(map[string]struct{})
	for _, name := range deleteSegmentsNames {
		deleteSegmentsNameSet[name] = struct{}{}
	}
	var users []models.User
	if err := db.Preload("Segments").Where("user_id = ?", ID).Find(&users).Error; err != nil {
		return err
	}
	for _, user := range users {
		for _, segment := range user.Segments {
			if _, ok := deleteSegmentsNameSet[segment.Name]; ok {
				if err := db.Model(&user).Association("Segments").Delete(&segment); err != nil {
					return err
				}
			}
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
		if !user.ExpirationDate.Before(time.Now()) {
			segments = append(segments, user.Segments...)
		} else {
			for _, segment := range user.Segments {
				if err := db.Model(&user).Association("Segments").Delete(&segment); err != nil {
					return []models.Segment{}, err
				}
			}
		}
	}
	return segments, nil
}
