package models

type CategoryPayload struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" binding:"required,min=2"`
	Description string `json:"description" binding:"omitempty"`
}

type UpdateCategoryPayload struct {
	Name        string `json:"name" binding:"omitempty,min=2"`
	Description string `json:"description" binding:"omitempty"`
}
