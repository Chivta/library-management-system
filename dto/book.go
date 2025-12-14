package dto

type BookCreateDTO struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

type BookUpdateDTO struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

type BookResponseDTO struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
}
