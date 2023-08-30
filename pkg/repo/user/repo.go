package user

import (
	"segment/pkg/dtos"
	"segment/pkg/storage/user"
)

type IRepository interface {
	GetUserSegmentsByUserID(ID int) (dtos.UserDTO, error)
	DeleteAddSegmentsToUser(deleteSegmentsNames, addSegmentsNames []string, ID int) error
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

func (r *Repository) DeleteAddSegmentsToUser(deleteSegmentsNames, addSegmentsNames []string, ID int) error {
	return r.storage.DeleteAddSegmentsToUser(deleteSegmentsNames, addSegmentsNames, ID)
}
