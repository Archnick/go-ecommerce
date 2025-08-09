package models

type ReviewPayload struct {
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment" binding:"omitempty"`
	UserID    uint   `json:"user_id" binding:"required"`
	ProductID uint   `json:"product_id" binding:"required"`
}

type UpdateReviewPayload struct {
	Rating  int    `json:"rating" binding:"omitempty,min=1,max=5"`
	Comment string `json:"comment" binding:"omitempty"`
}
