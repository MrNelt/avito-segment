package dtos

type UserResponseDTO struct {
	ID       uint     `json:"ID"`
	Segments []string `json:"segments"`
}

type UserRequestDTO struct {
	ID uint `json:"ID"`
}
