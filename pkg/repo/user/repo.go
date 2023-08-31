package user

import (
	"segment/pkg/dtos"
	"segment/pkg/storage/user"
	"time"
)

type IRepository interface {
	GetUserSegmentsByUserID(ID int) (dtos.UserDTO, error)
	DeleteAddSegmentsToUser(deleteSegmentsNames, addSegmentsNames []string, ID int, TTL uint, TTLUnit string) error
}

type Repository struct {
	storage user.IStorage
}

func NewRepository(storage user.IStorage) *Repository {
	return &Repository{storage: storage}
}

func (r *Repository) GetUserSegmentsByUserID(ID int) (dtos.UserDTO, error) {
	segments, err := r.storage.GetUserSegmentsByUserID(ID)
	if err != nil {
		return dtos.UserDTO{}, err
	}
	segmentNames := make(map[string]struct{})
	for i := range segments {
		segmentNames[segments[i].Name] = struct{}{}
	}
	var userDTO dtos.UserDTO
	userDTO.ID = ID
	userDTO.Segments = make([]string, len(segmentNames))
	i := 0
	for name := range segmentNames {
		userDTO.Segments[i] = name
		i++
	}
	return userDTO, nil
}

func (r *Repository) DeleteAddSegmentsToUser(deleteSegmentsNames, addSegmentsNames []string, ID int, TTL uint, TTLUnit string) error {
	date := time.Now().UTC()
	switch TTLUnit {
	case "SECONDS":
		date = date.Add(time.Second * time.Duration(TTL))
	case "MINUTES":
		date = date.Add(time.Minute * time.Duration(TTL))
	case "HOURS":
		date = date.Add(time.Hour * time.Duration(TTL))
	case "DAYS":
		date = date.Add(24 * time.Hour * time.Duration(TTL))
	}
	return r.storage.DeleteAddSegmentsToUser(deleteSegmentsNames, addSegmentsNames, ID, date)
}
