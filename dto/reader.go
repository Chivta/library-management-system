package dto

type ReaderCreateDTO struct {
	Name    string `json:"name" validate:"required,min=1,max=100"`
	Surname string `json:"surname" validate:"required,min=1,max=100"`
}

type ReaderUpdateDTO struct {
	Name    string `json:"name" validate:"required,min=1,max=100"`
	Surname string `json:"surname" validate:"required,min=1,max=100"`
}

type ReaderResponseDTO struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
