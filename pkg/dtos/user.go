package dtos

type UserDTO struct {
	ID       int      `json:"ID"`
	Segments []string `json:"segments"`
}
